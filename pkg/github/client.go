package github

import (
	"context"
	"fmt"

	gh "github.com/google/go-github/v59/github"
)

type Client struct {
	client *gh.Client
}

func NewClient(token string) *Client {
	return &Client{
		gh.NewClient(nil).WithAuthToken(token),
	}
}

func (c *Client) CurrentUser(ctx context.Context) (string, error) {
	user, _, err := c.client.Users.Get(ctx, "")
	if err != nil {
		return "", fmt.Errorf("get current user: %w", err)
	}
	return user.GetLogin(), nil
}

func (c *Client) GetPullRequest(ctx context.Context, owner, repo string, number uint32) (*PullRequest, error) {
	pr, _, err := c.client.PullRequests.Get(ctx, owner, repo, int(number))
	if err != nil {
		return nil, fmt.Errorf("get pull request: %w", err)
	}

	ret := &PullRequest{
		Number: pr.GetNumber(),
		Title:  pr.GetTitle(),
		Body:   pr.GetBody(),
		Draft:  pr.GetDraft(),
		Author: pr.User.GetLogin(),
		State:  pr.GetState(),
	}

	if ret.Reviewers, err = c.getPrReviewers(ctx, owner, repo, number); err != nil {
		return nil, fmt.Errorf("get pull request reviewers: %w", err)
	}

	if ret.Files, err = c.getPrFiles(ctx, owner, repo, number); err != nil {
		return nil, fmt.Errorf("get pull request files: %w", err)
	}
	return ret, nil
}

func (c *Client) SubmitReview(ctx context.Context, req SubmitReviewRequest) error {
	if !req.Event.IsValid() {
		return fmt.Errorf("invalid event: %s", req.Event)
	}
	event := string(req.Event)
	_, _, err := c.client.PullRequests.CreateReview(ctx, req.Owner, req.Repo, req.Number, &gh.PullRequestReviewRequest{
		Body:  &req.Body,
		Event: &event,
	})
	if err != nil {
		return fmt.Errorf("create review: %w", err)
	}
	return nil
}

// RequestReview requests a review for the specified pull request. If no reviewer is specified, the current user is used.
func (c *Client) RequestReview(ctx context.Context, owner, repo string, number uint32, reviewer ...string) error {
	if len(reviewer) == 0 { // use the current user
		user, err := c.CurrentUser(ctx)
		if err != nil {
			return fmt.Errorf("get current user: %w", err)
		}
		reviewer = append(reviewer, user)
	}
	_, _, err := c.client.PullRequests.RequestReviewers(ctx, owner, repo, int(number), gh.ReviewersRequest{
		Reviewers: reviewer,
	})
	if err != nil {
		return fmt.Errorf("request reviewers: %w", err)
	}
	return nil
}

func (c *Client) getPrReviewers(ctx context.Context, owner, repo string, number uint32) ([]string, error) {
	reviews, _, err := c.client.PullRequests.ListReviews(ctx, owner, repo, int(number), nil)
	if err != nil {
		return nil, fmt.Errorf("list reviews: %w", err)
	}

	var reviewers []string
	for _, review := range reviews {
		reviewers = append(reviewers, review.GetUser().GetLogin())
	}
	return reviewers, nil
}

func (c *Client) getPrFiles(ctx context.Context, owner, repo string, number uint32) ([]*File, error) {
	files, _, err := c.client.PullRequests.ListFiles(ctx, owner, repo, int(number), nil)
	if err != nil {
		return nil, fmt.Errorf("list files: %w", err)
	}

	var ret []*File
	for _, file := range files {
		if file.GetStatus() == "" || file.GetFilename() == "" {
			continue
		}
		ret = append(ret, &File{
			Name:   file.GetFilename(),
			Status: file.GetStatus(),
		})
	}
	return ret, nil
}

package cmd

import (
	"os"
	"strings"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/vietanhduong/auto-approve/pkg/aafile"
	"github.com/vietanhduong/auto-approve/pkg/github"
	"github.com/vietanhduong/auto-approve/pkg/logging"
)

func NewCommand() *cobra.Command {
	var ghToken string
	var aaFilePath string
	var githubLogFormat bool
	var prNumber int
	var githubRepo string
	var version bool
	var commentBody string

	cmd := &cobra.Command{
		Use:   "auto-approve",
		Short: "Auto-Approve the input Pull Request",
		Long: `Auto-Approve the input Pull Request
to bypass the "Require at least 1 reviewer" of GitHub.`,
		Example: `$ cat "$(git rev-parse --show-toplevel)/AUTOAPPROVE"
# Use similar GitHub CODEOWNERS syntax
* @vietanhduong

$ export GH_TOKEN=$GITHUB_TOKEN
$ export GITHUB_REPOSITORY="<owner>/<repo>"
$ auto-approve --pr $PR_NUMBER --comment "LGTM!"
`,
		Run: func(cmd *cobra.Command, args []string) {
			if version {
				printVersion()
				os.Exit(0)
			}
			if githubLogFormat {
				logging.EnableGitHubFormat()
			}
			if ghToken == "" { // use the default GITHUB_TOKEN if set
				ghToken = os.Getenv("GITHUB_TOKEN")
			}
			if ghToken == "" {
				logging.Errorf("GitHub token is required.")
				os.Exit(1)
			}
			if prNumber == -1 {
				logging.Errorf("Pull Request number is required.")
				os.Exit(1)
			}
			if githubRepo == "" {
				logging.Errorf("GitHub repository is required.")
				os.Exit(1)
			}

			owner, repo := parseRepo(githubRepo)

			if aaFilePath == "" {
				repoRoot, err := github.RepoRoot()
				if err != nil {
					logging.Errorf("Failed to determine repository root: %v", err)
					os.Exit(1)
				}
				aaFilePath = discoveryAaFile(repoRoot)
			}
			if aaFilePath == "" {
				logging.Notice("No AUTOAPPROVE file found. Skip auto-approving.")
				os.Exit(0)
			}
			gh := github.NewClient(ghToken)
			user, err := gh.CurrentUser(cmd.Context())
			if err != nil {
				logging.Errorf("Failed to get the Current GitHub User: %v", err)
				os.Exit(1)
			}

			pr, err := gh.GetPullRequest(cmd.Context(), owner, repo, uint32(prNumber))
			if err != nil {
				logging.Errorf("Failed to get Pull Request: %v", err)
				os.Exit(1)
			}

			if pr.State != "open" || pr.Draft {
				logging.Notice("The Pull Request is not in the open state, or it's a draft PR.")
				os.Exit(0)
			}

			aaFile, err := parseAAFile(aaFilePath)
			if err != nil {
				logging.Errorf("parse AUTOAPPROVE file: %v", err)
				os.Exit(1)
			}

			if !match(aaFile, pr) {
				logging.Debugf("The Pull Request doesn't match any rules in AAFile %s.", aaFilePath)
				os.Exit(0)
			}

			if !lo.Contains(pr.Reviewers, user) {
				if err = gh.RequestReview(cmd.Context(), owner, repo, uint32(prNumber)); err != nil {
					logging.Errorf("Failed to request review: %v", err)
					os.Exit(1)
				}
			}

			req := github.SubmitReviewRequest{
				Owner:  owner,
				Repo:   repo,
				Number: prNumber,
				Body:   commentBody,
				Event:  github.SubmitEventApprove,
			}
			if err = gh.SubmitReview(cmd.Context(), req); err != nil {
				logging.Errorf("Failed to submit review: %v", err)
				os.Exit(1)
			}
		},
	}
	cmd.PersistentFlags().BoolVarP(&version, "version", "v", false, "Print version info.")
	cmd.Flags().StringVarP(&ghToken, "gh-token", "t", os.Getenv("GH_TOKEN"), "GitHub token.")
	cmd.Flags().IntVarP(&prNumber, "pr", "p", -1, "Pull request number.")
	cmd.Flags().StringVarP(&githubRepo, "repo", "r", os.Getenv("GITHUB_REPOSITORY"), "GitHub repository. Format: owner/repo")
	cmd.Flags().StringVarP(&aaFilePath, "aafile", "f", "", "Auto-approve file. Leave empty for discovering the file in the repository.")
	cmd.Flags().BoolVar(&githubLogFormat, "github-log-format", false, "Output in GitHub Action log format.")
	cmd.Flags().StringVarP(&commentBody, "comment", "b", "Auto Approved", "Approve comment body.")
	return cmd
}

func parseRepo(ghRepo string) (owner, repo string) {
	parts := strings.Split(ghRepo, "/")
	if len(parts) != 2 {
		logging.Errorf("invalid repository format: %s", ghRepo)
		os.Exit(1)
	}
	return parts[0], parts[1]
}

func match(af aafile.AAFile, pr *github.PullRequest) bool {
	for _, f := range pr.Files {
		if len(af.Match(f.Name).MatchUser(pr.Author)) > 0 {
			return true
		}
	}
	return false
}
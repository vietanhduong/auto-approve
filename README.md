# Auto Approve

GitHub Auto-Approve PR to bypass the "Require at least 1 approver" rule.

## Usage

### Prepare `AUTOAPPROVE`

Before running the `auto-approve` command line, you must provide a `AUTOAPPROVE` file in your target repository. The `AUTOAPPROVE` file has similar syntax with the GitHub's `CODEOWNERS` file.

```
# Global rules
* @vietanhduong

# Wildcard rules
path/to/* @user1
```

The `AUTOAPPROVE` file must be placed in either `AUTOAPPROVE` or `.github/AUTOAPPROVE`. If both locations are specified, the `AUTOAPPROVE` in the root directory will be preferred.

### Permission

The `auto-approve` command requires a GitHub Personal Access Token (PAT) to `get` and `submit approve` Pull Requests in the target repository. Hence, the PAT must have the ability to:

- Read the pull request
  - `pull_requests:read`
  - `content:read`
- Write to the pull request
  - `pull_requests:write`

### `auto-approve` Command Usage

```console
$ auto-approve --help
Auto-Approve the input Pull Request
to bypass the "Require at least 1 reviewer" of GitHub.

Usage:
  auto-approve PR_NUMBER [flags]

Examples:
$ cat "$(git rev-parse --show-toplevel)/AUTOAPPROVE"
# Use similar GitHub CODEOWNERS syntax
* @vietanhduong

$ export GH_TOKEN=$GITHUB_TOKEN
$ export GITHUB_REPOSITORY="<owner>/<repo>"
$ auto-approve $PR_NUMBER --comment "LGTM!"


Flags:
  -f, --aafile string       Auto-approve file. Leave empty for discovering the file in the repository.
  -b, --comment string      Approve comment body. (default "Auto Approved")
  -t, --gh-token string     GitHub token.
      --github-log-format   Output in GitHub Action log format.
  -h, --help                help for auto-approve
  -r, --repo string         GitHub repository. Format: owner/repo
  -v, --version             Print version info.
```

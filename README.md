# GitHub PR comment resource

A Concourse resource for working with pull request comments. It uses the GitHub v4 API (GraphQL).

## Source Configuration

| Parameter                   | Required | Example                          | Description                                                                                                                                                                                                                                                                                |
|-----------------------------|----------|----------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `repository`                | Yes      | `test-owner/test-repository`     | The repository to target.                                                                                                                                                                                                                                                                  |
| `access_token`              | Yes      |                                  | A Github Access Token with repository access (required for setting status on commits). N.B. If you want github-pr-comment to work with a private repository. Set `repo:full` permissions on the access token you create on GitHub. If it is a public repository, `repo:status` is enough. |
| `comments`                  | No       | `^hello world$`                  | Only comments which match this regex will be returned as valid versions.
| `labels`                    | No       | `[wip, help]`                    | PRs which match at least one label will be considered valid.
| `latest_per_pr`             | No       | true                             | Returns only the latest comment per PR which matches previous criteria if any.

## Behavior

### `check`

Produces new versions for all comments that match any of the filters defined in `source`, chronologically sorted. 

A version is represented as follows:

- `comment_id`: The numerical ID of a comment.

## Work left

- Implementation of `in` and `out`.
- Push artifact to docker.io

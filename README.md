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

## Development environment


### Infrastructure setup

There is the `devenv` directory where you can fully test out the resource.

Run the `keys.sh` script in the directory to generate keys for the worker and web instances, and then
run `docker-compose up -d`. This should finish without any errors, and you should be able to see Concourse running
on http://localhost:8080.

### Updating/applying the pipeline

To update the pipeline, you will need to download `fly` and put it into your `PATH` environment variable. Additionally,
you will need to create `devenv/vars.yml` file the following contents:

```yaml
access_token: YOUR_GITHUB_TOKEN
```

Next, you will need to invoke `fly -t dev -c http://localhost:8080`, open the browser, and log in with the username
`test` and password `test`.

Finally, run `apply.sh` and reply with `y` on the prompt. Now you should see your new pipeline in the browser from
where you can unpause it.

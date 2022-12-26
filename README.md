# export-workflow-logs

The logs for a GitHub Actions workflow run are only [available for a limited time](https://docs.github.com/en/organizations/managing-organization-settings/configuring-the-retention-period-for-github-actions-artifacts-and-logs-in-your-organization) before they are automatically deleted. Use this GitHub Action to export the logs to longer term storage such as Amazon S3.

This Action uses the [Download workflow run logs API](https://docs.github.com/en/rest/actions/workflow-runs?apiVersion=2022-11-28#download-workflow-run-logs) to fetch the run logs. The logs are then saved as a `.zip` archive at the destination. For supported destinations, see [Usage](#usage).

## Quick Start

`.github/workflows/my-workflow.yml`

```yml
name: Hello World
on: push
jobs:
  hello-world:
    runs-on: ubuntu-latest
    steps:
      - name: Print Hello World
        run: echo "Hello World!"
```

`.github/workflows/export-my-workflow-logs-to-s3.yml`

```yml
name: Export Hello World Logs To S3
on:
  workflow_run:
    workflows: [Hello World]
    types: [completed]
jobs:
  export-hello-world-logs:
    runs-on: ubuntu-latest
    steps:
      - uses: timorthi/export-workflow-logs@v1
        with:
          repo-token: ${{ secrets.GITHUB_TOKEN }}
          run-id: ${{ github.event.workflow_run.id }}
          destination: s3
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets. AWS_SECRET_ACCESS_KEY }}
          aws-region: us-west-1
          s3-bucket-name: my-workflow-logs
          s3-key: ${{ github.event.workflow_run.name }}/${{ github.event.workflow_run.id }}.zip # You can take advantage of the `workflow_run` event payload to generate a unique name for the exported logs.
```

## Usage

Workflow run logs can only be downloaded on completion of that workflow. To export workflow logs, you will have to run this action in a separate workflow that runs after the conclusion of an upstream workflow (see the [`workflow_run`](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#workflow_run) event). Attempting to export the workflow run logs of an in-progress workflow will result in a 404 error from the GitHub API.

The following inputs are required regardless of the chosen destination:

| Name | Description |
| - | - |
| `repo-token` | Token to use to fetch workflow logs. Typically the `GITHUB_TOKEN` secret. |
| `run-id` | The workflow run ID for which to export logs. Typically obtained via the `github` context per the above example. |
| `destination` | The service to export workflow logs to. Supported values: `s3` |

Further inputs are required and they are dependent on the intended destination of the workflow logs.

### [Amazon S3](https://aws.amazon.com/s3/)

The S3 exporter uses the `S3PutObject` API to save the workflow logs file.

The following inputs are required if `destination` is `s3`:
| Name | Description |
| - | - |
| `aws-access-key-id` | Access Key ID to use to upload workflow logs to S3 |
| `aws-secret-access-key` | Secret Access Key to use to upload workflow logs to S3 |
| `aws-region` | Region of the S3 bucket to upload to. Example: `us-east-1`
| `s3-bucket-name` | Name of the S3 bucket to upload to
| `s3-key` | S3 key to save the workflow logs to

## Development

### Testing

To run unit tests, run `make test`.

To test changes in a GitHub Action, change `action.yml` to point to the Dockerfile:

```yml
runs:
  using: "docker"
  image: "Dockerfile"
  args:
    ...
```

Then, make sure the GitHub Actions workflow that calls this Action references the branch or commit in which you made this change.

```yml
- uses: timorthi/export-workflow-logs@my-feature-branch-1
  with:
    ...
```

This will force the workflow to build the image (and therefore the Go source code) on every run.

### Building & Releasing

TODO

# GitHub Actions Workflows

## sync-me-catalog.yml

**Purpose:** Keeps the ME (Managed Entity) catalog under `knowledge/me-catalog/` in sync with the upstream [`opencord/omci-lib-go`](https://github.com/opencord/omci-lib-go) library.

The workflow runs the ME catalog extractor script on a weekly schedule (every Monday at 06:00 UTC) and opens a pull request automatically if the generated JSON files differ from what is currently committed.

### Triggers

| Trigger | Details |
|---|---|
| `schedule` | Every Monday at 06:00 UTC (`0 6 * * 1`) |
| `workflow_dispatch` | Run manually from the Actions tab at any time |

### How to trigger manually

1. Go to **Actions → Sync ME Catalog** in this repository.
2. Click **Run workflow**, select the branch (usually `main`), and click **Run workflow**.

### Reviewing the generated PR

When changes are detected the workflow opens a PR with:

- Branch: `chore/sync-me-catalog-<run_id>`
- Title: `chore: sync ME catalog from opencord/omci-lib-go`
- Labels: `automation`, `me-catalog`

Review the diff in `knowledge/me-catalog/` to confirm the changes look reasonable (new MEs added, attributes changed, etc.) before merging.

### What to do if the extractor fails

| Symptom | Action |
|---|---|
| **Extractor script not found / TODO step fails** | The extractor has not been added to this repo yet.  Add the script and update the `Run ME catalog extractor` step in `sync-me-catalog.yml` with the real invocation command. |
| **Upstream API changed** | Check the [opencord/omci-lib-go releases](https://github.com/opencord/omci-lib-go/releases) for breaking changes and update the extractor accordingly. |
| **Go version mismatch** | Update the `go-version` input in the `Set up Go` step. |
| **Network error fetching upstream** | Re-run the workflow — transient network errors on GitHub-hosted runners are usually self-healing. |

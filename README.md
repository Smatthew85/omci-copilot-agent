# omci-copilot-agent
Knowledge base and tooling for a GitHub Copilot AI agent that analyzes ITU-T G.988 OMCI messages and helps diagnose failures. Includes a script to extract the Managed Entity (ME) catalog from opencord/omci-lib-go into downloadable JSON files.

## Automation

### ME Catalog Sync

The ME catalog under `knowledge/me-catalog/` is kept in sync with the upstream [`opencord/omci-lib-go`](https://github.com/opencord/omci-lib-go) library automatically.

A scheduled GitHub Actions workflow ([`.github/workflows/sync-me-catalog.yml`](.github/workflows/sync-me-catalog.yml)) runs every Monday at 06:00 UTC, re-extracts the catalog, and opens a pull request if anything has changed.  You can also trigger it manually from the [Actions tab](../../actions/workflows/sync-me-catalog.yml).

See [`.github/workflows/README.md`](.github/workflows/README.md) for details on reviewing generated PRs and troubleshooting extractor failures.

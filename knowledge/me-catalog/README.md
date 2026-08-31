# Managed Entity Catalog

This directory is intended to hold generated JSON files for OMCI Managed Entity (ME) definitions so Copilot can reference real class IDs, names, attributes, access rules, and supported actions.

## Current status

The catalog is **not populated on this branch** because the repository does not currently contain the extractor script or command referenced by the root `README.md`.

## Investigation summary

- **Expected source:** [`opencord/omci-lib-go`](https://github.com/opencord/omci-lib-go)
- **Expected extractor location:** not found in this branch
- **Expected language:** likely Go, based on the repository `.gitignore`
- **Exact invocation:** unavailable because no extractor command or script is checked in
- **Expected output directory:** `knowledge/me-catalog/`
- **Prerequisites:** undetermined; likely Go toolchain and network access to download `github.com/opencord/omci-lib-go`

## What was tried

- Searched the checked-in repository contents for extractor code and scripts.
- Fetched and inspected the other repository branches currently visible from `origin`.
- Searched repository code for `omci-lib-go`, `extractor`, and `sync-me-catalog`.

No runnable extractor implementation was found, so no JSON was generated or committed.

## Source and regeneration

The intended source is `opencord/omci-lib-go` via this repository's extractor, but the extractor script is not present in this branch yet, so there is currently no exact local regeneration command to document.

Once the extractor is added, this README should be updated with:

- the exact regeneration command
- the total Managed Entity count
- example JSON file links
- the committed JSON schema summary

If `.github/workflows/sync-me-catalog.yml` is present in a future branch or after merge, that workflow can keep this directory in sync on a weekly schedule.

## Intended JSON schema

Each generated JSON file is expected to describe one Managed Entity and include fields such as:

- class ID
- ME name
- attributes with index, name, size, access, and mandatory/optional status
- supported actions
- relationships to other Managed Entities

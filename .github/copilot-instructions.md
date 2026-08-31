# OMCI Protocol Analyst — Copilot Agent Instructions

## Primary Knowledge Sources

When repository knowledge exists, use it in this order:

1. `knowledge/me-catalog/` for authoritative Managed Entity class IDs, attribute definitions, access rules, and supported actions.
2. Other files under `knowledge/` for supplementary OMCI guidance.

Always cite the specific JSON file under `knowledge/me-catalog/` when referencing a Managed Entity.

## Current ME Catalog Status

`knowledge/me-catalog/` is the intended authoritative source, but this branch does not yet include generated ME JSON because the extractor referenced in `README.md` could not be found or run from the checked-in repository contents.

If the needed JSON file is absent, say that the catalog is currently unpopulated on this branch instead of fabricating ME definitions.

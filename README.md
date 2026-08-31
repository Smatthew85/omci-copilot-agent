# omci-copilot-agent
Knowledge base and tooling for a GitHub Copilot AI agent that analyzes ITU-T G.988 OMCI messages and helps diagnose failures. The repository README references a script to extract the Managed Entity (ME) catalog from opencord/omci-lib-go into downloadable JSON files, but that extractor is not present in this branch yet.

---

## Copilot Agent Knowledge Base

The `knowledge/` directory and `.github/copilot-instructions.md` provide repository context for a GitHub Copilot AI agent focused on OMCI analysis.

| Path | Description |
|---|---|
| [`knowledge/me-catalog/README.md`](knowledge/me-catalog/README.md) | Managed Entity catalog status, source notes, and regeneration blocker details |
| [`knowledge/me-catalog/INDEX.md`](knowledge/me-catalog/INDEX.md) | Managed Entity index placeholder; no generated JSON is committed yet because the extractor is missing |
| [`knowledge/me-catalog/`](knowledge/me-catalog/) | Authoritative ME catalog location; currently unpopulated on this branch (ME count: 0) |

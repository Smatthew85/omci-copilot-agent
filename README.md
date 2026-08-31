# omci-copilot-agent
Knowledge base and tooling for a GitHub Copilot AI agent that analyzes ITU-T G.988 OMCI messages and helps diagnose failures. Includes a script to extract the Managed Entity (ME) catalog from opencord/omci-lib-go into downloadable JSON files.

---

## Copilot Agent Knowledge Base

The `knowledge/` directory and `.github/copilot-instructions.md` together configure a **GitHub Copilot AI agent** specialized in OMCI protocol analysis and ONU provisioning failure diagnosis.

### Copilot Instructions

[`.github/copilot-instructions.md`](.github/copilot-instructions.md) defines the agent's persona, diagnostic methodology, response format, and guardrails. Copilot automatically applies these instructions when you ask questions in this repository.

### Knowledge Directory

| Path | Description |
|---|---|
| [`knowledge/README.md`](knowledge/README.md) | Index of the knowledge base layout |
| [`knowledge/me-catalog/`](knowledge/me-catalog/) | JSON ME definitions (produced by the extractor tooling) |
| [`knowledge/result-codes/README.md`](knowledge/result-codes/README.md) | G.988 result/reason code reference with typical causes |
| [`knowledge/provisioning-flows/standard-onu-provisioning.md`](knowledge/provisioning-flows/standard-onu-provisioning.md) | Standard ONU provisioning sequence, dependency graph, and failure notes |
| [`knowledge/failure-patterns/README.md`](knowledge/failure-patterns/README.md) | Symptom → likely cause → suggested action diagnostic table |

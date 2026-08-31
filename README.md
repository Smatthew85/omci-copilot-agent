# omci-copilot-agent
Knowledge base and tooling for a GitHub Copilot AI agent that analyzes ITU-T G.988 OMCI messages and helps diagnose failures. Includes a script to extract the Managed Entity (ME) catalog from opencord/omci-lib-go into downloadable JSON files.

## Copilot Agent Knowledge Base

The [`knowledge/`](knowledge/) directory contains reference documentation that grounds the Copilot agent's OMCI analysis:

| Section | Description |
|---|---|
| [`knowledge/message-format/`](knowledge/message-format/) | OMCI frame formats — [baseline 48-byte layout](knowledge/message-format/baseline-frame.md), [extended frame layout](knowledge/message-format/extended-frame.md), and a full [Message Type reference](knowledge/message-format/message-types.md) |
| [`knowledge/result-codes/`](knowledge/result-codes/) | G.988 result and reason codes |
| [`knowledge/provisioning-flows/`](knowledge/provisioning-flows/) | Standard ONU provisioning sequences |
| [`knowledge/failure-patterns/`](knowledge/failure-patterns/) | Symptom → root-cause diagnostic table |

See [`knowledge/README.md`](knowledge/README.md) for the full index.

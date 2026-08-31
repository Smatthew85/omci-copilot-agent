# omci-copilot-agent
Knowledge base and tooling for a GitHub Copilot AI agent that analyzes ITU-T G.988 OMCI messages and helps diagnose failures. Includes a script to extract the Managed Entity (ME) catalog from opencord/omci-lib-go into downloadable JSON files.

## Copilot Agent Knowledge Base

| Resource | Description |
|----------|-------------|
| [`knowledge/result-codes/`](knowledge/result-codes/README.md) | G.988 result codes and message type reference |
| [`knowledge/failure-patterns/`](knowledge/failure-patterns/README.md) | Symptom → root cause → remediation table |
| [`knowledge/provisioning-flows/`](knowledge/provisioning-flows/standard-onu-provisioning.md) | Standard ONU OMCI provisioning sequence |
| [`knowledge/examples/`](knowledge/examples/README.md) | Golden set of annotated OMCI failure cases for few-shot diagnosis |

See [`knowledge/README.md`](knowledge/README.md) for the full knowledge base index.

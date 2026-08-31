# OMCI Copilot Agent — Knowledge Base

This directory contains the reference material used by the Copilot agent to analyze
OMCI messages and diagnose provisioning failures.

---

## Index

| Directory | Contents |
|-----------|----------|
| [`result-codes/`](./result-codes/) | G.988 OMCI result/reason codes and message type reference |
| [`failure-patterns/`](./failure-patterns/) | Symptom → root cause → remediation diagnostic table |
| [`provisioning-flows/`](./provisioning-flows/) | Standard ONU OMCI provisioning sequence with ME dependencies |
| [`examples/`](./examples/) | Golden set of anonymized OMCI failure cases with expected diagnoses (few-shot references) |

---

## Contributing

- Add new result codes or message types to [`result-codes/README.md`](./result-codes/README.md).
- Add new failure patterns to [`failure-patterns/README.md`](./failure-patterns/README.md).
- Add new provisioning flow documentation under [`provisioning-flows/`](./provisioning-flows/).
- Add new golden examples following the guide in [`examples/README.md`](./examples/README.md).

# Copilot Agent Knowledge Base

This directory contains reference documentation used by the GitHub Copilot agent to analyze ITU-T G.988 OMCI messages and diagnose ONU provisioning failures.

## Contents

| Directory | Description |
|---|---|
| [`message-format/`](message-format/) | OMCI frame formats — baseline (48-byte) and extended layouts, Message Type reference |
| [`result-codes/`](result-codes/) | G.988 result and reason codes with diagnostic guidance |
| [`provisioning-flows/`](provisioning-flows/) | Standard ONU provisioning sequences and ME dependency order |
| [`failure-patterns/`](failure-patterns/) | Symptom → root-cause diagnostic table |

## How to Use

The Copilot agent is instructed (via [`.github/copilot-instructions.md`](../.github/copilot-instructions.md)) to consult this knowledge base when:

1. A user pastes a raw hex OMCI frame or a sequence of frames.
2. A user describes a provisioning failure symptom.
3. A user asks about a specific ME class, attribute, or result code.

Start with [`message-format/`](message-format/) to decode the frame structure, then cross-reference [`result-codes/`](result-codes/) for response interpretation, and consult [`failure-patterns/`](failure-patterns/) for root-cause guidance.

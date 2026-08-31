# OMCI AVC — Knowledge Base

This directory is the authoritative reference for interpreting OMCI **Attribute Value Change**
(AVC) notifications (Message Type 17, `0x11`) as defined in ITU-T G.988.

AVC messages are **autonomous ONU→OLT** notifications reporting that one or more ME
attributes have changed their values. They carry no AR/AK bits. Attributes must be
designated AVC-capable in the ME definition to appear in these messages.

---

## Index

| Document | Description |
|---|---|
| [`avc-message-format.md`](avc-message-format.md) | Baseline and extended frame layouts, attribute mask semantics, decoded example |
| [`common-avc-triggers.md`](common-avc-triggers.md) | Typical AVC-generating scenarios per ME (ANI-G, ONU-G, PPTP Eth UNI, Software Image) |
| [`interpretation-workflow.md`](interpretation-workflow.md) | Step-by-step recipe for the Copilot agent when a user pastes an AVC frame |

---

## Quick Reference

| Field | Value |
|---|---|
| Message Type byte | `0x11` (17 decimal) |
| Direction | ONU → OLT (autonomous; no AR/AK) |
| Frame size | 48 bytes (baseline) or variable (extended) |
| Attribute mask | 2 bytes (MSB-first, bit N = attribute N per G.988 §11.2.2) |
| New values | Follow the attribute mask in G.988 attribute order, padded to fill content area |

---

## Scope

- **In scope:** AVC notification format, attribute mask decoding, common AVC-generating
  scenarios for well-known MEs.
- **Out of scope:** Alarm notifications (see [`../alarms/`](../alarms/)),
  provisioning-time failures (see [`../failure-patterns/`](../failure-patterns/)).

---

## Contributing

- **New AVC trigger scenarios** — add rows to the relevant ME table in
  [`common-avc-triggers.md`](common-avc-triggers.md).
- **Per-ME AVC-capable attribute lists** — consult `knowledge/me-catalog/NNN-slug.json`
  (once populated) for the authoritative AVC capability flags per attribute.

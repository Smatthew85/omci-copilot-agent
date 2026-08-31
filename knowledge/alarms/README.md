# OMCI Alarms — Knowledge Base

This directory is the authoritative reference for interpreting OMCI **Alarm** notifications
(Message Type 16, `0x10`) as defined in ITU-T G.988.

Alarms are **autonomous ONU→OLT** messages that report fault conditions on Managed
Entities. They carry no AR/AK bits (no acknowledge handshake). The OLT must track alarm
state and sequence numbers to detect missed notifications.

---

## Index

| Document | Description |
|---|---|
| [`alarm-message-format.md`](alarm-message-format.md) | Baseline and extended frame layouts, alarm bitmap semantics, sequence number, decoded example |
| [`alarm-synchronization.md`](alarm-synchronization.md) | Sequence-number gap detection, Get All Alarms (MT 27/28) resync workflow |
| [`common-alarms.md`](common-alarms.md) | Per-ME alarm bit tables for high-traffic MEs (ANI-G, ONU-G, PPTP Eth UNI, Circuit Pack) |
| [`interpretation-workflow.md`](interpretation-workflow.md) | Step-by-step recipe for the Copilot agent when a user pastes an alarm frame |

---

## Quick Reference

| Field | Value |
|---|---|
| Message Type byte | `0x10` (16 decimal) |
| Direction | ONU → OLT (autonomous; no AR/AK) |
| Frame size | 48 bytes (baseline) or variable (extended) |
| Alarm bitmap | 28 bytes (224 alarm bits) within the content area |
| Sequence number | 1 byte, last content byte before the trailer |
| Resync messages | Get All Alarms (MT 27), Get All Alarms Next (MT 28) |

---

## Scope

- **In scope:** Alarm notification format, alarm bitmap decoding, sequence tracking,
  resync procedure, per-ME alarm definitions for commonly observed MEs.
- **Out of scope:** Provisioning-time failures (see [`../failure-patterns/`](../failure-patterns/)),
  AVC notifications (see [`../avc/`](../avc/)).

---

## Contributing

- **New per-ME alarm tables** — add a subsection to [`common-alarms.md`](common-alarms.md)
  following the existing table format.
- **New failure patterns** triggered by alarm conditions — cross-reference
  [`../failure-patterns/README.md`](../failure-patterns/README.md).
- Consult `knowledge/me-catalog/NNN-slug.json` (once populated) for the authoritative
  per-ME alarm bitmap definitions.

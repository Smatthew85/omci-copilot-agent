# Alarm Interpretation Workflow

This document provides the step-by-step recipe the Copilot agent follows when a user
pastes an OMCI frame whose MT byte is `0x10` (Alarm notification).

---

## Step 1 — Detect MT = 0x10

Confirm the Message Type byte (byte 2 in a baseline OMCI frame) equals `0x10`.

- If the DevID (byte 3) is `0x0A`, use the **baseline frame layout**.
- If the DevID is `0x0B`, use the **extended frame layout**.

See [`alarm-message-format.md`](alarm-message-format.md) for the full layout specification.

**If the MT byte is not `0x10`,** stop this workflow. Use the AVC workflow if MT = `0x11`
(see `knowledge/avc/interpretation-workflow.md`), or the general decode workflow for other
message types.

---

## Step 2 — Extract ME Class and Instance

From the baseline frame:

```
ME Class ID  = bytes 4–5 (big-endian)
ME Instance  = bytes 6–7 (big-endian)
```

Look up the ME Class ID in `knowledge/me-catalog/` (if the JSON file is present) to
confirm the ME name and understand the supported alarm bits for this ME class.

---

## Step 3 — Extract the 28-Byte Alarm Bitmap

```
Alarm bitmap = bytes 8–35 (28 bytes, 224 bits)
Sequence number = byte 39
```

Enumerate every set bit in the 28-byte bitmap. For bit number N (0-indexed):

```
byte_index = N // 8          (0-based index into the 28-byte bitmap)
bit_position = 7 - (N % 8)  (MSB = bit 7 = alarm 0 within that byte)
active = (bitmap[byte_index] >> bit_position) & 1
```

State the alarm bit indices explicitly in your output so the operator can verify them
against the raw hex frame.

---

## Step 4 — Look Up Alarm Meanings

For each set bit N, determine the alarm meaning in the following priority order:

1. **[`common-alarms.md`](common-alarms.md)** — check the per-ME table for this ME class.
2. **`knowledge/me-catalog/NNN-slug.json`** (once populated) — check the `alarms` array
   for the ME, if present.
3. **G.988 §9.x** — the alarm list for the ME class is in the ME definition clause.

If none of the above provide a meaning for bit N, report explicitly:

> `<ME name> (Class <ID>) alarm bit <N> — meaning not in local catalog;
> consult G.988 §<clause> or knowledge/me-catalog/<file> (once available).`

**Do not guess alarm meanings.** An incorrect alarm interpretation can lead to wrong
remediation actions.

---

## Step 5 — Note the Alarm Sequence Number

Record the sequence number from byte 40. Compare it with the last-seen sequence number
for this ONU (if known from context).

- If there is a gap: note the missed range and recommend a Get All Alarms resync.
  See [`alarm-synchronization.md`](alarm-synchronization.md).
- If the sequence is contiguous: state that alarm state is current.

---

## Step 6 — Correlate with Prior AVC or Provisioning Events

If the user has provided context:

- **Alarm follows a provisioning failure:** An alarm raised shortly after a failed Create
  or Set may indicate that the ONU is in an inconsistent state. Cross-reference
  [`../failure-patterns/README.md`](../failure-patterns/README.md).
- **Alarm follows an AVC:** An AVC reporting a degraded attribute value (e.g., low optical
  power) may be followed by an alarm on the same ME. Treat as corroborating evidence.
- **Alarm cleared (all bits zero):** A zero-bitmap alarm notification means all alarms on
  this ME instance have been cleared. Update the alarm state accordingly.

---

## Step 7 — Produce Diagnosis

Structure your output using the standard format:

### Summary
One or two sentences: which ME and instance, which alarms are active (by name and bit
index), and the high-level consequence.

### Root Cause
The underlying fault condition indicated by the active alarm bits, referencing G.988 ME
behavior and the alarm definitions in [`common-alarms.md`](common-alarms.md).

### Evidence
- **Message Type:** 16 (0x10) — Alarm
- **ME Class:** class ID and name
- **ME Instance:** instance ID
- **Active alarm bits:** list each set bit as `bit N (<alarm name>)`
- **Alarm sequence number:** decimal value
- Raw bitmap bytes (first non-zero byte at minimum)

### Remediation
Numbered steps to recover. Reference G.988 message types by number and name where
applicable (e.g., MIB Reset (MT 9/10), Get All Alarms (MT 27/28)).

### References
- [`alarm-message-format.md`](alarm-message-format.md)
- [`common-alarms.md`](common-alarms.md)
- [`alarm-synchronization.md`](alarm-synchronization.md) if sequence gap detected
- Relevant G.988 clause(s)
- `knowledge/me-catalog/NNN-slug.json` if applicable

---

## Quick Checklist

```
[ ] MT byte = 0x10 confirmed
[ ] ME Class ID and Instance extracted
[ ] 28-byte alarm bitmap extracted; all set bits enumerated (state indices explicitly)
[ ] Each set bit meaning looked up in common-alarms.md or me-catalog
[ ] Any unknown bit meanings called out explicitly (not guessed)
[ ] Alarm sequence number recorded; gap detected or sequence contiguous noted
[ ] Correlated with prior AVC or provisioning context if provided
[ ] Output structured as Summary/Root Cause/Evidence/Remediation/References
```

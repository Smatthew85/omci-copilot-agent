# AVC Interpretation Workflow

This document provides the step-by-step recipe the Copilot agent follows when a user
pastes an OMCI frame whose MT byte is `0x11` (Attribute Value Change notification).

---

## Step 1 — Detect MT = 0x11

Confirm the Message Type byte (byte 2 in a baseline OMCI frame) equals `0x11`.

- If the DevID (byte 3) is `0x0A`, use the **baseline frame layout**.
- If the DevID is `0x0B`, use the **extended frame layout**.

See [`avc-message-format.md`](avc-message-format.md) for the full layout specification.

**If the MT byte is not `0x11`,** stop this workflow. Use the alarm workflow if MT = `0x10`
(see `knowledge/alarms/interpretation-workflow.md`), or the general decode workflow for
other message types.

---

## Step 2 — Extract ME Class and Instance

From the baseline frame:

```
ME Class ID  = bytes 4–5 (big-endian)
ME Instance  = bytes 6–7 (big-endian)
```

Look up the ME Class ID in `knowledge/me-catalog/` (if the JSON file is present) to
confirm the ME name and understand which attributes are AVC-capable for this ME class.

---

## Step 3 — Extract the 2-Byte Attribute Mask

```
Attribute mask = bytes 8–9 (big-endian, 2 bytes)
```

Enumerate set bits using G.988 numbering (bit 15 of the 2-byte mask = attribute 1;
bit 0 = attribute 16):

```
for N in range(1, 17):
    if mask_bit_set(mask, N):
        attribute N is reported as changed
```

State the attribute indices explicitly in your output so the operator can cross-check
them against the raw frame.

---

## Step 4 — Decode New Attribute Values

Starting at byte 10, extract each reported attribute value in ascending attribute-index
order, at the natural size defined for that attribute in G.988 or in the ME catalog JSON.

For each set bit N:

1. Identify attribute N's name and size from:
   - `knowledge/me-catalog/NNN-slug.json` (preferred, once populated), or
   - G.988 §9.x for the ME class.
2. Read the appropriate number of bytes from the content area.
3. Interpret the raw bytes according to the attribute's encoding (unsigned integer,
   two's-complement, enumeration, bitmask, string, etc.).
4. If the attribute meaning is not in the local catalog or this document, report:
   > `<ME name> (Class <ID>) attribute <N> — meaning not in local catalog;
   > consult G.988 §<clause> or knowledge/me-catalog/<file> (once available).`

Reference [`common-avc-triggers.md`](common-avc-triggers.md) for well-known AVC-capable
attributes and their value semantics.

---

## Step 5 — Correlate with Recent Configuration or Events

If the user has provided context about recent OLT actions or provisioning events:

- **Expected AVC:** If the OLT recently pushed a `Set` to attribute N of the same ME
  instance, an AVC on that attribute is an expected confirmation, not a fault. Note this
  explicitly in your output.
- **Unexpected AVC:** If no prior Set was made, the ONU autonomously detected a value
  change (e.g., optical power fluctuation, link state change, software upgrade event).
  Treat this as informational or as a signal of an environmental/hardware event.
- **Software upgrade correlation:** AVCs on Software Image (ME 7) attributes 3–5 during
  an upgrade sequence are expected. See [`common-avc-triggers.md`](common-avc-triggers.md).

---

## Step 6 — Produce Diagnosis

Structure your output using the standard format:

### Summary
One or two sentences: which ME, which attribute(s) changed, and the high-level meaning.

### Root Cause
Why this AVC was generated — environmental trigger, OLT-initiated change, or upgrade
event.

### Evidence
- **Message Type:** 17 (0x11) — Attribute Value Change
- **ME Class:** class ID and name
- **ME Instance:** instance ID
- **Attribute mask:** hex value, then list each set bit as `attribute N (<name>)`
- **New values:** decoded value for each changed attribute, with units if applicable

### Remediation
Numbered steps only if action is needed (e.g., verify threshold settings, check optical
link). If the AVC is informational and expected, state that no action is required.

### References
- [`avc-message-format.md`](avc-message-format.md)
- [`common-avc-triggers.md`](common-avc-triggers.md)
- Relevant G.988 clause(s)
- `knowledge/me-catalog/NNN-slug.json` if applicable

---

## Quick Checklist

```
[ ] MT byte = 0x11 confirmed
[ ] ME Class ID and Instance extracted
[ ] Attribute mask decoded; set bits enumerated (state indices explicitly)
[ ] New values decoded for each set bit
[ ] Checked common-avc-triggers.md for known AVC scenarios
[ ] Correlated with any prior OLT Set or provisioning context
[ ] Output structured as Summary/Root Cause/Evidence/Remediation/References
[ ] Any unknown attribute meanings called out explicitly (not guessed)
```

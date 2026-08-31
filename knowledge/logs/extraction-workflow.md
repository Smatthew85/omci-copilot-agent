# OMCI Frame Extraction Workflow

This document describes the end-to-end process the agent follows when a user
pastes log text instead of pre-cleaned hex frames.

---

## Step-by-Step Recipe

### Step 1 — Detect Log Source

Examine the pasted text for structural clues:

| Observation | Conclusion |
|---|---|
| Lines are valid JSON with a `"logger"` field | VOLTHA log |
| `"logger"` contains `omci-cc`, `MibDownloadFsm`, `UniVlanConfigFsm` | `openonu-adapter-go` → use [`voltha-openonu-adapter.md`](voltha-openonu-adapter.md) |
| `"logger"` contains `openolt`, `flowMgr`, `core-proxy` | `openolt-adapter` → use [`voltha-openolt-adapter.md`](voltha-openolt-adapter.md) |
| Text lines with syslog-style timestamp and `OMCI:` / `PON OMCI Msg:` prefix | BAL / OpenOLT agent → use [`bal-openolt-agent.md`](bal-openolt-agent.md) |
| Format unrecognized | **Ask the user** to confirm the log source before proceeding |

### Step 2 — Scan for OMCI Carrier Fields

Once the source is identified, search each relevant line for the OMCI hex payload.
Check field names in priority order:

1. `omci-message`
2. `omci-msg`
3. `TxOmciMessage`
4. `RxOmciMessage`
5. `omciMsg`
6. `pkt`
7. `packet`
8. `raw`
9. Inline after `OMCI:`, `OMCI TX:`, `OMCI RX:`, or `PON OMCI Msg:` text prefix

### Step 3 — Extract Hex

For each carrier field or prefix found:

1. Take the raw string value.
2. Strip leading/trailing `"` characters (JSON quoting).
3. Remove all spaces (` `) and colons (`:`).
4. Remove any `0x` prefix if present.
5. Convert to lowercase for consistency.

### Step 4 — Validate Frame Length

| Hex string length | Frame type | Action |
|---|---|---|
| 96 chars (48 bytes) | Baseline OMCI frame | Decode with [`baseline-frame.md`](../message-format/baseline-frame.md) |
| > 96 chars | Possibly extended OMCI | Check bytes 8–11 for `0800` (extended indicator); if present decode with [`extended-frame.md`](../message-format/extended-frame.md) |
| < 96 chars | Truncated or non-OMCI | Skip; note truncation |
| Not hex characters | Non-OMCI field | Skip |

### Step 5 — Group by Device

Assign each valid frame to an ONU using available identifiers:

| Log Source | Grouping Key |
|---|---|
| VOLTHA (either adapter) | `device-id` field |
| BAL / OpenOLT agent | `PON[<intf-id>] ONU[<onu-id>]` pair |

If multiple ONUs appear in the same paste, process each group independently.

### Step 6 — Order Chronologically

Within each device group:

- **VOLTHA logs:** sort by `ts` (Unix float64, ascending).
- **BAL logs:** sort by syslog timestamp string (ascending).

### Step 7 — Correlate TX/RX Pairs

Each OMCI baseline request has a corresponding response (unless the ONU did not
reply or the log is incomplete). Correlate pairs by TCID:

1. Extract bytes 0–1 (first 4 hex chars) of each frame → TCID.
2. Match each TX frame (AR bit set; `msg` contains `tx`) with an RX frame
   (AK bit set; `msg` contains `rx`) that has the same TCID.
3. Flag TCIDs with a TX but no RX as potential timeouts or missing log lines.

AR/AK bit positions: Message Type byte (offset 2, hex chars 4–5):
- Bit 6 (0x40) set → AR (Action Request, expecting reply)
- Bit 5 (0x20) set → AK (Acknowledgement, this is the reply)

### Step 8 — Hand Off to Decoder

Pass each extracted and validated hex string to the appropriate decoder:

- Baseline (96 chars): [`knowledge/message-format/baseline-frame.md`](../message-format/baseline-frame.md)
- Extended (>96 chars): [`knowledge/message-format/extended-frame.md`](../message-format/extended-frame.md)

**Always report the extracted hex frame(s) explicitly before decoding**, so the
user can verify that extraction was correct.

### Step 9 — Produce Diagnosis

Follow the standard output structure from `.github/copilot-instructions.md`:

1. **Summary** — ME, action, result code, consequence.
2. **Root Cause** — underlying reason referencing G.988 behavior or provisioning constraints.
3. **Evidence** — Message Type, ME Class, ME Instance, Result Code, attribute details.
4. **Remediation** — numbered recovery steps.
5. **References** — knowledge base links and G.988 clauses.

---

## Worked Example

### Input — 3-line VOLTHA log paste

```json
{"level":"debug","ts":1693000010.123,"logger":"omci-cc","caller":"omci_cc.go:312","msg":"omci-message-tx","device-id":"0001010000000001","intf-id":0,"onu-id":1,"omci-message":"000100094f0000000000000000000000000000000000000000000000000000000000000000000000000000000000028a"}
{"level":"debug","ts":1693000010.456,"logger":"omci-cc","caller":"omci_cc.go:389","msg":"omci-message-rx","device-id":"0001010000000001","intf-id":0,"onu-id":1,"omci-message":"000129094f0000000000000000000000000000000000000000000000000000000000000000000000000000000000028a"}
{"level":"debug","ts":1693000015.100,"logger":"MibDownloadFsm","caller":"mib_download_fsm.go:245","msg":"MibDownloadFsm - send create TCONT","device-id":"0001010000000001","omci-message":"0005440a010600010000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00000028a"}
```

### Identification

- Source: VOLTHA `openonu-adapter-go` (JSON, `logger` = `omci-cc` / `MibDownloadFsm`)

### Extracted Frames

| # | ts | Direction | TCID | Hex (truncated for display) |
|---|---|---|---|---|
| 1 | 1693000010.123 | TX (OLT→ONU) | `0001` | `000100094f00…028a` |
| 2 | 1693000010.456 | RX (ONU→OLT) | `0001` | `000129094f00…028a` |
| 3 | 1693000015.100 | TX (OLT→ONU) | `0005` | `0005440a0106…028a` |

### TX/RX Correlation

- Frames 1 and 2 share TCID `0001` → request/response pair for MIB Reset.
- Frame 3 (TCID `0005`) has no matching RX in this excerpt → response not captured or not yet received.

### Decoded Summary

**Frame 1 (TX):** MIB Reset request — MT `0x09` (AR), ME ONU-G (class bytes `0x00 0x4f` = `0x004f` = 79 decimal = ONU-G, instance `0x0000`).

**Frame 2 (RX):** MIB Reset response — MT `0x29` (AK, `0x09` base), result byte at offset 8 = `0x00` → Success.

**Frame 3 (TX):** Create T-CONT (MT `0x44` = Create, ME class `0x0106` = 262 decimal = T-CONT, instance `0x0001`). Response not shown — provisioning may still be in progress or the response was not captured in this log excerpt.

### Diagnosis

The MIB Reset completed successfully (result `0x00`). The log excerpt ends mid-provisioning at the T-CONT Create step. No failure is evident in the provided frames; a longer log excerpt is needed to diagnose any provisioning issue.

# openonu-adapter-go Log Format

`openonu-adapter-go` is the primary source of OMCI traffic visibility in VOLTHA.
It originates and receives virtually every OMCI message exchanged with the ONU and
writes each one to a structured log line.

---

## Log Framework

`openonu-adapter-go` uses the `voltha-lib-go` logging package, which wraps
[zap](https://github.com/uber-go/zap) and emits **structured JSON** to stdout.

### Mandatory Fields Present on Every Log Line

| Field | Type | Description |
|---|---|---|
| `level` | string | `debug`, `info`, `warn`, `error` |
| `ts` | float64 | Unix timestamp with sub-second precision (e.g. `1693000012.345`) |
| `logger` | string | Subsystem that produced the line (see below) |
| `caller` | string | `<file>.go:<line>` of the calling function |
| `msg` | string | Human-readable event description |

### Common Contextual Fields

| Field | Type | Description |
|---|---|---|
| `device-id` | string | VOLTHA logical device ID for the ONU — groups all frames for one ONU |
| `parent-id` | string | OLT logical device ID |
| `onu-id` | uint32 | ONU-ID on the PON interface |
| `intf-id` | uint32 | PON interface (port) number |

---

## Common OMCI-Related Loggers

| `logger` Value | Subsystem |
|---|---|
| `omci-cc` | OMCI Communication Channel — sends and receives all OMCI frames |
| `omciCC` | Alias used in some releases |
| `omci-msg` | Lower-level OMCI message handler |
| `onuDeviceEntry` | ONU device state machine — drives MIB upload/download |
| `MibDownloadFsm` | MIB download finite state machine |
| `UniVlanConfigFsm` | UNI VLAN configuration state machine — sets Extended VLAN ME |
| `omci-flow-mgr` | OMCI-level flow manager (commonly seen in newer releases) |

---

## Where OMCI Hex Appears

OMCI frame bytes are logged as lowercase hex strings. Depending on the log call
site the field name varies:

| Field Name | Notes |
|---|---|
| `omci-message` | Most common; used by `omci-cc` send/receive paths |
| `omci-msg` | Alternate name in some `omci-cc` versions |
| `TxOmciMessage` | TX path in certain release branches |
| `RxOmciMessage` | RX path in certain release branches |
| `omciMsg` | Camel-case variant |
| `packet` | Used when the full Ethernet/OMCI encapsulation is included |
| `raw` | Raw bytes when logged via generic packet dump helpers |

> **Field names marked `TxOmciMessage`, `RxOmciMessage`, and `omciMsg` are
> "commonly seen" rather than authoritative** — exact names may differ across
> adapter releases. When in doubt, look for a 96-character (baseline) or longer
> hex string in any field of a line whose `logger` is one of the OMCI loggers
> above.

---

## Representative Log Line Shapes

All examples use placeholder device IDs and fabricated but structurally valid hex.

### 1. OMCI TX — OLT→ONU (MIB Reset request)

```json
{"level":"debug","ts":1693000010.123,"logger":"omci-cc","caller":"omci_cc.go:312","msg":"omci-message-tx","device-id":"0001010000000001","intf-id":0,"onu-id":1,"omci-message":"000100094f0000000000000000000000000000000000000000000000000000000000000000000000000000000000028a"}
```

- `msg` = `"omci-message-tx"` → outbound frame
- `omci-message` field holds the 96-char (48-byte baseline) hex string
- TCID = first 4 hex chars = `0001`

### 2. OMCI RX — ONU→OLT (MIB Reset response)

```json
{"level":"debug","ts":1693000010.456,"logger":"omci-cc","caller":"omci_cc.go:389","msg":"omci-message-rx","device-id":"0001010000000001","intf-id":0,"onu-id":1,"omci-message":"000129094f0000000000000000000000000000000000000000000000000000000000000000000000000000000000028a"}
```

- `msg` = `"omci-message-rx"` → inbound frame (AK bit set, MT byte `0x29` = MIB Reset response)
- Same TCID `0001` correlates to the TX above

### 3. MIB Upload Next response (RX)

```json
{"level":"debug","ts":1693000015.789,"logger":"omci-cc","caller":"omci_cc.go:389","msg":"omci-message-rx","device-id":"0001010000000001","intf-id":0,"onu-id":1,"seq-no":3,"omci-message":"00032e0a002d000003000f8000484741432d47393934344d0000000000000000000000000000000000000000000028a0"}
```

- `seq-no` = MIB Upload sequence number (useful for ordering out-of-order responses)
- ME class `0x002d` = 45 decimal = ONU-G

### 4. FSM state transition with embedded OMCI reference

```json
{"level":"info","ts":1693000020.001,"logger":"MibDownloadFsm","caller":"mib_download_fsm.go:245","msg":"MibDownloadFsm - send create TCONT","device-id":"0001010000000001","omci-message":"0005440a010600010000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00000028a"}
```

- `logger` = `"MibDownloadFsm"` indicates this is part of the MIB download sequence
- `msg` describes the semantic action; `omci-message` holds the frame

---

## Field Extraction Rules

1. **Identify OMCI log lines** — filter for lines where `logger` is one of the OMCI
   loggers listed above, or where `msg` contains `omci-message-tx`, `omci-message-rx`,
   `send`, or `receive`.
2. **Locate the hex field** — check fields in priority order: `omci-message` →
   `omci-msg` → `TxOmciMessage` → `RxOmciMessage` → `omciMsg` → `packet` → `raw`.
3. **Strip JSON quoting** — remove the surrounding `"` characters.
4. **Remove whitespace and colons** — some log paths insert spaces or colons between
   byte pairs; strip them all.
5. **Validate length**:
   - **96 hex chars** (48 bytes) → baseline OMCI frame
   - **≥ 10 hex chars with `0x0800` at offset 8–11** → extended OMCI frame (see
     [`knowledge/message-format/extended-frame.md`](../message-format/extended-frame.md))
   - Anything else → log truncation or non-OMCI field; skip.
6. **Note direction** — `msg` containing `tx` or `Tx` → OLT→ONU; `rx` or `Rx` →
   ONU→OLT.

---

## Correlation Tips

| Technique | How |
|---|---|
| **Group by ONU** | All frames with the same `device-id` belong to one ONU |
| **Correlate TX/RX pairs** | First 4 hex chars of the frame = TCID (Transaction Correlation ID); a TX and RX with the same TCID are a request/response pair |
| **Order frames** | Sort by `ts` (Unix float) within a `device-id` group |
| **Detect retransmits** | Two TX lines with the same TCID and no intervening RX → timeout/retransmit |

---

## How to Extract (Verbatim Recipe)

1. Split the pasted text into individual JSON lines.
2. Parse each line as JSON; skip lines that fail to parse.
3. Discard lines where `logger` is not one of the OMCI loggers listed above AND
   `msg` does not contain `omci`.
4. For each remaining line, find the hex value by checking field names in priority
   order: `omci-message`, `omci-msg`, `TxOmciMessage`, `RxOmciMessage`, `omciMsg`,
   `packet`, `raw`.
5. Strip `"`, whitespace, and `:` from the value.
6. If `len(hex) == 96` → baseline frame; if `len(hex) > 96` → check for extended
   marker; otherwise skip.
7. Record `{ts, device-id, direction, hex}` for each valid frame.
8. Sort records by `ts`.
9. Pass each `hex` string to the baseline or extended frame decoder.

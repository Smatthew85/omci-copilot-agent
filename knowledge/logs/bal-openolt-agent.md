# BAL / OpenOLT Agent Log Format

The OLT hardware-side logs come from the **OpenOLT agent** (the gRPC daemon running
on the OLT) and from the underlying **Broadcom Abstraction Layer (BAL)**. Unlike
VOLTHA logs, these are **text-based** (not JSON) and use a syslog-style format.

> **Vendor note:** Log format, prefix strings, and verbosity options vary across
> OLT hardware vendors (Edgecore, Adtran, Radisys, and others). The shapes below
> represent the general / most commonly observed patterns rather than one vendor's
> exact output. Treat field positions as approximate.

---

## Typical Log Locations on OLT Hardware

| File | Contents |
|---|---|
| `/var/log/openolt.log` | OpenOLT agent log (most relevant) |
| `/var/log/bal.log` or `/var/log/bal_core.log` | BAL core log (lower level) |

Logs may also be sent to syslog or a remote aggregator depending on the OLT
configuration.

---

## Line Format

```
<timestamp> <hostname> <process>[<pid>]: <level> <message>
```

Example (non-OMCI line for reference):

```
2023-08-25T10:00:05.123456+00:00 olt-01 openolt[1234]: INFO  PON[0] ONU[1] activated SN=HGAC1234ABCD
```

OMCI lines follow the same prefix pattern but include a hex payload.

---

## Recognizing OMCI Lines

Look for one of these prefixes in the `<message>` portion:

| Prefix | Direction | Notes |
|---|---|---|
| `OMCI:` | Either | Most common generic prefix |
| `PON OMCI Msg:` | Either | Seen in some Broadcom BAL releases |
| `OMCI TX:` | OLT→ONU | Explicitly tagged transmit |
| `OMCI RX:` | ONU→OLT | Explicitly tagged receive |
| `omci_indication` | ONU→OLT | Used in some OpenOLT agent debug builds |

---

## Example OMCI Line and Hex Extraction

### Raw log line

```
2023-08-25T10:00:12.456789+00:00 olt-01 openolt[1234]: DEBUG PON[0] ONU[1] OMCI TX: 00 01 00 09 4f 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 02 8a
```

### Extraction steps

1. Locate `OMCI TX:` (or whichever prefix matched).
2. Take everything after the prefix and the following space: `00 01 00 09 …`.
3. Remove all spaces: `000100094f000000…`.
4. Remove any colons: not present in this example, but some vendors use `00:01:00:09`.
5. Result: `000100094f0000000000000000000000000000000000000000000000000000000000000000000000000000000000028a`
6. Validate: 96 hex chars → 48-byte baseline frame. ✓

### Extracted hex

```
000100094f0000000000000000000000000000000000000000000000000000000000000000000000000000000000028a
```

This is a MIB Reset request (TCID=`0001`, MT=`0x09`, ME class bytes `0x00 0x4f` = `0x004f` = 79 decimal = ONU-G, instance=`0x0000`).

---

## PON Port and ONU ID

BAL logs identify the ONU by `PON[<intf-id>]` and `ONU[<onu-id>]` rather than a
VOLTHA `device-id`. Use these two values together to group frames for the same ONU:

```
PON[0] ONU[1]  →  intf-id=0, onu-id=1
```

If you also have the openolt-adapter VOLTHA log, cross-reference the `intf-id` and
`onu-id` fields to find the corresponding `device-id`.

---

## Common Non-OMCI Noise to Ignore

| Prefix / Content | What It Is |
|---|---|
| `PLOAM` | Physical Layer OAM — not OMCI, different protocol |
| `GEM port` / `GEM frame` | GEM layer events, not OMCI messages |
| `DBA` / `BW map` / `bandwidth map` | Dynamic Bandwidth Assignment — scheduler events |
| `OAM` without `OMCI` | IEEE 802.3ah OAM (used on some Ethernet-based OLTs) |
| `RSSI` / `optical power` | Physical layer measurements |
| `ALLOC-ID` / `alloc_id` | T-CONT alloc assignment events |

Skip these lines entirely during OMCI extraction.

---

## Direction Limitations

BAL and OpenOLT agent logs often show **only one direction** unless verbose logging
is explicitly enabled:

- Default: TX only (frames sent from OLT to ONU) are logged.
- With `--log-level=debug` or BAL `verbose=true`: both TX and RX are logged.

If only TX frames are visible, you can still identify the request type and infer
the expected response. Note this limitation explicitly when producing a diagnosis:
> "Only OLT→ONU (TX) frames are available; ONU responses are not visible in this
> log. Diagnosis is based on the TX sequence only."

---

## How to Extract (Verbatim Recipe)

1. Split the pasted text into lines.
2. For each line, check whether it contains one of the OMCI prefixes:
   `OMCI:`, `PON OMCI Msg:`, `OMCI TX:`, `OMCI RX:`, `omci_indication`.
3. Note the direction from the prefix (`TX` → OLT→ONU, `RX` → ONU→OLT,
   ambiguous → unknown).
4. Extract everything after the prefix (and any following whitespace).
5. Remove spaces, colons, and any trailing newline characters.
6. Validate length (96 chars → baseline; longer → check for extended marker).
7. Record `{timestamp, intf-id, onu-id, direction, hex}`.
8. Sort records by timestamp.
9. Pass each hex string to the decoder.

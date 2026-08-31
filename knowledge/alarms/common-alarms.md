# Common Alarms Reference

> **Note:** This list documents alarms **defined in ITU-T G.988**. Do not add alarm
> bit mappings that are not grounded in the spec or in a vendor document explicitly
> cited in the row. The per-ME alarm bitmap is authoritative in
> `knowledge/me-catalog/NNN-slug.json` (once populated) and in G.988 §9.x.
>
> This reference is **representative, not exhaustive**. Consult G.988 directly for the
> complete per-ME alarm list, especially for MEs not covered here.

---

## ANI-G (Class 263)

The ANI-G ME represents the ONU's physical PON interface. Its alarms report optical
layer fault conditions. G.988 §9.2.1.

| Bit | Alarm name | Typical cause |
|---|---|---|
| 0 | Low received optical power | Rx optical level below the low threshold |
| 1 | High received optical power | Rx optical level above the high threshold |
| 2 | Signal fail | BER has crossed the signal-fail threshold (severe degradation) |
| 3 | Signal degrade | BER has crossed the signal-degrade threshold (early warning) |
| 4 | Transmit optical level (low) | Tx output power below the low threshold |
| 5 | Transmit optical level (high) | Tx output power above the high threshold |
| 6 | Laser bias current (high) | Laser bias current exceeds the high threshold |

> Consult G.988 §9.2.1 and `knowledge/me-catalog/263-ani-g.json` (once available) for
> additional alarm bits and their exact threshold semantics.

---

## PPTP Ethernet UNI (Class 11)

The PPTP Ethernet UNI ME represents a physical Ethernet UNI port. G.988 §9.5.1.

| Bit | Alarm name | Typical cause |
|---|---|---|
| 0 | LAN loss of signal | No signal detected on the UNI port (link down) |

> Consult G.988 §9.5.1 and `knowledge/me-catalog/011-pptp-ethernet-uni.json` (once
> available) for any additional alarm bits.

---

## ONU-G (Class 256)

The ONU-G ME represents the ONU as a whole and carries system-level alarms.
G.988 §9.1.1.

| Bit | Alarm name | Typical cause |
|---|---|---|
| 0 | Equipment alarm | General hardware fault on the ONU |
| 1 | Powering alarm | Power supply anomaly (voltage or current out of range) |
| 2 | Battery missing | Battery backup absent when expected |
| 3 | Battery failure | Battery present but detected as failed |
| 4 | Battery low | Battery charge below acceptable threshold |
| 5 | Physical intrusion | Cabinet or housing intrusion detected |

> Consult G.988 §9.1.1 and `knowledge/me-catalog/256-onu-g.json` (once available) for
> additional alarm bits.

---

## Circuit Pack (Class 6)

The Circuit Pack ME represents a plug-in or integrated line card within the ONU.
G.988 §9.1.6.

| Bit | Alarm name | Typical cause |
|---|---|---|
| 0 | Equipment alarm | Hardware fault on the circuit pack |
| 1 | Powering alarm | Power supply problem specific to the pack |
| 2 | Self-test failure | Circuit pack self-test did not pass |
| 3 | Laser end of life | Laser nearing end-of-life threshold |
| 4 | Temperature yellow | Temperature entering warning range |
| 5 | Temperature red | Temperature in critical range |

> Consult G.988 §9.1.6 and `knowledge/me-catalog/006-circuit-pack.json` (once available)
> for additional alarm bits and vendor-specific extensions.

---

## GEM Interworking Termination Point (Class 266)

The GEM Interworking TP ME connects a GEM port to an upper-layer service. G.988 §9.2.4.

| Bit | Alarm name | Typical cause |
|---|---|---|
| 6 | Deprecated (consult G.988) | — |

> The GEM Interworking TP alarm list is sparse and depends on the interworking type
> (IP, Ethernet, VEIP). Consult G.988 §9.2.4 and
> `knowledge/me-catalog/266-gem-interworking-tp.json` (once available).

---

## Placeholder MEs

The following MEs have alarm definitions in G.988 but detailed bit tables have not yet
been added to this document. When you encounter alarms from these MEs, consult G.988
directly and add a table here.

| ME | Class ID | G.988 clause |
|---|---|---|
| ONU2-G | 257 | §9.1.2 |
| T-CONT | 262 | §9.2.2 |
| GEM Port Network CTP | 268 | §9.2.3 |
| Multicast GEM Interworking TP | 281 | §9.2.8 |
| PPTP POTS UNI | 53 | §9.8.1 |
| VEIP | 329 | §9.14 |

---

## References

- ITU-T G.988 §9.1.1 — ONU-G alarms
- ITU-T G.988 §9.1.6 — Circuit Pack alarms
- ITU-T G.988 §9.2.1 — ANI-G alarms
- ITU-T G.988 §9.2.4 — GEM Interworking TP alarms
- ITU-T G.988 §9.5.1 — PPTP Ethernet UNI alarms
- [`alarm-message-format.md`](alarm-message-format.md) — bitmap decoding
- [`interpretation-workflow.md`](interpretation-workflow.md) — agent recipe

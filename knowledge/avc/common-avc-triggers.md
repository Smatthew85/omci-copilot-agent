# Common AVC Triggers

> **Note:** This list documents AVC-generating scenarios that are **defined in ITU-T G.988**.
> Do not add AVC triggers that are not grounded in the spec or in a vendor document
> explicitly cited in the row. Per-attribute AVC capability is authoritative in
> `knowledge/me-catalog/NNN-slug.json` (once populated) and in G.988.

---

## ANI-G (Class 263 / `0x0107`)

The ANI-G ME represents the ONU's physical PON interface (Analogue Narrowband Interface
Group in G.988 terminology; also called the PON TC adapter interface).

| Attribute index | Attribute name | Typical AVC trigger | G.988 reference |
|---|---|---|---|
| 9 | Received optical signal level | Optical Rx power changes crossing an ONU-internal reporting threshold | G.988 §9.2.1 |

> **Note:** Additional ANI-G attributes may be AVC-capable depending on the ONU
> implementation. Consult `knowledge/me-catalog/263-ani-g.json` (once available) or
> G.988 §9.2.1 for the full per-attribute AVC flag list.

---

## ONU-G (Class 256 / `0x0100`)

The ONU-G ME is the root ME representing the ONU as a whole.

| Attribute index | Attribute name | Typical AVC trigger | G.988 reference |
|---|---|---|---|
| 8 | Administrative state | OLT or ONU changes administrative state (lock/unlock) | G.988 §9.1.1 |
| 9 | Operational state | ONU transitions between enabled and disabled | G.988 §9.1.1 |

---

## PPTP Ethernet UNI (Class 11 / `0x000B`)

The PPTP Ethernet UNI ME represents a physical Ethernet UNI port.

| Attribute index | Attribute name | Typical AVC trigger | G.988 reference |
|---|---|---|---|
| 2 | Expected type | Auto-detection result updates the attribute | G.988 §9.5.1 |
| 8 | Operational state | Link up/down on the UNI port | G.988 §9.5.1 |
| 9 | Administrative state | Administrative state lock/unlock | G.988 §9.5.1 |

---

## Software Image (Class 7 / `0x0007`)

The Software Image ME tracks the state of each software image slot on the ONU. AVC
notifications on this ME are especially important during ONU software upgrades.

| Attribute index | Attribute name | Typical AVC trigger | G.988 reference |
|---|---|---|---|
| 3 | Is committed | Changes at end of a successful software activate/commit sequence | G.988 §9.1.4 |
| 4 | Is active | Changes when the ONU boots from a new image after End Software Download + Activate | G.988 §9.1.4 |
| 5 | Is valid | Changes when a new image has been downloaded and its CRC is verified | G.988 §9.1.4 |

> **Upgrade correlation:** During a software upgrade (Start/End Download → Activate →
> Reboot → Commit), the OLT should expect AVC notifications on attributes 5, 4, and 3 in
> that order. An AVC with `is-valid = 1` arriving without a prior software download
> sequence may indicate an unexpected ONU-initiated image validation.

---

## Summary Table

| ME | Class ID | Attribute | Typical trigger |
|---|---|---|---|
| ANI-G | 263 | 9 (Rx optical level) | Optical power change |
| ONU-G | 256 | 8 (Admin state) | Administrative lock/unlock |
| ONU-G | 256 | 9 (Operational state) | ONU enabled/disabled transition |
| PPTP Ethernet UNI | 11 | 2 (Expected type) | Auto-detection update |
| PPTP Ethernet UNI | 11 | 8 (Operational state) | UNI link up/down |
| PPTP Ethernet UNI | 11 | 9 (Admin state) | Administrative lock/unlock |
| Software Image | 7 | 3 (Is committed) | Image commit during upgrade |
| Software Image | 7 | 4 (Is active) | Image activated after reboot |
| Software Image | 7 | 5 (Is valid) | Image CRC validated after download |

---

## References

- ITU-T G.988 §9.1.1 — ONU-G
- ITU-T G.988 §9.1.4 — Software Image
- ITU-T G.988 §9.2.1 — ANI-G
- ITU-T G.988 §9.5.1 — PPTP Ethernet UNI
- [`avc-message-format.md`](avc-message-format.md) — AVC frame encoding
- [`interpretation-workflow.md`](interpretation-workflow.md) — decoding recipe

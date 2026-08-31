# Vendor-Specific OMCI Quirks

This directory is the authoritative place to capture vendor-specific deviations from
ITU-T G.988 observed in the field. It covers ONU vendors primarily; OLT quirks that
affect OMCI behavior can be recorded under an `olt/` subfolder per vendor as needed.

---

## Scope

- **In scope:** ME support gaps, firmware bugs, non-standard attribute behavior, and
  workarounds that have been confirmed through reproducible testing, a bug ticket, vendor
  errata, a VOLTHA issue, or a dated field observation with device details.
- **Out of scope:** Speculation, marketing claims, unverified reports, or behavior that
  is clearly a site-specific misconfiguration rather than a firmware/ONU issue.

## Contribution Rule

> **Do not add speculative or unverified quirks.**
>
> Every quirk entry must cite at least one source: a vendor bug ID, a VOLTHA issue
> link, an internal ticket number, a reproducible test case, or a dated field
> observation that includes the device model and firmware version.

To add a new quirk, follow the guide in [`_template/README.md`](./_template/README.md).

---

## Vendor Index

| Folder | Vendor | Typical product families |
|--------|--------|--------------------------|
| [`adtran/`](./adtran/) | Adtran | Total Access, SDX, 4-series ONUs |
| [`nokia/`](./nokia/) | Nokia | G-010G, XS-010X, Fastmile, G-1426G ONUs |
| [`calix/`](./calix/) | Calix | GigaSpire, 700-series, 800-series ONUs |
| [`huawei/`](./huawei/) | Huawei | EchoLife HG, MA5671/5672, ONT series |
| [`zte/`](./zte/) | ZTE | F660, F670, ZXHN series ONUs |
| [`_other/`](./_other/) | Other / Unknown | Catch-all for vendors not listed above |

---

## How the Copilot Agent Uses This Folder

When a user mentions a specific ONU vendor or model, the agent consults the relevant
vendor subfolder for known deviations **before** applying generic G.988 rules. If a
matching quirk file is found, the agent cites it and prefers its diagnosis. If no
quirk matches, the agent proceeds with standard G.988 diagnosis and notes that the
vendor folder was checked.

See [`.github/copilot-instructions.md`](../../.github/copilot-instructions.md) for the
full agent instructions.

---

## See Also

- [`../README.md`](../README.md) — top-level knowledge base index
- [`_template/quirk.md`](./_template/quirk.md) — contribution template

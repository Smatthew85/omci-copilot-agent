# <Vendor> <Model/Family> — <one-line quirk summary>

- **Vendor:**
- **Model family:**
- **Firmware versions affected:**
- **First observed:** YYYY-MM-DD
- **Last verified:** YYYY-MM-DD
- **Source(s):** (vendor bug ID, VOLTHA issue link, internal ticket, field capture — at least one required)

## ME(s) involved

- Class ID / name (link to `../../me-catalog/NNN-slug.json` once catalog exists)

## Standard G.988 behavior

What the spec says should happen.

## Observed behavior

What actually happens on this vendor/firmware. Include hex frames or log excerpts where possible.

## Impact

What breaks for the operator (e.g., "VLAN filtering silently drops packets", "MIB sync loop",
"Create returns success but ME not usable").

## Workaround

Concrete steps the operator/adapter can take. If none exists, say so.

## Detection

How the agent can recognize this quirk from a user's paste (specific hex pattern, log signature,
result code + ME combination).

## References

Bug tickets, VOLTHA PRs, vendor docs, dated field notes.

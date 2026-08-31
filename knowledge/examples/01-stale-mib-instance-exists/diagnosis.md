# Diagnosis: Create MAC Bridge Service Profile returns 0x07 Instance Exists

## Summary

The OLT issued a Create request for MAC Bridge Service Profile (ME Class 45) instance 1.
The ONU responded with result code **0x07 — Instance Exists**, indicating the instance
was already present in the ONU MIB. The OLT's view of ONU state was out of sync, most
likely because a prior provisioning attempt left the instance behind and no MIB Reset was
performed before the retry.

## Root Cause

The ONU MIB retains a MAC Bridge Service Profile instance from a previous provisioning
cycle. Because the OLT did not perform a MIB Reset (MT 15) after the ONU reboot, its
local view of the MIB did not reflect the stale instance on the ONU side.

Per G.988, a Create action **must** fail with result code 0x07 when an instance with the
requested ME Class ID and Instance ID already exists. The ONU is behaving correctly; the
problem is the OLT–ONU state divergence.

## Evidence

- **Message Type**: MT=10 — Create (request byte `0x4A`, response byte `0x2A`)
- **ME Class**: 45 — MAC Bridge Service Profile
- **ME Instance**: 1
- **Result Code**: `0x07` — Instance Exists (response content byte 0)
- The Create request carries plausible SetByCreate attribute values; the ONU rejects the
  operation before processing any attribute, so attribute correctness is not the issue.

See also:
- [`knowledge/result-codes/README.md`](../../result-codes/README.md) — definition of 0x07
- [`knowledge/failure-patterns/README.md`](../../failure-patterns/README.md) — "Create returns Instance Exists" pattern

## Remediation

1. **MIB Reset** — send MT=15 (MIB Reset) to the ONU. This clears all dynamically
   provisioned MEs and returns the ONU MIB to its power-on state.
2. **MIB Upload** — send MT=13 (MIB Upload) to retrieve the ONU MIB snapshot count,
   followed by repeated MT=14 (MIB Upload Next) messages to enumerate all MIB entries
   and rebuild the OLT's local copy.
3. **Reissue Create** — once the OLT's MIB view is consistent with the ONU, reissue
   the Create for ME 45 instance 1 (and any other service MEs needed).

If the stale instance is intentional (e.g., the ONU survived across a soft restart),
consider issuing a Delete (MT=11) for the instance before re-creating it rather than
performing a full MIB Reset.

## References

- [`knowledge/provisioning-flows/`](../../provisioning-flows/) — standard ONU provisioning
  sequence including MIB Reset / MIB Upload steps
- [`knowledge/result-codes/README.md`](../../result-codes/README.md)
- [`knowledge/failure-patterns/README.md`](../../failure-patterns/README.md)
- ITU-T G.988, clause 9.3.3 — MAC Bridge Service Profile ME definition
- ITU-T G.988, clause 11.2.2 — MIB Reset and synchronization procedure

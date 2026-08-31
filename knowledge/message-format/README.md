# OMCI Message Format Reference

This directory contains reference documentation for ITU-T G.988 OMCI frame formats. It enables the Copilot agent to decode raw hex OMCI frames without external tools.

## Contents

| File | Description |
|---|---|
| [`baseline-frame.md`](baseline-frame.md) | 48-byte baseline OMCI frame layout — field offsets, the Message Type byte breakdown, and a worked decoding example |
| [`extended-frame.md`](extended-frame.md) | Extended OMCI frame layout — variable-length payload, when it is used, structural differences from baseline, and a worked example |
| [`message-types.md`](message-types.md) | Complete Message Type (MT) reference table — all standardised MT values, direction, AR/AK semantics |

## Purpose

Raw OMCI frames appear in OLT/ONU debug logs, VOLTHA traces, and packet captures as hex strings. These documents allow the Copilot agent to:

1. Identify the frame format (baseline vs. extended) from the Device Identifier byte.
2. Extract TCID, Message Type, ME Class, and ME Instance without running external tools.
3. Map the MT value to a human-readable operation name and determine whether a response is expected.
4. Decode the Message Contents field with the correct field layout for the given MT.

See also:

- [`../result-codes/`](../result-codes/) — G.988 result and reason codes
- [`../provisioning-flows/`](../provisioning-flows/) — standard ONU provisioning sequences
- [`../failure-patterns/`](../failure-patterns/) — symptom → root-cause diagnostic table

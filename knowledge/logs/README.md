# OMCI Log Format Reference

This directory documents the log formats where OMCI frames appear in VOLTHA and
OLT hardware environments. It enables the Copilot agent to accept a raw log paste
(not pre-cleaned hex) and still produce a full diagnosis by extracting, grouping,
and decoding the embedded OMCI frames.

---

## Contents

| File | Description |
|---|---|
| [`voltha-openonu-adapter.md`](voltha-openonu-adapter.md) | `openonu-adapter-go` structured JSON log format — the primary source of OMCI traffic on the ONU side |
| [`voltha-openolt-adapter.md`](voltha-openolt-adapter.md) | `openolt-adapter` structured JSON log format — OMCI proxy/forwarding and ONU activation events |
| [`bal-openolt-agent.md`](bal-openolt-agent.md) | OLT hardware-side logs (Broadcom BAL / OpenOLT agent) — text-based, syslog-style |
| [`extraction-workflow.md`](extraction-workflow.md) | End-to-end recipe: paste log → find OMCI lines → extract hex → hand to `knowledge/message-format/` |

---

## Purpose

Users of this agent will most often paste log output from VOLTHA or OLT hardware
rather than raw hex frames. These documents tell the agent how to:

1. Recognize which component produced the log (openonu-adapter, openolt-adapter, BAL).
2. Identify which log lines carry OMCI frames.
3. Extract the raw hex from the surrounding JSON or text context.
4. Group frames by device and order them chronologically.
5. Correlate TX/RX pairs by Transaction Correlation ID (TCID).
6. Hand the cleaned hex off to [`knowledge/message-format/baseline-frame.md`](../message-format/baseline-frame.md)
   or [`extended-frame.md`](../message-format/extended-frame.md) for decoding.

---

## Quick Reference — Log Source Detection

| Characteristic | Log Source |
|---|---|
| Structured JSON with `"logger"` field | VOLTHA (openonu-adapter or openolt-adapter) |
| `"logger"` value contains `omci-cc`, `MibDownloadFsm`, `UniVlanConfigFsm` | `openonu-adapter-go` |
| `"logger"` value contains `openolt`, `flowMgr`, `core-proxy` | `openolt-adapter` |
| Text lines with syslog-style timestamp, `OMCI:` or `PON OMCI Msg:` prefix | BAL / OpenOLT agent on OLT hardware |

If the format is ambiguous, ask the user to confirm the log source before extracting.

# openolt-adapter Log Format

`openolt-adapter` sits between the VOLTHA core and the OLT hardware. It proxies
OMCI frames between the core (which drives `openonu-adapter-go`) and the physical
OLT, and logs ONU discovery/activation events that precede any OMCI exchange.

---

## Log Framework

Same as `openonu-adapter-go`: structured JSON via `voltha-lib-go` / zap.
The mandatory fields (`level`, `ts`, `logger`, `caller`, `msg`) and common
contextual fields (`device-id`, `intf-id`, `onu-id`) are identical — see
[`voltha-openonu-adapter.md`](voltha-openonu-adapter.md) for field definitions.

---

## Key Loggers

| `logger` Value | Subsystem |
|---|---|
| `openolt` | Main adapter entry point |
| `core-proxy` | Proxy between VOLTHA core and adapter |
| `adapter-proxy` | Proxy between adapters |
| `flowMgr` | Flow manager — translates VOLTHA flows to OLT BAL calls |
| `resourceMgr` | Resource manager (ALLOC-ID, GEM-ID allocation) |

---

## How OMCI Appears in openolt-adapter Logs

The openolt-adapter does **not** originate OMCI messages — it forwards them.
OMCI bytes appear in two contexts:

### 1. OMCI Indication (ONU → OLT → core)

Lines with `msg="omci-indication"` or similar carry a received OMCI frame from
the OLT hardware up to the VOLTHA core.

```json
{"level":"debug","ts":1693000010.460,"logger":"openolt","caller":"omci.go:88","msg":"omci-indication","device-id":"0001010000000001","intf-id":0,"onu-id":1,"pkt":"000129094f0000000000000000000000000000000000000000000000000000000000000000000000000000000000028a"}
```

- `pkt` field holds the hex frame (strip surrounding quotes)
- Direction: ONU→OLT (inbound indication)

### 2. Proxy OMCI Message (core → OLT → ONU)

Lines with `msg="proxy-omci-message"` or `msg="send-omci-message"` carry a frame
being forwarded from the core to the OLT hardware.

```json
{"level":"debug","ts":1693000010.120,"logger":"openolt","caller":"omci.go:55","msg":"proxy-omci-message","device-id":"0001010000000001","intf-id":0,"onu-id":1,"omci-message":"000100094f0000000000000000000000000000000000000000000000000000000000000000000000000000000000028a"}
```

- `omci-message` field holds the hex frame
- Direction: OLT→ONU (outbound proxy)

### 3. ONU Discovery / Activation Events

These lines do **not** carry OMCI frames but are important context — they mark
the point before OMCI provisioning begins.

```json
{"level":"info","ts":1693000005.001,"logger":"openolt","caller":"device_handler.go:412","msg":"onu-discovery","device-id":"0001010000000001","intf-id":0,"onu-id":1,"serial-number":"HGAC1234ABCD"}
{"level":"info","ts":1693000006.200,"logger":"openolt","caller":"device_handler.go:550","msg":"onu-activation-completed","device-id":"0001010000000001","intf-id":0,"onu-id":1}
```

- `serial-number` is a placeholder; real values are ONU-vendor-specific.
- Timestamps show that OMCI provisioning starts only after `onu-activation-completed`.

### 4. Flow Provisioning Attempts

Flow add/remove events trigger downstream OMCI traffic. They appear as:

```json
{"level":"info","ts":1693000030.010,"logger":"flowMgr","caller":"flow_mgr.go:780","msg":"adding-flow","device-id":"0001010000000001","intf-id":0,"onu-id":1,"flow-id":1001,"flow-type":"unicast"}
```

---

## Field Extraction Rules

Apply the same extraction recipe as `openonu-adapter-go`
([`voltha-openonu-adapter.md` — How to Extract](voltha-openonu-adapter.md#how-to-extract-verbatim-recipe)),
with these differences:

| Step | openolt-adapter specifics |
|---|---|
| Logger filter | Use `openolt`, `core-proxy`, `adapter-proxy`, `flowMgr` |
| Hex field priority | `pkt` → `omci-message` → `packet` |
| Direction | `msg` containing `indication` → ONU→OLT; `proxy` or `send` → OLT→ONU |

---

## When Only openolt-adapter Logs Are Available

The openolt-adapter log shows the full byte stream but lacks the semantic context
(FSM state, ME name, attribute values) that `openonu-adapter-go` provides.

Recommended approach:

1. Extract frames from the openolt-adapter log using the rules above.
2. Decode each frame against [`knowledge/message-format/baseline-frame.md`](../message-format/baseline-frame.md).
3. Note which frames are requests (AR bit set) and which are responses (AK bit set).
4. Look up result codes in [`knowledge/result-codes/README.md`](../result-codes/README.md).
5. If the openonu-adapter log is also available, cross-reference by `device-id` and
   TCID to add FSM-level context.

> **Limitation:** The openolt-adapter log may not include OMCI frames if the OLT
> hardware suppresses them at the driver level. In that case, the openonu-adapter
> log is the more reliable source.

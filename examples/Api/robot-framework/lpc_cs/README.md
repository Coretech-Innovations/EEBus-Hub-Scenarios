# LPC Controllable-System examples

These suites drive the **Limitation of Power Consumption (LPC)** use case with **both
actors simulated in one hub**: the hub's built-in HEMS and a simulated EVSE acting as
the controllable system (CS). They communicate over the hub's own EEBUS (SHIP/SPINE)
stack and are driven through the Robot Framework Remote interface.

| File | Scenario |
|------|----------|
| `cs_setup.resource` | Shared setup keywords: reset the hub, create and pair the HEMS and the CS, and run the heartbeat handshake. |
| `heartbeat_gated_limit.robot` | A load-control limit is rejected before any heartbeat is received, then accepted once the CS has received one. |
| `limit_validation.robot` | A valid active limit is accepted (CS → Limited); a negative limit value is rejected and the state is unchanged. |

## Running

Prerequisites are the same as the other examples (a running hub with the Remote server
and a `robotFrameworkSupported` license — see `../README.md`).

```bash
python -m robot examples/robot-framework/lpc_cs/
```

**Timing note.** These are integration tests: two simulated devices pair over the hub's
SHIP/SPINE stack, which takes a few seconds (the first pairing after a hub start also
pays an mDNS/cert warmup). The setup keywords wait generously for the binding and the
first heartbeat. Each suite takes ~25 s; a full run ~50 s. Prefer running against a
freshly started hub.

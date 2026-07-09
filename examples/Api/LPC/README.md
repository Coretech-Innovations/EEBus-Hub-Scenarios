# LPC — Limitation of Power Consumption

## Description

A HEMS and an EVSE are created and connected. Once the connection is established, the
example reads the EVSE's initial limit state and then exercises a full limit cycle:

1. **Send an active limit** — the HEMS sends an active power consumption limit with a
   value and a duration.
2. **Observe the EVSE** — the EVSE applies the limit and reports the new active power
   limit and a `Limited` state.
3. **Watch it expire** — once the duration elapses, the EVSE releases the limit and
   returns to its `Unlimited/Controlled` state.

Throughout the run the example prints the EVSE's LPC state, active power limit, failsafe
values, and power constraints so you can follow each transition.

## How to Run

```bash
go run ./examples/Api/LPC/LPC.go
```

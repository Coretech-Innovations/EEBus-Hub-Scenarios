# Single EV at Full Current

## Description

A HEMS and an EVSE are created and connected, then an EV with a current range of
6–10 A per phase is added and plugged into the EVSE.

Because the site's available current (a 40 A fuse budget) far exceeds the EV's demand
and nothing else is drawing power, the HEMS lets the EV charge at its maximum:
**10 A on every phase**.

## How to Run

```bash
go run ./examples/Api/OPEV/single-ev-full-current
```

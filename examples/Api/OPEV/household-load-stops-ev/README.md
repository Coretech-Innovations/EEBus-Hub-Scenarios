# Household Load Curtails EV Charging

## Description

A HEMS and an EVSE are created and connected, then an EV with a current range of
6–10 A per phase is added. With a 40 A fuse budget and no other load, the EV would
charge at its maximum of 10 A per phase.

An uncontrollable household appliance is then added that consumes the entire available
current budget. To keep the site within its limit, the HEMS curtails the EV's charging
current down to **0 A**.

## How to Run

```bash
go run ./examples/Api/OPEV/household-load-stops-ev
```

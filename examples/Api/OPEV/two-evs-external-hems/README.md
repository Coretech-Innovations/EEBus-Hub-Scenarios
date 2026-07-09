# Two EVs Under an External HEMS

## Description

Two EVs are created, each plugged into its own EVSE, and both are connected to an
**external** HEMS identified by its SKI, which you pass on the command line.

Each EV receives its charging current limits from the external HEMS based on the current
the HEMS makes available.

## How to Run

Pass the external HEMS's SKI as an argument:

```bash
go run ./examples/Api/OPEV/two-evs-external-hems <remoteSKI>
```

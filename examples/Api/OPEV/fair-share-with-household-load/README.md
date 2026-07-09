# Fair-Share With a Household Load

## Description

A HEMS and an EVSE are created and connected, then two EVs are added:

- EV 1 — current range 6–10 A
- EV 2 — current range 10–20 A

An uncontrollable household load is added that draws 24 A of the 40 A fuse budget,
leaving 16 A available for charging. The HEMS applies a fair-share distribution of the
remaining current across the two EVs.

## How to Run

```bash
go run ./examples/Api/OPEV/fair-share-with-household-load
```

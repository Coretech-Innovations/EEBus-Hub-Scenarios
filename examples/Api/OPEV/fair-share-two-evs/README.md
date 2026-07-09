# Fair-Share Between Two EVs

## Description

A HEMS and an EVSE are created and connected, then two EVs are added:

- EV 1 — current range 6–18 A
- EV 2 — current range 10–30 A

Their combined maximum demand exceeds the site's available current, so the HEMS applies
a fair-share distribution, splitting the available current between the two EVs while
respecting each one's minimum and maximum.

## How to Run

```bash
go run ./examples/Api/OPEV/fair-share-two-evs
```

# Receiving Limits From an External HEMS

## Description

A simulated EV is created and connected to an **external** HEMS — a separate device
(real, or simulated in another hub) identified by its SKI, which you pass on the command
line.

The EV receives its charging current limits from the external HEMS according to the
current the HEMS makes available.

## How to Run

Pass the external HEMS's SKI as an argument:

```bash
go run ./examples/Api/OPEV/external-hems <remoteSKI>
```

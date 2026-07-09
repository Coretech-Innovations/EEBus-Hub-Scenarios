# Limiting an External EV

## Description

A HEMS is created and connected to an **external** EV — a separate device (real, or
simulated in another hub) identified by its SKI, which you pass on the command line.

The HEMS supplies the external EV with its requested current as long as it stays within
the site's available current (a 40 A budget in this example).

## How to Run

Pass the external EV's SKI as an argument:

```bash
go run ./examples/Api/OPEV/external-ev <remoteSKI>
```

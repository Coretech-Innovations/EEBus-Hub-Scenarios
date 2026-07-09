# Grid Connection Point Monitoring

## Description

This example demonstrates the **Monitoring of Grid Connection Point (MGCP)** use case.
A HEMS is connected to a smart meter at the building's grid connection point, and a
number of uncontrollable loads are added to the simulation.

The smart meter reports the grid connection point measurements — voltage, current,
power, and total consumed energy — which the HEMS reads and prints periodically.

The grid connection point is provided by an external device (real, or simulated in
another hub) identified by its SKI, which you pass on the command line.

## How to Run

Pass the grid connection point's SKI as an argument:

```bash
go run ./examples/Api/MGCP/grid-point-monitoring <gcpRemoteSki>
```

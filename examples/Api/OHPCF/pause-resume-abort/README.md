# Heat Pump: Pause, Resume, and Abort

## Description

A HEMS and a heat pump are created and connected. The HEMS then issues an OHPCF
(Optimization of Self-Consumption by Heat Pump Compressor Flexibility) announcement
describing a flexible operating window:

- power (good approximation): 15 000 W
- maximum power: 16 000 W
- minimum active duration: 5 s
- minimum pause duration: 5 s
- start time: 5 s

The HEMS sends the start time and then drives the heat pump through pause, resume, and
abort transitions once the minimum active and pause durations have elapsed. The example
waits for the heat pump to return to the inactive state after the abort.

## How to Run

```bash
go run ./examples/Api/OHPCF/pause-resume-abort
```

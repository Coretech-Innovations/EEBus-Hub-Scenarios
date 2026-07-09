# OPEV — Overload Protection by EV Current Curtailment (Web UI)

## Description

This example demonstrates the **Overload Protection by EV Current Curtailment (OPEV)**
use case using the EEBUS Hub web interface instead of the REST API.

A HEMS and an EVSE are created and connected, then an EV with a current range of 6–10 A
per phase is added and plugged into the EVSE.

Because the site's available current (a 40 A fuse budget) far exceeds the EV's demand and
nothing else is drawing power, the HEMS lets the EV charge at its maximum: **10 A on
every phase**.

## How to Run

Import the provided configuration file (`OPEV.json`) into the simulation from the web
interface, then start the simulation.

![OPEV simulation in the web interface](<Screenshot 2024-09-12 111917.png>)

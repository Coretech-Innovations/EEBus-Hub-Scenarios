# LPC — Limitation of Power Consumption (Web UI)

## Description

This example walks through the **Limitation of Power Consumption (LPC)** use case using
the EEBUS Hub web interface instead of the REST API.

A HEMS and an EVSE are created and connected. Once they are paired, the HEMS sends an
active power consumption limit with a value and a duration. The EVSE applies the limit
and reports a `Limited` state; when the duration elapses it releases the limit and
returns to its `Unlimited/Controlled` state.

## How to Run

Create the HEMS and the EVSE from the web interface, connect them, then send an active
power consumption limit and watch the EVSE's state and active power limit update.

The screenshots below show the sequence in order:

![Step 1](<Screenshot 2024-10-31 124518-1st.png>)

![Step 2](<Screenshot 2024-10-31 124757-2nd.png>)

![Step 3](<Screenshot 2024-10-31 124644-3rd.png>)

![Step 4](<Screenshot 2024-10-31 124843-4th.png>)

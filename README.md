# EEBus-Hub
EEBUS Hub is a testing and simulation framework for validating device integration over an EEBUS network. It builds on the [eebus-go](https://github.com/enbility/eebus-go) library for core EEBUS interactions and adds practical tooling to streamline real-world testing.

EEBUS Hub exposes APIs to control and orchestrate multiple EEBUS actors that participate in typical energy scenarios such as EV, EVSE, HEMS, Controlbox, Heatpump, Inverters and SMGW. You can run full simulations using built-in virtual devices, or plug in real hardware alongside simulated participants to accelerate integration, troubleshooting, and regression testing.

This repository includes examples that show how to use EEBUS Hub programmatically to automate your EEBUS testing procedures—ideal for CI pipelines and repeatable test suites.

Compliance support: EEBUS Hub helps manufacturers validate S14a compliance (and beyond) by enabling controlled, reproducible test scenarios and automated verification workflows.

Licensing: To obtain an EEBUS Hub license, contact us at eebus.hub@coretech-innovations.com

For researchers, academics, and early-stage startups you can obtain a non-commercial licence by joining our Coretech EEBUS Innovation Program. Contact us for more details.


## Supported Use Cases

| Use Case | Category | Scenarios | Server | Client |
|:---------|:---------|:---------:|:------:|:------:|
| Limitation of Power Consumption (LPC) | Grid | 1–4 | ✅ | ✅ |
| Limitation of Power Production (LPP) | Grid | 1–4 | ✅ | ✅ |
| Monitoring of Power Consumption (MPC) | Grid | 1–5 | ✅ | ✅ |
| Monitoring of Grid Connection Point (MGCP) | Grid | 1–7 | — | — |
| EV Commissioning and Configuration (EVCC) | E-mobility | 1,2,3,6,7,8 | ✅ | ✅ |
| EVSE Commissioning and Configuration (EVSECC) | E-mobility | 1–2 | ✅ | ✅ |
| Overload Protection by EV Current Curtailment (OPEV) | E-mobility | 1–3 | ✅ | ✅ |
| EV State of Charge (EVSOC) | E-mobility | 1–4 | ✅ | ✅ |
| EV Charging Electricity Measurement (EVCEM) | E-mobility | 1 | ✅ | ✅ |
| Optimization of Self Consumption During EV Charging (OSCEV) | E-mobility | 1,2,3,4,5,6 | ✅ | ✅ |
| Coordinated EV Charging (CEVC) | E-mobility | — | — | — |
| Monitoring of Inverter (MOI) | Inverter | 1–7 | ✅ | ✅ |
| Monitoring of Battery (MOB) | Inverter | 1–9 | ✅ | ✅ |
| Monitoring of PV String (MPS) | Inverter | – | - | - |
| Control of Battery (COB) | Inverter | 1-5 | ✅ | ✅ |
| Optimization of Self-Consumption by Heat Pump Compressor Flexibility (OHPCF) | HVAC | 1–2 | ✅ | ✅ |
| Incentive Table based Power Consumption Management (ITPCM) | HVAC | – | - | - |
| Node Identification (NID) | Generic | — | — | — |

✅ Supported &nbsp;&nbsp; — Not yet supported

## Clone the project

```bash
git clone https://github.com/Coretech-Innovations/EEBus-Hub.git

```

## How to Use

This is a simple API calls for adding EVSE and EV and connecting them with the CEM in the system

```bash
# adding new EVSE
curl -X POST http://localhost:8080/api/v1/evse/add  -H 'Content-Type: application/json'  -d '{"deviceName":"Coretech EVSE WLBX", "deviceCode":"0002","deviceModel":"Charging Station","brandName":"Coretech Innovations","vendor":{"name":"Coretech Innovations","code":"60745"},"serialNumber":"SN7640"}'

# Trusting the created EVSE from the CEM side
curl -X POST http://localhost:8080/api/v1/cem/trust -H 'Content-Type: application/json' -d '{"remoteSki": <EVSE Ski>}'

# Get the SKI of the CEM
curl -X GET http://localhost:8080/api/v1/cem 

# Pairing the EVSE with the CEM
curl -X POST http://localhost:8080/api/v1/evse/<EVSE ID>/cem -H 'Content-Type: application/json' -d '{"remoteski": <CEM Ski>}'

# adding new EV
curl -X POST http://localhost:8080/api/v1/ev/add -H 'Content-Type: application/json' -d '{"device": {"name": "Taycan", "code": "0003", "serialNumber": "SN1235"},"currentLimits": {"min": 5, "max": 10}, "asymmetricCharging": false}'

# adding the created EV to the EVSE we created before
curl -X POST http://localhost:8080/api/v1/ev/<EV ID>/evse/<EVSE ID>

# starting the simulation
curl -X POST http://localhost:8080/api/v1/sim -H 'Content-Type: application/json' -d '{"action": "start","tickRate": 1000,"simTimePerTick": 10}'
```

EVSE Ski: The Ski returned from creating the EVSE  
EVSE ID: The id returned from creating the EVSE  
CEM Ski: The Ski returned from calling the GET request on the CEM  
EV ID: The id returned from creating the EV


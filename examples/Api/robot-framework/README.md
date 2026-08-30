# EEBUS Hub — Robot Framework examples

Runnable example test suites showing how to drive the EEBUS Hub from
[Robot Framework](https://robotframework.org/) through its **Remote Library**
(served live from the Hub binary — no keyword source ships).

These files are meant to be shared with test authors as a starting point. Copy them,
adapt the entity fields, and build your own suites.

| File | What it demonstrates |
|------|----------------------|
| `common.resource` | Shared setup: imports the Hub Remote Library, defines the connection variables, and two thin helpers (`Should Be Successful`, `Reset Simulation To A Clean State`). |
| `01_evse_lifecycle.robot` | CRUD on a single entity: no-arg create with defaults, overriding only the fields you care about (friendly aliases + type coercion), path parameters, and reading the result dict. |
| `02_ev_evse_session.robot` | A multi-entity scenario: two entity types, a keyword with two path parameters (`Connect EV To EVSE`), polling an async pairing with `Wait Until Keyword Succeeds`, and asserting an expected failure (`ok=False`). |

## Prerequisites

1. **A running Hub** with a license that has `robotFrameworkSupported: true`. The Remote
   server is enabled by default:
   ```bash
   eebus-hub                          # serves the Remote Library on 0.0.0.0:8270
   ```
   Options: `--no-robotremote` to disable it, `--robotremote-port <port>` (default `8270`).
   The Remote Library binds to the hub's own `bindaddress` argument (default `0.0.0.0`, env:
   `BINDADDRESS`) — the same address the REST API uses, so the two never drift apart. The
   default lets you run the tests from a different machine; start the hub as
   `eebus-hub 8080 127.0.0.1` to restrict both to the hub's own host.

2. **Robot Framework** (in a virtualenv — no other dependency is needed; every Hub
   keyword is served by the binary):
   ```bash
   python3 -m venv .venv
   .venv/bin/pip install robotframework
   ```

## Run

```bash
.venv/bin/python -m robot examples/robot-framework/
```

Point at a non-default host/port by overriding the variables:

```bash
.venv/bin/python -m robot \
    --variable HUB_HOST:192.168.1.50 --variable HUB_RF_PORT:8270 \
    examples/robot-framework/
```

## The keyword contract (what every call returns)

Each Hub keyword returns a dictionary:

| Key | Meaning |
|-----|---------|
| `status_code` | HTTP status (e.g. `200`, `404`). |
| `json` | Parsed response body (dict / list), or `None` if not JSON. |
| `text` | Raw response body as a string. |
| `ok` | `True` for HTTP 2xx. |

A 4xx/5xx does **not** raise on its own — it comes back with `ok=False`. Assert on it
explicitly (`Should Be Successful` in `common.resource`, or check `status_code` directly
as in `02_ev_evse_session.robot`).

Request-body fields are passed as `alias=value` and default to "omit": send only what
you want to set. `Add *` keywords carry a default body so a no-arg create works. Nested
JSON fields have friendly aliases (e.g. `nominalPowerMax` → `nominalPower.max`), and a
value given as text is coerced to the schema type (`"11000"` and `${11000}` both work).

## Browse all keywords

Generate the reference docs (the Swagger-UI equivalent) from a running Hub:

```bash
.venv/bin/python -m robot.libdoc Remote::http://localhost:8270 rf-keywords.html
```

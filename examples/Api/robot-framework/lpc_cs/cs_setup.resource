*** Settings ***
Documentation     Self-contained setup for LPC controllable-system (CS) scenarios where
...    BOTH actors are simulated inside the same hub: the hub's built-in HEMS and a
...    simulated EVSE acting as the controllable system.
...
...    This resource is independent of the other examples — it imports the hub's Remote
...    Library directly and carries its own small helpers, so the lpc_cs suites can be
...    copied and run on their own.

Library          Remote    http://${HUB_HOST}:${HUB_RF_PORT}/


*** Variables ***
${HUB_HOST}       localhost
${HUB_RF_PORT}    8270


*** Keywords ***
Should Be Successful
    [Documentation]    Fail unless the hub returned a 2xx status. Every hub keyword returns
    ...    a dict: status_code / json / text / ok (true for HTTP 2xx).
    [Arguments]    ${response}
    Should Be True    ${response}[ok]
    ...    Hub call failed: HTTP ${response}[status_code] -> ${response}[text]

Reset Simulation To A Clean State
    [Documentation]    Remove all simulated entities so a suite starts fresh.
    ${r}=    Reset Simulation
    Should Be Successful    ${r}

Setup Env Init
    [Documentation]    Reset the hub, enable the LPC use case and create the controllable
    ...    system (a simulated EVSE); the HEMS already exists after a reset. The HEMS
    ...    heartbeat is stopped so the CS starts having NEVER received a heartbeat — the
    ...    equivalent of a controller that has not begun sending heartbeats yet. Returns the
    ...    HEMS SKI and the new EVSE id.
    Reset Simulation To A Clean State

    ${r}=    Update Use Cases Availability    lpcEnabled=${True}
    Should Be Successful    ${r}

    ${hems}=    Get HEMS Information
    Should Be Successful    ${hems}
    ${hems_ski}=    Set Variable    ${hems}[json][ski]

    ${evse}=    Add EVSE    lpcEnabled=${True}
    Should Be Successful    ${evse}
    ${evse_id}=    Set Variable    ${evse}[json][id]

    Stop The HEMS Heartbeat
    RETURN    ${hems_ski}    ${evse_id}

Connect CS To HEMS
    [Documentation]    Start the simulation, establish mutual SHIP trust, and wait until the
    ...    connection is fully established — i.e. the HEMS can read the CS's load-control
    ...    data. Returns the CS device address.
    ...
    ...    Waiting for the connection here (rather than retrying the first write) keeps the
    ...    two failure modes separate: if pairing does not complete, THIS keyword fails with
    ...    a connection error; a later write can then only fail on the CS's own decision.
    [Arguments]    ${evse_id}    ${hems_ski}
    ${r}=    Start Simulation
    Should Be Successful    ${r}
    Sleep    5s    reason=let the simulation come up before trusting

    # EVSE trusts the HEMS; the response carries the EVSE's own SKI.
    ${trust}=    Trust New Device On EVSE    ${evse_id}    remoteSKI=${hems_ski}
    Should Be Successful    ${trust}
    ${evse_ski}=    Set Variable    ${trust}[json][ski]

    # HEMS trusts the EVSE back.
    ${r}=    Trust New Device On HEMS    remoteSki=${evse_ski}
    Should Be Successful    ${r}

    # SHIP/SPINE pairing over loopback takes a few seconds (plus first-run mDNS/cert warmup),
    # so allow a generous window.
    Wait Until Keyword Succeeds    90x    2s    The HEMS Can Read The CS Load Control
    ${address}=    Get The Controllable System Address
    RETURN    ${address}

The HEMS Can Read The CS Load Control
    [Documentation]    Succeed once the HEMS can read the connected CS's active-power-limit
    ...    data. This read exercises the same descriptor path as a write, so once it works
    ...    the connection is established and a subsequent write can only be rejected by the
    ...    CS's own logic, not by an incomplete binding.
    ${r}=    Get Connected Devices Active Power Consumption Limit Data By HEMS
    Should Be Successful    ${r}
    Should Not Be Empty    ${r}[json]    The HEMS cannot read the CS load-control data yet (still binding)

Get The Controllable System Address
    [Documentation]    Return the device address of the first LPC device known to the HEMS.
    ...    (LpcDevicesData is the reliable source for the addressable device string.)
    ${r}=    Get Connected Devices LPC Data By HEMS
    Should Be Successful    ${r}
    Should Not Be Empty    ${r}[json]    No LPC device is bound to the HEMS
    RETURN    ${r}[json][0][deviceAddress]

Send HEMS Heartbeat
    [Documentation]    Start the HEMS heartbeat, wait until the CS has received one, then
    ...    STOP it again. Afterwards the CS has established control (received a heartbeat)
    ...    even though the heartbeat is no longer running — which is exactly what the next
    ...    write is evaluated against.
    [Arguments]    ${evse_id}
    Start The HEMS Heartbeat
    Wait For The CS To Receive A Heartbeat    ${evse_id}
    Stop The HEMS Heartbeat
    Sleep    2s    reason=let the CS settle after the heartbeat stops

Wait For The CS To Receive A Heartbeat
    [Documentation]    Block until the EVSE reports at least one received heartbeat.
    [Arguments]    ${evse_id}
    Wait Until Keyword Succeeds    45x    2s    CS Heartbeat Counter Is Present    ${evse_id}

CS Heartbeat Counter Is Present
    [Documentation]    Succeed once the EVSE's latest-received-heartbeat carries a counter.
    [Arguments]    ${evse_id}
    ${r}=    Get EVSE Latest Received Heartbeat    ${evse_id}
    Should Be Successful    ${r}
    Should Be True    ${r}[json][result][counter] >= 0    No heartbeat received by the CS yet

Stop The HEMS Heartbeat
    [Documentation]    Stop the HEMS heartbeat and confirm it is stopped.
    ${r}=    Update HEMS Heartbeat Status    action=stop
    Should Be Successful    ${r}
    ${hb}=    Get HEMS Heartbeat
    Should Be Equal    ${hb}[json][state]    Stopped

Start The HEMS Heartbeat
    [Documentation]    Start the HEMS heartbeat and confirm it is running.
    ${r}=    Update HEMS Heartbeat Status    action=start
    Should Be Successful    ${r}
    ${hb}=    Get HEMS Heartbeat
    Should Be Equal    ${hb}[json][state]    Running

CS LPC State Should Be
    [Documentation]    Assert the EVSE's LPC state machine is in the expected state.
    [Arguments]    ${evse_id}    ${expected}
    ${r}=    Get EVSE LPC State    ${evse_id}
    Should Be Successful    ${r}
    Should Be Equal    ${r}[json][currentState]    ${expected}

CS LPC State Should Not Be
    [Documentation]    Assert the EVSE's LPC state machine is NOT in the given state.
    [Arguments]    ${evse_id}    ${unexpected}
    ${r}=    Get EVSE LPC State    ${evse_id}
    Should Be Successful    ${r}
    Should Not Be Equal    ${r}[json][currentState]    ${unexpected}
    ...    msg=CS unexpectedly reached ${unexpected}

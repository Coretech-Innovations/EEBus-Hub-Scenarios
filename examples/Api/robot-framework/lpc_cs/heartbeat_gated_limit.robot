*** Settings ***
Documentation     Heartbeat-gated load-control acceptance for a controllable system (CS).
...
...    Both actors are simulated in one hub: the hub's built-in HEMS and an EVSE acting
...    as the controllable system (CS).
...
...    Intent: a controllable system must only accept a load-control write once the
...    controlling HEMS has established communication through its heartbeat.
...
...    Scenario:
...      Setup   : create and pair the HEMS and the CS, with no heartbeat sent yet.
...      Test 1  : send a deactivation limit (APCL active=false) -> REJECTED
...                ("No heartbeat received"); the CS is not controlled.
...      Test 2  : run the heartbeat handshake, then send the deactivation limit
...                (APCL active=false, value=0) -> ACCEPTED; the CS transitions to
...                Unlimited/Controlled.
...
...    Prerequisites: a hub running with the Remote server (default) and a license with
...    robotFrameworkSupported: true. See ../README.md.

Resource          cs_setup.resource
Suite Setup       Set Up And Connect The CS To The HEMS
Suite Teardown    Reset Simulation To A Clean State


*** Variables ***
${EVSE_ID}          ${None}
${DEVICE_ADDRESS}   ${EMPTY}


*** Test Cases ***
Deactivation Limit Is Rejected Before Any Heartbeat
    [Documentation]    The connection is already established (the suite setup waited
    ...    for it), so this is a single, deterministic write: with no heartbeat ever
    ...    received, the CS must reject it. A failure here means the CS wrongly accepted the
    ...    limit — not that the connection was not ready.
    ${r}=    Send Active Power Consumption Limit To Device By HEMS    ${DEVICE_ADDRESS}    active=${False}
    Should Be Equal    ${r}[json][success]    ${False}    The CS accepted a limit without a heartbeat
    Should Contain    ${r}[json][reason]    No heartbeat received    ${r}[json][reason]
    # The rejected write must not have moved the CS into a controlled state.
    CS LPC State Should Not Be    ${EVSE_ID}    Unlimited/Controlled

Deactivation Limit Is Accepted After A Heartbeat Was Received
    [Documentation]    Once the CS has received a heartbeat (even though it is
    ...    then stopped again), the same deactivation write is accepted and the CS moves
    ...    to Unlimited/Controlled.
    Send HEMS Heartbeat    ${EVSE_ID}
    ${r}=    Send Active Power Consumption Limit To Device By HEMS    ${DEVICE_ADDRESS}    active=${False}    value=${0}
    Should Be Equal    ${r}[json][success]    ${True}    The CS rejected the limit: ${r}[json][reason]
    CS LPC State Should Be    ${EVSE_ID}    Unlimited/Controlled


*** Keywords ***
Set Up And Connect The CS To The HEMS
    [Documentation]    Setup: create the two devices (heartbeat stopped) and pair them.
    ${hems_ski}    ${evse_id}=    Setup Env Init
    Set Suite Variable    ${EVSE_ID}    ${evse_id}
    ${address}=    Connect CS To HEMS    ${evse_id}    ${hems_ski}
    Set Suite Variable    ${DEVICE_ADDRESS}    ${address}

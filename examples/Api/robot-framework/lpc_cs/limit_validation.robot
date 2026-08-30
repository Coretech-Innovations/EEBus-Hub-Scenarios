*** Settings ***
Documentation     Active-limit validation for a controllable system (CS).
...
...    Both actors are simulated in one hub: the hub's built-in HEMS and an EVSE acting
...    as the controllable system (CS).
...
...    Intent: a controllable system accepts a valid active consumption limit (moving to
...    the Limited state) but must reject an invalid (negative) limit value, and a
...    rejected write must not change the already-applied state.
...
...    Scenario:
...      Setup   : create and pair the HEMS and the CS, then run the heartbeat handshake.
...      Test 1  : send an active limit with a valid value -> ACCEPTED; CS = Limited.
...      Test 2  : send an active limit with a NEGATIVE value -> REJECTED
...                ("invalid data received"); CS stays Limited.
...
...    Prerequisites: a hub running with the Remote server (default) and a license with
...    robotFrameworkSupported: true. See ../README.md.

Resource          cs_setup.resource
Suite Setup       Set Up Connect And Establish Heartbeat
Suite Teardown    Reset Simulation To A Clean State


*** Variables ***
${EVSE_ID}          ${None}
${DEVICE_ADDRESS}   ${EMPTY}
${LIMIT_WATTS}      ${5000}


*** Test Cases ***
Valid Active Limit Is Accepted And CS Becomes Limited
    [Documentation]    The connection is established and a heartbeat has been
    ...    received (both done in setup), so this is a single, deterministic write: a
    ...    well-formed active limit is accepted and drives the CS into the Limited state.
    ${r}=    Send Active Power Consumption Limit To Device By HEMS    ${DEVICE_ADDRESS}
    ...    active=${True}    value=${LIMIT_WATTS}    durationIndefinite=${True}
    Should Be Equal    ${r}[json][success]    ${True}    The CS rejected a valid limit: ${r}[json][reason]
    CS LPC State Should Be    ${EVSE_ID}    Limited

Negative Limit Value Is Rejected And State Is Unchanged
    [Documentation]    A negative limit value is invalid; the CS must reject it
    ...    and remain in the Limited state it already reached.
    ${r}=    Send Active Power Consumption Limit To Device By HEMS    ${DEVICE_ADDRESS}
    ...    active=${True}    value=${-1000}    durationIndefinite=${True}
    Should Be Equal    ${r}[json][success]    ${False}    The CS accepted a negative limit value
    Should Contain    ${r}[json][reason]    invalid data received
    CS LPC State Should Be    ${EVSE_ID}    Limited


*** Keywords ***
Set Up Connect And Establish Heartbeat
    [Documentation]    Setup: create the two devices, pair them, and run the heartbeat
    ...    handshake so the CS has established control before the limit writes.
    ${hems_ski}    ${evse_id}=    Setup Env Init
    Set Suite Variable    ${EVSE_ID}    ${evse_id}
    ${address}=    Connect CS To HEMS    ${evse_id}    ${hems_ski}
    Set Suite Variable    ${DEVICE_ADDRESS}    ${address}
    Send HEMS Heartbeat    ${evse_id}

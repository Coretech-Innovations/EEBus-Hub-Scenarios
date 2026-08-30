*** Settings ***
Documentation     Example 2 - EV and EVSE charging session.
...
...    A multi-entity scenario closer to a real EEBUS setup: create an EVSE and an
...    EV, connect the EV to the EVSE, and inspect the resulting state. It shows
...    patterns the first example does not:
...
...    - creating two different entity types and correlating their ids,
...    - a keyword with TWO path parameters (``Connect EV To EVSE`` takes evId + evseId),
...    - waiting on an asynchronous action with ``Wait Until Keyword Succeeds``,
...    - updating an entity with ``Modify EV`` and reading a sub-resource back,
...    - handling an expected failure (a bad id comes back as ``ok=False``, not a raise).
...
...    Prerequisites: a Hub running with ``--robotremote`` and a license that has
...    ``robotFrameworkSupported: true``. See README.md in this folder.

Resource          common.resource
Suite Setup       Set Up Fresh EV And EVSE
Suite Teardown    Reset Simulation To A Clean State


*** Variables ***
${EVSE_ID}    ${None}
${EV_ID}      ${None}


*** Test Cases ***
Connect The EV To The EVSE
    [Documentation]    ``Connect EV To EVSE`` takes two path parameters, passed
    ...    positionally in the order evId, evseId. Pairing between the two
    ...    simulated devices is asynchronous, so poll until it settles.
    ${r}=    Connect EV To EVSE    ${EV_ID}    ${EVSE_ID}
    Should Be Successful    ${r}
    Wait Until Keyword Succeeds    15x    2s    EV Reports A Connected EVSE

Update The EV And Read It Back
    [Documentation]    ``Modify EV`` accepts a partial body; only the named fields
    ...    are sent. Then confirm the change is reflected when reading the entity.
    ${updated}=    Modify EV    ${EV_ID}    charged=${42}    chargingEnable=${False}
    Should Be Successful    ${updated}

    ${view}=    Get EV    ${EV_ID}
    Should Be Successful    ${view}

Asking For A Nonexistent EV Fails Gracefully
    [Documentation]    A 4xx does not raise on its own; it returns with ``ok=False``.
    ...    That is how you assert on expected error paths.
    ${r}=    Get EV    999999
    Should Be Equal As Integers    ${r}[status_code]    404
    Should Not Be True    ${r}[ok]


*** Keywords ***
Set Up Fresh EV And EVSE
    [Documentation]    Reset the simulation, then create one EVSE and one EV and
    ...    remember their ids for the whole suite.
    Reset Simulation To A Clean State

    ${evse}=    Add EVSE    evseccEnabled=${True}
    Should Be Successful    ${evse}
    Set Suite Variable    ${EVSE_ID}    ${evse}[json][id]

    ${ev}=    Add EV
    Should Be Successful    ${ev}
    Set Suite Variable    ${EV_ID}    ${ev}[json][id]

EV Reports A Connected EVSE
    [Documentation]    Poll helper: succeeds once the EV listing shows this EV paired
    ...    to our EVSE. ``List EVs`` returns an array; each entry carries
    ...    ``connectedToEVSE`` and (once paired) an ``evseId``.
    ${listing}=    List EVs
    Should Be Successful    ${listing}
    ${matches}=    Evaluate
    ...    [e for e in $listing['json'] if e['evId'] == ${EV_ID}]
    Should Not Be Empty    ${matches}    EV ${EV_ID} not found in the listing
    ${ev}=    Set Variable    ${matches}[0]
    Should Be True    ${ev}[connectedToEVSE]    EV ${EV_ID} is not connected yet
    Should Be Equal As Integers    ${ev}[evseId]    ${EVSE_ID}

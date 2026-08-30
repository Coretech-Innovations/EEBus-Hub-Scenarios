*** Settings ***
Documentation     Example 1 - EVSE lifecycle (create / list / modify / delete).
...
...    A self-contained, deterministic tour of the CRUD pattern that applies to
...    every simulated entity. It demonstrates the four things a test author needs:
...
...    - creating an entity with NO arguments (the Hub fills a valid default body),
...    - overriding only the fields you care about (friendly aliases, e.g.
...      ``nominalPowerMax`` maps to the nested ``nominalPower.max`` JSON field),
...    - path parameters (the numeric entity id) passed positionally,
...    - reading the result dictionary (``status_code`` / ``json`` / ``text`` / ``ok``).
...
...    Prerequisites: a Hub running with ``--robotremote`` and a license that has
...    ``robotFrameworkSupported: true``. See README.md in this folder.

Resource          common.resource
Suite Setup       Reset Simulation To A Clean State


*** Test Cases ***
Add EVSE With Defaults Then Find It In The List
    [Documentation]    A no-argument create uses the Hub's built-in default body,
    ...    so a valid EVSE is created without spelling out ~20 fields.
    [Tags]    smoke
    ${created}=    Add EVSE    vendorCode=TestVendor
    Should Be Successful    ${created}
    ${evse_id}=    Set Variable    ${created}[json][id]

    ${listing}=    List EVSEs
    Should Be Successful    ${listing}
    ${ids}=    Evaluate    [item['evseId'] for item in $listing.get('json') or []]
    Should Contain    ${ids}    ${evse_id}    The created EVSE should appear in the list
    [Teardown]    Delete EVSE    ${evse_id}

Create EVSE Overriding Only The Fields I Care About
    [Documentation]    Pass just the fields you want to change; everything else keeps
    ...    its default. ``nominalPowerMax`` is a friendly alias for the nested
    ...    ``nominalPower.max`` field. A value given as text (``"11000"``) is coerced
    ...    to the schema's numeric type, so both ``11000`` and ``${11000}`` work.
    ${created}=    Add EVSE    deviceName=Garage Wallbox    nominalPowerMax=${11000}
    Should Be Successful    ${created}
    ${evse_id}=    Set Variable    ${created}[json][id]
    [Teardown]    Delete EVSE    ${evse_id}

Modify A Single Field On An Existing EVSE
    [Documentation]    PATCH tolerates a partial body: the numeric id is a path
    ...    parameter (positional), and only the named field is sent.
    ${created}=    Add EVSE
    ${evse_id}=    Set Variable    ${created}[json][id]

    ${updated}=    Modify EVSE    ${evse_id}    nominalPowerMax=${16000}
    Should Be Successful    ${updated}
    [Teardown]    Delete EVSE    ${evse_id}

Delete EVSE Removes It From The List
    [Documentation]    After deletion the entity id must no longer be listed.
    ${created}=    Add EVSE
    ${evse_id}=    Set Variable    ${created}[json][id]

    ${deleted}=    Delete EVSE    ${evse_id}
    Should Be Successful    ${deleted}

    ${listing}=    List EVSEs
    Should Be Successful    ${listing}
    ${ids}=    Evaluate    [item['evseId'] for item in $listing.get('json') or []]
    Should Not Contain    ${ids}    ${evse_id}    Deleted EVSE should be gone from the list

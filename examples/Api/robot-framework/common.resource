*** Settings ***
Documentation    Shared setup for the EEBUS Hub Robot Framework examples.
...
...    Imports the Hub's Remote Library (served live from the running binary) and
...    provides a couple of thin helpers the example suites reuse. This file contains
...    NO keyword implementations of its own for the Hub API — every Hub keyword
...    (``Add EVSE``, ``Post Simulation``, ...) is served remotely by the binary.

Library          Remote    http://${HUB_HOST}:${HUB_RF_PORT}/


*** Variables ***
${HUB_HOST}       localhost
${HUB_RF_PORT}    8270


*** Keywords ***
Should Be Successful
    [Documentation]    Fail unless the Hub returned a 2xx status.
    ...    Every Hub keyword returns a dictionary: ``status_code``, ``json``,
    ...    ``text`` and ``ok`` (true for HTTP 2xx). A 4xx/5xx does NOT raise on its
    ...    own — it comes back with ``ok=False`` — so assert on it explicitly.
    [Arguments]    ${response}
    Should Be True    ${response}[ok]
    ...    Hub call failed: HTTP ${response}[status_code] -> ${response}[text]

Reset Simulation To A Clean State
    [Documentation]    Remove all simulated entities so each suite starts fresh.
    ...    ``Reset Simulation`` is a composite keyword served by the Hub — it wraps
    ...    ``Post Simulation    action=reset``. See robot/composite-keywords.yaml.
    ${r}=    Reset Simulation
    Should Be Successful    ${r}
    RETURN    ${r}

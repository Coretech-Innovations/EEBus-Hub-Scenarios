# OPEV — Overload Protection by EV Current Curtailment

Runnable examples of the **Overload Protection by EV Current Curtailment (OPEV)** use
case, in which a HEMS keeps the total current drawn by one or more EVs within the
building's available capacity by curtailing charging current.

## Actors

- **HEMS** — enforces the site current limit and distributes the available current.
- **EV / EVSE** — the charging station and vehicle whose current is curtailed.

## Examples

| Example | What it demonstrates |
|---------|----------------------|
| `single-ev-full-current` | A single EV charges at its maximum current when the site has ample capacity. |
| `household-load-stops-ev` | An uncontrollable household load consumes the whole current budget, so the EV is curtailed to zero. |
| `fair-share-two-evs` | Two EVs share the available current when both fit under the site limit. |
| `fair-share-with-household-load` | Two EVs share a reduced current budget after a household load takes part of the capacity. |
| `external-ev` | A HEMS curtails the current of an external EV identified by its SKI. |
| `external-hems` | A simulated EV receives its current limits from an external HEMS identified by its SKI. |
| `two-evs-external-hems` | Two EV/EVSE pairs receive their limits from an external HEMS. |

Run any example from the repository root, for example:

```bash
go run ./examples/Api/OPEV/single-ev-full-current
```

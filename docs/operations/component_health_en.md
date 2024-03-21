# Component health

The component operator tracks the current health status of components installed in a multinode EcoSystem
in the status of the 'Component' resource.
The following health states are currently possible:
- Pending (`""`) - the health state has not yet been written.
- Available (`"available"`)
- Not available (`"unavailable"`)
- Unknown (`"unknown"`) - the health status of the component is unknown or cannot be guaranteed,
  [see below](#special-case-switching-off-the-component-operator).

## Special case: Switching off the component operator

If the component operator is terminated regularly (via `sigint` or `sigterm`),
it sets its own health to `unavailable` and that of all other components to `unknown`
in order to avoid misleading states.
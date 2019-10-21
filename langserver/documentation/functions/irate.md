## `irate()`

`irate(v range-vector)` calculates the per-second instant rate of increase of
the time series in the range vector. This is based on the last two data points.
Breaks in monotonicity (such as counter resets due to target restarts) are
automatically adjusted for.

The following example expression returns the per-second rate of HTTP requests
looking up to 5 minutes back for the two most recent data points, per time
series in the range vector:

```
irate(http_requests_total{job="api-server"}[5m])
```

`irate` should only be used when graphing volatile, fast-moving counters.
Use `rate` for alerts and slow-moving counters, as brief changes
in the rate can reset the `FOR` clause and graphs consisting entirely of rare
spikes are hard to read.

Note that when combining `irate()` with an
[aggregation operator](operators.md#aggregation-operators) (e.g. `sum()`)
or a function aggregating over time (any function ending in `_over_time`),
always take a `irate()` first, then aggregate. Otherwise `irate()` cannot detect
counter resets when your target restarts.

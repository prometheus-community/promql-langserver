## `increase()`

`increase(v range-vector)` calculates the increase in the
time series in the range vector. Breaks in monotonicity (such as counter
resets due to target restarts) are automatically adjusted for. The
increase is extrapolated to cover the full time range as specified
in the range vector selector, so that it is possible to get a
non-integer result even if a counter increases only by integer
increments.

The following example expression returns the number of HTTP requests as measured
over the last 5 minutes, per time series in the range vector:

```
increase(http_requests_total{job="api-server"}[5m])
```

`increase` should only be used with counters. It is syntactic sugar
for `rate(v)` multiplied by the number of seconds under the specified
time range window, and should be used primarily for human readability.
Use `rate` in recording rules so that increases are tracked consistently
on a per-second basis.

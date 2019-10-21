## `holt_winters()`

`holt_winters(v range-vector, sf scalar, tf scalar)` produces a smoothed value
for time series based on the range in `v`. The lower the smoothing factor `sf`,
the more importance is given to old data. The higher the trend factor `tf`, the
more trends in the data is considered. Both `sf` and `tf` must be between 0 and
1.

`holt_winters` should only be used with gauges.

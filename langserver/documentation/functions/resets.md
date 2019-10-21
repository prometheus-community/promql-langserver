## `resets()`

For each input time series, `resets(v range-vector)` returns the number of
counter resets within the provided time range as an instant vector. Any
decrease in the value between two consecutive samples is interpreted as a
counter reset.

`resets` should only be used with counters.

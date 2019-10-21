## `absent()`

`absent(v instant-vector)` returns an empty vector if the vector passed to it
has any elements and a 1-element vector with the value 1 if the vector passed to
it has no elements.

This is useful for alerting on when no time series exist for a given metric name
and label combination.

```
absent(nonexistent{job="myjob"})
# => {job="myjob"}

absent(nonexistent{job="myjob",instance=~".*"})
# => {job="myjob"}

absent(sum(nonexistent{job="myjob"}))
# => {}
```

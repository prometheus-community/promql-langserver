## `<aggregation>_over_time()`

The following functions allow aggregating each series of a given range vector
over time and return an instant vector with per-series aggregation results:

* `avg_over_time(range-vector)`: the average value of all points in the specified interval.
* `min_over_time(range-vector)`: the minimum value of all points in the specified interval.
* `max_over_time(range-vector)`: the maximum value of all points in the specified interval.
* `sum_over_time(range-vector)`: the sum of all values in the specified interval.
* `count_over_time(range-vector)`: the count of all values in the specified interval.
* `quantile_over_time(scalar, range-vector)`: the φ-quantile (0 ≤ φ ≤ 1) of the values in the specified interval.
* `stddev_over_time(range-vector)`: the population standard deviation of the values in the specified interval.
* `stdvar_over_time(range-vector)`: the population standard variance of the values in the specified interval.

Note that all values in the specified interval have the same weight in the
aggregation even if the values are not equally spaced throughout the interval.

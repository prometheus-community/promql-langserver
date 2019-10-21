## `histogram_quantile()`

`histogram_quantile(φ float, b instant-vector)` calculates the φ-quantile (0 ≤ φ
≤ 1) from the buckets `b` of a
[histogram](https://prometheus.io/docs/concepts/metric_types/#histogram). (See
[histograms and summaries](https://prometheus.io/docs/practices/histograms) for
a detailed explanation of φ-quantiles and the usage of the histogram metric type
in general.) The samples in `b` are the counts of observations in each bucket.
Each sample must have a label `le` where the label value denotes the inclusive
upper bound of the bucket. (Samples without such a label are silently ignored.)
The [histogram metric type](https://prometheus.io/docs/concepts/metric_types/#histogram)
automatically provides time series with the `_bucket` suffix and the appropriate
labels.

Use the `rate()` function to specify the time window for the quantile
calculation.

Example: A histogram metric is called `http_request_duration_seconds`. To
calculate the 90th percentile of request durations over the last 10m, use the
following expression:

    histogram_quantile(0.9, rate(http_request_duration_seconds_bucket[10m]))

The quantile is calculated for each label combination in
`http_request_duration_seconds`. To aggregate, use the `sum()` aggregator
around the `rate()` function. Since the `le` label is required by
`histogram_quantile()`, it has to be included in the `by` clause. The following
expression aggregates the 90th percentile by `job`:

    histogram_quantile(0.9, sum(rate(http_request_duration_seconds_bucket[10m])) by (job, le))

To aggregate everything, specify only the `le` label:

    histogram_quantile(0.9, sum(rate(http_request_duration_seconds_bucket[10m])) by (le))

The `histogram_quantile()` function interpolates quantile values by
assuming a linear distribution within a bucket. The highest bucket
must have an upper bound of `+Inf`. (Otherwise, `NaN` is returned.) If
a quantile is located in the highest bucket, the upper bound of the
second highest bucket is returned. A lower limit of the lowest bucket
is assumed to be 0 if the upper bound of that bucket is greater than
0. In that case, the usual linear interpolation is applied within that
bucket. Otherwise, the upper bound of the lowest bucket is returned
for quantiles located in the lowest bucket.

If `b` contains fewer than two buckets, `NaN` is returned. For φ < 0, `-Inf` is
returned. For φ > 1, `+Inf` is returned.

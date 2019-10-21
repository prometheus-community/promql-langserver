## `round()`

`round(v instant-vector, to_nearest=1 scalar)` rounds the sample values of all
elements in `v` to the nearest integer. Ties are resolved by rounding up. The
optional `to_nearest` argument allows specifying the nearest multiple to which
the sample values should be rounded. This multiple may also be a fraction.

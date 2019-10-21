# Progress

## 11/10/2019

* implemented basic text synchronization
* PR for position data in promql is ready: https://github.com/prometheus/prometheus/pull/6061

### TODO

* return hover information; expected monday evening
* consider adding child to implementation in languiage server
* documeting funcs: use a map from func name to struct
* for e2e tests: invest 1-2hrs researching gopls and check back

## 14/10/2019

* passing messages to language client but VSCode is ignoring them
* found ~5-line Vim script to connect language client to language server
* will test if Vim language client works

## 21/10/2019

* identified issue preventing requests from being sent successfully: request context was cancelled prematurely
* can now successfully send error messages to client and display them in VSCode
* will use go-bindata or statik to compile documentation strings into the binary
* have working example of `label_replace` hover text
* VSCode's monaco is >2mb minified, may need to find alternative language client implementation for Prometheus frontend

### TODO

* finish function documentation hover in the next few days
* implement lite parser ~1 week
* minimum e2e tests afterwards

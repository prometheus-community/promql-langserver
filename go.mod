module github.com/slrtbtfs/promql-lsp

go 1.13

require (
	github.com/alecthomas/units v0.0.0-20190924025748-f65c72e2690d // indirect
	github.com/json-iterator/go v1.1.8 // indirect
	github.com/prometheus/client_golang v1.2.1
	github.com/prometheus/procfs v0.0.6 // indirect
	github.com/prometheus/prometheus v0.0.0-20180315085919-58e2a31db8de
	github.com/prometheus/tsdb v0.10.0 // indirect
	github.com/rakyll/statik v0.1.7-0.20191104211043-6b2f3ee522b6
	github.com/slrtbtfs/go-tools-vendored v0.0.0-20191025135527-118fa9afc3f0
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/sys v0.0.0-20191110163157-d32e6e3b99c4 // indirect
	golang.org/x/tools v0.0.0-20190918214516-5a1a30219888
	golang.org/x/xerrors v0.0.0-20191011141410-1b5146add898
	gopkg.in/yaml.v3 v3.0.0-20191107175235-0b070bb63a18
)

replace github.com/prometheus/prometheus => github.com/slrtbtfs/prometheus v1.8.2-0.20191111173200-f07eddb4fa6d

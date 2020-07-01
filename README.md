[![CircleCI](https://circleci.com/gh/prometheus-community/promql-langserver.svg?style=svg)](https://circleci.com/gh/prometheus-community/promql-langserver)
[![Go Report Card](https://goreportcard.com/badge/github.com/prometheus-community/promql-langserver)](https://goreportcard.com/report/github.com/slrtbtfs/promql-lsp)
[![GoDoc](https://godoc.org/github.com/prometheus-community/promql-langserver?status.png)](https://pkg.go.dev/github.com/prometheus-community/promql-langserver)
![golangci-lint](https://github.com/prometheus-community/promql-langserver/workflows/golangci-lint/badge.svg)

# promql-lsp

Implementation of the Language Server Protocol for PromQL.

## Roadmap

- [x] Connect to IDEs and text editors over
  - [x] Stdio
  - [ ] Websocket
  - [ ] HTTP
- [x] Sync document content with the client
- [x] Support plain PromQL queries
- [x] Support queries inside yaml files (e.g. alertmanager configuration)
- [x] Connect to a prometheus instance to get label and metric data
- [x] Show error messages for incorrect queries in the client
- [ ] Show documentation on hover
  - [x] Type information
  - [x] Function documentation
  - [x] Aggregator documentation
  - [ ] Keyword documentation
  - [x] Metric and label help strings from a prometheus instance
- [ ] Completion
  - [x] Functions
  - [x] Metrics
  - [x] Recording Rules
  - [x] Aggregators
  - [x] Labels
  - [x] Label Values
  - [ ] Context sensitive, i.e respecting function argument types
- [x] Signature information for functions (while typing)
- [ ] (Linting)
- [ ] (Formatting)

## Some Screenshots

### Completion in VS Code

![Completion in VS Code](https://github.com/prometheus-community/promql-langserver/raw/master/screenshots/vscode_completion.png)

### Metric metadata from a Prometheus Server

![Metric Metadata in VS Code](https://github.com/prometheus-community/promql-langserver/raw/master/screenshots/vscode_hover2.png)

### Viewing documentation right from your editor

![Docs in VS Code](https://github.com/prometheus-community/promql-langserver/raw/master/screenshots/vscode_yaml.png)

### Vim and other editors are supported, too

![Vim](https://github.com/prometheus-community/promql-langserver/raw/master/screenshots/vim.png)

## Using the Language Server

A Language Server on its own is not very useful. You need some Language Client to use it with.

The following Language Clients have been tested with this language server. More will be added in the future. 

Feel free to reach out if you want to use it with another Editor/Tool.

Reading this [documentation](./doc/developing_editor.md) can help you in your work.

### VS Code

There exists a VS Code extension based on this language server: <https://github.com/slrtbtfs/vscode-prometheus>

It is used as the main test platform for this language server.

Since it isn't published on the Extensions Marketplace yet, you'll have to follow the somewhat more complicated installation steps described in the README there.

### (Neo)Vim 

With Vim, currently only PromQL queries inside YAML files work without additional support. Generally the experience with Vim is more buggy than with VS Code.

#### Setup

1. Install the [YouCompleteMe](https://github.com/ycm-core/YouCompleteMe) Plugin.
2. Put the configuration following configuration file for the language server in `.vim/promql-lsp.yaml`.

        # Change this adress to the address of the prometheus server you want to use for metadata
        prometheus_url: http://localhost:9090
        rpc_trace: text

3. Add the following to your `.vimrc`

        let g:ycm_language_server = [
          \   { 'name': 'promql',
          \     'filetypes': [ 'yaml' ],
          \     'cmdline': [ 'promql-langserver', '--config-file', expand('~/.vim/promql-lsp.yaml')]
          \   },
          \ ]

#### Debugging

The Vim command `:YcmDebugInfo` gives status information and points to logfiles.

### Sublime Text 3

1. Install package `LSP`, `LSP-promql` via `Package Control`.
2. Follow the [installation instruction](https://github.com/nevill/lsp-promql#installation).

## Contributing

Refer to [CONTRIBUTING.md](./CONTRIBUTING.md)

## License

Apache License 2.0, see [LICENSE](./LICENSE).

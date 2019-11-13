[![Go Report Card](https://goreportcard.com/badge/github.com/slrtbtfs/promql-lsp)](https://goreportcard.com/report/github.com/slrtbtfs/promql-lsp)
[![Build Status](https://cloud.drone.io/api/badges/slrtbtfs/promql-lsp/status.svg)](https://cloud.drone.io/slrtbtfs/promql-lsp)

**Warning:** This software is not in a very useful state yet.

# promql-lsp

Implementation of the Language Server Protocol for PromQL.

## Features

Most of the planned features are not implemented yet.

- [x] Connect to IDEs and text editors over stdio
- [ ] Connect to remote clients over websocket or http
- [x] Sync document content with the client
- [x] Support plain PromQL queries
- [x] Support queries inside yaml files
- [x] Connect to a prometheus instance to get label and metric data
- [x] Show error messages for incorrect queries in the client
- [ ] Show documentation on hover
  - [x] Type information
  - [x] Function documentation
  - [ ] Operator documentation
  - [x] Metric and label help strings from a prometheus instance
  - [ ] Even when the Query has syntax errors
- [ ] Completion
  - [ ] Functions
  - [x] Metrics
  - [ ] Labels
  - [ ] Context sensitive, i.e respecting function argument types
  - [ ] Even when the Query has syntax errors
- [ ] Signature information for functions (while typing)
- [ ] (Linting)
- [ ] (Formatting)

## Using the Language Server

A Language Server on its own is not very useful. You need some Language Client to use it with.

The following Language Clients have been tested with this language server. More will be added in the future. 

Feel free to reach out if you want to use it with another Editor/Tool.

### VS Code

There exists a VS Code extension based on this language server: <https://github.com/slrtbtfs/vscode-prometheus>

It is used as the main test platform for this langauge server.

Since it isn't published a on the Extensions Marketplace yet, you'll have to follow the somewhat more complicated installation steps described in the README there.

### (Neo)Vim 

With Vim, currently only PromQL queries inside YAML files work without additional support. Generally the experience with Vim is more buggy than with VS Code.

#### Setup

1. Install the [YouCompleteMe](https://github.com/ycm-core/YouCompleteMe) Plugin.
2. Put the configuration following configuration file for the langauge server in `.vim/promql-lsp.yaml`.

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
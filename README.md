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
  - [ ] Even when the has syntax errors
- [ ] Completion
  - [ ] Functions
  - [x] Metrics
  - [ ] Labels
  - [ ] Context sensitive, i.e respecting function argument types
  - [ ] Even when the has syntax errors
- [ ] Signature information for functions (while typing)
- [ ] (Linting)
- [ ] (Formatting)


---
title: PromQL Language Server
type: Proposal
menu: proposals
status: WIP
owner: slrtbtfs
---

## Summary

This document describes the motivation and design of the PromQL language server, as well as how it can be integrated into existing development environments.

## Motivation

Modern IDEs and editors (tools) provide a lot of language specific features like syntax highlighting, autocompletion or hover information. Traditionally for each combination of tool and language an own integration had to be written which leads a lot of redundant work, inconsistency and especially disadvantages tools without large communities.

To improve that situation Microsoft introduced the [Language Server Protocol (LSP)](https://microsoft.github.io/language-server-protocol/) which has since been adopted by a large set of [languages](https://microsoft.github.io/language-server-protocol/implementors/servers/) and [tools](https://microsoft.github.io/language-server-protocol/implementors/tools/) and is [endorsed by Red Hat](https://developers.redhat.com/blog/2016/06/27/a-common-interface-for-building-developer-tools/).

LSP is a standardized Protocol for communication between an tool (client) and a Language Server that implements all language specific logic. As a result implementing a server for a language makes it fairly easy to integrate it with most tool and vice versa.

PromQL does not have such a language Server yet. Some tools, like the Prometheus expression browser do have limited support for autocompletion tough.

## Technical summary of PromQL language server

To improve the general user experience of Prometheus a PromQL language server is proposed.

The PromQL language Server can either be included into IDE/editor plugins or run on a server and communicate with tools like the Prometheus Web UI over a network. PromQL Queries are supported both standalone and as part of a .yaml configuration file.

To provide autocompletion for labels it can optionally connect to an Prometheus Server and use it's label data.

For testing purposes some IDE plugins will be developed. These will include a [TextMate Grammar](https://macromates.com/manual/en/language_grammars) to enable syntax highlighting which is not supported by the LSP itself.

## Architecture

### Implemented Server Capabilities

The [LSP Specification](https://microsoft.github.io/language-server-protocol/specification) describes a rich set of server capabilities, not all of whom are useful for PromQL. It is possible for a server to only implement a subset and advertise them at initialization.

For capabilities marked with _maybe_ it is not yet clear, wether a use case exist. An exclamation mark denotes that these capabilities are prioritized.

The following capabilities will be implemented by the PromQL Language Server:

#### General Capabilities: _all!_

Necessary to establish and end communication with a client.

#### Window Capabilities: _all_

Enable the server to send notifications and log messages to the client.

#### Telemetry Capabilities: _maybe_

Enable the server to send telemetry events to the client. Might be useful once there are published IDE integrations.

#### Client Capabilities: _all!_

Enable a client to advertise it's capabilities to the server.

#### Workspace Capabilities: _maybe some_

For PromQL the concept of workspaces is not relevant. Only implemented if required by another capability.

#### Text Synchronization Capabilities: _all!_

Notify the server about File Changes. Mandatory to be able to inspect the content of unsaved files.

#### Diagnostics Capabilities: _all!_

Send Errors and Warnings to the Client. These are used to show syntax errors and linting. The linter warns about common errors such as `rate(sum(...))` or `http_requests_total{status="^5..$"}`.

#### Language Capabilities: _some!_

The core part of the language server. Some of these, e.g. Go to (definition|typeDefinition|declaration|implementation), renaming and folding are not useful for PromQL itself. Implemented are:

##### completion!, completion_resolve!

Give completions for functions, operators and, if a Prometheus server is attached, labels.

##### hover!

Show documentation for functions and operators.

##### signatureHelp

Show the type of expressions (`string`, `scalar`, `Instant vector`, `Range vector` and functions combining these).

##### codeAction: _maybe_

Enable QuickFixes. In case this is implemented, also some of the Workspace Capabilities would be required.

##### formatting, rangeFormatting, onTypeFormatting: _maybe_

There isn't that much formatting changes that could be done, other than ensuring there is a sane amount of white space.

### JSON-RPC

The LSP Protocol is based on the [JSON-RPC 2.0 Protocol](https://www.jsonrpc.org/specification). The transport layer is not specified by the protocol. For IDEs it's usually stdin/stdout, over a Network it could be HTTP. Thus the JSON-RPC implementation used by PromQL should abstract over the transport layer.

### Parsing

See proposal 2019_promql-parser-improvements.

## Open Questions

* How do we make this server compatible with the [YAML Language Server](https://github.com/redhat-developer/yaml-language-server)?
* Which JSON-RPC implementation to use (probably [the one from gopls](https://github.com/golang/tools/tree/master/internal/jsonrpc2))?
* Can we reuse [the Server Implementation of gopls](https://github.com/golang/tools/tree/master/internal/lsp/protocol)? Note that there is a lot of code generated by a typescript file and some black magic.

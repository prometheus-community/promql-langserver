# About [LSP](https://microsoft.github.io/language-server-protocol/)

From their homepage:

> The Language Server Protocol (LSP) defines the protocol used between an editor or IDE and a language server that provides language features like auto complete, go to definition, find all references etc.

* Standardized protocol supported by a lot of IDEs and languages
* Originally created for VS Code.
* Creating an LSP server for PromQL would add these features to most development environments.

## What can LSP do?
Most language insight features are not really useful for a query language.

## Useful for PromQL
* Diagnostics (e.g. linting and syntax errors)
* Auto completion
* Hover information

## Maybe useful for PromQL
* Find references (e.g. of a specific label)
* Document Formatting Requests (e.g. beautifying)

# About [LSP](https://microsoft.github.io/language-server-protocol/)

From their homepage:

> The Language Server Protocol (LSP) defines the protocol used between an editor or IDE and a language server that provides language features like auto complete, go to definition, find all references etc.

* Standardized protocol supported by a lot of IDEs and languages
* Originally created for VS Code.
* Red Hat announced adoption in 2016
* Creating an LSP server for PromQL would add these features to most development environments.

# What can LSP do?
Many language insight features are not really useful for a query language.

## Not useful for PromQL
* Go to (definition|typeDefinition|declaration|implementation)
* renaming
* folding
* Workspaces

## Useful for PromQL
* Diagnostics (e.g. linting and syntax errors)
* Auto completion
* Hover information

## Maybe useful for PromQL
* Find references (e.g. of a specific label)
* Document Formatting Requests (e.g. beautifying)

## Some concrete examples
* Complain about syntax errors before the request has been submitted
* Offer autocompletition for labels and values
* Complain if a request does not match anything
* Warn about request such as `http_requests_total{status="^4..$"}` or `rate(sum(...))`
* Show documentation of functions on hover

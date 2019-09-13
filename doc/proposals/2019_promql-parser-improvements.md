---
title: PromQL Language Server
type: Proposal
menu: proposals
status: WIP
owner: slrtbtfs
---

## PromQL Parser Improvements

### Motivation

Prometheus currently uses a hand written parser, which is generally in a good shape and reliably catches syntax errors. The error messages are mostly helpful. It hasn't seen major changes in the past few years.

The data stored in the generated abstract syntax tree isn't sufficient for the use case of a language server tough. The following issues need to be addressed.the

1. When an syntax error occurs, the parser aborts and a hard coded string is used to report an error to the user. For the language server errors should be represented by a more advanced data structure which can also store the position of the errors and possible Quick Fixes. The list of possible errors should be kept in one file and not be spread over several source files.
2. The Nodes in the abstract syntax tree are only required to store which substring they represent, not where in the code this substring is located.
3. For the purposes of autocompletion an AST should still be generated, when some closing parentheses are missing.
4. Some of the error messages should explain themselves better.

Since these changes would also benefit the upstream prometheus implementation and consistence between the language server and prometheus itself is desired, it is proposed that the PromQL language server does not implement it's own parser. Instead all necessary changes should go into the upstream parser.

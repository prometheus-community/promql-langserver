---
title: PromQL Parser Improvements
type: Proposal
menu: proposals
status: WIP
owner: slrtbtfs
---

## PromQL Parser Improvements

### Motivation

Prometheus currently uses a hand written parser, which is generally in a good shape and reliably catches syntax errors. The error messages are mostly helpful. It hasn't seen major changes in the past few years, so some bitrot might have set in.

For the planned PromQL language server (see Proposal 2019_promql-language-server) it is desirable to use the same parser as prometheus to ensure consistency and avoid code duplication.

The current Parser does is not sufficient for that use case though.

The proposed changes include adding the necessary features to the parser, improving error handling and cleaning up the existing code Base.

### Proposed Changes

*TODO*.

#### Change Title

##### Problem

...

##### Proposed solution

...

---

1. When an syntax error occurs, the parser aborts and a hard coded string is used to report an error to the user. For the language server errors should be represented by a more advanced data structure which can also store the position of the errors and possible Quick Fixes. The error messages are spread over several source files.
2. The Nodes in the abstract syntax tree are only required to store which substring they represent, not where in the code this substring is located.
3. For the purposes of autocompletion an AST should still be generated, when some closing parentheses are missing.
4. Some of the error messages should explain themselves better.

Since these changes would also benefit the upstream prometheus implementation and consistence between the language server and prometheus itself is desired, it is proposed that the PromQL language server does not implement it's own parser. Instead all necessary changes should go into the upstream parser.

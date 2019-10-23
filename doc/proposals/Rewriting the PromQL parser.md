# Rewriting the PromQL Parser

## Whats wrong with the existing parser? (rant)

A lot of things:

### The Abstract Syntax Tree is too flat

Look at query:

    metric{label=="value", foo~="ba."}[5m] offset 3h

One would expect a the structure of the respective syntax tree approximately look like this:

    |---- OffsetModifier: metric{label=="value", foo~="ba."}[5m] offset 3h
    . . . |---- MatrixSelector: metric{label=="value", foo~="ba."}[5m]
    . . . . . . |---- VectorSelector: metric{label=="value", foo~="ba."}
    . . . . . . . . . |---- LabelMatcher: label=="value"
    . . . . . . . . . . . . |---- Label: label
    . . . . . . . . . . . . |---- Value: "value"
    . . . . . . . . . |---- LabelMatcher: foo~="ba."
    . . . . . . . . . . . . |---- Label: foo 
    . . . . . . . . . . . . |---- Value: "ba."

This is how it actually looks like with the current parser:

    |---- MatrixSelector: metric{label=="value", foo~="ba."}[5m] offset 3h

This is ok if all we want to do is evaluating PromQL queries. 

For using the data from the parser for things like autocompletion that works anywhere in the code, generating helpful error messages, showing signature documentation and documentation when editing queries, etc., it isn't.

### The parser code is more complex than it would need to be

#### switch statements vs interfaces

When applying an operation to a lot of different data structures in go the idiomatic way is to use an interface:

    type interface foo {
        doThing()
    }

A much uglier alternative is:

    func doThing(obj foo) {
        switch obj.(type){
            case a:
                // doThing for type a
            case b:
                // doThing for type b
            default:
                panic("Forgot to implement something")
        }
    }

For some reason the current parser uses the second variant a lot. This also causes a lot of code duplication.


#### Inconsistent recursion

The basic structure of the parser is like this.

    func parseAorBorCorD() *node{
        // try to parse A if successful return it 
        ...
        // recognize B
        if ...{
            return parseB()
        }else if  ... {// recognize C
            return parseC()
        } else {
            return parseD()
        }
    }

    func parseB() *node{
        // parse B
        ...
        if ... {// recognize C
            C := parseC()
            // Apply some modifications to C that depend on B
            ...
            return C
        } 
    }

This makes it very hard to follow what is going on. To make things worse all these  `parse...()` functions don't follow a common convention what their arguments and return types are. Examples from the actual code include:

    func (p *parser) parseSeriesDesc() (m labels.Labels, vals []sequenceValue, err error)
    func (p *parser) unaryExpr() Expr
    func (p *parser) subqueryOrRangeSelector(expr Expr, checkRange bool) Expr
    func (p *parser) number(val string) float64
    func (p *parser) labels() []string

### The current parser does not return any results on incomplete expressions

In the language server I want to use the data from the parser to provide completion and function signatures while the user is typing. That's not possible with the current parser.

## Why not a generated parser?

Parser Generators are nice, because they drastically reduce the amount of work for building a parser. However, by design, they are only useful if for recognizing a specified formal grammar. 
If they should deal with errors, one has to specify a formal grammar that describes each error. The same applies for parsing incomplete expressions. 

There has been a decision to eventually replace the current PromQL parser by a generated one. The reason for this mainly was having a formal grammar that could be used by IDEs, language servers and syntax highlighting. This doesn't apply anymore.

    * For highlighting typically TextMate grammars are used. These are far less powerful as those for parser generators, so grammar reuse is not possible here. Also, such a TextMate Grammar for PromQL already exists.
    * For the language server I'd like to use the output from the parser, so I don't need a grammar here
    * IDE integrations are all gonna use the language server

Given that features that are required for the language server, like handling of incomplete Expressions and producing helpful error messages are really hard to do with generated grammars, I propose abandoning the aforementioned decision.

## How is the new parser supposed to look like?

## How do we get there?
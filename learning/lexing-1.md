In order to work easily with the source code, some transformations are needed before interpreting source code
the first transformation we'll do is lexing, lexing or tokenizer it transform the source code into "tokens", tokens are
a data structure which categorize  elements of the syntax, for example for the following source code

```
let x = 5 + 5;
```

it would produce the following result
[
LET,
IDENTIFIER("x"),
EQUAL_SIGN,
INTEGER_LITERAL(5),
PLUS_SIGN,
INTEGER_LITERAL(5),
SEMICOLON
]

this output would be the feeding to the parser later on, note that tokens has attached the values or information we care
about those tokens, identifier name "x" or value 5 for the INTEGER_LITERAL
function NewToken(type, text) {

  let this = {
    "type": type,
    "text": text,
  };

  this.is = fn(type) {
    return this.type == type
  }

  return this;
}

function NewLexer(input) {

  let pos = 0;
  let ch = input[0];
  let this = {};

  let bytetokens = {
    ";":  "SEMICOLON",
    ":    ": "COLON",
    "(":  "LPAREN",
    ")":  "RPAREN",
    "{":  "LBRACE",
    "}":  "RBRACE",
    "[":  "LBRACKET",
    "]":  "RBRACKET",
    "+":  "PLUS",
    "-":  "MINUS",
    "*":  "ASTERISK",
    "/":  "SLASH",
    ",":  "COMMA",
    ".":  "DOT",
    NULL: "EOF",
  }

  this.read = fn() {
    pos = pos + 1
    if pos >= len(input) {
      ch = null
    } else {
      ch = input[pos]
    }
  }

  this.peek = fn() {
    let next = pos + 1
    if next >= len(input) {
      return null
    } else {
      return input[next]
    }
  }

  this.whitespace = fn() {
    while ch == " " || ch == "\t" || ch == "\n" || ch == "\r" {
      this.read()
    }
  }

  this.next = fn() {
    this.whitespace()

    if bytetokens[ch] != NULL {
      let tok = bytetokens
    }
  }

  this.char = fn() { return ch }

  this.done = fn() { return ch == null }

  return this;
}

let lex = NewLexer(read("std.monkey"))

while !lex.done() {
  lex.read()
  print(lex.char(), lex.peek())
}

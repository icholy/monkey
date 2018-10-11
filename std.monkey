
let TokenType = {
	"ILLEGAL":   "ILLEGAL",
	"EOF":       "EOF",
	"IDENT":     "IDENT",
	"INT":       "INT",
	"ASSIGN":    "ASSIGN",
	"PLUS":      "PLUS",
	"MINUS":     "MINUS",
	"BANG":      "BANG",
	"ASTERISK":  "ASTERISK",
	"SLASH":     "SLASH",
	"GT":        "GT",
	"LT":        "LT",
	"EQ":        "EQ",
	"NE":        "NE",
	"GT_EQ":     "GT_EQ",
	"LT_EQ":     "LT_EQ",
	"DOT":       "DOT",
	"OR":        "OR",
	"AND":       "AND",
	"COMMA":     "COMMA",
	"SEMICOLON": "SEMICOLON",
	"COLON":     "COLON",
	"LPAREN":    "LPAREN",
	"RPAREN":    "RPAREN",
	"LBRACE":    "LBRACE",
	"RBRACE":    "RBRACE",
	"LBRACKET":  "LBRACKET",
	"RBRACKET":  "RBRACKET",
	"STRING":    "STRING",
	"FN":        "FN",
	"FUNCTION":  "FUNCTION",
	"LET":       "LET",
	"TRUE":      "TRUE",
	"FALSE":     "FALSE",
	"IF":        "IF",
	"ELSE":      "ELSE",
	"RETURN":    "RETURN",
	"IMPORT":    "IMPORT",
	"WHILE":     "WHILE",
	"PACKAGE":   "PACKAGE",
	"DEBUGGER":  "DEBUGGER",
	"NULL":      "NULL",
}

function NewToken(type, text) {

  let this = {
    "type": type,
    "text": text,
  };

  return this;
}

function NewLexer(input) {

  let pos = 0;
  let ch = input[0];
  let this = {};

  let simpletokens = {
    ";":  TokenType.SEMICOLON,
    ":":  TokenType.COLON,
    "(":  TokenType.LPAREN,
    ")":  TokenType.RPAREN,
    "{":  TokenType.LBRACE,
    "}":  TokenType.RBRACE,
    "[":  TokenType.LBRACKET,
    "]":  TokenType.RBRACKET,
    "+":  TokenType.PLUS,
    "-":  TokenType.MINUS,
    "*":  TokenType.ASTERISK,
    "/":  TokenType.SLASH,
    ",":  TokenType.COMMA,
    ".":  TokenType.DOT,
    null: TokenType.EOF,
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
    let tok = {}
    this.whitespace()

    if simpletokens[ch] != null {
      let type = simpletokens[ch]
      tok = NewToken(type, ch)
      this.read()
      return tok
    }

    if ch == "<" {
      if this.peek() == "=" {
        this.read()
        tok.type = TokenType.GT_EQ
        tok.text = "<="
      } else {
        tok = this.charTok(TokenType.GT)
      }
    } else {
      tok = this.charTok(TokenType.ILLEGAL)
    }

    this.read()
    return tok
  }

  this.charTok = fn(type) { return NewToken(type, ch) }

  this.char = fn() { return ch }

  this.done = fn() { return ch == null }

  return this;
}

let source = read("std.monkey")
let lex = NewLexer(source)

while !lex.done() {
  print(lex.next())
}


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
  "IN":        "IN",
  "SWITCH":    "SWITCH",
  "CASE":      "CASE",
  "DEFAULT":   "DEFAULT"
}

function NewToken(type, text) {

  let this = {
    "type": type,
    "text": text,
  };

  return this;
}

function NewSet(array) {
  let set = {};
  let i = 0;
  while i < len(array) {
    set[array[i]] = true
    i = i + 1
  }
  return set
}

function NewLexer(input) {

  let pos = 0;
  let ch = input[0];
  let this = {};

  let numbers = NewSet(["1", "2", "3", "4", "5", "6", "7", "8", "9"])

  let letters = NewSet([
    "a", "b", "c", "d", "e", "f", 
    "g", "h", "i", "j", "k", "l",
    "m", "n", "o", "p", "q", "r",
    "s", "t", "u", "v", "w", "x",
    "y", "z", "A", "B", "C", "D",
    "E", "F", "J", "K", "L", "M",
    "N", "O", "P", "Q", "R", "S",
    "T", "U", "V", "W", "X", "Y",
    "Z"
  ])

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
    null: TokenType.EOF
  }

  let keywords = {
    "fn":       TokenType.FN,
    "let":      TokenType.LET,
    "true":     TokenType.TRUE,
    "false":    TokenType.FALSE,
    "if":       TokenType.IF,
    "else":     TokenType.ELSE,
    "return":   TokenType.RETURN,
    "function": TokenType.FUNCTION,
    "import":   TokenType.IMPORT,
    "while":    TokenType.WHILE,
    "package":  TokenType.PACKAGE,
    "debugger": TokenType.DEBUGGER,
    "null":     TokenType.NULL,
    "in":       TokenType.IN,
    "switch":   TokenType.SWITCH,
    "case":     TokenType.CASE,
    "default":  TokenType.DEFAULT
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

    if ch in simpletokens {
      let type = simpletokens[ch]
      tok = NewToken(type, ch)
      this.read()
      return tok
    }

    switch ch {
    case "<":
      if this.peek() == "=" {
        this.read()
        tok.type = TokenType.LT_EQ
        tok.text = "<="
      } else {
        tok = this.charTok(TokenType.LT)
      }
    case ">":
      if this.peek() == "=" {
        this.read()
        tok.type = TokenType.GT_EQ
        tok.text = ">="
      } else {
        tok = this.charTok(TokenType.GT)
      }
    case "=":
      if this.peek() == "=" {
        this.read()
        tok.type = TokenType.EQ
        tok.text = "=="
      } else {
        tok = this.charTok(TokenType.ASSIGN)
      }
    case "!":
      if this.peek() == "=" {
        this.read()
        tok.type = TokenType.NE
        tok.text = "!="
      } else {
        tok = this.charTok(TokenType.BANG)
      }
    case "|":
      if this.peek() == "|" {
        this.read()
        tok.type = TokenType.OR
        tok.text = "||"
      } else {
        tok = this.charTok(TokenType.ILLEGAL)
      }
    case "&":
      if this.peek() == "&" {
        this.read()
        tok.type = TokenType.AND
        tok.text = "&&"
      } else {
        tok = this.charTok(TokenType.ILLEGAL)
      }
    case "\"":
      tok.text = this.str()
      tok.type = TokenType.STRING
      return tok
    default:
      if ch in letters {
        let ident = this.ident()
        tok.text = ident 
        if ident in keywords {
          tok.type = keywords[ident]
        } else {
          tok.type = TokenType.IDENT
        }
        return tok
      }
      if ch in numbers {
        tok.text = this.integer()
        tok.type = TokenType.INT
        return tok
      }
      tok = this.charTok(TokenType.ILLEGAL)
    }

    this.read()
    return tok
  }
  this.str = fn() {
    this.read()
    let escaped = false
    let str = ""
    while !this.done() {
      if escaped {
        switch ch {
        case "t":
          str = str + "\t"
        case "r":
          str = str + "\r"
        case "n":
          str = str + "\n"
        default:
          str = str + ch
        }
        escaped = false
      } else {
        if ch == "\"" {
          this.read()
          return str
        }
        if ch == "\\" {
          escaped = true
        } else {
          str = str + ch
        }
      }
      this.read()
    }
    this.read()
    return str
  }

  this.ident = fn() {
    let ident = "";
    while ch in letters || ch in numbers || ch == "_" {
      ident = ident + ch
      this.read()
    }
    return ident;
  }

  this.integer = fn() {
    let integer = ""
    while ch in numbers {
      integer = integer + ch
      this.read()
    }
    return integer
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

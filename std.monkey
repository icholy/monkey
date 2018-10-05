

function NewLexer(input) {

  let pos = 0;
  let ch = input[0];
  let export = {};

  export.next = fn() {
    pos = pos + 1;
    if (pos >= len(input)) {
      ch = "";
    } else {
      ch = input[pos];
    }
  };

  export.char = fn() {
    return ch;
  }

  export.done = fn() {
    return ch == "";
  }

  return export;
}

let lex = NewLexer(read("std.monkey"))

while (!lex.done()) {
  print(lex.char())
  lex.next()
}
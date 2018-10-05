

function empty(array) {
  return len(array) == 0
}

function reduce(array, f, acc) {
  let inner = fn(array, acc) {
    if (empty(array)) {
      return acc;
    }
    return inner(rest(array), f(first(array), acc))
  }
  return inner(array, acc)
}

function map(array, f) {
  reduce(array, fn(x, acc) { append(acc, f(x)) }, [])
}

function foreach(array, callback) {
  reduce(array, fn(x, _) { callback(x) }, 0)
  return;
}

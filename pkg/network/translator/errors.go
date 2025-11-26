package translator

import "errors"

var ErrNotFound = errors.New("translator: not found")
var ErrBadCast = errors.New("translator: bad cast")

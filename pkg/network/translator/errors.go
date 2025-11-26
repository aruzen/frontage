package translator

import "errors"

var ErrNotFound = errors.New("translator: not found")
var ErrBadCast = errors.New("translator: bad cast")
var ErrNewContextFailed = errors.New("translator: new context failed")

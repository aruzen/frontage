package translator

import "errors"

var ErrNotFound = errors.New("translator: not found")
var ErrBadCast = errors.New("translator: bad cast")
var ErrNewContextFailed = errors.New("translator: new context failed")
var ErrNewStateFailed = errors.New("translator: new state failed")
var ErrNilCardRepository = errors.New("translator: nil card repository")
var ErrNilMaterialsTranslator = errors.New("translator: nil materials translator")

package controller

import "errors"

var ErrUnsupportedPacketTag = errors.New("unsupported packet tag")
var ErrMissingHandler = errors.New("missing handler")

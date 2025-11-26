package repository

import "errors"

var ErrNotFound = errors.New("repository: not found")
var ErrDuplicateTag = errors.New("repository: duplicate tag")

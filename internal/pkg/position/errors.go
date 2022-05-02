package position

import "github.com/pkg/errors"

// ErrInstrumentMismatch indicates that the trade instrument does not match the expected value.
var ErrInstrumentMismatch error = errors.New("instrument mismatch")

// ErrNotSorted indicates that the trades are not sorted by timestamp.
var ErrNotSorted error = errors.New("not sorted")

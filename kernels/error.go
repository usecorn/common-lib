package kernels

import (
	"strings"

	"github.com/cockroachdb/errors"
)

var (
	ErrMissingStart         = errors.New("must have either startBlock or startTime")
	ErrNegativeRate         = errors.New("earn rate must be non-negative")
	ErrInvalidUserAddr      = errors.New("invalid user address")
	ErrEmptySource          = errors.New("source cannot be empty")
	ErrEmptySubSource       = errors.New("subSource cannot be empty")
	ErrEmptyCategory        = errors.New("category cannot be empty")
	ErrNegativeAmount       = errors.New("amount must be positive")
	ErrMissingGrantTime     = errors.New("grant time must be set")
	ErrEarnInf              = errors.New("earn rate cannot be infinite")
	ErrNonPostiveMultiplier = errors.New("multiplier must be positive")
	ErrNonPostiveStartBlock = errors.New("start block must be positive")
	ErrNonPostiveStartTime  = errors.New("start time must be positive")
	ErrInvalidEarnRate      = errors.New("invalid earn rate")
	ErrEmptyBatch           = errors.New("batch cannot be empty")
)

func IsErrTooOld(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "update starting_at to a value less than the previous starting_at")
}

type KernelError struct {
	Err string `json:"error"`
}

func (ke KernelError) Error() string {
	return ke.Err
}

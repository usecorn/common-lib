package kernels

import "fmt"

var (
	ErrMissingStart         = fmt.Errorf("must have either startBlock or startTime")
	ErrNegativeRate         = fmt.Errorf("earn rate must be non-negative")
	ErrInvalidUserAddr      = fmt.Errorf("invalid user address")
	ErrEmptySource          = fmt.Errorf("source cannot be empty")
	ErrEmptySubSource       = fmt.Errorf("subSource cannot be empty")
	ErrEmptyCategory        = fmt.Errorf("category cannot be empty")
	ErrNegativeAmount       = fmt.Errorf("amount must be positive")
	ErrMissingGrantTime     = fmt.Errorf("grant time must be set")
	ErrEarnInf              = fmt.Errorf("earn rate cannot be infinite")
	ErrNonPostiveMultiplier = fmt.Errorf("multiplier must be positive")
	ErrNonPostiveStartBlock = fmt.Errorf("start block must be positive")
	ErrNonPostiveStartTime  = fmt.Errorf("start time must be positive")
	ErrInvalidEarnRate      = fmt.Errorf("invalid earn rate")
	ErrEmptyBatch           = fmt.Errorf("batch cannot be empty")
)

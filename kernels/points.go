package kernels

import "github.com/pkg/errors"

type PointsEarnRequestFullBatch struct {
	EarnRequestFullBatch
	Program int64 `json:"program"`
}

func (b PointsEarnRequestFullBatch) IsPerBlock() bool {
	return false
}

func (b PointsEarnRequestFullBatch) Validate() error {
	if err := b.EarnRequestFullBatch.Validate(); err != nil {
		return err
	}
	if b.Program < 0 {
		return errors.New("program must be non-negative")
	}
	return nil
}

type PointsEarnRequestBatch struct {
	EarnRequestBatch
	Program int64 `json:"program"`
}

func (b PointsEarnRequestBatch) IsPerBlock() bool {
	return false
}

func (b PointsEarnRequestBatch) Validate() error {
	if err := b.EarnRequestBatch.Validate(); err != nil {
		return err
	}
	if b.Program < 0 {
		return errors.New("program must be non-negative")
	}
	return nil
}

type PointsEarnRequest struct {
	EarnRequest
	Program int64 `json:"program"`
}

func (b PointsEarnRequest) IsPerBlock() bool {
	return false
}

func (b PointsEarnRequest) Validate() error {
	if err := b.EarnRequest.Validate(); err != nil {
		return err
	}
	if b.Program < 0 {
		return errors.New("program must be non-negative")
	}
	return nil
}

package service

import (
	"context"
	"fmt"

	"github.com/callstats-io/ai-decision/service/src/grpc"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func validatePositiveInt(field string, val int32) error {
	if val <= 0 {
		return fmt.Errorf("%s: must be a positive integer", field)
	}
	return nil
}

func validateNonEmptyString(field string, val string) error {
	if val == "" {
		return fmt.Errorf("%s: cannot be empty", field)
	}
	return nil

}
func validateNonEmptyBytes(field string, data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("%s: cannot be empty", field)
	}
	return nil

}
func validateTimestamp(field string, gt *timestamp.Timestamp) error {
	if gt == nil {
		return fmt.Errorf("%s: cannot be nil", field)
	}
	if gt.Seconds <= 0 {
		return fmt.Errorf("%s: must have positive seconds", field)
	}
	return nil

}

// validate all errors are nil or return first error
func validate(ctx context.Context, errors ...error) error {
	for _, err := range errors {
		if err != nil {
			return grpc.ErrInvalidArgument(ctx, err)
		}
	}
	return nil
}

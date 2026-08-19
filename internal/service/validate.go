package service

import (
	"errors"
	"fmt"
	"rkv/internal/constants"
)

func validateKey(key string) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}
	if len(key) > constants.MaxKeySize {
		return fmt.Errorf("key exceeds max size of %d bytes",
			constants.MaxKeySize)
	}
	return nil
}

func validateValue(value []byte) error {
	if len(value) == 0 {
		return errors.New("value cannot be empty")
	}
	if len(value) > constants.MaxValueSize {
		return fmt.Errorf("value exceeds max size of %d bytes",
			constants.MaxValueSize)
	}
	return nil
}

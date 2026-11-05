//go:build !linux

package collector

import (
	"context"
	"fmt"
)

// Collect is unavailable on non-Linux systems until platform-specific support exists.
func Collect(context.Context) ([]string, error) {
	return nil, fmt.Errorf("local snapshot is currently supported only on Linux")
}

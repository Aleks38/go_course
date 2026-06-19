package api

import (
	"fmt"
	"time"
)

func newBatchID() string {
	return fmt.Sprintf("b_%d", time.Now().UnixNano())
}

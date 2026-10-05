package agent

import (
	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"time"
)

// optionalUnixMilli matches the Proxy mapper's absent/zero timestamp semantics.
func optionalUnixMilli(value *time.Time) *int64 {
	if value == nil {
		return nil
	}
	return pointy.GetUnixMilliPointer(value.UnixMilli())
}

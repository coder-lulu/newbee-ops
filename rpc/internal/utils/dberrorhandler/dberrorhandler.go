package dberrorhandler

import (
    "github.com/zeromicro/go-zero/core/logx"
)

// DefaultEntError provides a minimal error adapter for ent operations.
// In production, map to i18n/business errors; for now, just log and return original.
func DefaultEntError(logger logx.Logger, err error, _ interface{}) error {
    if err != nil {
        logger.Errorf("ent operation error: %v", err)
    }
    return err
}


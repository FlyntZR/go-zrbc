package wschannel

import (
	"fmt"

	"go-zrbc/pkg/xlog"
)

type Logger struct {
	ServiceID string `json:"service_id"`
	UUID      string `json:"uuid"`
	UserID    string `json:"u_id"`
}

func (logger *Logger) Prefix() string {
	return fmt.Sprintf("service_id(%s), uuid(%s), uid(%s)", logger.ServiceID, logger.UUID, logger.UserID)
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	xlog.ErrorfN(2, "%s: %s", logger.Prefix(), fmt.Sprintf(format, args...))
}

func (logger *Logger) Debugf(format string, args ...interface{}) {
	xlog.DebugfN(2, "%s: %s", logger.Prefix(), fmt.Sprintf(format, args...))
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	xlog.InfofN(2, "%s: %s", logger.Prefix(), fmt.Sprintf(format, args...))
}

func (logger *Logger) Error(args ...interface{}) {
	xlog.ErrorfN(2, "%s: %s", logger.Prefix(), fmt.Sprintf("%v", args))
}

func (logger *Logger) Debug(args ...interface{}) {
	xlog.DebugN(2, "%s: %s", logger.Prefix(), fmt.Sprintf("%v", args))
}

func (logger *Logger) Info(args ...interface{}) {
	xlog.InfofN(2, "%s: %s", logger.Prefix(), fmt.Sprintf("%v", args))
}

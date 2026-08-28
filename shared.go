package shared

import (
	"github.com/go-log4g/core"
	"github.com/go-log4g/shared/impl"
)

func init() {
	core.RegisterAppender("sharedRollingFile", impl.BuildSharedRollingFileAppender)
}

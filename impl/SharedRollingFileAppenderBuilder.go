package impl

import (
	"github.com/go-jang/go/lang"
	core "github.com/go-log4g/core/impl"
	"github.com/go-log4g/shared/impl/model"
)

func BuildSharedRollingFileAppender(name string, pending *core.PendingAppender) core.Appender {
	var definition model.SharedRollingFileAppenderDefinition
	lang.Assert(pending.Definition.Node.Decode(&definition) == nil, "cannot parse shared rolling file appender %q", name)

	config := pending.Builder.BuildRollingFileAppenderConfig(name, definition.AppenderDefinition)

	return NewSharedRollingFileAppender(
		config.File,
		config.FilePattern,
		config.Append,
		config.BufferSize,
		config.ImmediateFlush,
		config.Layout,
		config.Filter,
		config.StatusLogger,
		config.StartupPolicy,
		config.TriggerPolicy,
		config.Strategy,
		definition.FileLock,
	)
}

package impl

import (
	"fmt"

	"github.com/go-errr/go/err"
	core "github.com/go-log4g/core/impl"
	"github.com/go-log4g/core/impl/model"
	sharedmodel "github.com/go-log4g/shared/impl/model"
)

type SharedRollingFileAppender struct {
	prefix string
}

func NewSharedRollingFileAppender(name string, dynamicDefinition *model.DynamicDefinition) core.Appender {
	var definition sharedmodel.SharedRollingFileAppenderDefinition
	err.Assert(dynamicDefinition.Node.Decode(&definition), "cannot parse sharedRollingFile appender %q", name)

	return &SharedRollingFileAppender{
		prefix: name + ": ",
	}
}

// Implements core.Appender
func (this *SharedRollingFileAppender) Append(event *core.LogEvent) {
	fmt.Println(this.prefix + event.Record.Message)
}

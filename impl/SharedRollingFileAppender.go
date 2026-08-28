package impl

import (
	"os"
	"path/filepath"
	"time"

	"github.com/go-errr/go/err"
	filelock "github.com/go-jang/file-lock"
	"github.com/go-jang/go/lang"
	"github.com/go-jang/go/util/optional"
	core "github.com/go-log4g/core/impl"
	"github.com/go-log4g/core/impl/filter"
	"github.com/go-log4g/core/impl/rolling"
)

type SharedRollingFileAppender struct {
	*core.AbstractFileAppender

	file        string
	filePattern string
	fileLock    *filelock.FairFileLock
	policy      rolling.TriggeringPolicy
	strategy    rolling.RolloverStrategy
}

func NewSharedRollingFileAppender(file string, filePattern string, append bool, bufferSize int, immediateFlush bool, layout core.Layout, filter filter.Filter, statusLogger *core.StatusLogger,
	startupPolicy *rolling.OnStartupTriggeringPolicy, triggerPolicy rolling.TriggeringPolicy, strategy rolling.RolloverStrategy, fileLock string) *SharedRollingFileAppender {

	absoluteFile, e := filepath.Abs(file)
	err.Assert(e, "cannot resolve shared log file %q", file)

	err.Assert(os.MkdirAll(filepath.Dir(absoluteFile), 0755), "cannot create shared log directory for %q", absoluteFile)

	if fileLock == "" {
		fileLock = filepath.Join(filepath.Dir(absoluteFile), ".lock")
	}

	lockFactory := filelock.NewFactory(fileLock)

	result := &SharedRollingFileAppender{
		file:        absoluteFile,
		filePattern: filePattern,
		fileLock:    lockFactory.NewFairLock(absoluteFile),
		policy:      triggerPolicy,
		strategy:    strategy,
	}

	result.initialize(append, startupPolicy)
	result.AbstractFileAppender = core.NewAbstractFileAppender(absoluteFile, bufferSize, immediateFlush, layout, filter, statusLogger, result)

	return result
}

func (this *SharedRollingFileAppender) initialize(append bool, startupPolicy *rolling.OnStartupTriggeringPolicy) {
	this.fileLock.Lock()
	defer this.fileLock.Unlock()

	info, e := os.Stat(this.file)
	if e != nil && !os.IsNotExist(e) {
		err.Assert(e, "cannot stat shared log file %q", this.file)
	}

	if e == nil && startupPolicy != nil && startupPolicy.IsTriggered(info.Size()) {
		this.strategy.Rollover(this.file, this.filePattern, info.ModTime())
	}

	if append {
		return
	}

	handle, e := os.OpenFile(this.file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	err.Assert(e, "cannot truncate shared log file %q", this.file)
	err.Assert(handle.Close(), "cannot close shared log file %q", this.file)
}

// Implements core.FileAppenderDelegate
func (this *SharedRollingFileAppender) BeforeAppend(data []byte) {
}

// Implements core.FileAppenderDelegate
func (this *SharedRollingFileAppender) Write(data []byte) {
	this.fileLock.Lock()
	defer this.fileLock.Unlock()

	fileTime, size := this.fileState()
	now := time.Now()

	context := rolling.TriggeringContext{
		FileSize:  size,
		EventSize: len(data),
		Time:      now,
	}

	if size > 0 && this.policy != nil && this.policy.IsTriggered(context) {
		this.strategy.Rollover(this.file, this.filePattern, fileTime)
	}

	handle, e := os.OpenFile(this.file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	err.Assert(e, "cannot open shared log file %q", this.file)
	defer func() {
		err.Assert(handle.Close(), "cannot close shared log file %q", this.file)
	}()

	written := optional.OfCommaErr(handle.Write(data)).OrElsePanic("cannot write data to shared log file %q", this.file)
	lang.Assert(written == len(data), "cannot write complete data to shared log file %q", this.file)
}

// Implements core.FileAppenderDelegate
func (this *SharedRollingFileAppender) AfterAppend(data []byte) {
}

func (this *SharedRollingFileAppender) fileState() (time.Time, int64) {
	info, e := os.Stat(this.file)
	if os.IsNotExist(e) {
		return time.Now(), 0
	}

	err.Assert(e, "cannot stat shared log file %q", this.file)
	return info.ModTime(), info.Size()
}

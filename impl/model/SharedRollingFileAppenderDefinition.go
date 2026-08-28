package model

import coreModel "github.com/go-log4g/core/impl/model"

type SharedRollingFileAppenderDefinition struct {
	coreModel.AppenderDefinition `yaml:",inline"`

	FileLock string `yaml:"fileLock"`
}

package model

type SharedRollingFileAppenderDefinition struct {
	File           string `yaml:"file"`
	FilePattern    string `yaml:"filePattern"`
	Append         *bool  `yaml:"append"`
	BufferSize     *int   `yaml:"bufferSize"`
	ImmediateFlush *bool  `yaml:"immediateFlush"`
}

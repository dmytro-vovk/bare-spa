package boot

import (
	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Logger    Logger `json:"logger"`
	Webserver struct {
		Listen string `json:"listen"`
	} `json:"webserver"`
	Storage struct {
		File string `json:"file"`
	} `json:"storage"`
	Auth struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	} `json:"auth"`
	UART map[string]struct {
		Device string `json:"device"`
		Config string `json:"config"`
	} `json:"uart"`
	SystemBoard struct {
		Type string `json:"type"`
	} `json:"system_board"`
}

type Logger struct {
	Level            string            `json:"level"`
	TimestampFormat  string            `json:"timestamp_format"`
	RuntimeFormatter *runtimeFormatter `json:"formatter"`
	Rotor            *rotor            `json:"rotor"`
}

type runtimeFormatter struct {
	Line         bool `json:"line"`
	Package      bool `json:"package"`
	File         bool `json:"file"`
	BaseNameOnly bool `json:"base_name_only"`
}

type rotor struct {
	Filename   string `json:"filename"`
	MaxSize    int    `json:"max_size"`
	MaxAge     int    `json:"max_age"`
	MaxBackups int    `json:"max_backups"`
	LocalTime  bool   `json:"local_time"`
	Compress   bool   `json:"compress"`
}

func (l *Logger) Formatter(child log.Formatter) log.Formatter {
	f := l.RuntimeFormatter
	if f == nil || !(f.Line || f.Package || f.File) {
		return child
	}

	return &runtime.Formatter{
		ChildFormatter: child,
		Line:           f.Line,
		Package:        f.Package,
		File:           f.File,
		BaseNameOnly:   f.BaseNameOnly,
	}
}

func (l *Logger) FieldMap() log.FieldMap {
	return log.FieldMap{
		log.FieldKeyMsg: "message",
	}
}

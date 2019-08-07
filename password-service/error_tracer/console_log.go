package error_tracer

import log "github.com/sirupsen/logrus"

type ConsoleLog struct {
}

func (c *ConsoleLog) InfoLog(api string, tag string, message string) {
	log.Infof("%s\t==/\t%s\t==/\t%s\n", api, tag, message)
}

func (c *ConsoleLog) ErrorLog(api string, tag string, message string) {
	log.Warnf("%s\t==/\t%s\t==/\t%s\n", api, tag, message)
}

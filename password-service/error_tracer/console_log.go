package error_tracer

import log "github.com/sirupsen/logrus"

type ConsoleLog struct {
}

func (c *ConsoleLog) ErrorLog(api string, tag string, message string) {
	log.Warnf("%s\t==/\t%s\t==/\t==/\t\n", api, tag, message)
}

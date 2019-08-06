package error_tracer

type IErrorTracer interface {
	ErrorLog(api string, tag string, message string)
}

var Client IErrorTracer

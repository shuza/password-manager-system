package error_tracer

type IErrorTracer interface {
	InfoLog(api string, tag string, message string)
	ErrorLog(api string, tag string, message string)
}

var Client IErrorTracer

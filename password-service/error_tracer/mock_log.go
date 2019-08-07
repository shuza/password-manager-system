package error_tracer

type MockLog struct {
}

func (c *MockLog) InfoLog(api string, tag string, message string) {

}

func (c *MockLog) ErrorLog(api string, tag string, message string) {

}

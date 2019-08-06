package error_tracer

type MockLog struct {
}

func (c *MockLog) ErrorLog(api string, tag string, message string) {

}

//go:build examplemain
// +build examplemain

package main

type Logger struct{}

func (l *Logger) Log(items ...any) {}

func (l *Logger) Debug(items ...any) {}

func main() {
	logger := &Logger{}
	_, err := fmt.Println("hello world")

	if err != nil {
		logger.Log(err)
	} else {
		logger.Debug("debug log")
	}
}

package logger

import "fmt"

type Logger struct {
	logFile  string
	logLevel string
}

func New(fileName string, level string) *Logger {
	return &Logger{logFile: fileName, logLevel: level}
}

func (l Logger) Info(msg string) {
	fmt.Println(msg)
}

func (l Logger) Error(msg string) {
	fmt.Println(msg)

}

// TODO

//* log_file - путь к файлу логов;
//* log_level - уровень логирования (error / warn / info / debug);

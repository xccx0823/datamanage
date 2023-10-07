package log

func Info(args ...interface{}) {
	logger.Infoln(args...)
}
func InfoF(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Error(args ...interface{}) {
	logger.Errorln(args...)
}
func ErrorF(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Warn(args ...interface{}) {
	logger.Warnln(args...)
}
func WarnF(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

type Standard interface {
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatalln(args ...interface{})

	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Panicln(args ...interface{})

	Print(args ...interface{})
	Printf(format string, args ...interface{})
	Println(args ...interface{})
}

package logging

func Info(arg ...any) {
	instance.print("::info::", arg...)
}

func Infof(format string, arg ...any) {
	instance.printf("::info::", format, arg...)
}

func Debug(msg ...any) {
	instance.print("::debug::", msg...)
}

func Debugf(format string, arg ...any) {
	instance.printf("::debug::", format, arg...)
}

func Notice(msg ...any) {
	instance.print("::notice::", msg...)
}

func Noticef(format string, arg ...any) {
	instance.printf("::notice::", format, arg...)
}

func Warning(msg ...any) {
	instance.print("::warning::", msg...)
}

func Warningf(format string, arg ...any) {
	instance.printf("::warning::", format, arg...)
}

func Error(msg ...any) {
	instance.print("::error::", msg...)
}

func Errorf(format string, arg ...any) {
	instance.printf("::error::", format, arg...)
}

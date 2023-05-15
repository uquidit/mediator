package logger

import "log"

type FileLogger struct {
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
}

func NewFileLogger(filename string, append bool) (*FileLogger, error) {

	if out, err := openLogfile(filename, append); err != nil {
		return nil, err
	} else {
		l := FileLogger{
			infoLogger:    log.New(out, "[INFO] ", log.LstdFlags),
			warningLogger: log.New(out, "[WARNING] ", log.LstdFlags),
			errorLogger:   log.New(out, "[ERROR] ", log.LstdFlags),
		}
		return &l, nil
	}
}

func (l *FileLogger) Infof(format string, a ...any) {
	l.infoLogger.Printf(format, a...)
}
func (l *FileLogger) Warningf(format string, a ...any) {
	l.warningLogger.Printf(format, a...)
}
func (l *FileLogger) Errorf(format string, a ...any) {
	l.errorLogger.Printf(format, a...)
}

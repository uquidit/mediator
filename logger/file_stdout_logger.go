package logger

import (
	"fmt"
	"log"
	"path/filepath"
)

type FileStdoutLogger struct {
	fileWarningLogger *log.Logger
	fileInfoLogger    *log.Logger
	fileErrorLogger   *log.Logger

	stdOutInfoLogger    *log.Logger
	stdErrWarningLogger *log.Logger
	stdErrErrorLogger   *log.Logger

	is_init bool
}

func GetRootFileStdoutLogger() (*FileStdoutLogger, error) {
	// if root logger is not initialized, return an error
	if !is_init {
		return nil, fmt.Errorf("root logger not initialized")
	}

	//make a FileStdoutLogger that uses root loggers
	l := FileStdoutLogger{
		fileWarningLogger:   fileWarningLogger,
		fileInfoLogger:      fileInfoLogger,
		fileErrorLogger:     fileErrorLogger,
		stdOutInfoLogger:    stdOutInfoLogger,
		stdErrWarningLogger: stdErrWarningLogger,
		stdErrErrorLogger:   stdErrErrorLogger,
		is_init:             true,
	}
	return &l, nil
}

func NewFileStdoutLogger(fullfilename string, append bool) (*FileStdoutLogger, error) {

	if fullfilename == "" {
		return nil, fmt.Errorf("cannot init logger: no file name")
	}

	if absFileName, err := filepath.Abs(fullfilename); err != nil {
		return nil, err
	} else if file, err := openLogfile(absFileName, append); err != nil {
		return nil, err
	} else {
		l := FileStdoutLogger{}
		l.fileInfoLogger = log.New(file, "[INFO] ", log.LstdFlags)
		l.fileWarningLogger = log.New(file, "[WARNING] ", log.LstdFlags)
		l.fileErrorLogger = log.New(file, "[ERROR] ", log.LstdFlags)
		l.is_init = true
		return &l, nil
	}
}

func (l *FileStdoutLogger) Infof(format string, a ...any) {
	if l.is_init {
		if l.fileInfoLogger != nil {
			l.fileInfoLogger.Printf(format, a...)
		}
		if l.stdOutInfoLogger != nil {
			l.stdOutInfoLogger.Printf(format, a...)
		}
	} else {
		log.Fatalln("logger not initialized")
	}
}

func (l *FileStdoutLogger) Warning(err error) {
	l.Warningf(err.Error() + "\n")
}

func (l *FileStdoutLogger) Warningf(format string, a ...any) {
	if l.is_init {
		if l.fileWarningLogger != nil {
			l.fileWarningLogger.Printf(format, a...)
		}
		if l.stdErrWarningLogger != nil {
			l.stdErrWarningLogger.Printf(format, a...)
		}
	} else {
		log.Fatalln("logger not initialized")
	}
}

func (l *FileStdoutLogger) Error(err error) {
	l.Errorf(err.Error() + "\n")
}

func (l *FileStdoutLogger) Errorf(format string, a ...any) {
	if l.is_init {
		if l.fileErrorLogger != nil {
			l.fileErrorLogger.Printf(format, a...)
		}
		if l.stdErrErrorLogger != nil {
			l.stdErrErrorLogger.Fatalf(format, a...)
		}
	} else {
		log.Fatalln("logger not initialized")
	}
}

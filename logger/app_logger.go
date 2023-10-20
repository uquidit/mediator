package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var currentPid int
var logFile *os.File
var pidFlag bool
var functionFlag bool
var nanoSecondsFlag bool

func InitAppLoggerWithConfigFile(configFilePath string) error {
	//TODO define configuration file format
	//TODO parse configuration and call InitAppLogger(logInStdOut bool, logInFile bool, withPid bool, withFunction bool, appendFile bool, folder string, filename string)
	return nil
}

func InitAppLogger(logLevel logrus.Level, logInStdOut bool, logInFile bool, withPid bool, withFunction bool, appendFile bool, withNanoSeconds bool, folder string, filename string) error {
	if !logInFile && !logInStdOut {
		return fmt.Errorf("cannot init logger: choose log in Stdout or log in file or both")
	}
	pidFlag = withPid
	functionFlag = withFunction
	nanoSecondsFlag = withNanoSeconds
	currentPid = os.Getpid()

	if logInFile {
		logfolder := ""
		if folder == "" && filename == "" {
			return fmt.Errorf("cannot init logger: no folder and file name")
		}
		if filename == "" {
			return fmt.Errorf("cannot init logger: no file name")
		}
		// Folder can be empty if filename is a full path file
		if folder == "" {
			folder = filepath.Dir(filename)
			filename = filepath.Base(filename)
		}

		logf, err := filepath.Abs(folder)
		if err != nil {
			return fmt.Errorf("cannot init logger: invalid folder path, %w", err)
		} else {
			logfolder = logf
		}
		file, err := getLogfile(logfolder, filename, appendFile)
		if err != nil {
			return fmt.Errorf("cannot init logger: invalid file, %w", err)
		}
		if logInStdOut {
			logrus.SetOutput(io.MultiWriter(os.Stderr, file))
		} else {
			logrus.SetOutput(file)
		}
	} else if logInStdOut {
		logrus.SetOutput(os.Stderr)
	}
	logrus.SetReportCaller(true)
	logrus.SetFormatter(new(DefaultLogFormatter))
	logrus.SetLevel(logLevel)
	return nil
}

type DefaultLogFormatter struct{}

func (f *DefaultLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	logMessage := "["
	if nanoSecondsFlag {
		logMessage += time.Now().Format(time.RFC3339Nano)
	} else {
		logMessage += time.Now().Format(time.RFC3339)
	}
	logMessage += "][" + strings.ToUpper(entry.Level.String()) + "]"
	if pidFlag {
		logMessage += "[" + strconv.Itoa(currentPid) + "]"
	}
	if functionFlag && entry.Caller != nil {
		logMessage += "[" + entry.Caller.Function + "." + strconv.Itoa(entry.Caller.Line) + "]"
	}
	logMessage += " " + entry.Message + "\n"
	return []byte(logMessage), nil
}

func CloseLogFile() {
	if logFile != nil {
		logFile.Sync()
		logFile.Close()
	}
}

func getLogfile(folder string, filename string, flag_append bool) (*os.File, error) {
	name := ""
	if !filepath.IsAbs(filename) {
		name = filepath.Clean(filepath.Join(folder + "/" + filename))
	} else {
		name = filepath.Clean(filename)
	}
	flags := os.O_CREATE | os.O_WRONLY
	if flag_append {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}

	file, err := os.OpenFile(name, flags, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}

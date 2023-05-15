package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	fileWarningLogger *log.Logger
	fileInfoLogger    *log.Logger
	fileErrorLogger   *log.Logger

	stdOutInfoLogger    *log.Logger
	stdErrWarningLogger *log.Logger
	stdErrErrorLogger   *log.Logger

	is_init         bool = false
	logfolder       string
	logfileHandlers []*os.File
)

func openLogfile(filename string, flag_append bool) (*os.File, error) {
	name := ""

	// check if provided file name is absolute
	if !filepath.IsAbs(filename) {
		// it's not a full path. prepend log folder
		name = filepath.Clean(logfolder + "/" + filename)
	} else {
		name = filepath.Clean(filename)
	}

	// process flags
	flags := os.O_CREATE | os.O_WRONLY
	if flag_append {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}

	//open logfile
	if file, err := os.OpenFile(name, flags, 0666); err != nil {
		return nil, err
	} else {
		logfileHandlers = append(logfileHandlers, file)
		return file, nil
	}

}

func CloseAllLogfiles() {
	for _, file := range logfileHandlers {
		file.Sync()
		file.Close()
	}
}

func InitFullPath(fullfilename string, append bool, dumpfile bool, dumpstdout bool) error {
	if fullfilename == "" {
		return Init("", append, "", dumpfile, dumpstdout)
	} else {
		return Init(
			filepath.Base(fullfilename),
			append,
			filepath.Dir(fullfilename),
			dumpfile,
			dumpstdout,
		)
	}
}

func Init(filename string, append bool, folder string, dumpfile bool, dumpstdout bool) error {

	if dumpfile {
		if filename == "" {
			return fmt.Errorf("cannot init logger: no file name")
		}
		if folder == "" {
			return fmt.Errorf("cannot init logger: no folder name")
		}
		// check log folder
		if logf, err := filepath.Abs(folder); err != nil {
			return fmt.Errorf("cannot init logger: %w", err)
		} else {
			logfolder = logf // already cleant by Abs()
		}

		if file, err := openLogfile(filename, append); err != nil {
			return err
		} else {
			is_init = true
			fileInfoLogger = log.New(file, "[INFO] ", log.LstdFlags)
			fileWarningLogger = log.New(file, "[WARNING] ", log.LstdFlags)
			fileErrorLogger = log.New(file, "[ERROR] ", log.LstdFlags)
		}
	}

	if dumpstdout {
		is_init = true
		stdOutInfoLogger = log.New(os.Stderr, "[INFO] ", log.LstdFlags)
		stdErrWarningLogger = log.New(os.Stderr, "[WARNING] ", log.LstdFlags)
		stdErrErrorLogger = log.New(os.Stderr, "[ERROR] ", log.LstdFlags)
	}
	return nil

}

func Infof(format string, a ...any) {
	if is_init {
		if fileInfoLogger != nil {
			fileInfoLogger.Printf(format, a...)
		}
		if stdOutInfoLogger != nil {
			stdOutInfoLogger.Printf(format, a...)
		}
	} else {
		log.Fatalln("logger not initialized")
	}
}

func Warning(err error) {
	Warningf(err.Error() + "\n")
}

func Warningf(format string, a ...any) {
	if is_init {
		if fileWarningLogger != nil {
			fileWarningLogger.Printf(format, a...)
		}
		if stdErrWarningLogger != nil {
			stdErrWarningLogger.Printf(format, a...)
		}
	} else {
		log.Fatalln("logger not initialized")
	}
}

func Error(err error) {
	Errorf(err.Error() + "\n")
}

func Errorf(format string, a ...any) {
	if is_init {
		if fileErrorLogger != nil {
			fileErrorLogger.Printf(format, a...)
		}
		if stdErrErrorLogger != nil {
			stdErrErrorLogger.Fatalf(format, a...)
		}
	} else {
		log.Fatalln("logger not initialized")
	}
}

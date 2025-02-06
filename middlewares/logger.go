package middlewares

import (
	"io"
	"os"
	"path"
)

const LogPath = "/logs"
const LogFile = "access.log"

// LogWriter creates a log file at the specified path if it does not exist
// and returns an io.Writer that writes to both the log file and the standard output.
// It ensures that the log directory is created with appropriate permissions.

func LogWriter() io.Writer {
	_ = os.Mkdir(LogPath, 0770)
	logFilePath := path.Join(LogPath, LogFile)

	logFile, _ := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
	return io.MultiWriter(logFile, os.Stdout)
}

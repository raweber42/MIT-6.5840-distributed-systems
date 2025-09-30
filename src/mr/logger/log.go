package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Log is the exported Logrus instance. Other packages will import
// this variable to perform logging.
var Log *logrus.Logger

func init() {
	// 1. Create a new Logrus instance
	l := logrus.New()

	// 2. Configuration
	// Set the output to stdout (or a file)
	l.SetOutput(os.Stdout)

	l.SetFormatter(&logrus.TextFormatter{
		ForceColors: true, // Enable colors in the output

	})

	// 3. Assign the configured instance to the exported variable
	Log = l

	// Optional: Initial confirmation log
	Log.Debug("Logger initialized successfully.")
}

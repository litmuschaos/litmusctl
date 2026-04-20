/*
Copyright © 2021 The LitmusChaos Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package utils

import (
	"os"

	"github.com/litmuschaos/litmusctl/pkg/config"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()
	Log.SetOutput(os.Stdout)
	Log.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: false,
	})
	// Default to Error level; will be set to Debug when verbose mode is enabled
	Log.SetLevel(logrus.ErrorLevel)
}

// ConfigureVerboseLogging configures the logger based on verbose mode
func ConfigureVerboseLogging() {
	if config.VerboseMode {
		Log.SetLevel(logrus.DebugLevel)
		Log.Debug("Verbose mode enabled")
	} else {
		Log.SetLevel(logrus.ErrorLevel)
	}
}

// Debug logs a message at debug level (only shown in verbose mode)
func Debug(args ...interface{}) {
	if config.VerboseMode {
		Log.Debug(args...)
	}
}

// Debugf logs a formatted message at debug level (only shown in verbose mode)
func Debugf(format string, args ...interface{}) {
	if config.VerboseMode {
		Log.Debugf(format, args...)
	}
}

// Info logs a message at info level (only shown in verbose mode)
func Info(args ...interface{}) {
	if config.VerboseMode {
		Log.Info(args...)
	}
}

// Infof logs a formatted message at info level (only shown in verbose mode)
func Infof(format string, args ...interface{}) {
	if config.VerboseMode {
		Log.Infof(format, args...)
	}
}

// Warn logs a warning message (shown regardless of verbose mode)
func Warn(args ...interface{}) {
	Log.Warn(args...)
}

// Warnf logs a formatted warning message (shown regardless of verbose mode)
func Warnf(format string, args ...interface{}) {
	Log.Warnf(format, args...)
}

// Error logs an error message (shown regardless of verbose mode)
func Error(args ...interface{}) {
	Log.Error(args...)
}

// Errorf logs a formatted error message (shown regardless of verbose mode)
func Errorf(format string, args ...interface{}) {
	Log.Errorf(format, args...)
}

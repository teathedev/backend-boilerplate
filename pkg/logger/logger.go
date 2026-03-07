// Package logger provides a logger for the application.
package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/teathedev/backend-boilerplate/pkg/env"
	"github.com/teathedev/backend-boilerplate/pkg/utils"

	"github.com/fatih/color"
)

type colors func(format string, a ...any) string

type loggerColors struct {
	Trace   colors
	Info    colors
	Warning colors
	Error   colors
	Fatal   colors
}

var defaultColors = loggerColors{
	Trace:   color.CyanString,
	Info:    color.GreenString,
	Warning: color.YellowString,
	Error:   color.RedString,
	Fatal:   color.MagentaString,
}

type logLevels string

const (
	LOGLEVELTrace   logLevels = "TRACE"
	LOGLEVELInfo    logLevels = "INFO"
	LOGLEVELWarning logLevels = "WARNING"
	LOGLEVELError   logLevels = "ERROR"
	LOGLEVELFatal   logLevels = "FATAL"
)

type prodLog struct {
	App     string    `json:"app"`
	Module  string    `json:"module"`
	Level   logLevels `json:"level"`
	Message string    `json:"message"`
	Date    time.Time `json:"date"`
}

type Logger interface {
	Trace(string, ...LogParams)
	Info(string, ...LogParams)
	Warning(string, ...LogParams)
	Error(string, ...LogParams)
	Fatal(string, ...LogParams)
}

type logger struct {
	env        string
	dateFormat string
	logFormat  string
	app        string
	module     string
	colors     loggerColors
}

func (l *logger) getColor(level logLevels) colors {
	switch level {
	case LOGLEVELTrace:
		return l.colors.Trace
	case LOGLEVELInfo:
		return l.colors.Info
	case LOGLEVELWarning:
		return l.colors.Warning
	case LOGLEVELError:
		return l.colors.Error
	case LOGLEVELFatal:
		return l.colors.Fatal
	default:
		return color.GreenString
	}
}

type LogParams = map[string]any

func (l *logger) serializeLog(message string, level logLevels, args ...LogParams) string {
	now := time.Now()
	if l.env == "production" {
		// Create base log map with camelCase fields
		logMap := map[string]any{
			"app":     l.app,
			"module":  l.module,
			"level":   level,
			"message": message,
			"date":    now,
		}

		// Merge all LogParams directly into the log map
		for _, arg := range args {
			if reflect.TypeFor[LogParams]().Kind() == reflect.Map {
				maps.Copy(logMap, arg)
			} else if utils.IsStruct(reflect.TypeFor[LogParams]()) {
				// Handle structs by converting them to JSON and adding as "data"
				jsonData, err := json.Marshal(arg)
				if err != nil {
					logMap["structError"] = err.Error()
				} else {
					logMap["data"] = string(jsonData)
				}
			} else {
				// Handle other types as "extra"
				logMap["extra"] = arg
			}
		}

		msg, err := json.Marshal(logMap)
		if err != nil {
			fmt.Println(err)
		}

		return string(msg) + "\n"
	}

	color := l.getColor(level)
	msg := color("[%s] [%s]", l.app, level)
	temp := fmt.Sprintf(l.logFormat, now.Format(l.dateFormat), l.module)
	msg = fmt.Sprintf("%s %s", msg, temp)

	// Add the main message
	msg = fmt.Sprintf("%s [Message=%s]", msg, message)

	// Add additional arguments
	for _, arg := range args {
		if reflect.TypeFor[LogParams]().Kind() == reflect.Map {
			for k, v := range arg {
				msg = fmt.Sprintf("%s [%s=%v]", msg, k, v)
			}
		} else if utils.IsStruct(reflect.TypeFor[LogParams]()) {
			// Handle structs by converting them to JSON
			jsonData, err := json.Marshal(arg)
			if err != nil {
				msg = fmt.Sprintf("%s [%s=%v]", msg, "error", err)
			} else {
				msg = fmt.Sprintf("%s [%s=%s]", msg, "data", string(jsonData))
			}
		} else {
			// Handle other types
			msg = fmt.Sprintf("%s [%s=%v]", msg, "extra", arg)
		}
	}

	return msg + "\n"
}

func (l *logger) log(format string, level logLevels, args ...LogParams) {
	// Format the main message without including the extra arguments
	message := fmt.Sprint(format)

	// Serialize the log with the main message and additional arguments
	log := l.serializeLog(message, level, args...)
	io.WriteString(os.Stdout, log)
}

func (l *logger) Trace(format string, args ...LogParams) {
	l.log(format, LOGLEVELTrace, args...)
}

func (l *logger) Info(format string, args ...LogParams) {
	l.log(format, LOGLEVELInfo, args...)
}

func (l *logger) Warning(format string, args ...LogParams) {
	l.log(format, LOGLEVELWarning, args...)
}

func (l *logger) Error(format string, args ...LogParams) {
	l.log(format, LOGLEVELError, args...)
}

func (l *logger) Fatal(format string, args ...LogParams) {
	l.log(format, LOGLEVELFatal, args...)
	os.Exit(1)
}

var application string

func New(module string) *logger {
	env := strings.ToLower(env.GetString("GO_ENV", "development"))

	return &logger{
		app:        application,
		module:     module,
		dateFormat: "2006-01-02T15:04:05.00000000000Z07:00",
		colors:     defaultColors,
		logFormat:  "[%s] [%s]",
		env:        env,
	}
}

func init() {
	application = env.GetString("APP_NAME", "")

	if len(application) == 0 {
		fmt.Println("!!! APP_NAME is missing !!!")
	}
}

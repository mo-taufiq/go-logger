package gologger

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cast"
)

type LogConfiguration struct {
	TimeZone            string
	TimeFormat          string
	Path                string
	CreateLogFile       bool
	DebugMode           bool
	NestedLocationLevel int
	LogFuncName         bool
	NestedFuncLevel     int
	RuntimeCallerSkip   int
}

// Default config
var (
	LogConf = LogConfiguration{
		TimeZone:            "Asia/Jakarta",
		TimeFormat:          "2006-01-02T15:04:05-0700",
		CreateLogFile:       true,
		Path:                "./logs",
		DebugMode:           true,
		NestedLocationLevel: 1,
		LogFuncName:         true,
		NestedFuncLevel:     1,
		RuntimeCallerSkip:   2,
	}
)

const (
	ansiResetColor = "\u001b[0m"

	// black = "8;2;0;0;0"
	// white = "8;2;255;255;255"

	// zshGrey = "8;2;191;176;160"
	// zshCyan = "8;2;69;132;136"

	// redGoogle    = "8;2;219;68;55"
	// blueGoogle   = "8;2;76;139;245"
	// yellowGoogle = "8;2;244;180;0"
	// greenGoogle  = "8;2;15;157;88"

	black   = "0"
	red     = "8;5;1"
	green   = "8;5;2"
	yellow  = "8;5;3"
	blue    = "8;5;4"
	magenta = "8;5;5"
	cyan    = "8;5;6"
	grey    = "8;5;7"
	white   = "8;5;231"
)

type Tag struct {
	BackgroundColor string
	Color           string
}

type LogFormat struct {
	Level    string
	Message  string
	FuncName string
	Location string
	Time     time.Time
	Tag      Tag
}

func Warning(message string) {
	timezone, _ := time.LoadLocation(LogConf.TimeZone)
	log := LogFormat{
		Level:    "WARNING",
		Message:  message,
		FuncName: getFuncName(LogConf.NestedFuncLevel),
		Location: getPathAndLineNumber(LogConf.NestedLocationLevel),
		Time:     time.Now().In(timezone),
		Tag: Tag{
			BackgroundColor: yellow,
			Color:           black,
		},
	}
	print(log)
}

func Error(message string) {
	timezone, _ := time.LoadLocation(LogConf.TimeZone)
	log := LogFormat{
		Level:    "ERROR",
		Message:  message,
		FuncName: getFuncName(LogConf.NestedFuncLevel),
		Location: getPathAndLineNumber(LogConf.NestedLocationLevel),
		Time:     time.Now().In(timezone),
		Tag: Tag{
			BackgroundColor: red,
			Color:           white,
		},
	}
	print(log)
}

func Info(message string) {
	timezone, _ := time.LoadLocation(LogConf.TimeZone)
	log := LogFormat{
		Level:    "INFO",
		Message:  message,
		FuncName: getFuncName(LogConf.NestedFuncLevel),
		Location: getPathAndLineNumber(LogConf.NestedLocationLevel),
		Time:     time.Now().In(timezone),
		Tag: Tag{
			BackgroundColor: green,
			Color:           white,
		},
	}
	print(log)
}

func PrintJSONIndent(title string, data interface{}) {
	LogConf.RuntimeCallerSkip = LogConf.RuntimeCallerSkip + 1
	json, err := json.MarshalIndent(data, " ", "    ")
	if err != nil {
		Error(err.Error())
		LogConf.RuntimeCallerSkip = LogConf.RuntimeCallerSkip - 1
		return
	}
	Info(fmt.Sprintf("%s:\n%s", title, json))
	LogConf.RuntimeCallerSkip = LogConf.RuntimeCallerSkip - 1
}

func writeLogFile(content string) {
	if LogConf.CreateLogFile {
		timezone, _ := time.LoadLocation(LogConf.TimeZone)
		timeNow := time.Now().In(timezone).Format("2006-01-02")

		path := filepath.Join(LogConf.Path, timeNow+".log")
		WriteLogFile(path, content+"\n")
	}
}

func print(log LogFormat) {
	timeNow := log.Time.Format(LogConf.TimeFormat)

	// write log to log file
	logFuncName := "[" + log.FuncName + "]"
	if !LogConf.LogFuncName {
		logFuncName = ""
	}
	rawLog := "[" + timeNow + "][" + log.Level + "]" + logFuncName + "[" + log.Location + "] " + log.Message
	writeLogFile(rawLog)

	// print log to terminal with style
	logTime := createTextBackground(grey, black, timeNow)

	logType := "|" + createTextBackground(log.Tag.BackgroundColor, log.Tag.Color, fmt.Sprintf("%-7s", log.Level))

	logFuncName = "|" + createTextBackground(blue, white, fmt.Sprintf("%-30s", log.FuncName))
	if !LogConf.LogFuncName {
		logFuncName = ""
	}

	logLocation := "|" + createTextBackground(cyan, white, fmt.Sprintf("%-40s", log.Location)) + "|"

	if LogConf.DebugMode {
		strLog := fmt.Sprintf("%s\t%s\t%s\t%s\t %s", logTime, logType, logFuncName, logLocation, log.Message)
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 0, 0, '\t', 0)
		fmt.Fprintln(w, strLog)
		w.Flush()
		// fmt.Println(strLog)
	}
}

func createTextBackground(backgroundColor, textColor, text string) string {
	label := fmt.Sprintf(" %s ", text)
	return createBackgroundColor(backgroundColor) + createColor(textColor) + label + ansiResetColor
}

func createBackgroundColor(color string) string {
	return "\u001b[4" + color + "m"
}

func createColor(color string) string {
	return "\u001b[3" + color + "m"
}

func getFuncName(n int) string {
	pc := make([]uintptr, 1)
	runtime.Callers(LogConf.RuntimeCallerSkip+1, pc) // adjust the number
	f := runtime.FuncForPC(pc[0])
	funcName := f.Name()

	arr := strings.Split(funcName, "/")
	sliceArr := arr
	if n <= 0 {
		sliceArr = arr
	} else {
		s := len(arr) - n
		if s < 0 {
			sliceArr = arr
		} else {
			sliceArr = arr[s:]
		}
	}

	lastFunc := strings.Join(sliceArr, "/")
	return lastFunc
}

func getPathAndLineNumber(n int) string {
	_, filePath, lineNumber, ok := runtime.Caller(LogConf.RuntimeCallerSkip) // adjust the number
	if !ok {
		return "?:?"
	}

	arr := strings.Split(filePath, "/")
	sliceArr := arr
	if n <= 0 {
		sliceArr = arr
	} else {
		s := len(arr) - n
		if s < 0 {
			sliceArr = arr
		} else {
			sliceArr = arr[s:]
		}
	}

	selectedPath := strings.Join(sliceArr, "/")
	return selectedPath + ":" + cast.ToString(lineNumber)
}

func WriteLogFile(fullFilepath, content string) {
	// jika path belum ada
	if !IsPathExist(fullFilepath) {
		file := filepath.Base(fullFilepath)
		dirPath := strings.Replace(fullFilepath, file, "", 1)

		fmt.Printf("create a new folder to store log files: %s\n", dirPath)
		CreateNewNestedDirectory(dirPath)
	}

	// ubah permission file jadi writeable
	if IsPathExist(fullFilepath) {
		SetWritable(fullFilepath)
	}

	f, err := os.OpenFile(fullFilepath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println("error opening file:", err)
	}

	defer f.Close()

	if _, err = f.WriteString(content); err != nil {
		fmt.Println("error writing file:", err)
	}

	// ubah permission file jadi read-only
	SetReadOnly(fullFilepath)
}

func SetWritable(filepath string) error {
	err := os.Chmod(filepath, 0222)
	if err != nil {
		fmt.Println("error change permission file to writeable:", err)
	}
	return err
}

func SetReadOnly(filepath string) error {
	err := os.Chmod(filepath, 0444)
	if err != nil {
		fmt.Println("error change permission file to read-only:", err)
	}
	return err
}

func CreateNewNestedDirectory(folderPath string) error {
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		fmt.Println("error creating a new directory or file:", err)
	}
	return err
}

func IsPathExist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

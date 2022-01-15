package gologger

import (
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
	Path                string
	DebugMode           bool
	NestedLocationLevel int
	NestedFuncLevel     int
}

var (
	LogConf = LogConfiguration{
		TimeZone:            "Asia/Jakarta",
		Path:                "./logs",
		DebugMode:           true,
		NestedLocationLevel: 2,
		NestedFuncLevel:     2,
	}
)

const (
	ansiReset = "\u001b[0m"

	black   = "0"
	red     = "1"
	green   = "2"
	yellow  = "3"
	blue    = "4"
	magenta = "5"
	cyan    = "6"
	white   = "7"
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

func writeLogFile(content string) {
	timezone, _ := time.LoadLocation(LogConf.TimeZone)
	timeNow := time.Now().In(timezone).Format("2006-01-02")

	path := filepath.Join(LogConf.Path, timeNow+".log")
	WriteLogFile(path, content+"\n")
}

func print(log LogFormat) {
	logType := createTextBackground(log.Tag.BackgroundColor, log.Tag.Color, log.Level)

	timeNow := log.Time.Format("2006-01-02T15:04:05-0700")

	// write log to log file
	rawLog := "[" + timeNow + "][" + log.Level + "][" + log.FuncName + "][" + log.Location + "] " + log.Message
	writeLogFile(rawLog)

	// print log to terminal with style
	logTime := createTextBackground(white, black, timeNow)

	logFuncName := createTextBackground(black, white, log.FuncName)

	logLocation := createTextBackground(blue, white, log.Location)

	if LogConf.DebugMode {
		strLog := fmt.Sprintf("%s\t%s\t%s\t%s %s", logTime, logType, logFuncName, logLocation, log.Message)
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 0, 0, '\t', tabwriter.Debug)
		fmt.Fprintln(w, strLog)
		w.Flush()
		// fmt.Println(strLog)
	}
}

func createTextBackground(backgroundColor, textColor, text string) string {
	label := fmt.Sprintf(" %7s ", text)
	return createBackgroundColor(backgroundColor, false) + createColor(textColor, false) + label + ansiReset
}

func createBackgroundColor(color string, isBright bool) string {
	str := ""
	if isBright {
		str = ";1"
	}

	return "\u001b[4" + color + str + "m"
}

func createColor(color string, isBright bool) string {
	str := ""
	if isBright {
		str = ";1"
	}

	return "\u001b[3" + color + str + "m"
}

func getFuncName(n int) string {
	pc := make([]uintptr, 1)
	runtime.Callers(3, pc) // adjust the number
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
	_, filePath, lineNumber, ok := runtime.Caller(2) // adjust the number
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

		fmt.Println("creating a new nested file path")
		CreateNewNestedDirectory(dirPath)
	}

	// ubah permission file jadi writeable
	SetWritable(fullFilepath)

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

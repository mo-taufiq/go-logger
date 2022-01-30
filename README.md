# **go-logger**

### **Installation**
```zsh
go get github.com/mo-taufiq/go-logger
```

### **Quickstart**
```go
package main

import (
	gologger "github.com/mo-taufiq/go-logger"
)

func main() {
	gologger.LogConf.TimeZone = "Asia/Jakarta"
	gologger.LogConf.TimeFormat = "2006-01-02T15:04:05-0700"
	gologger.LogConf.CreateLogFile = true
	gologger.LogConf.Path = "./logs"
	gologger.LogConf.DebugMode = true
	gologger.LogConf.NestedLocationLevel = 1
	gologger.LogConf.LogFuncName = true
	gologger.LogConf.NestedFuncLevel = 1
	gologger.LogConf.RuntimeCallerSkip = 2

	gologger.Info("log info")
	gologger.Warning("log warning")
	gologger.Error("log error")
}
```

### **Preview**
![Preview go-logger](https://media.giphy.com/media/5kQpeJJkTwHHp1CZNl/giphy.gif)
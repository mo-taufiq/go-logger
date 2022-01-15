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
	gologger.LogConf = gologger.LogConfiguration{
		DebugMode:           true,
		Path:                "./logs",
		TimeZone:            "Asia/Jakarta",
		NestedFuncLevel:     1,
		NestedLocationLevel: 2,
	}

	gologger.Info("log error")
	gologger.Warning("log warning")
	gologger.Error("log error")
}
```

### **Preview**
![Preview go-logger](https://media.giphy.com/media/5kQpeJJkTwHHp1CZNl/giphy.gif)
# Flog

[![Build Status](https://travis-ci.org/ivpusic/grpool.svg?branch=master)](https://github.com/infinitasx/easi-go-aws)

## ⚙ Installation

```text
go get -u github.com/eininst/flog
```

## ⚡ Quickstart

```go
package main

import (
	"github.com/eininst/flog"
)

func main() {
    flog.Trace("Something very low level.")
    flog.Debug("Useful debugging information.")
    flog.Info("Something noteworthy happened!")
    flog.Warn("You should probably take a look at this.")
    flog.Error("Something failed but I'm not quitting.")

    flog.Info("1", "2", "3")
    flog.Info(flog.Sprintf("My name is ${name}", flog.H{
        "name": "jack",
    }))
    flog.Infof("My name is %s", "tomi")

    //Calls os.Exit(1) after logging
    flog.Fatal("Bye.")
    //Calls panic() after logging
    flog.Panic("I'm bailing.")
}
```

```
2022/08/28 17:44:13 [TRACE] test.go:8 Something very low level. 
2022/08/28 17:44:13 [DEBUG] test.go:9 Useful debugging information.     
2022/08/28 17:44:13 [INFO] test.go:10 Something noteworthy happened!    
2022/08/28 17:44:13 [WARN] test.go:11 You should probably take a look at this.  
2022/08/28 17:44:13 [ERROR] test.go:12 Something failed but I'm not quitting.   
2022/08/28 17:44:13 [INFO] test.go:14 1 2 3     
2022/08/28 17:44:13 [INFO] test.go:15 My name is jack   
2022/08/28 17:44:13 [INFO] test.go:18 My name is tomi   
2022/08/28 17:44:13 [FATAL] test.go:21 Bye. 
```

> With Fields

```go
package main

import (
	"github.com/eininst/flog"
)

func main() {
    flog.Debug("Start with fields！")

    flog.With(flog.Fields{
        "animal": "walrus",
        "size":   10,
    }).Info("A group of walrus emerges from the ocean")

    flog.With(flog.Fields{
        "omg":    true,
        "number": 122,
    }).Warn("The group's number increased tremendously!")

    flog.With(flog.Fields{
        "name":   "wzq",
        "omg":    true,
        "number": 100,
    }).Error("The ice breaks!")

    contextLogger := flog.With(flog.Fields{
        "common": "this is a common field",
        "other":  "I also should be logged always",
    })

    contextLogger.Info("I'll be logged with common and other field")
    contextLogger.Info("Me too")
}
```

```textmate
2022/08/29 20:55:41 [DEBUG] main.go:12 Start with fields！
2022/08/29 20:55:41 [INFO] main.go:17 A group of walrus emerges from the ocean  animal=walrus size=10
2022/08/29 20:55:41 [WARN] main.go:22 The group's number increased tremendously!  number=122 omg=true
2022/08/29 20:55:41 [ERROR] main.go:28 The ice breaks!  name=wzq number=100 omg=true
2022/08/29 20:55:41 [INFO] main.go:35 I'll be logged with common and other field  common=this is a common field other=I also should be logged always
2022/08/29 20:55:41 [INFO] main.go:36 Me too  common=this is a common field other=I also should be logged always
```

> You can customize it all you want:

```go
func init() {
    flog.SetLevel(flog.InfoLevel)
    
    flog.SetTimeFormat("2006.01.02 15:04:05.000")

    logf := "%s[${pid}]%s ${time} ${level} ${path} ${msg}"
    flog.SetFormat(fmt.Sprintf(logf, flog.Cyan, flog.Reset))

    flog.SetFullPath(true)
    
    flog.SetMsgMinLen(50)
}
```

```text
[79338] 2022.09.04 10:02:55.140 [INFO] /Users/wangziqing/go/flog/examples/main.go:27 A group of walrus emerges from the ocean            animal=walrus size=10
[79338] 2022.09.04 10:02:55.140 [WARN] /Users/wangziqing/go/flog/examples/main.go:32 The group's number increased tremendously!          number=122 omg=true
[79338] 2022.09.04 10:02:55.140 [ERROR] /Users/wangziqing/go/flog/examples/main.go:38 The ice breaks!                                     name=wzq number=100 omg=true
[79338] 2022.09.04 10:02:55.140 [INFO] /Users/wangziqing/go/flog/examples/main.go:45 I'll be logged with common and other field          common=this is a common field other=I also should be logged always
[79338] 2022.09.04 10:02:55.140 [INFO] /Users/wangziqing/go/flog/examples/main.go:46 Me too                                              common=this is a common field other=I also should be logged always
```

> Dump Json

```go
func init() {
    flog.DumpJson()
}
```

```textmate
{"level":"DEBUG","msg":"Start with fields！","path":"/Users/wangziqing/go/flog/examples/main.go:24","pid":"79546","time":"2022.09.04 10:04:59.105"}
{"animal":"walrus","level":"INFO","msg":"A group of walrus emerges from the ocean","path":"/Users/wangziqing/go/flog/examples/main.go:29","pid":"79546","size":10,"time":"2022.09.04 10:04:59.105"}
{"level":"WARN","msg":"The group's number increased tremendously!","number":122,"omg":true,"path":"/Users/wangziqing/go/flog/examples/main.go:34","pid":"79546","time":"2022.09.04 10:04:59.105"}
{"level":"ERROR","msg":"The ice breaks!","name":"wzq","number":100,"omg":true,"path":"/Users/wangziqing/go/flog/examples/main.go:40","pid":"79546","time":"2022.09.04 10:04:59.105"}
{"common":"this is a common field","level":"INFO","msg":"I'll be logged with common and other field","other":"I also should be logged always","path":"/Users/wangziqing/go/flog/examples/main.go:47","pid":"79546","time":"2022.09.04 10:04:59.105"}
{"common":"this is a common field","level":"INFO","msg":"Me too","other":"I also should be logged always","path":"/Users/wangziqing/go/flog/examples/main.go:48","pid":"79546","time":"2022.09.04 10:04:59.105"}
```

> See [examples](/examples)

## License

*MIT*
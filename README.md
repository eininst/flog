# Flog

[![Build Status](https://travis-ci.org/ivpusic/grpool.svg?branch=master)](https://github.com/infinitasx/easi-go-aws)

## Installation
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

#### With Fields
```go
package main

import (
	"github.com/eininst/flog"
)

func main() {
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
2022/08/28 19:09:41 [INFO] main.go:11 A group of walrus emerges from the ocean          animal=walrus size=10
2022/08/28 19:09:41 [WARN] main.go:16 The group's number increased tremendously!        omg=true number=122
2022/08/28 19:09:41 [ERROR] main.go:22 The ice breaks!                                  name=wzq omg=true number=100
2022/08/28 19:09:41 [INFO] main.go:29 I'll be logged with common and other field        common=this is a common field other=I also should be logged always
2022/08/28 19:09:41 [INFO] main.go:30 Me too                                            other=I also should be logged always common=this is a common field
```

> You can customize it all you want:
```go
func init() {
    flog.SetLevel(flog.InfoLevel)
    flog.SetFormat("${level} ${time} ${path} ${msg}\t${fields}")
    flog.SetTimeFormat("2006.01.02 15:04:05.000")
    flog.SetFullPath(true)
}
```

```text
[INFO] 2022.08.28 19:15:42.039 /Users/wangziqing/go/flog/test/main.go:17 Something noteworthy happened! 
[WARN] 2022.08.28 19:15:42.039 /Users/wangziqing/go/flog/test/main.go:18 You should probably take a look at this.       
[ERROR] 2022.08.28 19:15:42.039 /Users/wangziqing/go/flog/test/main.go:19 Something failed but I'm not quitting.        
[INFO] 2022.08.28 19:15:42.039 /Users/wangziqing/go/flog/test/main.go:21 1 2 3  
[INFO] 2022.08.28 19:15:42.039 /Users/wangziqing/go/flog/test/main.go:22 My name is jack        
[INFO] 2022.08.28 19:15:42.039 /Users/wangziqing/go/flog/test/main.go:25 My name is tomi        
[INFO] 2022.08.28 19:15:42.039 /Users/wangziqing/go/flog/test/main.go:30 A group of walrus emerges from the ocean       size=10 animal=walrus
[FATAL] 2022.08.28 19:15:42.039 /Users/wangziqing/go/flog/test/main.go:33 Bye.  
```


> Dump Json
```go
func init() {
    flog.DumpJson()
}
```

```textmate
{"level":"INFO","msg":"Something noteworthy happened!","path":"/Users/wangziqing/go/flog/test/main.go:18","time":"2022.08.28 19:17:45.394"}
{"level":"WARN","msg":"You should probably take a look at this.","path":"/Users/wangziqing/go/flog/test/main.go:19","time":"2022.08.28 19:17:45.395"}
{"level":"ERROR","msg":"Something failed but I'm not quitting.","path":"/Users/wangziqing/go/flog/test/main.go:20","time":"2022.08.28 19:17:45.395"}
{"level":"INFO","msg":"1 2 3","path":"/Users/wangziqing/go/flog/test/main.go:22","time":"2022.08.28 19:17:45.395"}
{"level":"INFO","msg":"My name is jack","path":"/Users/wangziqing/go/flog/test/main.go:23","time":"2022.08.28 19:17:45.395"}
{"level":"INFO","msg":"My name is tomi","path":"/Users/wangziqing/go/flog/test/main.go:26","time":"2022.08.28 19:17:45.395"}
{"animal":"walrus","level":"INFO","msg":"A group of walrus emerges from the ocean  ","path":"/Users/wangziqing/go/flog/test/main.go:31","size":10,"time":"2022.08.28 19:17:45.395"}
{"level":"FATAL","msg":"Bye.","path":"/Users/wangziqing/go/flog/test/main.go:34","time":"2022.08.28 19:17:45.395"}
```
## License
*MIT*
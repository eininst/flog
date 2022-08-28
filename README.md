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
	flog "github.com/eininst/fastgo-log"
)

func main() {
	flog.Trace("Something very low level.")
	flog.Debug("Useful debugging information.")
	flog.Info("Something noteworthy happened!")
	flog.Warn("You should probably take a look at this.")
	flog.Error("Something failed but I'm not quitting.")

	flog.Info(flog.Sprintf("My name is ${name}", flog.H{
		"name": "jack",
	}))

	//Calls os.Exit(1) after logging
	flog.Fatal("Bye.")
	//Calls panic() after logging
	flog.Panic("I'm bailing.")
}
```


```
2022/08/28 02:13:18 [TRACE] test.go:15 Something very low level.
2022/08/28 02:13:18 [DEBUG] test.go:16 Useful debugging information.
2022/08/28 02:13:18 [INFO] test.go:17 Something noteworthy happened!
2022/08/28 02:13:18 [WARN] test.go:18 You should probably take a look at this.
2022/08/28 02:13:18 [ERROR] test.go:19 Something failed but I'm not quitting.
2022/08/28 02:13:18 [INFO] test.go:21 My name is jack
2022/08/28 02:13:18 [FATAL] test.go:26 Bye.
```

> You can customize it all you want:
```go
func init() {
	flog.SetLevel(flog.InfoLevel)
	flog.SetFormat("${Level} ${Time} ${Path} ${Msg}")
	flog.SetTimeFormat("2006.01.02 15:04:05.000")
	flog.SetFullPath(true)
}
```

```text
[INFO] 2022.08.28 02:24:01.219 /Users/wangziqing/go/fastgo/cmd/consumer/test.go:16 Something noteworthy happened!
[WARN] 2022.08.28 02:24:01.219 /Users/wangziqing/go/fastgo/cmd/consumer/test.go:17 You should probably take a look at this.
[ERROR] 2022.08.28 02:24:01.219 /Users/wangziqing/go/fastgo/cmd/consumer/test.go:18 Something failed but I'm not quitting.
[INFO] 2022.08.28 02:24:01.219 /Users/wangziqing/go/fastgo/cmd/consumer/test.go:20 My name is jack
[FATAL] 2022.08.28 02:24:01.219 /Users/wangziqing/go/fastgo/cmd/consumer/test.go:25 Bye.
```
## License
*MIT*
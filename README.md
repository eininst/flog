# Log

[![Build Status](https://travis-ci.org/ivpusic/grpool.svg?branch=master)](https://github.com/infinitasx/easi-go-aws)

## Use
```go
flog.Trace("Something very low level.")
flog.Debug("Useful debugging information.")
flog.Info("Something noteworthy happened!")
flog.Warn("You should probably take a look at this.")
flog.Error("Something failed but I'm not quitting.")

flog.Info(color.Green(flog.Sprintf("My name is {{name}}", flog.H{
    "name": "jack",
})))
```

## License
*MIT*
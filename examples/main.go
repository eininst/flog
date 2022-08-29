package main

import (
	"github.com/eininst/flog"
)

func init() {
	flog.SetLevel(flog.InfoLevel)

	flog.SetTimeFormat("2006.01.02 15:04:05.000")

	flog.SetFormat("${level} ${time} ${path} ${msg}")

	flog.SetFullPath(true)

	flog.SetMsgMinLen(42)
}

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

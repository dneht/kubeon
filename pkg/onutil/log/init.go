package log

import (
	"flag"
	"fmt"
	"k8s.io/klog/v2"
	"strconv"
)

var logLevel = 1

func Init(level int) {
	if level < 0 {
		level = 0
	} else if level > 8 {
		level = 8
	}
	logLevel = level
	flagSet := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(flagSet)
	err := flagSet.Parse([]string{"-v", strconv.FormatInt(int64(level), 10)})
	if nil != err {
		fmt.Printf("Set log level failed: %v\n", err)
	}
	flagSet.Parsed()
}

func Level() int {
	return logLevel
}

func IsDebug() bool {
	return logLevel >= 6
}

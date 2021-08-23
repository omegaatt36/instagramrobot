package helpers

import (
	"os"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"
)

func RegisterLogger() {
	log.SetReportCaller(true)
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		ForceQuote:      true,
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.Kitchen,
		PadLevelText:    true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			// Show file name + line
			// c := f.Func.Name()
			// c = c[strings.LastIndex(c, "."):]
			// m := f.File
			// m = m[strings.LastIndex(m, "/")+1 : len(m)-3]
			// return "[" + m + c + ":" + strconv.Itoa(f.Line) + "]", ""
			return "", ""
		},
	})
}

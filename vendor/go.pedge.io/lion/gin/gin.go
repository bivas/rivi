/*
Package ginlion defines functionality to integrate lion with gin.

https://github.com/gin-gonic/gin

Full typical use:

	import (
		"go.pedge.io/lion/env"
		"go.pedge.io/lion/gin"
	)

	func setupGin() error {
		if err := envlion.Setup(); err != nil {
			return err
		}
		engine := ginlion.Default()
		...
	}

Some of the code here is copied from the gin repository.
This code is under the MIT License that can be found at https://github.com/gin-gonic/gin/blob/master/LICENSE.
*/
package ginlion // import "go.pedge.io/lion/gin"

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.pedge.io/lion/proto"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// Default is the equivalent of gin's Default function.
func Default() *gin.Engine {
	engine := gin.New()
	engine.Use(GlobalLoggerAndRecovery())
	return engine
}

// GlobalLoggerAndRecovery returns LoggerAndRecovery on the global proto Logger.
func GlobalLoggerAndRecovery() gin.HandlerFunc {
	return LoggerAndRecovery(protolion.GlobalLogger)
}

// LoggerAndRecovery is the equivalent of both gin's Logger and Recovery middlewares.
func LoggerAndRecovery(protoLoggerProvider func() protolion.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		var path string
		if c.Request.URL != nil {
			path = c.Request.URL.Path
		}
		defer func() {
			call := &Call{
				Method:      c.Request.Method,
				Path:        path,
				UserAgent:   c.Request.Header.Get("User-Agent"),
				RequestForm: valuesMap(c.Request.Form),
				ClientIp:    c.ClientIP(),
				StatusCode:  uint32(statusCode(c.Writer.Status())),
				Error:       c.Errors.Errors(),
			}
			if c.Request.URL != nil {
				call.Query = valuesMap(c.Request.URL.Query())
			}
			call.Duration = fmt.Sprintf("%v", time.Since(start))
			protoLogger := protoLoggerProvider()
			if recoverErr := recover(); recoverErr != nil {
				stack := stack(3)
				call.Error = append(call.Error, fmt.Sprintf("panic: %s\n%s", recoverErr, string(stack)))
				protoLogger.Error(call)
				c.AbortWithStatus(http.StatusInternalServerError)
			} else {
				protoLogger.Info(call)
			}
		}()
		c.Next()
	}
}

// stack returns a nicely formated stack frame, skipping skip frames
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

func valuesMap(values map[string][]string) map[string]string {
	if values == nil {
		return nil
	}
	m := make(map[string]string)
	for key, value := range values {
		m[key] = strings.Join(value, " ")
	}
	return m
}

func statusCode(code int) int {
	if code == 0 {
		return http.StatusOK
	}
	return code
}

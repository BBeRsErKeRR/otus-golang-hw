package logger

import (
	"fmt"
	"testing"

	"bytes"
	"io"
	"os"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	testcases := []struct {
		Name     string
		LogLevel string
		Result   string
		Action   func(logger Logger)
	}{
		{
			Name:     "test info",
			LogLevel: "info",
			Result:   "Info message",
			Action: func(log Logger) {
				log.Info("Info message")
			},
		},
		{
			Name:     "test debug",
			LogLevel: "debug",
			Result:   "Debug message",
			Action: func(log Logger) {
				log.Debug("Debug message")
			},
		},
		{
			Name:     "test error",
			LogLevel: "error",
			Result:   "Error message",
			Action: func(log Logger) {
				log.Error("Error message")
			},
		},
		{
			Name:     "test warn",
			LogLevel: "warn",
			Result:   "Warn message",
			Action: func(log Logger) {
				log.Warn("Warn message")
			},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			level := testcase.LogLevel

			conf := &LoggerConf{
				Level:    level,
				OutPaths: []string{"stdout"},
				ErrPaths: []string{"stderr"},
			}
			bcp := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			logger, err := New(conf)
			require.Nil(t, err)
			require.NotNil(t, logger)
			testcase.Action(logger)

			outC := make(chan string)

			go func() {
				var buf bytes.Buffer
				io.Copy(&buf, r)
				outC <- buf.String()
			}()
			w.Close()
			os.Stdout = bcp
			output := <-outC

			require.Contains(t, string(output), fmt.Sprintf(`"level":"%s"`, level))
			require.Contains(t, string(output), fmt.Sprintf(`"msg":"%s"`, testcase.Result))
		})
	}

	regresscases := []struct {
		Name     string
		Expected string
		Config   *LoggerConf
	}{
		{
			Name: "bad level value case",
			Config: &LoggerConf{
				Level:    "bad_level",
				OutPaths: []string{"stdout"},
				ErrPaths: []string{"stderr"},
			},
			Expected: `unrecognized level: "bad_level"`,
		},
		{
			Name: "bad output paths",
			Config: &LoggerConf{
				Level:    "info",
				OutPaths: []string{"http://100"},
				ErrPaths: []string{"stderr"},
			},
			Expected: `open sink "http://100": no sink found for scheme "http"`,
		},
	}
	for _, testcase := range regresscases {
		t.Run(testcase.Name, func(t *testing.T) {
			log, err := New(testcase.Config)
			require.Nil(t, log)
			require.Equal(t, testcase.Expected, err.Error())
		})
	}
}

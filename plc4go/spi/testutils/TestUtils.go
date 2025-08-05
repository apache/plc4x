/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package testutils

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ajankovic/xdiff"
	"github.com/ajankovic/xdiff/parser"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/stretchr/testify/assert"

	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/pool"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	"github.com/apache/plc4x/plc4go/spi/transports/test"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

func CompareResults(t *testing.T, actualString []byte, referenceString []byte) error {
	localLog := ProduceTestingLogger(t)
	// Now parse the xml strings of the actual and the reference in xdiff's dom
	p := parser.New()
	actual, err := p.ParseBytes(actualString)
	if err != nil {
		return errors.Wrap(err, "Error parsing actual input")
	}
	reference, err := p.ParseBytes(referenceString)
	if err != nil {
		return errors.Wrap(err, "Error parsing reference input")
	}
	// Use XDiff to actually do the comparison
	diff, err := xdiff.Compare(actual, reference)
	if err != nil {
		return errors.Wrap(err, "Error comparing xml trees")
	}
	if diff == nil {
		// All good
		return nil
	}
	cleanDiff := make([]xdiff.Delta, 0)
	for _, delta := range diff {
		if delta.Operation == xdiff.Delete && delta.Subject.Value == nil || delta.Operation == xdiff.Insert && delta.Subject.Value == nil || delta.Operation == xdiff.InsertSubtree && delta.Subject.Value == nil {
			localLog.Info().Interface("delta", delta).Msg("We ignore empty elements which should be deleted")
			continue
		}
		// Workaround for different precisions on float
		if delta.Operation == xdiff.Update &&
			string(delta.Subject.Parent.FirstChild.Name) == "dataType" &&
			string(delta.Subject.Parent.FirstChild.Value) == "float" &&
			string(delta.Object.Parent.FirstChild.Name) == "dataType" &&
			string(delta.Object.Parent.FirstChild.Value) == "float" {
			if strings.Contains(string(delta.Subject.Value), string(delta.Object.Value)) || strings.Contains(string(delta.Object.Value), string(delta.Subject.Value)) {
				localLog.Info().Interface("delta", delta).Msg("We ignore precision diffs")
				continue
			}
		}
		if delta.Operation == xdiff.Update &&
			string(delta.Subject.Parent.FirstChild.Name) == "dataType" &&
			string(delta.Subject.Parent.FirstChild.Value) == "string" &&
			string(delta.Object.Parent.FirstChild.Name) == "dataType" &&
			string(delta.Object.Parent.FirstChild.Value) == "string" {
			if diff, err := xdiff.Compare(delta.Subject, delta.Object); diff == nil && err == nil {
				localLog.Info().Interface("delta", delta).Msg("We ignore newline diffs")
				continue
			}
		}
		cleanDiff = append(cleanDiff, delta)
	}

	enc := xdiff.NewTextEncoder(os.Stdout)
	if err := enc.Encode(diff); err != nil {
		return errors.Wrap(err, "Error outputting results")
	}
	if len(cleanDiff) <= 0 {
		localLog.Warn().Msg("We only found non relevant changes")
		return nil
	}

	assert.Equal(t, string(referenceString), string(actualString))
	asciiBoxWriter := utils.NewAsciiBoxWriter()
	expectedBox := asciiBoxWriter.BoxString(string(referenceString), utils.WithAsciiBoxName("expected"))
	gotBox := asciiBoxWriter.BoxString(string(actualString), utils.WithAsciiBoxName("got"))
	boxSideBySide := asciiBoxWriter.BoxSideBySide(expectedBox, gotBox)
	_ = boxSideBySide // TODO: xml too distorted, we need a don't center option
	return errors.New("there were differences: Expected: \n" + string(referenceString) + "\nBut Got: \n" + string(actualString))
}

// TestContext produces a context which is getting cleaned up by testing.T
func TestContext(t *testing.T) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	ctx = ProduceTestingLogger(t).WithContext(ctx)
	return ctx
}

var (
	muteLog                             bool
	highLogPrecision                    bool
	passLoggerToModel                   bool
	receiveTimeout                      time.Duration
	traceTransactionManagerWorkers      bool
	traceTransactionManagerTransactions bool
	traceDefaultMessageCodecWorker      bool
	traceExecutorWorkers                bool
	traceTestTransportInstance          bool
)

func init() {
	getOrLeaveBool("PLC4X_TEST_MUTE_LOG", &muteLog)
	getOrLeaveBool("PLC4X_TEST_HIGH_TEST_LOG_PRECISION", &highLogPrecision)
	if highLogPrecision {
		zerolog.TimeFieldFormat = time.RFC3339Nano
	}
	getOrLeaveBool("PLC4X_TEST_PASS_LOGGER_TO_MODEL", &passLoggerToModel)
	receiveTimeout = 3 * time.Second
	getOrLeaveDuration("PLC4X_TEST_RECEIVE_TIMEOUT_MS", &receiveTimeout)
	getOrLeaveBool("PLC4X_TEST_TRACE_TRANSACTION_MANAGER_WORKERS", &traceTransactionManagerWorkers)
	getOrLeaveBool("PLC4X_TEST_TRACE_TRANSACTION_MANAGER_TRANSACTIONS", &traceTransactionManagerTransactions)
	getOrLeaveBool("PLC4X_TEST_TRACE_DEFAULT_MESSAGE_CODEC_WORKER", &traceDefaultMessageCodecWorker)
	getOrLeaveBool("PLC4X_TEST_TRACE_EXECUTOR_WORKERS", &traceExecutorWorkers)
	getOrLeaveBool("PLC4X_TEST_TEST_TRANSPORT_INSTANCE", &traceTestTransportInstance)
}

func getOrLeaveBool(key string, setting *bool) {
	if env, ok := os.LookupEnv(key); ok {
		*setting = strings.EqualFold(env, "true")
	}
}

func getOrLeaveDuration(key string, setting *time.Duration) {
	if env, ok := os.LookupEnv(key); ok && env != "" {
		parsedDuration, err := strconv.ParseInt(env, 10, 64)
		if err != nil {
			panic(err)
		}
		*setting = time.Duration(parsedDuration) * time.Millisecond
	}
}

func shouldNoColor() bool {
	if _, forceColorEnv := os.LookupEnv("FORCE_COLOR"); forceColorEnv {
		color.NoColor = false // Apparently the color.NoColor is a bit to eager
		return false
	}
	noColor := false
	{
		_, noColorEnv := os.LookupEnv("NO_COLOR")
		onJenkins := os.Getenv("JENKINS_URL") != ""
		onGithubAction := os.Getenv("GITHUB_ACTIONS") != ""
		onCI := os.Getenv("CI") != ""
		if noColorEnv || onJenkins || onGithubAction || onCI {
			noColor = true
		}
	}
	if !noColor {
		color.NoColor = false // Apparently the color.NoColor is a bit to eager
	}
	return noColor
}

type TestingLog interface {
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Helper()
}

// ProduceTestingLogger produces a logger which redirects to testing.T
func ProduceTestingLogger(t TestingLog) zerolog.Logger {
	if muteLog {
		return zerolog.Nop()
	}
	noColor := shouldNoColor()
	consoleWriter := zerolog.NewConsoleWriter(
		zerolog.ConsoleTestWriter(t),
		func(w *zerolog.ConsoleWriter) {
			w.NoColor = noColor
		},
		func(w *zerolog.ConsoleWriter) {
			if highLogPrecision {
				w.TimeFormat = time.StampNano
			}
		},
		func(w *zerolog.ConsoleWriter) {
			w.FormatFieldValue = func(i interface{}) string {
				switch i := i.(type) {
				case string:
					if strings.Contains(i, "\\n") {
						if noColor {
							return "see below"
						} else {
							return fmt.Sprintf("\x1b[%dm%v\x1b[0m", 31, "see below")
						}
					}
				case []uint8:
					if len(i) > 4 && i[0] == '[' && i[len(i)-1] == ']' && strings.Contains(string(i), "\\n") {
						if noColor {
							return "see below"
						} else {
							return fmt.Sprintf("\x1b[%dm%v\x1b[0m", 31, "see below")
						}
					}
				}
				return fmt.Sprintf("%s", i)
			}
			w.FormatExtra = func(m map[string]interface{}, buffer *bytes.Buffer) error {
				for key, i := range m {
					switch i := i.(type) {
					case string:
						if strings.Contains(i, "\n") {
							buffer.WriteString("\n")
							if noColor {
								buffer.WriteString("field " + key)
							} else {
								buffer.WriteString(fmt.Sprintf("\x1b[%dm%v\x1b[0m", 32, "field "+key))
							}
							buffer.WriteString(":\n" + i)
						}
					case []any:
						allStrings := false
						containsNewLine := false
						stringsElems := make([]string, len(i))
						for i, elem := range i {
							if aString, ok := elem.(string); ok {
								containsNewLine = containsNewLine || strings.Contains(aString, "\n")
								allStrings = true
								stringsElems[i] = strings.Trim(aString, "\n")
							} else {
								allStrings = false
								break
							}
						}
						if allStrings && containsNewLine {
							buffer.WriteString("\n")
							if noColor {
								buffer.WriteString("field " + key)
							} else {
								buffer.WriteString(fmt.Sprintf("\x1b[%dm%v\x1b[0m", 32, "field "+key))
							}
							var sb strings.Builder
							for j, elem := range stringsElems {
								sb.WriteString(strconv.Itoa(j) + ":\n")
								sb.WriteString(elem)
								sb.WriteString("\n")
							}
							buffer.WriteString(":\n" + sb.String())
						}
					}
				}
				return nil
			}
		},
	)
	logger := zerolog.New(
		consoleWriter,
	)
	if highLogPrecision {
		logger = logger.With().Timestamp().Logger()
	}
	stackSetter.Do(func() {
		zerolog.ErrorStackMarshaler = func(err error) interface{} {
			if err == nil {
				return nil
			}
			var r strings.Builder
			stack := pkgerrors.MarshalStack(err)
			if stack == nil {
				return nil
			}
			stackMap := stack.([]map[string]string)
			for _, entry := range stackMap {
				stackSourceFileName := entry[pkgerrors.StackSourceFileName]
				stackSourceLineName := entry[pkgerrors.StackSourceLineName]
				stackSourceFunctionName := entry[pkgerrors.StackSourceFunctionName]
				r.WriteString(fmt.Sprintf("\tat %v (%v:%v)\n", stackSourceFunctionName, stackSourceFileName, stackSourceLineName))
			}
			return r.String()
		}
	})
	logger = logger.With().Stack().Logger()
	return logger
}

var stackSetter sync.Once

// EnrichOptionsWithOptionsForTesting appends options useful for testing to config.WithOption s
func EnrichOptionsWithOptionsForTesting(t *testing.T, _options ...options.WithOption) []options.WithOption {
	if extractedTraceWorkers, found := options.ExtractTracerWorkers(_options...); found {
		traceExecutorWorkers = extractedTraceWorkers
	}
	_options = append(_options,
		options.WithCustomLogger(ProduceTestingLogger(t)),
		options.WithPassLoggerToModel(passLoggerToModel),
		options.WithReceiveTimeout(receiveTimeout),
		options.WithTraceTransactionManagerWorkers(traceTransactionManagerWorkers),
		options.WithTraceTransactionManagerTransactions(traceTransactionManagerTransactions),
		options.WithTraceDefaultMessageCodecWorker(traceDefaultMessageCodecWorker),
		options.WithExecutorOptionTracerWorkers(traceExecutorWorkers),
		test.WithTraceTransportInstance(traceTestTransportInstance),
	)
	// We always create a custom executor to ensure shared executor for transaction manager is not used for tests
	testSharedExecutorInstance := pool.NewFixedSizeExecutor(
		runtime.NumCPU(),
		100,
		_options...,
	)
	testSharedExecutorInstance.Start()
	t.Cleanup(testSharedExecutorInstance.Stop)
	_options = append(_options,
		transactions.WithCustomExecutor(testSharedExecutorInstance),
	)
	return _options
}

type _explodingGlobalLogger struct {
	hardExplode bool
}

func (e _explodingGlobalLogger) Write(_ []byte) (_ int, err error) {
	if e.hardExplode {
		debug.PrintStack()
		panic("found a global log usage")
	}
	return 0, errors.New("found a global log usage")
}

// ExplodingGlobalLogger Useful to find unredirected logs
func ExplodingGlobalLogger(hardExplode bool) {
	log.Logger = zerolog.New(_explodingGlobalLogger{hardExplode: hardExplode})
}

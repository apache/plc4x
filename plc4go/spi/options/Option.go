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

package options

import (
	"context"
	"github.com/rs/zerolog"
)

// WithOption is a marker interface for options supplied by the builders like WithDefaultTtl
type WithOption interface {
	isOption() bool
}

// Option is a marker struct which can be used to mark an option
type Option struct {
}

func (_ Option) isOption() bool {
	return true
}

// WithCustomLogger is a global option to supply a custom logger
func WithCustomLogger(logger zerolog.Logger) WithOption {
	return withCustomLogger{logger: logger}
}

// ExtractCustomLogger can be used to extract the custom logger
func ExtractCustomLogger(options ...WithOption) (customLogger zerolog.Logger) {
	for _, option := range options {
		switch option := option.(type) {
		case withCustomLogger:
			customLogger = option.logger
			return
		}
	}
	return
}

// WithPassLoggerToModel enables passing of log to the model
func WithPassLoggerToModel(passLogger bool) WithOption {
	return withPassLoggerToModel{passLogger: passLogger}
}

// ExtractPassLoggerToModel to extract the flag indicating that model should be passed to Model
func ExtractPassLoggerToModel(options ...WithOption) bool {
	for _, option := range options {
		switch option := option.(type) {
		case withPassLoggerToModel:
			return option.passLogger
		}
	}
	return false
}

// GetLoggerContextForModel returns a log context if the WithPassLoggerToModel WithOption is set
func GetLoggerContextForModel(ctx context.Context, log zerolog.Logger, options ...WithOption) context.Context {
	passToModel := false
optionsSearch:
	for _, option := range options {
		switch option := option.(type) {
		case withPassLoggerToModel:
			passToModel = option.passLogger
			break optionsSearch
		}
	}
	if passToModel {
		return log.WithContext(ctx)
	}
	return ctx
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type withCustomLogger struct {
	Option
	logger zerolog.Logger
}

type withPassLoggerToModel struct {
	Option
	passLogger bool
}

//
//
///////////////////////////////////////
///////////////////////////////////////

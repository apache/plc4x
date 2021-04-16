//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

package logging

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var oldLogger zerolog.Logger

// init is used for _ imports for easy log config
func init() {
	oldLogger = log.Logger
	log.Logger = log.Logger.Level(zerolog.ErrorLevel)
}

// ErrorLevel configures zerolog to WarnLevel
func ErrorLevel() {
	log.Logger = log.Logger.Level(zerolog.ErrorLevel)
}

// WarnLevel configures zerolog to WarnLevel
func WarnLevel() {
	log.Logger = log.Logger.Level(zerolog.WarnLevel)
}

// InfoLevel configures zerolog to InfoLevel
func InfoLevel() {
	log.Logger = log.Logger.Level(zerolog.InfoLevel)
}

// DebugLevel configures zerolog to DebugLevel
func DebugLevel() {
	log.Logger = log.Logger.Level(zerolog.DebugLevel)
}

// TraceLevel configures zerolog to TraceLevel
func TraceLevel() {
	log.Logger = log.Logger.Level(zerolog.TraceLevel)
}

// ResetLogging can be used to reset to the old log settings
func ResetLogging() {
	log.Logger = oldLogger
}

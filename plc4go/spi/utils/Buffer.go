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

package utils

// WithReaderWriterArgs is a marker interface for reader args supplied by the builders
type WithReaderWriterArgs interface {
	WithReaderArgs
	WithWriterArgs
}

// WithAdditionalStringRepresentation can be used by e.g. enums to supply an additional string representation
func WithAdditionalStringRepresentation(stringRepresentation string) WithReaderWriterArgs {
	return withAdditionalStringRepresentation{readerWriterArg: readerWriterArg{WithReaderArgs: readerArg{}, WithWriterArgs: writerArg{}}, stringRepresentation: stringRepresentation}
}

// WithRenderAsList indicates that an element can be rendered as list
func WithRenderAsList(renderAsList bool) WithReaderWriterArgs {
	return withRenderAsList{readerWriterArg: readerWriterArg{WithReaderArgs: readerArg{}, WithWriterArgs: writerArg{}}, renderAsList: renderAsList}
}

// WithEncoding specifies an encoding
func WithEncoding(encoding string) WithReaderWriterArgs {
	return withEncoding{readerWriterArg: readerWriterArg{WithReaderArgs: readerArg{}, WithWriterArgs: writerArg{}}, encoding: encoding}
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type readerWriterArg struct {
	WithReaderArgs
	WithWriterArgs
}

func (r readerWriterArg) isReaderArgs() bool {
	return r.WithReaderArgs != nil
}

func (r readerWriterArg) isWriterArgs() bool {
	return r.WithWriterArgs != nil
}

type withAdditionalStringRepresentation struct {
	readerWriterArg
	stringRepresentation string
}

type withRenderAsList struct {
	readerWriterArg
	renderAsList bool
}

type withEncoding struct {
	readerWriterArg
	encoding string
}

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////

func UpcastReaderArgs(args ...WithReaderArgs) []WithReaderWriterArgs {
	result := make([]WithReaderWriterArgs, len(args))
	for i, arg := range args {
		result[i] = readerWriterArg{arg, writerArg{}}
	}
	return result
}

func UpcastWriterArgs(args ...WithWriterArgs) []WithReaderWriterArgs {
	result := make([]WithReaderWriterArgs, len(args))
	for i, arg := range args {
		result[i] = readerWriterArg{readerArg{}, arg}
	}
	return result
}

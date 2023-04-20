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

import "container/list"

const (
	rwDataTypeKey             = "dataType"
	rwBitLengthKey            = "bitLength"
	rwStringRepresentationKey = "stringRepresentation"
	rwBitKey                  = "bit"
	rwByteKey                 = "byte"
	rwUintKey                 = "uint"
	rwIntKey                  = "int"
	rwFloatKey                = "float"
	rwStringKey               = "string"
	rwEncodingKey             = "encoding"
	rwIsListKey               = "isList"
)

type BufferCommons struct {
}

func (b BufferCommons) SanitizeLogicalName(logicalName string) string {
	if logicalName == "" {
		return "value"
	}
	return logicalName
}

func (b BufferCommons) IsToBeRenderedAsList(readerWriterArgs ...WithReaderWriterArgs) bool {
	for _, arg := range readerWriterArgs {
		if !arg.isWriterArgs() && !arg.isReaderArgs() {
			panic("not a reader or writer arg")
		}
		switch arg.(type) {
		case withRenderAsList:
			return arg.(withRenderAsList).renderAsList
		}
	}
	return false
}

func (b BufferCommons) ExtractAdditionalStringRepresentation(readerWriterArgs ...WithReaderWriterArgs) string {
	for _, arg := range readerWriterArgs {
		if !arg.isWriterArgs() && !arg.isReaderArgs() {
			panic("not a reader or writer arg")
		}
		switch arg.(type) {
		case withAdditionalStringRepresentation:
			return arg.(withAdditionalStringRepresentation).stringRepresentation
		}
	}
	return ""
}

type Stack struct {
	list.List
}

func (s *Stack) Push(value any) any {
	s.PushBack(value)
	return value
}

func (s *Stack) Pop() any {
	if s.Len() <= 0 {
		return nil
	}
	element := s.Back()
	if element == nil {
		return nil
	}
	s.Remove(element)
	return element.Value
}

func (s *Stack) Peek() any {
	back := s.Back()
	if back == nil {
		return nil
	}
	return back.Value
}

func (s *Stack) Empty() bool {
	return s.Len() == 0
}

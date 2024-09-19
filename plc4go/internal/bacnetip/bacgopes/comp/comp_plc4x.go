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

package comp

import (
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
)

// AddRootMessageIfAbundant can be used to add a spi.Message if not present
func AddRootMessageIfAbundant(options []Option, rootMessage spi.Message) []Option {
	return AddIfAbundant(options, rootMessage)
}

// WithRootMessage can be used to bubble up a spi.Message
func WithRootMessage(rootMessage spi.Message) Option {
	return AsOption(rootMessage)
}

// ExtractRootMessage gets a spi.Message if present
func ExtractRootMessage(options []Option) (rootMessage spi.Message) {
	rootMessage, _ = ExtractIfPresent[spi.Message](options)
	return
}

// AddNLMIfAbundant can be used to add a readWriteModel.NLM if not present
func AddNLMIfAbundant(options []Option, nlm readWriteModel.NLM) []Option {
	return AddIfAbundant(options, nlm)
}

// WithNLM can be used to bubble up a readWriteModel.NLM
func WithNLM(nlm readWriteModel.NLM) Option {
	return AsOption(nlm)
}

// ExtractNLM gets a readWriteModel.NLM if present
func ExtractNLM(options []Option) (nlm readWriteModel.NLM) {
	nlm, _ = ExtractIfPresent[readWriteModel.NLM](options)
	return nil
}

// WithAPDU can be used to bubble up a readWriteModel.APDU
func WithAPDU(apdu readWriteModel.APDU) Option {
	return AsOption(apdu)
}

// ExtractAPDU gets a readWriteModel.APDU if present
func ExtractAPDU(options []Option) (apdu readWriteModel.APDU) {
	apdu, _ = ExtractIfPresent[readWriteModel.APDU](options)
	return nil
}

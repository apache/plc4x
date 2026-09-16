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

package modbus

import (
	"fmt"

	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi"
)

// readRequestBuilder assembles a read request in which one rejected address costs only its own
// tag. The shared spiModel.DefaultPlcReadRequestBuilder parses the addresses in Build and returns
// nothing at all as soon as one of them doesn't parse, so a single bad address takes down the
// whole scan - and modbus addresses regularly come from somewhere other than the caller (a device
// profile handed down from a server, say), where one typo among hundreds of good points is
// ordinary rather than exceptional.
//
// plc4j does the same thing one layer up, in the shared builder: DefaultPlcReadRequest.Builder
// keeps the tag in the request carrying INVALID_ADDRESS and a null tag, and the driver reads that
// back with getTagResponseCode. The equivalent move here would change the builder every driver
// uses, so the rejection is carried by a tag type of this driver's own instead: the tag takes part
// in the request, the reader recognizes it (see modbusTagsOf) and reports it as INVALID_ADDRESS,
// and nothing of it ever reaches the wire.
type readRequestBuilder struct {
	// The default builder does the rest of the work, and it only ever sees tags that are already
	// parsed - which is why its Build can no longer fail.
	apiModel.PlcReadRequestBuilder

	tagHandler spi.PlcTagHandler
	log        zerolog.Logger
}

// newReadRequestBuilder wraps the default builder, which is handed the reader and assembles the
// request itself.
func newReadRequestBuilder(delegate apiModel.PlcReadRequestBuilder, tagHandler spi.PlcTagHandler, log zerolog.Logger) apiModel.PlcReadRequestBuilder {
	return &readRequestBuilder{
		PlcReadRequestBuilder: delegate,
		tagHandler:            tagHandler,
		log:                   log,
	}
}

// AddTagAddress parses the address right away, so that the tag the request ends up with is either
// a modbus tag or the record of why there isn't one. The default builder is only ever handed
// already parsed tags, so it never parses anything of its own.
func (b *readRequestBuilder) AddTagAddress(tagName string, tagAddress string) apiModel.PlcReadRequestBuilder {
	tag, err := b.tagHandler.ParseTag(tagAddress)
	if err != nil {
		// The response carries only the code, so this is the one place the reason is known.
		b.log.Warn().Err(err).Str("tagName", tagName).Str("tagAddress", tagAddress).
			Msg("Couldn't parse an address, so this one tag is answered with INVALID_ADDRESS")
		b.PlcReadRequestBuilder.AddTag(tagName, unparsedTag{address: tagAddress, err: err})
		return b
	}
	b.PlcReadRequestBuilder.AddTag(tagName, tag)
	return b
}

// AddTag keeps the chain on this builder; a tag that is handed over already parsed needs nothing
// else done to it.
func (b *readRequestBuilder) AddTag(tagName string, tag apiModel.PlcTag) apiModel.PlcReadRequestBuilder {
	b.PlcReadRequestBuilder.AddTag(tagName, tag)
	return b
}

var _ apiModel.PlcTag = unparsedTag{}

// unparsedTag is a tag whose address this driver couldn't make sense of. It is a tag only so far
// as it can take part in a request and be reported back; it carries no address the driver could
// read, which is the whole reason it exists.
type unparsedTag struct {
	address string
	err     error
}

func (u unparsedTag) GetAddressString() string { return u.address }

// GetValueType is NULL, the value a tag that can't be read comes back with (see valueOrNull).
func (u unparsedTag) GetValueType() apiValues.PlcValueType { return apiValues.NULL }

func (u unparsedTag) GetArrayInfo() []apiModel.ArrayInfo { return nil }

func (u unparsedTag) String() string {
	return fmt.Sprintf("unparsedTag{address: %s, err: %v}", u.address, u.err)
}

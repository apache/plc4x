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

package model

import (
	"context"
	"runtime/debug"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/interceptors"
	"github.com/pkg/errors"
)

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=DefaultPlcReadRequestBuilder
type DefaultPlcReadRequestBuilder struct {
	reader                 spi.PlcReader     `ignore:"true"`
	tagHandler             spi.PlcTagHandler `ignore:"true"`
	tagNames               []string
	tagAddresses           map[string]string
	tags                   map[string]apiModel.PlcTag
	readRequestInterceptor interceptors.ReadRequestInterceptor `ignore:"true"`
}

func NewDefaultPlcReadRequestBuilder(tagHandler spi.PlcTagHandler, reader spi.PlcReader) apiModel.PlcReadRequestBuilder {
	return NewDefaultPlcReadRequestBuilderWithInterceptor(tagHandler, reader, nil)
}

func NewDefaultPlcReadRequestBuilderWithInterceptor(tagHandler spi.PlcTagHandler, reader spi.PlcReader, readRequestInterceptor interceptors.ReadRequestInterceptor) apiModel.PlcReadRequestBuilder {
	return &DefaultPlcReadRequestBuilder{
		reader:                 reader,
		tagHandler:             tagHandler,
		tagNames:               make([]string, 0),
		tagAddresses:           map[string]string{},
		tags:                   map[string]apiModel.PlcTag{},
		readRequestInterceptor: readRequestInterceptor,
	}
}

func (d *DefaultPlcReadRequestBuilder) AddTagAddress(name string, query string) apiModel.PlcReadRequestBuilder {
	d.tagNames = append(d.tagNames, name)
	d.tagAddresses[name] = query
	return d
}

func (d *DefaultPlcReadRequestBuilder) AddTag(name string, tag apiModel.PlcTag) apiModel.PlcReadRequestBuilder {
	d.tagNames = append(d.tagNames, name)
	d.tags[name] = tag
	return d
}

func (d *DefaultPlcReadRequestBuilder) Build() (apiModel.PlcReadRequest, error) {
	for _, name := range d.tagNames {
		if tagAddress, ok := d.tagAddresses[name]; ok {
			tag, err := d.tagHandler.ParseTag(tagAddress)
			if err != nil {
				return nil, errors.Wrapf(err, "Error parsing tag query: %s", tagAddress)
			}
			d.tags[name] = tag
		}
	}
	// Reset the queries
	d.tagAddresses = map[string]string{}

	return NewDefaultPlcReadRequest(d.tags, d.tagNames, d.reader, d.readRequestInterceptor), nil
}

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=DefaultPlcReadRequest
type DefaultPlcReadRequest struct {
	*DefaultPlcTagRequest
	reader                 spi.PlcReader                       `ignore:"true"`
	readRequestInterceptor interceptors.ReadRequestInterceptor `ignore:"true"`
}

func NewDefaultPlcReadRequest(
	tags map[string]apiModel.PlcTag,
	tagNames []string,
	reader spi.PlcReader,
	readRequestInterceptor interceptors.ReadRequestInterceptor,
) apiModel.PlcReadRequest {
	return &DefaultPlcReadRequest{
		DefaultPlcTagRequest:   NewDefaultPlcTagRequest(tags, tagNames),
		reader:                 reader,
		readRequestInterceptor: readRequestInterceptor,
	}
}

func (d *DefaultPlcReadRequest) GetReader() spi.PlcReader {
	return d.reader
}

func (d *DefaultPlcReadRequest) GetReadRequestInterceptor() interceptors.ReadRequestInterceptor {
	return d.readRequestInterceptor
}
func (d *DefaultPlcReadRequest) Execute() <-chan apiModel.PlcReadRequestResult {
	return d.ExecuteWithContext(context.TODO())
}

func (d *DefaultPlcReadRequest) ExecuteWithContext(ctx context.Context) <-chan apiModel.PlcReadRequestResult {
	if d.readRequestInterceptor != nil {
		return d.ExecuteWithContextAndInterceptor(ctx)
	}

	return d.reader.Read(ctx, d)
}

func (d *DefaultPlcReadRequest) ExecuteWithContextAndInterceptor(ctx context.Context) <-chan apiModel.PlcReadRequestResult {
	// Split the requests up into multiple ones.
	readRequests := d.readRequestInterceptor.InterceptReadRequest(ctx, d)

	// Shortcut for single-request-requests
	if len(readRequests) == 1 {
		return d.reader.Read(ctx, readRequests[0])
	}
	// Create a sub-result-channel slice
	var subResultChannels []<-chan apiModel.PlcReadRequestResult

	// Iterate over all requests and add the result-channels to the list
	for _, subRequest := range readRequests {
		subResultChannels = append(subResultChannels, d.reader.Read(ctx, subRequest))
		// TODO: Replace this with a real queueing of requests. Later on we need throttling. At the moment this avoids race condition as the read above writes to fast on the line which is a problem for the test
		time.Sleep(time.Millisecond * 4)
	}

	// Create a new result-channel, which completes as soon as all sub-result-channels have returned
	resultChannel := make(chan apiModel.PlcReadRequestResult, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				resultChannel <- NewDefaultPlcReadRequestResult(d, nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack()))
			}
		}()
		var subResults []apiModel.PlcReadRequestResult
		// Iterate over all sub-results
		for _, subResultChannel := range subResultChannels {
			select {
			case <-ctx.Done():
				resultChannel <- NewDefaultPlcReadRequestResult(d, nil, ctx.Err())
				return
			case subResult := <-subResultChannel:
				subResults = append(subResults, subResult)
			}
		}
		// As soon as all are done, process the results
		result := d.readRequestInterceptor.ProcessReadResponses(ctx, d, subResults)
		// Return the final result
		resultChannel <- result
	}()
	return resultChannel
}

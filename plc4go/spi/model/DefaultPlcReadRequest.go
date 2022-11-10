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
	"time"

	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/interceptors"
	"github.com/pkg/errors"
)

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=DefaultPlcReadRequestBuilder
type DefaultPlcReadRequestBuilder struct {
	reader                 spi.PlcReader
	tagHandler             spi.PlcTagHandler
	tagNames               []string
	tagAddresses           map[string]string
	tags                   map[string]model.PlcTag
	readRequestInterceptor interceptors.ReadRequestInterceptor
}

func NewDefaultPlcReadRequestBuilder(tagHandler spi.PlcTagHandler, reader spi.PlcReader) *DefaultPlcReadRequestBuilder {
	return NewDefaultPlcReadRequestBuilderWithInterceptor(tagHandler, reader, nil)
}

func NewDefaultPlcReadRequestBuilderWithInterceptor(tagHandler spi.PlcTagHandler, reader spi.PlcReader, readRequestInterceptor interceptors.ReadRequestInterceptor) *DefaultPlcReadRequestBuilder {
	return &DefaultPlcReadRequestBuilder{
		reader:                 reader,
		tagHandler:             tagHandler,
		tagNames:               make([]string, 0),
		tagAddresses:           map[string]string{},
		tags:                   map[string]model.PlcTag{},
		readRequestInterceptor: readRequestInterceptor,
	}
}

func (d *DefaultPlcReadRequestBuilder) AddTagAddress(name string, query string) model.PlcReadRequestBuilder {
	d.tagNames = append(d.tagNames, name)
	d.tagAddresses[name] = query
	return d
}

func (d *DefaultPlcReadRequestBuilder) AddTag(name string, tag model.PlcTag) model.PlcReadRequestBuilder {
	d.tagNames = append(d.tagNames, name)
	d.tags[name] = tag
	return d
}

func (d *DefaultPlcReadRequestBuilder) Build() (model.PlcReadRequest, error) {
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
	DefaultPlcTagRequest
	reader                 spi.PlcReader
	readRequestInterceptor interceptors.ReadRequestInterceptor
}

func NewDefaultPlcReadRequest(tags map[string]model.PlcTag, tagNames []string, reader spi.PlcReader, readRequestInterceptor interceptors.ReadRequestInterceptor) model.PlcReadRequest {
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
func (d *DefaultPlcReadRequest) Execute() <-chan model.PlcReadRequestResult {
	return d.ExecuteWithContext(context.TODO())
}

func (d *DefaultPlcReadRequest) ExecuteWithContext(ctx context.Context) <-chan model.PlcReadRequestResult {
	// Shortcut, if no interceptor is defined
	if d.readRequestInterceptor == nil {
		return d.reader.Read(ctx, d)
	}

	// Split the requests up into multiple ones.
	readRequests := d.readRequestInterceptor.InterceptReadRequest(ctx, d)
	// Shortcut for single-request-requests
	if len(readRequests) == 1 {
		return d.reader.Read(ctx, readRequests[0])
	}
	// Create a sub-result-channel slice
	var subResultChannels []<-chan model.PlcReadRequestResult

	// Iterate over all requests and add the result-channels to the list
	for _, subRequest := range readRequests {
		subResultChannels = append(subResultChannels, d.reader.Read(ctx, subRequest))
		// TODO: Replace this with a real queueing of requests. Later on we need throttling. At the moment this avoids race condition as the read above writes to fast on the line which is a problem for the test
		time.Sleep(time.Millisecond * 4)
	}

	// Create a new result-channel, which completes as soon as all sub-result-channels have returned
	resultChannel := make(chan model.PlcReadRequestResult)
	go func() {
		var subResults []model.PlcReadRequestResult
		// Iterate over all sub-results
		for _, subResultChannel := range subResultChannels {
			select {
			case <-ctx.Done():
				resultChannel <- &DefaultPlcReadRequestResult{Request: d, Err: ctx.Err()}
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

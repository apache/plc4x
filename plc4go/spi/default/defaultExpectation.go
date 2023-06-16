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

package _default

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/plc4x/plc4go/spi"
)

type defaultExpectation struct {
	Context        context.Context
	CreationTime   time.Time
	Expiration     time.Time
	AcceptsMessage spi.AcceptsMessage
	HandleMessage  spi.HandleMessage
	HandleError    spi.HandleError
}

func (d *defaultExpectation) GetContext() context.Context {
	return d.Context
}

func (d *defaultExpectation) GetCreationTime() time.Time {
	return d.CreationTime
}

func (d *defaultExpectation) GetExpiration() time.Time {
	return d.Expiration
}

func (d *defaultExpectation) GetAcceptsMessage() spi.AcceptsMessage {
	return d.AcceptsMessage
}

func (d *defaultExpectation) GetHandleMessage() spi.HandleMessage {
	return d.HandleMessage
}

func (d *defaultExpectation) GetHandleError() spi.HandleError {
	return d.HandleError
}

func (d *defaultExpectation) String() string {
	return fmt.Sprintf("Expectation(expires at %v)", d.Expiration)
}

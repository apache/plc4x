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
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/interceptors"
)

// Note this file is a Helper for mockery to generate use mocks from other package

// Deprecated: don't use it in productive code
type PlcBrowser interface {
	spi.PlcBrowser
}

// Deprecated: don't use it in productive code
type PlcReader interface {
	spi.PlcReader
}

// Deprecated: don't use it in productive code
type PlcWriter interface {
	spi.PlcWriter
}

// Deprecated: don't use it in productive code
type PlcSubscriber interface {
	spi.PlcSubscriber
}

// Deprecated: don't use it in productive code
type PlcTagHandler interface {
	spi.PlcTagHandler
}

// Deprecated: don't use it in productive code
type PlcValueHandler interface {
	spi.PlcValueHandler
}

// Deprecated: don't use it in productive code
type PlcTag interface {
	apiModel.PlcTag
}

// Deprecated: don't use it in productive code
type PlcBrowseItem interface {
	apiModel.PlcBrowseItem
}

// Deprecated: don't use it in productive code
type PlcQuery interface {
	apiModel.PlcQuery
}

// Deprecated: don't use it in productive code
type ReadRequestInterceptor interface {
	interceptors.ReadRequestInterceptor
}

// Deprecated: don't use it in productive code
type WriteRequestInterceptor interface {
	interceptors.WriteRequestInterceptor
}

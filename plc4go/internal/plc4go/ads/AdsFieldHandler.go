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
package ads

import (
	"errors"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	apiModel "github.com/apache/plc4x/plc4go/pkg/plc4go/model"
)

type FieldHandler struct {
	spi.PlcFieldHandler
}

func NewFieldHandler() FieldHandler {
	return FieldHandler{}
}

func (m FieldHandler) ParseQuery(query string) (apiModel.PlcField, error) {
	return nil, errors.New("Invalid address format for address '" + query + "'")
}

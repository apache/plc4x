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

package capability

import (
	"fmt"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type Capability interface {
	fmt.Stringer
	utils.Serializable
	getFN(fn string) GenericFunction
	getZIndex() int
}

//go:generate go tool plc4xGenerator -type=capability -prefix=capability_
type capability struct {
	_zindex int
}

var _ Capability = (*capability)(nil)

func NewCapability(options ...Option) Capability {
	return &capability{_zindex: 99}
}

func (c *capability) getFN(_ string) GenericFunction {
	return nil
}

func (c *capability) getZIndex() int {
	return c._zindex
}

func (c *capability) IsCollector() bool {
	return false
}

func (c *capability) IsCapability() bool {
	return true
}

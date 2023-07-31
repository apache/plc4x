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
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
)

var _ apiModel.PlcSubscriptionTag = &DefaultPlcSubscriptionTag{}

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=DefaultPlcSubscriptionTag
type DefaultPlcSubscriptionTag struct {
	plcTag              apiModel.PlcTag              `stringer:"true"`
	plcSubscriptionType apiModel.PlcSubscriptionType `stringer:"true"`
	duration            time.Duration                `stringer:"true"`
}

func NewDefaultPlcSubscriptionTag(plcSubscriptionType apiModel.PlcSubscriptionType, plcTag apiModel.PlcTag, duration time.Duration) *DefaultPlcSubscriptionTag {
	return &DefaultPlcSubscriptionTag{
		plcTag:              plcTag,
		plcSubscriptionType: plcSubscriptionType,
		duration:            duration,
	}
}

func (d *DefaultPlcSubscriptionTag) GetAddressString() string {
	return d.plcTag.GetAddressString()
}

func (d *DefaultPlcSubscriptionTag) GetValueType() apiValues.PlcValueType {
	return d.plcTag.GetValueType()
}

func (d *DefaultPlcSubscriptionTag) GetArrayInfo() []apiModel.ArrayInfo {
	return d.plcTag.GetArrayInfo()
}

func (d *DefaultPlcSubscriptionTag) GetPlcSubscriptionType() apiModel.PlcSubscriptionType {
	return d.plcSubscriptionType
}

func (d *DefaultPlcSubscriptionTag) GetDuration() time.Duration {
	return d.duration
}

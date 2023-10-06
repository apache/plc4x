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
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
)

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=DefaultPlcSubscriptionEventItem
type DefaultPlcSubscriptionEventItem struct {
	code             apiModel.PlcResponseCode `stringer:"true"`
	tag              apiModel.PlcTag
	subscriptionType apiModel.PlcSubscriptionType `stringer:"true"`
	interval         time.Duration                `stringer:"true"`
	value            apiValues.PlcValue
}

func NewDefaultPlcSubscriptionEventItem(code apiModel.PlcResponseCode, tag apiModel.PlcTag, subscriptionType apiModel.PlcSubscriptionType, interval time.Duration, value apiValues.PlcValue) *DefaultPlcSubscriptionEventItem {
	return &DefaultPlcSubscriptionEventItem{
		code:             code,
		tag:              tag,
		subscriptionType: subscriptionType,
		interval:         interval,
		value:            value,
	}
}

func (d *DefaultPlcSubscriptionEventItem) GetCode() apiModel.PlcResponseCode {
	return d.code
}

func (d *DefaultPlcSubscriptionEventItem) GetTag() apiModel.PlcTag {
	return d.tag
}

func (d *DefaultPlcSubscriptionEventItem) GetSubscriptionType() apiModel.PlcSubscriptionType {
	return d.subscriptionType
}

func (d *DefaultPlcSubscriptionEventItem) GetInterval() time.Duration {
	return d.interval
}

func (d *DefaultPlcSubscriptionEventItem) GetValue() apiValues.PlcValue {
	return d.value
}

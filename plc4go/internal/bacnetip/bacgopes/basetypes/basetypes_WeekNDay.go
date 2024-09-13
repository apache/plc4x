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

package basetypes

import (
	"fmt"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"
)

type WeekNDay struct {
	*OctetString
}

func NewWeekNDay(arg Arg) (*WeekNDay, error) {
	s := &WeekNDay{}
	panic("implement me")
	return s, nil
}

func (w WeekNDay) String() string {
	value := w.GetValue()
	if len(value) != 3 {
		return "WeekNDay(?): " + w.OctetString.String()
	} else {
		return fmt.Sprintf("WeekNDay(%d, %d, %d)", value[0], value[1], value[2])
	}
}

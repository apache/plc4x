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
	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"
)

type Segmentation struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewSegmentation(arg Arg) (*Segmentation, error) {
	s := &Segmentation{
		enumerations: map[string]uint64{
			"segmentedBoth":     0,
			"segmentedTransmit": 1,
			"segmentedReceive":  2,
			"noSegmentation":    3,
		},
	}
	var err error
	args := NoArgs
	if !IsNil(arg) {
		args = append(args, arg)
	}
	s.Enumerated, err = NewEnumerated(args)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create enumerated")
	}
	return s, nil
}

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

package s7

type ControllerType int

const (
	ControllerType_UNKNOWN = -0x1
	ControllerType_ANY     = 0x0
	ControllerType_S7_300  = 0x1
	ControllerType_S7_400  = 0x2
	ControllerType_S7_1200 = 0x3
	ControllerType_S7_1500 = 0x4
	ControllerType_LOGO    = 0x5
)

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

import plc4go "github.com/apache/plc4x/plc4go/pkg/api"

type DefaultPlcConnectionConnectResult interface {
	plc4go.PlcConnectionConnectResult
}

func NewDefaultPlcConnectionConnectResult(connection plc4go.PlcConnection, err error) DefaultPlcConnectionConnectResult {
	return &defaultPlcConnectionConnectResult{
		connection: connection,
		err:        err,
	}
}

//go:generate plc4xGenerator -type=defaultPlcConnectionConnectResult
type defaultPlcConnectionConnectResult struct {
	connection plc4go.PlcConnection
	err        error
}

func (d *defaultPlcConnectionConnectResult) GetConnection() plc4go.PlcConnection {
	return d.connection
}

func (d *defaultPlcConnectionConnectResult) GetErr() error {
	return d.err
}

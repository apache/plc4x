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

package s7

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/s7/readwrite/model"
)

func TestParseFromOptions_controllerType(t *testing.T) {
	t.Run("default is ANY", func(t *testing.T) {
		configuration, err := ParseFromOptions(zerolog.Nop(), map[string][]string{})
		require.NoError(t, err)
		assert.Equal(t, readWriteModel.ControllerType_ANY, configuration.controllerType)
	})
	t.Run("explicit type is parsed", func(t *testing.T) {
		configuration, err := ParseFromOptions(zerolog.Nop(), map[string][]string{
			"controller-type": {"S7_1200"},
		})
		require.NoError(t, err)
		assert.Equal(t, readWriteModel.ControllerType_S7_1200, configuration.controllerType)
	})
	t.Run("LOGO is distinct from S7_200", func(t *testing.T) {
		logoConfiguration, err := ParseFromOptions(zerolog.Nop(), map[string][]string{
			"controller-type": {"LOGO"},
		})
		require.NoError(t, err)
		s7200Configuration, err := ParseFromOptions(zerolog.Nop(), map[string][]string{
			"controller-type": {"S7_200"},
		})
		require.NoError(t, err)
		assert.NotEqual(t, logoConfiguration.controllerType, s7200Configuration.controllerType)
	})
	t.Run("unknown type errors", func(t *testing.T) {
		_, err := ParseFromOptions(zerolog.Nop(), map[string][]string{
			"controller-type": {"S7_9000"},
		})
		assert.Error(t, err)
	})
}

func TestNewDriverContext_logoPduSizeOverride(t *testing.T) {
	configuration, err := ParseFromOptions(zerolog.Nop(), map[string][]string{
		"controller-type": {"LOGO"},
	})
	require.NoError(t, err)
	driverContext, err := NewDriverContext(configuration)
	require.NoError(t, err)
	// LOGO devices only support small PDUs; the default of 1024 gets clamped to 480,
	// which then snaps to the nearest COTP TPDU size minus header overhead.
	assert.LessOrEqual(t, driverContext.PduSize, uint16(496))
	assert.Equal(t, readWriteModel.ControllerType_LOGO, driverContext.ControllerType)
}

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

package pool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_workItem_String(t *testing.T) {
	type fields struct {
		workItemId       int32
		runnable         Runnable
		completionFuture *future
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Simple test",
			fields: fields{
				completionFuture: &future{},
			},
			want: `
╔═workItem══════════════════════════════════════════════════════════════════════╗
║╔═workItemId═╗╔═completionFuture══════════════════════════════════════════════╗║
║║0x00000000 0║║╔═future══════════════════════════════════════════════════════╗║║
║╚════════════╝║║╔═cancelRequested╗╔═interruptRequested╗╔═completed╗╔═errored╗║║║
║              ║║║    b0 false    ║║     b0 false      ║║ b0 false ║║b0 false║║║║
║              ║║╚════════════════╝╚═══════════════════╝╚══════════╝╚════════╝║║║
║              ║╚═════════════════════════════════════════════════════════════╝║║
║              ╚═══════════════════════════════════════════════════════════════╝║
╚═══════════════════════════════════════════════════════════════════════════════╝`[1:],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &workItem{
				workItemId:       tt.fields.workItemId,
				runnable:         tt.fields.runnable,
				completionFuture: tt.fields.completionFuture,
			}
			assert.Equalf(t, tt.want, w.String(), "String()")
		})
	}
}

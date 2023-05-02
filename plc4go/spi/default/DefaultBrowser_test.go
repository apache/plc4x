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

import (
	"context"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
)

func TestNewDefaultBrowser(t *testing.T) {
	type args struct {
		defaultBrowserRequirements DefaultBrowserRequirements
	}
	tests := []struct {
		name string
		args args
		want DefaultBrowser
	}{
		{
			name: "just create it",
			want: &defaultBrowser{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDefaultBrowser(tt.args.defaultBrowserRequirements); !assert.Equal(t, tt.want, got) {
				t.Errorf("NewDefaultBrowser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultBrowser_Browse(t *testing.T) {
	type fields struct {
		DefaultBrowserRequirements DefaultBrowserRequirements
	}
	type args struct {
		ctx           context.Context
		browseRequest apiModel.PlcBrowseRequest
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantAsserter func(t *testing.T, results <-chan apiModel.PlcBrowseRequestResult) bool
	}{
		{
			name: "Browse empty",
			args: args{
				ctx:           context.Background(),
				browseRequest: spiModel.NewDefaultPlcBrowseRequest(nil, nil, nil),
			},
			wantAsserter: func(t *testing.T, results <-chan apiModel.PlcBrowseRequestResult) bool {
				timeout := time.NewTimer(2 * time.Second)
				utils.CleanupTimer(timeout)
				select {
				case result := <-results:
					assert.NotNil(t, result)
				case <-timeout.C:
					t.Error("timeout")
					return false
				}
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &defaultBrowser{
				DefaultBrowserRequirements: tt.fields.DefaultBrowserRequirements,
			}
			assert.Truef(t, tt.wantAsserter(t, m.Browse(tt.args.ctx, tt.args.browseRequest)), "Browse(%v, %v)", tt.args.ctx, tt.args.browseRequest)
		})
	}
}

func Test_defaultBrowser_BrowseWithInterceptor(t *testing.T) {
	type fields struct {
		DefaultBrowserRequirements DefaultBrowserRequirements
	}
	type args struct {
		ctx           context.Context
		browseRequest apiModel.PlcBrowseRequest
		interceptor   func(result apiModel.PlcBrowseItem) bool
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantAsserter func(t *testing.T, results <-chan apiModel.PlcBrowseRequestResult) bool
	}{
		{
			name: "Browse empty",
			args: args{
				ctx:           context.Background(),
				browseRequest: spiModel.NewDefaultPlcBrowseRequest(nil, nil, nil),
			},
			wantAsserter: func(t *testing.T, results <-chan apiModel.PlcBrowseRequestResult) bool {
				timeout := time.NewTimer(2 * time.Second)
				utils.CleanupTimer(timeout)
				select {
				case result := <-results:
					assert.NotNil(t, result)
				case <-timeout.C:
					t.Error("timeout")
					return false
				}
				return true
			},
		},
		{
			name: "Browse something",
			fields: fields{
				DefaultBrowserRequirements: testBrowser{},
			},
			args: args{
				ctx: context.Background(),
				browseRequest: spiModel.NewDefaultPlcBrowseRequest(
					map[string]apiModel.PlcQuery{
						"test": testQuery{},
					},
					[]string{"test"},
					nil,
				),
			},
			wantAsserter: func(t *testing.T, results <-chan apiModel.PlcBrowseRequestResult) bool {
				timeout := time.NewTimer(2 * time.Second)
				utils.CleanupTimer(timeout)
				select {
				case result := <-results:
					assert.NotNil(t, result)
				case <-timeout.C:
					t.Error("timeout")
					return false
				}
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &defaultBrowser{
				DefaultBrowserRequirements: tt.fields.DefaultBrowserRequirements,
			}
			assert.Truef(t, tt.wantAsserter(t, m.Browse(tt.args.ctx, tt.args.browseRequest)), "BrowseWithInterceptor(%v, %v, func())", tt.args.ctx, tt.args.browseRequest)
		})
	}
}

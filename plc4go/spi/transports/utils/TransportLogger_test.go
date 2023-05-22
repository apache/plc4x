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

package utils

import (
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestNewTransportLogger(t *testing.T) {
	type args struct {
		source  io.ReadWriteCloser
		options []Option
	}
	tests := []struct {
		name string
		args args
		want *TransportLogger
	}{
		{
			name: "create it",
			want: &TransportLogger{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTransportLogger(tt.args.source, tt.args.options...); !assert.Equal(t, tt.want, got) {
				t.Errorf("NewTransportLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransportLogger_Close(t1 *testing.T) {
	type fields struct {
		source io.ReadWriteCloser
		log    zerolog.Logger
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TransportLogger{
				source: tt.fields.source,
				log:    tt.fields.log,
			}
			if err := t.Close(); (err != nil) != tt.wantErr {
				t1.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransportLogger_Read(t1 *testing.T) {
	type fields struct {
		source io.ReadWriteCloser
		log    zerolog.Logger
	}
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TransportLogger{
				source: tt.fields.source,
				log:    tt.fields.log,
			}
			got, err := t.Read(tt.args.p)
			if (err != nil) != tt.wantErr {
				t1.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t1.Errorf("Read() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransportLogger_Write(t1 *testing.T) {
	type fields struct {
		source io.ReadWriteCloser
		log    zerolog.Logger
	}
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TransportLogger{
				source: tt.fields.source,
				log:    tt.fields.log,
			}
			got, err := t.Write(tt.args.p)
			if (err != nil) != tt.wantErr {
				t1.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t1.Errorf("Write() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithLogger(t *testing.T) {
	type args struct {
		log zerolog.Logger
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithLogger(tt.args.log); !assert.Equal(t, tt.want, got) {
				t.Errorf("WithLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}

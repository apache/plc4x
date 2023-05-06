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

package options

import (
	"reflect"
	"testing"
)

func TestFilterDiscoveryOptionProtocolSpecific(t *testing.T) {
	type args struct {
		options []WithDiscoveryOption
	}
	tests := []struct {
		name string
		args args
		want []DiscoveryOptionProtocolSpecific
	}{
		{
			name: "nothing",
		},
		{
			name: "find it",
			args: args{options: []WithDiscoveryOption{
				discoveryOptionProtocolSpecific{},
				discoveryOptionProtocolSpecific{},
				discoveryOptionDeviceName{},
			}},
			want: []DiscoveryOptionProtocolSpecific{
				discoveryOptionProtocolSpecific{},
				discoveryOptionProtocolSpecific{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterDiscoveryOptionProtocolSpecific(tt.args.options); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterDiscoveryOptionProtocolSpecific() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterDiscoveryOptionsDeviceName(t *testing.T) {
	type args struct {
		options []WithDiscoveryOption
	}
	tests := []struct {
		name string
		args args
		want []DiscoveryOptionDeviceName
	}{
		{
			name: "nothing",
		},
		{
			name: "find it",
			args: args{options: []WithDiscoveryOption{
				discoveryOptionProtocolSpecific{},
				discoveryOptionDeviceName{},
				discoveryOptionDeviceName{},
			}},
			want: []DiscoveryOptionDeviceName{
				discoveryOptionDeviceName{},
				discoveryOptionDeviceName{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterDiscoveryOptionsDeviceName(tt.args.options); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterDiscoveryOptionsDeviceName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterDiscoveryOptionsLocalAddress(t *testing.T) {
	type args struct {
		options []WithDiscoveryOption
	}
	tests := []struct {
		name string
		args args
		want []DiscoveryOptionLocalAddress
	}{
		{
			name: "nothing",
		},
		{
			name: "find it",
			args: args{options: []WithDiscoveryOption{
				discoveryOptionProtocolSpecific{},
				discoveryOptionLocalAddress{},
				discoveryOptionLocalAddress{},
			}},
			want: []DiscoveryOptionLocalAddress{
				discoveryOptionLocalAddress{},
				discoveryOptionLocalAddress{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterDiscoveryOptionsLocalAddress(tt.args.options); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterDiscoveryOptionsLocalAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterDiscoveryOptionsProtocol(t *testing.T) {
	type args struct {
		options []WithDiscoveryOption
	}
	tests := []struct {
		name string
		args args
		want []DiscoveryOptionProtocol
	}{

		{
			name: "nothing",
		},
		{
			name: "find it",
			args: args{options: []WithDiscoveryOption{
				discoveryOptionProtocolSpecific{},
				discoveryOptionProtocol{},
				discoveryOptionProtocol{},
			}},
			want: []DiscoveryOptionProtocol{
				discoveryOptionProtocol{},
				discoveryOptionProtocol{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterDiscoveryOptionsProtocol(tt.args.options); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterDiscoveryOptionsProtocol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterDiscoveryOptionsRemoteAddress(t *testing.T) {
	type args struct {
		options []WithDiscoveryOption
	}
	tests := []struct {
		name string
		args args
		want []DiscoveryOptionRemoteAddress
	}{

		{
			name: "nothing",
		},
		{
			name: "find it",
			args: args{options: []WithDiscoveryOption{
				discoveryOptionProtocolSpecific{},
				discoveryOptionRemoteAddress{},
				discoveryOptionRemoteAddress{},
			}},
			want: []DiscoveryOptionRemoteAddress{
				discoveryOptionRemoteAddress{},
				discoveryOptionRemoteAddress{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterDiscoveryOptionsRemoteAddress(tt.args.options); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterDiscoveryOptionsRemoteAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterDiscoveryOptionsTransport(t *testing.T) {
	type args struct {
		options []WithDiscoveryOption
	}
	tests := []struct {
		name string
		args args
		want []DiscoveryOptionTransport
	}{

		{
			name: "nothing",
		},
		{
			name: "find it",
			args: args{options: []WithDiscoveryOption{
				discoveryOptionProtocolSpecific{},
				discoveryOptionTransport{},
				discoveryOptionTransport{},
			}},
			want: []DiscoveryOptionTransport{
				discoveryOptionTransport{},
				discoveryOptionTransport{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterDiscoveryOptionsTransport(tt.args.options); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterDiscoveryOptionsTransport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDiscoveryOptionDeviceName(t *testing.T) {
	type args struct {
		deviceName string
	}
	tests := []struct {
		name string
		args args
		want WithDiscoveryOption
	}{
		{
			name: "something",
			args: args{deviceName: "something"},
			want: discoveryOptionDeviceName{deviceName: "something"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDiscoveryOptionDeviceName(tt.args.deviceName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDiscoveryOptionDeviceName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDiscoveryOptionLocalAddress(t *testing.T) {
	type args struct {
		localAddress string
	}
	tests := []struct {
		name string
		args args
		want WithDiscoveryOption
	}{
		{
			name: "something",
			args: args{localAddress: "something"},
			want: discoveryOptionLocalAddress{localAddress: "something"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDiscoveryOptionLocalAddress(tt.args.localAddress); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDiscoveryOptionLocalAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDiscoveryOptionProtocol(t *testing.T) {
	type args struct {
		protocolName string
	}
	tests := []struct {
		name string
		args args
		want WithDiscoveryOption
	}{
		{
			name: "something",
			args: args{protocolName: "something"},
			want: discoveryOptionProtocol{protocolName: "something"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDiscoveryOptionProtocol(tt.args.protocolName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDiscoveryOptionProtocol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDiscoveryOptionProtocolSpecific(t *testing.T) {
	type args struct {
		key   string
		value any
	}
	tests := []struct {
		name string
		args args
		want WithDiscoveryOption
	}{
		{
			name: "something",
			args: args{key: "something", value: "something"},
			want: discoveryOptionProtocolSpecific{key: "something", value: "something"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDiscoveryOptionProtocolSpecific(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDiscoveryOptionProtocolSpecific() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDiscoveryOptionRemoteAddress(t *testing.T) {
	type args struct {
		remoteAddress string
	}
	tests := []struct {
		name string
		args args
		want WithDiscoveryOption
	}{
		{
			name: "something",
			args: args{remoteAddress: "something"},
			want: discoveryOptionRemoteAddress{remoteAddress: "something"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDiscoveryOptionRemoteAddress(tt.args.remoteAddress); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDiscoveryOptionRemoteAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDiscoveryOptionTransport(t *testing.T) {
	type args struct {
		transportName string
	}
	tests := []struct {
		name string
		args args
		want WithDiscoveryOption
	}{
		{
			name: "something",
			args: args{transportName: "something"},
			want: discoveryOptionTransport{transportName: "something"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDiscoveryOptionTransport(tt.args.transportName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDiscoveryOptionTransport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_discoveryOptionDeviceName_GetDeviceName(t *testing.T) {
	type fields struct {
		discoveryOption discoveryOption
		deviceName      string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "something",
			fields: fields{
				deviceName: "something",
			},
			want: "something",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := discoveryOptionDeviceName{
				discoveryOption: tt.fields.discoveryOption,
				deviceName:      tt.fields.deviceName,
			}
			if got := d.GetDeviceName(); got != tt.want {
				t.Errorf("GetDeviceName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_discoveryOptionLocalAddress_GetLocalAddress(t *testing.T) {
	type fields struct {
		discoveryOption discoveryOption
		localAddress    string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "something",
			fields: fields{
				localAddress: "something",
			},
			want: "something",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := discoveryOptionLocalAddress{
				discoveryOption: tt.fields.discoveryOption,
				localAddress:    tt.fields.localAddress,
			}
			if got := d.GetLocalAddress(); got != tt.want {
				t.Errorf("GetLocalAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_discoveryOptionProtocolSpecific_GetKey(t *testing.T) {
	type fields struct {
		discoveryOption discoveryOption
		key             string
		value           any
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "something",
			fields: fields{
				key: "something",
			},
			want: "something",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := discoveryOptionProtocolSpecific{
				discoveryOption: tt.fields.discoveryOption,
				key:             tt.fields.key,
				value:           tt.fields.value,
			}
			if got := d.GetKey(); got != tt.want {
				t.Errorf("GetKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_discoveryOptionProtocolSpecific_GetValue(t *testing.T) {
	type fields struct {
		discoveryOption discoveryOption
		key             string
		value           any
	}
	tests := []struct {
		name   string
		fields fields
		want   any
	}{
		{
			name: "something",
			fields: fields{
				value: "something",
			},
			want: "something",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := discoveryOptionProtocolSpecific{
				discoveryOption: tt.fields.discoveryOption,
				key:             tt.fields.key,
				value:           tt.fields.value,
			}
			if got := d.GetValue(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_discoveryOptionProtocol_GetProtocolName(t *testing.T) {
	type fields struct {
		discoveryOption discoveryOption
		protocolName    string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "something",
			fields: fields{
				protocolName: "something",
			},
			want: "something",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := discoveryOptionProtocol{
				discoveryOption: tt.fields.discoveryOption,
				protocolName:    tt.fields.protocolName,
			}
			if got := d.GetProtocolName(); got != tt.want {
				t.Errorf("GetProtocolName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_discoveryOptionRemoteAddress_GetRemoteAddress(t *testing.T) {
	type fields struct {
		discoveryOption discoveryOption
		remoteAddress   string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "something",
			fields: fields{
				remoteAddress: "something",
			},
			want: "something",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := discoveryOptionRemoteAddress{
				discoveryOption: tt.fields.discoveryOption,
				remoteAddress:   tt.fields.remoteAddress,
			}
			if got := d.GetRemoteAddress(); got != tt.want {
				t.Errorf("GetRemoteAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_discoveryOptionTransport_GetTransportName(t *testing.T) {
	type fields struct {
		discoveryOption discoveryOption
		transportName   string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "something",
			fields: fields{
				transportName: "something",
			},
			want: "something",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := discoveryOptionTransport{
				discoveryOption: tt.fields.discoveryOption,
				transportName:   tt.fields.transportName,
			}
			if got := d.GetTransportName(); got != tt.want {
				t.Errorf("GetTransportName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_discoveryOption_isDiscoveryOption(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "something",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			di := discoveryOption{}
			if got := di.isDiscoveryOption(); got != tt.want {
				t.Errorf("isDiscoveryOption() = %v, want %v", got, tt.want)
			}
		})
	}
}

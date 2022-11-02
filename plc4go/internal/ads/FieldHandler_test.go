package ads

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/apache/plc4x/plc4go/pkg/api/model"
	model2 "github.com/apache/plc4x/plc4go/protocols/ads/readwrite/model"
)

func TestFieldHandler_ParseQuery(t *testing.T) {
	type fields struct {
		directAdsStringField *regexp.Regexp
		directAdsField       *regexp.Regexp
		symbolicAdsField     *regexp.Regexp
	}
	type args struct {
		query string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.PlcField
		wantErr bool
	}{
		// All tests without any array notation.
		{
			name: "simple direct numeric address",
			args: args{
				query: "1234/5678:BOOL",
			},
			want: DirectPlcField{
				PlcField: PlcField{
					NumElements:  NONE,
					StartElement: NONE,
					EndElement:   NONE,
				},
				IndexGroup:   1234,
				IndexOffset:  5678,
				AdsDatatype:  model2.AdsDataType_BOOL,
				StringLength: NONE,
			},
		},
		{
			name: "simple direct hex address",
			args: args{
				query: "0x04D2/0x162E:BOOL",
			},
			want: DirectPlcField{
				PlcField: PlcField{
					NumElements:  NONE,
					StartElement: NONE,
					EndElement:   NONE,
				},
				IndexGroup:   1234,
				IndexOffset:  5678,
				AdsDatatype:  model2.AdsDataType_BOOL,
				StringLength: NONE,
			},
		},
		{
			name: "simple direct numeric string address",
			args: args{
				query: "1234/5678:STRING(80)",
			},
			want: DirectPlcField{
				PlcField: PlcField{
					NumElements:  NONE,
					StartElement: NONE,
					EndElement:   NONE,
				},
				IndexGroup:   1234,
				IndexOffset:  5678,
				AdsDatatype:  model2.AdsDataType_STRING,
				StringLength: 80,
			},
		},
		{
			name: "simple direct hex string address",
			args: args{
				query: "0x04D2/0x162E:WSTRING(80)",
			},
			want: DirectPlcField{
				PlcField: PlcField{
					NumElements:  NONE,
					StartElement: NONE,
					EndElement:   NONE,
				},
				IndexGroup:   1234,
				IndexOffset:  5678,
				AdsDatatype:  model2.AdsDataType_WSTRING,
				StringLength: 80,
			},
		},
		{
			name: "simple symbolic address",
			args: args{
				query: "MAIN.testVariable",
			},
			want: SymbolicPlcField{
				PlcField: PlcField{
					NumElements:  NONE,
					StartElement: NONE,
					EndElement:   NONE,
				},
				SymbolicAddress: "MAIN.testVariable",
			},
		},
		// All tests with simple array notation.
		{
			name: "simple array direct numeric address",
			args: args{
				query: "1234/5678:BOOL[42]",
			},
			want: DirectPlcField{
				PlcField: PlcField{
					NumElements:  42,
					StartElement: NONE,
					EndElement:   NONE,
				},
				IndexGroup:   1234,
				IndexOffset:  5678,
				AdsDatatype:  model2.AdsDataType_BOOL,
				StringLength: NONE,
			},
		},
		{
			name: "simple array direct hex address",
			args: args{
				query: "0x04D2/0x162E:BOOL[42]",
			},
			want: DirectPlcField{
				PlcField: PlcField{
					NumElements:  42,
					StartElement: NONE,
					EndElement:   NONE,
				},
				IndexGroup:   1234,
				IndexOffset:  5678,
				AdsDatatype:  model2.AdsDataType_BOOL,
				StringLength: NONE,
			},
		},
		{
			name: "simple array direct numeric string address",
			args: args{
				query: "1234/5678:STRING(80)[42]",
			},
			want: DirectPlcField{
				PlcField: PlcField{
					NumElements:  42,
					StartElement: NONE,
					EndElement:   NONE,
				},
				IndexGroup:   1234,
				IndexOffset:  5678,
				AdsDatatype:  model2.AdsDataType_STRING,
				StringLength: 80,
			},
		},
		{
			name: "simple array direct hex string address",
			args: args{
				query: "0x04D2/0x162E:WSTRING(80)[42]",
			},
			want: DirectPlcField{
				PlcField: PlcField{
					NumElements:  42,
					StartElement: NONE,
					EndElement:   NONE,
				},
				IndexGroup:   1234,
				IndexOffset:  5678,
				AdsDatatype:  model2.AdsDataType_WSTRING,
				StringLength: 80,
			},
		},
		{
			name: "simple array symbolic address",
			args: args{
				query: "MAIN.testVariable[42]",
			},
			want: SymbolicPlcField{
				PlcField: PlcField{
					NumElements:  42,
					StartElement: NONE,
					EndElement:   NONE,
				},
				SymbolicAddress: "MAIN.testVariable",
			},
		},
		// All tests with range array notation.
		{
			name: "range array direct numeric address",
			args: args{
				query: "1234/5678:BOOL[23..42]",
			},
			want: DirectPlcField{
				PlcField: PlcField{
					NumElements:  NONE,
					StartElement: 23,
					EndElement:   42,
				},
				IndexGroup:   1234,
				IndexOffset:  5678,
				AdsDatatype:  model2.AdsDataType_BOOL,
				StringLength: NONE,
			},
		},
		{
			name: "range array direct hex address",
			args: args{
				query: "0x04D2/0x162E:BOOL[23..42]",
			},
			want: DirectPlcField{
				PlcField: PlcField{
					NumElements:  NONE,
					StartElement: 23,
					EndElement:   42,
				},
				IndexGroup:   1234,
				IndexOffset:  5678,
				AdsDatatype:  model2.AdsDataType_BOOL,
				StringLength: NONE,
			},
		},
		{
			name: "range array direct numeric string address",
			args: args{
				query: "1234/5678:STRING(80)[23..42]",
			},
			want: DirectPlcField{
				PlcField: PlcField{
					NumElements:  NONE,
					StartElement: 23,
					EndElement:   42,
				},
				IndexGroup:   1234,
				IndexOffset:  5678,
				AdsDatatype:  model2.AdsDataType_STRING,
				StringLength: 80,
			},
		},
		{
			name: "range array direct hex string address",
			args: args{
				query: "0x04D2/0x162E:WSTRING(80)[23..42]",
			},
			want: DirectPlcField{
				PlcField: PlcField{
					NumElements:  NONE,
					StartElement: 23,
					EndElement:   42,
				},
				IndexGroup:   1234,
				IndexOffset:  5678,
				AdsDatatype:  model2.AdsDataType_WSTRING,
				StringLength: 80,
			},
		},
		{
			name: "range array symbolic address",
			args: args{
				query: "MAIN.testVariable[23..42]",
			},
			want: SymbolicPlcField{
				PlcField: PlcField{
					NumElements:  NONE,
					StartElement: 23,
					EndElement:   42,
				},
				SymbolicAddress: "MAIN.testVariable",
			},
		},
		// All tests with array with offset notation.
		{
			name: "array with offset direct numeric address",
			args: args{
				query: "1234/5678:BOOL[23:42]",
			},
			want: DirectPlcField{
				PlcField: PlcField{
					NumElements:  42,
					StartElement: 23,
					EndElement:   NONE,
				},
				IndexGroup:   1234,
				IndexOffset:  5678,
				AdsDatatype:  model2.AdsDataType_BOOL,
				StringLength: NONE,
			},
		},
		{
			name: "array with offset direct hex address",
			args: args{
				query: "0x04D2/0x162E:BOOL[23:42]",
			},
			want: DirectPlcField{
				PlcField: PlcField{
					NumElements:  42,
					StartElement: 23,
					EndElement:   NONE,
				},
				IndexGroup:   1234,
				IndexOffset:  5678,
				AdsDatatype:  model2.AdsDataType_BOOL,
				StringLength: NONE,
			},
		},
		{
			name: "array with offset direct numeric string address",
			args: args{
				query: "1234/5678:STRING(80)[23:42]",
			},
			want: DirectPlcField{
				PlcField: PlcField{
					NumElements:  42,
					StartElement: 23,
					EndElement:   NONE,
				},
				IndexGroup:   1234,
				IndexOffset:  5678,
				AdsDatatype:  model2.AdsDataType_STRING,
				StringLength: 80,
			},
		},
		{
			name: "array with offset direct hex string address",
			args: args{
				query: "0x04D2/0x162E:WSTRING(80)[23:42]",
			},
			want: DirectPlcField{
				PlcField: PlcField{
					NumElements:  42,
					StartElement: 23,
					EndElement:   NONE,
				},
				IndexGroup:   1234,
				IndexOffset:  5678,
				AdsDatatype:  model2.AdsDataType_WSTRING,
				StringLength: 80,
			},
		},
		{
			name: "array with offset symbolic address",
			args: args{
				query: "MAIN.testVariable[23:42]",
			},
			want: SymbolicPlcField{
				PlcField: PlcField{
					NumElements:  42,
					StartElement: 23,
					EndElement:   NONE,
				},
				SymbolicAddress: "MAIN.testVariable",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewFieldHandler()
			got, err := m.ParseQuery(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseQuery() got = %v, want %v", got, tt.want)
			}
		})
	}
}

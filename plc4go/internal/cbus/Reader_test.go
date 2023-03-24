package cbus

import (
	"context"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"reflect"
	"testing"
)

func TestNewReader(t *testing.T) {
	type args struct {
		tpduGenerator *AlphaGenerator
		messageCodec  spi.MessageCodec
		tm            *spi.RequestTransactionManager
	}
	tests := []struct {
		name string
		args args
		want *Reader
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewReader(tt.args.tpduGenerator, tt.args.messageCodec, tt.args.tm); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewReader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReader_Read(t *testing.T) {
	type fields struct {
		alphaGenerator *AlphaGenerator
		messageCodec   spi.MessageCodec
		tm             *spi.RequestTransactionManager
	}
	type args struct {
		ctx         context.Context
		readRequest apiModel.PlcReadRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   <-chan apiModel.PlcReadRequestResult
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Reader{
				alphaGenerator: tt.fields.alphaGenerator,
				messageCodec:   tt.fields.messageCodec,
				tm:             tt.fields.tm,
			}
			if got := m.Read(tt.args.ctx, tt.args.readRequest); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReader_readSync(t *testing.T) {
	type fields struct {
		alphaGenerator *AlphaGenerator
		messageCodec   spi.MessageCodec
		tm             *spi.RequestTransactionManager
	}
	type args struct {
		ctx         context.Context
		readRequest apiModel.PlcReadRequest
		result      chan apiModel.PlcReadRequestResult
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Reader{
				alphaGenerator: tt.fields.alphaGenerator,
				messageCodec:   tt.fields.messageCodec,
				tm:             tt.fields.tm,
			}
			m.readSync(tt.args.ctx, tt.args.readRequest, tt.args.result)
		})
	}
}

func TestReader_mapEncodedReply(t *testing.T) {
	type fields struct {
		alphaGenerator *AlphaGenerator
		messageCodec   spi.MessageCodec
		tm             *spi.RequestTransactionManager
	}
	type args struct {
		transaction     *spi.RequestTransaction
		encodedReply    readWriteModel.EncodedReply
		tagName         string
		addResponseCode func(name string, responseCode apiModel.PlcResponseCode)
		addPlcValue     func(name string, plcValue apiValues.PlcValue)
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "empty input",
			args: args{
				transaction: func() *spi.RequestTransaction {
					transactionManager := spi.NewRequestTransactionManager(1)
					transaction := transactionManager.StartTransaction()
					transaction.Submit(func() {
						// NO-OP
					})
					return transaction
				}(),
				encodedReply:    nil,
				tagName:         "",
				addResponseCode: nil,
				addPlcValue:     nil,
			},
		},
		{
			name: "CALDataStatus",
			args: args{
				transaction: func() *spi.RequestTransaction {
					transactionManager := spi.NewRequestTransactionManager(1)
					transaction := transactionManager.StartTransaction()
					transaction.Submit(func() {
						// NO-OP
					})
					return transaction
				}(),
				encodedReply: func() readWriteModel.EncodedReplyCALReply {
					calDataStatus := readWriteModel.NewCALDataStatus(readWriteModel.ApplicationIdContainer_LIGHTING_3A, 0, nil, readWriteModel.CALCommandTypeContainer_CALCommandStatus_0Bytes, nil, nil)
					calReplyShort := readWriteModel.NewCALReplyShort(0, calDataStatus, nil, nil)
					return readWriteModel.NewEncodedReplyCALReply(calReplyShort, 0, nil, nil)
				}(),
				tagName: "someTag",
				addResponseCode: func(name string, responseCode apiModel.PlcResponseCode) {
					// TODO: add assertions
				},
				addPlcValue: func(name string, plcValue apiValues.PlcValue) {
					// TODO: add assertions
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Reader{
				alphaGenerator: tt.fields.alphaGenerator,
				messageCodec:   tt.fields.messageCodec,
				tm:             tt.fields.tm,
			}
			if err := m.mapEncodedReply(tt.args.transaction, tt.args.encodedReply, tt.args.tagName, tt.args.addResponseCode, tt.args.addPlcValue); (err != nil) != tt.wantErr {
				t.Errorf("mapEncodedReply() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

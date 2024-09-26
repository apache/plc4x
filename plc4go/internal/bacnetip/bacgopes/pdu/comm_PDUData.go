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

package pdu

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
)

// Note: upstream this belongs to comm but that would create a circular dependency

type PDUData interface {
	Copyable
	SetPduData([]byte)
	GetPduData() []byte
	Get() (byte, error)
	GetShort() (int16, error)
	GetLong() (int64, error)
	GetData(dlen int) ([]byte, error)
	Put(byte)
	PutData(...byte)
	PutShort(uint16)
	PutLong(uint32)
}

// _PDUData is basically a bridge to spi.Message
type _PDUData struct {
	data []byte
}

var _ PDUData = (*_PDUData)(nil)

func NewPDUData(args Args, kwArgs KWArgs, _ ...Option) PDUData {
	data, ok := GAO[any](args, 0, nil)
	if ok {
		args = args[1:]
	}
	if _debug != nil {
		_debug("__init__ %r %r %r", data, args, kwArgs)
	}
	p := &_PDUData{}
	switch data := data.(type) {
	case []byte:
		p.data = make([]byte, len(data))
		copy(p.data, data)
	case PDU:
		otherData := data.GetPduData()
		p.data = make([]byte, len(otherData))
		copy(p.data, otherData)
	case PDUData:
		otherData := data.GetPduData()
		p.data = make([]byte, len(otherData))
		copy(p.data, otherData)
	}
	return p
}

func (d *_PDUData) SetPduData(data []byte) {
	d.data = data
}

func (d *_PDUData) GetPduData() []byte {
	return d.data
}

func (d *_PDUData) Get() (byte, error) {
	if d.data == nil || len(d.data) == 0 {
		return 0, io.EOF
	}
	octet := d.data[0]
	d.data = d.data[1:]
	return octet, nil
}

func (d *_PDUData) GetData(dlen int) ([]byte, error) {
	if len(d.data) < dlen {
		return nil, io.EOF
	}
	data := d.data[:dlen]
	d.data = d.data[dlen:]
	if len(data) == 0 {
		return nil, nil
	}
	return data, nil
}

func (d *_PDUData) GetShort() (int16, error) {
	data, err := d.GetData(2)
	if err != nil {
		return 0, err
	}
	return int16(binary.BigEndian.Uint16(data)), nil
}

func (d *_PDUData) GetLong() (int64, error) {
	data, err := d.GetData(4)
	if err != nil {
		return 0, err
	}
	return int64(binary.BigEndian.Uint64(data)), nil
}

func (d *_PDUData) Put(n byte) {
	d.data = append(d.data, n)
}

func (d *_PDUData) PutData(n ...byte) {
	d.data = append(d.data, n...)
}

func (d *_PDUData) PutShort(n uint16) {
	ba := make([]byte, 2)
	binary.BigEndian.PutUint16(ba, n)
	d.data = append(d.data, ba...)
}

func (d *_PDUData) PutLong(n uint32) {
	ba := make([]byte, 4)
	binary.BigEndian.PutUint32(ba, n)
	d.data = append(d.data, ba...)
}

func (d *_PDUData) deepCopy() *_PDUData {
	copyPDUData := *d
	copyPDUData.data = make([]byte, len(d.data))
	copy(copyPDUData.data, d.data)
	return &copyPDUData
}

func (d *_PDUData) DeepCopy() any {
	return d.deepCopy()
}

func (d *_PDUData) PrintDebugContents(indent int, file io.Writer, _ids []uintptr) {
	if indent == 0 {
		indent = 1
	}
	if file == nil {
		file = os.Stderr
	}

	var hexed string
	if len(d.data) > 20 {
		hexed = Btox(d.data[:20], ".") + "..."
	} else {
		hexed = Btox(d.data, ".")
	}
	_, _ = fmt.Fprintf(file, "%spduData = x'%s'\n", strings.Repeat("    ", indent), hexed)
}

/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.java.s7.messages.s7.params;

import org.apache.plc4x.java.exceptions.PlcException;
import org.apache.plc4x.java.exceptions.PlcIoException;
import org.apache.plc4x.java.s7.messages.s7.types.*;

import java.io.DataOutputStream;
import java.io.IOException;

public class ReadVarParameter implements S7Parameter {

    private MemoryArea memoryArea;
    private TransportSize transportSize;
    private final short numElements;
    private short dataBlockNumber;
    private short byteOffset;
    private byte bitOffset;

    public ReadVarParameter(short numElements) {
        this.numElements = numElements;
    }

    public ReadVarParameter(MemoryArea memoryArea, TransportSize transportSize, short numElements, short dataBlockNumber, short byteOffset, byte bitOffset) {
        this.memoryArea = memoryArea;
        this.transportSize = transportSize;
        this.numElements = numElements;
        this.dataBlockNumber = dataBlockNumber;
        this.byteOffset = byteOffset;
        this.bitOffset = bitOffset;
    }

    @Override
    public short getLength() {
        return 14;
    }

    @Override
    public void serialize(DataOutputStream dos) throws PlcException {
        try {
            dos.writeByte(Function.READ_VAR.getCode());
            dos.writeByte((byte) 0x01);    // Item count (Read one variable at a time)
            dos.writeByte(SpecificationType.VARIABLE_SPECIFICATION.getCode());
            dos.writeByte((byte) 0x0a);    // Length of this item (excluding spec type and length)
            dos.writeByte(VariableAddressingMode.S7ANY.getCode());
            dos.writeByte(transportSize.getCode());
            dos.writeShort(numElements);
            dos.writeShort(dataBlockNumber);
            dos.writeByte(memoryArea.getCode());
            dos.writeShort(byteOffset);
            dos.writeByte(bitOffset);
        } catch (IOException e) {
            throw new PlcIoException("Error serializing message", e);
        }
    }

}

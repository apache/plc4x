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
package org.apache.plc4x.java.s7.netty.model.params.items;

import org.apache.plc4x.java.s7.netty.model.types.MemoryArea;
import org.apache.plc4x.java.s7.netty.model.types.SpecificationType;
import org.apache.plc4x.java.s7.netty.model.types.TransportSize;
import org.apache.plc4x.java.s7.netty.model.types.VariableAddressingMode;

public class S7AnyVarParameterItem implements VarParameterItem {

    private final SpecificationType specificationType;
    private final MemoryArea memoryArea;
    private final TransportSize transportSize;
    private final short numElements;
    private final short dataBlockNumber;
    private final short byteOffset;
    private final byte bitOffset;

    public S7AnyVarParameterItem(SpecificationType specificationType, MemoryArea memoryArea, TransportSize transportSize, short numElements, short dataBlockNumber, short byteOffset, byte bitOffset) {
        this.specificationType = specificationType;
        this.memoryArea = memoryArea;
        this.transportSize = transportSize;
        this.numElements = numElements;
        this.dataBlockNumber = dataBlockNumber;
        this.byteOffset = byteOffset;
        this.bitOffset = bitOffset;
    }

    @Override
    public VariableAddressingMode getAddressingMode() {
        return VariableAddressingMode.S7ANY;
    }

    public SpecificationType getSpecificationType() {
        return specificationType;
    }

    public MemoryArea getMemoryArea() {
        return memoryArea;
    }

    public TransportSize getTransportSize() {
        return transportSize;
    }

    public short getNumElements() {
        return numElements;
    }

    public short getDataBlockNumber() {
        return dataBlockNumber;
    }

    public short getByteOffset() {
        return byteOffset;
    }

    public byte getBitOffset() {
        return bitOffset;
    }

}

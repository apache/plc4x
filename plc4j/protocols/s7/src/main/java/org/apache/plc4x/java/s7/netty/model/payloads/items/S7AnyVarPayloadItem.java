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
package org.apache.plc4x.java.s7.netty.model.payloads.items;

import org.apache.plc4x.java.s7.netty.model.params.items.*;
import org.apache.plc4x.java.s7.netty.model.types.MemoryArea;
import org.apache.plc4x.java.s7.netty.model.types.SpecificationType;
import org.apache.plc4x.java.s7.netty.model.types.TransportSize;
import org.apache.plc4x.java.s7.netty.model.types.VariableAddressingMode;

/**
 * "Low-level" description of S7 Address range and the necessary size for transportation of values.
 * Is used as Arguments of {@link org.apache.plc4x.java.s7.netty.model.params.VarParameter} object.
 *
 * Contains the information to read one or more sequential values of the same datatype starting from given offset in a memory region
 * and also contains the transportation size of the datatype read.
 *
 * In detail:
 * <ul>
 *     <li>transportSize - {@link TransportSize} of the datatype</li>
 *     <li>numElements - number of consecutive elements to be read</li>
 *     <li>dataBlockNumber - number of the datablock</li>
 *     <li>bit / byteOffset where the adress starts</li>
 * </ul>
 */
public class S7AnyVarPayloadItem implements VarParameterItem {

    private final SpecificationType specificationType;
    private final MemoryArea memoryArea;
    private final TransportSize dataType;
    private final int numElements;
    private final int dataBlockNumber;
    private final int byteOffset;
    private final byte bitOffset;

    public S7AnyVarPayloadItem(SpecificationType specificationType, MemoryArea memoryArea, TransportSize dataType, int numElements, int dataBlockNumber, int byteOffset, byte bitOffset) {
        this.specificationType = specificationType;
        this.memoryArea = memoryArea;
        this.dataType = dataType;
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

    public TransportSize getDataType() {
        return dataType;
    }

    public int getNumElements() {
        return numElements;
    }

    public int getDataBlockNumber() {
        return dataBlockNumber;
    }

    public int getByteOffset() {
        return byteOffset;
    }

    public byte getBitOffset() {
        return bitOffset;
    }

}

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
package org.apache.plc4x.java.s7.model;

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.s7.netty.model.types.MemoryArea;
import org.apache.plc4x.java.s7.netty.model.types.TransportSize;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class S7Field implements PlcField {

    private static final Pattern ADDRESS_PATTERN =
        Pattern.compile("^%(?<memoryArea>.)(?<transferSizeCode>[XBWD]?)(?<byteOffset>\\d{1,5})(.(?<bitOffset>[0-7]))?:(?<dataType>[a-zA-Z_]+)(\\[(?<numElements>\\d+)])?");
    private static final Pattern DATA_BLOCK_ADDRESS_PATTERN =
        Pattern.compile("^%DB(?<blockNumber>\\d{1,4}).DB(?<transferSizeCode>[XBWD]?)(?<byteOffset>\\d{1,5})(.(?<bitOffset>[0-7]))?:(?<dataType>[a-zA-Z_]+)(\\[(?<numElements>\\d+)])?");

    private static final String DATA_TYPE = "dataType";
    private static final String TRANSFER_SIZE_CODE = "transferSizeCode";
    private static final String BLOCK_NUMBER = "blockNumber";
    private static final String BYTE_OFFSET = "byteOffset";
    private static final String BIT_OFFSET = "bitOffset";
    private static final String NUM_ELEMENTS = "numElements";
    private static final String MEMORY_AREA = "memoryArea";

    private final TransportSize dataType;
    private final MemoryArea memoryArea;
    private final short blockNumber;
    private final short byteOffset;
    private final short bitOffset;
    private final int numElements;

    private S7Field(TransportSize dataType, MemoryArea memoryArea, short blockNumber, short byteOffset, short bitOffset, int numElements) {
        this.dataType = dataType;
        this.memoryArea = memoryArea;
        this.blockNumber = blockNumber;
        this.byteOffset = byteOffset;
        this.bitOffset = bitOffset;
        this.numElements = numElements;
    }

    public TransportSize getDataType() {
        return dataType;
    }

    public MemoryArea getMemoryArea() {
        return memoryArea;
    }

    public short getBlockNumber() {
        return blockNumber;
    }

    public short getByteOffset() {
        return byteOffset;
    }

    public short getBitOffset() {
        return bitOffset;
    }

    public int getNumElements() {
        return numElements;
    }

    public static boolean matches(String fieldString) {
        return DATA_BLOCK_ADDRESS_PATTERN.matcher(fieldString).matches() ||
            ADDRESS_PATTERN.matcher(fieldString).matches();
    }

    public static S7Field of(String fieldString) {
        Matcher matcher = DATA_BLOCK_ADDRESS_PATTERN.matcher(fieldString);
        if(matcher.matches()) {
            TransportSize dataType = TransportSize.valueOf(matcher.group(DATA_TYPE));
            MemoryArea memoryArea = MemoryArea.DATA_BLOCKS;
            String transferSizeCode = matcher.group(TRANSFER_SIZE_CODE);
            short blockNumber = Short.parseShort(matcher.group(BLOCK_NUMBER));
            short byteOffset = Short.parseShort(matcher.group(BYTE_OFFSET));
            short bitOffset = 0;
            if(matcher.group(BIT_OFFSET) != null) {
                bitOffset = Short.parseShort(matcher.group(BIT_OFFSET));
            } else if(dataType == TransportSize.BOOL) {
                throw new PlcInvalidFieldException("Expected bit offset for BOOL parameters.");
            }
            int numElements = 1;
            if(matcher.group(NUM_ELEMENTS) != null) {
                numElements = Integer.parseInt(matcher.group(NUM_ELEMENTS));
            }
            if(!transferSizeCode.isEmpty() && !dataType.getSizeCode().equals(transferSizeCode)) {
                throw new PlcInvalidFieldException("Transfer size code '" + transferSizeCode +
                    "' doesn't match specified data type '" + dataType.name() + "'");
            }
            return new S7Field(dataType, memoryArea, blockNumber, byteOffset, bitOffset, numElements);
        } else {
            matcher = ADDRESS_PATTERN.matcher(fieldString);
            if (matcher.matches()) {
                TransportSize dataType = TransportSize.valueOf(matcher.group(DATA_TYPE));
                MemoryArea memoryArea = MemoryArea.valueOfShortName(matcher.group(MEMORY_AREA));
                String transferSizeCode = matcher.group(TRANSFER_SIZE_CODE);
                short byteOffset = Short.parseShort(matcher.group(BYTE_OFFSET));
                short bitOffset = 0;
                if(matcher.group(BIT_OFFSET) != null) {
                    bitOffset = Short.parseShort(matcher.group(BIT_OFFSET));
                } else if(dataType == TransportSize.BOOL) {
                    throw new PlcInvalidFieldException("Expected bit offset for BOOL parameters.");
                }
                int numElements = 1;
                if(matcher.group(NUM_ELEMENTS) != null) {
                    numElements = Integer.parseInt(matcher.group(NUM_ELEMENTS));
                }
                if(!transferSizeCode.isEmpty() && !dataType.getSizeCode().equals(transferSizeCode)) {
                    throw new PlcInvalidFieldException("Transfer size code '" + transferSizeCode +
                        "' doesn't match specified data type '" + dataType.name() + "'");
                }
                return new S7Field(dataType, memoryArea, (short) 0, byteOffset, bitOffset, numElements);
            }
        }
        throw new PlcInvalidFieldException("Unable to parse address: " + fieldString);
    }

}

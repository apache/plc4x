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
package org.apache.plc4x.java.s7.tag;

import org.apache.plc4x.java.s7.readwrite.MemoryArea;
import org.apache.plc4x.java.s7.readwrite.TransportSize;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;

import java.nio.charset.StandardCharsets;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;

public class S7StringVarLengthTag extends S7Tag {
    
    public static final Pattern DATA_BLOCK_STRING_VAR_LENGTH_ADDRESS_PATTERN =
        Pattern.compile("^%DB(?<blockNumber>\\d{1,5}).DB(?<transferSizeCode>[XBWD]?)(?<byteOffset>\\d{1,7})(.(?<bitOffset>[0-7]))?:(?<dataType>STRING|WSTRING)(\\[(?<numElements>\\d+)])?");

    public static final Pattern DATA_BLOCK_STRING_VAR_LENGTH_SHORT_PATTERN =
        Pattern.compile("^%DB(?<blockNumber>\\d{1,5}):(?<byteOffset>\\d{1,7})(.(?<bitOffset>[0-7]))?:(?<dataType>STRING|WSTRING)(\\[(?<numElements>\\d+)])?");
     

    protected S7StringVarLengthTag(TransportSize dataType, MemoryArea memoryArea,
                                   int blockNumber, int byteOffset,
                                   byte bitOffset, int numElements) {
        super(dataType, memoryArea, blockNumber, byteOffset, bitOffset, numElements);
    }

    public static boolean matches(String address) {
        return  DATA_BLOCK_STRING_VAR_LENGTH_ADDRESS_PATTERN.matcher(address).matches() ||
                DATA_BLOCK_STRING_VAR_LENGTH_SHORT_PATTERN.matcher(address).matches();
    }

    @Override
    public String toString() {
        return "S7StringVarLengthTag{" +
            "dataType=" + getDataType() +
            ", memoryArea=" + getMemoryArea() +
            ", blockNumber=" + getBlockNumber() +
            ", byteOffset=" + getByteOffset() +
            ", bitOffset=" + getBitOffset() +
            ", numElements=" + getNumberOfElements() +
            '}';
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws BufferException {
        writeBuffer.pushContext();

        String memoryArea = getMemoryArea().name();
        writeBuffer.writeString(memoryArea.getBytes(StandardCharsets.UTF_8).length * 8,
            memoryArea, WithOption.WithEncoding("UTF8"));

        writeBuffer.writeUnsignedInt(16, getBlockNumber());
        writeBuffer.writeUnsignedInt(16, getByteOffset());
        writeBuffer.writeUnsignedInt(8, getBitOffset());
        writeBuffer.writeUnsignedInt(16, getNumberOfElements());

        String dataType = getDataType().name();
        writeBuffer.writeString(dataType.getBytes(StandardCharsets.UTF_8).length * 8,
            dataType, WithOption.WithEncoding("UTF8"));

        writeBuffer.popContext();
    }
    
    public static S7StringVarLengthTag of(String address) {
        Matcher matcher;
        
        if ((matcher = DATA_BLOCK_STRING_VAR_LENGTH_ADDRESS_PATTERN.matcher(address)).matches()) {
            TransportSize dataType = TransportSize.valueOf(matcher.group(DATA_TYPE));
            MemoryArea memoryArea = MemoryArea.DATA_BLOCKS;
            int blockNumber = checkDataBlockNumber(Integer.parseInt(matcher.group(BLOCK_NUMBER)));
            Short transferSizeCode = getSizeCode(matcher.group(TRANSFER_SIZE_CODE));
            int byteOffset = checkByteOffset(Integer.parseInt(matcher.group(BYTE_OFFSET)));
            byte bitOffset = 0;
            if (matcher.group(BIT_OFFSET) != null) {
                bitOffset = Byte.parseByte(matcher.group(BIT_OFFSET));
            } else if (dataType == TransportSize.BOOL) {
                throw new PlcInvalidTagException("Expected bit offset for BOOL parameters.");
            }
            int numElements = 1;
            if (matcher.group(NUM_ELEMENTS) != null) {
                numElements = Integer.parseInt(matcher.group(NUM_ELEMENTS));
            }

            if ((transferSizeCode != null) && (dataType.getShortName() != transferSizeCode)) {
                throw new PlcInvalidTagException("Transfer size code '" + transferSizeCode +
                    "' doesn't match specified data type '" + dataType.name() + "'");
            }

            return new S7StringVarLengthTag(dataType, memoryArea, blockNumber, byteOffset, bitOffset, numElements);
        } else if ((matcher = DATA_BLOCK_STRING_VAR_LENGTH_SHORT_PATTERN.matcher(address)).matches()) {
            TransportSize dataType = TransportSize.valueOf(matcher.group(DATA_TYPE));
            MemoryArea memoryArea = MemoryArea.DATA_BLOCKS;
            int blockNumber = checkDataBlockNumber(Integer.parseInt(matcher.group(BLOCK_NUMBER)));
            int byteOffset = checkByteOffset(Integer.parseInt(matcher.group(BYTE_OFFSET)));
            byte bitOffset = 0;
            int numElements = 1;
            if (matcher.group(NUM_ELEMENTS) != null) {
                numElements = Integer.parseInt(matcher.group(NUM_ELEMENTS));
            }

            return new S7StringVarLengthTag(dataType, memoryArea, blockNumber,
                byteOffset, bitOffset, numElements);
        }
        
        return null;
    }    

}

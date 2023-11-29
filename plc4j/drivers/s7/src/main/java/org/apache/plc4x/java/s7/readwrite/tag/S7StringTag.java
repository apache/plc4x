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
package org.apache.plc4x.java.s7.readwrite.tag;

import org.apache.plc4x.java.s7.readwrite.MemoryArea;
import org.apache.plc4x.java.s7.readwrite.TransportSize;
import org.apache.plc4x.java.spi.codegen.WithOption;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

import java.nio.charset.StandardCharsets;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import static org.apache.plc4x.java.s7.readwrite.tag.S7Tag.BIT_OFFSET;
import static org.apache.plc4x.java.s7.readwrite.tag.S7Tag.BLOCK_NUMBER;
import static org.apache.plc4x.java.s7.readwrite.tag.S7Tag.BYTE_OFFSET;
import static org.apache.plc4x.java.s7.readwrite.tag.S7Tag.DATA_TYPE;
import static org.apache.plc4x.java.s7.readwrite.tag.S7Tag.NUM_ELEMENTS;
import static org.apache.plc4x.java.s7.readwrite.tag.S7Tag.STRING_LENGTH;
import static org.apache.plc4x.java.s7.readwrite.tag.S7Tag.TRANSFER_SIZE_CODE;
import static org.apache.plc4x.java.s7.readwrite.tag.S7Tag.checkByteOffset;
import static org.apache.plc4x.java.s7.readwrite.tag.S7Tag.checkDataBlockNumber;
import static org.apache.plc4x.java.s7.readwrite.tag.S7Tag.getSizeCode;

public class S7StringTag extends S7Tag {
    
    private static final Pattern DATA_BLOCK_STRING =
        Pattern.compile("^%DB(?<blockNumber>\\d{1,5}):(?<byteOffset>\\d{1,7})(.(?<bitOffset>[0-7]))?:(?<dataType>STRING|WSTRING)?"); 

    private static final Pattern DATA_BLOCK_STRING_SHORT =
        Pattern.compile("^%DB(?<blockNumber>\\d{1,5}).DB(?<transferSizeCode>[XBWD]?)(?<byteOffset>\\d{1,7})(.(?<bitOffset>[0-7]))?:(?<dataType>STRING|WSTRING)?");      
    
    private static final Pattern DATA_BLOCK_STRING_ADDRESS_PATTERN =
        Pattern.compile("^%DB(?<blockNumber>\\d{1,5}).DB(?<transferSizeCode>[XBWD]?)(?<byteOffset>\\d{1,7})(.(?<bitOffset>[0-7]))?:(?<dataType>STRING|WSTRING)\\((?<stringLength>\\d{1,3})\\)(\\[(?<numElements>\\d+)])?");
    
    private static final Pattern DATA_BLOCK_STRING_SHORT_PATTERN =
        Pattern.compile("^%DB(?<blockNumber>\\d{1,5}):(?<byteOffset>\\d{1,7})(.(?<bitOffset>[0-7]))?:(?<dataType>STRING|WSTRING)\\((?<stringLength>\\d{1,3})\\)(\\[(?<numElements>\\d+)])?");   
    
     

    private final int stringLength;

    protected S7StringTag(TransportSize dataType, MemoryArea memoryArea,
                          int blockNumber, int byteOffset,
                          byte bitOffset, int numElements,
                          int stringLength) {
        super(dataType, memoryArea, blockNumber, byteOffset, bitOffset, numElements);
        this.stringLength = stringLength;
    }

    public int getStringLength() {
        return stringLength;
    }

    public static boolean matches(String address) {
        return  DATA_BLOCK_STRING.matcher(address).matches() ||
                DATA_BLOCK_STRING_SHORT.matcher(address).matches() ||
                DATA_BLOCK_STRING_ADDRESS_PATTERN.matcher(address).matches() ||
                DATA_BLOCK_STRING_SHORT_PATTERN.matcher(address).matches();
    }    
    
    
    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext(getClass().getSimpleName());

        String memoryArea = getMemoryArea().name();
        writeBuffer.writeString("memoryArea",
            memoryArea.getBytes(StandardCharsets.UTF_8).length * 8,
            memoryArea, WithOption.WithEncoding(StandardCharsets.UTF_8.name()));

        writeBuffer.writeUnsignedInt("blockNumber", 16, getBlockNumber());
        writeBuffer.writeUnsignedInt("byteOffset", 16, getByteOffset());
        writeBuffer.writeUnsignedInt("bitOffset", 8, getBitOffset());
        writeBuffer.writeUnsignedInt("numElements", 16, getNumberOfElements());
        writeBuffer.writeInt("stringLength", 64, getStringLength());

        String dataType = getDataType().name();
        writeBuffer.writeString("dataType",
            dataType.getBytes(StandardCharsets.UTF_8).length * 8,
            dataType, WithOption.WithEncoding(StandardCharsets.UTF_8.name()));

        writeBuffer.popContext(getClass().getSimpleName());
    }
    
    public static S7StringTag of(String address) {
        Matcher matcher;
        
        if ((matcher = DATA_BLOCK_STRING_SHORT.matcher(address)).matches()) {
            TransportSize dataType = TransportSize.valueOf(matcher.group(DATA_TYPE));
            int stringLength = -1;
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
//            if (matcher.group(NUM_ELEMENTS) != null) {
//                numElements = Integer.parseInt(matcher.group(NUM_ELEMENTS));
//            }

            if ((transferSizeCode != null) && (dataType.getShortName() != transferSizeCode)) {
                throw new PlcInvalidTagException("Transfer size code '" + transferSizeCode +
                    "' doesn't match specified data type '" + dataType.name() + "'");
            }

            return new S7StringTag(dataType, memoryArea, blockNumber, byteOffset, bitOffset, numElements, stringLength);
                    
        } else if ((matcher = DATA_BLOCK_STRING.matcher(address)).matches()) {
            TransportSize dataType = TransportSize.valueOf(matcher.group(DATA_TYPE));
            int stringLength = Integer.parseInt(matcher.group(STRING_LENGTH));
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

            return new S7StringTag(dataType, memoryArea, blockNumber, byteOffset, bitOffset, numElements, stringLength);
                    
        } else if ((matcher = DATA_BLOCK_STRING_ADDRESS_PATTERN.matcher(address)).matches()) {            
            TransportSize dataType = TransportSize.valueOf(matcher.group(DATA_TYPE));
            int stringLength = Integer.parseInt(matcher.group(STRING_LENGTH));
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

            return new S7StringTag(dataType, memoryArea, blockNumber, byteOffset, bitOffset, numElements, stringLength);
        } else if ((matcher = DATA_BLOCK_STRING_SHORT_PATTERN.matcher(address)).matches()) {
            TransportSize dataType = TransportSize.valueOf(matcher.group(DATA_TYPE));
            int stringLength = Integer.parseInt(matcher.group(STRING_LENGTH));
            MemoryArea memoryArea = MemoryArea.DATA_BLOCKS;
            int blockNumber = checkDataBlockNumber(Integer.parseInt(matcher.group(BLOCK_NUMBER)));
            int byteOffset = checkByteOffset(Integer.parseInt(matcher.group(BYTE_OFFSET)));
            byte bitOffset = 0;
            int numElements = 1;
            if (matcher.group(NUM_ELEMENTS) != null) {
                numElements = Integer.parseInt(matcher.group(NUM_ELEMENTS));
            }

            return new S7StringTag(dataType, memoryArea, blockNumber,
                byteOffset, bitOffset, numElements, stringLength);   
        }
        
        return null;
    }    

}

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

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.s7.readwrite.MemoryArea;
import org.apache.plc4x.java.s7.readwrite.TransportSize;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;

import java.nio.charset.StandardCharsets;
import java.util.Objects;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class S7StringFixedLengthTag extends S7Tag {

    public static final Pattern DATA_BLOCK_STRING_FIXED_LENGTH_ADDRESS_PATTERN =
        Pattern.compile("^%DB(?<blockNumber>\\d{1,5}).DB(?<transferSizeCode>[XBWD]?)(?<byteOffset>\\d{1,7})(.(?<bitOffset>[0-7]))?" + ARRAY_EXPRESSION + ":(?<dataType>STRING|WSTRING)\\((?<stringLength>\\d{1,3})\\)");

    public static final Pattern DATA_BLOCK_STRING_FIXED_LENGTH_SHORT_PATTERN =
        Pattern.compile("^%DB(?<blockNumber>\\d{1,5}):(?<byteOffset>\\d{1,7})(.(?<bitOffset>[0-7]))?" + ARRAY_EXPRESSION + ":(?<dataType>STRING|WSTRING)\\((?<stringLength>\\d{1,3})\\)");

    private final int stringLength;

    public S7StringFixedLengthTag(TransportSize dataType, MemoryArea memoryArea,
                                     int blockNumber, int byteOffset,
                                     byte bitOffset, int numElements,
                                     int stringLength) {
        super(dataType, memoryArea, blockNumber, byteOffset, bitOffset, numElements);
        this.stringLength = stringLength;
    }

    public int getStringLength() {
        return stringLength;
    }

    /**
     * Spells the tag the way {@link #of(String)} parses it back, which means carrying the
     * declared length: without it the address reads as a variable-length string, and the length
     * is what decides the layout when the string is read or written.
     */
    @Override
    public String getAddressString() {
        // The base form ends in ":STRING"; the declared length belongs directly behind it.
        return super.getAddressString() + "(" + stringLength + ")";
    }

    public static boolean matches(String address) {
        return  DATA_BLOCK_STRING_FIXED_LENGTH_ADDRESS_PATTERN.matcher(address).matches() ||
            DATA_BLOCK_STRING_FIXED_LENGTH_SHORT_PATTERN.matcher(address).matches();
    }

    @Override
    public String toString() {
        return "S7StringFixedLengthTag{" +
            "dataType=" + getDataType() +
            ", memoryArea=" + getMemoryArea() +
            ", blockNumber=" + getBlockNumber() +
            ", byteOffset=" + getByteOffset() +
            ", bitOffset=" + getBitOffset() +
            ", numElements=" + getNumberOfElements() +
            ", stringLength=" + stringLength +
            '}';
    }

    @Override
    public boolean equals(Object o) {
        if (o == null || getClass() != o.getClass()) return false;
        S7StringFixedLengthTag s7StringFixedLengthTag = (S7StringFixedLengthTag) o;
        return getBlockNumber() == s7StringFixedLengthTag.getBlockNumber() && getByteOffset() == s7StringFixedLengthTag.getByteOffset() && getBitOffset() == s7StringFixedLengthTag.getBitOffset() && getNumberOfElements() == s7StringFixedLengthTag.getNumberOfElements() && getDataType() == s7StringFixedLengthTag.getDataType() && getMemoryArea() == s7StringFixedLengthTag.getMemoryArea() && stringLength == s7StringFixedLengthTag.getStringLength();
    }

    @Override
    public int hashCode() {
        return Objects.hash(getDataType(), getMemoryArea(), getBlockNumber(), getByteOffset(), getBitOffset(), getNumberOfElements(), getStringLength());
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
        writeBuffer.writeSignedInt(32, getStringLength());

        String dataType = getDataType().name();
        writeBuffer.writeString(dataType.getBytes(StandardCharsets.UTF_8).length * 8,
            dataType, WithOption.WithEncoding("UTF8"));

        writeBuffer.popContext();
    }

    /**
     * What one string element occupies: its characters plus the two bytes of length that precede
     * them, doubled for a wide string. This is the same arithmetic the optimizer does when it
     * sizes the request, which is where an unchecked count would otherwise overflow.
     */
    private static int bytesPerString(TransportSize dataType, int stringLength) {
        return (stringLength + 2) * (dataType == TransportSize.WSTRING ? 2 : 1);
    }

    public static S7StringFixedLengthTag of(String address) {
        Matcher matcher;
        
        if ((matcher = DATA_BLOCK_STRING_FIXED_LENGTH_ADDRESS_PATTERN.matcher(address)).matches()) {
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
            int[] selection = selectionOf(matcher, address);
            byteOffset += selection[0] * bytesPerString(dataType, stringLength);
            int numElements = checkNumElements(selection[1],
                bytesPerString(dataType, stringLength),
                dataType.name() + "(" + stringLength + ")");

            if ((transferSizeCode != null) && (dataType.getShortName() != transferSizeCode)) {
                throw new PlcInvalidTagException("Transfer size code '" + transferSizeCode +
                    "' doesn't match specified data type '" + dataType.name() + "'");
            }

            return new S7StringFixedLengthTag(dataType, memoryArea, blockNumber, byteOffset, bitOffset, numElements, stringLength);
        } else if ((matcher = DATA_BLOCK_STRING_FIXED_LENGTH_SHORT_PATTERN.matcher(address)).matches()) {
            TransportSize dataType = TransportSize.valueOf(matcher.group(DATA_TYPE));
            int stringLength = Integer.parseInt(matcher.group(STRING_LENGTH));
            MemoryArea memoryArea = MemoryArea.DATA_BLOCKS;
            int blockNumber = checkDataBlockNumber(Integer.parseInt(matcher.group(BLOCK_NUMBER)));
            int byteOffset = checkByteOffset(Integer.parseInt(matcher.group(BYTE_OFFSET)));
            byte bitOffset = 0;
            int[] selection = selectionOf(matcher, address);
            byteOffset += selection[0] * bytesPerString(dataType, stringLength);
            int numElements = checkNumElements(selection[1],
                bytesPerString(dataType, stringLength),
                dataType.name() + "(" + stringLength + ")");

            return new S7StringFixedLengthTag(dataType, memoryArea, blockNumber,
                byteOffset, bitOffset, numElements, stringLength);
        }

        throw new PlcInvalidTagException("Unable to parse address: " + address);
    }

}

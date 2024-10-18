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

import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.s7.readwrite.MemoryArea;
import org.apache.plc4x.java.s7.readwrite.S7Address;
import org.apache.plc4x.java.s7.readwrite.S7AddressAny;
import org.apache.plc4x.java.s7.readwrite.TransportSize;
import org.apache.plc4x.java.spi.codegen.WithOption;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.spi.model.DefaultArrayInfo;
import org.apache.plc4x.java.spi.utils.Serializable;

import java.nio.charset.StandardCharsets;
import java.util.Collections;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class S7Tag implements PlcTag, Serializable {

    //byteOffset theoretically can reach up to 2097151 ... see checkByteOffset() below --> 7digits
    private static final Pattern ADDRESS_PATTERN =
        Pattern.compile("^%(?<memoryArea>.)(?<transferSizeCode>[XBWD]?)(?<byteOffset>\\d{1,7})(.(?<bitOffset>[0-7]))?:(?<dataType>(S5)?[a-zA-Z_]+)(\\[(?<numElements>\\d+)])?");

    //blockNumber usually has its max hat around 64000 --> 5digits
    private static final Pattern DATA_BLOCK_ADDRESS_PATTERN =
        Pattern.compile("^%DB(?<blockNumber>\\d{1,5}).DB(?<transferSizeCode>[XBWD]?)(?<byteOffset>\\d{1,7})(.(?<bitOffset>[0-7]))?:(?<dataType>(S5)?[a-zA-Z_]+)(\\[(?<numElements>\\d+)])?");

    private static final Pattern DATA_BLOCK_SHORT_PATTERN =
        Pattern.compile("^%DB(?<blockNumber>\\d{1,5}):(?<byteOffset>\\d{1,7})(.(?<bitOffset>[0-7]))?:(?<dataType>(S5)?[a-zA-Z_]+)(\\[(?<numElements>\\d+)])?");

    private static final Pattern PLC_PROXY_ADDRESS_PATTERN =
        Pattern.compile("[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}");

    protected static final String DATA_TYPE = "dataType";
    protected static final String STRING_LENGTH = "stringLength";
    protected static final String TRANSFER_SIZE_CODE = "transferSizeCode";
    protected static final String BLOCK_NUMBER = "blockNumber";
    protected static final String BYTE_OFFSET = "byteOffset";
    protected static final String BIT_OFFSET = "bitOffset";
    protected static final String NUM_ELEMENTS = "numElements";
    protected static final String MEMORY_AREA = "memoryArea";

    private final TransportSize dataType;
    private final MemoryArea memoryArea;
    private final int blockNumber;
    private final int byteOffset;
    private final byte bitOffset;
    private final int numElements;

    public S7Tag(TransportSize dataType, MemoryArea memoryArea,
                    int blockNumber, int byteOffset,
                    byte bitOffset, int numElements) {
        this.dataType = dataType;
        this.memoryArea = memoryArea;
        this.blockNumber = blockNumber;
        this.byteOffset = byteOffset;
        this.bitOffset = bitOffset;
        this.numElements = numElements;
    }

    @Override
    public String getAddressString() {
        return null;
    }

    @Override
    public PlcValueType getPlcValueType() {
        // Translate non-standard tag names.
        switch (dataType.name()) {
            case "S5TIME":
                return PlcValueType.TIME;
            case "LDT":
                return PlcValueType.DATE_AND_LTIME;
            case "DTL":
                return PlcValueType.DATE_AND_LTIME;
            default:
                return PlcValueType.valueOf(dataType.name());
        }
    }

    @Override
    public List<ArrayInfo> getArrayInfo() {
        if (numElements != 1) {
            return Collections.singletonList(new DefaultArrayInfo(0, numElements - 1));
        }
        return Collections.emptyList();
    }

    public TransportSize getDataType() {
        return dataType;
    }

    public String getPlcDataType() {
        return dataType.toString();
    }

    public MemoryArea getMemoryArea() {
        return memoryArea;
    }

    public int getBlockNumber() {
        return blockNumber;
    }

    public int getByteOffset() {
        return byteOffset;
    }

    public byte getBitOffset() {
        return bitOffset;
    }

    public int getNumberOfElements() {
        return numElements;
    }

    public static boolean matches(String tagString) {
        return DATA_BLOCK_ADDRESS_PATTERN.matcher(tagString).matches() ||
            DATA_BLOCK_SHORT_PATTERN.matcher(tagString).matches() ||
            PLC_PROXY_ADDRESS_PATTERN.matcher(tagString).matches() ||
            ADDRESS_PATTERN.matcher(tagString).matches();
    }

    public static S7Tag of(String tagString) {
        Matcher matcher;
        if ((matcher = DATA_BLOCK_ADDRESS_PATTERN.matcher(tagString)).matches()) {
            String dataTypeName = matcher.group(DATA_TYPE);
            if("RAW_BYTE_ARRAY".equals(dataTypeName)) {
                dataTypeName = "BYTE";
            }
            TransportSize dataType = TransportSize.valueOf(dataTypeName);
            MemoryArea memoryArea = MemoryArea.DATA_BLOCKS;
            Short transferSizeCode = getSizeCode(matcher.group(TRANSFER_SIZE_CODE));
            int blockNumber = checkDataBlockNumber(Integer.parseInt(matcher.group(BLOCK_NUMBER)));
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

            return new S7Tag(dataType, memoryArea, blockNumber, byteOffset, bitOffset, numElements);
        } else if ((matcher = DATA_BLOCK_SHORT_PATTERN.matcher(tagString)).matches()) {
            String dataTypeName = matcher.group(DATA_TYPE);
            if("RAW_BYTE_ARRAY".equals(dataTypeName)) {
                dataTypeName = "BYTE";
            }
            TransportSize dataType = TransportSize.valueOf(dataTypeName);
            MemoryArea memoryArea = MemoryArea.DATA_BLOCKS;
            int blockNumber = checkDataBlockNumber(Integer.parseInt(matcher.group(BLOCK_NUMBER)));
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

            return new S7Tag(dataType, memoryArea, blockNumber, byteOffset, bitOffset, numElements);
        } else if (PLC_PROXY_ADDRESS_PATTERN.matcher(tagString).matches()) {
            try {
                byte[] addressData = Hex.decodeHex(tagString.replaceAll("[-]", ""));
                ReadBuffer rb = new ReadBufferByteBased(addressData);
                final S7Address s7Address = S7Address.staticParse(rb);
                if (s7Address instanceof S7AddressAny) {
                    S7AddressAny s7AddressAny = (S7AddressAny) s7Address;

                    if ((s7AddressAny.getTransportSize() != TransportSize.BOOL) && s7AddressAny.getBitAddress() != 0) {
                        throw new PlcInvalidTagException("A bit offset other than 0 is only supported for type BOOL");
                    }

                    return new S7Tag(s7AddressAny.getTransportSize(), s7AddressAny.getArea(),
                        s7AddressAny.getDbNumber(), s7AddressAny.getByteAddress(),
                        s7AddressAny.getBitAddress(), s7AddressAny.getNumberOfElements());
                } else {
                    throw new PlcInvalidTagException("Unsupported address type.");
                }
            } catch (ParseException | DecoderException e) {
                throw new PlcInvalidTagException("Unable to parse address: " + tagString);
            }
        } else if ((matcher = ADDRESS_PATTERN.matcher(tagString)).matches()) {
            TransportSize dataType = TransportSize.valueOf(matcher.group(DATA_TYPE));
            MemoryArea memoryArea = getMemoryAreaForShortName(matcher.group(MEMORY_AREA));
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
            if ((dataType != TransportSize.BOOL) && bitOffset != 0) {
                throw new PlcInvalidTagException("A bit offset other than 0 is only supported for type BOOL");
            }

            return new S7Tag(dataType, memoryArea, (short) 0, byteOffset, bitOffset, numElements);
        }
        throw new PlcInvalidTagException("Unable to parse address: " + tagString);
    }

    /**
     * checks if DataBlockNumber of S7Tag is in valid range
     *
     * @param blockNumber given DataBlockNumber
     * @return given blockNumber if Ok, throws PlcInvalidTagException otherwise
     */
    protected static int checkDataBlockNumber(int blockNumber) {
        // TODO: check the value or add reference - limit eventually depending on active S7 --> make a case selection
        if (blockNumber > 64000 || blockNumber < 1) {
            throw new PlcInvalidTagException("DataBlock numbers larger than 64000 or smaller than 1 are not supported.");
        }
        return blockNumber;
    }

    /**
     * checks if ByteOffset from S7Tag is in valid range
     *
     * @param byteOffset given byteOffset
     * @return given byteOffset if Ok, throws PlcInvalidTagException otherwise
     */
    protected static int checkByteOffset(int byteOffset) {
        // TODO: check the value or add reference
        if (byteOffset > 2097151 || byteOffset < 0) {
            throw new PlcInvalidTagException("ByteOffset must be smaller than 2097151 and positive.");
        }
        return byteOffset;
    }

    protected static Short getSizeCode(String value) {
        if ((value == null) || value.isEmpty()) {
            return null;
        }
        if (value.length() > 1) {
            return null;
        }
        return (short) value.getBytes()[0];
    }

    protected static MemoryArea getMemoryAreaForShortName(String shortName) {
        for (MemoryArea memoryArea : MemoryArea.values()) {
            if (memoryArea.getShortName().equals(shortName)) {
                return memoryArea;
            }
        }
        return null;
    }

    @Override
    public String toString() {
        return "S7Tag{" +
            "dataType=" + dataType +
            ", memoryArea=" + memoryArea +
            ", blockNumber=" + blockNumber +
            ", byteOffset=" + byteOffset +
            ", bitOffset=" + bitOffset +
            ", numElements=" + numElements +
            '}';
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

        String dataType = getDataType().name();
        writeBuffer.writeString("dataType",
            dataType.getBytes(StandardCharsets.UTF_8).length * 8,
            dataType, WithOption.WithEncoding(StandardCharsets.UTF_8.name()));

        writeBuffer.popContext(getClass().getSimpleName());
    }

}

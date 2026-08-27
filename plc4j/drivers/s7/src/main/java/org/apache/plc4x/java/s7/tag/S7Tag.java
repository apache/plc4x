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
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.s7.readwrite.MemoryArea;
import org.apache.plc4x.java.s7.readwrite.S7Address;
import org.apache.plc4x.java.s7.readwrite.S7AddressAny;
import org.apache.plc4x.java.s7.readwrite.TransportSize;
import org.apache.plc4x.java.spi.buffers.api.ReadBuffer;
import org.apache.plc4x.java.spi.buffers.api.Serializable;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.drivers.model.AddressConstraints;
import org.apache.plc4x.java.spi.drivers.model.ArrayNotationParser;
import org.apache.plc4x.java.spi.drivers.model.DefaultArrayInfo;

import java.nio.charset.StandardCharsets;
import java.util.Collections;
import java.util.List;
import java.util.Objects;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class S7Tag implements PlcTag, Serializable {

    /** The shared array notation, which sits between the address and the type. */
    protected static final String ARRAY_EXPRESSION = ArrayNotationParser.ARRAY_GROUP;

    protected static final String ARRAY = "array";

    //byteOffset theoretically can reach up to 2097151 ... see checkByteOffset() below --> 7digits
    private static final Pattern ADDRESS_PATTERN =
        Pattern.compile("^%(?<memoryArea>.)(?<transferSizeCode>[XBWD]?)(?<byteOffset>\\d{1,7})(.(?<bitOffset>[0-7]))?" + ARRAY_EXPRESSION + ":(?<dataType>(S5)?[a-zA-Z_]+)");

    //blockNumber usually has its max hat around 64000 --> 5digits
    private static final Pattern DATA_BLOCK_ADDRESS_PATTERN =
        Pattern.compile("^%DB(?<blockNumber>\\d{1,5}).DB(?<transferSizeCode>[XBWD]?)(?<byteOffset>\\d{1,7})(.(?<bitOffset>[0-7]))?" + ARRAY_EXPRESSION + ":(?<dataType>(S5)?[a-zA-Z_]+)");

    private static final Pattern DATA_BLOCK_SHORT_PATTERN =
        Pattern.compile("^%DB(?<blockNumber>\\d{1,5}):(?<byteOffset>\\d{1,7})(.(?<bitOffset>[0-7]))?" + ARRAY_EXPRESSION + ":(?<dataType>(S5)?[a-zA-Z_]+)");

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

    /** highest byte offset an S7 address can name; see checkByteOffset() below */
    protected static final int MAX_BYTE_OFFSET = 2097151;

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
        this.numElements = numElements;
        
        //TODO: Should this address conversion be done in the mspec?
        switch (dataType) {
            case COUNTER: {
                this.bitOffset = (byte) ((byteOffset) & 0x0007);
                this.byteOffset = (byteOffset >> 3);
                break;
            }
            default :{
                this.byteOffset = byteOffset;
                this.bitOffset = bitOffset;                
            }
            
        }
        
    }

    /**
     * Spells the tag the way {@link #of(String)} parses it back, so a tag can be carried as a
     * string - which is what a log line, a browse result or a serialized request needs.
     *
     * <p>The optional transfer size code is left out: it only repeats what the type already
     * says, and an address without it parses to the same tag.</p>
     */
    @Override
    public String getAddressString() {
        StringBuilder sb = new StringBuilder();
        if (memoryArea == MemoryArea.DATA_BLOCKS) {
            // Rendering a data block through its short name would produce "%D100", which parses
            // back as block 0 of the data-block area - a different address.
            sb.append("%DB").append(blockNumber).append(".DB").append(addressedByteOffset());
        } else {
            sb.append('%').append(memoryArea.getShortName()).append(addressedByteOffset());
        }
        // A bit offset is only part of an address for BOOL, and is required there.
        if (dataType == TransportSize.BOOL) {
            sb.append('.').append(bitOffset);
        }
        sb.append(ArrayNotationParser.render(getArrayInfo()));
        return sb.append(':').append(dataType.name()).toString();
    }

    /**
     * The byte offset as the address writes it. A COUNTER address names a counter, which the
     * constructor splits across the byte and bit offsets, so rendering one has to put it back
     * together - otherwise "%DB1.DB100:COUNTER" comes back as counter 12.
     */
    protected int addressedByteOffset() {
        if (dataType == TransportSize.COUNTER) {
            return (byteOffset << 3) | bitOffset;
        }
        return byteOffset;
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
            case "COUNTER":
                return PlcValueType.WORD;                
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
        return S7StringFixedLengthTag.matches(tagString) ||
            S7StringVarLengthTag.matches(tagString) ||
            DATA_BLOCK_ADDRESS_PATTERN.matcher(tagString).matches() ||
            DATA_BLOCK_SHORT_PATTERN.matcher(tagString).matches() ||
            PLC_PROXY_ADDRESS_PATTERN.matcher(tagString).matches() ||
            ADDRESS_PATTERN.matcher(tagString).matches();
    }

    public static S7Tag of(String tagString) {
        Matcher matcher;
        // STRING/WSTRING addresses carry their length - "%DB69:68:STRING(20)" - which none of the
        // patterns below accept. Delegate those to the dedicated subtypes, the same way
        // S7PlcTagHandler.parseTag() does, so both entry points parse the same address space and
        // the declared length isn't lost (it decides the string layout when reading/writing).
        if (S7StringFixedLengthTag.matches(tagString)) {
            return S7StringFixedLengthTag.of(tagString);
        }
        if (S7StringVarLengthTag.matches(tagString)) {
            return S7StringVarLengthTag.of(tagString);
        }
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
            int[] selection = selectionOf(matcher, tagString);
            byteOffset += selection[0] * dataType.getSizeInBytes();
            int numElements = checkNumElements(dataType, selection[1]);

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
            int[] selection = selectionOf(matcher, tagString);
            byteOffset += selection[0] * dataType.getSizeInBytes();
            int numElements = checkNumElements(dataType, selection[1]);

            return new S7Tag(dataType, memoryArea, blockNumber, byteOffset, bitOffset, numElements);
        } else if (PLC_PROXY_ADDRESS_PATTERN.matcher(tagString).matches()) {
            try {
                String hex = tagString.replace("-", "");
                byte[] addressData = new byte[hex.length() / 2];
                for (int i = 0; i < addressData.length; i++) {
                    addressData[i] = (byte) Integer.parseInt(hex.substring(i * 2, i * 2 + 2), 16);
                }
                ReadBuffer rb = new ReadBufferByteBased(addressData);
                final S7Address s7Address = S7Address.staticParse(rb);
                if (s7Address instanceof S7AddressAny s7AddressAny) {
                    if ((s7AddressAny.getTransportSize() != TransportSize.BOOL) && s7AddressAny.getBitAddress() != 0) {
                        throw new PlcInvalidTagException("A bit offset other than 0 is only supported for type BOOL");
                    }
                    return new S7Tag(s7AddressAny.getTransportSize(), s7AddressAny.getArea(),
                        s7AddressAny.getDbNumber(), s7AddressAny.getByteAddress(),
                        s7AddressAny.getBitAddress(),
                        checkNumElements(s7AddressAny.getTransportSize(),
                            s7AddressAny.getNumberOfElements()));
                }
                throw new PlcInvalidTagException("Unsupported address type.");
            } catch (BufferException | NumberFormatException e) {
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
            int[] selection = selectionOf(matcher, tagString);
            byteOffset += selection[0] * dataType.getSizeInBytes();
            int numElements = checkNumElements(dataType, selection[1]);

            if ((transferSizeCode != null) && (dataType.getShortName() != transferSizeCode)) {
                throw new PlcInvalidTagException("Transfer size code '" + transferSizeCode +
                    "' doesn't match specified data type '" + dataType.name() + "'");
            }
            if ((dataType != TransportSize.BOOL) && bitOffset != 0) {
                throw new PlcInvalidTagException("A bit offset other than 0 is only supported for type BOOL");
            }

            return new S7Tag(dataType, memoryArea, (short) 0, byteOffset, bitOffset, numElements);
        }
        throw ArrayNotationParser.invalidAddress(tagString,
            "%{area}{offset}[selection]:{TYPE} - for example %DB42:28.0[0..3]:BYTE");
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
     * checks that the number of elements of an S7Tag spans no more than the addressable area.
     * <p>
     * A count is a request to allocate, so it has to be answerable before anything is allocated
     * for it. The largest span that can be addressed at all is one byte past the highest legal
     * byte offset, and a tag reaching beyond that is asking for memory no device could return.
     *
     * @param dataType    data type of the tag, which fixes what one element costs
     * @param numElements given number of elements
     * @return given numElements if Ok, throws PlcInvalidTagException otherwise
     */
    /**
     * Resolves the address's array expression to the offset of the first element and the number
     * of elements. An absent expression selects one element at the address itself.
     *
     * <p>The offset is in elements; a caller scales it by the data type's size to reach a byte
     * offset.
     *
     * @return {@code {elementOffset, numElements}}
     */
    protected static int[] selectionOf(Matcher matcher, String address) {
        String expression = matcher.group(ARRAY);
        if (expression == null) {
            return new int[]{0, 1};
        }
        ArrayInfo dimension = ArrayNotationParser
            .parse(expression, address, AddressConstraints.SINGLE_DIMENSION).get(0);
        return new int[]{dimension.getLowerBound() - dimension.getBase(), dimension.getSize()};
    }

    protected static int checkNumElements(TransportSize dataType, int numElements) {
        return checkNumElements(numElements, dataType.getSizeInBytes(), dataType.name());
    }

    /**
     * checks that the number of elements spans no more than the addressable area, for a caller
     * that knows what one of its elements costs - a string's element size comes from its declared
     * length rather than from the transport size alone.
     *
     * @param numElements    given number of elements
     * @param bytesPerElement what one element occupies
     * @param describe       how to name the type in the failure
     * @return given numElements if Ok, throws PlcInvalidTagException otherwise
     */
    protected static int checkNumElements(int numElements, int bytesPerElement, String describe) {
        if (numElements < 1) {
            throw new PlcInvalidTagException("The number of elements must be greater than zero.");
        }
        int elementSize = Math.max(1, bytesPerElement);
        if (numElements > (MAX_BYTE_OFFSET + 1) / elementSize) {
            throw new PlcInvalidTagException("A tag of " + numElements + " elements of type " +
                describe + " spans more than the addressable " + (MAX_BYTE_OFFSET + 1) +
                " bytes.");
        }
        return numElements;
    }

    /**
     * checks if ByteOffset from S7Tag is in valid range
     *
     * @param byteOffset given byteOffset
     * @return given byteOffset if Ok, throws PlcInvalidTagException otherwise
     */
    protected static int checkByteOffset(int byteOffset) {
        // TODO: check the value or add reference
        if (byteOffset > MAX_BYTE_OFFSET || byteOffset < 0) {
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
    public boolean equals(Object o) {
        if (o == null || getClass() != o.getClass()) return false;
        S7Tag s7Tag = (S7Tag) o;
        return blockNumber == s7Tag.blockNumber && byteOffset == s7Tag.byteOffset && bitOffset == s7Tag.bitOffset && numElements == s7Tag.numElements && dataType == s7Tag.dataType && memoryArea == s7Tag.memoryArea;
    }

    @Override
    public int hashCode() {
        return Objects.hash(dataType, memoryArea, blockNumber, byteOffset, bitOffset, numElements);
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

}

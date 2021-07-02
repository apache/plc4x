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
package org.apache.plc4x.java.s7.readwrite.field;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeInfo;
import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Hex;
import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.s7.readwrite.S7Address;
import org.apache.plc4x.java.s7.readwrite.S7AddressAny;
import org.apache.plc4x.java.s7.readwrite.io.S7AddressIO;
import org.apache.plc4x.java.s7.readwrite.types.MemoryArea;
import org.apache.plc4x.java.s7.readwrite.types.TransportSize;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.utils.XmlSerializable;
import org.w3c.dom.Document;
import org.w3c.dom.Element;

import java.nio.charset.StandardCharsets;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class S7Field implements PlcField, XmlSerializable {

    //byteOffset theoretically can reach up to 2097151 ... see checkByteOffset() below --> 7digits
    private static final Pattern ADDRESS_PATTERN =
        Pattern.compile("^%(?<memoryArea>.)(?<transferSizeCode>[XBWD]?)(?<byteOffset>\\d{1,7})(.(?<bitOffset>[0-7]))?:(?<dataType>[a-zA-Z_]+)(\\[(?<numElements>\\d+)])?");

    //blockNumber usually has its max hat around 64000 --> 5digits
    private static final Pattern DATA_BLOCK_ADDRESS_PATTERN =
        Pattern.compile("^%DB(?<blockNumber>\\d{1,5}).DB(?<transferSizeCode>[XBWD]?)(?<byteOffset>\\d{1,7})(.(?<bitOffset>[0-7]))?:(?<dataType>[a-zA-Z_]+)(\\[(?<numElements>\\d+)])?");

    private static final Pattern DATA_BLOCK_SHORT_PATTERN =
        Pattern.compile("^%DB(?<blockNumber>\\d{1,5}):(?<byteOffset>\\d{1,7})(.(?<bitOffset>[0-7]))?:(?<dataType>[a-zA-Z_]+)(\\[(?<numElements>\\d+)])?");

    private static final Pattern DATA_BLOCK_STRING_ADDRESS_PATTERN =
        Pattern.compile("^%DB(?<blockNumber>\\d{1,5}).DB(?<transferSizeCode>[XBWD]?)(?<byteOffset>\\d{1,7})(.(?<bitOffset>[0-7]))?:(?<dataType>STRING|WSTRING)\\((?<stringLength>\\d{1,3})\\)(\\[(?<numElements>\\d+)])?");

    private static final Pattern DATA_BLOCK_STRING_SHORT_PATTERN =
        Pattern.compile("^%DB(?<blockNumber>\\d{1,5}):(?<byteOffset>\\d{1,7})(.(?<bitOffset>[0-7]))?:(?<dataType>STRING|WSTRING)\\((?<stringLength>\\d{1,3})\\)(\\[(?<numElements>\\d+)])?");

    private static final Pattern PLC_PROXY_ADDRESS_PATTERN =
        Pattern.compile("[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}");

    private static final String DATA_TYPE = "dataType";
    private static final String STRING_LENGTH = "stringLength";
    private static final String TRANSFER_SIZE_CODE = "transferSizeCode";
    private static final String BLOCK_NUMBER = "blockNumber";
    private static final String BYTE_OFFSET = "byteOffset";
    private static final String BIT_OFFSET = "bitOffset";
    private static final String NUM_ELEMENTS = "numElements";
    private static final String MEMORY_AREA = "memoryArea";

    private final TransportSize dataType;
    private final MemoryArea memoryArea;
    private final int blockNumber;
    private final int byteOffset;
    private final byte bitOffset;
    private final int numElements;

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    protected S7Field(@JsonProperty("dataType") TransportSize dataType, @JsonProperty("memoryArea") MemoryArea memoryArea,
                      @JsonProperty("blockNumber") int blockNumber, @JsonProperty("byteOffset") int byteOffset,
                      @JsonProperty("bitOffset") byte bitOffset, @JsonProperty("numElements") int numElements) {
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

    public static boolean matches(String fieldString) {
        return
            DATA_BLOCK_STRING_ADDRESS_PATTERN.matcher(fieldString).matches() ||
                DATA_BLOCK_STRING_SHORT_PATTERN.matcher(fieldString).matches() ||
                DATA_BLOCK_ADDRESS_PATTERN.matcher(fieldString).matches() ||
                DATA_BLOCK_SHORT_PATTERN.matcher(fieldString).matches() ||
                PLC_PROXY_ADDRESS_PATTERN.matcher(fieldString).matches() ||
                ADDRESS_PATTERN.matcher(fieldString).matches();
    }

    /**
     * @return Java type of expected response.
     * <p>
     * TODO validate all Methods existing are implemented
     */
    @Override
    public Class<?> getDefaultJavaType() {
        switch (dataType) {
            case STRING:
                return String.class;
            case USINT:
            case SINT:
            case UINT:
            case INT:
            case DINT:
                return Integer.class;
            case UDINT:
            case ULINT:
            case LINT:
                return Long.class;
            case BOOL:
                return Boolean.class;
            case REAL:
            case LREAL:
                return Double.class;
            case DATE_AND_TIME:
                return LocalDateTime.class;
            case DATE:
                return LocalDate.class;
            case TIME_OF_DAY:
                return LocalTime.class;
            default:
                throw new NotImplementedException("The response type for datatype " + dataType + " is not yet implemented");
        }
    }

    public static S7Field of(String fieldString) {
        Matcher matcher;
        if ((matcher = DATA_BLOCK_STRING_ADDRESS_PATTERN.matcher(fieldString)).matches()) {
            TransportSize dataType = TransportSize.valueOf(matcher.group(DATA_TYPE));
            int stringLength = Integer.parseInt(matcher.group(STRING_LENGTH));
            MemoryArea memoryArea = MemoryArea.DATA_BLOCKS;
            Short transferSizeCode = getSizeCode(matcher.group(TRANSFER_SIZE_CODE));
            int byteOffset = checkByteOffset(Integer.parseInt(matcher.group(BYTE_OFFSET)));
            byte bitOffset = 0;
            if (matcher.group(BIT_OFFSET) != null) {
                bitOffset = Byte.parseByte(matcher.group(BIT_OFFSET));
            } else if (dataType == TransportSize.BOOL) {
                throw new PlcInvalidFieldException("Expected bit offset for BOOL parameters.");
            }
            int numElements = 1;
            if (matcher.group(NUM_ELEMENTS) != null) {
                numElements = Integer.parseInt(matcher.group(NUM_ELEMENTS));
            }

            if ((transferSizeCode != null) && (dataType.getShortName() != transferSizeCode)) {
                throw new PlcInvalidFieldException("Transfer size code '" + transferSizeCode +
                    "' doesn't match specified data type '" + dataType.name() + "'");
            }

            return new S7StringField(dataType, memoryArea, (short) 0, byteOffset, bitOffset, numElements, stringLength);
        } else if ((matcher = DATA_BLOCK_STRING_SHORT_PATTERN.matcher(fieldString)).matches()) {
            TransportSize dataType = TransportSize.valueOf(matcher.group(DATA_TYPE));
            int stringLength = Integer.parseInt(matcher.group(STRING_LENGTH));
            MemoryArea memoryArea = MemoryArea.DATA_BLOCKS;
            int blockNumber = checkDatablockNumber(Integer.parseInt(matcher.group(BLOCK_NUMBER)));
            int byteOffset = checkByteOffset(Integer.parseInt(matcher.group(BYTE_OFFSET)));
            byte bitOffset = 0;
            int numElements = 1;
            if (matcher.group(NUM_ELEMENTS) != null) {
                numElements = Integer.parseInt(matcher.group(NUM_ELEMENTS));
            }

            return new S7StringField(dataType, memoryArea, blockNumber,
                byteOffset, bitOffset, numElements, stringLength);
        } else if ((matcher = DATA_BLOCK_ADDRESS_PATTERN.matcher(fieldString)).matches()) {
            TransportSize dataType = TransportSize.valueOf(matcher.group(DATA_TYPE));
            MemoryArea memoryArea = MemoryArea.DATA_BLOCKS;
            Short transferSizeCode = getSizeCode(matcher.group(TRANSFER_SIZE_CODE));
            int blockNumber = checkDatablockNumber(Integer.parseInt(matcher.group(BLOCK_NUMBER)));
            int byteOffset = checkByteOffset(Integer.parseInt(matcher.group(BYTE_OFFSET)));
            byte bitOffset = 0;
            if (matcher.group(BIT_OFFSET) != null) {
                bitOffset = Byte.parseByte(matcher.group(BIT_OFFSET));
            } else if (dataType == TransportSize.BOOL) {
                throw new PlcInvalidFieldException("Expected bit offset for BOOL parameters.");
            }
            int numElements = 1;
            if (matcher.group(NUM_ELEMENTS) != null) {
                numElements = Integer.parseInt(matcher.group(NUM_ELEMENTS));
            }

            if ((transferSizeCode != null) && (dataType.getShortName() != transferSizeCode)) {
                throw new PlcInvalidFieldException("Transfer size code '" + transferSizeCode +
                    "' doesn't match specified data type '" + dataType.name() + "'");
            }

            return new S7Field(dataType, memoryArea, blockNumber, byteOffset, bitOffset, numElements);
        } else if ((matcher = DATA_BLOCK_SHORT_PATTERN.matcher(fieldString)).matches()) {
            TransportSize dataType = TransportSize.valueOf(matcher.group(DATA_TYPE));
            MemoryArea memoryArea = MemoryArea.DATA_BLOCKS;
            int blockNumber = checkDatablockNumber(Integer.parseInt(matcher.group(BLOCK_NUMBER)));
            int byteOffset = checkByteOffset(Integer.parseInt(matcher.group(BYTE_OFFSET)));
            byte bitOffset = 0;
            if (matcher.group(BIT_OFFSET) != null) {
                bitOffset = Byte.parseByte(matcher.group(BIT_OFFSET));
            } else if (dataType == TransportSize.BOOL) {
                throw new PlcInvalidFieldException("Expected bit offset for BOOL parameters.");
            }
            int numElements = 1;
            if (matcher.group(NUM_ELEMENTS) != null) {
                numElements = Integer.parseInt(matcher.group(NUM_ELEMENTS));
            }

            return new S7Field(dataType, memoryArea, blockNumber, byteOffset, bitOffset, numElements);
        } else if (PLC_PROXY_ADDRESS_PATTERN.matcher(fieldString).matches()) {
            try {
                byte[] addressData = Hex.decodeHex(fieldString.replaceAll("[-]", ""));
                ReadBuffer rb = new ReadBufferByteBased(addressData);
                final S7Address s7Address = S7AddressIO.staticParse(rb);
                if (s7Address instanceof S7AddressAny) {
                    S7AddressAny s7AddressAny = (S7AddressAny) s7Address;

                    if ((s7AddressAny.getTransportSize() != TransportSize.BOOL) && s7AddressAny.getBitAddress() != 0) {
                        throw new PlcInvalidFieldException("A bit offset other than 0 is only supported for type BOOL");
                    }

                    return new S7Field(s7AddressAny.getTransportSize(), s7AddressAny.getArea(),
                        s7AddressAny.getDbNumber(), s7AddressAny.getByteAddress(),
                        s7AddressAny.getBitAddress(), s7AddressAny.getNumberOfElements());
                } else {
                    throw new PlcInvalidFieldException("Unsupported address type.");
                }
            } catch (ParseException | DecoderException e) {
                throw new PlcInvalidFieldException("Unable to parse address: " + fieldString);
            }
        } else if ((matcher = ADDRESS_PATTERN.matcher(fieldString)).matches()) {
            TransportSize dataType = TransportSize.valueOf(matcher.group(DATA_TYPE));
            MemoryArea memoryArea = getMemoryAreaForShortName(matcher.group(MEMORY_AREA));
            Short transferSizeCode = getSizeCode(matcher.group(TRANSFER_SIZE_CODE));
            int byteOffset = checkByteOffset(Integer.parseInt(matcher.group(BYTE_OFFSET)));
            byte bitOffset = 0;
            if (matcher.group(BIT_OFFSET) != null) {
                bitOffset = Byte.parseByte(matcher.group(BIT_OFFSET));
            } else if (dataType == TransportSize.BOOL) {
                throw new PlcInvalidFieldException("Expected bit offset for BOOL parameters.");
            }
            int numElements = 1;
            if (matcher.group(NUM_ELEMENTS) != null) {
                numElements = Integer.parseInt(matcher.group(NUM_ELEMENTS));
            }

            if ((transferSizeCode != null) && (dataType.getShortName() != transferSizeCode)) {
                throw new PlcInvalidFieldException("Transfer size code '" + transferSizeCode +
                    "' doesn't match specified data type '" + dataType.name() + "'");
            }
            if ((dataType != TransportSize.BOOL) && bitOffset != 0) {
                throw new PlcInvalidFieldException("A bit offset other than 0 is only supported for type BOOL");
            }

            return new S7Field(dataType, memoryArea, (short) 0, byteOffset, bitOffset, numElements);
        }
        throw new PlcInvalidFieldException("Unable to parse address: " + fieldString);
    }

    /**
     * checks if DatablockNumber of S7Field is in valid range
     *
     * @param blockNumber given DatablockNumber
     * @return given blockNumber if Ok, throws PlcInvalidFieldException otherwise
     */
    private static int checkDatablockNumber(int blockNumber) {
        //ToDo check the value or add reference - limit eventually depending on active S7 --> make a case selection
        if (blockNumber > 64000 || blockNumber < 1) {
            throw new PlcInvalidFieldException("Datablock numbers larger than 64000 or smaller than 1 are not supported.");
        }
        return blockNumber;
    }

    /**
     * checks if ByteOffset from S7Field is in valid range
     *
     * @param byteOffset given byteOffset
     * @return given byteOffset if Ok, throws PlcInvalidFieldException otherwise
     */
    private static int checkByteOffset(int byteOffset) {
        //ToDo check the value or add reference
        if (byteOffset > 2097151 || byteOffset < 0) {
            throw new PlcInvalidFieldException("ByteOffset must be smaller than 2097151 and positive.");
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
        return "S7Field{" +
            "dataType=" + dataType +
            ", memoryArea=" + memoryArea +
            ", blockNumber=" + blockNumber +
            ", byteOffset=" + byteOffset +
            ", bitOffset=" + bitOffset +
            ", numElements=" + numElements +
            '}';
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws ParseException {
        writeBuffer.pushContext(getClass().getSimpleName());

        String memoryArea = getMemoryArea().name();
        writeBuffer.writeString("memoryArea", memoryArea.getBytes(StandardCharsets.UTF_8).length * 8, StandardCharsets.UTF_8.name(), memoryArea);

        writeBuffer.writeInt("blockNumber", 64, getBlockNumber());
        writeBuffer.writeInt("byteOffset", 64, getByteOffset());
        writeBuffer.writeInt("bitOffset", 64, getBitOffset());
        writeBuffer.writeInt("numElements", 64, getNumberOfElements());

        String dataType = getDataType().name();
        writeBuffer.writeString("dataType", dataType.getBytes(StandardCharsets.UTF_8).length * 8, StandardCharsets.UTF_8.name(), dataType);

        writeBuffer.popContext(getClass().getSimpleName());
    }

    @Override
    public void xmlSerialize(Element parent) {
        Document doc = parent.getOwnerDocument();
        Element messageElement = doc.createElement(getClass().getSimpleName());
        parent.appendChild(messageElement);

        Element memoryAreaElement = doc.createElement("memoryArea");
        memoryAreaElement.appendChild(doc.createTextNode(getMemoryArea().name()));
        messageElement.appendChild(memoryAreaElement);

        Element blockNumberElement = doc.createElement("blockNumber");
        blockNumberElement.appendChild(doc.createTextNode(Integer.toString(getBlockNumber())));
        messageElement.appendChild(blockNumberElement);

        Element byteOffsetElement = doc.createElement("byteOffset");
        byteOffsetElement.appendChild(doc.createTextNode(Integer.toString(getByteOffset())));
        messageElement.appendChild(byteOffsetElement);

        Element bitOffsetElement = doc.createElement("bitOffset");
        bitOffsetElement.appendChild(doc.createTextNode(Integer.toString(getBitOffset())));
        messageElement.appendChild(bitOffsetElement);

        Element numElementsElement = doc.createElement("numElements");
        numElementsElement.appendChild(doc.createTextNode(Integer.toString(getNumberOfElements())));
        messageElement.appendChild(numElementsElement);

        Element dataTypeElement = doc.createElement("dataType");
        dataTypeElement.appendChild(doc.createTextNode(getDataType().name()));
        messageElement.appendChild(dataTypeElement);
    }

}

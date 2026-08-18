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
package org.apache.plc4x.java.modbus.base;

import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.modbus.readwrite.DataItem;
import org.apache.plc4x.java.modbus.readwrite.ModbusDataType;
import org.apache.plc4x.java.spi.buffers.api.ReadBuffer;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.values.PlcList;

import java.util.ArrayList;
import java.util.List;

/**
 * Reads and writes values as they are laid out in Modbus registers.
 * <p>
 * {@link DataItem} encodes one value at its natural width and knows nothing about registers. What
 * it does not cover is that a register is 16 bits wide: a lone value narrower than that occupies a
 * whole register and is padded, while several of them are packed - two bytes per register, eight
 * bits per byte. Which half of the register a padded value sits in depends on the byte order. That
 * is Modbus knowledge, so it lives here rather than in the protocol description.
 */
public class ModbusRegisterCodec {

    /** Padding carries no value, but the buffers still insist on knowing the encoding. */
    private static final WithOption UNSIGNED_BINARY =
        WithOption.WithUnsignedIntegerEncoding("unsigned-binary");

    private ModbusRegisterCodec() {
    }

    /**
     * Reads {@code numberOfValues} values of the given type.
     *
     * @return the value itself for a single value, a {@link PlcList} for several.
     */
    public static PlcValue parse(ReadBuffer readBuffer, ModbusDataType dataType, int numberOfValues,
                                 boolean bigEndian, int stringLength) throws BufferException {
        if (numberOfValues == 1) {
            int padding = paddingBits(dataType);
            if (padding == 0) {
                return DataItem.staticParse(readBuffer, dataType, stringLength);
            }
            // A value narrower than a register is padded to fill it. Big endian puts the padding
            // first, little endian last - this is what the ',true'/',false' cases used to express.
            if (bigEndian) {
                readBuffer.readUnsignedInt(padding, UNSIGNED_BINARY);
                return DataItem.staticParse(readBuffer, dataType, stringLength);
            }
            int leading = leadingPaddingBitsLittleEndian(dataType);
            if (leading > 0) {
                readBuffer.readUnsignedInt(leading, UNSIGNED_BINARY);
            }
            PlcValue value = DataItem.staticParse(readBuffer, dataType, stringLength);
            readBuffer.readUnsignedInt(padding - leading, UNSIGNED_BINARY);
            return value;
        }

        // Several values are packed without padding between them. A trailing pad, when the count
        // does not fill the last register, is simply left unread.
        List<PlcValue> values = new ArrayList<>(numberOfValues);
        for (int i = 0; i < numberOfValues; i++) {
            values.add(DataItem.staticParse(readBuffer, dataType, stringLength));
        }
        return new PlcList(values);
    }

    /**
     * Writes a value, or every element of a list, with the same layout {@link #parse} expects.
     */
    public static void serialize(WriteBuffer writeBuffer, PlcValue value, ModbusDataType dataType,
                                 int numberOfValues, boolean bigEndian, int stringLength) throws BufferException {
        if (numberOfValues == 1) {
            int padding = paddingBits(dataType);
            if (padding == 0) {
                DataItem.staticSerialize(writeBuffer, value, dataType, stringLength);
                return;
            }
            if (bigEndian) {
                writeBuffer.writeUnsignedInt(padding, 0, UNSIGNED_BINARY);
                DataItem.staticSerialize(writeBuffer, value, dataType, stringLength);
            } else {
                int leading = leadingPaddingBitsLittleEndian(dataType);
                if (leading > 0) {
                    writeBuffer.writeUnsignedInt(leading, 0, UNSIGNED_BINARY);
                }
                DataItem.staticSerialize(writeBuffer, value, dataType, stringLength);
                writeBuffer.writeUnsignedInt(padding - leading, 0, UNSIGNED_BINARY);
            }
            return;
        }

        for (int i = 0; i < numberOfValues; i++) {
            DataItem.staticSerialize(writeBuffer, elementAt(value, i), dataType, stringLength);
        }
        // Pad the last register when an odd number of sub-register values leaves it half filled.
        int trailing = trailingPaddingBits(dataType, numberOfValues);
        if (trailing > 0) {
            writeBuffer.writeUnsignedInt(trailing, 0, UNSIGNED_BINARY);
        }
    }

    /**
     * How many bytes {@link #serialize} will write for the given value.
     */
    public static int lengthInBytes(PlcValue value, ModbusDataType dataType, int numberOfValues,
                                    int stringLength) {
        if (numberOfValues == 1) {
            return (widthBits(dataType, stringLength) + paddingBits(dataType) + 7) / 8;
        }
        int bits = (widthBits(dataType, stringLength) * numberOfValues)
            + trailingPaddingBits(dataType, numberOfValues);
        return (bits + 7) / 8;
    }

    /**
     * The bits a lone value of this type needs on top of its own width to fill a register. Types
     * that are a whole register or wider need none, and CHAR is deliberately left unpadded, as it
     * always has been.
     */
    private static int paddingBits(ModbusDataType dataType) {
        switch (dataType) {
            case BOOL:
                return 15;
            case BYTE:
            case SINT:
            case USINT:
                return 8;
            default:
                return 0;
        }
    }

    /**
     * Where the padding sits for a lone value when the byte order is little endian. A byte simply
     * comes first and is padded behind, but a BOOL sits between seven bits and eight - which is
     * what the ',false' BOOL case used to spell out.
     */
    private static int leadingPaddingBitsLittleEndian(ModbusDataType dataType) {
        return (dataType == ModbusDataType.BOOL) ? 7 : 0;
    }

    /** The width of one value of this type, in bits. */
    private static int widthBits(ModbusDataType dataType, int stringLength) {
        switch (dataType) {
            case BOOL:
                return 1;
            case BYTE:
            case SINT:
            case USINT:
            case CHAR:
                return 8;
            case WORD:
            case INT:
            case UINT:
            case WCHAR:
                return 16;
            case DWORD:
            case DINT:
            case UDINT:
            case REAL:
                return 32;
            case LWORD:
            case LINT:
            case ULINT:
            case LREAL:
                return 64;
            case STRING:
                return stringLength * 8;
            case WSTRING:
                return stringLength * 16;
            default:
                throw new IllegalArgumentException("Unsupported Modbus data type " + dataType);
        }
    }

    /** The bits needed to round a packed run of values up to a whole register. */
    private static int trailingPaddingBits(ModbusDataType dataType, int numberOfValues) {
        int width = widthBits(dataType, 1);
        if (width >= 16) {
            return 0;
        }
        int usedBits = width * numberOfValues;
        int remainder = usedBits % 16;
        return (remainder == 0) ? 0 : 16 - remainder;
    }

    private static PlcValue elementAt(PlcValue value, int index) {
        return value.isList() ? value.getList().get(index) : value;
    }
}

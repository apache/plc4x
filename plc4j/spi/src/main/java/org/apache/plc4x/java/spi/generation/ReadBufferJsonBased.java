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
package org.apache.plc4x.java.spi.generation;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;

import java.io.IOException;
import java.io.InputStream;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.*;

@SuppressWarnings({"rawtypes", "unchecked"})
public class ReadBufferJsonBased implements ReadBuffer, BufferCommons {

    public static final String REQUIRED_ELEMENT_NOT_FOUND = "Required element %s not found in %s";
    public static final String REQUIRED_CONTEXT_NOT_FOUND = "Required context %s not found in %s";

    private final Deque stack;
    private final Object rootElement;
    private final boolean doValidateAttr;

    private int pos;

    public ReadBufferJsonBased(InputStream is) {
        this(is, true);
    }

    public ReadBufferJsonBased(InputStream is, boolean doValidateAttr) {
        this.doValidateAttr = doValidateAttr;
        pos = 1;
        stack = new ArrayDeque();
        // JsonParser here would be overkill as json is by definition not deterministic (key/value)
        ObjectMapper mapper = new ObjectMapper();
        try {
            rootElement = mapper.readValue(is, Map.class);
        } catch (IOException e) {
            throw new PlcRuntimeException(e);
        }
    }

    @Override
    public int getPos() {
        return pos / 8;
    }

    @Override
    public void reset(int pos) {
        throw new NotImplementedException();
    }

    @Override
    public boolean hasMore(int numBits) {
        return false;
    }

    @Override
    public ByteOrder getByteOrder() {
        // NO OP
        return ByteOrder.BIG_ENDIAN;
    }

    @Override
    public void setByteOrder(ByteOrder byteOrder) {
        // NO OP
    }

    @Override
    public void pullContext(String logicalName, WithReaderArgs... readerArgs) {
        logicalName = sanitizeLogicalName(logicalName);
        if (stack.isEmpty()) {
            if (!(rootElement instanceof Map)) {
                throw new PlcRuntimeException(String.format(REQUIRED_CONTEXT_NOT_FOUND, logicalName, rootElement));
            }
            Object context = ((Map) rootElement).get(logicalName);
            if (context == null) {
                throw new PlcRuntimeException(String.format(REQUIRED_CONTEXT_NOT_FOUND, logicalName, rootElement));
            }
            stack.push(context);
            return;
        }
        Object peek = stack.peek();
        if (peek instanceof List) {
            List contextList = (List) stack.pop();
            Object context = contextList.get(0);
            if (contextList.size() < 2) {
                stack.push(Collections.emptyList());
            } else {
                contextList.remove(0);
                stack.push(contextList);
            }
            Object subContext = ((Map) context).get(logicalName);
            if (subContext == null) {
                throw new PlcRuntimeException(String.format(REQUIRED_CONTEXT_NOT_FOUND, logicalName, peek));
            }
            stack.push(subContext);
            return;
        }
        if (!(peek instanceof Map)) {
            throw new PlcRuntimeException("Invalid parser state");
        }
        Map map = (Map) peek;
        Object context = map.get(logicalName);
        if (context == null) {
            throw new PlcRuntimeException(String.format(REQUIRED_CONTEXT_NOT_FOUND, logicalName, peek));
        }
        stack.push(context);
    }

    @Override
    public boolean readBit(String logicalName, WithReaderArgs... readerArgs) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        move(1);
        Map element = getElement(logicalName);
        validateAttr(logicalName, element, rwBitKey, 1);
        Boolean value = (Boolean) element.get(logicalName);
        if (value == null) {
            throw new PlcRuntimeException(String.format(REQUIRED_ELEMENT_NOT_FOUND, logicalName, stack.peek()));
        }
        return value;
    }

    @Override
    public byte readByte(String logicalName, WithReaderArgs... readerArgs) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        move(8);
        Map element = getElement(logicalName);
        validateAttr(logicalName, element, rwByteKey, 8);
        String hexString = (String) element.get(logicalName);
        if (hexString == null) {
            throw new PlcRuntimeException(String.format(REQUIRED_ELEMENT_NOT_FOUND, logicalName, stack.peek()));
        }
        if (!hexString.startsWith("0x")) {
            throw new PlcRuntimeException(String.format("Hex string should start with 0x. Actual value %s", hexString));
        }
        hexString = hexString.substring(2);
        return Byte.parseByte(hexString, 16);
    }

    @Override
    public byte[] readByteArray(String logicalName, int numberOfBytes, WithReaderArgs... readerArgs) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        move(8 * numberOfBytes);
        Map element = getElement(logicalName);
        validateAttr(logicalName, element, rwByteKey, 8 * numberOfBytes);
        String hexString = (String) element.get(logicalName);
        if (hexString == null) {
            throw new PlcRuntimeException(String.format(REQUIRED_ELEMENT_NOT_FOUND, logicalName, stack.peek()));
        }
        if (!hexString.startsWith("0x")) {
            throw new PlcRuntimeException(String.format("Hex string should start with 0x. Actual value %s", hexString));
        }
        hexString = hexString.substring(2);
        byte[] bytes = new byte[numberOfBytes];
        for (int i = 0; i < hexString.length(); i = i + 2) {
            bytes[i / 2] = Byte.parseByte(hexString.substring(i, i + 2), 16);
        }
        return bytes;
    }

    @Override
    public byte readUnsignedByte(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        move(bitLength);
        Map element = getElement(logicalName);
        validateAttr(logicalName, element, rwUintKey, bitLength);
        Integer value = (Integer) element.get(logicalName);
        if (value == null) {
            throw new PlcRuntimeException(String.format(REQUIRED_ELEMENT_NOT_FOUND, logicalName, stack.peek()));
        }
        return value.byteValue();
    }

    @Override
    public short readUnsignedShort(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        move(bitLength);
        Map element = getElement(logicalName);
        validateAttr(logicalName, element, rwUintKey, bitLength);
        Integer value = (Integer) element.get(logicalName);
        if (value == null) {
            throw new PlcRuntimeException(String.format(REQUIRED_ELEMENT_NOT_FOUND, logicalName, stack.peek()));
        }
        return value.shortValue();
    }

    @Override
    public int readUnsignedInt(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        move(bitLength);
        Map element = getElement(logicalName);
        validateAttr(logicalName, element, rwUintKey, bitLength);
        Integer value = (Integer) element.get(logicalName);
        if (value == null) {
            throw new PlcRuntimeException(String.format(REQUIRED_ELEMENT_NOT_FOUND, logicalName, stack.peek()));
        }
        return value;
    }

    @Override
    public long readUnsignedLong(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        move(bitLength);
        Map element = getElement(logicalName);
        validateAttr(logicalName, element, rwUintKey, bitLength);
        Integer value = (Integer) element.get(logicalName);
        if (value == null) {
            throw new PlcRuntimeException(String.format(REQUIRED_ELEMENT_NOT_FOUND, logicalName, stack.peek()));
        }
        return Long.valueOf(value);
    }

    @Override
    public BigInteger readUnsignedBigInteger(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        move(bitLength);
        Map element = getElement(logicalName);
        validateAttr(logicalName, element, rwUintKey, bitLength);
        Integer value = (Integer) element.get(logicalName);
        if (value == null) {
            throw new PlcRuntimeException(String.format(REQUIRED_ELEMENT_NOT_FOUND, logicalName, stack.peek()));
        }
        return BigInteger.valueOf(value);
    }

    @Override
    public byte readSignedByte(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        move(bitLength);
        Map element = getElement(logicalName);
        validateAttr(logicalName, element, rwIntKey, bitLength);
        Integer value = (Integer) element.get(logicalName);
        if (value == null) {
            throw new PlcRuntimeException(String.format(REQUIRED_ELEMENT_NOT_FOUND, logicalName, stack.peek()));
        }
        return value.byteValue();
    }

    @Override
    public short readShort(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        move(bitLength);
        Map element = getElement(logicalName);
        validateAttr(logicalName, element, rwIntKey, bitLength);
        Integer value = (Integer) element.get(logicalName);
        if (value == null) {
            throw new PlcRuntimeException(String.format(REQUIRED_ELEMENT_NOT_FOUND, logicalName, stack.peek()));
        }
        return value.shortValue();
    }

    @Override
    public int readInt(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        move(bitLength);
        Map element = getElement(logicalName);
        validateAttr(logicalName, element, rwIntKey, bitLength);
        Integer value = (Integer) element.get(logicalName);
        if (value == null) {
            throw new PlcRuntimeException(String.format(REQUIRED_ELEMENT_NOT_FOUND, logicalName, stack.peek()));
        }
        return value;
    }

    @Override
    public long readLong(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        move(bitLength);
        Map element = getElement(logicalName);
        validateAttr(logicalName, element, rwIntKey, bitLength);
        Integer value = (Integer) element.get(logicalName);
        if (value == null) {
            throw new PlcRuntimeException(String.format(REQUIRED_ELEMENT_NOT_FOUND, logicalName, stack.peek()));
        }
        return Long.valueOf(value);
    }

    @Override
    public BigInteger readBigInteger(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        move(bitLength);
        Map element = getElement(logicalName);
        validateAttr(logicalName, element, rwIntKey, bitLength);
        Integer value = (Integer) element.get(logicalName);
        if (value == null) {
            throw new PlcRuntimeException(String.format(REQUIRED_ELEMENT_NOT_FOUND, logicalName, stack.peek()));
        }
        return BigInteger.valueOf(value);
    }

    @Override
    public float readFloat(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        move(bitLength);
        Map element = getElement(logicalName);
        validateAttr(logicalName, element, rwFloatKey, bitLength);
        Float value = (Float) element.get(logicalName);
        if (value == null) {
            throw new PlcRuntimeException(String.format(REQUIRED_ELEMENT_NOT_FOUND, logicalName, stack.peek()));
        }
        return value;
    }

    @Override
    public double readDouble(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        move(bitLength);
        Map element = getElement(logicalName);
        validateAttr(logicalName, element, rwFloatKey, bitLength);
        Float value = (Float) element.get(logicalName);
        if (value == null) {
            throw new PlcRuntimeException(String.format(REQUIRED_ELEMENT_NOT_FOUND, logicalName, stack.peek()));
        }
        return Double.valueOf(value);
    }

    @Override
    public BigDecimal readBigDecimal(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        move(bitLength);
        Map element = getElement(logicalName);
        validateAttr(logicalName, element, rwFloatKey, bitLength);
        Float value = (Float) element.get(logicalName);
        if (value == null) {
            throw new PlcRuntimeException(String.format(REQUIRED_ELEMENT_NOT_FOUND, logicalName, stack.peek()));
        }
        return BigDecimal.valueOf(value);
    }

    @Override
    public String readString(String logicalName, int bitLength, String encoding, WithReaderArgs... readerArgs) {
        logicalName = sanitizeLogicalName(logicalName);
        move(bitLength);
        Map element = getElement(logicalName);
        validateAttr(logicalName, element, rwStringKey, bitLength);
        String value = (String) element.get(logicalName);
        if (value == null) {
            throw new PlcRuntimeException(String.format(REQUIRED_ELEMENT_NOT_FOUND, logicalName, stack.peek()));
        }
        return value;
    }

    @Override
    public void closeContext(String logicalName, WithReaderArgs... readerArgs) {
        logicalName = sanitizeLogicalName(logicalName);
        if (stack.isEmpty()) {
            throw new PlcRuntimeException(String.format("Required context close %s not found in %s", logicalName, rootElement));
        }
        // Delete us from stack
        stack.pop();
        if (stack.isEmpty()) {
            return;
        }
        Object peek = stack.peek();
        if (peek instanceof List) {
            return;
        }
        if (!(peek instanceof Map)) {
            throw new PlcRuntimeException("Invalid parser state");
        }
        Map map = (Map) peek;
        if (map.get(logicalName) == null) {
            throw new PlcRuntimeException(String.format(REQUIRED_CONTEXT_NOT_FOUND, logicalName, peek));
        }
        map.remove(logicalName);
    }


    private Map getElement(String logicalName) {
        logicalName = sanitizeLogicalName(logicalName);
        Object peek = stack.peek();
        Map element;
        if (peek instanceof List) {
            List elementList = (List) stack.pop();
            element = (Map) elementList.get(0);
            if (elementList.size() < 2) {
                stack.push(Collections.emptyList());
            } else {
                elementList.remove(0);
                stack.push(elementList);
            }
            return element;
        } else if (peek instanceof Map) {
            return (Map) peek;
        } else {
            throw new PlcRuntimeException(String.format("Invalid state at %s with %s", logicalName, peek));
        }
    }


    private void validateAttr(String logicalName, Map element, String dataType, int bitLength) {
        if (!doValidateAttr) {
            return;
        }
        String renderedKeyDataLengthKey = String.format("%s__plc4x_%s", logicalName, "dataType");
        String actualDataType = (String) element.get(renderedKeyDataLengthKey);
        if (!dataType.equals(actualDataType)) {
            throw new PlcRuntimeException(String.format("Unexpected %s :%s. Want %s", renderedKeyDataLengthKey, actualDataType, dataType));
        }
        String renderedBitLengthKey = String.format("%s__plc4x_%s", logicalName, "bitLength");
        Integer actualBitLength = (Integer) element.get(renderedBitLengthKey);
        if (bitLength != actualBitLength) {
            throw new PlcRuntimeException(String.format("Unexpected %s :%s. Want %s", renderedBitLengthKey, actualBitLength, bitLength));
        }
    }

    private void move(int bits) {
        pos += bits;
    }
}

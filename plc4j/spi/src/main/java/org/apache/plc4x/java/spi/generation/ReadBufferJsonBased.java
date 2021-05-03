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

package org.apache.plc4x.java.spi.generation;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;

import java.io.IOException;
import java.io.InputStream;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.*;

@SuppressWarnings({"rawtypes", "unchecked"})
public class ReadBufferJsonBased implements ReadBuffer {

    Stack stack;

    Object rootElement;

    int pos;

    boolean doValidateAttr;

    public ReadBufferJsonBased(InputStream is) {
        this(is, true);
    }

    public ReadBufferJsonBased(InputStream is, boolean doValidateAttr) {
        this.doValidateAttr = doValidateAttr;
        pos = 1;
        stack = new Stack<>();
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
    public boolean hasMore(int numBits) {
        return false;
    }

    @Override
    public void pullContext(String logicalName, WithReaderArgs... readerArgs) {
        logicalName = sanitizeLogicalName(logicalName);
        if (stack.empty()) {
            if (!(rootElement instanceof Map)) {
                throw new PlcRuntimeException(String.format("Required context %s not found in %s", logicalName, rootElement));
            }
            Object context = ((Map)rootElement).get(logicalName);
            if (context == null) {
                throw new PlcRuntimeException(String.format("Required context %s not found in %s", logicalName, rootElement));
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
                throw new PlcRuntimeException(String.format("Required context %s not found in %s", logicalName, peek));
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
            throw new PlcRuntimeException(String.format("Required context %s not found in %s", logicalName, peek));
        }
        stack.push(context);
    }

    @Override
    public boolean readBit(String logicalName, WithReaderArgs... readerArgs) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        move(1);
        Map element = getElement(logicalName);
        validateAttr(logicalName, element, "bit", 1);
        Boolean value = (Boolean) element.get(logicalName);
        if (value == null) {
            throw new PlcRuntimeException(String.format("Required element %s not found in %s", logicalName, stack.peek()));
        }
        return value;
    }

    @Override
    public byte readUnsignedByte(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        move(bitLength);
        Map element = getElement(logicalName);
        validateAttr(logicalName, element, "uint", bitLength);
        Integer value = (Integer) element.get(logicalName);
        if (value == null) {
            throw new PlcRuntimeException(String.format("Required element %s not found in %s", logicalName, stack.peek()));
        }
        return value.byteValue();
    }

    @Override
    public short readUnsignedShort(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        move(bitLength);
        Map element = getElement(logicalName);
        validateAttr(logicalName, element, "uint", bitLength);
        Integer value = (Integer) element.get(logicalName);
        if (value == null) {
            throw new PlcRuntimeException(String.format("Required element %s not found in %s", logicalName, stack.peek()));
        }
        return value.shortValue();
    }

    @Override
    public int readUnsignedInt(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        move(bitLength);
        Map element = getElement(logicalName);
        validateAttr(logicalName, element, "uint", bitLength);
        Integer value = (Integer) element.get(logicalName);
        if (value == null) {
            throw new PlcRuntimeException(String.format("Required element %s not found in %s", logicalName, stack.peek()));
        }
        return value;
    }

    @Override
    public long readUnsignedLong(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        move(bitLength);
        Map element = getElement(logicalName);
        validateAttr(logicalName, element, "uint", bitLength);
        Integer value = (Integer) element.get(logicalName);
        if (value == null) {
            throw new PlcRuntimeException(String.format("Required element %s not found in %s", logicalName, stack.peek()));
        }
        return Long.valueOf(value);
    }

    @Override
    public BigInteger readUnsignedBigInteger(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        move(bitLength);
        Map element = getElement(logicalName);
        validateAttr(logicalName, element, "uint", bitLength);
        Integer value = (Integer) element.get(logicalName);
        if (value == null) {
            throw new PlcRuntimeException(String.format("Required element %s not found in %s", logicalName, stack.peek()));
        }
        return BigInteger.valueOf(value);
    }

    @Override
    public byte readByte(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        move(bitLength);
        Map element = getElement(logicalName);
        validateAttr(logicalName, element, "int", bitLength);
        Integer value = (Integer) element.get(logicalName);
        if (value == null) {
            throw new PlcRuntimeException(String.format("Required element %s not found in %s", logicalName, stack.peek()));
        }
        return value.byteValue();
    }

    @Override
    public short readShort(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        move(bitLength);
        Map element = getElement(logicalName);
        validateAttr(logicalName, element, "int", bitLength);
        Integer value = (Integer) element.get(logicalName);
        if (value == null) {
            throw new PlcRuntimeException(String.format("Required element %s not found in %s", logicalName, stack.peek()));
        }
        return value.shortValue();
    }

    @Override
    public int readInt(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        move(bitLength);
        Map element = getElement(logicalName);
        validateAttr(logicalName, element, "int", bitLength);
        Integer value = (Integer) element.get(logicalName);
        if (value == null) {
            throw new PlcRuntimeException(String.format("Required element %s not found in %s", logicalName, stack.peek()));
        }
        return value;
    }

    @Override
    public long readLong(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        move(bitLength);
        Map element = getElement(logicalName);
        validateAttr(logicalName, element, "int", bitLength);
        Integer value = (Integer) element.get(logicalName);
        if (value == null) {
            throw new PlcRuntimeException(String.format("Required element %s not found in %s", logicalName, stack.peek()));
        }
        return Long.valueOf(value);
    }

    @Override
    public BigInteger readBigInteger(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        move(bitLength);
        Map element = getElement(logicalName);
        validateAttr(logicalName, element, "int", bitLength);
        Integer value = (Integer) element.get(logicalName);
        if (value == null) {
            throw new PlcRuntimeException(String.format("Required element %s not found in %s", logicalName, stack.peek()));
        }
        return BigInteger.valueOf(value);
    }

    @Override
    public float readFloat(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        move(bitLength);
        Map element = getElement(logicalName);
        validateAttr(logicalName, element, "float", bitLength);
        Float value = (Float) element.get(logicalName);
        if (value == null) {
            throw new PlcRuntimeException(String.format("Required element %s not found in %s", logicalName, stack.peek()));
        }
        return value;
    }

    @Override
    public double readDouble(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        move(bitLength);
        Map element = getElement(logicalName);
        validateAttr(logicalName, element, "float", bitLength);
        Float value = (Float) element.get(logicalName);
        if (value == null) {
            throw new PlcRuntimeException(String.format("Required element %s not found in %s", logicalName, stack.peek()));
        }
        return Double.valueOf(value);
    }

    @Override
    public BigDecimal readBigDecimal(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        move(bitLength);
        Map element = getElement(logicalName);
        validateAttr(logicalName, element, "float", bitLength);
        Float value = (Float) element.get(logicalName);
        if (value == null) {
            throw new PlcRuntimeException(String.format("Required element %s not found in %s", logicalName, stack.peek()));
        }
        return BigDecimal.valueOf(value);
    }

    @Override
    public String readString(String logicalName, int bitLength, String encoding, WithReaderArgs... readerArgs) {
        logicalName = sanitizeLogicalName(logicalName);
        move(bitLength);
        Map element = getElement(logicalName);
        validateAttr(logicalName, element, "int", bitLength);
        String value = (String) element.get(logicalName);
        if (value == null) {
            throw new PlcRuntimeException(String.format("Required element %s not found in %s", logicalName, stack.peek()));
        }
        return value;
    }

    @Override
    public void closeContext(String logicalName, WithReaderArgs... readerArgs) {
        logicalName = sanitizeLogicalName(logicalName);
        if (stack.empty()) {
            throw new PlcRuntimeException(String.format("Required context close %s not found in %s", logicalName, rootElement));
        }
        // Delete us from stack
        stack.pop();
        if (stack.empty()) {
            return;
        }
        Object peek =stack.peek();
        if (peek instanceof List) {
            return;
        }
        if (!(peek instanceof Map)) {
            throw new PlcRuntimeException("Invalid parser state");
        }
        Map map = (Map) peek;
        if (map.get(logicalName)==null) {
            throw new PlcRuntimeException(String.format("Required context %s not found in %s", logicalName, peek));
        }
        map.remove(logicalName);
    }


    private Map getElement(String logicalName) {
        logicalName = sanitizeLogicalName(logicalName);
        Object peek = stack.peek();
        Map element;
        if (peek instanceof List) {
            List elementList = (List)stack.pop();
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
        if (bitLength !=actualBitLength) {
            throw new PlcRuntimeException(String.format("Unexpected %s :%s. Want %s", renderedBitLengthKey, actualBitLength, bitLength));
        }
    }


    private String sanitizeLogicalName(String logicalName) {
        if (logicalName.equals("")) {
            return "value";
        }
        return logicalName;
    }


    private void move(int bits) {
        pos += bits;
    }
}

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

import org.apache.commons.lang3.StringUtils;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;

import javax.xml.stream.XMLEventFactory;
import javax.xml.stream.XMLEventWriter;
import javax.xml.stream.XMLOutputFactory;
import javax.xml.stream.XMLStreamException;
import javax.xml.stream.events.*;

import java.io.ByteArrayOutputStream;
import java.io.UnsupportedEncodingException;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.Stack;

public class WriteBufferXmlBased implements WriteBuffer {

    Stack<String> stack;

    ByteArrayOutputStream byteArrayOutputStream;

    XMLEventFactory xmlEventFactory;

    XMLEventWriter xmlEventWriter;

    int pos = 1;

    int depth = 0;

    public WriteBufferXmlBased() {
        byteArrayOutputStream = new ByteArrayOutputStream();
        XMLOutputFactory xmlOutputFactory = XMLOutputFactory.newInstance();
        xmlEventFactory = XMLEventFactory.newInstance();
        try {
            xmlEventWriter = xmlOutputFactory.createXMLEventWriter(byteArrayOutputStream);
        } catch (XMLStreamException e) {
            throw new PlcRuntimeException(e);
        }
        this.stack = new Stack<>();
    }

    @Override
    public int getPos() {
        return pos / 8;
    }

    @Override
    public void pushContext(String logicalName) {
        try {
            indent();
            depth++;
            StartElement startElement = xmlEventFactory.createStartElement("", "", logicalName);
            xmlEventWriter.add(startElement);
            newLine();
        } catch (XMLStreamException e) {
            throw new PlcRuntimeException(e);
        }
        stack.push(logicalName);
    }

    @Override
    public void writeBit(String logicalName, boolean value) throws ParseException {
        String dataType = "bit";
        int bitLength = 1;
        String data = Boolean.valueOf(value).toString();
        createAndAppend(logicalName, dataType, bitLength, data);
        move(1);
    }

    @Override
    public void writeUnsignedByte(String logicalName, int bitLength, byte value) throws ParseException {
        createAndAppend(logicalName, "uint", bitLength, Byte.valueOf(value).toString());
        move(bitLength);
    }

    @Override
    public void writeUnsignedShort(String logicalName, int bitLength, short value) throws ParseException {
        createAndAppend(logicalName, "uint", bitLength, Short.valueOf(value).toString());
        move(bitLength);
    }

    @Override
    public void writeUnsignedInt(String logicalName, int bitLength, int value) throws ParseException {
        createAndAppend(logicalName, "uint", bitLength, Integer.valueOf(value).toString());
        move(bitLength);
    }

    @Override
    public void writeUnsignedLong(String logicalName, int bitLength, long value) throws ParseException {
        createAndAppend(logicalName, "uint", bitLength, Long.valueOf(value).toString());
        move(bitLength);
    }

    @Override
    public void writeUnsignedBigInteger(String logicalName, int bitLength, BigInteger value) throws ParseException {
        createAndAppend(logicalName, "uint", bitLength, value.toString());
        move(bitLength);
    }

    @Override
    public void writeByte(String logicalName, int bitLength, byte value) throws ParseException {
        createAndAppend(logicalName, "int", bitLength, Byte.valueOf(value).toString());
        move(bitLength);
    }

    @Override
    public void writeShort(String logicalName, int bitLength, short value) throws ParseException {
        createAndAppend(logicalName, "int", bitLength, Short.valueOf(value).toString());
        move(bitLength);
    }

    @Override
    public void writeInt(String logicalName, int bitLength, int value) throws ParseException {
        createAndAppend(logicalName, "int", bitLength, Integer.valueOf(value).toString());
        move(bitLength);
    }

    @Override
    public void writeLong(String logicalName, int bitLength, long value) throws ParseException {
        createAndAppend(logicalName, "int", bitLength, Long.valueOf(value).toString());
        move(bitLength);
    }

    @Override
    public void writeBigInteger(String logicalName, int bitLength, BigInteger value) throws ParseException {
        createAndAppend(logicalName, "int", bitLength, value.toString());
        move(bitLength);
    }

    @Override
    public void writeFloat(String logicalName, float value, int bitsExponent, int bitsMantissa) throws ParseException {
        int bitLength = (value < 0 ? 1 : 0) + bitsExponent + bitsMantissa;
        createAndAppend(logicalName, "float", bitLength, Float.valueOf(value).toString());
        move(bitLength);
    }

    @Override
    public void writeDouble(String logicalName, double value, int bitsExponent, int bitsMantissa) throws ParseException {
        int bitLength = (value < 0 ? 1 : 0) + bitsExponent + bitsMantissa;
        createAndAppend(logicalName, "float", bitLength, Double.valueOf(value).toString());
        move(bitLength);
    }

    @Override
    public void writeBigDecimal(String logicalName, int bitLength, BigDecimal value) throws ParseException {
        createAndAppend(logicalName, "float", bitLength, value.toString());
        move(bitLength);
    }

    @Override
    public void writeString(String logicalName, int bitLength, String encoding, String value) throws ParseException {
        createAndAppend(logicalName, "string", bitLength, value);
        move(bitLength);
    }

    @Override
    public void popContext(String logicalName) {
        try {
            depth--;
            indent();
            EndElement endElement = xmlEventFactory.createEndElement("", "", logicalName);
            xmlEventWriter.add(endElement);
            newLine();
        } catch (XMLStreamException e) {
            throw new PlcRuntimeException(e);
        }

        String context = stack.pop();
        if (!context.equals(logicalName)) {
            throw new PlcRuntimeException("Unexpected pop context '" + context + '\'');
        }
        if (stack.isEmpty()) {
            try {
                xmlEventWriter.close();
            } catch (XMLStreamException e) {
                throw new PlcRuntimeException(e);
            }
        }
    }

    public String getXmlString() {
        try {
            return byteArrayOutputStream.toString("UTF-8");
        } catch (UnsupportedEncodingException e) {
            throw new PlcRuntimeException(e);
        }
    }

    private void move(int bits) {
        pos += bits;
    }

    private void newLine() throws XMLStreamException {
        Characters newLine = xmlEventFactory.createCharacters("\n");
        xmlEventWriter.add(newLine);
    }

    private void indent() throws XMLStreamException {
        Characters indent = xmlEventFactory.createCharacters(StringUtils.repeat("  ", depth));
        xmlEventWriter.add(indent);
    }

    private void createAndAppend(String logicalName, String dataType, int bitLength, String data) {
        try {
            indent();
            StartElement startElement = xmlEventFactory.createStartElement("", "", sanitizeLogicalName(logicalName));
            xmlEventWriter.add(startElement);
            Attribute dataTypeAttribute = xmlEventFactory.createAttribute("dataType", dataType);
            xmlEventWriter.add(dataTypeAttribute);
            Attribute bitLengthAttribute = xmlEventFactory.createAttribute("bitLength", String.valueOf(bitLength));
            xmlEventWriter.add(bitLengthAttribute);
            Characters dataCharacters = xmlEventFactory.createCharacters(data);
            xmlEventWriter.add(dataCharacters);
            EndElement endElement = xmlEventFactory.createEndElement("", "", sanitizeLogicalName(logicalName));
            xmlEventWriter.add(endElement);
            newLine();
        } catch (XMLStreamException e) {
            throw new PlcRuntimeException(e);
        }
    }

    private String sanitizeLogicalName(String logicalName) {
        if (logicalName.equals("")) {
            return "value";
        }
        return logicalName;
    }
}

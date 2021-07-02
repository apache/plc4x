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
import javax.xml.stream.events.Attribute;
import javax.xml.stream.events.Characters;
import javax.xml.stream.events.EndElement;
import javax.xml.stream.events.StartElement;
import java.io.ByteArrayOutputStream;
import java.io.UnsupportedEncodingException;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.ArrayDeque;
import java.util.Deque;
import java.util.Stack;

public class WriteBufferXmlBased implements WriteBuffer, BufferCommons {

    Deque<String> stack;

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
        this.stack = new ArrayDeque<>();
    }

    @Override
    public int getPos() {
        return pos / 8;
    }

    @Override
    public void pushContext(String logicalName, WithWriterArgs... writerArgs) {
        try {
            indent();
            depth++;
            StartElement startElement = xmlEventFactory.createStartElement("", "", logicalName);
            xmlEventWriter.add(startElement);
            if (isToBeRenderedAsList(writerArgs)) {
                Attribute isListAttribute = xmlEventFactory.createAttribute(rwIsListKey, "true");
                xmlEventWriter.add(isListAttribute);
            }
            newLine();
        } catch (XMLStreamException e) {
            throw new PlcRuntimeException(e);
        }
        stack.push(logicalName);
    }

    @Override
    public void writeBit(String logicalName, boolean value, WithWriterArgs... writerArgs) throws ParseException {
        String dataType = "bit";
        int bitLength = 1;
        String data = Boolean.toString(value);
        createAndAppend(logicalName, dataType, bitLength, data, writerArgs);
        move(1);
    }

    @Override
    public void writeByte(String logicalName, byte value, WithWriterArgs... writerArgs) throws ParseException {
        createAndAppend(logicalName, rwByteKey, 8, String.format("0x%02x", value), writerArgs);
        move(8);
    }

    @Override
    public void writeByteArray(String logicalName, byte[] bytes, WithWriterArgs... writerArgs) throws ParseException {
        StringBuilder hexString = new StringBuilder("0x");
        for (byte aByte : bytes) {
            hexString.append(String.format("%02x", aByte));
        }
        createAndAppend(logicalName, rwByteKey, bytes.length * 8, hexString.toString(), writerArgs);
        move(8 * bytes.length);
    }

    @Override
    public void writeUnsignedByte(String logicalName, int bitLength, byte value, WithWriterArgs... writerArgs) throws ParseException {
        createAndAppend(logicalName, rwUintKey, bitLength, Byte.toString(value), writerArgs);
        move(bitLength);
    }

    @Override
    public void writeUnsignedShort(String logicalName, int bitLength, short value, WithWriterArgs... writerArgs) throws ParseException {
        createAndAppend(logicalName, rwUintKey, bitLength, Short.toString(value), writerArgs);
        move(bitLength);
    }

    @Override
    public void writeUnsignedInt(String logicalName, int bitLength, int value, WithWriterArgs... writerArgs) throws ParseException {
        createAndAppend(logicalName, rwUintKey, bitLength, Integer.toString(value), writerArgs);
        move(bitLength);
    }

    @Override
    public void writeUnsignedLong(String logicalName, int bitLength, long value, WithWriterArgs... writerArgs) throws ParseException {
        createAndAppend(logicalName, rwUintKey, bitLength, Long.toString(value), writerArgs);
        move(bitLength);
    }

    @Override
    public void writeUnsignedBigInteger(String logicalName, int bitLength, BigInteger value, WithWriterArgs... writerArgs) throws ParseException {
        createAndAppend(logicalName, rwUintKey, bitLength, value.toString(), writerArgs);
        move(bitLength);
    }

    @Override
    public void writeSignedByte(String logicalName, int bitLength, byte value, WithWriterArgs... writerArgs) throws ParseException {
        createAndAppend(logicalName, rwIntKey, bitLength, Byte.toString(value), writerArgs);
        move(bitLength);
    }

    @Override
    public void writeShort(String logicalName, int bitLength, short value, WithWriterArgs... writerArgs) throws ParseException {
        createAndAppend(logicalName, rwIntKey, bitLength, Short.toString(value), writerArgs);
        move(bitLength);
    }

    @Override
    public void writeInt(String logicalName, int bitLength, int value, WithWriterArgs... writerArgs) throws ParseException {
        createAndAppend(logicalName, rwIntKey, bitLength, Integer.toString(value), writerArgs);
        move(bitLength);
    }

    @Override
    public void writeLong(String logicalName, int bitLength, long value, WithWriterArgs... writerArgs) throws ParseException {
        createAndAppend(logicalName, rwIntKey, bitLength, Long.toString(value), writerArgs);
        move(bitLength);
    }

    @Override
    public void writeBigInteger(String logicalName, int bitLength, BigInteger value, WithWriterArgs... writerArgs) throws ParseException {
        createAndAppend(logicalName, rwIntKey, bitLength, value.toString(), writerArgs);
        move(bitLength);
    }

    @Override
    public void writeFloat(String logicalName, float value, int bitsExponent, int bitsMantissa, WithWriterArgs... writerArgs) throws ParseException {
        int bitLength = (value < 0 ? 1 : 0) + bitsExponent + bitsMantissa;
        createAndAppend(logicalName, rwFloatKey, bitLength, Float.toString(value), writerArgs);
        move(bitLength);
    }

    @Override
    public void writeDouble(String logicalName, double value, int bitsExponent, int bitsMantissa, WithWriterArgs... writerArgs) throws ParseException {
        int bitLength = (value < 0 ? 1 : 0) + bitsExponent + bitsMantissa;
        createAndAppend(logicalName, rwFloatKey, bitLength, Double.toString(value), writerArgs);
        move(bitLength);
    }

    @Override
    public void writeBigDecimal(String logicalName, int bitLength, BigDecimal value, WithWriterArgs... writerArgs) throws ParseException {
        createAndAppend(logicalName, rwFloatKey, bitLength, value.toString(), writerArgs);
        move(bitLength);
    }

    @Override
    public void writeString(String logicalName, int bitLength, String encoding, String value, WithWriterArgs... writerArgs) throws ParseException {
        createAndAppend(logicalName, rwStringKey, bitLength, value, writerArgs);
        move(bitLength);
    }

    @Override
    public void popContext(String logicalName, WithWriterArgs... writerArgs) {
        try {
            depth--;
            indent();
            EndElement endElement = xmlEventFactory.createEndElement("", "", logicalName);
            xmlEventWriter.add(endElement);
            if (depth != 0) {
                // We don't want an extra newline at the end so we write only if we are not at the end
                newLine();
            }
        } catch (XMLStreamException e) {
            throw new PlcRuntimeException(e);
        }

        String context = stack.pop();
        if (!context.equals(logicalName)) {
            throw new PlcRuntimeException("Unexpected pop context '" + context + '\''+ ". Expected '"+logicalName+ '\'');
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

    private void createAndAppend(String logicalName, String dataType, int bitLength, String data, WithWriterArgs... writerArgs) {
        try {
            indent();
            StartElement startElement = xmlEventFactory.createStartElement("", "", sanitizeLogicalName(logicalName));
            xmlEventWriter.add(startElement);
            Attribute dataTypeAttribute = xmlEventFactory.createAttribute(rwDataTypeKey, dataType);
            xmlEventWriter.add(dataTypeAttribute);
            Attribute bitLengthAttribute = xmlEventFactory.createAttribute(rwBitLengthKey, String.valueOf(bitLength));
            xmlEventWriter.add(bitLengthAttribute);
            String additionalStringRepresentation = extractAdditionalStringRepresentation(writerArgs);
            if (additionalStringRepresentation != null) {
                Attribute additionalStringRepresentationAttribute = xmlEventFactory.createAttribute(rwStringRepresentationKey, additionalStringRepresentation);
                xmlEventWriter.add(additionalStringRepresentationAttribute);
            }
            Characters dataCharacters = xmlEventFactory.createCharacters(data);
            xmlEventWriter.add(dataCharacters);
            EndElement endElement = xmlEventFactory.createEndElement("", "", sanitizeLogicalName(logicalName));
            xmlEventWriter.add(endElement);
            newLine();
        } catch (XMLStreamException e) {
            throw new PlcRuntimeException(e);
        }
    }

}

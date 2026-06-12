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
package org.apache.plc4x.java.spi.buffers.xmlbased;

import org.apache.plc4x.java.spi.buffers.api.AbstractBuffer;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;

import javax.xml.stream.XMLEventFactory;
import javax.xml.stream.XMLEventWriter;
import javax.xml.stream.XMLOutputFactory;
import javax.xml.stream.XMLStreamException;
import javax.xml.stream.events.Attribute;
import javax.xml.stream.events.Characters;
import javax.xml.stream.events.EndElement;
import javax.xml.stream.events.StartElement;
import java.io.ByteArrayOutputStream;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.util.Optional;

public class WriteBufferXmlBased extends AbstractBuffer implements WriteBuffer, XmlBuffer {

    private final ByteArrayOutputStream byteArrayOutputStream;

    private final XMLEventFactory xmlEventFactory;

    private final XMLEventWriter xmlEventWriter;

    private final int sizeInBits;
    private int positionInBits;

    public WriteBufferXmlBased() {
        byteArrayOutputStream = new ByteArrayOutputStream();
        XMLOutputFactory xmlOutputFactory = XMLOutputFactory.newInstance();
        xmlEventFactory = XMLEventFactory.newInstance();
        // TODO: Set this ....
        sizeInBits = 0;
        positionInBits = 0;
        try {
            xmlEventWriter = xmlOutputFactory.createXMLEventWriter(byteArrayOutputStream);
        } catch (XMLStreamException e) {
            throw new IllegalArgumentException(e);
        }
    }

    @Override
    public void pushContext(WithOption... options) throws BufferException {
        String name = getName(WithOption.AddOptions(getContext(), options));
        try {
            indent();
            StartElement startElement = xmlEventFactory.createStartElement("", "", name);
            xmlEventWriter.add(startElement);
            if (isToBeRenderedAsList(options)) {
                Attribute isListAttribute = xmlEventFactory.createAttribute(rwIsListKey, "true");
                xmlEventWriter.add(isListAttribute);
                // For simple elements, make sure the elements are called "value".
                options = WithOption.UpdateOptions(options, WithOption.WithName("value"));
            }
            newLine();
        } catch (XMLStreamException e) {
            throw new BufferException("Error generating xml", e);
        }
        super.pushContext(options);
    }

    @Override
    public void popContext(WithOption... options) throws BufferException {
        String name = getName(WithOption.AddOptions(getContext(), options));
        try {
            indent();
            EndElement endElement = xmlEventFactory.createEndElement("", "", name);
            xmlEventWriter.add(endElement);
            newLine();
        } catch (XMLStreamException e) {
            throw new BufferException("Error generating xml", e);
        }
        super.popContext(options);

        if (context.isEmpty()) {
            try {
                xmlEventWriter.close();
            } catch (XMLStreamException e) {
                throw new BufferException("Error generating xml", e);
            }
        }
    }

    @Override
    public void writeBit(boolean value, WithOption... options) throws BufferException {
        String name = getName(options);
        String dataType = "bit";
        int bitLength = 1;
        String data = Boolean.toString(value);
        createAndAppend(name, dataType, bitLength, data, options);
        move(1);
    }

    @Override
    public void writeBits(int numBits, byte[] value, WithOption... options) throws BufferException {
        String name = getName(options);
        StringBuilder hexString = new StringBuilder("0x");
        int numberOfBytes = (numBits + 7) / 8;
        if (numberOfBytes > value.length) {
            throw new BufferException("Not enough bytes to write " + numBits + " bits");
        }
        for (int i = 0; i < numberOfBytes; i++) {
            hexString.append(String.format("%02x", value[i]));
        }
        createAndAppend(name, rwByteKey, value.length * 8, hexString.toString(), options);
        move(8 * value.length);
    }

    @Override
    public void writeUnsignedByte(int bitLength, byte value, WithOption... options) throws BufferException {
        String name = getName(options);
        createAndAppend(name, rwUintKey, bitLength, Byte.toString(value), options);
        move(bitLength);
    }

    @Override
    public void writeUnsignedShort(int bitLength, short value, WithOption... options) throws BufferException {
        String name = getName(options);
        createAndAppend(name, rwUintKey, bitLength, Short.toString(value), options);
        move(bitLength);
    }

    @Override
    public void writeUnsignedInt(int bitLength, int value, WithOption... options) throws BufferException {
        String name = getName(options);
        createAndAppend(name, rwUintKey, bitLength, Integer.toString(value), options);
        move(bitLength);
    }

    @Override
    public void writeUnsignedLong(int bitLength, long value, WithOption... options) throws BufferException {
        String name = getName(options);
        createAndAppend(name, rwUintKey, bitLength, Long.toString(value), options);
        move(bitLength);
    }

    @Override
    public void writeUnsignedBigInteger(int bitLength, BigInteger value, WithOption... options) throws BufferException {
        String name = getName(options);
        createAndAppend(name, rwUintKey, bitLength, value.toString(), options);
        move(bitLength);
    }

    @Override
    public void writeSignedByte(int bitLength, byte value, WithOption... options) throws BufferException {
        String name = getName(options);
        createAndAppend(name, rwIntKey, bitLength, Byte.toString(value), options);
        move(bitLength);
    }

    @Override
    public void writeSignedShort(int bitLength, short value, WithOption... options) throws BufferException {
        String name = getName(options);
        createAndAppend(name, rwIntKey, bitLength, Short.toString(value), options);
        move(bitLength);
    }

    @Override
    public void writeSignedInt(int bitLength, int value, WithOption... options) throws BufferException {
        String name = getName(options);
        createAndAppend(name, rwIntKey, bitLength, Integer.toString(value), options);
        move(bitLength);
    }

    @Override
    public void writeSignedLong(int bitLength, long value, WithOption... options) throws BufferException {
        String name = getName(options);
        createAndAppend(name, rwIntKey, bitLength, Long.toString(value), options);
        move(bitLength);
    }

    @Override
    public void writeSignedBigInteger(int bitLength, BigInteger value, WithOption... options) throws BufferException {
        String name = getName(options);
        createAndAppend(name, rwIntKey, bitLength, value.toString(), options);
        move(bitLength);
    }

    @Override
    public void writeFloat(int bitLength, float value, WithOption... options) throws BufferException {
        String name = getName(options);
        createAndAppend(name, rwFloatKey, bitLength, Float.toString(value), options);
        move(bitLength);
    }

    @Override
    public void writeDouble(int bitLength, double value, WithOption... options) throws BufferException {
        String name = getName(options);
        createAndAppend(name, rwFloatKey, bitLength, Double.toString(value), options);
        move(bitLength);
    }

    @Override
    public void writeBigDecimal(int bitLength, BigDecimal value, WithOption... options) throws BufferException {
        String name = getName(options);
        createAndAppend(name, rwFloatKey, bitLength, value.toString(), options);
        move(bitLength);
    }

    @Override
    public void writeString(int bitLength, String value, WithOption... options) throws BufferException {
        String name = getName(options);
        // Try getting the concrete encoding first.
        String encoding = WithOption.extractStringEncoding(options, getContext()).orElse(null);
        // Then try with the global encoding.
        if (encoding == null) {
            encoding = WithOption.extractEncoding(options, getContext()).orElse("UTF8");
        }
        String cleanedUpString = (value != null ? value : "").trim().replaceAll("[^\t\r\n -\uD7FF\uE000-\uFFFD\ud800\udc00-\udbff\udfff]", "");
        createAndAppend(name, rwStringKey, bitLength, cleanedUpString, encoding, options);
        move(bitLength);
    }

    @Override
    public WriteBuffer createSubBuffer(int numBits, WithOption... options) {
        return null;
    }

    @Override
    public int getPositionInBits() {
        return positionInBits;
    }

    @Override
    public int getRemainingBits() {
        return sizeInBits - positionInBits;
    }

    private void move(int bits) {
        positionInBits += bits;
    }

    @Override
    public byte[] getBytes() {
        try {
            xmlEventWriter.flush();
        } catch (XMLStreamException e) {
            throw new IllegalArgumentException(e);
        }
        return byteArrayOutputStream.toByteArray();
    }

    public String getXmlString() {
        try {
            xmlEventWriter.flush();
        } catch (XMLStreamException e) {
            throw new IllegalArgumentException(e);
        }
        return byteArrayOutputStream.toString(StandardCharsets.UTF_8);
    }

    private void newLine() throws XMLStreamException {
        Characters newLine = xmlEventFactory.createCharacters("\n");
        xmlEventWriter.add(newLine);
    }

    private void indent() throws XMLStreamException {
        for (int i = 1; i < context.size(); i++) {
            xmlEventWriter.add(xmlEventFactory.createCharacters("  "));
        }
    }

    private void createAndAppend(String name, String dataType, int bitLength, String data, WithOption... options) throws BufferException {
        createAndAppend(name, dataType, bitLength, data, null, options);
    }

    private void createAndAppend(String name, String dataType, int bitLength, String data, String encoding, WithOption... options) throws BufferException {
        try {
            indent();
            StartElement startElement = xmlEventFactory.createStartElement("", "", sanitizeLogicalName(name));
            xmlEventWriter.add(startElement);
            Attribute dataTypeAttribute = xmlEventFactory.createAttribute(rwDataTypeKey, dataType);
            xmlEventWriter.add(dataTypeAttribute);
            Attribute bitLengthAttribute = xmlEventFactory.createAttribute(rwBitLengthKey, String.valueOf(bitLength));
            xmlEventWriter.add(bitLengthAttribute);
            Optional<String> additionalStringRepresentation = WithOption.extractAdditionalStringRepresentation(options);
            if (additionalStringRepresentation.isPresent()) {
                Attribute additionalStringRepresentationAttribute = xmlEventFactory.createAttribute(rwStringRepresentationKey, additionalStringRepresentation.get());
                xmlEventWriter.add(additionalStringRepresentationAttribute);
            }
            if (encoding != null) {
                Attribute encodingAttribute = xmlEventFactory.createAttribute(rwEncodingKey, encoding);
                xmlEventWriter.add(encodingAttribute);
            }
            Characters dataCharacters = xmlEventFactory.createCharacters(data);
            xmlEventWriter.add(dataCharacters);
            EndElement endElement = xmlEventFactory.createEndElement("", "", sanitizeLogicalName(name));
            xmlEventWriter.add(endElement);
            newLine();
        } catch (XMLStreamException e) {
            throw new BufferException("Error generating xml", e);
        }
    }

    String sanitizeLogicalName(String logicalName) {
        if (logicalName == null || logicalName.isBlank()) {
            return "value";
        }
        return logicalName;
    }

    boolean isToBeRenderedAsList(WithOption... options) {
        Optional<Boolean> renderAsList = WithOption.extractRenderAsList(options);
        return renderAsList.orElse(false);
    }

    protected String getName(WithOption... options) throws BufferException {
        Optional<String> name = WithOption.extractName(options, getContext());
        if (name.isEmpty()) {
            throw new BufferException("Missing 'name' option.");
        }
        return name.get();
    }

}

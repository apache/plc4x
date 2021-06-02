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

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;

import javax.xml.XMLConstants;
import javax.xml.stream.XMLEventReader;
import javax.xml.stream.XMLInputFactory;
import javax.xml.stream.XMLStreamException;
import javax.xml.stream.events.*;
import java.io.InputStream;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.Iterator;

public class ReadBufferXmlBased implements ReadBuffer, BufferCommons {

    XMLEventReader reader;

    int pos = 1;

    public ReadBufferXmlBased(InputStream is) {
        XMLInputFactory xmlInputFactory = XMLInputFactory.newInstance();
        xmlInputFactory.setProperty(XMLConstants.ACCESS_EXTERNAL_DTD, "");
        xmlInputFactory.setProperty(XMLConstants.ACCESS_EXTERNAL_SCHEMA, "");
        try {
            reader = xmlInputFactory.createXMLEventReader(is);
        } catch (XMLStreamException e) {
            throw new PlcRuntimeException(e);
        }
    }

    @Override
    public int getPos() {
        return pos / 8;
    }

    @Override
    public boolean hasMore(int numBits) {
        return true;
    }

    @Override
    public void pullContext(String logicalName, WithReaderArgs... readerArgs) {
        StartElement startElement = travelToNextStartElement();
        String elementName = startElement.getName().getLocalPart();
        if (!elementName.equals(logicalName)) {
            throw new PlcRuntimeException(String.format("Unexpected Start element '%s'. Expected '%s'", elementName, logicalName));
        }
    }

    @Override
    public boolean readBit(String logicalName, WithReaderArgs... readerArgs) throws ParseException {
        String bit = decode(logicalName, rwBitKey, 1);
        move(1);
        return bit.equals("true");
    }

    @Override
    public byte readByte(String logicalName, WithReaderArgs... readerArgs) throws ParseException {
        move(8);
        String hexString = decode(logicalName, rwByteKey, 8);
        if (!hexString.startsWith("0x")) {
            throw new PlcRuntimeException(String.format("Hex string should start with 0x. Actual value %s", hexString));
        }
        hexString = hexString.substring(2);
        return Byte.parseByte(hexString, 16);
    }

    @Override
    public byte[] readByteArray(String logicalName, int numberOfBytes, WithReaderArgs... readerArgs) throws ParseException {
        move(8 * numberOfBytes);
        String hexString = decode(logicalName, rwByteKey, 8 * numberOfBytes);
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
        move(bitLength);
        return Byte.parseByte(decode(logicalName, rwUintKey, bitLength));
    }

    @Override
    public short readUnsignedShort(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        move(bitLength);
        return Short.parseShort(decode(logicalName, rwUintKey, bitLength));
    }

    @Override
    public int readUnsignedInt(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        move(bitLength);
        return Integer.parseInt(decode(logicalName, rwUintKey, bitLength));
    }

    @Override
    public long readUnsignedLong(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        move(bitLength);
        return Long.parseLong(decode(logicalName, rwUintKey, bitLength));
    }

    @Override
    public BigInteger readUnsignedBigInteger(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        move(bitLength);
        throw new PlcRuntimeException("not implemented yet");
    }

    @Override
    public byte readSignedByte(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        move(bitLength);
        return Byte.parseByte(decode(logicalName, rwIntKey, bitLength));
    }

    @Override
    public short readShort(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        move(bitLength);
        return Short.parseShort(decode(logicalName, rwIntKey, bitLength));
    }

    @Override
    public int readInt(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        move(bitLength);
        return Integer.parseInt(decode(logicalName, rwIntKey, bitLength));
    }

    @Override
    public long readLong(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        move(bitLength);
        return Long.parseLong(decode(logicalName, rwIntKey, bitLength));
    }

    @Override
    public BigInteger readBigInteger(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        move(bitLength);
        return new BigInteger(decode(logicalName, rwIntKey, bitLength));
    }

    @Override
    public float readFloat(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        move(bitLength);
        return Float.parseFloat(decode(logicalName, rwFloatKey, bitLength));
    }

    @Override
    public double readDouble(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        move(bitLength);
        return Double.parseDouble(decode(logicalName, rwFloatKey, bitLength));
    }

    @Override
    public BigDecimal readBigDecimal(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        move(bitLength);
        return new BigDecimal(decode(logicalName, rwFloatKey, bitLength));
    }

    @Override
    public String readString(String logicalName, int bitLength, String encoding, WithReaderArgs... readerArgs) {
        move(bitLength);
        return decode(logicalName, rwStringKey, bitLength);
    }

    @Override
    public void closeContext(String logicalName, WithReaderArgs... readerArgs) {
        EndElement endElement = travelToNextEndElement();
        if (!endElement.getName().getLocalPart().equals(logicalName)) {
            throw new PlcRuntimeException(String.format("Unexpected End element '%s'. Expected '%s'", endElement.getName().getLocalPart(), logicalName));
        }
    }

    private void move(int bits) {
        pos += bits;
    }

    private StartElement travelToNextStartElement() {
        while (reader.hasNext()) {
            XMLEvent xmlEvent;
            try {
                xmlEvent = reader.nextEvent();
            } catch (XMLStreamException e) {
                throw new PlcRuntimeException(e);
            }
            if (xmlEvent.isStartElement()) {
                return xmlEvent.asStartElement();
            } else if (xmlEvent.isEndElement()) {
                throw new PlcRuntimeException(String.format("Unexpected End element %s", xmlEvent.asEndElement().getName().getLocalPart()));
            }
        }
        throw new PlcRuntimeException("EOF");
    }

    private EndElement travelToNextEndElement() {
        while (reader.hasNext()) {
            XMLEvent xmlEvent;
            try {
                xmlEvent = reader.nextEvent();
            } catch (XMLStreamException e) {
                throw new PlcRuntimeException(e);
            }
            if (xmlEvent.isStartElement()) {
                throw new PlcRuntimeException(String.format("Unexpected Start element %s", xmlEvent.asStartElement().getName().getLocalPart()));
            } else if (xmlEvent.isEndElement()) {
                return xmlEvent.asEndElement();
            }
        }
        throw new PlcRuntimeException("EOF");
    }

    private String decode(String logicalName, String dataType, int bitLength) {
        StartElement startElement = travelToNextStartElement();
        validateStartElement(startElement, logicalName, dataType, bitLength);
        Characters characters;
        try {
            characters = reader.nextEvent().asCharacters();
        } catch (XMLStreamException e) {
            throw new PlcRuntimeException(e);
        }
        String data = characters.getData();
        try {
            reader.nextEvent().asEndElement();
        } catch (XMLStreamException e) {
            throw new PlcRuntimeException(e);
        }
        return data.trim();
    }

    private void validateStartElement(StartElement startElement, String logicalName, String dataType, int bitLength) {
        logicalName = sanitizeLogicalName(logicalName);
        if (!startElement.getName().getLocalPart().equals(logicalName)) {
            throw new PlcRuntimeException(String.format("unexpected element '%s'. Expected '%s'", startElement.getName().getLocalPart(), logicalName));
        }
        validateAttr(logicalName, startElement.getAttributes(), dataType, bitLength);
    }

    private void validateAttr(String logicalName, Iterator<Attribute> attr, String dataType, int bitLength) {
        boolean dataTypeValidated = false;
        boolean bitLengthValidate = false;
        while (attr.hasNext()) {
            Attribute attribute = attr.next();
            if (attribute.getName().getLocalPart().equals(rwDataTypeKey)) {
                if (!attribute.getValue().equals(dataType)) {
                    throw new PlcRuntimeException(String.format("%s: Unexpected dataType :%s. Want %s", logicalName, attribute.getValue(), dataType));
                }
                dataTypeValidated = true;
            } else if (attribute.getName().getLocalPart().equals(rwBitLengthKey)) {
                if (!attribute.getValue().equals(Integer.toString(bitLength))) {
                    throw new PlcRuntimeException(String.format("%s: Unexpected bitLength '%s'. Want '%d'", logicalName, attribute.getValue(), bitLength));
                }
                bitLengthValidate = true;
            }
        }
        if (!dataTypeValidated) {
            throw new PlcRuntimeException(String.format("%s: required attribute %s missing", logicalName, rwDataTypeKey));
        }
        if (!bitLengthValidate) {
            throw new PlcRuntimeException(String.format("%s: required attribute %s missing", logicalName, rwBitLengthKey));
        }
    }

}

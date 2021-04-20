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

import javax.xml.stream.XMLEventReader;
import javax.xml.stream.XMLInputFactory;
import javax.xml.stream.XMLStreamException;
import javax.xml.stream.events.*;
import java.io.InputStream;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.Iterator;

public class ReadBufferXmlBased implements ReadBuffer {

    XMLEventReader reader;

    int pos = 1;

    public ReadBufferXmlBased(InputStream is) {
        XMLInputFactory xmlInputFactory = XMLInputFactory.newInstance();
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
    public void pullContext(String logicalName) {
        StartElement startElement = travelToNextStartElement();
        String elementName = startElement.getName().getLocalPart();
        if (!elementName.equals(logicalName)) {
            throw new PlcRuntimeException(String.format("Unexpected Start element '%s'. Expected '%s'", elementName, logicalName));
        }
    }

    @Override
    public boolean readBit(String logicalName) throws ParseException {
        String bit = decode(logicalName, "bit", 1);
        move(1);
        return bit.equals("1");
    }

    @Override
    public byte readUnsignedByte(String logicalName, int bitLength) throws ParseException {
        move(bitLength);
        return Byte.parseByte(decode(logicalName, "uint8", bitLength));
    }

    @Override
    public short readUnsignedShort(String logicalName, int bitLength) throws ParseException {
        move(bitLength);
        return Short.parseShort(decode(logicalName, "uint16", bitLength));
    }

    @Override
    public int readUnsignedInt(String logicalName, int bitLength) throws ParseException {
        move(bitLength);
        return Integer.parseInt(decode(logicalName, "uint32", bitLength));
    }

    @Override
    public long readUnsignedLong(String logicalName, int bitLength) throws ParseException {
        move(bitLength);
        return Long.parseLong(decode(logicalName, "uint32", bitLength));
    }

    @Override
    public BigInteger readUnsignedBigInteger(String logicalName, int bitLength) throws ParseException {
        move(bitLength);
        throw new PlcRuntimeException("not implemented yet");
    }

    @Override
    public byte readByte(String logicalName, int bitLength) throws ParseException {
        move(bitLength);
        return Byte.parseByte(decode(logicalName, "int8", bitLength));
    }

    @Override
    public short readShort(String logicalName, int bitLength) throws ParseException {
        move(bitLength);
        return Short.parseShort(decode(logicalName, "int16", bitLength));
    }

    @Override
    public int readInt(String logicalName, int bitLength) throws ParseException {
        move(bitLength);
        return Integer.parseInt(decode(logicalName, "int32", bitLength));
    }

    @Override
    public long readLong(String logicalName, int bitLength) throws ParseException {
        move(bitLength);
        return Long.parseLong(decode(logicalName, "int64", bitLength));
    }

    @Override
    public BigInteger readBigInteger(String logicalName, int bitLength) throws ParseException {
        move(bitLength);
        throw new PlcRuntimeException("not implemented yet");
    }

    @Override
    public float readFloat(String logicalName, int bitLength) throws ParseException {
        move(bitLength);
        return Float.parseFloat(decode(logicalName, "float32", bitLength));
    }

    @Override
    public double readDouble(String logicalName, int bitLength) throws ParseException {
        move(bitLength);
        return Double.parseDouble(decode(logicalName, "float64", bitLength));
    }

    @Override
    public BigDecimal readBigDecimal(String logicalName, int bitLength) throws ParseException {
        move(bitLength);
        throw new PlcRuntimeException("not implemented yet");
    }

    @Override
    public String readString(String logicalName, int bitLength, String encoding) {
        move(bitLength);
        return decode(logicalName, "string", bitLength);
    }

    @Override
    public void closeContext(String logicalName) {
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
                throw new PlcRuntimeException("Unexpected Start element" + xmlEvent.asEndElement().getName().getLocalPart());
            } else if (xmlEvent.isEndElement()) {
                return xmlEvent.asEndElement();
            }
        }
        throw new PlcRuntimeException("EOF");
    }

    private String decode(String logicalName, String dataType, int bitLength) {
        StartElement startElement = travelToNextStartElement();
        validateStartElement(startElement, logicalName, dataType, bitLength);
        Characters characters = null;
        try {
            characters = reader.nextEvent().asCharacters();
        } catch (XMLStreamException e) {
            throw new PlcRuntimeException(e);
        }
        String data = characters.getData();
        try {
            XMLEvent endEvent = reader.nextEvent().asEndElement();
        } catch (XMLStreamException e) {
            throw new PlcRuntimeException(e);
        }
        return data;
    }

    private boolean validateStartElement(StartElement startElement, String logicalName, String dataType, int bitLength) {
        logicalName = sanitizeLogicalName(logicalName);
        if (!startElement.getName().getLocalPart().equals(logicalName)) {
            throw new PlcRuntimeException(String.format("unexpected element '%s'. Expected '%s'", startElement.getName().getLocalPart(), logicalName));
        } else if (!validateAttr(startElement.getAttributes(), dataType, bitLength)) {
            throw new PlcRuntimeException("Error validating Attributes");
        }
        return true;
    }

    private boolean validateAttr(Iterator<Attribute> attr, String dataType, int bitLength) {
        boolean dataTypeValidated = false;
        boolean bitLengthValidate = false;
        while (attr.hasNext()) {
            Attribute attribute = attr.next();
            if (attribute.getName().getLocalPart().equals("dataType")) {
                if (!attribute.getValue().equals(dataType)) {
                    throw new PlcRuntimeException(String.format("Unexpected dataType :%s. Want %s", attribute.getValue(), dataType));
                }
                dataTypeValidated = true;
            } else if (attribute.getName().getLocalPart().equals("bitLength")) {
                if (!attribute.getValue().equals(Integer.valueOf(bitLength).toString())) {
                    throw new PlcRuntimeException(String.format("Unexpected bitLength '%s'. Want '%d'", attribute.getValue(), bitLength));
                }
                bitLengthValidate = true;
            }
        }
        if (!dataTypeValidated) {
            throw new PlcRuntimeException("required attribute dataType missing");
        }
        if (!bitLengthValidate) {
            throw new PlcRuntimeException("required attribute bitLength missing");
        }
        return true;
    }

    private String sanitizeLogicalName(String logicalName) {
        if (logicalName.equals("")) {
            return "value";
        }
        return logicalName;
    }
}

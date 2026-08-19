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
import org.apache.plc4x.java.spi.buffers.api.ReadBuffer;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;

import javax.xml.stream.XMLEventReader;
import javax.xml.stream.XMLInputFactory;
import javax.xml.stream.XMLStreamException;
import javax.xml.stream.events.*;
import java.io.InputStream;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.Iterator;
import java.util.Optional;

public class ReadBufferXmlBased extends AbstractBuffer implements ReadBuffer, XmlBuffer {

    private final int sizeInBits;
    private int positionInBits;

    XMLEventReader reader;

    public ReadBufferXmlBased(InputStream is) {
        XMLInputFactory xmlInputFactory = XMLInputFactory.newInstance();
        xmlInputFactory.setProperty(XMLInputFactory.SUPPORT_DTD, false);
        // TODO: Set this ....
        sizeInBits = 0;
        positionInBits = 0;
        try {
            reader = xmlInputFactory.createXMLEventReader(is);
        } catch (XMLStreamException e) {
            throw new IllegalArgumentException(e);
        }
    }

    @Override
    public void pushContext(WithOption... options) throws BufferException {
        String name = getName(options);
        StartElement startElement = travelToNextStartElement();
        String elementName = startElement.getName().getLocalPart();
        if (!elementName.equals(name)) {
            throw new BufferException(String.format("Unexpected start element '%s'. Expected '%s'", elementName, name));
        }
        super.pushContext(options);
    }

    @Override
    public void popContext(WithOption... options) throws BufferException {
        String name = getName(WithOption.AddOptions(getContext(), options));
        EndElement endElement = travelToNextEndElement();
        String elementName = endElement.getName().getLocalPart();
        if (!elementName.equals(name)) {
            throw new BufferException(String.format("Unexpected end element '%s'. Expected '%s'", elementName, name));
        }
        super.popContext(options);
    }

    @Override
    public boolean readBit(WithOption... options) throws BufferException {
        String name = getName(WithOption.AddOptions(getContext(), options));
        String bitString = decode(name, rwBitKey, 1);
        move(1);
        if (!"true".equalsIgnoreCase(bitString) && !"false".equalsIgnoreCase(bitString)) {
            throw new BufferException(String.format("Unexpected bit value '%s'. Expected 'true' or 'false'", bitString));
        }
        return Boolean.parseBoolean(bitString);
    }

    @Override
    public byte[] readBits(int numBits, WithOption... options) throws BufferException {
        String name = getName(options);
        String hexString = decode(name, rwByteKey, numBits);
        if (!hexString.startsWith("0x")) {
            throw new BufferException(String.format("Hex string should start with 0x. Actual value %s", hexString));
        }
        hexString = hexString.substring(2);
        int numberOfBytes = (numBits + 7) / 8;
        if (hexString.length() != numberOfBytes * 2) {
            throw new BufferException(String.format("Hex string should be %d bytes long. Actual value %s", numberOfBytes, hexString));
        }
        byte[] bytes = new byte[numberOfBytes];
        for (int i = 0; i < hexString.length(); i = i + 2) {
            // Without this hack, we can't parse values such as "0x80" into a byte.
            bytes[i / 2] = (byte) Short.parseShort(hexString.substring(i, i + 2), 16);
        }
        return bytes;
    }

    @Override
    public byte readUnsignedByte(int numBits, WithOption... options) throws BufferException {
        String name = getScalarName(options);
        move(numBits);
        return Byte.parseByte(decode(name, rwUintKey, numBits));
    }

    @Override
    public short readUnsignedShort(int numBits, WithOption... options) throws BufferException {
        String name = getScalarName(options);
        move(numBits);
        return Short.parseShort(decode(name, rwUintKey, numBits));
    }

    @Override
    public int readUnsignedInt(int numBits, WithOption... options) throws BufferException {
        String name = getScalarName(options);
        move(numBits);
        return Integer.parseInt(decode(name, rwUintKey, numBits));
    }

    @Override
    public long readUnsignedLong(int numBits, WithOption... options) throws BufferException {
        String name = getScalarName(options);
        move(numBits);
        return Long.parseLong(decode(name, rwUintKey, numBits));
    }

    @Override
    public BigInteger readUnsignedBigInteger(int numBits, WithOption... options) throws BufferException {
        String name = getScalarName(options);
        move(numBits);
        String bigIntString = decode(name, rwUintKey, numBits);
        return new BigInteger(bigIntString);
    }

    @Override
    public byte readSignedByte(int numBits, WithOption... options) throws BufferException {
        String name = getScalarName(options);
        move(numBits);
        return Byte.parseByte(decode(name, rwIntKey, numBits));
    }

    @Override
    public short readSignedShort(int numBits, WithOption... options) throws BufferException {
        String name = getScalarName(options);
        move(numBits);
        return Short.parseShort(decode(name, rwIntKey, numBits));
    }

    @Override
    public int readSignedInt(int numBits, WithOption... options) throws BufferException {
        String name = getScalarName(options);
        move(numBits);
        return Integer.parseInt(decode(name, rwIntKey, numBits));
    }

    @Override
    public long readSignedLong(int numBits, WithOption... options) throws BufferException {
        String name = getScalarName(options);
        move(numBits);
        return Long.parseLong(decode(name, rwIntKey, numBits));
    }

    @Override
    public BigInteger readSignedBigInteger(int numBits, WithOption... options) throws BufferException {
        String name = getScalarName(options);
        move(numBits);
        return new BigInteger(decode(name, rwIntKey, numBits));
    }

    @Override
    public float readFloat(int numBits, WithOption... options) throws BufferException {
        String name = getScalarName(options);
        move(numBits);
        return Float.parseFloat(decode(name, rwFloatKey, numBits));
    }

    @Override
    public double readDouble(int numBits, WithOption... options) throws BufferException {
        String name = getScalarName(options);
        move(numBits);
        return Double.parseDouble(decode(name, rwFloatKey, numBits));
    }

    @Override
    public BigDecimal readBigDecimal(int numBits, WithOption... options) throws BufferException {
        String name = getScalarName(options);
        move(numBits);
        return new BigDecimal(decode(name, rwFloatKey, numBits));
    }

    @Override
    public String readString(int numBits, WithOption... options) throws BufferException {
        String name = getScalarName(options);
        move(numBits);
        return decode(name, rwStringKey, numBits);
    }

    @Override
    public ReadBuffer createSubBuffer(int numBits, WithOption... options) {
        // TODO: Find out how to do this.
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

    @Override
    public void setPositionInBits(int positionInBits) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    private void move(int bits) {
        positionInBits += bits;
    }

    private StartElement travelToNextStartElement() throws BufferException {
        while (reader.hasNext()) {
            XMLEvent xmlEvent;
            try {
                xmlEvent = reader.nextEvent();
            } catch (XMLStreamException e) {
                throw new BufferException("Error reading next element", e);
            }
            if (xmlEvent.isStartElement()) {
                return xmlEvent.asStartElement();
            } else if (xmlEvent.isEndElement()) {
                throw new BufferException(String.format("Unexpected End element '%s'", xmlEvent.asEndElement().getName().getLocalPart()));
            }
        }
        throw new BufferException("EOF");
    }

    private EndElement travelToNextEndElement() throws BufferException {
        while (reader.hasNext()) {
            XMLEvent xmlEvent;
            try {
                xmlEvent = reader.nextEvent();
            } catch (XMLStreamException e) {
                throw new BufferException("Error reading next element", e);
            }
            if (xmlEvent.isStartElement()) {
                throw new BufferException(String.format("Unexpected Start element '%s'", xmlEvent.asStartElement().getName().getLocalPart()));
            } else if (xmlEvent.isEndElement()) {
                return xmlEvent.asEndElement();
            }
        }
        throw new BufferException("EOF");
    }

    private String decode(String name, String dataType, int numBits) throws BufferException {
        StartElement startElement = travelToNextStartElement();
        validateStartElement(startElement, name, dataType, numBits);
        Characters characters;
        try {
            characters = reader.nextEvent().asCharacters();
        } catch (XMLStreamException e) {
            throw new BufferException("Error reading characters.", e);
        }
        String data = characters.getData();
        try {
            XMLEvent xmlEvent = reader.nextEvent();
            if (!xmlEvent.isEndElement()) {
                throw new BufferException("Expected end element.");
            }
        } catch (XMLStreamException e) {
            throw new BufferException("Error reading end element.", e);
        }
        return data.trim();
    }

    private void validateStartElement(StartElement startElement, String name, String dataType, int numBits) throws BufferException {
        //name = sanitizeLogicalName(name);
        String elementName = startElement.getName().getLocalPart();
        // "*" matches any element name: used for elements whose serialized name is only known
        // dynamically (e.g. enum fields are written as <fieldName><EnumTypeName ...>).
        if (!"*".equals(name) && !elementName.equals(name) && !isReservedFieldName(name, elementName)) {
            throw new BufferException(String.format("unexpected element '%s'. Expected '%s'", elementName, name));
        }
        validateAttr(startElement.getAttributes(), name, dataType, numBits);
    }

    private boolean isReservedFieldName(String expectedName, String elementName) {
        // The generated parsers qualify reserved fields as "<Type>.reserved<N>" while the
        // generated serializers write them plainly as "reserved" - accept the serializer's
        // form so XML written by WriteBufferXmlBased can be read back.
        return "reserved".equals(elementName) && expectedName.matches("(.*\\.)?reserved\\d*");
    }

    private void validateAttr(Iterator<Attribute> attr, String name, String dataType, int numBits) throws BufferException {
        boolean dataTypeValidated = false;
        boolean numBitsValidate = false;
        while (attr.hasNext()) {
            Attribute attribute = attr.next();
            if (attribute.getName().getLocalPart().equals(rwDataTypeKey)) {
                if (!attribute.getValue().equals(dataType)) {
                    throw new BufferException(String.format("%s: Unexpected dataType '%s'. Want '%s'", name, attribute.getValue(), dataType));
                }
                dataTypeValidated = true;
            } else if (attribute.getName().getLocalPart().equals(rwBitLengthKey)) {
                if (!attribute.getValue().equals(Integer.toString(numBits))) {
                    throw new BufferException(String.format("%s: Unexpected numBits '%s'. Want '%d'", name, attribute.getValue(), numBits));
                }
                numBitsValidate = true;
            }
        }
        if (!dataTypeValidated) {
            throw new BufferException(String.format("%s: required attribute '%s' missing", name, rwDataTypeKey));
        }
        if (!numBitsValidate) {
            throw new BufferException(String.format("%s: required attribute '%s' missing", name, rwBitLengthKey));
        }
    }

    /**
     * Resolves the element name for a scalar read.
     *
     * <p>Mirrors {@code FieldWriterArray.writeSimpleTypeArrayField}, which serializes every simple
     * list item as {@code <value>} and discards the field's own name. Reads carrying the
     * render-as-list marker are exactly those items, so without this a list serialized as
     * {@code <classId isList="true"><value/>...} fails to read back with
     * "unexpected element 'value'. Expected 'classId'". Only scalar reads go through here:
     * byte arrays are written as a single element, and complex list items push their own type
     * name as context.
     */
    private String getScalarName(WithOption... options) throws BufferException {
        if (WithOption.extractRenderAsList(options).orElse(false)) {
            return "value";
        }
        return getName(options);
    }

    protected String getName(WithOption... options) throws BufferException {
        Optional<String> name = WithOption.extractName(options, getContext());
        if (name.isEmpty()) {
            throw new BufferException("Missing 'name' option.");
        }
        return name.get();
    }

}

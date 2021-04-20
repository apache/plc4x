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
import org.w3c.dom.Document;
import org.w3c.dom.Element;

import javax.xml.parsers.DocumentBuilderFactory;
import javax.xml.parsers.ParserConfigurationException;
import javax.xml.transform.OutputKeys;
import javax.xml.transform.Transformer;
import javax.xml.transform.TransformerException;
import javax.xml.transform.TransformerFactory;
import javax.xml.transform.dom.DOMSource;
import javax.xml.transform.stream.StreamResult;
import java.io.ByteArrayOutputStream;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.Stack;

public class WriteBufferXmlBased implements WriteBuffer {

    Stack<Element> stack;

    Document document;

    public WriteBufferXmlBased() {
        DocumentBuilderFactory documentBuilderFactory = DocumentBuilderFactory.newInstance();
        try {
            this.document = documentBuilderFactory.newDocumentBuilder().newDocument();
        } catch (ParserConfigurationException e) {
            throw new PlcRuntimeException(e);
        }
        this.stack = new Stack<>();
    }

    @Override
    public int getPos() {
        return 0;
    }

    @Override
    public void pushContext(String logicalName) {
        stack.push(document.createElement(logicalName));
    }

    @Override
    public void writeBit(String logicalName, boolean value) throws ParseException {
        Element element = document.createElement(sanitizeLogicalName(logicalName));
        element.setAttribute("dataType", "bit");
        element.setAttribute("bitLength", "1");
        element.appendChild(document.createTextNode(Boolean.valueOf(value).toString()));
        stack.peek().appendChild(element);
    }

    @Override
    public void writeUnsignedByte(String logicalName, int bitLength, byte value) throws ParseException {
        Element element = document.createElement(sanitizeLogicalName(logicalName));
        element.setAttribute("dataType", "uint8");
        element.setAttribute("bitLength", String.valueOf(bitLength));
        element.appendChild(document.createTextNode(Byte.valueOf(value).toString()));
        stack.peek().appendChild(element);
    }

    @Override
    public void writeUnsignedShort(String logicalName, int bitLength, short value) throws ParseException {
        Element element = document.createElement(sanitizeLogicalName(logicalName));
        element.setAttribute("dataType", "uint16");
        element.setAttribute("bitLength", String.valueOf(bitLength));
        element.appendChild(document.createTextNode(Short.valueOf(value).toString()));
        stack.peek().appendChild(element);
    }

    @Override
    public void writeUnsignedInt(String logicalName, int bitLength, int value) throws ParseException {
        Element element = document.createElement(sanitizeLogicalName(logicalName));
        element.setAttribute("dataType", "uint32");
        element.setAttribute("bitLength", String.valueOf(bitLength));
        element.appendChild(document.createTextNode(Integer.valueOf(value).toString()));
        stack.peek().appendChild(element);
    }

    @Override
    public void writeUnsignedLong(String logicalName, int bitLength, long value) throws ParseException {
        Element element = document.createElement(sanitizeLogicalName(logicalName));
        element.setAttribute("dataType", "uint64");
        element.setAttribute("bitLength", String.valueOf(bitLength));
        element.appendChild(document.createTextNode(Long.valueOf(value).toString()));
        stack.peek().appendChild(element);
    }

    @Override
    public void writeUnsignedBigInteger(String logicalName, int bitLength, BigInteger value) throws ParseException {
        Element element = document.createElement(sanitizeLogicalName(logicalName));
        element.setAttribute("dataType", "bigInt");
        element.setAttribute("bitLength", String.valueOf(bitLength));
        element.appendChild(document.createTextNode(value.toString()));
        stack.peek().appendChild(element);
    }

    @Override
    public void writeByte(String logicalName, int bitLength, byte value) throws ParseException {
        Element element = document.createElement(sanitizeLogicalName(logicalName));
        element.setAttribute("dataType", "int8");
        element.setAttribute("bitLength", String.valueOf(bitLength));
        element.appendChild(document.createTextNode(Byte.valueOf(value).toString()));
        stack.peek().appendChild(element);
    }

    @Override
    public void writeShort(String logicalName, int bitLength, short value) throws ParseException {
        Element element = document.createElement(sanitizeLogicalName(logicalName));
        element.setAttribute("dataType", "int16");
        element.setAttribute("bitLength", String.valueOf(bitLength));
        element.appendChild(document.createTextNode(Short.valueOf(value).toString()));
        stack.peek().appendChild(element);
    }

    @Override
    public void writeInt(String logicalName, int bitLength, int value) throws ParseException {
        Element element = document.createElement(sanitizeLogicalName(logicalName));
        element.setAttribute("dataType", "int32");
        element.setAttribute("bitLength", String.valueOf(bitLength));
        element.appendChild(document.createTextNode(Integer.valueOf(value).toString()));
        stack.peek().appendChild(element);
    }

    @Override
    public void writeLong(String logicalName, int bitLength, long value) throws ParseException {
        Element element = document.createElement(sanitizeLogicalName(logicalName));
        element.setAttribute("dataType", "int64");
        element.setAttribute("bitLength", String.valueOf(bitLength));
        element.appendChild(document.createTextNode(Long.valueOf(value).toString()));
        stack.peek().appendChild(element);
    }

    @Override
    public void writeBigInteger(String logicalName, int bitLength, BigInteger value) throws ParseException {
        Element element = document.createElement(sanitizeLogicalName(logicalName));
        element.setAttribute("dataType", "bigInt");
        element.setAttribute("bitLength", String.valueOf(bitLength));
        element.appendChild(document.createTextNode(value.toString()));
        stack.peek().appendChild(element);
    }

    @Override
    public void writeFloat(String logicalName, float value, int bitsExponent, int bitsMantissa) throws ParseException {
        Element element = document.createElement(sanitizeLogicalName(logicalName));
        element.setAttribute("dataType", "float32");
        element.setAttribute("bitLength", String.valueOf((value < 0 ? 1 : 0) + bitsExponent + bitsMantissa));
        element.appendChild(document.createTextNode(Float.valueOf(value).toString()));
        stack.peek().appendChild(element);
    }

    @Override
    public void writeDouble(String logicalName, double value, int bitsExponent, int bitsMantissa) throws ParseException {
        Element element = document.createElement(sanitizeLogicalName(logicalName));
        element.setAttribute("dataType", "float64");
        element.setAttribute("bitLength", String.valueOf((value < 0 ? 1 : 0) + bitsExponent + bitsMantissa));
        element.appendChild(document.createTextNode(Double.valueOf(value).toString()));
        stack.peek().appendChild(element);
    }

    @Override
    public void writeBigDecimal(String logicalName, int bitLength, BigDecimal value) throws ParseException {
        Element element = document.createElement(sanitizeLogicalName(logicalName));
        element.setAttribute("dataType", "bigFloat");
        element.setAttribute("bitLength", String.valueOf(bitLength));
        element.appendChild(document.createTextNode(value.toString()));
        stack.peek().appendChild(element);
    }

    @Override
    public void writeString(String logicalName, int bitLength, String encoding, String value) throws ParseException {
        Element element = document.createElement(sanitizeLogicalName(logicalName));
        element.setAttribute("dataType", "string");
        element.setAttribute("bitLength", String.valueOf(bitLength));
        element.appendChild(document.createTextNode(value));
        stack.peek().appendChild(element);
    }

    @Override
    public void popContext(String logicalName) {
        Element currentContext = stack.pop();
        if (!currentContext.getTagName().equals(logicalName)) {
            throw new PlcRuntimeException("Unexpected pop context '" + currentContext.getTagName() + '\'');
        }
        if (stack.isEmpty()) {
            document.appendChild(currentContext);
            return;
        }
        stack.peek().appendChild(currentContext);
    }

    public String getXmlString() {
        try {
            Transformer transformer = TransformerFactory.newInstance().newTransformer();
            transformer.setOutputProperty(OutputKeys.INDENT, "yes");
            transformer.setOutputProperty(OutputKeys.METHOD, "xml");
            transformer.setOutputProperty(OutputKeys.ENCODING, "UTF-8");
            transformer.setOutputProperty(OutputKeys.OMIT_XML_DECLARATION, "yes");
            //transformer.setOutputProperty(OutputKeys.DOCTYPE_SYSTEM, "roles.dtd");
            transformer.setOutputProperty("{http://xml.apache.org/xslt}indent-amount", "4");

            ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream();
            transformer.transform(new DOMSource(document),
                new StreamResult(byteArrayOutputStream));

            return byteArrayOutputStream.toString();
        } catch (TransformerException e) {
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

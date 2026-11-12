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
package org.apache.plc4x.java.profinet.gsdml;

import org.w3c.dom.Attr;
import org.w3c.dom.Element;
import org.w3c.dom.NamedNodeMap;
import org.w3c.dom.Node;
import org.w3c.dom.Text;
import org.xml.sax.ErrorHandler;
import org.xml.sax.SAXException;
import org.xml.sax.SAXParseException;

import javax.xml.XMLConstants;
import javax.xml.parsers.DocumentBuilder;
import javax.xml.parsers.DocumentBuilderFactory;
import javax.xml.parsers.ParserConfigurationException;
import java.io.IOException;
import java.io.InputStream;
import java.util.ArrayList;
import java.util.List;
import java.util.regex.Pattern;

/**
 * Reads GSDML documents with the JDK's DOM parser, without reflection.
 * <p>
 * Each GSDML class reads its own properties from {@link #properties(Element, boolean)}. The rules
 * are the ones the jackson-dataformat-xml mapping of these classes followed, so device descriptions
 * keep loading the way they did:
 * <ul>
 *     <li>Attributes and child elements are matched by local name; namespaces and prefixes are
 *     ignored. Either can supply any property. Attributes come first, then the child elements in
 *     document order, and a later occurrence replaces an earlier one (no merging).</li>
 *     <li>Consecutive child elements of the same name form one run. An unwrapped list is such a
 *     run, so a second run of the same name replaces the first.</li>
 *     <li>A wrapped list is one element whose child elements, whatever their name, are the items.</li>
 *     <li>The value of a scalar child element is its text before its first child element.</li>
 *     <li>{@code xsi:nil="true"} makes a property {@code null}, and a list item an empty item.</li>
 *     <li>Missing int and boolean properties are 0 and false. Ints are decimal; booleans are
 *     {@code true}/{@code false} in lower, title or upper case; blank values read as 0 and false.</li>
 *     <li>A class that is "strict" rejects any property it does not know, and any text content.</li>
 * </ul>
 * Document type declarations are rejected outright. GSDML files do not use them, and refusing them
 * rules out entity expansion and any fetching of external resources while parsing.
 */
final class GsdmlXml {

    private static final Pattern DECIMAL_INT = Pattern.compile("[+-]?[0-9]+");

    private static final ErrorHandler FAIL_ON_ERRORS = new ErrorHandler() {
        @Override
        public void warning(SAXParseException exception) {
            // Warnings do not make a document unusable.
        }

        @Override
        public void error(SAXParseException exception) throws SAXException {
            throw exception;
        }

        @Override
        public void fatalError(SAXParseException exception) throws SAXException {
            throw exception;
        }
    };

    private GsdmlXml() {
    }

    /**
     * Reads one GSDML class from an element.
     */
    @FunctionalInterface
    interface Reader<T> {
        T read(Element element) throws IOException;
    }

    /**
     * Parses a document and returns its root element.
     */
    static Element parse(InputStream in) throws IOException {
        try {
            return newDocumentBuilder().parse(in).getDocumentElement();
        } catch (SAXException e) {
            throw new IOException("Malformed GSDML document: " + e.getMessage(), e);
        }
    }

    private static DocumentBuilder newDocumentBuilder() {
        try {
            DocumentBuilderFactory factory = DocumentBuilderFactory.newDefaultInstance();
            factory.setNamespaceAware(true);
            factory.setFeature(XMLConstants.FEATURE_SECURE_PROCESSING, true);
            factory.setFeature("http://apache.org/xml/features/disallow-doctype-decl", true);
            factory.setFeature("http://xml.org/sax/features/external-general-entities", false);
            factory.setFeature("http://xml.org/sax/features/external-parameter-entities", false);
            factory.setFeature("http://apache.org/xml/features/nonvalidating/load-external-dtd", false);
            factory.setAttribute(XMLConstants.ACCESS_EXTERNAL_DTD, "");
            factory.setAttribute(XMLConstants.ACCESS_EXTERNAL_SCHEMA, "");
            factory.setXIncludeAware(false);
            factory.setExpandEntityReferences(false);
            DocumentBuilder builder = factory.newDocumentBuilder();
            builder.setErrorHandler(FAIL_ON_ERRORS);
            return builder;
        } catch (ParserConfigurationException e) {
            throw new IllegalStateException("The JDK XML parser does not support the settings GSDML parsing relies on", e);
        }
    }

    /**
     * Lists the properties of an element in the order they are applied: attributes first, then runs
     * of same-named child elements in document order.
     *
     * @param strict whether text content is an error, for classes that reject unknown properties
     */
    static List<Property> properties(Element element, boolean strict) throws IOException {
        List<Property> properties = new ArrayList<>();
        NamedNodeMap attributes = element.getAttributes();
        for (int i = 0; i < attributes.getLength(); i++) {
            Attr attribute = (Attr) attributes.item(i);
            if (!isNamespaceDeclaration(attribute) && !isNilAttribute(attribute)) {
                properties.add(new Property(attribute.getLocalName(), attribute, null));
            }
        }
        Property run = null;
        for (Node child = element.getFirstChild(); child != null; child = child.getNextSibling()) {
            if (child instanceof Element childElement) {
                String name = childElement.getLocalName();
                if (run != null && run.name.equals(name)) {
                    run.elements.add(childElement);
                } else {
                    run = new Property(name, null, new ArrayList<>(List.of(childElement)));
                    properties.add(run);
                }
            } else if (strict && child instanceof Text text && !text.getData().isBlank()) {
                throw new IOException("Unexpected text content in GSDML element " + element.getLocalName());
            }
        }
        return properties;
    }

    private static boolean isNamespaceDeclaration(Attr attribute) {
        return XMLConstants.XMLNS_ATTRIBUTE_NS_URI.equals(attribute.getNamespaceURI());
    }

    private static boolean isNilAttribute(Attr attribute) {
        return XMLConstants.W3C_XML_SCHEMA_INSTANCE_NS_URI.equals(attribute.getNamespaceURI())
            && "nil".equals(attribute.getLocalName());
    }

    private static boolean isNil(Element element) {
        return "true".equals(element.getAttributeNS(XMLConstants.W3C_XML_SCHEMA_INSTANCE_NS_URI, "nil").trim());
    }

    /**
     * The text of an element up to its first child element, or {@code null} for a nil element.
     */
    private static String leadingText(Element element) {
        if (isNil(element)) {
            return null;
        }
        StringBuilder text = new StringBuilder();
        for (Node child = element.getFirstChild(); child != null && !(child instanceof Element); child = child.getNextSibling()) {
            if (child instanceof Text textNode) {
                text.append(textNode.getData());
            }
        }
        return text.toString();
    }

    private static boolean hasChildElements(Element element) {
        for (Node child = element.getFirstChild(); child != null; child = child.getNextSibling()) {
            if (child instanceof Element) {
                return true;
            }
        }
        return false;
    }

    private static boolean hasAttributes(Element element) {
        NamedNodeMap attributes = element.getAttributes();
        for (int i = 0; i < attributes.getLength(); i++) {
            if (!isNamespaceDeclaration((Attr) attributes.item(i))) {
                return true;
            }
        }
        return false;
    }

    private static boolean hasText(Element element) {
        for (Node child = element.getFirstChild(); child != null; child = child.getNextSibling()) {
            if (child instanceof Text text && !text.getData().isBlank()) {
                return true;
            }
        }
        return false;
    }

    private static Element emptyElement(Node node) {
        return node.getOwnerDocument().createElementNS(null, node.getLocalName());
    }

    /**
     * One occurrence of a property: either an attribute or a run of same-named child elements.
     */
    static final class Property {

        private final String name;
        private final Attr attribute;
        private final List<Element> elements;

        private Property(String name, Attr attribute, List<Element> elements) {
            this.name = name;
            this.attribute = attribute;
            this.elements = elements;
        }

        String name() {
            return name;
        }

        String string() {
            if (attribute != null) {
                return attribute.getValue();
            }
            return leadingText(elements.get(elements.size() - 1));
        }

        int intValue() throws IOException {
            if (attribute != null) {
                return toInt(attribute.getValue());
            }
            int value = 0;
            for (Element element : elements) {
                value = toInt(leadingText(element));
            }
            return value;
        }

        boolean booleanValue() throws IOException {
            if (attribute != null) {
                return toBoolean(attribute.getValue());
            }
            boolean value = false;
            for (Element element : elements) {
                value = toBoolean(leadingText(element));
            }
            return value;
        }

        <T> T object(Reader<T> reader) throws IOException {
            if (attribute != null) {
                return reader.read(emptyFromBlankAttribute());
            }
            T value = null;
            for (Element element : elements) {
                value = isNil(element) ? null : readObject(element, reader);
            }
            return value;
        }

        <T> List<T> wrappedList(Reader<T> reader) throws IOException {
            if (attribute != null) {
                emptyFromBlankAttribute();
                return new ArrayList<>();
            }
            List<T> items = null;
            for (Element wrapper : elements) {
                if (isNil(wrapper)) {
                    items = null;
                    continue;
                }
                if (hasText(wrapper)) {
                    throw new IOException("Unexpected text content in GSDML list " + name);
                }
                items = new ArrayList<>();
                for (Node child = wrapper.getFirstChild(); child != null; child = child.getNextSibling()) {
                    if (child instanceof Element item) {
                        items.add(readItem(item, reader));
                    }
                }
            }
            return items;
        }

        <T> List<T> unwrappedList(Reader<T> reader) throws IOException {
            if (attribute != null) {
                throw new IOException("GSDML list " + name + " cannot be given as an attribute");
            }
            List<T> items = new ArrayList<>();
            for (Element element : elements) {
                items.add(readItem(element, reader));
            }
            return items;
        }

        /**
         * The error for a property a strict class does not know.
         */
        IOException unknownIn(Class<?> type) {
            return new IOException("Unrecognized GSDML " + (attribute != null ? "attribute" : "element")
                + " \"" + name + "\" in " + type.getSimpleName());
        }

        private Element emptyFromBlankAttribute() throws IOException {
            if (!attribute.getValue().isBlank()) {
                throw new IOException("GSDML property " + name + " needs an element, not the attribute value \""
                    + attribute.getValue() + "\"");
            }
            return emptyElement(attribute);
        }

        private <T> T readItem(Element element, Reader<T> reader) throws IOException {
            return isNil(element) ? reader.read(emptyElement(element)) : readObject(element, reader);
        }

        private <T> T readObject(Element element, Reader<T> reader) throws IOException {
            if (!hasAttributes(element) && !hasChildElements(element) && hasText(element)) {
                throw new IOException("GSDML element " + name + " needs attributes or child elements, not text");
            }
            return reader.read(element);
        }

        private int toInt(String value) throws IOException {
            String trimmed = value == null ? "" : value.trim();
            if (trimmed.isEmpty()) {
                return 0;
            }
            if (DECIMAL_INT.matcher(trimmed).matches()) {
                try {
                    return Integer.parseInt(trimmed);
                } catch (NumberFormatException e) {
                    // Out of range, reported below.
                }
            }
            throw new IOException("GSDML property " + name + " is not a valid int: \"" + value + "\"");
        }

        private boolean toBoolean(String value) throws IOException {
            String trimmed = value == null ? "" : value.trim();
            return switch (trimmed) {
                case "", "false", "False", "FALSE" -> false;
                case "true", "True", "TRUE" -> true;
                default -> throw new IOException("GSDML property " + name + " is not a valid boolean: \"" + value + "\"");
            };
        }
    }
}

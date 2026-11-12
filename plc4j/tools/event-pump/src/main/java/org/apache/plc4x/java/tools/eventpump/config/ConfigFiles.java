/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.apache.plc4x.java.tools.eventpump.config;

import org.apache.plc4x.java.spi.utils.json.Json;
import org.snakeyaml.engine.v2.api.Dump;
import org.snakeyaml.engine.v2.api.DumpSettings;
import org.snakeyaml.engine.v2.api.Load;
import org.snakeyaml.engine.v2.api.LoadSettings;
import org.snakeyaml.engine.v2.common.FlowStyle;
import org.snakeyaml.engine.v2.exceptions.YamlEngineException;
import org.w3c.dom.Document;
import org.w3c.dom.Element;
import org.w3c.dom.Node;
import org.w3c.dom.NodeList;
import org.xml.sax.SAXException;

import javax.xml.XMLConstants;
import javax.xml.parsers.DocumentBuilder;
import javax.xml.parsers.DocumentBuilderFactory;
import javax.xml.parsers.ParserConfigurationException;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.regex.Pattern;

/**
 * Reads and writes the three file formats of an event pump configuration. Reading turns any of
 * them into the same plain tree (maps, lists, scalars) that the configuration classes bind from;
 * writing does the reverse.
 * <p>
 * The XML layout is the one Jackson's XmlMapper used to read and write, so existing files keep
 * working: the root element's name doesn't matter, every property is an element, and a list is
 * a wrapper element holding one element per entry (both named after the property):
 * <pre>{@code
 * <EventPumpConfiguration>
 *   <connections>
 *     <connections><id>plc1</id><url>...</url></connections>
 *   </connections>
 * </EventPumpConfiguration>
 * }</pre>
 */
final class ConfigFiles {

    static final String XML_ROOT = "EventPumpConfiguration";

    /** Element names this writer can produce (a subset of XML names, which is all a config needs). */
    private static final Pattern XML_NAME = Pattern.compile("[A-Za-z_][A-Za-z0-9_.-]*");

    private ConfigFiles() {
    }

    static Object readJson(InputStream inputStream) throws IOException {
        return Json.parse(new InputStreamReader(inputStream, StandardCharsets.UTF_8));
    }

    static Object readYaml(InputStream inputStream) throws IOException {
        LoadSettings settings = LoadSettings.builder().setLabel("event pump configuration").build();
        try {
            return new Load(settings).loadFromInputStream(inputStream);
        } catch (YamlEngineException e) {
            throw new IOException("Invalid YAML: " + e.getMessage(), e);
        }
    }

    /**
     * Reads XML into the plain tree: an element with child elements becomes a map (children of the
     * same name become a list), an element without becomes its text, or null if it has none.
     */
    static Object readXml(InputStream inputStream) throws IOException {
        Document document;
        try {
            document = newDocumentBuilder().parse(inputStream);
        } catch (SAXException e) {
            throw new IOException("Invalid XML: " + e.getMessage(), e);
        }
        return xmlToTree(document.getDocumentElement());
    }

    /**
     * A parser that refuses DTDs, so no external entities or entity expansion.
     */
    static DocumentBuilder newDocumentBuilder() throws IOException {
        try {
            DocumentBuilderFactory factory = DocumentBuilderFactory.newInstance();
            factory.setNamespaceAware(true);
            factory.setFeature(XMLConstants.FEATURE_SECURE_PROCESSING, true);
            factory.setFeature("http://apache.org/xml/features/disallow-doctype-decl", true);
            factory.setFeature("http://xml.org/sax/features/external-general-entities", false);
            factory.setFeature("http://xml.org/sax/features/external-parameter-entities", false);
            factory.setXIncludeAware(false);
            factory.setExpandEntityReferences(false);
            DocumentBuilder builder = factory.newDocumentBuilder();
            // Keep the default handler from printing parse errors to stderr; they are thrown anyway.
            builder.setErrorHandler(null);
            return builder;
        } catch (ParserConfigurationException e) {
            throw new IOException("No usable XML parser", e);
        }
    }

    private static Object xmlToTree(Element element) {
        XmlObject children = new XmlObject();
        boolean hasChildElements = false;
        NodeList nodes = element.getChildNodes();
        for (int i = 0; i < nodes.getLength(); i++) {
            if (nodes.item(i) instanceof Element child) {
                hasChildElements = true;
                String name = child.getLocalName() != null ? child.getLocalName() : child.getTagName();
                Object value = xmlToTree(child);
                if (children.containsKey(name)) {
                    Object existing = children.get(name);
                    if (existing instanceof RepeatedElements repeated) {
                        repeated.add(value);
                    } else {
                        RepeatedElements repeated = new RepeatedElements();
                        repeated.add(existing);
                        repeated.add(value);
                        children.put(name, repeated);
                    }
                } else {
                    children.put(name, value);
                }
            }
        }
        if (hasChildElements) {
            return children;
        }
        String text = element.getTextContent();
        return text.isEmpty() ? null : text;
    }

    /**
     * An element with child elements. Its own type so the binding knows the map came from XML,
     * where a list is a wrapper element - in JSON or YAML that reading would hide a mistake.
     */
    static final class XmlObject extends LinkedHashMap<String, Object> {
    }

    /**
     * Sibling elements of the same name.
     */
    static final class RepeatedElements extends ArrayList<Object> {
    }

    static String writeJson(Object tree) throws IOException {
        return Json.toPrettyJson(tree) + "\n";
    }

    static String writeYaml(Object tree) {
        DumpSettings settings = DumpSettings.builder()
            .setDefaultFlowStyle(FlowStyle.BLOCK)
            .setIndent(2)
            .build();
        return new Dump(settings).dumpToString(tree);
    }

    static String writeXml(Map<String, Object> tree) throws IOException {
        StringBuilder out = new StringBuilder();
        writeXmlElement(out, XML_ROOT, tree, 0);
        return out.toString();
    }

    private static void writeXmlElement(StringBuilder out, String name, Object value, int level) throws IOException {
        if (!XML_NAME.matcher(name).matches()) {
            throw new IOException("'" + name + "' can't be written as an XML element name");
        }
        String indent = "  ".repeat(level);
        if (value instanceof Map<?, ?> map) {
            if (map.isEmpty()) {
                out.append(indent).append('<').append(name).append("/>\n");
                return;
            }
            out.append(indent).append('<').append(name).append(">\n");
            for (Map.Entry<?, ?> entry : map.entrySet()) {
                writeXmlElement(out, String.valueOf(entry.getKey()), entry.getValue(), level + 1);
            }
            out.append(indent).append("</").append(name).append(">\n");
        } else if (value instanceof List<?> list) {
            // Jackson's wrapped list: a wrapper, and one element per entry, all of the same name.
            Map<String, Object> wrapper = new LinkedHashMap<>();
            if (list.isEmpty()) {
                writeXmlElement(out, name, wrapper, level);
                return;
            }
            out.append(indent).append('<').append(name).append(">\n");
            for (Object item : list) {
                writeXmlElement(out, name, item, level + 1);
            }
            out.append(indent).append("</").append(name).append(">\n");
        } else if (value == null) {
            out.append(indent).append('<').append(name).append("/>\n");
        } else {
            out.append(indent).append('<').append(name).append('>');
            appendXmlText(out, String.valueOf(value));
            out.append("</").append(name).append(">\n");
        }
    }

    private static void appendXmlText(StringBuilder out, String text) throws IOException {
        for (int i = 0; i < text.length(); i++) {
            char c = text.charAt(i);
            switch (c) {
                case '&' -> out.append("&amp;");
                case '<' -> out.append("&lt;");
                case '>' -> out.append("&gt;");
                case '\t', '\n', '\r' -> out.append(c);
                default -> {
                    if (c < 0x20 || c == 0xFFFE || c == 0xFFFF) {
                        throw new IOException("Character U+" + String.format("%04X", (int) c)
                            + " can't be written in XML");
                    }
                    out.append(c);
                }
            }
        }
    }

}

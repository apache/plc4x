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
package org.apache.plc4x.java.utils.testutils.utils.xml.comparison;

import org.apache.plc4x.java.utils.testutils.utils.migration.TestCasePatcher;
import org.apache.commons.lang3.RegExUtils;
import org.dom4j.Document;
import org.dom4j.Element;
import org.dom4j.io.OutputFormat;
import org.dom4j.io.XMLWriter;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.File;
import java.io.FileWriter;
import java.io.IOException;
import java.net.URI;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.regex.Pattern;

/**
 * Handles auto-migration of test XML documents.
 * Supports both DOM-based and string-based migration strategies.
 */
public class XmlAutoMigrator {

    private static final Logger LOGGER = LoggerFactory.getLogger(XmlAutoMigrator.class);

    /**
     * Migrates a test document using DOM element replacement.
     * This is more reliable for very large documents.
     *
     * @param expectedElement the expected XML element to replace
     * @param actualElement the actual XML element to use
     * @param documentUri the URI of the test document
     * @throws XmlMigrationException if migration fails
     */
    public static void migrateDom(Element expectedElement, Element actualElement, URI documentUri) {
        try {
            // Get the document reference BEFORE we start modifying the tree
            Document document = expectedElement.getDocument();
            if (document == null) {
                throw new XmlMigrationException("Expected element is not attached to a document");
            }

            // Get the parent element to replace the old response with the new one
            Element parent = expectedElement.getParent();
            if (parent != null) {
                // Create a detached copy of the actual element
                Element newElement = actualElement.createCopy();

                // Get the index of the old element before removing it
                int index = parent.indexOf(expectedElement);

                // Remove the old element
                parent.remove(expectedElement);

                // Add the new element at the same position using DOM4J's content API
                if (index >= 0 && index <= parent.content().size()) {
                    parent.content().add(index, newElement);
                } else {
                    // Fallback: just append if index is invalid
                    parent.add(newElement);
                }
            }

            // Write the updated document back to the file (only if on the filesystem, not in a JAR)
            if (documentUri != null && "file".equals(documentUri.getScheme())) {
                File file = Paths.get(documentUri).toFile();
                try (FileWriter fileWriter = new FileWriter(file)) {
                    OutputFormat format = OutputFormat.createPrettyPrint();
                    format.setIndentSize(2);
                    format.setEncoding("UTF-8");
                    XMLWriter writer = new XMLWriter(fileWriter, format);
                    writer.write(document);
                    writer.close();
                }
                LOGGER.info("Successfully migrated test document using DOM: {}", file.getAbsolutePath());
            } else {
                // Resource is inside a JAR or unavailable — update the in-memory DOM
                // (comparison will pass on this run) but skip writing to disk
                LOGGER.warn("Cannot write auto-migrated XML back to non-file URI: {}. " +
                    "In-memory DOM updated for this run. To persist, run from source.", documentUri);
            }
        } catch (Exception e) {
            throw new XmlMigrationException("Failed to migrate using DOM method", e);
        }
    }

    /**
     * Migrates a test document using string/regex replacement.
     * This preserves formatting better but may have issues with very large documents.
     *
     * @param expectedXml the expected XML string to replace
     * @param actualXml the actual XML string to use
     * @param documentUri the URI of the test document
     * @throws XmlMigrationException if migration fails
     */
    public static void migrateString(String expectedXml, String actualXml, URI documentUri) {
        try {
            Path path = Paths.get(documentUri);
            LOGGER.info("Migrating {} using string replacement", path);

            // Read the current file content
            String content = Files.readString(path, StandardCharsets.UTF_8);
            // Make sure this also works on Windows
            content = content.replaceAll("\r\n", "\n");

            // Determine indent and adjust new XML
            String indent = TestCasePatcher.determineIndent(content, expectedXml);
            String indentedActualXml = TestCasePatcher.indent(actualXml, indent);

            // Create pattern and replace
            Pattern pattern = TestCasePatcher.getPatternForFragment(expectedXml);
            if (!pattern.matcher(content).find()) {
                throw new XmlMigrationException(
                    "Auto migration failed: Can't match content. Pattern matching failed.\n" +
                    "Try to copy the XML manually or use DOM-based migration.");
            }

            content = RegExUtils.replaceFirst((CharSequence) content, pattern, indentedActualXml + "\n");

            // Write back to file
            Files.writeString(path, content, StandardCharsets.UTF_8);

            LOGGER.info("Successfully migrated test document using string replacement: {}", path);
        } catch (IOException e) {
            throw new XmlMigrationException("Failed to migrate using string method", e);
        }
    }

    /**
     * Exception thrown when XML migration fails.
     */
    public static class XmlMigrationException extends RuntimeException {
        public XmlMigrationException(String message) {
            super(message);
        }

        public XmlMigrationException(String message, Throwable cause) {
            super(message, cause);
        }
    }
}

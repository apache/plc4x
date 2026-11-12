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

import org.w3c.dom.Element;

import java.io.IOException;
import java.io.InputStream;
import java.nio.file.Files;
import java.nio.file.Path;

/**
 * Root of a GSDML device description (the {@code ISO15745Profile} element).
 */
public class ProfinetISO15745Profile {

    private ProfinetProfileHeader profileHeader;

    private ProfinetProfileBody profileBody;

    /**
     * Parses a GSDML file.
     *
     * @throws IOException if the file cannot be read, is not well-formed XML (or declares a DTD), or
     *                     holds a value that does not fit the GSDML model
     */
    public static ProfinetISO15745Profile parse(Path file) throws IOException {
        try (InputStream in = Files.newInputStream(file)) {
            return parse(in);
        }
    }

    /**
     * Parses a GSDML document.
     *
     * @throws IOException if the stream cannot be read, is not well-formed XML (or declares a DTD), or
     *                     holds a value that does not fit the GSDML model
     */
    public static ProfinetISO15745Profile parse(InputStream in) throws IOException {
        return fromXml(GsdmlXml.parse(in));
    }

    public ProfinetProfileHeader getProfileHeader() {
        return profileHeader;
    }

    public ProfinetProfileBody getProfileBody() {
        return profileBody;
    }

    static ProfinetISO15745Profile fromXml(Element element) throws IOException {
        ProfinetISO15745Profile result = new ProfinetISO15745Profile();
        for (GsdmlXml.Property property : GsdmlXml.properties(element, false)) {
            switch (property.name()) {
                case "ProfileHeader" -> result.profileHeader = property.object(ProfinetProfileHeader::fromXml);
                case "ProfileBody" -> result.profileBody = property.object(ProfinetProfileBody::fromXml);
                default -> {
                    // Not used by the driver.
                }
            }
        }
        return result;
    }
}

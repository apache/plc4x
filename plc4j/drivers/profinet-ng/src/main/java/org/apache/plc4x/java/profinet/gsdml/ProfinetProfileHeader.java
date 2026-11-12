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

public class ProfinetProfileHeader {

    private String profileIdentification;

    private String profileRevision;

    private String profileName;

    private String profileSource;

    private String profileClassID;

    private ProfinetISO15745Reference iso15745Reference;

    public String getProfileIdentification() {
        return profileIdentification;
    }

    public String getProfileRevision() {
        return profileRevision;
    }

    public String getProfileName() {
        return profileName;
    }

    public String getProfileSource() {
        return profileSource;
    }

    public String getProfileClassID() {
        return profileClassID;
    }

    public ProfinetISO15745Reference getIso15745Reference() {
        return iso15745Reference;
    }

    static ProfinetProfileHeader fromXml(Element element) throws IOException {
        ProfinetProfileHeader result = new ProfinetProfileHeader();
        for (GsdmlXml.Property property : GsdmlXml.properties(element, false)) {
            switch (property.name()) {
                case "ProfileIdentification" -> result.profileIdentification = property.string();
                case "ProfileRevision" -> result.profileRevision = property.string();
                case "ProfileName" -> result.profileName = property.string();
                case "ProfileSource" -> result.profileSource = property.string();
                case "ProfileClassID" -> result.profileClassID = property.string();
                case "ISO15745Reference" -> result.iso15745Reference = property.object(ProfinetISO15745Reference::fromXml);
                default -> {
                    // Not used by the driver.
                }
            }
        }
        return result;
    }
}

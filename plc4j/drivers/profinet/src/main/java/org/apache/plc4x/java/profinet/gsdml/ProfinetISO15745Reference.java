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

public class ProfinetISO15745Reference {

    private int iso15745Part;

    private int iso15745Edition;

    private String profileTechnology;

    public int getIso15745Part() {
        return iso15745Part;
    }

    public int getIso15745Edition() {
        return iso15745Edition;
    }

    public String getProfileTechnology() {
        return profileTechnology;
    }

    static ProfinetISO15745Reference fromXml(Element element) throws IOException {
        ProfinetISO15745Reference result = new ProfinetISO15745Reference();
        for (GsdmlXml.Property property : GsdmlXml.properties(element, false)) {
            switch (property.name()) {
                case "ISO15745Part" -> result.iso15745Part = property.intValue();
                case "ISO15745Edition" -> result.iso15745Edition = property.intValue();
                case "ProfileTechnology" -> result.profileTechnology = property.string();
                default -> {
                    // Not used by the driver.
                }
            }
        }
        return result;
    }
}

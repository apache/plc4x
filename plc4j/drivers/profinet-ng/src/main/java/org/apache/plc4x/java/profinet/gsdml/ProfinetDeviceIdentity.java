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

public class ProfinetDeviceIdentity {


    private String vendorId;

    private String deviceID;

    private ProfinetTextId infoText;

    private ProfinetValue vendorName;

    public String getVendorId() {
        return vendorId;
    }

    public String getDeviceID() {
        return deviceID;
    }

    public ProfinetTextId getInfoText() {
        return infoText;
    }

    public ProfinetValue getVendorName() {
        return vendorName;
    }

    static ProfinetDeviceIdentity fromXml(Element element) throws IOException {
        ProfinetDeviceIdentity result = new ProfinetDeviceIdentity();
        for (GsdmlXml.Property property : GsdmlXml.properties(element, false)) {
            switch (property.name()) {
                case "VendorID" -> result.vendorId = property.string();
                case "DeviceID" -> result.deviceID = property.string();
                case "InfoText" -> result.infoText = property.object(ProfinetTextId::fromXml);
                case "VendorName" -> result.vendorName = property.object(ProfinetValue::fromXml);
                default -> {
                    // Not used by the driver.
                }
            }
        }
        return result;
    }
}

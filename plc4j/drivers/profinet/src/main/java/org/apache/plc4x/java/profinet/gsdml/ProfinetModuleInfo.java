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

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonRootName;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonRootName("DeviceAccessPointList")
public class ProfinetModuleInfo {

    @JacksonXmlProperty(localName="Name")
    private ProfinetTextId name;

    @JacksonXmlProperty(localName="InfoText")
    private ProfinetTextId infoText;

    @JacksonXmlProperty(localName="VendorName")
    private ProfinetValue vendorName;

    @JacksonXmlProperty(localName="OrderNumber")
    private ProfinetValue orderNumber;

    @JacksonXmlProperty(localName="HardwareRelease")
    private ProfinetValue hardwareRelease;

    @JacksonXmlProperty(localName="SoftwareRelease")
    private ProfinetValue softwareRelease;

    public ProfinetTextId getName() {
        return name;
    }

    public ProfinetTextId getInfoText() {
        return infoText;
    }

    public ProfinetValue getVendorName() {
        return vendorName;
    }

    public ProfinetValue getOrderNumber() {
        return orderNumber;
    }

    public ProfinetValue getHardwareRelease() {
        return hardwareRelease;
    }

    public ProfinetValue getSoftwareRelease() {
        return softwareRelease;
    }
}

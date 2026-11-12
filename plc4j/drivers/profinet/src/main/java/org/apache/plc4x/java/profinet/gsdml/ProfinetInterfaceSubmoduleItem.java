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

public class ProfinetInterfaceSubmoduleItem {

    private String id;

    private String submoduleIdentNumber;

    private int subslotNumber;

    private String textId;

    private String supportedRtClasses;

    private String supportedProtocols;

    private boolean networkComponentDiagnosisSupported;

    private boolean ptpBoundarySupported;

    private boolean dcpBoundarySupported;

    private ProfinetApplicationRelations applicationRelations;

    public String getId() {
        return id;
    }

    public String getSubmoduleIdentNumber() {
        return submoduleIdentNumber;
    }

    public int getSubslotNumber() {
        return subslotNumber;
    }

    public String getTextId() {
        return textId;
    }

    public String getSupportedRtClasses() {
        return supportedRtClasses;
    }

    public String getSupportedProtocols() {
        return supportedProtocols;
    }

    public boolean isNetworkComponentDiagnosisSupported() {
        return networkComponentDiagnosisSupported;
    }

    public boolean isPtpBoundarySupported() {
        return ptpBoundarySupported;
    }

    public boolean isDcpBoundarySupported() {
        return dcpBoundarySupported;
    }

    public ProfinetApplicationRelations getApplicationRelations() {
        return applicationRelations;
    }

    static ProfinetInterfaceSubmoduleItem fromXml(Element element) throws IOException {
        ProfinetInterfaceSubmoduleItem result = new ProfinetInterfaceSubmoduleItem();
        for (GsdmlXml.Property property : GsdmlXml.properties(element, false)) {
            switch (property.name()) {
                case "ID" -> result.id = property.string();
                case "SubmoduleIdentNumber" -> result.submoduleIdentNumber = property.string();
                case "SubslotNumber" -> result.subslotNumber = property.intValue();
                case "TextId" -> result.textId = property.string();
                case "SupportedRT_Classes" -> result.supportedRtClasses = property.string();
                case "SupportedProtocols" -> result.supportedProtocols = property.string();
                case "NetworkComponentDiagnosisSupported" -> result.networkComponentDiagnosisSupported = property.booleanValue();
                case "PTP_BoundarySupported" -> result.ptpBoundarySupported = property.booleanValue();
                case "DCP_BoundarySupported" -> result.dcpBoundarySupported = property.booleanValue();
                case "ApplicationRelations" -> result.applicationRelations = property.object(ProfinetApplicationRelations::fromXml);
                default -> {
                    // Not used by the driver.
                }
            }
        }
        return result;
    }
}


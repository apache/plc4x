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
import java.util.List;

public class ProfinetDeviceAccessPointItem implements ProfinetDeviceItem {

    private String id;

    private String pnioVersion;

    private String physicalSlots;

    private String moduleIdentNumber;

    private int minDeviceInterval;

    private String dnsCompatibleName;

    private int fixedInSlots;

    private int objectUUIDLocalIndex;

    private boolean deviceAccessSupported;

    private boolean multipleWriteSupported;

    private boolean checkDeviceIDAllowed;

    private boolean nameOfStationNotTransferable;

    private boolean lldpNodSupported;

    private String resetToFactoryModes;

    private ProfinetModuleInfo moduleInfo;

    private ProfinetCertificationInfo certificationInfo;

    private ProfinetIOConfigData ioConfigData;

    private List<ProfinetModuleItemRef> useableModules;

    private List<ProfinetVirtualSubmoduleItem> virtualSubmoduleList;

    private ProfinetSystemDefinedSubmoduleList systemDefinedSubmoduleList;

    private ProfinetGraphics graphics;

    public String getId() {
        return id;
    }

    public String getPnioVersion() {
        return pnioVersion;
    }

    public String getPhysicalSlots() {
        return physicalSlots;
    }

    public String getModuleIdentNumber() {
        return moduleIdentNumber;
    }

    @Override
    public Integer getInputDataLength() {
        return 0;
    }

    @Override
    public Integer getOutputDataLength() {
        return 0;
    }

    public int getMinDeviceInterval() {
        return minDeviceInterval;
    }

    public String getDnsCompatibleName() {
        return dnsCompatibleName;
    }

    public int getFixedInSlots() {
        return fixedInSlots;
    }

    public int getObjectUUIDLocalIndex() {
        return objectUUIDLocalIndex;
    }

    public boolean isDeviceAccessSupported() {
        return deviceAccessSupported;
    }

    public boolean isMultipleWriteSupported() {
        return multipleWriteSupported;
    }

    public boolean isCheckDeviceIDAllowed() {
        return checkDeviceIDAllowed;
    }

    public boolean isNameOfStationNotTransferable() {
        return nameOfStationNotTransferable;
    }

    public boolean isLldpNodSupported() {
        return lldpNodSupported;
    }

    public String getResetToFactoryModes() {
        return resetToFactoryModes;
    }

    public ProfinetModuleInfo getModuleInfo() {
        return moduleInfo;
    }

    public ProfinetCertificationInfo getCertificationInfo() {
        return certificationInfo;
    }

    public ProfinetIOConfigData getIoConfigData() {
        return ioConfigData;
    }

    public List<ProfinetModuleItemRef> getUseableModules() {
        return useableModules;
    }

    public List<ProfinetVirtualSubmoduleItem> getVirtualSubmoduleList() {
        return virtualSubmoduleList;
    }

    public ProfinetSystemDefinedSubmoduleList getSystemDefinedSubmoduleList() {
        return systemDefinedSubmoduleList;
    }

    public ProfinetGraphics getGraphics() {
        return graphics;
    }

    static ProfinetDeviceAccessPointItem fromXml(Element element) throws IOException {
        ProfinetDeviceAccessPointItem result = new ProfinetDeviceAccessPointItem();
        for (GsdmlXml.Property property : GsdmlXml.properties(element, false)) {
            switch (property.name()) {
                case "ID" -> result.id = property.string();
                case "PNIO_Version" -> result.pnioVersion = property.string();
                case "PhysicalSlots" -> result.physicalSlots = property.string();
                case "ModuleIdentNumber" -> result.moduleIdentNumber = property.string();
                case "MinDeviceInterval" -> result.minDeviceInterval = property.intValue();
                case "DNS_CompatibleName" -> result.dnsCompatibleName = property.string();
                case "FixedInSlots" -> result.fixedInSlots = property.intValue();
                case "ObjectUUID_LocalIndex" -> result.objectUUIDLocalIndex = property.intValue();
                case "DeviceAccessSupported" -> result.deviceAccessSupported = property.booleanValue();
                case "MultipleWriteSupported" -> result.multipleWriteSupported = property.booleanValue();
                case "CheckDeviceID_Allowed" -> result.checkDeviceIDAllowed = property.booleanValue();
                case "NameOfStationNotTransferable" -> result.nameOfStationNotTransferable = property.booleanValue();
                case "LLDP_NoD_Supported" -> result.lldpNodSupported = property.booleanValue();
                case "ResetToFactoryModes" -> result.resetToFactoryModes = property.string();
                case "ModuleInfo" -> result.moduleInfo = property.object(ProfinetModuleInfo::fromXml);
                case "CertificationInfo" -> result.certificationInfo = property.object(ProfinetCertificationInfo::fromXml);
                case "IOConfigData" -> result.ioConfigData = property.object(ProfinetIOConfigData::fromXml);
                case "UseableModules" -> result.useableModules = property.wrappedList(ProfinetModuleItemRef::fromXml);
                case "VirtualSubmoduleList" -> result.virtualSubmoduleList = property.wrappedList(ProfinetVirtualSubmoduleItem::fromXml);
                case "SystemDefinedSubmoduleList" -> result.systemDefinedSubmoduleList = property.object(ProfinetSystemDefinedSubmoduleList::fromXml);
                case "Graphics" -> result.graphics = property.object(ProfinetGraphics::fromXml);
                default -> {
                    // Not used by the driver.
                }
            }
        }
        return result;
    }
}

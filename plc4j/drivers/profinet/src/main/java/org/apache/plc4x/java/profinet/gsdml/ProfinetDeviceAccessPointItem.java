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

import java.util.List;

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonRootName("DeviceAccessPointItem")
public class ProfinetDeviceAccessPointItem implements ProfinetDeviceItem {

    @JacksonXmlProperty(isAttribute=true, localName="ID")
    private String id;

    @JacksonXmlProperty(isAttribute=true, localName="PNIO_Version")
    private String pnioVersion;

    @JacksonXmlProperty(isAttribute=true, localName="PhysicalSlots")
    private String physicalSlots;

    @JacksonXmlProperty(isAttribute=true, localName="ModuleIdentNumber")
    private String moduleIdentNumber;

    @JacksonXmlProperty(isAttribute=true, localName="MinDeviceInterval")
    private int minDeviceInterval;

    @JacksonXmlProperty(isAttribute=true, localName="DNS_CompatibleName")
    private String dnsCompatibleName;

    @JacksonXmlProperty(isAttribute=true, localName="FixedInSlots")
    private int fixedInSlots;

    @JacksonXmlProperty(isAttribute=true, localName="ObjectUUID_LocalIndex")
    private int objectUUIDLocalIndex;

    @JacksonXmlProperty(isAttribute=true, localName="DeviceAccessSupported")
    private boolean deviceAccessSupported;

    @JacksonXmlProperty(isAttribute=true, localName="MultipleWriteSupported")
    private boolean multipleWriteSupported;

    @JacksonXmlProperty(isAttribute=true, localName="CheckDeviceID_Allowed")
    private boolean checkDeviceIDAllowed;

    @JacksonXmlProperty(isAttribute=true, localName="NameOfStationNotTransferable")
    private boolean nameOfStationNotTransferable;

    @JacksonXmlProperty(isAttribute=true, localName="LLDP_NoD_Supported")
    private boolean lldpNodSupported;

    @JacksonXmlProperty(isAttribute=true, localName="ResetToFactoryModes")
    private String resetToFactoryModes;

    @JacksonXmlProperty(localName="ModuleInfo")
    private ProfinetModuleInfo moduleInfo;

    @JacksonXmlProperty(localName="CertificationInfo")
    private ProfinetCertificationInfo certificationInfo;

    @JacksonXmlProperty(localName="IOConfigData")
    private ProfinetIOConfigData ioConfigData;

    @JacksonXmlProperty(localName="UseableModules")
    private List<ProfinetModuleItemRef> useableModules;

    @JacksonXmlProperty(localName="VirtualSubmoduleList")
    private List<ProfinetVirtualSubmoduleItem> virtualSubmoduleList;

    @JacksonXmlProperty(localName="SystemDefinedSubmoduleList")
    private ProfinetSystemDefinedSubmoduleList systemDefinedSubmoduleList;

    @JacksonXmlProperty(localName="Graphics")
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
}

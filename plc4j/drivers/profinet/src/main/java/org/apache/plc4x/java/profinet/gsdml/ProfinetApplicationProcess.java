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

import com.fasterxml.jackson.annotation.JsonGetter;
import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonRootName;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlProperty;

import java.util.List;

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonRootName("ApplicationProcess")
public class ProfinetApplicationProcess {

    @JacksonXmlProperty(localName="DeviceAccessPointList")
    private List<ProfinetDeviceAccessPointItem> deviceAccessPointList;

    @JacksonXmlProperty(localName="ModuleList")
    private List<ProfinetModuleItem> moduleList;

    @JacksonXmlProperty(localName="LogBookEntryList")
    private List<ProfinetLogBookEntryItem> logBookEntryList;

    @JacksonXmlProperty(localName="GraphicsList")
    private List<ProfinetGraphicItem> graphicsList;

    @JacksonXmlProperty(localName="ExternalTextList")
    private ProfinetExternalTextList externalTextList;

    public List<ProfinetDeviceAccessPointItem> getDeviceAccessPointList() {
        return deviceAccessPointList;
    }

    public List<ProfinetModuleItem> getModuleList() {
        return moduleList;
    }

    public List<ProfinetLogBookEntryItem> getLogBookEntryList() {
        return logBookEntryList;
    }

    public List<ProfinetGraphicItem> getGraphicsList() {
        return graphicsList;
    }

    public ProfinetExternalTextList getExternalTextList() {
        return externalTextList;
    }
}

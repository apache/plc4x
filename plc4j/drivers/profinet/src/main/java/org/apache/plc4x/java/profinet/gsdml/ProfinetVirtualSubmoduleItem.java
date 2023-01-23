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
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlElementWrapper;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlProperty;

import java.util.List;
import java.util.Map;

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonRootName("VirtualSubmoduleItem")
public class ProfinetVirtualSubmoduleItem {

    @JacksonXmlProperty(isAttribute=true, localName="ID")
    private String id;

    @JacksonXmlProperty(isAttribute=true, localName="SubmoduleIdentNumber")
    private String submoduleIdentNumber;

    @JacksonXmlProperty(isAttribute=true, localName="SubslotNumber")
    private int subslotNumber = 1;

    @JacksonXmlProperty(isAttribute=true, localName="Writeable_IM_Records")
    private String writeableImRecords;

    @JacksonXmlProperty(isAttribute=true, localName="MayIssueProcessAlarm")
    private boolean mayIssueProcessAlarm;

    @JacksonXmlProperty(localName="IOData")
    private ProfinetIoData ioData;

    @JacksonXmlProperty(localName="ModuleInfo")
    private ProfinetModuleInfo moduleInfo;

    @JacksonXmlProperty(localName="RecordDataList")
    private List<ProfinetParameterRecordDataItem> recordDataList;

    public String getId() {
        return id;
    }

    public String getSubmoduleIdentNumber() {
        return submoduleIdentNumber;
    }

    public String getWriteableImRecords() {
        return writeableImRecords;
    }

    public boolean isMayIssueProcessAlarm() {
        return mayIssueProcessAlarm;
    }

    public ProfinetIoData getIoData() {
        return ioData;
    }

    public ProfinetModuleInfo getModuleInfo() {
        return moduleInfo;
    }

    public List<ProfinetParameterRecordDataItem> getRecordDataList() {
        return recordDataList;
    }

    public int getSubslotNumber() {
        return subslotNumber;
    }
}

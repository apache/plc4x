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
}

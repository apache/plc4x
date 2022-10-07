package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonRootName;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlProperty;

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

}

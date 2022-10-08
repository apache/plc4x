package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonRootName;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlProperty;

import java.util.List;

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonRootName("LogBookEntryItem")
public class ProfinetLogBookEntryItem {

    @JacksonXmlProperty(isAttribute=true, localName="Status")
    private String status;

    @JacksonXmlProperty(localName="ErrorCode2Value")
    private ProfinetErrorCode2Value errorCode2Value;

    public String getStatus() {
        return status;
    }

    public ProfinetErrorCode2Value getErrorCode2Value() {
        return errorCode2Value;
    }
}



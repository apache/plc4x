package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonRootName;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlProperty;

import java.util.List;

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonRootName("ErrorCode2Value")
public class ProfinetErrorCode2Value {

    @JacksonXmlProperty(isAttribute=true, localName="Name")
    private ProfinetTextId name;

    public ProfinetTextId getName() {
        return name;
    }
}



package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonRootName;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonRootName("IOConfigData")
public class ProfinetIOConfigData {

    @JacksonXmlProperty(isAttribute=true, localName="MaxInputLength")
    private int maxInputLength;

    @JacksonXmlProperty(isAttribute=true, localName="MaxOutputLength")
    private int maxOutputLength;

}

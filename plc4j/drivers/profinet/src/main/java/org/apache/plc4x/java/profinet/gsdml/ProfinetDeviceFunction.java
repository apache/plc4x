package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonRootName;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonRootName("DeviceFunction")
public class ProfinetDeviceFunction {

    @JacksonXmlProperty(localName="Family")
    private ProfinetFamily family;

    public ProfinetFamily getFamily() {
        return family;
    }
}

package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonRootName;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonRootName("ApplicationRelations")
public class ProfinetApplicationRelations {

    @JacksonXmlProperty(isAttribute=true, localName="StartupMode")
    private String startupMode;

    @JacksonXmlProperty(localName="TimingProperties")
    private ProfinetTimingProperties TimingProperties;

    public String getStartupMode() {
        return startupMode;
    }

    public ProfinetTimingProperties getTimingProperties() {
        return TimingProperties;
    }
}

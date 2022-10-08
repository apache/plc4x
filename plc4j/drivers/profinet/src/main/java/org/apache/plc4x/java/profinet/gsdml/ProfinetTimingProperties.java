package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonRootName;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonRootName("TimingProperties")
public class ProfinetTimingProperties {

    @JacksonXmlProperty(isAttribute=true, localName="SendClock")
    private String sendClock;

    @JacksonXmlProperty(isAttribute=true, localName="ReductionRatio")
    private String reductionRatio;

    public String getSendClock() {
        return sendClock;
    }

    public String getReductionRatio() {
        return reductionRatio;
    }
}

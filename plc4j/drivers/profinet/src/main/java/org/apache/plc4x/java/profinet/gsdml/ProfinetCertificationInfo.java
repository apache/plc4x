package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonRootName;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonRootName("CertificationInfo")
public class ProfinetCertificationInfo {

    @JacksonXmlProperty(isAttribute=true, localName="ConformanceClass")
    private String conformanceClass;

    @JacksonXmlProperty(isAttribute=true, localName="ApplicationClass")
    private String applicationClass;

    @JacksonXmlProperty(isAttribute=true, localName="NetloadClass")
    private String netloadClass;

    public String getConformanceClass() {
        return conformanceClass;
    }

    public String getApplicationClass() {
        return applicationClass;
    }

    public String getNetloadClass() {
        return netloadClass;
    }
}

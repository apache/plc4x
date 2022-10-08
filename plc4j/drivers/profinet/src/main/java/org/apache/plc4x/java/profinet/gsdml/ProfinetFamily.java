package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonRootName;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonRootName("Family")
public class ProfinetFamily {

    @JacksonXmlProperty(isAttribute=true, localName="MainFamily")
    private String mainFamily;

    @JacksonXmlProperty(isAttribute=true, localName="ProductFamily")
    private String productFamily;

    public String getMainFamily() {
        return mainFamily;
    }

    public String getProductFamily() {
        return productFamily;
    }
}

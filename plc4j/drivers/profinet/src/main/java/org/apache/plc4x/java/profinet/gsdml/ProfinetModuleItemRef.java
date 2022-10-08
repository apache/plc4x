package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonRootName;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonRootName("ModuleItemRef")
public class ProfinetModuleItemRef {

    @JacksonXmlProperty(isAttribute=true, localName="ModuleItemTarget")
    private String moduleItemTarget;

    @JacksonXmlProperty(isAttribute=true, localName="AllowedInSlots")
    private String allowedInSlots;

    public String getModuleItemTarget() {
        return moduleItemTarget;
    }

    public String getAllowedInSlots() {
        return allowedInSlots;
    }
}

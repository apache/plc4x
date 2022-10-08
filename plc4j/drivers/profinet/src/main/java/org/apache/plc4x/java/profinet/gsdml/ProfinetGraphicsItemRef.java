package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonRootName;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonRootName("GraphicsItemRef")
public class ProfinetGraphicsItemRef {

    @JacksonXmlProperty(isAttribute=true, localName="Type")
    private String type;

    @JacksonXmlProperty(isAttribute=true, localName="GraphicItemTarget")
    private String graphicItemTarget;

    public String getType() {
        return type;
    }

    public String getGraphicItemTarget() {
        return graphicItemTarget;
    }
}

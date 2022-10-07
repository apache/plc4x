package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonRootName;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonRootName("Graphics")
public class ProfinetGraphics {

    @JacksonXmlProperty(localName="GraphicItemRef")
    private ProfinetGraphicsItemRef GraphicItemRef;

}

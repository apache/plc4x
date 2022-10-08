package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonRootName;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonRootName("GraphicsItem")
public class ProfinetGraphicItem {

    @JacksonXmlProperty(isAttribute=true, localName="ID")
    private String id;

    @JacksonXmlProperty(isAttribute=true, localName="GraphicFile")
    private String graphicFile;

}

package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonRootName;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonRootName("InfoText")
public class ProfinetTextIdValue {

    @JacksonXmlProperty(isAttribute=true, localName="TextId")
    private String textId;

    @JacksonXmlProperty(isAttribute=true, localName="Value")
    private String value;

}

package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonRootName;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonRootName("Ref")
public class ProfinetRef {

    @JacksonXmlProperty(isAttribute=true, localName="DataType")
    private String dataType;

    @JacksonXmlProperty(isAttribute=true, localName="ByteOffset")
    private int byteOffset;

    @JacksonXmlProperty(isAttribute=true, localName="DefaultValue")
    private String defaultValue;

    @JacksonXmlProperty(isAttribute=true, localName="AllowedValues")
    private String allowedValues;

    @JacksonXmlProperty(isAttribute=true, localName="Changeable")
    private boolean changeable;

    @JacksonXmlProperty(isAttribute=true, localName="Visible")
    private boolean visible;

    @JacksonXmlProperty(isAttribute=true, localName="TextId")
    private String textId;


}

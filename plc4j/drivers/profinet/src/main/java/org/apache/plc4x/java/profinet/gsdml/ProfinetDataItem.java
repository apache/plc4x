package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonRootName;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlProperty;

import java.util.List;

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonRootName("CertificationInfo")
public class ProfinetDataItem {

    @JacksonXmlProperty(isAttribute=true, localName="DataType")
    private String dataType;

    @JacksonXmlProperty(isAttribute=true, localName="TextId")
    private String textId;

    @JacksonXmlProperty(isAttribute=true, localName="UseAsBits")
    private boolean UseAsBits;

    @JacksonXmlProperty(isAttribute=true, localName="BitDataItem")
    private List<ProfinetBitDataItem> bitDataItem;

}



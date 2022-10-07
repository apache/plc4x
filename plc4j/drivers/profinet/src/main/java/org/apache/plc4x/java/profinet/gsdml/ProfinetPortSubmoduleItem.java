package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonRootName;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlProperty;

import java.util.List;

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonRootName("PortSubmoduleItem")
public class ProfinetPortSubmoduleItem extends ProfinetSystemDefinedSubmoduleItem{

    @JacksonXmlProperty(isAttribute=true, localName="ID")
    private String id;

    @JacksonXmlProperty(isAttribute=true, localName="SubmoduleIdentNumber")
    private String submoduleIdentNumber;

    @JacksonXmlProperty(isAttribute=true, localName="SubslotNumber")
    private int subslotNumber;

    @JacksonXmlProperty(isAttribute=true, localName="TextId")
    private String textId;

    @JacksonXmlProperty(isAttribute=true, localName="MaxPortRxDelay")
    private int maxPortRxDelay;

    @JacksonXmlProperty(isAttribute=true, localName="MaxPortTxDelay")
    private int maxPortTxDelay;

    @JacksonXmlProperty(localName="MAUTypeList")
    private List<ProfinetValue> mauTypeItem;

}


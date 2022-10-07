package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonRootName;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlProperty;

import java.util.List;

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonRootName("Input")
public class ProfinetIoDataInput extends ProfinetIoDataItem {

    @JacksonXmlProperty(isAttribute=true, localName="Consistency")
    private String consistency;

    @JacksonXmlProperty(localName="DataItem")
    private List<ProfinetDataItem> dataItemList;

}

package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.annotation.*;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlElementWrapper;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlProperty;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlRootElement;

import java.util.List;

@JsonTypeName("Output")
public class ProfinetIoDataOutput {

    @JacksonXmlProperty(isAttribute=true, localName="Consistency")
    private String consistency;

    @JacksonXmlProperty(localName="DataItem")
    @JacksonXmlElementWrapper(useWrapping = false)
    private List<ProfinetDataItem> dataItemList;

    public String getConsistency() {
        return consistency;
    }

    public List<ProfinetDataItem> getDataItemList() {
        return dataItemList;
    }
}

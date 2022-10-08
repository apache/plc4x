package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonRootName;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlProperty;

import java.util.List;

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonRootName("ParameterRecordDataItem")
public class ProfinetParameterRecordDataItem {

    @JacksonXmlProperty(isAttribute=true, localName="Index")
    private int index;

    @JacksonXmlProperty(isAttribute=true, localName="Length")
    private int length;

    @JacksonXmlProperty(localName="Name")
    private ProfinetTextId name;

    @JacksonXmlProperty(localName="Ref")
    private ProfinetRef ref;

    public int getIndex() {
        return index;
    }

    public int getLength() {
        return length;
    }

    public ProfinetTextId getName() {
        return name;
    }

    public ProfinetRef getRef() {
        return ref;
    }
}

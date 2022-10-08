package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.annotation.JsonAlias;
import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonRootName;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlElementWrapper;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlProperty;

import java.util.List;


public class ProfinetIoData {

    @JacksonXmlProperty(localName="Input")
    @JacksonXmlElementWrapper(useWrapping = false)
    private List<ProfinetIoDataInput> input;

    @JacksonXmlProperty(localName="Output")
    @JacksonXmlElementWrapper(useWrapping = false)
    private List<ProfinetIoDataOutput> output;

    public List<ProfinetIoDataInput> getInput() {
        return input;
    }

    public List<ProfinetIoDataOutput> getOutput() {
        return output;
    }
}

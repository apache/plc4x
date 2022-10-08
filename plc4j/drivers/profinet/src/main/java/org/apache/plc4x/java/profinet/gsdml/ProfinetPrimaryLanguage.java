package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonRootName;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlElementWrapper;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlProperty;

import java.util.List;

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonRootName("PrimaryLanguage")
public class ProfinetPrimaryLanguage {

    @JacksonXmlProperty(localName="Text")
    @JacksonXmlElementWrapper(useWrapping = false)
    private List<ProfinetTextIdValue> text;

    public List<ProfinetTextIdValue> getText() {
        return text;
    }
}

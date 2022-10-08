package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.annotation.*;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlProperty;

@JsonTypeInfo(use= JsonTypeInfo.Id.CLASS)
@JsonSubTypes({
    @JsonSubTypes.Type(value = ProfinetIoDataInput.class, name = "Input"),
    @JsonSubTypes.Type(value = ProfinetIoDataOutput.class, name = "Output")
})
public abstract class ProfinetIoDataItem {

}

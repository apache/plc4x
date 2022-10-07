package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.annotation.JsonGetter;
import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonRootName;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlProperty;

import java.util.List;

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonRootName("ApplicationProcess")
public class ProfinetApplicationProcess {

    @JacksonXmlProperty(localName="DeviceAccessPointList")
    private List<ProfinetDeviceAccessPointItem> DeviceAccessPointList;
}

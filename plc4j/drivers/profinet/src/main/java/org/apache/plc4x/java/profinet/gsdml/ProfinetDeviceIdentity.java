package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonRootName;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonRootName("DeviceIdentity")
public class ProfinetDeviceIdentity {


    @JacksonXmlProperty(isAttribute=true, localName="ModuleItemTarget")
    private String ModuleItemTarget;

    @JacksonXmlProperty(isAttribute=true, localName="DeviceID")
    private String deviceID;

    @JacksonXmlProperty(localName="InfoText")
    private ProfinetTextId infoText;

    @JacksonXmlProperty(localName="VendorName")
    private ProfinetValue vendorName;

}

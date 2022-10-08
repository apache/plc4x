package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonRootName;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonRootName("DeviceAccessPointList")
public class ProfinetModuleInfo {

    @JacksonXmlProperty(localName="Name")
    private ProfinetTextId name;

    @JacksonXmlProperty(localName="InfoText")
    private ProfinetTextId infoText;

    @JacksonXmlProperty(localName="VendorName")
    private ProfinetValue vendorName;

    @JacksonXmlProperty(localName="OrderNumber")
    private ProfinetValue orderNumber;

    @JacksonXmlProperty(localName="HardwareRelease")
    private ProfinetValue hardwareRelease;

    @JacksonXmlProperty(localName="SoftwareRelease")
    private ProfinetValue softwareRelease;

    public ProfinetTextId getName() {
        return name;
    }

    public ProfinetTextId getInfoText() {
        return infoText;
    }

    public ProfinetValue getVendorName() {
        return vendorName;
    }

    public ProfinetValue getOrderNumber() {
        return orderNumber;
    }

    public ProfinetValue getHardwareRelease() {
        return hardwareRelease;
    }

    public ProfinetValue getSoftwareRelease() {
        return softwareRelease;
    }
}

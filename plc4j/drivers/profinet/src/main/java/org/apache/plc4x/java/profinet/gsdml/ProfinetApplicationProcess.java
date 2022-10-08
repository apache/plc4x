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
    private List<ProfinetDeviceAccessPointItem> deviceAccessPointList;

    @JacksonXmlProperty(localName="ModuleList")
    private List<ProfinetModuleItem> moduleList;

    @JacksonXmlProperty(localName="LogBookEntryList")
    private List<ProfinetLogBookEntryItem> logBookEntryList;

    @JacksonXmlProperty(localName="GraphicsList")
    private List<ProfinetGraphicItem> graphicsList;

    @JacksonXmlProperty(localName="ExternalTextList")
    private ProfinetExternalTextList externalTextList;

    public List<ProfinetDeviceAccessPointItem> getDeviceAccessPointList() {
        return deviceAccessPointList;
    }

    public List<ProfinetModuleItem> getModuleList() {
        return moduleList;
    }

    public List<ProfinetLogBookEntryItem> getLogBookEntryList() {
        return logBookEntryList;
    }

    public List<ProfinetGraphicItem> getGraphicsList() {
        return graphicsList;
    }

    public ProfinetExternalTextList getExternalTextList() {
        return externalTextList;
    }
}

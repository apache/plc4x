package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonRootName;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonRootName("InterfaceSubmoduleItem")
public class ProfinetInterfaceSubmoduleItem extends ProfinetSystemDefinedSubmoduleItem{

    @JacksonXmlProperty(isAttribute=true, localName="ID")
    private String id;

    @JacksonXmlProperty(isAttribute=true, localName="SubmoduleIdentNumber")
    private String submoduleIdentNumber;

    @JacksonXmlProperty(isAttribute=true, localName="SubslotNumber")
    private int subslotNumber;

    @JacksonXmlProperty(isAttribute=true, localName="TextId")
    private String textId;

    @JacksonXmlProperty(isAttribute=true, localName="SupportedRT_Classes")
    private String supportedRtClasses;

    @JacksonXmlProperty(isAttribute=true, localName="SupportedProtocols")
    private String supportedProtocols;

    @JacksonXmlProperty(isAttribute=true, localName="NetworkComponentDiagnosisSupported")
    private boolean networkComponentDiagnosisSupported;

    @JacksonXmlProperty(isAttribute=true, localName="PTP_BoundarySupported")
    private boolean ptpBoundarySupported;

    @JacksonXmlProperty(isAttribute=true, localName="DCP_BoundarySupported")
    private boolean dcpBoundarySupported;

    @JacksonXmlProperty(localName="ApplicationRelations")
    private ProfinetApplicationRelations applicationRelations;

}


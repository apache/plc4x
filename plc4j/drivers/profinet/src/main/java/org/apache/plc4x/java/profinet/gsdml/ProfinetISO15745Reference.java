package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonRootName;

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonRootName("ISO15745Reference")
public class ProfinetISO15745Reference {

    @JsonProperty("ISO15745Part")
    private int iso15745Part;

    @JsonProperty("ISO15745Edition")
    private int iso15745Edition;

    @JsonProperty("ProfileTechnology")
    private String profileTechnology;

}

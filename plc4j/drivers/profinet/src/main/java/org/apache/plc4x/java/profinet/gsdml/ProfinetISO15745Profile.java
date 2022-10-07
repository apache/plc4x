package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonRootName;

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonRootName("ISO15745Profile")
public class ProfinetISO15745Profile {

    @JsonProperty("ProfileHeader")
    private ProfinetProfileHeader profileHeader;

    @JsonProperty("ProfileBody")
    private ProfinetProfileBody profileBody;
}

package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonRootName;

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonRootName("ProfileHeader")
public class ProfinetProfileHeader {

    @JsonProperty("ProfileIdentification")
    private String profileIdentification;

    @JsonProperty("ProfileRevision")
    private String profileRevision;

    @JsonProperty("ProfileName")
    private String profileName;

    @JsonProperty("ProfileSource")
    private String profileSource;

    @JsonProperty("ProfileClassID")
    private String profileClassID;

    @JsonProperty("ISO15745Reference")
    private ProfinetISO15745Reference iso15745Reference;

    public String getProfileIdentification() {
        return profileIdentification;
    }

    public String getProfileRevision() {
        return profileRevision;
    }

    public String getProfileName() {
        return profileName;
    }

    public String getProfileSource() {
        return profileSource;
    }

    public String getProfileClassID() {
        return profileClassID;
    }

    public ProfinetISO15745Reference getIso15745Reference() {
        return iso15745Reference;
    }
}

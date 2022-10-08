package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonRootName;

@JsonRootName("ProfileBody")
public class ProfinetProfileBody {

    @JsonProperty("DeviceIdentity")
    private ProfinetDeviceIdentity deviceIdentity;

    @JsonProperty("DeviceFunction")
    private ProfinetDeviceFunction deviceFunction;

    @JsonProperty("ApplicationProcess")
    private ProfinetApplicationProcess applicationProcess;

    public ProfinetDeviceIdentity getDeviceIdentity() {
        return deviceIdentity;
    }

    public ProfinetDeviceFunction getDeviceFunction() {
        return deviceFunction;
    }

    public ProfinetApplicationProcess getApplicationProcess() {
        return applicationProcess;
    }
}

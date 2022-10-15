package org.apache.plc4x.java.spi.configuration;

import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.StringDefaultValue;

public class BaseConfiguration implements Configuration {

    @ConfigurationParameter("transportConfig")
    @StringDefaultValue("")
    public String transportConfig;

    public String getTransportConfig() {
        return transportConfig;
    }

    public void setTransportConfig(String transportConfig) {
        this.transportConfig = transportConfig;
    }
}

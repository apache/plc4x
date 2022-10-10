package org.apache.plc4x.java.profinet.gsdml;

import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.profinet.config.ProfinetConfiguration;
import org.apache.plc4x.java.profinet.context.ProfinetDriverContext;
import org.apache.plc4x.java.profinet.device.ProfinetDevice;
import org.apache.plc4x.java.profinet.readwrite.MacAddress;
import org.apache.plc4x.java.spi.configuration.ConfigurationFactory;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;

import java.util.Map;

import static org.junit.jupiter.api.Assertions.assertEquals;

@TestInstance(TestInstance.Lifecycle.PER_CLASS)
public class ProfinetConfigurationTests {

    /*
        Profinet GSD File Directory Configuration Test
     */
    @Test
    public void readGsdDirectory()  {

        String directory = "/home/plc4x/gsd_directory";
        ProfinetConfiguration configuration = (ProfinetConfiguration) new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "gsddirectory=" + directory);

        assertEquals(directory, configuration.getGsdDirectory());
    }

    /*
        Profinet GSD File Directory Configuration Test
     */
    @Test
    public void readProfinetDevices() throws DecoderException, PlcConnectionException {

        String[] macAddresses = new String[] {"CA:FE:00:00:00:01"};
        ProfinetConfiguration configuration = (ProfinetConfiguration) new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "devices=[" + String.join(",", macAddresses) + "]");

        ProfinetDriverContext context = new ProfinetDriverContext();
        context.setConfiguration(configuration);

        Map<String, ProfinetDevice> devices = configuration.getConfiguredDevices();


        for (String mac : macAddresses) {
            assert(devices.containsKey(mac.replace(":", "")));
        }
    }

    @Test
    public void readProfinetDevicesMultiple() throws DecoderException, PlcConnectionException {

        String[] macAddresses = new String[] {"CA:FE:00:00:00:01","CA:FE:00:00:00:02","CA:FE:00:00:00:03"};
        ProfinetConfiguration configuration = (ProfinetConfiguration) new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "devices=[" + String.join(",", macAddresses) + "]");

        ProfinetDriverContext context = new ProfinetDriverContext();
        context.setConfiguration(configuration);

        Map<String, ProfinetDevice> devices = configuration.getConfiguredDevices();

        for (String mac : macAddresses) {
            assert(devices.containsKey(mac.replace(":", "")));
        }
    }

    @Test
    public void readProfinetLowerCase() throws DecoderException, PlcConnectionException {

        String[] macAddresses = new String[] {"00:0c:29:75:25:67"};
        ProfinetConfiguration configuration = (ProfinetConfiguration) new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "devices=[" + String.join(",", macAddresses) + "]");

        ProfinetDriverContext context = new ProfinetDriverContext();
        context.setConfiguration(configuration);

        Map<String, ProfinetDevice> devices = configuration.getConfiguredDevices();

        for (String mac : macAddresses) {
            assert(devices.containsKey(mac.replace(":", "")));
        }
    }


}

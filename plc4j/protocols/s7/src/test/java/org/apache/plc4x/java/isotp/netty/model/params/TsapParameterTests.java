package org.apache.plc4x.java.isotp.netty.model.params;

import org.apache.plc4x.java.isotp.netty.model.types.DeviceGroup;
import org.apache.plc4x.java.isotp.netty.model.types.ParameterCode;
import org.junit.jupiter.api.*;

import static org.junit.jupiter.api.Assertions.*;

class TsapParameterTests {
    private TsapParameter tsapParameter;

    @BeforeEach
    void setUp() {
    }

    @AfterEach
    void tearDown() {
        tsapParameter = null;
    }

    @Test
    @Tag("fast")
    void calledPartameterTest() {
        DeviceGroup deviceGroup = DeviceGroup.valueOf((byte)0);
        tsapParameter = new CalledTsapParameter(deviceGroup, (byte)1, (byte)4);

        assertTrue(tsapParameter.getDeviceGroup() == DeviceGroup.valueOf((byte)0), "Device group incorrect");
        assertTrue(tsapParameter.getRackNumber() == (byte)1, "Rack number not correct");
        assertTrue(tsapParameter.getSlotNumber() == (byte)4, "Slot number not coorect");
        assertTrue(tsapParameter.getType() == ParameterCode.CALLED_TSAP);
    }

    @Test
    @Tag("fast")
    void callingPartameterTest() {
        DeviceGroup deviceGroup = DeviceGroup.valueOf((byte)0);
        tsapParameter = new CallingTsapParameter(deviceGroup, (byte)2, (byte)5);

        assertTrue(tsapParameter.getDeviceGroup() == DeviceGroup.valueOf((byte)0), "Device group incorrect");
        assertTrue(tsapParameter.getRackNumber() == (byte)2, "Rack number not correct");
        assertTrue(tsapParameter.getSlotNumber() == (byte)5, "Slot number not coorect");
        assertTrue(tsapParameter.getType() == ParameterCode.CALLING_TSAP);
    }

}
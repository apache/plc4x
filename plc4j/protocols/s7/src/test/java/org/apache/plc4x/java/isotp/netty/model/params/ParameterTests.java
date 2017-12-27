package org.apache.plc4x.java.isotp.netty.model.params;

import org.apache.plc4x.java.isotp.netty.model.types.DeviceGroup;
import org.apache.plc4x.java.isotp.netty.model.types.ParameterCode;
import org.apache.plc4x.java.isotp.netty.model.types.TpduSize;
import org.junit.jupiter.api.*;

import static org.junit.jupiter.api.Assertions.*;

class ParameterTest {
    private ParameterTest tsapParaameter;

    @BeforeEach
    void setUp() {
    }

    @AfterEach
    void tearDown() {
        tsapParaameter = null;
    }

    @Test
    @Tag("fast")
    void checksumPartameter() {
        ChecksumParameter checksumParameter = new ChecksumParameter((byte)1);

        assertTrue(checksumParameter.getChecksum() == (byte)1, "Checksum incorrect");
        assertTrue(checksumParameter.getType() == ParameterCode.CHECKSUM);
    }

    @Test
    @Tag("fast")
    void disconnectAdditionalInformationParameter() {
        byte[] data = {(byte)1, (byte)2};
        DisconnectAdditionalInformationParameter disconnectParameter = new DisconnectAdditionalInformationParameter(data);

        assertTrue(disconnectParameter.getData()[0] == (byte)1, "Return parameter incorrect");
        assertTrue(disconnectParameter.getData()[1] == (byte)2, "Return parameter incorrect");
        assertTrue(disconnectParameter.getType() == ParameterCode.DISCONNECT_ADDITIONAL_INFORMATION);
    }

    @Test
    @Tag("fast")
    void tpduSizeParameter() {
        TpduSizeParameter tpduSizeParameter = new TpduSizeParameter(TpduSize.SIZE_512);

        assertTrue(tpduSizeParameter.getTpduSize() == TpduSize.SIZE_512, "Tpdu size incorrect");
        assertTrue(tpduSizeParameter.getType() == ParameterCode.TPDU_SIZE);
    }
}
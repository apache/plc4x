package org.apache.plc4x.java.isotp.netty.model;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import org.apache.plc4x.java.isotp.netty.model.params.Parameter;
import org.apache.plc4x.java.isotp.netty.model.tpdus.ErrorTpdu;
import org.apache.plc4x.java.isotp.netty.model.types.RejectCause;
import org.apache.plc4x.java.isotp.netty.model.types.TpduCode;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class IsoTPMessageTests {

    @Test
    @Tag("fast")
    void testIsoTPMessage() {
        short destinationReference = 0x1;
        RejectCause rejectCause = RejectCause.REASON_NOT_SPECIFIED;
        List<Parameter> parameters = Collections.emptyList();
        ByteBuf errorData = Unpooled.buffer();
        ErrorTpdu tpdu = new ErrorTpdu(destinationReference, rejectCause, parameters, errorData);
        ByteBuf userData = Unpooled.buffer();
        IsoTPMessage isoTpMessage = new IsoTPMessage(tpdu, userData);

        errorData.writeByte(0x72);
        userData.writeByte(0x32);

        assertTrue(isoTpMessage.getTpdu().getTpduCode() == TpduCode.TPDU_ERROR, "Unexpected Tpdu");
        // Question: do we need two user data fields?
        assertTrue(tpdu.getUserData().readByte() == (byte) 0x72, "Unexpected user data");
        assertTrue(isoTpMessage.getUserData().readByte() == (byte) 0x32, "Unexpected user data");
    }


}
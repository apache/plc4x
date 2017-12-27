package org.apache.plc4x.java.s7.netty.model.messages;

import org.apache.plc4x.java.s7.netty.model.params.S7Parameter;
import org.apache.plc4x.java.s7.netty.model.payloads.S7Payload;
import org.apache.plc4x.java.s7.netty.model.types.MessageType;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import java.util.ArrayList;

import static org.junit.jupiter.api.Assertions.*;

class S7MessageTests {

    @Test
    @Tag("fast")
    void setupCommunictionsRequestMessage() {
        short tpduReference = 1;
        short maxAmqCaller = 4;
        short maxAmqCallee = 8;
        short pduLength = 128;

        SetupCommunicationRequestMessage setupMessage = new SetupCommunicationRequestMessage(tpduReference, maxAmqCaller, maxAmqCallee, pduLength);

        assertTrue(setupMessage.getTpduReference() == tpduReference, "Unexpected tpdu value");
        assertTrue(setupMessage.getMessageType() == MessageType.JOB, "Unexpected message type");
    }

    @Test
    @Tag("fast")
    void s7RequestMessage() {
        MessageType messageType = MessageType.USER_DATA;
        short tpduReference = 1;
        ArrayList<S7Parameter> s7Parameters = null;
        ArrayList<S7Payload> s7Payloads = null;

        S7RequestMessage message = new S7RequestMessage(messageType, tpduReference, s7Parameters, s7Payloads);

        assertTrue(message.getTpduReference() == tpduReference, "Unexpected tpdu value");
        assertTrue(message.getMessageType() == MessageType.USER_DATA, "Unexpected message type");
    }

    @Test
    @Tag("fast")
    void s7ResponseMessage() {
        MessageType messageType = MessageType.USER_DATA;
        short tpduReference = 1;
        ArrayList<S7Parameter> s7Parameters = null;
        ArrayList<S7Payload> s7Payloads = null;
        byte errorClass = 0x1;
        byte errorCode = 0x23;

        S7ResponseMessage message = new S7ResponseMessage(messageType, tpduReference, s7Parameters, s7Payloads, errorClass, errorCode);

        assertTrue(message.getTpduReference() == tpduReference, "Unexpected tpdu value");
        assertTrue(message.getMessageType() == MessageType.USER_DATA, "Unexpected message type");
        assertTrue(message.getErrorClass() == 0x1, "Unexpected error class");
        assertTrue(message.getErrorCode() == 0x23, "Unexpected error code");
    }


}
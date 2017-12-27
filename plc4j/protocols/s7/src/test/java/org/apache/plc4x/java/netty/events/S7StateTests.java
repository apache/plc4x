package org.apache.plc4x.java.netty.events;

import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class S7StateTests {

    @Test
    @Tag("fast")
    void testInitialS7ConnectionEvent() {
        S7ConnectionEvent s7event = new S7ConnectionEvent();

        assertTrue(s7event.getState() == S7ConnectionState.INITIAL, "Initial state not initial");
    }

    @Test
    @Tag("fast")
    void testS7ConnectionEvent() {
        S7ConnectionEvent s7event = new S7ConnectionEvent(S7ConnectionState.SETUP_COMPLETE);

        assertTrue(s7event.getState() == S7ConnectionState.SETUP_COMPLETE, "State not correct");
    }

}
package org.apache.plc4x.java.s7.connection;

import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.model.Address;
import org.apache.plc4x.java.s7.model.S7Address;
import org.apache.plc4x.java.s7.model.S7BitAddress;
import org.apache.plc4x.java.s7.model.S7DataBlockAddress;
import org.apache.plc4x.java.s7.netty.model.types.MemoryArea;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class S7PlcConnectionTests {

    private S7PlcConnection  s7PlcConnection;

    @BeforeEach
    void setUp() {
        s7PlcConnection = new S7PlcConnection("localhost", 1, 2);
    }

    @AfterEach
    void tearDown() {
        s7PlcConnection = null;
    }

    @Test
    void initialState() {
        assertTrue(s7PlcConnection.getHostName().equalsIgnoreCase("localhost"), "Hostname is incorrect");
        assertTrue(s7PlcConnection.getRack() == 1, "Rack is incorrect");
        assertTrue(s7PlcConnection.getSlot() == 2, "Slot is incorrect");
        assertTrue(s7PlcConnection.getPduSize() == 1024, "Pdu size is incorrect"); // Why is this hard coded?
    }

    @Test
    void emptyParseAddress() {
        try {
            s7PlcConnection.parseAddress("");
        }
        catch (PlcException exception) {
            assertTrue(exception.getMessage().startsWith("Address string doesn't match"), "Unexpected exception");
        }
    }

    @Test
    void parseDatablockAddress() {
        try {
            S7DataBlockAddress address = (S7DataBlockAddress) s7PlcConnection.parseAddress("DATA_BLOCKS/20/100");
            assertTrue(address.getDataBlockNumber() == 20, "unexpected data block");
            assertTrue(address.getByteOffset() == 100, "unexpected byte offset");
        }
        catch (PlcException exception) {
            fail("valid data block address");
        }
    }

    @Test
    void parseAddressAddress() {
        try {
            S7Address address = (S7Address) s7PlcConnection.parseAddress("TIMERS/10");
            assertTrue(address.getMemoryArea() == MemoryArea.TIMERS, "unexpected memory area");
            assertTrue(address.getByteOffset() == 10, "unexpected byte offset");
        }
        catch (PlcException exception) {
            fail("valid timer block address");
        }
    }

    @Test
    void parseAddressBitAddress() {
        try {
            S7BitAddress address = (S7BitAddress) s7PlcConnection.parseAddress("TIMERS/10/4");
            assertTrue(address.getMemoryArea() == MemoryArea.TIMERS, "unexpected memory area");
            assertTrue(address.getByteOffset() == 10, "unexpected byte offset");
            assertTrue(address.getBitOffset() == 4, "unexpected but offset");
        }
        catch (PlcException exception) {
            fail("valid timer block bit address");
        }
    }

    // TODO more tests for connect, close, read and write
}
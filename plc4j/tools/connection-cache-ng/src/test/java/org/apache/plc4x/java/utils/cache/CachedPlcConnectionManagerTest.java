package org.apache.plc4x.java.utils.cache;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcConnectionManager;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

public class CachedPlcConnectionManagerTest {

    @Test
    public void testSingleConnectionRequestTest() {
        PlcConnectionManager mockConnectionManager = Mockito.mock(PlcConnectionManager.class);
        CachedPlcConnectionManager connectionManager = CachedPlcConnectionManager.getBuilder(mockConnectionManager).build();
        try(PlcConnection connection = connectionManager.getConnection("test")) {
            Assertions.assertInstanceOf(LeasedPlcConnection.class, connection);
        } catch (Exception e) {
            Assertions.fail("Not expecting an exception here", e);
        }
        // TODO: Check getConnection was called on the mockConnectionManager instance ...
    }

}

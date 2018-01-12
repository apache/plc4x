package org.apache.plc4x.camel;


import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;

public class ConstantsTest {

    @Test
    public void testConstantsNotInstanceable() throws Exception {
        Assertions.assertThrows(IllegalAccessException.class, Constants.class::newInstance);
    }
}
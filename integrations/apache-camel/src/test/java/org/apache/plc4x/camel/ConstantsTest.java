package org.apache.plc4x.camel;


import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;

import java.lang.reflect.Constructor;

public class ConstantsTest {

    @Test
    public void testConstantsNotInstanceable() throws Exception {
        Assertions.assertThrows(IllegalStateException.class, () -> {
            try {
                Constructor<Constants> constructor = Constants.class.getDeclaredConstructor();
                constructor.setAccessible(true);
                constructor.newInstance();
            } catch (Exception e) {
                throw e.getCause();
            }
        });
    }
}
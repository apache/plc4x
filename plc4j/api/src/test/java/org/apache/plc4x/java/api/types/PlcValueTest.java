package org.apache.plc4x.java.api.types;

import org.apache.plc4x.java.api.exceptions.PlcIncompatibleDatatypeException;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.api.value.PlcValues;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

public class PlcValueTest {

    @Nested
    class Complex {

        @Test
        void complexTestCase_isComplex() {
            PlcValue value = PlcValues.of("Entry 1", PlcValues.of(
                PlcValues.of(true),
                PlcValues.of("Pimmel"),
                PlcValues.of(false),
                PlcValues.of("Arsch"),
                PlcValues.of(1278391)
            ));

            System.out.println(value);

            assertThrows(PlcIncompatibleDatatypeException.class, value::getBoolean);
            assertTrue(value.getValue("Entry 1").getIndex(0).getBoolean());
        }
    }

}
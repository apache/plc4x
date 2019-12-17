package org.apache.plc4x.java.api.value;

import org.apache.plc4x.java.api.value.PlcSimpleValue;

public class PlcInteger extends PlcSimpleValue<Integer> {

    public PlcInteger(Integer value) {
        super(value, true);
    }

    public PlcInteger(int value) {
        super(value, false);
    }

    @Override public boolean isString() {
        return true;
    }

    @Override public String getString() {
        return value.toString();
    }

    @Override public boolean getBoolean() {
        // We like C
        return !(value == 0);
    }

    @Override public double getDouble() {
        return (double)value;
    }

    @Override public float getFloat() {
        return (float)value;
    }

    @Override public long getLong() {
        return value;
    }

    @Override public int getInteger() {
        return value;
    }

    @Override public String toString() {
        return String.valueOf(value);
    }
}

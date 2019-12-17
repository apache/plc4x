package org.apache.plc4x.java.api.value;

public class PlcBoolean extends PlcSimpleValue<Boolean> {

    public PlcBoolean(Boolean value) {
        super(value, true);
    }

    public PlcBoolean(boolean bool) {
        super(bool, false);
    }

    @Override public boolean isBoolean() {
        return true;
    }

    @Override public boolean getBoolean() {
        return value;
    }

    @Override public double getDouble() {
        return value ? 1.0 : 0.0;
    }

    @Override public float getFloat() {
        return value ? 1.0f : 0.0f;
    }

    @Override public long getLong() {
        return value ? 1 : 0;
    }

    @Override public int getInteger() {
        return value ? 1 : 0;
    }

    @Override public String getString() {
        return value ? "true" : "false";
    }

    @Override public String toString() {
        return Boolean.toString(value);
    }
}

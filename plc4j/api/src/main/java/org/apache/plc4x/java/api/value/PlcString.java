package org.apache.plc4x.java.api.value;

public class PlcString extends PlcSimpleValue<String> {

    PlcString(String value) {
        super(value, true);
    }

    @Override public boolean isString() {
        return true;
    }

    @Override public String getString() {
        return value;
    }

    @Override public boolean getBoolean() {
        return Boolean.parseBoolean(value);
    }

    @Override public double getDouble() {
        return Double.parseDouble(value);
    }

    @Override public float getFloat() {
        return Float.parseFloat(value);
    }

    @Override public long getLong() {
        return Long.parseLong(value);
    }

    @Override public int getInteger() {
        return Integer.parseInt(value);
    }

    @Override public String toString() {
        return "\"" + value + "\"";
    }
}

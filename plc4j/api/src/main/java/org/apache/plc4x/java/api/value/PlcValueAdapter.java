package org.apache.plc4x.java.api.value;

import org.apache.plc4x.java.api.exceptions.PlcIncompatibleDatatypeException;

import java.util.List;
import java.util.Map;

public class PlcValueAdapter implements PlcValue {

    @Override public boolean isSimple() {
        return false;
    }

    @Override public boolean isNullable() {
        return false;
    }

    @Override public boolean isNull() {
        return false;
    }

    @Override public boolean is(Class<?> clazz) {
        return false;
    }

    @Override public boolean isConvertibleTo(Class<?> clazz) {
        return false;
    }

    @Override public <T> T get(Class<T> clazz) {
        throw new PlcIncompatibleDatatypeException("");
    }

    @Override public boolean isBoolean() {
        return false;
    }

    @Override public boolean getBoolean() {
        throw new PlcIncompatibleDatatypeException("");
    }

    @Override public boolean isDouble() {
        return false;
    }

    @Override public double getDouble() {
        throw new PlcIncompatibleDatatypeException("");
    }

    @Override public boolean isFloat() {
        return false;
    }

    @Override public float getFloat() {
        throw new PlcIncompatibleDatatypeException("");
    }

    @Override public boolean isLong() {
        return false;
    }

    @Override public long getLong() {
        throw new PlcIncompatibleDatatypeException("");
    }

    @Override public boolean isInteger() {
        return false;
    }

    @Override public int getInteger() {
        throw new PlcIncompatibleDatatypeException("");
    }

    @Override public boolean isString() {
        return false;
    }

    @Override public String getString() {
        throw new PlcIncompatibleDatatypeException("");
    }

    @Override public byte[] getRaw() {
        throw new PlcIncompatibleDatatypeException("");
    }

    @Override public boolean isList() {
        return false;
    }

    @Override public int length() {
        throw new PlcIncompatibleDatatypeException("");
    }

    @Override public PlcValue getIndex(int i) {
        throw new PlcIncompatibleDatatypeException("");
    }

    @Override public List<? extends PlcValue> getList() {
        throw new PlcIncompatibleDatatypeException("");
    }

    @Override public boolean isStruct() {
        return false;
    }

    @Override public boolean hasKey(String key) {
        return false;
    }

    @Override public PlcValue getValue(String key) {
        throw new PlcIncompatibleDatatypeException("");
    }

    @Override public Map<String, ? extends PlcValue> getStruct() {
        throw new PlcIncompatibleDatatypeException("");
    }

}

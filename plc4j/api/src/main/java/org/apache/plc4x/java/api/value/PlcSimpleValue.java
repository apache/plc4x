package org.apache.plc4x.java.api.value;

public abstract class PlcSimpleValue<T> extends PlcValueAdapter {

    final T value;
    final boolean isNullable;

    PlcSimpleValue(T value, boolean isNullable) {
        this.value = value;
        this.isNullable = isNullable;
    }

    @Override
    public boolean isSimple() {
        return true;
    }

    @Override
    public boolean isNullable() {
        return isNullable;
    }

    @Override
    public boolean isNull() {
        return isNullable && value == null;
    }

}

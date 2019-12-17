package org.apache.plc4x.java.api.value;

import org.apache.plc4x.java.api.value.PlcValueAdapter;

public abstract class PlcSimpleValue<T> extends PlcValueAdapter {

    protected final T value;
    protected final boolean isNullable;

    public PlcSimpleValue(T value, boolean isNullable) {
        this.value = value;
        this.isNullable = isNullable;
    }

    @Override public boolean isSimple() {
        return true;
    }

    @Override public boolean isNullable() {
        return isNullable;
    }

    @Override public boolean isNull() {
        return isNullable && value == null;
    }

}

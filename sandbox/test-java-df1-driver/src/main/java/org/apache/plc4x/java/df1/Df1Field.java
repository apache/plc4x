package org.apache.plc4x.java.df1;

import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.df1.fields.DataType;

public class Df1Field implements PlcField {

    private final int address;
    private final int size;
    private final DataType dataType;

    public Df1Field(int address, int size, DataType dataType) {
        this.address = address;
        this.size = size;
        this.dataType = dataType;
    }

    public int getAddress() {
        return address;
    }

    public int getSize() {
        return size;
    }

    public DataType getDataType() {
        return dataType;
    }

    public static PlcField of(String fieldQuery) {
        return new Df1Field(11, 2, DataType.INTEGER);
    }
}

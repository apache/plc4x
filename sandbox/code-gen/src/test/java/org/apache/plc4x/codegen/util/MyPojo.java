package org.apache.plc4x.codegen.util;

public class MyPojo {

    private Double field1;

    private Double field2;

    private Double field3;

    public MyPojo() {
    }

    public Double getField1() {
        return this.field1;
    }

    public Double getField2() {
        return this.field2;
    }

    public Double getField3() {
        return this.field3;
    }

    public void setField1(Double field1) {
        this.field1 = field1;
    }

    public void setField2(Double field2) {
        this.field2 = field2;
    }

    public void setField3(Double field3) {
        this.field3 = field3;
    }

    public void encode(org.apache.plc4x.codegen.api.Buffer buffer) {
    }

    public static MyPojo decode(org.apache.plc4x.codegen.api.Buffer buffer) {
        MyPojo instance = new MyPojo();
        buffer.readUint8();
        return instance;
    }

}
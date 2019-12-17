package org.apache.plc4x.java.api.value;

import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.Map;

public class PlcValues {

    public static PlcValue of(String s) {
        return new PlcString(s);
    }

    public static PlcValue of(Integer i) {
        return new PlcInteger(i);
    }

    public static PlcValue of(int i) {
        return new PlcInteger(i);
    }

    public static PlcValue of(Boolean b) {
        return new PlcBoolean(b);
    }

    public static PlcValue of(boolean b) {
        return new PlcBoolean(b);
    }

    public static PlcValue of(List<PlcValue> list) {
        return new PlcList(list);
    }

    public static PlcValue of(PlcValue... items) {
        return new PlcList(Arrays.asList(items));
    }

    public static PlcValue of(String key, PlcValue value) {
        return new PlcStruct(Collections.singletonMap(key, value));
    }

    public static PlcValue of(Map<String, PlcValue> map) {
        return new PlcStruct(map);
    }
}

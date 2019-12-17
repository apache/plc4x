package org.apache.plc4x.java.api.value;

import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.api.value.PlcValueAdapter;

import java.util.Collections;
import java.util.Map;
import java.util.stream.Collectors;

public class PlcStruct extends PlcValueAdapter {

    private final Map<String, PlcValue> map;

    public PlcStruct(Map<String, PlcValue> map) {
        this.map = Collections.unmodifiableMap(map);
    }

    @Override public boolean isStruct() {
        return true;
    }

    @Override public boolean hasKey(String key) {
        return map.containsKey(key);
    }

    @Override public PlcValue getValue(String key) {
        return map.get(key);
    }

    @Override public Map<String, ? extends PlcValue> getStruct() {
        return map;
    }

    @Override public String toString() {
        return "{" + map.entrySet().stream().map(entry -> String.format("\"%s\": %s", entry.getKey(), entry.getValue())).collect(Collectors.joining(",")) + "}";
    }
}

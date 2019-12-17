package org.apache.plc4x.java.api.value;

import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.api.value.PlcValueAdapter;

import java.util.Collections;
import java.util.List;
import java.util.stream.Collectors;

public class PlcList extends PlcValueAdapter {

    private final List<PlcValue> listItems;

    PlcList(List<PlcValue> listItems) {
        this.listItems = Collections.unmodifiableList(listItems);
    }

    @Override public boolean isList() {
        return true;
    }

    @Override public int length() {
        return listItems.size();
    }

    @Override public PlcValue getIndex(int i) {
        return listItems.get(i);
    }

    @Override public List<? extends PlcValue> getList() {
        return listItems;
    }

    @Override public String toString() {
        return "[" + listItems.stream().map(PlcValue::toString).collect(Collectors.joining(",")) + "]";
    }
}

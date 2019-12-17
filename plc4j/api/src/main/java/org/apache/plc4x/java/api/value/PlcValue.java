package org.apache.plc4x.java.api.value;

import java.util.List;
import java.util.Map;
import java.util.Set;

/**
 * Base Type of all Types.
 */
public interface PlcValue {

    // Simple Types

    boolean isSimple();

    boolean isNullable();

    boolean isNull();

    // Generic (\o/ Sebastian)

    boolean is(Class<?> clazz);

    boolean isConvertibleTo(Class<?> clazz);

    <T> T get(Class<T> clazz);

    // Boolean

    boolean isBoolean();

    boolean getBoolean();

    // Integer

    boolean isLong();

    long getLong();

    boolean isInteger();

    int getInteger();

    // Floating Point

    boolean isDouble();

    double getDouble();

    boolean isFloat();

    float getFloat();

    // String

    boolean isString();

    String getString();

    // Raw Access

    byte[] getRaw();

    // List Methods

    boolean isList();

    int length();

    PlcValue getIndex(int i);

    List<? extends PlcValue> getList();

    // Struct Methods

    boolean isStruct();

    Set<String> getKeys();

    boolean hasKey(String key);

    PlcValue getValue(String key);

    Map<String, ? extends PlcValue> getStruct();

}

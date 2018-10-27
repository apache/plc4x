package org.apache.plc4x.java.opm;

import java.lang.reflect.Field;

/**
 * Utility methods for usage with OPM.
 */
public class OpmUtils {

    static <T> PlcEntity getPlcEntityAndCheckPreconditions(Class<T> clazz) {
        PlcEntity annotation = clazz.getAnnotation(PlcEntity.class);
        if (annotation == null) {
            throw new IllegalArgumentException("Given Class is no Plc Entity, i.e., not annotated with @PlcEntity");
        }
        // Check if default constructor exists
        try {
            clazz.getConstructor();
        } catch (NoSuchMethodException e) {
            throw new IllegalArgumentException("Cannot use PlcEntity without default constructor", e);
        }
        return annotation;
    }

    static String extractAddress(Object proxy) throws OPMException {
        String address;
        try {
            Field field = proxy.getClass().getDeclaredField(PlcEntityManager.PLC_ADDRESS_FIELD_NAME);
            field.setAccessible(true);
            address = (String) field.get(proxy);
        } catch (IllegalAccessException | NoSuchFieldException e) {
            throw new OPMException("Problem with accessing internal plc address", e);
        }
        return address;
    }

}

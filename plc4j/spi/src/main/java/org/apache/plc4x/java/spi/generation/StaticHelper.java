/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.spi.generation;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.value.PlcValue;

import java.util.Collection;
import java.util.List;

public class StaticHelper {

    public static int ARRAY_SIZE_IN_BYTES(Object obj) {
        if (obj instanceof List) {
            List<?> list = (List<?>) obj;
            int numBytes = 0;
            for (Object element : list) {
                if (!(element instanceof Message)) {
                    throw new RuntimeException(
                        "Array elements for array size in bytes must implement Message interface");
                }
                numBytes += ((Message) element).getLengthInBytes();
            }
            return numBytes;
        }
        if (obj.getClass().isArray() && !obj.getClass().getComponentType().isPrimitive()) {
            Object[] arr = (Object[]) obj;
            int numBytes = 0;
            for (Object element : arr) {
                if (!(element instanceof Message)) {
                    throw new RuntimeException(
                        "Array elements for array size in bytes must implement Message interface");
                }
                numBytes += ((Message) element).getLengthInBytes();
            }
            return numBytes;
        }
        throw new RuntimeException("Unable to calculate array size in bytes for type " + obj.getClass().getName());
    }

    public static int COUNT(Object obj) {
        if (obj == null) {
            return 0;
        }
        if (obj.getClass().isArray()) {
            if (obj.getClass().getComponentType() != null && obj.getClass().getComponentType().isPrimitive()) {
                if (obj.getClass().getComponentType() == boolean.class) {
                    boolean[] arr = (boolean[]) obj;
                    return arr.length;
                }
                if (obj.getClass().getComponentType() == byte.class) {
                    byte[] arr = (byte[]) obj;
                    return arr.length;
                }
                if (obj.getClass().getComponentType() == short.class) {
                    short[] arr = (short[]) obj;
                    return arr.length;
                }
                if (obj.getClass().getComponentType() == int.class) {
                    int[] arr = (int[]) obj;
                    return arr.length;
                }
                if (obj.getClass().getComponentType() == long.class) {
                    long[] arr = (long[]) obj;
                    return arr.length;
                }
                if (obj.getClass().getComponentType() == float.class) {
                    float[] arr = (float[]) obj;
                    return arr.length;
                }
                if (obj.getClass().getComponentType() == double.class) {
                    double[] arr = (double[]) obj;
                    return arr.length;
                }
            } else {
                Object[] arr = (Object[]) obj;
                return arr.length;
            }
        } else if (obj instanceof Collection) {
            Collection<?> col = (Collection<?>) obj;
            return col.size();
        }
        throw new PlcRuntimeException("Unable to count object of type " + obj.getClass().getName());
    }

    public static int STR_LEN(Object str) {
        if (str == null) {
            return 0;
        }
        if (str instanceof PlcValue) {
            PlcValue plcValue = (PlcValue) str;
            return plcValue.getString().length();
        }
        return str.toString().length();
    }

    public static <T> T CAST(Object obj, Class<T> clazz) {
        try {
            return clazz.cast(obj);
        } catch (ClassCastException e) {
            throw new PlcRuntimeException("Unable to cast object of type " + obj.getClass().getName() + " to " + clazz.getName());
        }
    }

    public static int CEIL(double value) {
        return (int) Math.ceil(value);
    }
    
    public static int PADCOUNT(Object obj, boolean hasNext) {
        return hasNext ? COUNT(obj) : 0;
    }    

}

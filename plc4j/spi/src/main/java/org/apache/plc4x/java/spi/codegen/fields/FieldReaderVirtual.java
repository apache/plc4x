/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.spi.codegen.fields;

import org.apache.plc4x.java.spi.codegen.FieldCommons;
import org.apache.plc4x.java.spi.codegen.io.DataReader;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.WithReaderArgs;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.math.BigInteger;

public class FieldReaderVirtual<T> implements FieldCommons {

    private static final Logger LOGGER = LoggerFactory.getLogger(FieldReaderVirtual.class);

    @SuppressWarnings({"unused", "unchecked"})
    public T readVirtualField(String logicalName, Class<T> type, Object valueExpression, WithReaderArgs... readerArgs) throws ParseException {
        LOGGER.debug("reading field {}", logicalName);
        if (type.isPrimitive()) {
            // for primitives, we need to cast to the primitive as this does autoboxing
            if (type == boolean.class) {
                return (T) (Boolean) (boolean) valueExpression;
            } else if (type == byte.class) {
                Number valueExpressionNumber = (Number) valueExpression;
                long longValue = valueExpressionNumber.longValue();
                // 0xFF is a special case and a legit value for (unsigned) byte
                if ((byte) longValue != longValue && longValue != 0xFF) {
                    throw new ArithmeticException("byte overflow");
                }
                return (T) (Byte) valueExpressionNumber.byteValue();
            } else if (type == short.class) {
                Number valueExpressionNumber = (Number) valueExpression;
                long longValue = valueExpressionNumber.longValue();
                if ((short) longValue != longValue) {
                    throw new ArithmeticException("short overflow");
                }
                return (T) (Short) valueExpressionNumber.shortValue();
            } else if (type == int.class) {
                Number valueExpressionNumber = (Number) valueExpression;
                long longValue = valueExpressionNumber.longValue();
                if ((int) longValue != longValue) {
                    throw new ArithmeticException("integer overflow");
                }
                return (T) (Integer) valueExpressionNumber.intValue();
            } else if (type == long.class) {
                Number valueExpressionNumber = (Number) valueExpression;
                return (T) (Long) valueExpressionNumber.longValue();
            } else if (type == char.class) {
                return (T) (Character) (char) valueExpression;
            } else if (type == float.class) {
                Number valueExpressionNumber = (Number) valueExpression;
                return (T) (Float) valueExpressionNumber.floatValue();
            } else if (type == double.class) {
                Number valueExpressionNumber = (Number) valueExpression;
                return (T) (Double) valueExpressionNumber.doubleValue();
            } else {
                throw new IllegalStateException("Unmapped primitive " + type);
            }
        }
        if (type == String.class) {
            return type.cast(String.valueOf(valueExpression));
        }
        if (type == BigInteger.class) {
            long longValue = valueExpression instanceof Long ? (Long) valueExpression : ((Number)valueExpression).longValue();
            return (T) BigInteger.valueOf(longValue);
        }
        return type.cast(valueExpression);
    }

}

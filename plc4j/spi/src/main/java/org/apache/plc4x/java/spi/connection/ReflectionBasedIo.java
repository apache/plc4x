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

package org.apache.plc4x.java.spi.connection;

import org.apache.commons.lang3.reflect.MethodUtils;
import org.apache.plc4x.java.spi.generation.Message;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;

public class ReflectionBasedIo<BASE extends Message> implements Parser<BASE>, Serializer<BASE> {

    private final Method parseMethod;
    private final Method serializeMethod;

    public ReflectionBasedIo(Class<BASE> clazz) {
        String className = clazz.getSimpleName() + "IO";
        String fqcn = clazz.getPackage().getName() + ".io." + className;
        try {
            Class<?> ioClass = Class.forName(fqcn);
            parseMethod = MethodUtils.getMatchingMethod(ioClass, "parse", ReadBuffer.class);
            serializeMethod = MethodUtils.getMatchingMethod(ioClass, "serialize", WriteBuffer.class, clazz);
        } catch (ClassNotFoundException e) {
            throw new IllegalStateException(
                String.format("Unnable to get suitable IO Class for given Message. Expected IO Class '%s' for Class '%s'", clazz.getName(), fqcn));
        }
    }

    @Override public BASE parse(ReadBuffer io) {
        try {
            return (BASE) parseMethod.invoke(null, io);
        } catch (IllegalAccessException | InvocationTargetException e) {
            throw new IllegalStateException("Unable to use the parse Method!", e);
        }
    }

    @Override public void serialize(WriteBuffer io, BASE message) {
        try {
            serializeMethod.invoke(null, io, message);
        } catch (IllegalAccessException | InvocationTargetException e) {
            throw new IllegalStateException("Unable to use the serialize Method!", e);
        }
    }
}

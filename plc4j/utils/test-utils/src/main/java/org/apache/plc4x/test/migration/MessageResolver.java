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
package org.apache.plc4x.test.migration;

import java.util.Map;
import org.apache.commons.lang3.StringUtils;
import org.apache.plc4x.java.spi.generation.MessageIO;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.test.driver.exceptions.DriverTestsuiteException;
import org.apache.plc4x.test.parserserializer.exceptions.ParserSerializerTestsuiteException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.lang.reflect.Modifier;
import java.util.LinkedList;
import java.util.List;

public class MessageResolver {

    private final static Logger LOGGER = LoggerFactory.getLogger(MessageResolver.class);

    /**
     * Returns the messageIO class based on a configured package. convention out of {@code protocolName} {@code outputFlavor} and {@code name}.
     * If this fails its tries a fallback using the deprecated attribute {@code className}
     *
     * @param options Test framework options
     * @param name    name of the message
     * @return the found MessageIO
     * @throws DriverTestsuiteException if a MessageIO couldn't be found.
     */
    @SuppressWarnings("rawtypes")
    public static MessageIO getMessageIO(Map<String, String> options, String name) throws DriverTestsuiteException {
        try {
            return MessageResolver.getMessageIOType(options, name).getMessageIo().newInstance();
        } catch (InstantiationException | IllegalAccessException e) {
            throw new DriverTestsuiteException(e);
        }
    }

    public static MessageIO getMessageIOStaticLinked(Map<String, String> options, String typeName) throws ParserSerializerTestsuiteException {
        try {
            TypePair typePair = getMessageIOType(options, typeName);
            Class<?> ioRootClass = typePair.getType();
            Class<?> ioClass = typePair.getMessageIo();
            Method staticParseMethod = null;
            Method staticSerializeMethod = null;
            final List<Class<?>> parameterTypes = new LinkedList<>();
            for (Method method : ioClass.getMethods()) {
                if (method.getName().equals("staticParse") && Modifier.isStatic(method.getModifiers()) &&
                    (method.getReturnType() == ioRootClass)) {
                    staticParseMethod = method;

                    // Get a list of additional parameter types for the parser.
                    for (int i = 1; i < method.getParameterCount(); i++) {
                        Class<?> parameterType = staticParseMethod.getParameterTypes()[i];
                        parameterTypes.add(parameterType);
                    }
                }
                if (method.getName().equals("staticSerialize") && Modifier.isStatic(method.getModifiers()) &&
                    (method.getParameterTypes()[1] == ioRootClass)) {
                    staticSerializeMethod = method;
                }
            }
            if ((staticParseMethod == null) || (staticSerializeMethod == null)) {
                throw new ParserSerializerTestsuiteException(
                    "Unable to instantiate IO component. Missing static parse or serialize methods.");
            }
            final Method parseMethod = staticParseMethod;
            final Method serializeMethod = staticSerializeMethod;
            return new MessageIO() {
                @Override
                public Object parse(ReadBuffer io, Object... args) throws ParseException {
                    try {
                        Object[] argValues = new Object[args.length + 1];
                        argValues[0] = io;
                        for (int i = 1; i <= args.length; i++) {
                            String parameterValue = (String) args[i - 1];
                            Class<?> parameterType = parameterTypes.get(i - 1);
                            if (parameterType == Boolean.class) {
                                argValues[i] = Boolean.parseBoolean(parameterValue);
                            } else if (parameterType == Byte.class) {
                                argValues[i] = Byte.parseByte(parameterValue);
                            } else if (parameterType == Short.class) {
                                argValues[i] = Short.parseShort(parameterValue);
                            } else if (parameterType == Integer.class) {
                                argValues[i] = Integer.parseInt(parameterValue);
                            } else if (parameterType == Long.class) {
                                argValues[i] = Long.parseLong(parameterValue);
                            } else if (parameterType == Float.class) {
                                argValues[i] = Float.parseFloat(parameterValue);
                            } else if (parameterType == Double.class) {
                                argValues[i] = Double.parseDouble(parameterValue);
                            } else if (parameterType == String.class) {
                                argValues[i] = parameterValue;
                            } else if (Enum.class.isAssignableFrom(parameterType)) {
                                argValues[i] = Enum.valueOf((Class<? extends Enum>) parameterType, parameterValue);
                            } else {
                                throw new ParseException("Currently unsupported parameter type");
                            }
                        }

                        return parseMethod.invoke(null, argValues);
                    } catch (IllegalAccessException | InvocationTargetException e) {
                        throw new ParseException("error parsing", e);
                    }
                }

                @Override
                public void serialize(WriteBuffer io, Object value, Object... args) throws ParseException {
                    try {
                        serializeMethod.invoke(null, io, value);
                    } catch (IllegalAccessException | InvocationTargetException e) {
                        throw new ParseException("error serializing", e);
                    }
                }
            };
        } catch (DriverTestsuiteException e) {
            throw new ParserSerializerTestsuiteException("Unable to instantiate IO component", e.getCause());
        }
    }

    @SuppressWarnings("unchecked")
    private static TypePair getMessageIOType(Map<String, String> options, String typeName) throws DriverTestsuiteException {
        String extraMessage = "";
        if (options.containsKey("package")) {
            try {
                return lookup(options.get("package"), typeName);
            } catch (ClassNotFoundException e) {
                extraMessage = "custom package '" + options.get("package") + "' and ";
            }
        }

        String protocolName = options.get("protocolName");
        String outputFlavor = options.get("outputFlavor");
        String classPackage = String.format("org.apache.plc4x.java.%s.%s", protocolName, StringUtils.replace(outputFlavor, "-", ""));
        try {
            return lookup(classPackage, typeName);
        } catch (ClassNotFoundException e) {
            throw new DriverTestsuiteException("Could not find " + typeName + " in " + extraMessage + "standard package '" + classPackage + "'");
        }
    }

    private static TypePair lookup(String driverPackage, String typeName) throws ClassNotFoundException {
        try {
            Package.getPackage(driverPackage);
        } catch (RuntimeException e) {
            throw new DriverTestsuiteException("Invalid or non existent package detected: " + driverPackage, e);
        }
        String ioRootClassName = driverPackage + "." + typeName;
        String ioClassName = driverPackage + ".io." + typeName + "IO";
        // make sure both type and it's IO are present
        return new TypePair(
            Class.forName(ioRootClassName),
            (Class<? extends MessageIO<?, ?>>) Class.forName(ioClassName)
        );
    }

    static class TypePair {
        private final Class<?> type;
        private final Class<? extends MessageIO<?, ?>> messageIo;

        TypePair(Class<?> type, Class<? extends MessageIO<?, ?>> messageIo) {
            this.type = type;
            this.messageIo = messageIo;
        }
        Class<?> getType() {
            return type;
        }
        Class<? extends MessageIO<?, ?>> getMessageIo() {
            return messageIo;
        }
    }

}

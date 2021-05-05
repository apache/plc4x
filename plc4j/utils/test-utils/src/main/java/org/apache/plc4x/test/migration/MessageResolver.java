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
     * Returns the messageIO class based on a convention out of {@code protocolName} {@code outputFlavor} and {@code name}.
     * If this fails its tries a fallback using the deprecated attribute {@code className}
     *
     * @param protocolName name of the protocol
     * @param outputFlavor flavor of the output (e.g read-write)
     * @param name         name of the message
     * @param className    deprecated fallback classname attribute
     * @return the found MessageIO
     * @throws DriverTestsuiteException if a MessageIO couldn't be found.
     */
    @SuppressWarnings("rawtypes")
    public static MessageIO getMessageIO(String protocolName, String outputFlavor, String name, @Deprecated String className) throws DriverTestsuiteException {
        try {
            return MessageResolver.getMessageIOType(protocolName, outputFlavor, name).newInstance();
        } catch (InstantiationException | IllegalAccessException ignore) {
            LOGGER.warn("\n!!!Un-migrated test!!!\n");
            try {
                return MessageResolver.getMessageIOTypeFallback(className).newInstance();
            } catch (InstantiationException | IllegalAccessException e) {
                throw new DriverTestsuiteException(e);
            }
        }
    }

    public static MessageIO getMessageIOStaticLinked(String protocolName, String outputFlavor, String typeName, @Deprecated String classNameAttributeValue) throws ParserSerializerTestsuiteException {
        String ioClassName, ioRootClassName;
        try {
            String classPackage = String.format("org.apache.plc4x.java.%s.%s", protocolName, StringUtils.replace(outputFlavor, "-", ""));
            try {
                Package.getPackage(classPackage);
            } catch (RuntimeException e) {
                throw new RuntimeException("fallback to old", e);
            }
            ioRootClassName = classPackage + "." + typeName;
            ioClassName = classPackage + ".io." + typeName + "IO";
            try {
                Class.forName(ioRootClassName);
                Class.forName(ioClassName);
            } catch (ClassNotFoundException e) {
                throw new RuntimeException("fallback to old", e);
            }
        } catch (RuntimeException e) {
            LOGGER.error("Error in serializer", e);
            ioRootClassName = classNameAttributeValue.substring(0, classNameAttributeValue.lastIndexOf('.') + 1) + typeName;
            ioClassName = classNameAttributeValue.substring(0, classNameAttributeValue.lastIndexOf('.') + 1) + "io." + typeName + "IO";
        }
        try {
            Class<?> ioRootClass = Class.forName(ioRootClassName);
            Class<?> ioClass = Class.forName(ioClassName);
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
        } catch (ClassNotFoundException e) {
            throw new ParserSerializerTestsuiteException("Unable to instantiate IO component", e);
        }
    }

    @SuppressWarnings("unchecked")
    private static Class<? extends MessageIO<?, ?>> getMessageIOType(String protocolName, String outputFlavor, String typeName) throws DriverTestsuiteException {
        String classPackage = String.format("org.apache.plc4x.java.%s.%s", protocolName, StringUtils.replace(outputFlavor, "-", ""));
        try {
            Package.getPackage(classPackage);
        } catch (RuntimeException e) {
            throw new RuntimeException("fallback to old", e);
        }
        String ioRootClassName = classPackage + "." + typeName;
        String ioClassName = classPackage + ".io." + typeName + "IO";
        try {
            Class.forName(ioRootClassName);
            return (Class<? extends MessageIO<?, ?>>) Class.forName(ioClassName);
        } catch (ClassNotFoundException e) {
            throw new RuntimeException("fallback to old", e);
        }
    }

    @SuppressWarnings("unchecked")
    private static Class<? extends MessageIO<?, ?>> getMessageIOTypeFallback(String messageClassName) throws DriverTestsuiteException {
        String ioClassName = messageClassName.substring(0, messageClassName.lastIndexOf('.')) + ".io." +
            messageClassName.substring(messageClassName.lastIndexOf('.') + 1) + "IO";
        try {
            final Class<?> ioClass = Class.forName(ioClassName);
            if (MessageIO.class.isAssignableFrom(ioClass)) {
                return (Class<? extends MessageIO<?, ?>>) ioClass;
            }
            throw new DriverTestsuiteException("IO class muss implement MessageIO interface");
        } catch (ClassNotFoundException e) {
            throw new DriverTestsuiteException("Error loading io class", e);
        }
    }

}

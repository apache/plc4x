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
package org.apache.plc4x.test.migration;

import java.util.Arrays;
import java.util.Map;

import org.apache.commons.lang3.StringUtils;
import org.apache.plc4x.java.spi.generation.*;
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
    public static MessageInput<?> getMessageInput(Map<String, String> options, String name) throws DriverTestsuiteException {
        try {
            return MessageResolver.getMessageIOType(options, name).getMessageInput();
        } catch (ClassNotFoundException e) {
            throw new DriverTestsuiteException(e);
        }
    }

    public static MessageInput<?> getMessageIOStaticLinked(Map<String, String> options, String typeName) throws ParserSerializerTestsuiteException {
        try {
            TypeMessageInput typeMessageInput = getMessageIOType(options, typeName);
            final List<Class<?>> parameterTypes = new LinkedList<>();
            for (Method method : typeMessageInput.type.getMethods()) {
                int parameterCount = method.getParameterCount();
                boolean isNonGenericParse = parameterCount > 1 && method.getParameterTypes()[parameterCount - 1] != Object[].class;
                if (method.getName().equals("staticParse") && Modifier.isStatic(method.getModifiers()) && isNonGenericParse) {
                    // Get a list of additional parameter types for the parser.
                    parameterTypes.addAll(Arrays.asList(method.getParameterTypes()).subList(1, parameterCount));
                    break;
                }
            }
            return (io, args) -> {
                Object[] argValues = new Object[args.length];
                for (int i = 0; i < args.length; i++) {
                    String parameterValue = (String) args[i];
                    Class<?> parameterType = parameterTypes.get(i);
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
                return typeMessageInput.getMessageInput().parse(io, argValues);
            };
        } catch (DriverTestsuiteException | ClassNotFoundException e) {
            throw new ParserSerializerTestsuiteException("Unable to instantiate IO component", e);
        }
    }

    @SuppressWarnings("unchecked")
    private static TypeMessageInput getMessageIOType(Map<String, String> options, String typeName) throws DriverTestsuiteException, ClassNotFoundException {
        String extraMessage = "";
        if (options.containsKey("package")) {
            try {
                return lookup(options.get("package"), typeName);
            } catch (NoSuchMethodException e) {
                extraMessage = "custom package '" + options.get("package") + "' and ";
            }
        }

        String protocolName = options.get("protocolName");
        String outputFlavor = options.get("outputFlavor");
        String classPackage = String.format("org.apache.plc4x.java.%s.%s", protocolName, StringUtils.replace(outputFlavor, "-", ""));
        try {
            return lookup(classPackage, typeName);
        } catch (NoSuchMethodException e) {
            throw new DriverTestsuiteException("Could not find " + typeName + " in " + extraMessage + "standard package '" + classPackage + "'");
        }
    }

    private static TypeMessageInput lookup(String driverPackage, String typeName) throws ClassNotFoundException, NoSuchMethodException {
        try {
            Package.getPackage(driverPackage);
        } catch (RuntimeException e) {
            throw new DriverTestsuiteException("Invalid or non existent package detected: " + driverPackage, e);
        }
        String ioRootClassName = driverPackage.replace("-", "") + "." + typeName;
        // make sure both type and it's IO are present
        Class<? extends Message> messageType = (Class<? extends Message>) Class.forName(ioRootClassName);
        Method staticParse = messageType.getMethod("staticParse", ReadBuffer.class, Object[].class);
        return new TypeMessageInput(
            messageType,
            (io, args) -> {
                try {
                    return staticParse.invoke(null, io, args);
                } catch (IllegalAccessException | InvocationTargetException e) {
                    throw new RuntimeException(e);
                }
            }
        );
    }

    static class TypeMessageInput {
        private final Class<? extends Message> type;
        private final MessageInput<?> messageInput;

        TypeMessageInput(Class<? extends Message> type, MessageInput<?> messageInput) {
            this.type = type;
            this.messageInput = messageInput;
        }

        Class<?> getType() {
            return type;
        }

        MessageInput<?> getMessageInput() {
            return messageInput;
        }
    }

}

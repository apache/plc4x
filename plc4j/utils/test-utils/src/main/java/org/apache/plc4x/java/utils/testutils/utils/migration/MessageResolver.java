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
package org.apache.plc4x.java.utils.testutils.utils.migration;

import org.apache.commons.lang3.Strings;
import org.apache.plc4x.java.spi.buffers.api.Message;
import org.apache.plc4x.java.spi.buffers.api.MessageInput;
import org.apache.plc4x.java.spi.buffers.api.ReadBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.utils.testutils.driver.exceptions.DriverTestsuiteException;
import org.apache.plc4x.java.utils.testutils.parserserializer.exceptions.ParserSerializerTestsuiteException;

import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.lang.reflect.Modifier;
import java.util.List;
import java.util.Map;

public class MessageResolver {

    /**
     * Returns the messageIO class based on a configured package. convention out of {@code protocolName} {@code outputFlavor} and {@code name}.
     * If this fails its tries a fallback using the deprecated attribute {@code className}
     *
     * @param options         Test framework options
     * @param name            name of the message
     * @param parserArguments the parser arguments to create an instance of the message
     * @return the found MessageIO
     * @throws DriverTestsuiteException if a MessageIO couldn't be found.
     */
    public static MessageInput<?> getMessageInput(Map<String, String> options, String name, List<String> parserArguments) throws DriverTestsuiteException {
        try {
            return MessageResolver.getMessageIOType(options, name, parserArguments);
        } catch (ClassNotFoundException e) {
            throw new DriverTestsuiteException(e);
        }
    }

    public static MessageInput<?> getMessageIOStaticLinked(Map<String, String> options, String typeName, List<String> parserArguments) throws ParserSerializerTestsuiteException {
        try {
            MessageInput<?> typeMessageInput = getMessageIOType(options, typeName, parserArguments);
//            final List<Class<?>> parameterTypes = new LinkedList<>();
//            for (Method method : typeMessageInput.type.getMethods()) {
//                int parameterCount = method.getParameterCount();
//                boolean isNonGenericParse = parameterCount > 1 && method.getParameterTypes()[parameterCount - 1] != Object[].class;
//                if (method.getName().equals("staticParse") && Modifier.isStatic(method.getModifiers()) && isNonGenericParse) {
//                    // Get a list of additional parameter types for the parser.
//                    parameterTypes.addAll(Arrays.asList(method.getParameterTypes()).subList(1, parameterCount));
//                    break;
//                }
//            }
            return typeMessageInput;
        } catch (DriverTestsuiteException | ClassNotFoundException e) {
            throw new ParserSerializerTestsuiteException("Unable to instantiate IO component", e);
        }
    }

    private static MessageInput<?> getMessageIOType(Map<String, String> options, String typeName, List<String> parserArguments) throws DriverTestsuiteException, ClassNotFoundException {
        String extraMessage = "";//750-375 V02.08.xx (FW 11)
        if (options.containsKey("package")) {
            try {
                return lookup(options.get("package"), typeName, parserArguments);
            } catch (NoSuchMethodException | ClassNotFoundException e) {
                // Class not found in custom package — fall through to standard package lookup
                extraMessage = "custom package '" + options.get("package") + "' and ";
            }
        }

        String protocolName = options.get("protocol-name");
        String outputFlavor = options.get("output-flavor");
        String classPackage = String.format("org.apache.plc4x.java.%s.%s", protocolName, Strings.CS.replace(outputFlavor, "-", ""));
        try {
            return lookup(classPackage, typeName, parserArguments);
        } catch (NoSuchMethodException e) {
            throw new DriverTestsuiteException("Could not find " + typeName + " in " + extraMessage + "standard package '" + classPackage + "'");
        }
    }

    @SuppressWarnings("unchecked")
    private static MessageInput<?> lookup(String driverPackage, String typeName, List<String> parserArguments) throws ClassNotFoundException, NoSuchMethodException {
        try {
            Package.getPackage(driverPackage);
        } catch (RuntimeException e) {
            throw new DriverTestsuiteException("Invalid or non existent package detected: " + driverPackage, e);
        }
        String ioRootClassName = driverPackage.replace("-", "") + "." + typeName;
        // make sure both type and it's IO are present
        Class<? extends Message> messageType = (Class<? extends Message>) Class.forName(ioRootClassName);
        Method parseMethod = null;
        for (Method method : messageType.getMethods()) {
            if (Modifier.isStatic(method.getModifiers()) && "staticParse".equals(method.getName()) && Message.class.isAssignableFrom(method.getReturnType())) {
                // because we still have var-arg and non var-arg methods we have to be careful
                // below is additional verification of staticParse method which refuses var-arg variant
                if (method.getParameterCount() == 1) {
                    parseMethod = method;
                    break;
                } else if (method.getParameterCount() >= 2 && method.getParameterTypes()[1] != Object[].class) {
                    // TODO above if statement can be removed later on - when we get rid of var-args
                    parseMethod = method;
                    break;
                }
            }
        }
        if (parseMethod == null) {
            throw new DriverTestsuiteException("Could not find static parse method for " + typeName);
        }
        return new DeferredMessageInput(parseMethod, parserArguments);
    }

    public static class DeferredMessageInput implements MessageInput<Message> {

        private final Method method;
        private final List<String> args;

        public DeferredMessageInput(Method method, List<String> args) {
            this.method = method;
            this.args = args;
        }

        @Override
        @SuppressWarnings({"unchecked"})
        public Message parse(ReadBuffer io) throws BufferException {
            Object[] argValues = new Object[args.size() + 1];
            int index = 0;
            argValues[index++] = io;
            Class<?>[] parameterTypes = method.getParameterTypes();

            if (parameterTypes.length - 1 != args.size()) {
                String type = method.getDeclaringClass().getName();
                int totalArgs = parameterTypes.length - 1;
                int given = args.size();
                throw new BufferException("Invalid parameters detected. Type " + type + " expected " + totalArgs + ", not " + given);
            }

            for (String parameterValue : args) {
                Class<?> parameterType = parameterTypes[index];
                if (parameterType == Boolean.class || parameterType == boolean.class) {
                    argValues[index++] = Boolean.parseBoolean(parameterValue);
                } else if (parameterType == Byte.class || parameterType == byte.class) {
                    argValues[index++] = Byte.parseByte(parameterValue);
                } else if (parameterType == Short.class || parameterType == short.class) {
                    argValues[index++] = Short.parseShort(parameterValue);
                } else if (parameterType == Integer.class || parameterType == int.class) {
                    argValues[index++] = Integer.parseInt(parameterValue);
                } else if (parameterType == Long.class || parameterType == long.class) {
                    argValues[index++] = Long.parseLong(parameterValue);
                } else if (parameterType == Float.class || parameterType == float.class) {
                    argValues[index++] = Float.parseFloat(parameterValue);
                } else if (parameterType == Double.class || parameterType == double.class) {
                    argValues[index++] = Double.parseDouble(parameterValue);
                } else if (parameterType == String.class) {
                    argValues[index++] = parameterValue;
                } else if (Enum.class.isAssignableFrom(parameterType)) {
                    argValues[index++] = Enum.valueOf((Class<? extends Enum>) parameterType, parameterValue);
                } else {
                    throw new BufferException("Currently unsupported parameter type");
                }
            }

            try {
                return (Message) method.invoke(null, argValues);
            } catch (IllegalAccessException | InvocationTargetException e) {
                throw new BufferException("Could not parse payload", e);
            }

        }
    }
}

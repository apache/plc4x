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
package org.apache.plc4x.java.utils.testutils.driver.internal.handlers;

import org.apache.plc4x.java.spi.buffers.api.Message;
import org.apache.plc4x.java.spi.buffers.api.ReadBuffer;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.spi.buffers.xmlbased.ReadBufferXmlBased;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.utils.StaticHelper;
import org.apache.plc4x.java.utils.testutils.driver.internal.utils.ChannelUtil;
import org.apache.plc4x.java.utils.testutils.driver.internal.utils.Delay;
import org.apache.commons.lang3.ClassUtils;
import org.dom4j.Element;
import org.dom4j.QName;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.ByteArrayInputStream;
import java.lang.reflect.*;
import java.util.Arrays;
import java.util.Optional;

/**
 * Handler for injecting incoming PLC messages into the channel.
 * NOTE: Full message parsing/injection not yet implemented.
 * Use incoming-plc-bytes test steps instead for byte-level testing.
 */
public class IncomingPlcMessageHandler {

    private static final Logger LOGGER = LoggerFactory.getLogger(IncomingPlcMessageHandler.class);

    private final String basePackage;
    private final Element messageXml;

    public IncomingPlcMessageHandler(String basePackage, Element messageXml) {
        this.basePackage = basePackage;
        this.messageXml = messageXml;
    }

    /**
     * Executes the incoming PLC message injection.
     *
     * @param transportInstance the transport instance
     * @param byteOrder         the byte order name
     */
    public void executeIncomingPlcMessage(TransportInstance<?> transportInstance, String byteOrder) {
        Delay.shortDelay();

        // Get the name of the root message type
        Optional<Element> messageElementOptional = messageXml.elements().stream().filter(
            e -> !e.getQName().equals(new QName("parser-arguments"))).findFirst();
        if (messageElementOptional.isEmpty()) {
            throw new RuntimeException("No message element found in reference XML");
        }
        Element messageElement = messageElementOptional.get();
        String rootMessageTypeName = messageElement.getName();
        String className = basePackage + "." + rootMessageTypeName;

        // Parse the message.
        // In general, we take the first element, that's not "parser-argument", add that to the
        // package name, resolve the staticParse method, extract its parameters, parse and convert
        // the content in the "parser-arguments" element and invoke ith with that information.
        // The result should be a Message object representing the parsed message.
        Class<?> messageTypeClass;
        try {
            messageTypeClass = Thread.currentThread().getContextClassLoader().loadClass(className);

            // Get the method named "staticParse"
            Method staticParseMethod = Arrays.stream(messageTypeClass.getDeclaredMethods())
                .filter(method -> method.getName().equals("staticParse") &&
                    Modifier.isStatic(method.getModifiers()) &&
                    method.getParameterTypes()[0] == ReadBuffer.class)
                .findFirst()
                .orElseThrow(() -> new RuntimeException("No staticParse method found in class " + className));

            ReadBufferXmlBased readBuffer = new ReadBufferXmlBased(new ByteArrayInputStream(messageElement.asXML().getBytes()));

            Element parserArgumentsElement = messageXml.element(new QName("parser-arguments"));

            // Get the parameter types for the staticParseMethod
            int numArgs = staticParseMethod.getParameters().length;
            Object[] argValues = new Object[numArgs];
            for (int i = 0; i < numArgs; i++) {
                Parameter parameter = staticParseMethod.getParameters()[i];
                if (i == 0) {
                    if (!parameter.getName().equals("readBuffer")) {
                        throw new RuntimeException("Invalid parameter name for first argument. Expected 'readBuffer' but got: " + parameter.getName());
                    }
                    if (!parameter.getType().equals(ReadBuffer.class)) {
                        throw new RuntimeException("Invalid parameter type for first argument. Expected 'ReadBuffer' but got: " + parameter.getType());
                    }
                    argValues[i] = readBuffer;
                } else {
                    String parameterName = parameter.getName();
                    Class<?> parameterType = parameter.getType();
                    if ((parserArgumentsElement == null) || (parserArgumentsElement.element(new QName(parameterName)) == null)) {
                        throw new RuntimeException("No parser-arguments element or parameterName element found for parameter " + parameterName);
                    } else {
                        String parameterStringValue = parserArgumentsElement.element(new QName(parameterName)).getTextTrim();
                        argValues[i] = parseDynamic(parameterType, parameterStringValue);
                    }
                }
            }

            // Actually parse the message
            Object parsed = staticParseMethod.invoke(null, argValues);
            LOGGER.debug("Parsed message: {}", parsed);

            if (!(parsed instanceof Message message)) {
                throw new RuntimeException("Parsed message is not an instance of Message");
            }

            try {
                // The default integer/float encodings match what every driver codec configures;
                // generated serializers (e.g. S7) fail without them on fields that don't pass
                // explicit per-field options.
                WriteBufferByteBased writeBuffer = new WriteBufferByteBased(new byte[message.getLengthInBytes()],
                    WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
                    WithOption.WithSignedIntegerEncoding("twos-complement"),
                    WithOption.WithFloatEncoding("IEEE754"));
                writeBuffer.writeMessage(message);
                byte[] bytes = writeBuffer.getBytes();
                LOGGER.info("Sending serialized message: {}", StaticHelper.ENCODE_HEX(bytes));

                // Send out the bytes to the transport
                ChannelUtil.writeInboundBytes(transportInstance, bytes);
            } catch (BufferException e) {
                throw new RuntimeException("Failed to serialize message", e);
            }

        } catch (ClassNotFoundException e) {
            throw new RuntimeException("Failed to load class " + className, e);
        } catch (InvocationTargetException e) {
            throw new RuntimeException("Failed to invoke staticParse method in class " + className, e.getCause() != null ? e.getCause() : e);
        } catch (IllegalAccessException e) {
            throw new RuntimeException("Failed to access staticParse method in class " + className, e);
        }
    }


    public static Object parseDynamic(Class<?> type, String value) {
        try {
            // Try valueOf(String)
            try {
                Method m = type.getMethod("valueOf", String.class);
                if (Modifier.isStatic(m.getModifiers())) {
                    return m.invoke(null, value);
                }
            } catch (NoSuchMethodException ignored) {}

            // If the type is a primitive, get the corresponding object-type.
            if (type.isPrimitive()) {
                type = ClassUtils.primitiveToWrapper(type);
            }

            // Try parse*(String)
            for (Method m : type.getMethods()) {
                if (Modifier.isStatic(m.getModifiers())
                    && m.getParameterCount() == 1
                    && m.getParameterTypes()[0] == String.class
                    && m.getName().startsWith("parse")) {
                    return m.invoke(null, value);
                }
            }

            // String constructor
            try {
                Constructor<?> c = type.getConstructor(String.class);
                return c.newInstance(value);
            } catch (NoSuchMethodException ignored) {}

            throw new IllegalArgumentException(
                "Don't know how to parse type: " + type.getName()
            );

        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

}

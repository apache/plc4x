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
package org.apache.plc4x.java.ads.api.serial;

import org.apache.plc4x.java.ads.api.util.ByteReadable;
import org.junit.runners.Parameterized;

import java.lang.reflect.Method;
import java.lang.reflect.Modifier;
import java.util.Arrays;
import java.util.Collection;
import java.util.List;
import java.util.stream.Collectors;
import java.util.stream.Stream;

import static org.hamcrest.Matchers.*;
import static org.junit.Assert.assertThat;
import static org.mockito.Mockito.*;

//@RunWith(Parameterized.class)
public class SerialFactoryMethodTest {

    @Parameterized.Parameter
    public Class<?> clazz;

    @Parameterized.Parameters(name = "{index} {0}")
    public static Collection<Object[]> data() {
        return Stream.of(
            AmsSerialAcknowledgeFrame.class,
            AmsSerialFrame.class,
            AmsSerialResetFrame.class
        ).map(clazz -> new Object[]{clazz}).collect(Collectors.toList());
    }

    // TODO: Commented out as it was causing problems with Java 11
    //@Test
    public void testOf() throws Exception {
        List<Method> getters = Arrays
            .stream(clazz.getDeclaredMethods())
            .filter(method -> (
                method.getName().startsWith("get") || method.getName().startsWith("is"))
                && Modifier.isPublic(method.getModifiers())
                && method.getParameterCount() == 0)
            .collect(Collectors.toList());
        for (Method method : clazz.getDeclaredMethods()) {
            if (!method.getName().equals("of")) {
                continue;
            }
            Object invoke = method.invoke(null, Arrays.stream(method.getParameterTypes()).map(aClass -> {
                if (ByteReadable.class.isAssignableFrom(aClass)) {
                    ByteReadable mock = (ByteReadable) mock(aClass, RETURNS_DEEP_STUBS);
                    when(mock.getBytes()).thenReturn(new byte[0]);
                    return mock;
                } else {
                    return mock(aClass, RETURNS_DEEP_STUBS);
                }
            }).toArray());
            assertThat(invoke, notNullValue());
            assertThat(invoke, instanceOf(clazz));
            assertThat(invoke.toString(), not(isEmptyString()));
            // Testing getters for the coverage (sonar)
            for (Method getter : getters) {
                getter.invoke(invoke);
            }
        }
    }
}

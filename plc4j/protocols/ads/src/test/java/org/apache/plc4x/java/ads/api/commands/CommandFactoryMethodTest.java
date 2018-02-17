/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */
package org.apache.plc4x.java.ads.api.commands;

import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.Parameterized;

import java.lang.reflect.Method;
import java.util.Arrays;
import java.util.Collection;
import java.util.stream.Collectors;
import java.util.stream.Stream;

import static org.mockito.Mockito.RETURNS_DEEP_STUBS;
import static org.mockito.Mockito.mock;

@RunWith(Parameterized.class)
public class CommandFactoryMethodTest {

    @Parameterized.Parameter
    public Class<?> clazz;

    @Parameterized.Parameters(name = "{index} {0}")
    public static Collection<Object[]> data() {
        return Stream.of(
            ADSAbstractRequest.class,
            ADSAbstractResponse.class,
            ADSAddDeviceNotificationRequest.class,
            ADSAddDeviceNotificationResponse.class,
            ADSDeleteDeviceNotificationRequest.class,
            ADSDeleteDeviceNotificationResponse.class,
            ADSDeviceNotificationRequest.class,
            ADSDeviceNotificationResponse.class,
            ADSReadDeviceInfoRequest.class,
            ADSReadDeviceInfoResponse.class,
            ADSReadRequest.class,
            ADSReadResponse.class,
            ADSReadStateRequest.class,
            ADSReadStateResponse.class,
            ADSReadWriteRequest.class,
            ADSReadWriteResponse.class,
            ADSWriteControlRequest.class,
            ADSWriteControlResponse.class,
            ADSWriteRequest.class,
            ADSWriteResponse.class,
            UnknownCommand.class
        ).map(clazz -> new Object[]{clazz}).collect(Collectors.toList());
    }

    @Test
    public void testOf() throws Exception {
        for (Method method : clazz.getDeclaredMethods()) {
            if (method.getName().equals("of")) {
                method.invoke(null, Arrays.stream(method.getParameterTypes()).map(aClass -> mock(aClass, RETURNS_DEEP_STUBS)).toArray());
            }
        }
    }
}

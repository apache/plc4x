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
package org.apache.plc4x.java.ads.api.commands.types;

import org.apache.plc4x.java.ads.api.util.UnsignedIntLEByteValue;
import org.apache.plc4x.java.ads.api.util.UnsignedShortLEByteValue;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.Parameterized;

import java.lang.reflect.Method;
import java.util.Collection;
import java.util.stream.Collectors;
import java.util.stream.Stream;

import static org.hamcrest.Matchers.instanceOf;
import static org.junit.Assume.assumeThat;

@RunWith(Parameterized.class)
public class CommandTypesFactoryMethodTest {

    @Parameterized.Parameter
    public Class<?> clazz;

    @Parameterized.Parameters(name = "{index} {0}")
    public static Collection<Object[]> data() {
        return Stream.of(
            ADSState.class,
            CycleTime.class,
            Data.class,
            Device.class,
            DeviceState.class,
            IndexGroup.class,
            IndexOffset.class,
            Length.class,
            MajorVersion.class,
            MaxDelay.class,
            MinorVersion.class,
            NotificationHandle.class,
            ReadLength.class,
            Result.class,
            Samples.class,
            SampleSize.class,
            Stamps.class,
            TimeStamp.class,
            TransmissionMode.class,
            Version.class,
            WriteLength.class
        ).map(clazz -> new Object[]{clazz}).collect(Collectors.toList());
    }

    @Test
    public void testOfInt() throws Exception {
        assumeThat(clazz, instanceOf(UnsignedShortLEByteValue.class));
        Method ofMethod = clazz.getDeclaredMethod("of", int.class);
        ofMethod.invoke(null, 1);
    }

    @Test
    public void testOfLong() throws Exception {
        assumeThat(clazz, instanceOf(UnsignedIntLEByteValue.class));
        Method ofMethod = clazz.getDeclaredMethod("of", long.class);
        ofMethod.invoke(null, 1L);
    }

    @Test
    public void testOfString() throws Exception {
        Method ofMethod = clazz.getDeclaredMethod("of", String.class);
        ofMethod.invoke(null, "1");
    }
}

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
package org.apache.plc4x.java.ads.api.commands.types;

import org.apache.plc4x.java.ads.api.util.UnsignedIntLEByteValue;
import org.apache.plc4x.java.ads.api.util.UnsignedShortLEByteValue;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.Parameterized;

import java.lang.reflect.Field;
import java.lang.reflect.Method;
import java.lang.reflect.Modifier;
import java.math.BigInteger;
import java.nio.charset.Charset;
import java.util.Arrays;
import java.util.Collection;
import java.util.List;
import java.util.stream.Collectors;
import java.util.stream.Stream;

import static org.hamcrest.Matchers.*;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assume.assumeThat;

@RunWith(Parameterized.class)
public class CommandTypesFactoryMethodTest {

    @Parameterized.Parameter
    public Class<?> clazz;

    @Parameterized.Parameters(name = "{index} {0}")
    public static Collection<Object[]> data() {
        return Stream.of(
            AdsReturnCode.class,
            AdsStampHeader.class,
            AdsState.class,
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
    public void innerReservedClasses() {
        assertNotNull(IndexGroup.SystemServiceGroups.SYSTEMSERVICE_CHANGENETID);
        assertNotNull(AdsState.DefinedValues.ADSSTATE_CONFIG);
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
        assumeThat(clazz, not(AdsStampHeader.class));
        Method ofMethod = clazz.getDeclaredMethod("of", String.class);
        ofMethod.invoke(null, clazz != AdsReturnCode.class ? "1" : "ADS_CODE_0");
    }

    @Test
    public void testGetter() throws Exception {
        assumeThat(clazz, not(AdsStampHeader.class));
        List<Method> getters = Arrays
            .stream(clazz.getDeclaredMethods())
            .filter(method -> (
                method.getName().startsWith("get") || method.getName().startsWith("is"))
                && Modifier.isPublic(method.getModifiers())
                && method.getParameterCount() == 0)
            .collect(Collectors.toList());
        Method ofMethod = clazz.getDeclaredMethod("of", String.class);
        Object invoke = ofMethod.invoke(null, clazz != AdsReturnCode.class ? "1" : "ADS_CODE_0");
        // Testing getters for the coverage (sonar)
        for (Method getter : getters) {
            getter.invoke(invoke);
        }
    }

    @Test
    public void testOfStringCharset() throws Exception {
        assumeThat(clazz, isOneOf(Device.class, Data.class));
        Method ofMethod = clazz.getDeclaredMethod("of", String.class, Charset.class);
        ofMethod.invoke(null, "1", Charset.defaultCharset());
    }

    @Test
    public void testOfWintime() throws Exception {
        assumeThat(clazz, isOneOf(TimeStamp.class));
        {
            Method ofMethod = clazz.getDeclaredMethod("ofWinTime", BigInteger.class);
            ofMethod.invoke(null, BigInteger.valueOf(1));
        }
        {
            Method ofMethod = clazz.getDeclaredMethod("ofWinTime", String.class);
            ofMethod.invoke(null, "1");
        }
        {
            Method ofMethod = clazz.getDeclaredMethod("ofWinTime", long.class);
            ofMethod.invoke(null, 1L);
        }
    }

    @Test
    public void testOfBytes() throws Exception {
        assumeThat(clazz, not(AdsStampHeader.class));
        assumeThat(clazz, not(Data.class));
        assumeThat(clazz, not(AdsReturnCode.class));
        Field num_bytes_field = clazz.getDeclaredField("NUM_BYTES");
        Integer numberOfBytes = (Integer) num_bytes_field.get(null);
        Method ofMethod = clazz.getDeclaredMethod("of", byte[].class);
        ofMethod.invoke(null, (Object) new byte[numberOfBytes]);
    }
}

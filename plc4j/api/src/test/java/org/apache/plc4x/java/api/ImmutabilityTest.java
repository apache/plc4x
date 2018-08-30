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
package org.apache.plc4x.java.api;

import org.apache.plc4x.java.api.authentication.PlcUsernamePasswordAuthentication;
import org.apache.plc4x.java.api.messages.*;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.Parameterized;

import java.util.Collection;
import java.util.Set;
import java.util.stream.Collectors;
import java.util.stream.Stream;

import static org.hamcrest.Matchers.hasItem;
import static org.hamcrest.Matchers.not;
import static org.junit.Assume.assumeThat;
import static org.mutabilitydetector.unittesting.AllowedReason.allowingForSubclassing;
import static org.mutabilitydetector.unittesting.MutabilityAssert.assertInstancesOf;
import static org.mutabilitydetector.unittesting.MutabilityMatchers.areImmutable;

@RunWith(Parameterized.class)
public class ImmutabilityTest {

    private static Set<Class<?>> NOT_YET_IMMUTABLE = Stream.of(
        /*PlcReadResponseItem.class,
        SubscriptionEventItem.class,
        SubscriptionRequestCyclicItem.class,
        SubscriptionRequestItem.class,
        SubscriptionResponseItem.class,
        UnsubscriptionRequestItem.class,
        PlcWriteRequestItem.class,
        PlcProprietaryRequest.class,
        PlcProprietaryResponse.class,*/
        PlcSubscriptionEvent.class,
        PlcUnsubscriptionRequest.class
    ).collect(Collectors.toSet());

    @Parameterized.Parameter
    public Class<?> clazz;

    @Parameterized.Parameters(name = "{index} {0}")
    public static Collection<Object[]> data() {
        return Stream.of(
            PlcUsernamePasswordAuthentication.class,
            /*PlcReadRequestItem.class,
            PlcReadResponseItem.class,
            SubscriptionEventItem.class,
            SubscriptionRequestChangeOfStateItem.class,
            SubscriptionRequestCyclicItem.class,
            SubscriptionRequestEventItem.class,
            SubscriptionRequestItem.class,
            SubscriptionResponseItem.class,
            UnsubscriptionRequestItem.class,
            PlcWriteRequestItem.class,
            PlcWriteResponseItem.class,
            TypeSafePlcReadRequest.class,
            TypeSafePlcReadResponse.class,
            TypeSafePlcWriteRequest.class,
            TypeSafePlcWriteResponse.class,
            PlcProprietaryRequest.class,
            PlcProprietaryResponse.class,*/
            PlcSubscriptionEvent.class,
            PlcSubscriptionRequest.class,
            PlcSubscriptionResponse.class,
            PlcUnsubscriptionRequest.class,
            PlcUnsubscriptionResponse.class,
            PlcWriteRequest.class,
            PlcWriteResponse.class
        ).map(aClass -> new Object[]{aClass})
            .collect(Collectors.toList());
    }

    @Test
    public void immutability() {
        assumeThat(clazz + " not yet immutable", NOT_YET_IMMUTABLE, not(hasItem(clazz)));
        assertInstancesOf(clazz,
            areImmutable(),
            allowingForSubclassing());
    }
}

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

package org.apache.plc4x.java.base.messages;

import org.apache.plc4x.java.api.messages.PlcUnsubscriptionRequest;
import org.apache.plc4x.java.base.model.InternalPlcSubscriptionHandle;
import org.junit.Assert;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.Arrays;
import java.util.Collection;
import java.util.Collections;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.core.IsEqual.equalTo;
import static org.hamcrest.core.IsNull.notNullValue;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class DefaultPlcUnsubscriptionRequestTest {

    @Mock
    private PlcSubscriber subscriber;
    @Mock
    private InternalPlcSubscriptionHandle subscriptionHandle;
    private DefaultPlcUnsubscriptionRequest SUT;

    @BeforeEach
    void setUp() {
        SUT = new DefaultPlcUnsubscriptionRequest(subscriber, Collections.singletonList(subscriptionHandle));
    }

    @Test
    void execute() {
        SUT.execute();
        verify(subscriber, times(1)).unsubscribe(any());
    }

    @Test
    void getInternalPlcSubscriptionHandles() {
        Collection<? extends InternalPlcSubscriptionHandle> handles = SUT.getInternalPlcSubscriptionHandles();
        assertThat(handles, notNullValue());
        assertThat(handles.size(), equalTo(1));
        assertThat(handles.iterator().next(), equalTo(subscriptionHandle));
    }

    @Test
    void builder() {
        DefaultPlcUnsubscriptionRequest.Builder builder = new DefaultPlcUnsubscriptionRequest.Builder(subscriber);
        builder.addHandles(subscriptionHandle);
        PlcUnsubscriptionRequest unsubscriptionRequest = builder.build();
        Assert.assertThat(unsubscriptionRequest, notNullValue());

        InternalPlcSubscriptionHandle secondHandle = mock(InternalPlcSubscriptionHandle.class);
        InternalPlcSubscriptionHandle thirdHandle = mock(InternalPlcSubscriptionHandle.class);

        builder = new DefaultPlcUnsubscriptionRequest.Builder(subscriber);
        builder.addHandles(subscriptionHandle, secondHandle, thirdHandle);
        unsubscriptionRequest = builder.build();
        Assert.assertThat(unsubscriptionRequest, notNullValue());

        builder = new DefaultPlcUnsubscriptionRequest.Builder(subscriber);
        builder.addHandles(Arrays.asList(subscriptionHandle, secondHandle, thirdHandle));
        unsubscriptionRequest = builder.build();
        Assert.assertThat(unsubscriptionRequest, notNullValue());
    }


}
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

package org.apache.plc4x.java.base.model;

import nl.jqno.equalsverifier.EqualsVerifier;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.base.messages.PlcSubscriber;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.Collection;
import java.util.function.Consumer;

import static org.hamcrest.core.IsEqual.equalTo;
import static org.hamcrest.core.IsNull.notNullValue;
import static org.junit.Assert.assertThat;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.times;
import static org.mockito.Mockito.verify;

@ExtendWith(MockitoExtension.class)
class DefaultPlcConsumerRegistrationTest {

    @Mock
    private PlcSubscriber subscriber;
    @Mock
    private Consumer<PlcSubscriptionEvent> consumer;
    @Mock
    private InternalPlcSubscriptionHandle subscriptionHandle;
    private DefaultPlcConsumerRegistration SUT;

    @BeforeEach
    void setUp() {
        SUT = new DefaultPlcConsumerRegistration(subscriber, consumer, subscriptionHandle);
    }

    @Test
    void getConsumerHash() {
        int consumerHash = SUT.getConsumerHash();
        assertThat(consumerHash, equalTo(consumer.hashCode()));
    }

    @Test
    void getAssociatedHandles() {
        Collection<InternalPlcSubscriptionHandle> associatedHandles = SUT.getAssociatedHandles();
        assertThat(associatedHandles, notNullValue());
        assertThat(associatedHandles.size(), equalTo(1));
        assertThat(associatedHandles.iterator().next(), equalTo(subscriptionHandle));
    }

    @Test
    public void unregister() {
        SUT.unregister();
        verify(subscriber, times(1)).unregister(any());
    }

    @Test
    public void testToString() {
        String s = SUT.toString();
        assertThat(s, notNullValue());
    }

    @Test
    void equals() {
        EqualsVerifier.forClass(DefaultPlcConsumerRegistration.class).verify();
    }

}
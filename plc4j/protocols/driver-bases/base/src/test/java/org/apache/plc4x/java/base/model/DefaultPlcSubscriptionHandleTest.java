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
import nl.jqno.equalsverifier.Warning;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.base.messages.PlcSubscriber;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.Collections;
import java.util.function.Consumer;

import static org.hamcrest.core.IsNull.notNullValue;
import static org.junit.Assert.assertThat;
import static org.mockito.Mockito.times;
import static org.mockito.Mockito.verify;

@ExtendWith(MockitoExtension.class)
class DefaultPlcSubscriptionHandleTest {

    @Mock
    private PlcSubscriber subscriber;
    @Mock
    private Consumer<PlcSubscriptionEvent> consumer;
    private DefaultPlcSubscriptionHandle SUT;

    @BeforeEach
    void setUp() {
        SUT = new DefaultPlcSubscriptionHandle(subscriber);
    }

    @Test
    void register() {
        SUT.register(consumer);
        verify(subscriber, times(1)).register(consumer, Collections.singletonList(SUT));
    }

    @Test
    public void testToString() {
        String s = SUT.toString();
        assertThat(s, notNullValue());
    }

    @Test
    void equals() {
        EqualsVerifier.forClass(DefaultPlcSubscriptionHandle.class).usingGetClass()
            .suppress(Warning.STRICT_INHERITANCE).suppress(Warning.IDENTICAL_COPY)
            .suppress(Warning.ALL_FIELDS_SHOULD_BE_USED).verify();
    }

}
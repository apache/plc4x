/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.s7.userdata;

import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.s7.readwrite.MemoryArea;
import org.apache.plc4x.java.s7.readwrite.TransportSize;
import org.apache.plc4x.java.s7.tag.S7Tag;
import org.apache.plc4x.java.spi.drivers.functions.PlcSubscriber;
import org.junit.jupiter.api.Test;
import org.mockito.ArgumentCaptor;

import java.util.Collection;
import java.util.List;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyList;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

class S7SubscriptionHandleTest {

    @Test
    void alarmSubscriptionHandle_exposesTagAndDelegatesRegister() {
        PlcSubscriber subscriber = mock(PlcSubscriber.class);
        PlcConsumerRegistration reg = mock(PlcConsumerRegistration.class);
        S7AlarmSubscriptionHandle handle = new S7AlarmSubscriptionHandle(subscriber, "alarms");

        when(subscriber.registerConsumer(any(), any())).thenReturn(reg);

        assertThat(handle.getTagName()).isEqualTo("alarms");
        assertThat(handle.toString()).contains("alarms");

        PlcConsumerRegistration result = handle.register(evt -> {});
        assertThat(result).isSameAs(reg);

        @SuppressWarnings({"unchecked", "rawtypes"})
        ArgumentCaptor<List> handles = ArgumentCaptor.forClass(List.class);
        verify(subscriber).registerConsumer(any(), handles.capture());
        assertThat(handles.getValue()).containsExactly(handle);
    }

    @Test
    void alarmSubscriptionHandle_rejectsNullArgs() {
        PlcSubscriber subscriber = mock(PlcSubscriber.class);
        assertThatThrownBy(() -> new S7AlarmSubscriptionHandle(null, "x"))
            .isInstanceOf(NullPointerException.class);
        assertThatThrownBy(() -> new S7AlarmSubscriptionHandle(subscriber, null))
            .isInstanceOf(NullPointerException.class);
    }

    @Test
    void alarmQueryHandle_keepsPayloadAndReportsLengthInToString() {
        PlcSubscriber subscriber = mock(PlcSubscriber.class);
        byte[] payload = new byte[]{1, 2, 3, 4};
        S7AlarmQueryHandle handle = new S7AlarmQueryHandle(subscriber, "buffered", payload);

        assertThat(handle.getTagName()).isEqualTo("buffered");
        assertThat(handle.getPayload()).containsExactly((byte) 1, (byte) 2, (byte) 3, (byte) 4);
        assertThat(handle.toString()).contains("buffered").contains("4 bytes");
    }

    @Test
    void alarmQueryHandle_treatsNullPayloadAsEmpty() {
        PlcSubscriber subscriber = mock(PlcSubscriber.class);
        S7AlarmQueryHandle handle = new S7AlarmQueryHandle(subscriber, "buffered", null);
        assertThat(handle.getPayload()).isEmpty();
        assertThat(handle.toString()).contains("0 bytes");
    }

    @Test
    void alarmQueryHandle_register_delegatesToSubscriber() {
        PlcSubscriber subscriber = mock(PlcSubscriber.class);
        PlcConsumerRegistration reg = mock(PlcConsumerRegistration.class);
        when(subscriber.registerConsumer(any(), any())).thenReturn(reg);

        S7AlarmQueryHandle handle = new S7AlarmQueryHandle(subscriber, "buffered", new byte[]{0});
        assertThat(handle.register(e -> {})).isSameAs(reg);
        verify(subscriber).registerConsumer(any(), anyList());
    }

    @Test
    void cyclicSubscriptionHandle_exposesAllFieldsAndRegisters() {
        PlcSubscriber subscriber = mock(PlcSubscriber.class);
        PlcConsumerRegistration reg = mock(PlcConsumerRegistration.class);
        when(subscriber.registerConsumer(any(), any())).thenReturn(reg);

        S7Tag tag = new S7Tag(TransportSize.WORD, MemoryArea.FLAGS_MARKERS, 0, 10, (byte) 0, 1);
        S7CyclicSubscriptionHandle handle =
            new S7CyclicSubscriptionHandle(subscriber, "mw10", tag, (short) 7, 2);

        assertThat(handle.getTagName()).isEqualTo("mw10");
        assertThat(handle.getTag()).isSameAs(tag);
        assertThat(handle.getJobId()).isEqualTo((short) 7);
        assertThat(handle.getItemIndex()).isEqualTo(2);
        assertThat(handle.toString()).contains("mw10").contains("job=7").contains("idx=2");

        assertThat(handle.register(e -> {})).isSameAs(reg);
        verify(subscriber).registerConsumer(any(), anyList());
    }

    @Test
    void cyclicSubscriptionHandle_rejectsNullArgs() {
        PlcSubscriber subscriber = mock(PlcSubscriber.class);
        S7Tag tag = new S7Tag(TransportSize.WORD, MemoryArea.FLAGS_MARKERS, 0, 10, (byte) 0, 1);
        assertThatThrownBy(() -> new S7CyclicSubscriptionHandle(null, "x", tag, (short) 0, 0))
            .isInstanceOf(NullPointerException.class);
        assertThatThrownBy(() -> new S7CyclicSubscriptionHandle(subscriber, null, tag, (short) 0, 0))
            .isInstanceOf(NullPointerException.class);
        assertThatThrownBy(() -> new S7CyclicSubscriptionHandle(subscriber, "x", null, (short) 0, 0))
            .isInstanceOf(NullPointerException.class);
    }

}

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
package org.apache.plc4x.java.canopen;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.canopen.configuration.CANOpenConfiguration;
import org.apache.plc4x.java.canopen.conversation.ConversationPredicates;
import org.apache.plc4x.java.canopen.readwrite.CANOpenFrame;
import org.apache.plc4x.java.canopen.readwrite.CANOpenSDOResponse;
import org.apache.plc4x.java.canopen.readwrite.CANOpenService;
import org.apache.plc4x.java.canopen.readwrite.SDOInitiateExpeditedUploadResponse;
import org.apache.plc4x.java.canopen.readwrite.SDOInitiateUploadResponsePayload;
import org.apache.plc4x.java.canopen.readwrite.utils.StaticHelper;
import org.apache.plc4x.java.canopen.tag.CANOpenHeartbeatTag;
import org.apache.plc4x.java.canopen.tag.CANOpenTag;
import org.apache.plc4x.java.canopen.tag.CANOpenTagHandler;
import org.apache.plc4x.java.canopen.transport.CANOpenAbortException;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

class CANOpenSmallSurfaceTest {

    @Test
    void configAccessorsRoundtrip() {
        CANOpenConfiguration cfg = new CANOpenConfiguration();
        cfg.setNodeId(7);
        cfg.setHeartbeat(true);
        cfg.setRequestTimeout(5000);
        assertThat(cfg.getNodeId()).isEqualTo(7);
        assertThat(cfg.isHeartbeat()).isTrue();
        assertThat(cfg.getRequestTimeout()).isEqualTo(5000);
    }

    @Test
    void tagHandlerDispatchesToConcreteTagTypes() {
        CANOpenTagHandler handler = new CANOpenTagHandler();
        assertThat(handler.parseTag("HEARTBEAT")).isInstanceOf(CANOpenHeartbeatTag.class);
    }

    @Test
    void tagHandlerDoesNotSupportBrowsing() {
        assertThatThrownBy(() -> new CANOpenTagHandler().parseQuery("anything"))
            .isInstanceOf(UnsupportedOperationException.class);
    }

    @Test
    void tagOfThrowsForUnknownAddress() {
        assertThatThrownBy(() -> CANOpenTag.of("not-a-canopen-address"))
            .isInstanceOf(PlcInvalidTagException.class);
    }

    @Test
    void heartbeatTagAccessors() {
        CANOpenHeartbeatTag tag = CANOpenHeartbeatTag.of("HEARTBEAT");
        assertThat(tag.getNodeId()).isZero();
        assertThat(tag.isWildcard()).isTrue();
        assertThat(tag.getService()).isEqualTo(CANOpenService.HEARTBEAT);
        assertThat(tag.getPlcValueType()).isEqualTo(PlcValueType.NULL);
        assertThat(tag.getArrayInfo()).isNotNull();
        assertThat(tag.getAddressString()).isEqualTo("HEARTBEAT:0");
    }

    @Test
    void heartbeatTagMatchesHelper() {
        assertThat(CANOpenHeartbeatTag.matches("HEARTBEAT")).isTrue();
        assertThat(CANOpenHeartbeatTag.matches("not-heartbeat")).isFalse();
    }

    @Test
    void heartbeatGetMatcherThrowsOnUnmatchedInput() {
        assertThatThrownBy(() -> CANOpenHeartbeatTag.getMatcher("nope"))
            .isInstanceOf(PlcInvalidTagException.class);
    }

    @Test
    void abortExceptionCarriesAbortCode() {
        CANOpenAbortException byMessage = new CANOpenAbortException("boom", 0x05040001L);
        assertThat(byMessage.getMessage()).isEqualTo("boom");
        assertThat(byMessage.getAbortCode()).isEqualTo(0x05040001L);

        Throwable cause = new IllegalStateException("inner");
        CANOpenAbortException byCause = new CANOpenAbortException(cause, 0x05030000L);
        assertThat(byCause.getCause()).isSameAs(cause);
        assertThat(byCause.getAbortCode()).isEqualTo(0x05030000L);
    }

    @Test
    void conversationPredicateMatchesOnlySdoTransmitFromTargetNode() {
        CANOpenFrame match = Mockito.mock(CANOpenFrame.class);
        Mockito.when(match.getNodeId()).thenReturn((short) 7);
        Mockito.when(match.getService()).thenReturn(CANOpenService.TRANSMIT_SDO);
        Mockito.when(match.getPayload()).thenReturn(Mockito.mock(CANOpenSDOResponse.class));

        CANOpenFrame wrongNode = Mockito.mock(CANOpenFrame.class);
        Mockito.when(wrongNode.getNodeId()).thenReturn((short) 8);
        Mockito.when(wrongNode.getService()).thenReturn(CANOpenService.TRANSMIT_SDO);
        Mockito.when(wrongNode.getPayload()).thenReturn(Mockito.mock(CANOpenSDOResponse.class));

        CANOpenFrame wrongService = Mockito.mock(CANOpenFrame.class);
        Mockito.when(wrongService.getNodeId()).thenReturn((short) 7);
        Mockito.when(wrongService.getService()).thenReturn(CANOpenService.HEARTBEAT);

        var predicate = ConversationPredicates.sdoTransmitFrom(7);
        assertThat(predicate.test(match)).isTrue();
        assertThat(predicate.test(wrongNode)).isFalse();
        assertThat(predicate.test(wrongService)).isFalse();
    }

    @Test
    void staticHelperCountReturnsUnusedBytesForExpeditedIndicatedResponse() {
        SDOInitiateExpeditedUploadResponse payload = Mockito.mock(SDOInitiateExpeditedUploadResponse.class);
        Mockito.when(payload.getData()).thenReturn(new byte[]{1, 2});
        assertThat(StaticHelper.count(true, true, payload)).isEqualTo(2);
    }

    @Test
    void staticHelperCountIsZeroForNonExpeditedOrUnindicatedResponses() {
        SDOInitiateUploadResponsePayload payload = Mockito.mock(SDOInitiateUploadResponsePayload.class);
        assertThat(StaticHelper.count(false, true, payload)).isZero();
        assertThat(StaticHelper.count(true, false, payload)).isZero();
        // Indicated+expedited but non-expedited subtype also returns 0.
        assertThat(StaticHelper.count(true, true, payload)).isZero();
    }
}

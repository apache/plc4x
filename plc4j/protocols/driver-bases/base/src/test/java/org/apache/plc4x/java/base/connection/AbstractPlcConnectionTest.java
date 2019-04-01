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

package org.apache.plc4x.java.base.connection;

import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.api.exceptions.PlcUnsupportedOperationException;
import org.apache.plc4x.java.base.messages.DefaultPlcReadRequest;
import org.apache.plc4x.java.base.messages.PlcReader;
import org.assertj.core.api.WithAssertions;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.Mockito.mock;

class AbstractPlcConnectionTest implements WithAssertions {

    AbstractPlcConnection SUT = new AbstractPlcConnection() {
        @Override
        public void connect() {
            throw new NotImplementedException("not used");
        }

        @Override
        public boolean isConnected() {
            throw new NotImplementedException("not used");
        }

        @Override
        public void close() {
            throw new NotImplementedException("not used");
        }
    };

    @Test
    void getMetadata() {
        assertThat(SUT.getMetadata()).isNotNull().isSameAs(SUT);
    }

    @Test
    void canRead() {
        assertThat(SUT.canRead()).isFalse();
    }

    @Test
    void canWrite() {
        assertThat(SUT.canWrite()).isFalse();
    }

    @Test
    void canSubscribe() {
        assertThat(SUT.canSubscribe()).isFalse();
    }

    @Test
    void readRequestBuilder() {
        assertThatThrownBy(() -> SUT.readRequestBuilder()).isInstanceOf(PlcUnsupportedOperationException.class);
    }

    @Test
    void writeRequestBuilder() {
        assertThatThrownBy(() -> SUT.writeRequestBuilder()).isInstanceOf(PlcUnsupportedOperationException.class);
    }

    @Test
    void subscriptionRequestBuilder() {
        assertThatThrownBy(() -> SUT.subscriptionRequestBuilder()).isInstanceOf(PlcUnsupportedOperationException.class);
    }

    @Test
    void unsubscriptionRequestBuilder() {
        assertThatThrownBy(() -> SUT.unsubscriptionRequestBuilder()).isInstanceOf(PlcUnsupportedOperationException.class);
    }

    @Test
    void checkInternalTest() {
        assertThrows(IllegalArgumentException.class, () -> SUT.checkInternal("Test", DefaultPlcReadRequest.class));
        DefaultPlcReadRequest readRequest = SUT.checkInternal(
            new DefaultPlcReadRequest.Builder(mock(PlcReader.class), mock(PlcFieldHandler.class)).build(),
            DefaultPlcReadRequest.class);
        assertThat(readRequest).isNotNull();
    }


}
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

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.junit.jupiter.MockitoExtension;

import static org.hamcrest.core.IsEqual.equalTo;
import static org.hamcrest.core.IsNull.nullValue;
import static org.junit.Assert.assertThat;

@ExtendWith(MockitoExtension.class)
class PlcRawMessageTest {

    private ByteBuf rawData;
    private PlcRawMessage SUT;

    @BeforeEach
    void setUp() {
        rawData = Unpooled.wrappedBuffer("foo".getBytes());
        SUT = new PlcRawMessage(rawData);
    }

    @Test
    void getUserData() {
        ByteBuf userData = SUT.getUserData();
        assertThat(userData, equalTo(rawData));
    }

    @Test
    void getParent() {
        PlcProtocolMessage parent = SUT.getParent();
        assertThat(parent, nullValue());
    }

}
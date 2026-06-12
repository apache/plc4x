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
package org.apache.plc4x.java.s7.userdata;

import org.apache.plc4x.java.s7.readwrite.S7Message;
import org.apache.plc4x.java.s7.readwrite.S7MessageRequest;
import org.apache.plc4x.java.s7.readwrite.S7MessageUserData;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserData;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserDataItem;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserDataItemCPUFunctions;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserData;
import org.junit.jupiter.api.Test;

import java.io.Closeable;
import java.io.IOException;
import java.util.Collections;
import java.util.concurrent.atomic.AtomicInteger;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;

class S7UserDataPushDispatcherTest {

    @Test
    void dispatch_invokesMatchingHandler() {
        S7UserDataPushDispatcher dispatcher = new S7UserDataPushDispatcher();
        AtomicInteger calls = new AtomicInteger();
        dispatcher.register(0x02, 0x00, 0x05, msg -> calls.incrementAndGet());

        assertTrue(dispatcher.dispatch(pushMessage(0x02, 0x00, 0x05)));
        assertEquals(1, calls.get());
    }

    @Test
    void dispatch_returnsFalseForUnregisteredTriple() {
        S7UserDataPushDispatcher dispatcher = new S7UserDataPushDispatcher();
        dispatcher.register(0x02, 0x00, 0x05, msg -> {});

        assertFalse(dispatcher.dispatch(pushMessage(0x07, 0x00, 0x01)));
    }

    @Test
    void dispatch_returnsFalseForNonUserDataMessage() {
        S7UserDataPushDispatcher dispatcher = new S7UserDataPushDispatcher();
        dispatcher.register(0x02, 0x00, 0x05, msg -> {});

        S7Message reqMessage = new S7MessageRequest(1, null, null);
        assertFalse(dispatcher.dispatch(reqMessage));
    }

    @Test
    void unregister_stopsFurtherDispatches() throws IOException {
        S7UserDataPushDispatcher dispatcher = new S7UserDataPushDispatcher();
        AtomicInteger calls = new AtomicInteger();
        Closeable registration = dispatcher.register(0x02, 0x00, 0x05, msg -> calls.incrementAndGet());

        assertTrue(dispatcher.dispatch(pushMessage(0x02, 0x00, 0x05)));
        registration.close();
        assertFalse(dispatcher.dispatch(pushMessage(0x02, 0x00, 0x05)));
        assertEquals(1, calls.get());
        assertEquals(0, dispatcher.registeredHandlerCount());
    }

    private static S7Message pushMessage(int group, int type, int subfunction) {
        S7ParameterUserDataItem param = new S7ParameterUserDataItemCPUFunctions(
            (short) 0x12, (byte) type, (byte) group, (short) subfunction,
            (short) 0x00, null, null, null);
        return new S7MessageUserData(0,
            new S7ParameterUserData(Collections.singletonList(param)),
            new S7PayloadUserData(Collections.emptyList()));
    }
}

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
import org.apache.plc4x.java.s7.readwrite.S7MessageUserData;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserData;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserDataItem;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserDataItemCPUFunctions;

import java.io.Closeable;
import java.util.Map;
import java.util.Objects;
import java.util.concurrent.ConcurrentHashMap;
import java.util.function.Consumer;

/**
 * Routes unsolicited S7Comm UserData push messages (cyclic data, alarm pushes, mode-change
 * notifications) to the right service handler. Dispatching keys on the
 * {@code (cpuFunctionGroup, cpuFunctionType, cpuSubfunction)} triple from the UserData
 * parameter — that's how the PLC tags the kind of push it's sending.
 *
 * <p>Foundation only at this round. Concrete services (alarm, cyclic) will register
 * their handlers when they're ported.
 */
public final class S7UserDataPushDispatcher {

    private final Map<Key, Consumer<S7Message>> handlers = new ConcurrentHashMap<>();

    /** Register a handler for a push triple. The returned {@link Closeable} unregisters it. */
    public Closeable register(int cpuFunctionGroup, int cpuFunctionType, int cpuSubfunction,
                              Consumer<S7Message> handler) {
        Objects.requireNonNull(handler, "handler");
        Key key = new Key(cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
        handlers.put(key, handler);
        return () -> handlers.remove(key, handler);
    }

    /**
     * Try to dispatch an incoming UserData push to a registered handler. Returns {@code true}
     * iff the message was a UserData push that matched a registered handler. The connection's
     * receive loop should call this for any incoming message that didn't correlate against a
     * pending request.
     */
    public boolean dispatch(S7Message message) {
        if (!(message instanceof S7MessageUserData userData)) {
            return false;
        }
        if (!(userData.getParameter() instanceof S7ParameterUserData param)) {
            return false;
        }
        for (S7ParameterUserDataItem pItem : param.getItems()) {
            if (pItem instanceof S7ParameterUserDataItemCPUFunctions cpu) {
                Key key = new Key(cpu.getCpuFunctionGroup() & 0xFF,
                    cpu.getCpuFunctionType() & 0xFF, cpu.getCpuSubfunction() & 0xFFFF);
                Consumer<S7Message> handler = handlers.get(key);
                if (handler != null) {
                    handler.accept(message);
                    return true;
                }
            }
        }
        return false;
    }

    /** Visible for tests. */
    int registeredHandlerCount() {
        return handlers.size();
    }

    private record Key(int group, int type, int subfunction) {
    }

}

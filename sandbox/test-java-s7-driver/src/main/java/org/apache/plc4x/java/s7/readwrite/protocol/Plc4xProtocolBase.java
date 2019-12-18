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

package org.apache.plc4x.java.s7.readwrite.protocol;

import org.apache.plc4x.java.base.messages.PlcRequestContainer;
import org.apache.plc4x.java.s7.readwrite.TPKTPacket;

import java.util.function.Consumer;

public abstract class Plc4xProtocolBase<T> {

    protected final Plc4xProtocolContext<T> context;

    public Plc4xProtocolBase(Class<T> clazz) {
        // TODO create here?
        this.context = new Plc4xProtocolContext<T>(this, clazz);
    }

    public Plc4xProtocolContext<T> getContext() {
        return context;
    }

    public void onConnect() {
        // Intentionally do nothing here
    }

    protected abstract void encode(PlcRequestContainer msg, Consumer<T> sendHandler) throws Exception;

    protected abstract void decode(T msg) throws Exception;
}

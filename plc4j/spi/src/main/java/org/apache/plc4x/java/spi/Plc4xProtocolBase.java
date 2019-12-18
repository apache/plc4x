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

package org.apache.plc4x.java.spi;

import org.apache.plc4x.java.spi.messages.PlcRequestContainer;

public abstract class Plc4xProtocolBase<T> {

    protected final Plc4xNettyWrapper<T> context;

    public Plc4xProtocolBase(Class<T> clazz) {
        // TODO create here?
        this.context = new Plc4xNettyWrapper<T>(this, clazz);
    }

    public Plc4xNettyWrapper<T> getContext() {
        return context;
    }

    public void onConnect(ConversationContext<T> context) {
        // Intentionally do nothing here
    }

    protected abstract void encode(ConversationContext<T> context, PlcRequestContainer msg) throws Exception;

    protected abstract void decode(ConversationContext<T> context, T msg) throws Exception;
}

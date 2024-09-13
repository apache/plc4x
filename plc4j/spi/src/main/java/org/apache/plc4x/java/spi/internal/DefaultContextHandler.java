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
package org.apache.plc4x.java.spi.internal;

import java.util.concurrent.ExecutionException;
import java.util.concurrent.Future;
import org.apache.plc4x.java.spi.ConversationContext;

class DefaultContextHandler implements ConversationContext.ContextHandler {

    private final Future<Void> awaitable;
    private final Runnable cancel;

    public DefaultContextHandler(Future<Void> awaitable, Runnable cancel) {
        this.awaitable = awaitable;
        this.cancel = cancel;
    }

    @Override
    public boolean isDone() {
        return this.awaitable.isDone();
    }

    @Override
    public void cancel() {
        this.cancel.run();
    }

    @Override
    public void await() throws InterruptedException, ExecutionException {
        this.awaitable.get();
    }
}

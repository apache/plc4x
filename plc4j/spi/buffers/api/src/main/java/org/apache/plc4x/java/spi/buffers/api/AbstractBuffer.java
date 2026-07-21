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
package org.apache.plc4x.java.spi.buffers.api;

import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;

import java.util.ArrayDeque;
import java.util.Deque;
import java.util.HashMap;
import java.util.Map;

public abstract class AbstractBuffer implements Buffer {

    // Used as a stack (push/pop/peek). ArrayDeque instead of java.util.Stack: a buffer is a
    // single-threaded, per-message scratch object (positionInBits and the backing array are
    // themselves unsynchronized), so Stack's synchronization (it extends the synchronized Vector)
    // guards nothing here while adding a monitor enter/exit to getContext() on every field.
    protected final Deque<WithOption[]> context;

    public AbstractBuffer(WithOption... options) {
        context = new ArrayDeque<>();
        context.push(options);
    }

    @Override
    public void pushContext(WithOption... options) throws BufferException {
        Map<Class<?>, WithOption> newOptions = new HashMap<>();

        // Add all new options.
        for (WithOption option : options) {
            newOptions.put(option.getClass(), option);
        }

        // Add all inherited options.
        if (!context.isEmpty()) {
            for (WithOption option : context.peek()) {
                if (option.isSticky() && !newOptions.containsKey(option.getClass())) {
                    newOptions.put(option.getClass(), option);
                }
            }
        }

        // Switch the context to the next one.
        context.push(newOptions.values().toArray(new WithOption[0]));
    }

    @Override
    public void popContext(WithOption... options) throws BufferException {
        if (!context.isEmpty()) {
            context.pop();
        }
    }

    @Override
    public WithOption[] getContext() {
        if (!context.isEmpty()) {
            return context.peek();
        }
        return new WithOption[0];
    }

}

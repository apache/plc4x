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
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayDeque;
import java.util.Deque;
import java.util.HashMap;
import java.util.Map;

public abstract class AbstractBuffer implements Buffer {

    private static final Logger LOGGER = LoggerFactory.getLogger(AbstractBuffer.class);

    /**
     * How deeply a message may nest its types. Every type pushes a context on the way in and pops
     * it on the way out, so the depth of this stack is the depth of the message being read or
     * written, and a type that contains itself turns the sender's byte count into our call depth.
     * <p>
     * The deepest message in the project's own test suites nests 36 contexts, and the parser runs
     * out of stack somewhere above four thousand, so this sits far above anything a real device
     * sends and far below the point where the recursion stops being reported as a parse failure.
     * <p>
     * This is a backstop, not a substitute for a bound a protocol can state more precisely: where
     * a type knows its own sensible depth it says so in the mspec, as the ADS data type table
     * does. Such a bound has to fit inside this one - the ADS budget of 255 entries nests around
     * 512 contexts - so keep this comfortably above the largest of them.
     */
    public static final int DEFAULT_MAX_NESTING_DEPTH = 1024;

    /**
     * Name of the environment variable that overrides {@link #DEFAULT_MAX_NESTING_DEPTH}, for the
     * device whose messages genuinely nest deeper than anything we have seen.
     */
    public static final String MAX_NESTING_DEPTH_ENV = "PLC4X_MAX_NESTING_DEPTH";

    private static final int MAX_NESTING_DEPTH = resolveMaxNestingDepth(System.getenv(MAX_NESTING_DEPTH_ENV));

    // Used as a stack (push/pop/peek). ArrayDeque instead of java.util.Stack: a buffer is a
    // single-threaded, per-message scratch object (positionInBits and the backing array are
    // themselves unsynchronized), so Stack's synchronization (it extends the synchronized Vector)
    // guards nothing here while adding a monitor enter/exit to getContext() on every field.
    protected final Deque<WithOption[]> context;

    public AbstractBuffer(WithOption... options) {
        context = new ArrayDeque<>();
        context.push(options);
    }

    /**
     * @return the nesting depth in force, which is {@link #DEFAULT_MAX_NESTING_DEPTH} unless
     * {@value #MAX_NESTING_DEPTH_ENV} says otherwise.
     */
    public static int getMaxNestingDepth() {
        return MAX_NESTING_DEPTH;
    }

    /**
     * @param configured the value of {@value #MAX_NESTING_DEPTH_ENV}, or null when it is unset
     * @return the depth that value asks for, or {@link #DEFAULT_MAX_NESTING_DEPTH} when it asks
     * for nothing usable
     */
    public static int resolveMaxNestingDepth(String configured) {
        if (configured == null || configured.trim().isEmpty()) {
            return DEFAULT_MAX_NESTING_DEPTH;
        }
        int depth;
        try {
            depth = Integer.parseInt(configured.trim());
        } catch (NumberFormatException e) {
            LOGGER.warn("{} is not a number ({}), keeping the maximum nesting depth of {}",
                MAX_NESTING_DEPTH_ENV, configured, DEFAULT_MAX_NESTING_DEPTH);
            return DEFAULT_MAX_NESTING_DEPTH;
        }
        if (depth < 1) {
            LOGGER.warn("{} must be positive but was {}, keeping the maximum nesting depth of {}",
                MAX_NESTING_DEPTH_ENV, depth, DEFAULT_MAX_NESTING_DEPTH);
            return DEFAULT_MAX_NESTING_DEPTH;
        }
        return depth;
    }

    @Override
    public void pushContext(WithOption... options) throws BufferException {
        // Refuse before growing the stack rather than after, so the depth we report is a depth we
        // never actually occupied.
        if (context.size() >= MAX_NESTING_DEPTH) {
            throw new BufferException("Maximum nesting depth of " + MAX_NESTING_DEPTH + " exceeded");
        }

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
    public int getContextDepth() {
        return context.size();
    }

    @Override
    public void resetContextDepth(int contextDepth) {
        while (context.size() > contextDepth) {
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

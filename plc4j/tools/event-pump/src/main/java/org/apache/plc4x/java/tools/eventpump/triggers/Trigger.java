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

package org.apache.plc4x.java.tools.eventpump.triggers;

/**
 * Interface for triggers that determine when a tag batch should be fetched.
 * <p>
 * Triggers can be time-based, subscription-based, or expression-based.
 * When a trigger fires, it notifies its listener to fetch the associated batch.
 */
public interface Trigger extends AutoCloseable {

    /**
     * Start the trigger.
     * After starting, the trigger will fire events according to its configuration.
     *
     * @param listener The listener to notify when the trigger fires
     * @throws IllegalStateException if the trigger is already started
     */
    void start(TriggerListener listener);

    /**
     * Stop the trigger.
     * After stopping, the trigger will no longer fire events.
     * The trigger can be restarted by calling start() again.
     */
    void stop();

    /**
     * Check if the trigger is currently running.
     *
     * @return true if the trigger is running, false otherwise
     */
    boolean isRunning();

    /**
     * Get the trigger type (for debugging/logging).
     *
     * @return A string describing the trigger type
     */
    String getType();

    /**
     * Close the trigger and release all resources.
     * After closing, the trigger cannot be restarted.
     */
    @Override
    void close();

    /**
     * Listener interface for trigger events.
     */
    @FunctionalInterface
    interface TriggerListener {
        /**
         * Called when the trigger fires.
         *
         * @param trigger The trigger that fired
         */
        void onTrigger(Trigger trigger);
    }
}

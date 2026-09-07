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

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * A trigger that fires when a subscribed PLC tag value changes.
 * <p>
 * <b>NOTE:</b> This trigger implementation is currently a placeholder.
 * Full subscription support requires protocol-specific implementations.
 * <p>
 * For now, use {@link TimerTrigger} for polling-based data collection.
 * Subscription support will be added in a future release.
 * <p>
 * Example usage (when implemented):
 * <pre>
 * // Trigger on any value change
 * Trigger trigger = new SubscriptionTrigger(connection, "triggerTag", "MAIN.trigger");
 * </pre>
 */
public class SubscriptionTrigger implements Trigger {

    private static final Logger LOGGER = LoggerFactory.getLogger(SubscriptionTrigger.class);

    private final String tagName;
    private final String tagAddress;

    private volatile boolean running = false;
    private volatile boolean closed = false;

    /**
     * Create a subscription trigger that fires on any value change.
     *
     * @param tagName The name for the subscribed tag
     * @param tagAddress The tag address to subscribe to
     */
    public SubscriptionTrigger(String tagName, String tagAddress) {
        if (tagName == null || tagName.trim().isEmpty()) {
            throw new IllegalArgumentException("Tag name cannot be null or empty");
        }
        if (tagAddress == null || tagAddress.trim().isEmpty()) {
            throw new IllegalArgumentException("Tag address cannot be null or empty");
        }

        this.tagName = tagName;
        this.tagAddress = tagAddress;

        LOGGER.debug("Created SubscriptionTrigger placeholder for tag '{}' at address '{}'", tagName, tagAddress);
    }

    @Override
    public void start(TriggerListener listener) {
        throw new UnsupportedOperationException(
            "SubscriptionTrigger is not yet fully implemented. " +
            "Use TimerTrigger for polling-based data collection. " +
            "Subscription support will be added in a future release."
        );
    }

    @Override
    public void stop() {
        // No-op for placeholder implementation
    }

    @Override
    public boolean isRunning() {
        return false;  // Placeholder never runs
    }

    @Override
    public String getType() {
        return "Subscription [PLACEHOLDER] (tag=" + tagName + ", address=" + tagAddress + ")";
    }

    @Override
    public void close() {
        // No-op for placeholder implementation
    }

    /**
     * Get the tag name being subscribed to.
     *
     * @return The tag name
     */
    public String getTagName() {
        return tagName;
    }

    /**
     * Get the tag address being subscribed to.
     *
     * @return The tag address
     */
    public String getTagAddress() {
        return tagAddress;
    }
}

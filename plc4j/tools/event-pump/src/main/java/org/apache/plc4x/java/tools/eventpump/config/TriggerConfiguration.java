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

package org.apache.plc4x.java.tools.eventpump.config;

import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * Configuration for a trigger.
 * <p>
 * The type field determines the trigger type:
 * - "timer": Time-based trigger (requires either intervalSeconds or intervalMillis, not both)
 * - "subscription": Subscription-based trigger (requires tagName and tagAddress)
 */
public class TriggerConfiguration {

    @JsonProperty("type")
    private String type;

    // Timer trigger properties
    @JsonProperty("intervalSeconds")
    private Long intervalSeconds;

    // Sub-second polling is common, so the interval can alternatively be given in
    // milliseconds. Setting both units on the same trigger is rejected when the trigger is
    // built, rather than silently picking one of them.
    @JsonProperty("intervalMillis")
    private Long intervalMillis;

    @JsonProperty("initialDelaySeconds")
    private Long initialDelaySeconds;

    @JsonProperty("initialDelayMillis")
    private Long initialDelayMillis;

    // Subscription trigger properties
    @JsonProperty("tagName")
    private String tagName;

    @JsonProperty("tagAddress")
    private String tagAddress;

    @JsonProperty("condition")
    private String condition; // Optional: Expression for filtering events

    /**
     * Get the trigger type.
     *
     * @return The trigger type ("timer" or "subscription")
     */
    public String getType() {
        return type;
    }

    /**
     * Set the trigger type.
     *
     * @param type The trigger type
     */
    public void setType(String type) {
        this.type = type;
    }

    /**
     * Get the interval in seconds (for timer triggers).
     *
     * @return The interval in seconds
     */
    public Long getIntervalSeconds() {
        return intervalSeconds;
    }

    /**
     * Set the interval in seconds (for timer triggers). Mutually exclusive with
     * {@link #setIntervalMillis(Long)} - setting both is a configuration error.
     *
     * @param intervalSeconds The interval in seconds
     */
    public void setIntervalSeconds(Long intervalSeconds) {
        this.intervalSeconds = intervalSeconds;
    }

    /**
     * Get the interval in milliseconds (for timer triggers).
     *
     * @return The interval in milliseconds, or null if it is configured in seconds
     */
    public Long getIntervalMillis() {
        return intervalMillis;
    }

    /**
     * Set the interval in milliseconds (for timer triggers). Mutually exclusive with
     * {@link #setIntervalSeconds(Long)} - setting both is a configuration error.
     *
     * @param intervalMillis The interval in milliseconds
     */
    public void setIntervalMillis(Long intervalMillis) {
        this.intervalMillis = intervalMillis;
    }

    /**
     * Get the initial delay in seconds (for timer triggers).
     *
     * @return The initial delay in seconds, or null if not set
     */
    public Long getInitialDelaySeconds() {
        return initialDelaySeconds;
    }

    /**
     * Set the initial delay in seconds (for timer triggers). Mutually exclusive with
     * {@link #setInitialDelayMillis(Long)} - setting both is a configuration error.
     *
     * @param initialDelaySeconds The initial delay in seconds
     */
    public void setInitialDelaySeconds(Long initialDelaySeconds) {
        this.initialDelaySeconds = initialDelaySeconds;
    }

    /**
     * Get the initial delay in milliseconds (for timer triggers).
     *
     * @return The initial delay in milliseconds, or null if it is configured in seconds
     */
    public Long getInitialDelayMillis() {
        return initialDelayMillis;
    }

    /**
     * Set the initial delay in milliseconds (for timer triggers). Mutually exclusive with
     * {@link #setInitialDelaySeconds(Long)} - setting both is a configuration error.
     *
     * @param initialDelayMillis The initial delay in milliseconds
     */
    public void setInitialDelayMillis(Long initialDelayMillis) {
        this.initialDelayMillis = initialDelayMillis;
    }

    /**
     * Get the tag name (for subscription triggers).
     *
     * @return The tag name
     */
    public String getTagName() {
        return tagName;
    }

    /**
     * Set the tag name (for subscription triggers).
     *
     * @param tagName The tag name
     */
    public void setTagName(String tagName) {
        this.tagName = tagName;
    }

    /**
     * Get the tag address (for subscription triggers).
     *
     * @return The tag address
     */
    public String getTagAddress() {
        return tagAddress;
    }

    /**
     * Set the tag address (for subscription triggers).
     *
     * @param tagAddress The tag address
     */
    public void setTagAddress(String tagAddress) {
        this.tagAddress = tagAddress;
    }

    /**
     * Get the condition expression (for subscription triggers).
     * <p>
     * This is an optional expression that can be used to filter subscription events.
     * The syntax depends on the expression evaluator used.
     *
     * @return The condition expression, or null if not set
     */
    public String getCondition() {
        return condition;
    }

    /**
     * Set the condition expression (for subscription triggers).
     *
     * @param condition The condition expression
     */
    public void setCondition(String condition) {
        this.condition = condition;
    }
}

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
package org.apache.plc4x.java.utils.testutils.driver.internal.utils;

/**
 * Utility class for introducing delays in test execution.
 */
public class Delay {

    private static final long SHORT_DELAY_MS = 10;
    private static final long MEDIUM_DELAY_MS = 50;
    private static final long LONG_DELAY_MS = 100;

    private Delay() {
        // Utility class
    }

    /**
     * Introduces a short delay (10ms).
     */
    public static void shortDelay() {
        delay(SHORT_DELAY_MS);
    }

    /**
     * Introduces a medium delay (50ms).
     */
    public static void mediumDelay() {
        delay(MEDIUM_DELAY_MS);
    }

    /**
     * Introduces a long delay (100ms).
     */
    public static void longDelay() {
        delay(LONG_DELAY_MS);
    }

    /**
     * Introduces a custom delay.
     *
     * @param milliseconds the delay in milliseconds
     */
    public static void delay(long milliseconds) {
        try {
            Thread.sleep(milliseconds);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
    }
}

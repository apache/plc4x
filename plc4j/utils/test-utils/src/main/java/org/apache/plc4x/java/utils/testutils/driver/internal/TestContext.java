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
package org.apache.plc4x.java.utils.testutils.driver.internal;

import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;

import java.net.URI;
import java.util.concurrent.CompletableFuture;

/**
 * Context for sharing state between test steps.
 * Used to pass futures from api-request steps to api-response steps.
 */
public class TestContext {

    private CompletableFuture<?> pendingResponse;
    private CompletableFuture<PlcSubscriptionEvent> subscriptionEventCompletableFuture;
    private boolean autoMigrate;
    private URI testsuiteUri;

    /**
     * Stores a pending response future from an API request.
     *
     * @param future the CompletableFuture to store
     */
    public void setPendingResponse(CompletableFuture<?> future) {
        this.pendingResponse = future;
    }

    /**
     * Retrieves and clears the pending response future.
     *
     * @return the stored CompletableFuture, or null if none was set
     */
    public CompletableFuture<?> getPendingResponse() {
        CompletableFuture<?> result = this.pendingResponse;
        this.pendingResponse = null;
        return result;
    }

    /**
     * Checks if there's a pending response.
     *
     * @return true if a response is pending
     */
    public boolean hasPendingResponse() {
        return this.pendingResponse != null;
    }

    public void addSubscriptionEvent(PlcSubscriptionEvent subscriptionEvent) {
        // If the event came in before the testsuite asked for it, save it for later.
        if(subscriptionEventCompletableFuture == null) {
            subscriptionEventCompletableFuture = CompletableFuture.completedFuture(subscriptionEvent);
            return;
        }
        // Otherwise just return the received event.
        subscriptionEventCompletableFuture.complete(subscriptionEvent);
        subscriptionEventCompletableFuture = null;
    }

    public CompletableFuture<PlcSubscriptionEvent> getSubscriptionEvent() {
        if (subscriptionEventCompletableFuture == null) {
            subscriptionEventCompletableFuture = new CompletableFuture<>();
        }
        return subscriptionEventCompletableFuture;
    }

    /**
     * Sets whether auto-migration is enabled.
     *
     * @param autoMigrate true to enable auto-migration
     */
    public void setAutoMigrate(boolean autoMigrate) {
        this.autoMigrate = autoMigrate;
    }

    /**
     * Gets whether auto-migration is enabled.
     *
     * @return true if auto-migration is enabled
     */
    public boolean isAutoMigrate() {
        return autoMigrate;
    }

    /**
     * Sets the testsuite document URI.
     *
     * @param testsuiteUri the URI of the testsuite XML file
     */
    public void setTestsuiteUri(URI testsuiteUri) {
        this.testsuiteUri = testsuiteUri;
    }

    /**
     * Gets the testsuite document URI.
     *
     * @return the URI of the testsuite XML file
     */
    public URI getTestsuiteUri() {
        return testsuiteUri;
    }
}

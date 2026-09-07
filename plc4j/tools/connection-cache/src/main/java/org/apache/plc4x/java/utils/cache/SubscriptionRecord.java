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
package org.apache.plc4x.java.utils.cache;

import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;

import java.util.Collection;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.function.Consumer;

/**
 * Records the details of an active subscription so it can be re-established when the cached
 * connection is invalidated and recreated.
 * <p>
 * It keeps two per-tag handle maps so a client that still holds a handle from <em>before</em> a
 * reconnection can be transparently mapped to the new handle (e.g. for unsubscription):
 * <ul>
 *   <li><b>original handles</b> — the handle the client first received for each tag (never changes);</li>
 *   <li><b>current handles</b> — the handle currently valid for each tag (updated on every reconnect).</li>
 * </ul>
 * Handles are compared by object identity (PLC4X subscription handles are identity objects and do not
 * override {@code equals}).
 */
class SubscriptionRecord {

    private final PlcSubscriptionRequest originalRequest;
    private final Consumer<PlcSubscriptionEvent> consumer;

    // tag name -> the handle the client first received (immutable for the life of the record)
    private final Map<String, PlcSubscriptionHandle> originalHandles = new ConcurrentHashMap<>();
    // tag name -> the currently valid handle (re-pointed after each reconnection)
    private final Map<String, PlcSubscriptionHandle> currentHandles = new ConcurrentHashMap<>();

    // The most recent subscription response (updated after each re-subscription)
    private volatile PlcSubscriptionResponse currentResponse;

    SubscriptionRecord(PlcSubscriptionRequest originalRequest,
                       Consumer<PlcSubscriptionEvent> consumer,
                       PlcSubscriptionResponse initialResponse) {
        this.originalRequest = originalRequest;
        this.consumer = consumer;
        this.currentResponse = initialResponse;

        if (initialResponse != null) {
            for (String tagName : initialResponse.getTagNames()) {
                PlcSubscriptionHandle handle = initialResponse.getSubscriptionHandle(tagName);
                if (handle != null) {
                    originalHandles.put(tagName, handle);
                    currentHandles.put(tagName, handle);
                }
            }
        }
    }

    /** The subscription request used to create this subscription — replayed (tag-by-tag) to restore it. */
    PlcSubscriptionRequest getOriginalRequest() {
        return originalRequest;
    }

    /** The request-level event consumer, if one was set via {@code setConsumer(...)} (may be {@code null}). */
    Consumer<PlcSubscriptionEvent> getConsumer() {
        return consumer;
    }

    /** The current subscription response (updated after reconnection). */
    PlcSubscriptionResponse getCurrentResponse() {
        return currentResponse;
    }

    /**
     * Re-points this record's per-tag <em>current</em> handles to those of the new subscription response
     * after a reconnection. The original handles (which the client may still hold) are left untouched so
     * {@link #getCurrentHandle(PlcSubscriptionHandle)} can keep mapping old → new.
     *
     * @param newResponse the response from re-subscribing on the recreated connection.
     */
    void updateAfterReconnection(PlcSubscriptionResponse newResponse) {
        if (newResponse == null) {
            return;
        }
        this.currentResponse = newResponse;
        for (String tagName : newResponse.getTagNames()) {
            PlcSubscriptionHandle newHandle = newResponse.getSubscriptionHandle(tagName);
            if (newHandle != null) {
                currentHandles.put(tagName, newHandle);
            }
        }
    }

    /**
     * Maps a handle the client is holding (possibly an original handle from before a reconnection) to the
     * currently valid handle. Returns the input unchanged if it is not tracked by this record.
     */
    PlcSubscriptionHandle getCurrentHandle(PlcSubscriptionHandle clientHandle) {
        if (clientHandle == null) {
            return null;
        }
        for (Map.Entry<String, PlcSubscriptionHandle> entry : originalHandles.entrySet()) {
            if (entry.getValue() == clientHandle) {
                PlcSubscriptionHandle current = currentHandles.get(entry.getKey());
                return current != null ? current : clientHandle;
            }
        }
        // Already a current handle (or unknown) — return as-is.
        return clientHandle;
    }

    /** All currently valid subscription handles. */
    Collection<PlcSubscriptionHandle> getCurrentHandles() {
        return currentHandles.values();
    }

    /** Whether the given handle is tracked by this record, as either an original or a current handle. */
    boolean containsHandle(PlcSubscriptionHandle handle) {
        if (handle == null) {
            return false;
        }
        for (PlcSubscriptionHandle h : originalHandles.values()) {
            if (h == handle) {
                return true;
            }
        }
        for (PlcSubscriptionHandle h : currentHandles.values()) {
            if (h == handle) {
                return true;
            }
        }
        return false;
    }
}

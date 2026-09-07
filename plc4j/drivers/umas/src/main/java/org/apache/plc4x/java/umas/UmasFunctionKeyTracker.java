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
package org.apache.plc4x.java.umas;

import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

/**
 * Tracks UMAS request function keys for response discrimination during parsing.
 * <p>
 * UMAS responses use a generic function key (0xFE) and require the original
 * request's function key to determine the correct response subtype. The connection
 * records the function key before sending each request, and the parser
 * retrieves it by peeking at the transaction ID from the MBAP header bytes.
 * <p>
 * One of these belongs to one connection. It used to be static, so every connection
 * in the process shared one map keyed by transaction id alone - and transaction ids
 * are a counter per connection, not per process. Two connections talking to two
 * devices reach the same id in the ordinary course of things, and the second to
 * arrive overwrote the first: a response then got parsed under another
 * conversation's function key, which decides the response subtype. Nothing about
 * that is visible in a log; it comes out as a parse failure or as the wrong shape of
 * data. Keeping it per connection is what makes the id meaningful.
 */
public class UmasFunctionKeyTracker {

    private final Map<Integer, Short> pendingFunctionKeys = new ConcurrentHashMap<>();

    /**
     * Records that a request with the given transaction ID used the specified
     * UMAS function key. Called by the connection before sending a request.
     */
    public void trackRequest(int transactionId, short functionKey) {
        pendingFunctionKeys.put(transactionId, functionKey);
    }

    /**
     * Retrieves and removes the function key for the given transaction ID.
     * Called by the parser lambda after peeking the transaction ID from the
     * raw MBAP header bytes.
     *
     * @return the function key, or 0 if not tracked (e.g. unsolicited message)
     */
    public short consumeFunctionKey(int transactionId) {
        Short fk = pendingFunctionKeys.remove(transactionId);
        return fk != null ? fk : (short) 0;
    }

    /**
     * Forgets a request that will not be answered - it was never sent, or it gave up waiting.
     * Without this an unanswered request leaves its key behind for as long as the connection
     * lives, and a later request reaching the same id finds somebody else's.
     */
    public void forget(int transactionId) {
        pendingFunctionKeys.remove(transactionId);
    }

    /** Forgets everything, for a connection that is closing. */
    public void clear() {
        pendingFunctionKeys.clear();
    }

    int trackedCount() {
        return pendingFunctionKeys.size();
    }

}

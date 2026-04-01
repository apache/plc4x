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
package org.apache.plc4x.java.umas.readwrite;

import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

/**
 * Tracks UMAS request function keys for response discrimination during parsing.
 * <p>
 * UMAS responses use a generic function key (0xFE) and require the original
 * request's function key to determine the correct response subtype. The protocol
 * logic records the function key before sending each request, and the parser
 * retrieves it by peeking at the transaction ID from the MBAP header bytes.
 */
public class UmasFunctionKeyTracker {

    private static final Map<Integer, Short> PENDING_FUNCTION_KEYS = new ConcurrentHashMap<>();

    /**
     * Records that a request with the given transaction ID used the specified
     * UMAS function key. Called by the protocol logic before sending a request.
     */
    public static void trackRequest(int transactionId, short functionKey) {
        PENDING_FUNCTION_KEYS.put(transactionId, functionKey);
    }

    /**
     * Retrieves and removes the function key for the given transaction ID.
     * Called by the parser lambda after peeking the transaction ID from the
     * raw MBAP header bytes.
     *
     * @return the function key, or 0 if not tracked (e.g. unsolicited message)
     */
    public static short consumeFunctionKey(int transactionId) {
        Short fk = PENDING_FUNCTION_KEYS.remove(transactionId);
        return fk != null ? fk : (short) 0;
    }

}

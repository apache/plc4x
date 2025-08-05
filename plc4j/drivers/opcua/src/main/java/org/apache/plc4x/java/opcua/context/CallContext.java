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

package org.apache.plc4x.java.opcua.context;

import java.util.function.Supplier;
import org.apache.plc4x.java.opcua.protocol.chunk.ChunkStorage;
import org.apache.plc4x.java.opcua.readwrite.SecurityHeader;

public class CallContext {

    private final SecurityHeader sequenceHeader;
    private final Supplier<Integer> sequenceSupplier;
    private final int requestId;

    public CallContext(SecurityHeader sequenceHeader, Supplier<Integer> sequenceSupplier, int requestId) {
        this.sequenceHeader = sequenceHeader;
        this.sequenceSupplier = sequenceSupplier;
        this.requestId = requestId;
    }

    public SecurityHeader getSecurityHeader() {
        return sequenceHeader;
    }

    public int getNextSequenceNumber() {
        return sequenceSupplier.get();
    }

    public int getRequestId() {
        return requestId;
    }
}

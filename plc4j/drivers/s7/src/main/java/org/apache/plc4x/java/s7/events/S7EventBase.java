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

package org.apache.plc4x.java.s7.events;

import java.time.Instant;
import org.apache.plc4x.java.api.messages.PlcMetadataKeys;
import org.apache.plc4x.java.api.metadata.Metadata;
import org.apache.plc4x.java.spi.metadata.DefaultMetadata;
import org.apache.plc4x.java.api.metadata.time.TimeSource;

public abstract class S7EventBase implements S7Event {

    private final Instant timestamp;
    private final Metadata metadata;

    S7EventBase() {
        this(Instant.now());
    }

    S7EventBase(Instant timestamp) {
        this(timestamp, new DefaultMetadata.Builder()
            .put(PlcMetadataKeys.TIMESTAMP, timestamp.getEpochSecond())
            .put(PlcMetadataKeys.TIMESTAMP_SOURCE, TimeSource.HARDWARE) // event triggered by PLC itself
            .build()
        );
    }

    S7EventBase(Instant timestamp, Metadata metadata) {
        this.timestamp = timestamp;
        this.metadata = metadata;
    }

    @Override
    public Metadata getTagMetadata(String name) {
        return metadata;
    }

    @Override
    public Instant getTimestamp() {
        return timestamp;
    }
}

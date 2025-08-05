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
package org.apache.plc4x.java.spi.messages;

import java.util.Collections;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.metadata.Metadata;
import org.apache.plc4x.java.spi.metadata.DefaultMetadata;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.messages.utils.PlcResponseItem;

import java.time.Instant;
import java.util.Collection;
import java.util.Map;

public class DefaultPlcSubscriptionEvent extends DefaultPlcReadResponse implements PlcSubscriptionEvent {

    public final Instant timestamp;

    public DefaultPlcSubscriptionEvent(Instant timestamp,
        Map<String, PlcResponseItem<PlcValue>> tags) {
        this(timestamp, tags, Collections.emptyMap());
    }

    public DefaultPlcSubscriptionEvent(Instant timestamp,
                                       Map<String, PlcResponseItem<PlcValue>> tags,
                                       Map<String, Metadata> metadata) {
        super(null, tags, metadata);
        this.timestamp = timestamp;
    }

    @Override
    public Collection<String> getTagNames() {
        return getValues().keySet();
    }

    @Override
    public PlcTag getTag(String name) {
        throw new UnsupportedOperationException("getTag('" + name + "') not supported on " + this.getClass());
    }

    @Override
    public Instant getTimestamp() {
        return timestamp;
    }

}

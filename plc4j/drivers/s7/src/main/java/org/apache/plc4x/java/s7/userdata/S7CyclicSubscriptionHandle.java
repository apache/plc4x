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
package org.apache.plc4x.java.s7.userdata;

import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.s7.tag.S7Tag;
import org.apache.plc4x.java.spi.drivers.functions.PlcSubscriber;

import java.util.Collections;
import java.util.Objects;
import java.util.function.Consumer;

/**
 * Handle for one tag inside an active cyclic subscription. Multiple handles can share the
 * same {@code jobId} when a single subscribe request bundled several tags — the connection
 * matches incoming pushes against {@code jobId} and dispatches values to all sibling
 * handles in the same job.
 */
public final class S7CyclicSubscriptionHandle implements PlcSubscriptionHandle {

    private final PlcSubscriber subscriber;
    private final String tagName;
    private final S7Tag tag;
    private final short jobId;
    private final int itemIndex;

    public S7CyclicSubscriptionHandle(PlcSubscriber subscriber, String tagName, S7Tag tag,
                                      short jobId, int itemIndex) {
        this.subscriber = Objects.requireNonNull(subscriber, "subscriber");
        this.tagName = Objects.requireNonNull(tagName, "tagName");
        this.tag = Objects.requireNonNull(tag, "tag");
        this.jobId = jobId;
        this.itemIndex = itemIndex;
    }

    public String getTagName()  { return tagName; }
    public S7Tag   getTag()      { return tag; }
    public short   getJobId()    { return jobId; }
    public int     getItemIndex() { return itemIndex; }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer) {
        return subscriber.registerConsumer(consumer, Collections.singletonList(this));
    }

    @Override
    public String toString() {
        return "S7CyclicSubscriptionHandle[" + tagName + "@job=" + jobId + ",idx=" + itemIndex + "]";
    }
}

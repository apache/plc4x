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

package org.apache.plc4x.java.profinet.device;

import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.spi.messages.PlcSubscriber;
import org.apache.plc4x.java.spi.model.DefaultPlcConsumerRegistration;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionHandle;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Set;
import java.util.function.Consumer;

public class ProfinetSubscriptionHandle extends DefaultPlcSubscriptionHandle {

    private static final Logger LOGGER = LoggerFactory.getLogger(ProfinetSubscriptionHandle.class);
    private Set<Consumer<PlcSubscriptionEvent>> consumers;
    private PlcSubscriber plcSubscriber;

    public ProfinetSubscriptionHandle(PlcSubscriber plcSubscriber) {
        super(plcSubscriber);
        this.plcSubscriber = plcSubscriber;
    }

    /**
     * Registers a new Consumer, this allows multiple PLC4X consumers to use the same subscription.
     *
     * @param consumer - Consumer to be used to send any returned values.
     * @return PlcConsumerRegistration - return the important information back to the client.
     */
    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer) {
        LOGGER.info("Registering a new Profinet subscription consumer");
        consumers.add(consumer);
        return new DefaultPlcConsumerRegistration(plcSubscriber, consumer, this);
    }

    public Set<Consumer<PlcSubscriptionEvent>> getConsumers() {
        return consumers;
    }
}

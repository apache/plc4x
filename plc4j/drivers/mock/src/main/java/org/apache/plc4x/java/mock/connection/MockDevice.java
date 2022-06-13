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
package org.apache.plc4x.java.mock.connection;

import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;

import java.util.Collection;
import java.util.function.Consumer;

/**
 * Mock Object to do assertions on.
 */
public interface MockDevice {

    ResponseItem<PlcValue> read(String fieldQuery);

    PlcResponseCode write(String fieldQuery, Object value);

    ResponseItem<PlcSubscriptionHandle> subscribe(String fieldQuery);

    void unsubscribe();

    // TODO: this might not be right here as you are not really register at the device, rather on the connection
    PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> handles);

    // TODO: this might not be right here as you are not really register at the device, rather on the connection
    void unregister(PlcConsumerRegistration registration);

}

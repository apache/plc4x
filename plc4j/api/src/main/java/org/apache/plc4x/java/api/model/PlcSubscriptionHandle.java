/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.java.api.model;

import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;

import java.util.function.Consumer;

/**
 * When subscribing to remote resources, depending on the used protocol
 * different data is used to identify a subscription. This interface is
 * to be implemented in the individual Driver implementations to contain
 * all information needed to pull or unsubscribe any form of subscription.
 * <p>
 * For every subscribed item, a separate {@link PlcSubscriptionHandle} object is
 * returned in order to allow fine granular unsubscriptions.
 */
public interface PlcSubscriptionHandle {

    PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer);
}

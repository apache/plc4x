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
package org.apache.plc4x.java.base.messages;

import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.model.PlcField;

import java.time.Duration;
import java.util.LinkedHashSet;
import java.util.LinkedList;

public class DefaultPlcSubscriptionRequest implements InternalPlcSubscriptionRequest {

    @Override
    public int getNumberOfFields() {
        return 0;
    }

    @Override
    public LinkedHashSet<String> getFieldNames() {
        return null;
    }

    @Override
    public PlcField getField(String name) {
        return null;
    }

    @Override
    public LinkedList<PlcField> getFields() {
        return null;
    }

    @Override
    public Builder builder() {
        return new Builder() {
            @Override
            public Builder addCyclicField(String name, String fieldQuery, Duration pollingInterval) {
                return null;
            }

            @Override
            public Builder addChangeOfStateField(String name, String fieldQuery) {
                return null;
            }

            @Override
            public Builder addEventField(String name, String fieldQuery) {
                return null;
            }

            @Override
            public PlcSubscriptionRequest build() {
                return null;
            }
        };
    }

}

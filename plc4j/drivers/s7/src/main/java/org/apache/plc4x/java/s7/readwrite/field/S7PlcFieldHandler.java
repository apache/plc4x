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
package org.apache.plc4x.java.s7.readwrite.field;

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.spi.connection.PlcFieldHandler;

public class S7PlcFieldHandler implements PlcFieldHandler {

    @Override
    public PlcField createField(String fieldQuery) {
        if (S7Field.matches(fieldQuery)) {
            return S7Field.of(fieldQuery);
        } else if (S7SubscriptionField.matches(fieldQuery)){
            return S7SubscriptionField.of(fieldQuery);
        } else if (S7SslField.matches(fieldQuery)){
            return S7SslField.of(fieldQuery);
        }
        throw new PlcInvalidFieldException(fieldQuery);
    }

}

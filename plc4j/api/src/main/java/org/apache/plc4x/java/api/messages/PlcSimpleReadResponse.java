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
package org.apache.plc4x.java.api.messages;

import org.apache.plc4x.java.api.types.Value;

/**
 * A response for a simple read request.
 *
 * @param <T> the type of variable being requested.
 */
public class PlcSimpleReadResponse<T extends Value> extends AbstractPlcSimpleResourceMessage<T> implements PlcResponse {

    private final T value;

    public PlcSimpleReadResponse(Class<T> datatype, Address address, T value) {
        this(datatype, address, 1, value);
    }

    public PlcSimpleReadResponse(Class<T> datatype, Address address, int size, T value) {
        super(datatype, address, size);
        this.value = value;
    }

    public T getValue() {
        return value;
    }

}

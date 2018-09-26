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

import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;

import java.util.Collection;
import java.util.Map;

public class DefaultPlcWriteResponse implements InternalPlcWriteResponse {

    private final InternalPlcWriteRequest request;
    private final Map<String, PlcResponseCode> values;

    public DefaultPlcWriteResponse(InternalPlcWriteRequest request, Map<String, PlcResponseCode> values) {
        this.request = request;
        this.values = values;
    }

    @Override
    public Map<String, PlcResponseCode> getValues() {
        return values;
    }

    @Override
    public InternalPlcWriteRequest getRequest() {
        return request;
    }

    @Override
    public Collection<String> getFieldNames() {
        return request.getFieldNames();
    }

    @Override
    public PlcField getField(String name) {
        return request.getField(name);
    }

    @Override
    public PlcResponseCode getResponseCode(String name) {
        return values.get(name);
    }

}

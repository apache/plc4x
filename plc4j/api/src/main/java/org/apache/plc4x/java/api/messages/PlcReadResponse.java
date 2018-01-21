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

import org.apache.plc4x.java.api.messages.items.ReadRequestItem;
import org.apache.plc4x.java.api.messages.items.ReadResponseItem;

import java.util.Collections;
import java.util.List;
import java.util.Optional;

public class PlcReadResponse extends PlcResponse<PlcReadRequest, ReadResponseItem<?>, ReadRequestItem<?>> {

    public PlcReadResponse(PlcReadRequest request, ReadResponseItem<?> responseItems) {
        super(request, Collections.singletonList(responseItems));
    }

    public PlcReadResponse(PlcReadRequest request, List<? extends ReadResponseItem<?>> responseItems) {
        super(request, responseItems);
    }

    @SuppressWarnings("unchecked")
    public <T> Optional<ReadResponseItem<T>> getValue(ReadRequestItem<T> item) {
        return (Optional) super.getValue(item);
    }
}

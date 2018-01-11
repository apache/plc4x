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
package org.apache.plc4x.java.api.messages.specific;

import org.apache.plc4x.java.api.messages.items.ReadRequestItem;
import org.apache.plc4x.java.api.model.Address;

/**
 * @deprecated methods integrated into super type
 */
@Deprecated
public class SinglePlcReadRequest<T> extends TypeSafePlcReadRequest<T> {

    public SinglePlcReadRequest(Class<T> dataType) {
        super(dataType);
    }

    public SinglePlcReadRequest(Class<T> dataType, ReadRequestItem<T> readRequestItem) {
        super(dataType, readRequestItem);
    }

    public SinglePlcReadRequest(Class<T> dataType, Address address) {
        super(dataType, address);
    }

    public SinglePlcReadRequest(Class<T> dataType, Address address, int size) {
        super(dataType, address, size);
    }
}

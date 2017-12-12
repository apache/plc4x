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
package org.apache.plc4x.java.api.messages.items;

import org.apache.plc4x.java.api.model.Address;

public class WriteRequestItem {

    private final Class datatype;

    private final Address address;

    private final Object[] values;

    public WriteRequestItem(Class datatype, Address address, Object value) {
        this.datatype = datatype;
        this.address = address;
        this.values = new Object[]{value};
    }

    public WriteRequestItem(Class datatype, Address address, Object[] values) {
        this.datatype = datatype;
        this.address = address;
        this.values = values;
    }

    public Class getDatatype() {
        return datatype;
    }

    public Address getAddress() {
        return address;
    }

    public Object[] getValues() {
        return values;
    }

}

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
package org.apache.plc4x.java.s7.messages.items;

import org.apache.plc4x.java.base.messages.items.DefaultStringFieldItem;
import org.apache.plc4x.java.s7.types.S7DataType;

public class S7StringFieldItem extends DefaultStringFieldItem {

    private final S7DataType naturalDataType;

    public S7StringFieldItem(S7DataType naturalDataType, String... values) {
        super(values);
        this.naturalDataType = naturalDataType;
    }

    @Override
    public Object getObject(int index) {
        switch (naturalDataType) {
            case CHAR:
            case WCHAR:
            case STRING:
            case WSTRING:
                return getString(index);
            default:
                return null;
        }
    }

}


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
package org.apache.plc4x.java.s7.netty.model.params;

import org.apache.plc4x.java.s7.netty.model.params.items.VarParameterItem;
import org.apache.plc4x.java.s7.netty.model.types.ParameterType;

import java.util.List;

public class VarParameter implements S7Parameter {

    private final ParameterType type;
    private final List<VarParameterItem> items;

    public VarParameter(ParameterType type, List<VarParameterItem> items) {
        this.type = type;
        this.items = items;
    }

    @Override
    public ParameterType getType() {
        return type;
    }

    public List<VarParameterItem> getItems() {
        return items;
    }

}

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

package org.apache.plc4x.plugins.codegenerator.model.fields;

import org.apache.plc4x.plugins.codegenerator.model.Type;

public class OptionalField implements Field {

    private final Type type;
    private final String name;
    private final String conditionExpression;

    public OptionalField(Type type, String name, String conditionExpression) {
        this.type = type;
        this.name = name;
        this.conditionExpression = conditionExpression;
    }

    public Type getType() {
        return type;
    }

    public String getName() {
        return name;
    }

    public String getConditionExpression() {
        return conditionExpression;
    }

}

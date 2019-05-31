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

import org.apache.plc4x.language.fields.ArrayField;
import org.apache.plc4x.language.references.TypeReference;

public class DefaultArrayField implements ArrayField {

    private final TypeReference type;
    private final String name;
    private final ArrayField.LengthType lengthType;
    private final String lengthExpression;
    private final String[] params;

    public DefaultArrayField(TypeReference type, String name, ArrayField.LengthType lengthType, String lengthExpression, String[] params) {
        this.type = type;
        this.name = name;
        this.lengthType = lengthType;
        this.lengthExpression = lengthExpression;
        this.params = params;
    }

    public TypeReference getType() {
        return type;
    }

    public String getName() {
        return name;
    }

    public ArrayField.LengthType getLengthType() {
        return lengthType;
    }

    public String getLengthExpression() {
        return lengthExpression;
    }

    @Override
    public String[] getParams() {
        return params;
    }

    public static enum LengthType {
        COUNT,
        LENGTH
    }

}

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

public class ArrayField implements Field {

    private final String typeName;
    private final String name;
    private final LengthType lengthType;
    private final String lengthExpression;


    public ArrayField(String typeName, String name, LengthType lengthType, String lengthExpression) {
        this.typeName = typeName;
        this.name = name;
        this.lengthType = lengthType;
        this.lengthExpression = lengthExpression;
    }

    public String getTypeName() {
        return typeName;
    }

    public String getName() {
        return name;
    }

    public LengthType getLengthType() {
        return lengthType;
    }

    public String getLengthExpression() {
        return lengthExpression;
    }

    public static enum LengthType {
        COUNT,
        LENGTH
    }

}

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

package org.apache.plc4x.plugins.codegenerator.model.types;

import org.apache.plc4x.plugins.codegenerator.model.fields.Field;
import org.apache.plc4x.plugins.codegenerator.model.fields.SimpleField;

import javax.management.openmbean.SimpleType;
import java.util.List;
import java.util.stream.Collectors;

public class ComplexType extends Type {

    private final boolean isAbstract;
    private final List<Field> fields;

    public ComplexType(String name, boolean isAbstract, List<Field> fields) {
        super(name);
        this.isAbstract = isAbstract;
        this.fields = fields;
    }

    public boolean isAbstract() {
        return isAbstract;
    }

    public List<Field> getFields() {
        return fields;
    }

    public List<SimpleField> getSimpleFields() {
        return fields.stream().filter(field -> field instanceof SimpleField).map(
            field -> (SimpleField) field).collect(Collectors.toList());
    }

}

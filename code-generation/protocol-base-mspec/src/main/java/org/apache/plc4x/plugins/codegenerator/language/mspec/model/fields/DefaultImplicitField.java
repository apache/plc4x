/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.plugins.codegenerator.language.mspec.model.fields;

import org.apache.plc4x.plugins.codegenerator.types.fields.ImplicitField;
import org.apache.plc4x.plugins.codegenerator.types.references.TypeReference;
import org.apache.plc4x.plugins.codegenerator.types.terms.Term;

import java.util.*;

public class DefaultImplicitField extends DefaultField implements ImplicitField {

    private final TypeReference type;
    private final String name;
    private final Term serializeExpression;

    public DefaultImplicitField(Map<String, Term> attributes, TypeReference type, String name, Term serializeExpression) {
        super(attributes);
        this.type = Objects.requireNonNull(type);
        this.name = Objects.requireNonNull(name);
        this.serializeExpression = Objects.requireNonNull(serializeExpression);
    }

    public TypeReference getType() {
        return type;
    }

    public String getName() {
        return name;
    }

    public Term getSerializeExpression() {
        return serializeExpression;
    }

}

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

package org.apache.plc4x.plugins.codegenerator.language.mspec.model.definitions;


import org.apache.plc4x.plugins.codegenerator.types.definitions.Argument;
import org.apache.plc4x.plugins.codegenerator.types.definitions.TypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.references.DefaultComplexTypeReference;
import org.apache.plc4x.plugins.codegenerator.types.references.TypeReference;

public abstract class DefaultTypeDefinition {

    private final String name;
    private final Argument[] parserArguments;
    private final String[] tags;
    private TypeDefinition parentType;

    public DefaultTypeDefinition(String name, Argument[] parserArguments, String[] tags) {
        this.name = name;
        this.parserArguments = parserArguments;
        this.tags = tags;
        this.parentType = null;
    }

    public String getName() {
        return name;
    }

    public Argument[] getParserArguments() {
        return parserArguments;
    }

    public String[] getTags() {
        return tags;
    }

    public TypeDefinition getParentType() {
        return parentType;
    }

    public void setParentType(TypeDefinition parentType) {
        this.parentType = parentType;
    }

    public TypeReference getTypeReference() {
        return new DefaultComplexTypeReference(getName());
    }

}

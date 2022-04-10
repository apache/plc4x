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
package org.apache.plc4x.codegen.ast;

public class Primitive extends TypeDefinition {

    // Shorthands
    public static final Primitive VOID = new Primitive(DataType.VOID);
    public static final Primitive BOOLEAN = new Primitive(DataType.BOOLEAN);
    public static final Primitive INTEGER = new Primitive(DataType.INTEGER);
    public static final Primitive LONG = new Primitive(DataType.LONG);
    public static final Primitive DOUBLE = new Primitive(DataType.DOUBLE);
    public static final Primitive STRING = new Primitive(DataType.STRING);

    private final DataType type;

    public Primitive(DataType dataType) {
        super(dataType.toString());
        this.type = dataType;
    }

    public DataType getType() {
        return type;
    }

    @Override
    public <T> T accept(NodeVisitor<T> visitor) {
        return null;
    }

    @Override
    public void write(Generator generator) {
        generator.generatePrimitive(this.getType());
    }

    public enum DataType {
        VOID,
        BOOLEAN,
        INTEGER,
        LONG,
        DOUBLE,
        STRING
    }
}

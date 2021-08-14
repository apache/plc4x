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
package org.apache.plc4x.codegen.util;

import org.apache.plc4x.codegen.ast.Expressions;
import org.apache.plc4x.codegen.ast.Method;
import org.apache.plc4x.codegen.ast.Primitive;
import org.apache.plc4x.codegen.ast.TypeDefinition;

import java.util.Collections;

/**
 * This class defines constants necessary for the code generation related to the
 * "Buffer API" which has to be implemented natively.
 */
public class BufferUtil {

    static final TypeDefinition BUFFER_TYPE = Expressions.typeOf("org.apache.plc4x.codegen.api.Buffer");

    // Read Methods
    static final Method READ_UINT8 = new Method(BUFFER_TYPE, "readUint8", Primitive.INTEGER, Collections.emptyList(), Collections.emptyList());
    static final Method READ_UINT16 = new Method(BUFFER_TYPE, "readUint16", Primitive.INTEGER, Collections.emptyList(), Collections.emptyList());
    static final Method READ_UINT32 = new Method(BUFFER_TYPE, "readUint32", Primitive.LONG, Collections.emptyList(), Collections.emptyList());

    private BufferUtil() {
    }


}

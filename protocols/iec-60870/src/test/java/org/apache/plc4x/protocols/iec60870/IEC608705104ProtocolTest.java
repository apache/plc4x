/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.protocols.iec60870;

import org.apache.plc4x.plugins.codegenerator.protocol.TypeContext;
import org.apache.plc4x.protocol.iec608705104.IEC608705104Protocol;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertSame;

class IEC608705104ProtocolTest {

    @Test
    void getTypeContext() throws Exception {
        TypeContext typeContext = new IEC608705104Protocol().getTypeContext();
        assertNotNull(typeContext);
        assertNotNull(typeContext.getUnresolvedTypeReferences());
        assertSame(0, typeContext.getUnresolvedTypeReferences().size());
    }

}

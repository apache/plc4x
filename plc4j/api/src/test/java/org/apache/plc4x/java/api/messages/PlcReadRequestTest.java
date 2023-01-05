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
package org.apache.plc4x.java.api.messages;

import static org.mockito.Mockito.mock;

public class PlcReadRequestTest {

/*    PlcField dummyField;

    @Before
    public void setUp() {
        dummyField = mock(PlcField.class);
    }

    @Test
    public void constructor() {
        new PlcReadRequest();
        new PlcReadRequest(new PlcReadRequestItem<>(String.class, dummyField));
        new PlcReadRequest(String.class, dummyField);
        new PlcReadRequest(String.class, dummyField, 13);
        new PlcReadRequest(Collections.singletonList(new PlcReadRequestItem<>(String.class, dummyField)));
    }

    @Test
    public void builder() {
        { // empty
            try {
                PlcReadRequest.builder().build();
                fail("An empty builder should not be allowed to build a request");
            } catch (IllegalStateException e) {
                // expected
            }
        }
        { // one item
            PlcReadRequest.builder()
                .addField(String.class, dummyField)
                .build();
        }
        { // one item sized
            PlcReadRequest.builder()
                .addField(String.class, dummyField, 13)
                .build();
        }
        { // one item prebuild
            PlcReadRequest.builder()
                .addField(new PlcReadRequestItem<>(String.class, dummyField))
                .build();
        }
        { // two different item
            PlcReadRequest.builder()
                .addField(String.class, dummyField)
                .addField(Byte.class, dummyField)
                .build();
        }
        { // two different item typeSafe
            try {
                PlcReadRequest.builder()
                    .addField(String.class, dummyField)
                    .addField(Byte.class, dummyField)
                    .build(String.class);
                fail("Should not succeed in building with mixed types.");
            } catch (IllegalStateException e) {
                // expected
            }
        }
        { // two different item typeSafe
            try {
                PlcReadRequest.builder()
                    .addField(String.class, dummyField)
                    .addField(Byte.class, dummyField)
                    .build(Byte.class);
                fail("Should not succeed in building with mismatch of types.");
            } catch (ClassCastException e) {
                // expected
            }
        }
        { // two equal item typeSafe
            PlcReadRequest.builder()
                .addField(Byte.class, dummyField)
                .addField(Byte.class, dummyField)
                .build(Byte.class);
        }
    }*/

}
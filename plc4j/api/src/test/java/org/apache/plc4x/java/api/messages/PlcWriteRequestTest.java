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

public class PlcWriteRequestTest {

/*    PlcField dummyField;

    @Before
    public void setUp() {
        dummyField = mock(PlcField.class);
    }

    @Test
    public void constructor() {
        new PlcWriteRequest();
        new PlcWriteRequest(new PlcWriteRequestItem<>(String.class, dummyField, ""));
        new PlcWriteRequest(String.class, dummyField);
        new PlcWriteRequest(Collections.singletonList(new PlcWriteRequestItem<>(String.class, dummyField)));
    }

    @Test
    public void builder() {
        { // empty
            try {
                PlcWriteRequest.builder().build();
                fail("An empty build should fail.");
            } catch (IllegalStateException e) {
                // expected
            }
        }
        { // one item implicit type
            PlcWriteRequest.builder()
                .addField(dummyField, "")
                .build();
        }
        { // one item
            PlcWriteRequest.builder()
                .addField(String.class, dummyField)
                .build();
        }
        { // one item prebuild
            PlcWriteRequest.builder()
                .addField(new PlcWriteRequestItem<>(String.class, dummyField))
                .build();
        }
        { // two different item
            PlcWriteRequest.builder()
                .addField(String.class, dummyField)
                .addField(Byte.class, dummyField)
                .build();
        }
        { // two different item typeSafe
            try {
                PlcWriteRequest.builder()
                    .addField(String.class, dummyField)
                    .addField(Byte.class, dummyField)
                    .build(String.class);
                fail("Mixed types build should fail.");
            } catch (IllegalStateException e) {
                // expected
            }
        }
        { // two different item typeSafe
            try {
                PlcWriteRequest.builder()
                    .addField(String.class, dummyField)
                    .addField(Byte.class, dummyField)
                    .build(Byte.class);
                fail("Mismatch of types should have failed.");
            } catch (ClassCastException e) {
                // expected
            }
        }
        { // two equal item typeSafe
            PlcWriteRequest.builder()
                .addField(Byte.class, dummyField)
                .addField(Byte.class, dummyField)
                .build(Byte.class);
        }
    }

    @Test
    public void testToString() {
        new PlcWriteRequest().toString();
    }*/

}
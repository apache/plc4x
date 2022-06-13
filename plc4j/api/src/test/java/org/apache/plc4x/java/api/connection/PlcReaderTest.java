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
package org.apache.plc4x.java.api.connection;

import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;

import java.util.Collections;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;

import static org.hamcrest.Matchers.equalTo;
import static org.hamcrest.Matchers.instanceOf;
import static org.mockito.Mockito.mock;

public class PlcReaderTest {

/*    @Test
    public void read() throws Exception {
        ((PlcReader) readRequest -> CompletableFuture.completedFuture(new PlcReadResponse(readRequest, Collections.emptyList())))
            .read(new TypeSafePlcReadRequest<>(String.class, mock(PlcField.class))).get();
    }

    @Test
    public void readWrongType() throws Exception {
        try {
            ((PlcReader) readRequest -> CompletableFuture.completedFuture(new PlcReadResponse(readRequest, (List) Collections.singletonList(new PlcReadResponseItem(readRequest.getRequestItem().get(), PlcResponseCode.OK, 1)))))
                .read(new TypeSafePlcReadRequest<>(String.class, mock(PlcField.class))).get();
            fail("Should throw an exception");
        } catch (IllegalArgumentException e) {
            assertThat(e.getMessage(), equalTo("Datatype of 1 doesn't match required datatype of class java.lang.String"));
        } catch (ExecutionException e) {
            assertThat(e.getCause(), instanceOf(IllegalArgumentException.class));
            assertThat(e.getCause().getMessage(), equalTo("Unexpected data type class java.lang.Integer on readRequestItem. Expected class java.lang.String"));
        }
    }*/

}
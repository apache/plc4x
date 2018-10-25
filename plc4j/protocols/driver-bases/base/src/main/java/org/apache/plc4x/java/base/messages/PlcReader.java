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
package org.apache.plc4x.java.base.messages;

import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;

import java.util.concurrent.CompletableFuture;
import java.util.function.Consumer;

/**
 * Interface implemented by all PlcConnections that are able to read from remote resources.
 */
public interface PlcReader {

    /**
     * Reads a requested value from a PLC.
     *
     * @param readRequest object describing the type and location of the value.
     * @return a {@link CompletableFuture} giving async access to the returned value.
     */
    CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest);

}

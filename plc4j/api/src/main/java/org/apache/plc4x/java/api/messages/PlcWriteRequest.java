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
package org.apache.plc4x.java.api.messages;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.concurrent.CompletableFuture;

public interface PlcWriteRequest extends PlcFieldRequest {

    @Override
    CompletableFuture<? extends PlcWriteResponse> execute();

    int getNumberOfValues(String name);

    interface Builder extends PlcRequestBuilder {

        @Override
        PlcWriteRequest build();

        PlcWriteRequest.Builder addItem(String name, String fieldQuery, Boolean... values);

        PlcWriteRequest.Builder addItem(String name, String fieldQuery, Byte... values);

        PlcWriteRequest.Builder addItem(String name, String fieldQuery, Short... values);

        PlcWriteRequest.Builder addItem(String name, String fieldQuery, Integer... values);

        PlcWriteRequest.Builder addItem(String name, String fieldQuery, BigInteger... values);

        PlcWriteRequest.Builder addItem(String name, String fieldQuery, Long... values);

        PlcWriteRequest.Builder addItem(String name, String fieldQuery, Float... values);

        PlcWriteRequest.Builder addItem(String name, String fieldQuery, Double... values);

        PlcWriteRequest.Builder addItem(String name, String fieldQuery, BigDecimal... values);

        PlcWriteRequest.Builder addItem(String name, String fieldQuery, String... values);

        PlcWriteRequest.Builder addItem(String name, String fieldQuery, LocalTime... values);

        PlcWriteRequest.Builder addItem(String name, String fieldQuery, LocalDate... values);

        PlcWriteRequest.Builder addItem(String name, String fieldQuery, LocalDateTime... values);

        PlcWriteRequest.Builder addItem(String name, String fieldQuery, byte[]... values);

        PlcWriteRequest.Builder addItem(String name, String fieldQuery, Byte[]... values);

        <T> PlcWriteRequest.Builder addItem(String name, String fieldQuery, T... values);
    }

}

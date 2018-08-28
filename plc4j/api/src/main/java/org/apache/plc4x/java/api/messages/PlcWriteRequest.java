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

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;

public interface PlcWriteRequest extends PlcFieldRequest {

    interface Builder extends PlcMessageBuilder<PlcWriteRequest> {

        /* NOT Quite sure about these two ...
        PlcReadRequest.Builder addItem(String name, String fieldQuery, byte[]... values);
        PlcReadRequest.Builder addItem(String name, String fieldQuery, Object... values);*/

        PlcReadRequest.Builder addItem(String name, String fieldQuery, Boolean... values);

        PlcReadRequest.Builder addItem(String name, String fieldQuery, Byte... values);

        PlcReadRequest.Builder addItem(String name, String fieldQuery, Short... values);

        PlcReadRequest.Builder addItem(String name, String fieldQuery, Integer... values);

        PlcReadRequest.Builder addItem(String name, String fieldQuery, Long... values);

        PlcReadRequest.Builder addItem(String name, String fieldQuery, Float... values);

        PlcReadRequest.Builder addItem(String name, String fieldQuery, Double... values);

        PlcReadRequest.Builder addItem(String name, String fieldQuery, String... values);

        PlcReadRequest.Builder addItem(String name, String fieldQuery, LocalTime... values);

        PlcReadRequest.Builder addItem(String name, String fieldQuery, LocalDate... values);

        PlcReadRequest.Builder addItem(String name, String fieldQuery, LocalDateTime... values);
    }

}

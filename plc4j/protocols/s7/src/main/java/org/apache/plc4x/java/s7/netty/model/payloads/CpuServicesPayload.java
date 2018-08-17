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
package org.apache.plc4x.java.s7.netty.model.payloads;

import org.apache.plc4x.java.s7.netty.model.payloads.ssls.SslDataRecord;
import org.apache.plc4x.java.s7.netty.model.types.DataTransportErrorCode;
import org.apache.plc4x.java.s7.netty.model.types.ParameterType;
import org.apache.plc4x.java.s7.netty.model.types.SslId;

import java.util.Collections;
import java.util.List;

public class CpuServicesPayload implements S7Payload {

    private DataTransportErrorCode returnCode;
    private SslId sslId;
    private short sslIndex;
    private List<SslDataRecord> sslDataRecords;

    public CpuServicesPayload(DataTransportErrorCode returnCode, SslId sslId, short sslIndex) {
        this.returnCode = returnCode;
        this.sslId = sslId;
        this.sslIndex = sslIndex;
        this.sslDataRecords = Collections.emptyList();
    }

    public CpuServicesPayload(DataTransportErrorCode returnCode, SslId sslId, short sslIndex, List<SslDataRecord> sslDataRecords) {
        this.returnCode = returnCode;
        this.sslId = sslId;
        this.sslIndex = sslIndex;
        this.sslDataRecords = sslDataRecords;
    }

    @Override
    public ParameterType getType() {
        return ParameterType.CPU_SERVICES;
    }

    public DataTransportErrorCode getReturnCode() {
        return returnCode;
    }

    public SslId getSslId() {
        return sslId;
    }

    public short getSslIndex() {
        return sslIndex;
    }

    public List<SslDataRecord> getSslDataRecords() {
        return sslDataRecords;
    }

}

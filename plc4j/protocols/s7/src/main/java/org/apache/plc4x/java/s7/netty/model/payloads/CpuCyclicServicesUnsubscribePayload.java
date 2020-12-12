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

import org.apache.plc4x.java.s7.netty.model.types.DataTransportErrorCode;
import org.apache.plc4x.java.s7.netty.model.types.DataTransportSize;
import org.apache.plc4x.java.s7.netty.model.types.ParameterType;

/**
 *
 * @author cgarcia
 */
public class CpuCyclicServicesUnsubscribePayload implements S7Payload {
    private final DataTransportErrorCode returnCode;
    private final DataTransportSize dataTransportSize;
    private final int length;
    private final byte function;
    private final byte jobId;

    public CpuCyclicServicesUnsubscribePayload(DataTransportErrorCode returnCode, 
            DataTransportSize dataTransportSize, 
            int length, 
            byte function, 
            byte jobId) {
        this.returnCode = returnCode;
        this.dataTransportSize = dataTransportSize;
        this.length = length;
        this.function = function;
        this.jobId = jobId;
    }

    public DataTransportErrorCode getReturnCode() {
        return returnCode;
    }

    public DataTransportSize getDataTransportSize() {
        return dataTransportSize;
    }

    public int getLength() {
        return length;
    }

    public byte getFunction() {
        return function;
    }

    public byte getJobId() {
        return jobId;
    }

    




    @Override
    public ParameterType getType() {
        return ParameterType.CPU_SERVICES;
    }

}

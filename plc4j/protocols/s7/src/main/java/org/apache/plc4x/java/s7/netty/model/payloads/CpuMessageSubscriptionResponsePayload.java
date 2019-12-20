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

import org.apache.plc4x.java.s7.netty.model.types.AlarmType;
import org.apache.plc4x.java.s7.netty.model.types.DataTransportErrorCode;
import org.apache.plc4x.java.s7.netty.model.types.DataTransportSize;
import org.apache.plc4x.java.s7.netty.model.types.ParameterType;

/**
 *
 * @author cgarcia
 */
public class CpuMessageSubscriptionResponsePayload implements S7Payload  {

    private final DataTransportErrorCode returnCode;
    private final DataTransportSize dataTransportSize;
    private final int Length;
    private final byte Result;
    private final AlarmType alarm;

    public CpuMessageSubscriptionResponsePayload(DataTransportErrorCode returnCode, 
            DataTransportSize dataTransportSize, 
            int Length, 
            byte Result, 
            AlarmType alarm) {
        
        this.returnCode = returnCode;
        this.dataTransportSize = dataTransportSize;
        this.Length = Length;
        this.Result = Result;
        this.alarm = alarm;
    }
    
    
    
    @Override
    public ParameterType getType() {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }
    
}

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
public class CpuMessageSubscriptionServicePayload implements S7Payload {

    private final DataTransportErrorCode returnCode;
    private final DataTransportSize dataTransportSize;
    private final byte SubscribedEvents;
    private final String Id; //8 bytes string
    private final AlarmType alarm;


    public CpuMessageSubscriptionServicePayload (DataTransportErrorCode returnCode,
            DataTransportSize dataTransportSize, 
            byte SubscribedEvents, 
            String Id, 
            AlarmType alarm) {
        
        this.returnCode = returnCode;
        this.dataTransportSize = dataTransportSize;
        this.SubscribedEvents = SubscribedEvents;
        this.Id = Id;
        this.alarm = alarm;
    }
    
    @Override
    public ParameterType getType() {
        return ParameterType.CPU_SERVICES;
    }

    public DataTransportErrorCode getReturnCode() {
        return returnCode;
    }

    public DataTransportSize getDataTransportSize() {
        return dataTransportSize;
    }

    public byte getSubscribedEvents() {
        return SubscribedEvents;
    }

    public String getId() {
        return Id;
    }
  
    public AlarmType getAlarm() {
        return alarm;
    }    
    
}

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

import org.apache.plc4x.java.s7.netty.model.messages.S7PushMessage;
import org.apache.plc4x.java.s7.netty.model.payloads.items.AlarmMessageItem;
import org.apache.plc4x.java.s7.netty.model.types.DataTransportErrorCode;
import org.apache.plc4x.java.s7.netty.model.types.DataTransportSize;
import org.apache.plc4x.java.s7.netty.model.types.ParameterType;

/**
 *
 * @author cgarcia
 */
public class AlarmMessagePayload implements S7Payload, S7PushMessage {
        
    private final DataTransportErrorCode returnCode;
    private final DataTransportSize dataTransportSize;
    private final Object msgtype; 
    private final Integer length;
    private final AlarmMessageItem msg;    

    public AlarmMessagePayload(DataTransportErrorCode returnCode, 
            DataTransportSize dataTransportSize, 
            Integer length,
            AlarmMessageItem msg) {
        this.returnCode = returnCode;
        this.dataTransportSize = dataTransportSize;
        this.msgtype = null;
        this.length = length;
        this.msg = msg;
    }    
    
    public AlarmMessagePayload(DataTransportErrorCode returnCode, 
            DataTransportSize dataTransportSize, 
            Object msgtype, 
            Integer length,
            AlarmMessageItem msg) {
        this.returnCode = returnCode;
        this.dataTransportSize = dataTransportSize;
        this.msgtype = msgtype;
        this.length = length;
        this.msg = msg;
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

    public Object getMsgtype() {
        return msgtype;
    }

    public Integer getLength() {
        return length;
    }

    public AlarmMessageItem getMsg() {
        return msg;
    }

    @Override
    public String toString() {
        return "AlarmMessagePayload{" + "returnCode=" + returnCode 
                + ", dataTransportSize=" + dataTransportSize 
                + ", msgtype=" + msgtype 
                + ", length=" + length 
                + ", msg=" + msg 
                + '}';
    }
    
    
    
}

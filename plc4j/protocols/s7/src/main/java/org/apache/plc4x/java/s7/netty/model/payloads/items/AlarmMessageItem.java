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

package org.apache.plc4x.java.s7.netty.model.payloads.items;

import java.time.LocalDateTime;
import java.util.List;
import org.apache.plc4x.java.s7.netty.model.types.DataTransportErrorCode;
import org.apache.plc4x.java.s7.netty.model.types.DataTransportSize;

/**
 *
 * @author cgarcia
 */
public class AlarmMessageItem {
    private final DataTransportErrorCode returnCode;
    private final DataTransportSize dataTransportSize;
    private final LocalDateTime timestamp;
    private final byte function;
    private final byte objects;
    private final int Length;
    private final List<MessageObjectItem> msgItems; 

    public AlarmMessageItem(LocalDateTime timestamp, 
            byte function, 
            byte objects, 
            List<MessageObjectItem> msgItems) {
        
        this.timestamp = timestamp;
        this.function = function;
        this.objects = objects;
        this.msgItems = msgItems;
        this.returnCode = null;
        this.dataTransportSize = null;
        this.Length = 0;
    }
    
    public AlarmMessageItem(byte function,
            byte objects,
            DataTransportErrorCode returnCode,
            DataTransportSize dataTransportSize,
            int Length,
            List<MessageObjectItem> msgItems){
        this.function = function;
        this.objects = objects;
        this.returnCode = returnCode;
        this.dataTransportSize = dataTransportSize;
        this.Length = Length;
        this.timestamp = null; 
        this.msgItems = msgItems;
    };
            
            

    public DataTransportErrorCode getReturnCode() {
        return returnCode;
    }

    public DataTransportSize getDataTransportSize() {
        return dataTransportSize;
    }

    public LocalDateTime getTimestamp() {
        return timestamp;
    }

    public byte getFunction() {
        return function;
    }

    public byte getObjects() {
        return objects;
    }

    public int getLength() {
        return Length;
    }

    public List<MessageObjectItem> getMsgItems() {
        return msgItems;
    }

    @Override
    public String toString() {
        return "AlarmMessageItem{" + "returnCode=" + returnCode 
                + ", dataTransportSize=" + dataTransportSize 
                + ", timestamp=" + timestamp 
                + ", function=" + function 
                + ", objects=" + objects 
                + ", Length=" + Length 
                + ", msgItems=" + msgItems 
                + '}';
    }
    
    
}

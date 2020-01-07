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

/**
 * Diagnostic message from PLC, is transfer by payload.
 * EventID: System specific messages W#16#8xyz and W#16#9xyz.
 *          User messages W#16#Axyz and W#16#Bxyz
 *          where x=0 output event, x=1 input event
 *          and yz hex number in message handler.
 * 
 * @author cgarcia
 */
public class CpuDiagnosticMessageItem {
    
    private final Short EventID;
    private final Byte PriorityClass;
    private final Byte ObNumber;
    private final Short DatID;   
    private final Short Info1;   
    private final Integer Info2; 
    private LocalDateTime TimeStamp;
    
    public CpuDiagnosticMessageItem(Short EventID, Byte PriorityClass, Byte ObNumber, Short DatID, Short Info1, Integer Info2, LocalDateTime TimeStamp) {
        this.EventID = EventID;
        this.PriorityClass = PriorityClass;
        this.ObNumber = ObNumber;
        this.DatID = DatID;
        this.Info1 = Info1;
        this.Info2 = Info2;
        this.TimeStamp = TimeStamp;
    }

    public Short getEventID() {
        return EventID;
    }

    public Byte getPriorityClass() {
        return PriorityClass;
    }

    public Byte getObNumber() {
        return ObNumber;
    }

    public Short getDatID() {
        return DatID;
    }

    public Short getInfo1() {
        return Info1;
    }

    public Integer getInfo2() {
        return Info2;
    }

    public LocalDateTime getTimeStamp() {
        return TimeStamp;
    }

    @Override
    public String toString() {
        return "CpuDiagnosticMessageItem{" + "EventID=" + EventID 
                + ", PriorityClass=" + PriorityClass 
                + ", ObNumber=" + ObNumber 
                + ", DatID=" + DatID 
                + ", Info1=" + Info1 
                + ", Info2=" + Info2 
                + ", TimeStamp=" + TimeStamp 
                + '}';
    }



}

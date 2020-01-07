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

import java.time.format.DateTimeFormatter;
import java.util.Calendar;
import java.util.List;

/**
 *
 * @author cgarcia
 */
public class MessageObjectItem {
    
    private final byte VariableSpecification;
    private final byte Length;
    private final Object SyntaxID;
    private final byte NumberOfValues;
    private final int EventID;
    private final byte EventState;
    private final byte State;
    private final byte AckStateGoing;
    private final byte AckStateComming;
    private final Calendar TimestampComing;
    private final Calendar TimestampGoing;
    private final List<AssociatedValueItem> ComingValues;
    private final List<AssociatedValueItem> GoingValues;     

    public MessageObjectItem(byte VariableSpecification, 
            byte Length, 
            Object SyntaxID, 
            byte NumberOfValues, 
            int EventID, 
            byte EventState, 
            byte State, 
            byte AckStateGoing, 
            byte AckStateComming, 
            List<AssociatedValueItem> Values) {
        
        this.VariableSpecification = VariableSpecification;
        this.Length = Length;
        this.SyntaxID = SyntaxID;
        this.NumberOfValues = NumberOfValues;
        this.EventID = EventID;
        this.EventState = EventState;
        this.State = State;
        this.AckStateGoing = AckStateGoing;
        this.AckStateComming = AckStateComming;
        this.TimestampComing = null;
        this.ComingValues = Values;
        this.TimestampGoing = null;               
        this.GoingValues = null;
    }
    
    public MessageObjectItem(byte VariableSpecification, 
            byte Length, 
            Object SyntaxID, 
            byte NumberOfValues, 
            int EventID, 
            byte EventState, 
            byte State, 
            byte AckStateGoing, 
            byte AckStateComming,
            Calendar TimestampComing,
            List<AssociatedValueItem> ComingValues,
            Calendar TimestampGoing,
            List<AssociatedValueItem> GoingValues) {
        
        this.VariableSpecification = VariableSpecification;
        this.Length = Length;
        this.SyntaxID = SyntaxID;
        this.NumberOfValues = NumberOfValues;
        this.EventID = EventID;
        this.EventState = EventState;
        this.State = State;
        this.AckStateGoing = AckStateGoing;
        this.AckStateComming = AckStateComming;
        this.TimestampComing = TimestampComing;
        this.ComingValues = ComingValues;
        this.TimestampGoing = TimestampGoing;               
        this.GoingValues = GoingValues;
    }    

    public byte getVariableSpecification() {
        return VariableSpecification;
    }

    public byte getLength() {
        return Length;
    }

    public Object getSyntaxID() {
        return SyntaxID;
    }

    public byte getNumberOfValues() {
        return NumberOfValues;
    }

    public int getEventID() {
        return EventID;
    }

    public byte getEventState() {
        return EventState;
    }

    public byte getState() {
        return State;
    }

    public byte getAckStateGoing() {
        return AckStateGoing;
    }

    public byte getAckStateComming() {
        return AckStateComming;
    }

    public Calendar getTimestampComing() {
        return TimestampComing;
    }

    public Calendar getTimestampGoing() {
        return TimestampGoing;
    }

    public List<AssociatedValueItem> getComingValues() {
        return ComingValues;
    }

    public List<AssociatedValueItem> getGoingValues() {
        return GoingValues;
    }

    @Override
    public String toString() {
        DateTimeFormatter formatter = DateTimeFormatter.ofPattern("yyyy/MM/dd HH:mm:ss.SSS");
        return "MessageObjectItem{" + "VariableSpecification=" + VariableSpecification 
                + ", Length=" + Length 
                + ", SyntaxID=" + SyntaxID 
                + ", NumberOfValues=" + NumberOfValues 
                + ", EventID=" + EventID 
                + ", EventState=" + EventState 
                + ", State=" + State 
                + ", AckStateGoing=" + AckStateGoing 
                + ", AckStateComming=" + AckStateComming 
                + ", TimestampComing=" + TimestampComing
                + ", TimestampGoing="  + TimestampGoing 
                + ", ComingValues=" + ComingValues 
                + ", GoingValues=" + GoingValues 
                + '}';
    }



    
 
    
    
}

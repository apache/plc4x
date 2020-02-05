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
import org.apache.plc4x.java.s7.netty.model.types.AlarmQueryType;
import org.apache.plc4x.java.s7.netty.model.types.QueryType;
import org.apache.plc4x.java.s7.netty.model.types.VariableAddressingMode;

/**
 *
 * @author cgarcia
 */
public class MessageObjectItem {
    
    private final byte VariableSpecification;
    private final byte Length;
    private final VariableAddressingMode SyntaxID;
    private final byte NumberOfValues;
    private final int EventID;
    private final byte EventState;
    private final byte State;
    private final byte AckStateGoing;
    private final byte AckStateComming;
    private final byte EventGoing;
    private final byte EventComming;
    private final byte EventLastChange;
    private final LocalDateTime TimestampComing;
    private final LocalDateTime TimestampGoing;
    private final List<AssociatedValueItem> ComingValues;
    private final List<AssociatedValueItem> GoingValues;  
    //For Query alarms.
    private final QueryType querytype;
    private final AlarmQueryType alarmtype;
     
    //For Query request
    public MessageObjectItem(byte VariableSpecification, 
            byte Length, 
            VariableAddressingMode SyntaxID, 
            QueryType querytype, 
            AlarmQueryType alarmtype) {
        
        this.VariableSpecification = VariableSpecification;
        this.Length = Length;
        this.SyntaxID = SyntaxID;
        this.NumberOfValues = 0x00;
        this.EventID = 0x0000;
        this.EventState = 0x00;
        this.State = 0x00;
        this.AckStateGoing = 0x00;
        this.AckStateComming = 0x00;
        this.EventGoing = 0x00;
        this.EventComming = 0x00;
        this.EventLastChange = 0x00;
        this.TimestampComing = null;
        this.ComingValues = null;
        this.TimestampGoing = null;               
        this.GoingValues = null;
        this.querytype = querytype;
        this.alarmtype = alarmtype;
    }   
    
    //For Query reply
    public MessageObjectItem(byte Length, 
            AlarmQueryType alarmtype, 
            int EventID, 
            byte EventState, 
            byte AckStateGoing, 
            byte AckStateComming, 
            LocalDateTime TimestampComing,
            List<AssociatedValueItem> ComingValues,
            LocalDateTime TimestampGoing, 
            List<AssociatedValueItem> GoingValues) {
        
        this.VariableSpecification = 0x00;
        this.Length = Length;
        this.SyntaxID = null;
        this.NumberOfValues = 0x00;
        this.EventID = EventID;
        this.EventState = EventState;
        this.State = 0x00;
        this.AckStateGoing = AckStateGoing;
        this.AckStateComming = AckStateComming;
        this.EventGoing = 0x00;
        this.EventComming = 0x00;
        this.EventLastChange = 0x00;        
        this.TimestampComing = TimestampComing;
        this.ComingValues = ComingValues;
        this.TimestampGoing = TimestampGoing;               
        this.GoingValues = GoingValues;
        this.querytype = null;
        this.alarmtype = alarmtype;        
    }    

    public MessageObjectItem(byte VariableSpecification, 
            byte Length, 
            VariableAddressingMode SyntaxID, 
            byte NumberOfValues, 
            int EventID, 
            byte AckStateGoing, 
            byte AckStateComming) {
        
        this.VariableSpecification = VariableSpecification;
        this.Length = Length;
        this.SyntaxID = SyntaxID;
        this.NumberOfValues = NumberOfValues;
        this.EventID = EventID;
        this.EventState = 0x00;
        this.State = 0x00;
        this.AckStateGoing = AckStateGoing;
        this.AckStateComming = AckStateComming;
        this.EventGoing = 0x00;
        this.EventComming = 0x00;
        this.EventLastChange = 0x00;        
        this.TimestampComing = null;
        this.ComingValues = null;
        this.TimestampGoing = null;               
        this.GoingValues = null;
        this.querytype = null;
        this.alarmtype = null;
    }    
    
    public MessageObjectItem(byte VariableSpecification, 
            byte Length, 
            VariableAddressingMode SyntaxID, 
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
        this.EventGoing = 0x00;
        this.EventComming = 0x00;
        this.EventLastChange = 0x00;        
        this.TimestampComing = null;
        this.ComingValues = Values;
        this.TimestampGoing = null;               
        this.GoingValues = null;
        this.querytype = null;
        this.alarmtype = null;        
    }
    
    public MessageObjectItem(byte VariableSpecification, 
            byte Length, 
            VariableAddressingMode SyntaxID, 
            byte NumberOfValues, 
            int EventID, 
            byte EventState, 
            byte State, 
            byte AckStateGoing, 
            byte AckStateComming,
            LocalDateTime TimestampComing,
            List<AssociatedValueItem> ComingValues,
            LocalDateTime TimestampGoing,
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
        this.EventGoing = 0x00;
        this.EventComming = 0x00;
        this.EventLastChange = 0x00;        
        this.TimestampComing = TimestampComing;
        this.ComingValues = ComingValues;
        this.TimestampGoing = TimestampGoing;               
        this.GoingValues = GoingValues;
        this.querytype = null;
        this.alarmtype = null;        
    }    
    
    public MessageObjectItem(byte VariableSpecification, 
            byte Length, 
            VariableAddressingMode SyntaxID, 
            byte NumberOfValues, 
            int EventID, 
            byte EventState, 
            byte State, 
            byte AckStateGoing, 
            byte AckStateComming,
            byte EventGoing,
            byte EventComming,
            byte EventLastChange,
            List<AssociatedValueItem> ComingValues) {
        
        this.VariableSpecification = VariableSpecification;
        this.Length = Length;
        this.SyntaxID = SyntaxID;
        this.NumberOfValues = NumberOfValues;
        this.EventID = EventID;
        this.EventState = EventState;
        this.State = State;
        this.AckStateGoing = AckStateGoing;
        this.AckStateComming = AckStateComming;
        this.EventGoing = 0x00;
        this.EventComming = 0x00;
        this.EventLastChange = 0x00;        
        this.TimestampComing = null;
        this.ComingValues = ComingValues;
        this.TimestampGoing = null;               
        this.GoingValues = null;
        this.querytype = null;
        this.alarmtype = null;        
    }     
    
    public byte getVariableSpecification() {
        return VariableSpecification;
    }

    public byte getLength() {
        return Length;
    }

    public VariableAddressingMode getSyntaxID() {
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

    public LocalDateTime getTimestampComing() {
        return TimestampComing;
    }

    public byte getEventGoing() {
        return EventGoing;
    }

    public byte getEventComming() {
        return EventComming;
    }

    public byte getEventLastChange() {
        return EventLastChange;
    }
    
    public LocalDateTime getTimestampGoing() {
        return TimestampGoing;
    }

    public List<AssociatedValueItem> getComingValues() {
        return ComingValues;
    }

    public List<AssociatedValueItem> getGoingValues() {
        return GoingValues;
    }

    public QueryType getQuerytype() {
        return querytype;
    }

    public AlarmQueryType getAlarmtype() {
        return alarmtype;
    }

    @Override
    public String toString() {
        return "MessageObjectItem{" + "VariableSpecification=" + VariableSpecification +
                ", Length=" + Length +
                ", SyntaxID=" + SyntaxID +
                ", NumberOfValues=" + NumberOfValues +
                ", EventID=" + EventID +
                ", EventState=" + EventState +
                ", State=" + State +
                ", AckStateGoing=" + AckStateGoing +
                ", AckStateComming=" + AckStateComming +
                ", TimestampComing=" + TimestampComing +
                ", TimestampGoing=" + TimestampGoing +
                ", ComingValues=" + ComingValues +
                ", GoingValues=" + GoingValues +
                ", querytype=" + querytype +
                ", alarmtype=" + alarmtype + '}';
    }
        


    
}

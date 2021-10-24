/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.s7.readwrite.field;

import java.util.ArrayList;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.s7.readwrite.types.AlarmType;
import org.apache.plc4x.java.s7.readwrite.types.EventType;
import org.apache.plc4x.java.s7.readwrite.types.S7SubscriptionFieldType;
import org.apache.plc4x.java.s7.readwrite.types.TimeBase;

/**
 *
 * @author cgarcia
 */
public class S7SubscriptionField implements PlcField {

    //Event Subscription 
    private static final Pattern EVENT_SUBSCRIPTION_TYPE_PATTERN = 
        //Pattern.compile("(^MODE)|(^SYS)|(^USR)|(^ALM_S)|(^ALM_8)");
        Pattern.compile("(^MODE)|(^SYS)|(^USR)|(^ALM)");
    //Event ack
    private static final Pattern EVENT_ALARM_ACK_PATTERN = 
        Pattern.compile("(^XACK:)(((?:,{0,1})(16#[0-9a-fA-F]{8})(;([0-9a-fA-F]{2})))+)");
    
    //Query alarms from PLC.
    //TODO: Query SCAN 
    private static final Pattern EVENT_ALARM_QUERY_PATTERN = 
        Pattern.compile("(^QUERY:)((ALARM_S)|(ALARM_8))");    
        
      

    
    //byteOffset theoretically can reach up to 2097151 ... see checkByteOffset() below --> 7digits
    private static final Pattern ADDRESS_PATTERN =
        Pattern.compile("%(?<memoryArea>.)(?<transferSizeCode>[XBWD]?)(?<byteOffset>\\d{1,7})(.(?<bitOffset>[0-7]))?:(?<dataType>[a-zA-Z_]+)(\\[(?<numElements>\\d+)])?");

    //blockNumber usually has its max hat around 64000 --> 5digits
    private static final Pattern DATA_BLOCK_ADDRESS_PATTERN =
        Pattern.compile("%DB(?<blockNumber>\\d{1,5}).DB(?<transferDBSizeCode>[XBWD]?)(?<byteDBOffset>\\d{1,7})(.(?<bitDBOffset>[0-7]))?:(?<dataDBType>[a-zA-Z_]+)(\\[(?<numDBElements>\\d+)])?"); 
    
    //All fields index 9
    private static final Pattern EVENT_SUBSCRIPTION_S7ANY_QUERY_PATTERN =
        Pattern.compile("(^CYC(\\((?<timebase>((B01SEC)|(B1SEC)|(B10SEC)):(?<multiplier>[1-99]))\\)):)(((?:,{0,1})(("+ ADDRESS_PATTERN + ")|("+ DATA_BLOCK_ADDRESS_PATTERN +")))+)");      
       
    private static final Pattern EVENT_SUBSCRIPTION_DB_QUERY_PATTERN =     
        Pattern.compile("(^CYC(\\((?<timebase>((B01SEC)|(B1SEC)|(B10SEC)):(?<multiplier>[1-99]))\\)):)(((?:,{0,1})(%DB(?<blockNumber>\\d{1,5}).DB(?<transferSizeCode>[B]?)(?<byteOffset>\\d{1,7})(\\[(?<numElements>\\d+)]))?)+)");
    
    private final S7SubscriptionFieldType fieldtype;
    private final EventType eventtype;
    private final S7Field[] s7fields;    
    private final ArrayList<Integer> ackalarms;
    private final AlarmType alarmquerytype;
    private final TimeBase timebase;
    private final short multiplier;
    
    public S7SubscriptionField(S7SubscriptionFieldType fieldtype, EventType eventtype){
        this.fieldtype  = fieldtype;
        this.eventtype  = eventtype;
        this.s7fields   = null;
        this.ackalarms  = null;
        this.alarmquerytype = null;
        this.timebase   = null;
        this.multiplier = 0;
    }
       
    public S7SubscriptionField(S7SubscriptionFieldType fieldtype, ArrayList<Integer> ackalarms){
        this.fieldtype  = fieldtype;
        this.eventtype  = null;
        this.s7fields   = null;
        this.ackalarms  = ackalarms;
        this.alarmquerytype = null;
        this.timebase   = null;
        this.multiplier = 0;        
    }    

    public S7SubscriptionField(S7SubscriptionFieldType fieldtype, AlarmType alarmquerytype){
        this.fieldtype  = fieldtype;
        this.eventtype  = null;
        this.s7fields   = null;
        this.ackalarms  = null;  
        this.alarmquerytype = alarmquerytype;  
        this.timebase   = null;
        this.multiplier = 0;        
    }     
    
    public S7SubscriptionField(S7SubscriptionFieldType fieldtype, S7Field[] s7field, TimeBase timebase, short multiplier){
        this.fieldtype  = fieldtype;
        this.eventtype  = null;
        this.s7fields   = s7field;
        this.ackalarms  = null;  
        this.alarmquerytype = null;
        this.timebase   = timebase;
        this.multiplier = multiplier;
    }     

    public S7SubscriptionFieldType getFieldtype() {
        return fieldtype;
    }

    public EventType getEventtype() {
        return eventtype;
    }

    public S7Field[] getS7field() {
        return s7fields;
    }

    public ArrayList<Integer> getAckalarms() {
        return ackalarms;
    }

    public AlarmType getAlarmquerytype() {
        return alarmquerytype;
    }
    
    public TimeBase getTimeBase() {
        return timebase;
    }

    public short getMultiplier() {
        return multiplier;
    }    
    
           
    public static boolean matches(String fieldString) {
        return EVENT_SUBSCRIPTION_TYPE_PATTERN.matcher(fieldString).matches() ||
            EVENT_ALARM_ACK_PATTERN.matcher(fieldString).matches() ||
            EVENT_ALARM_QUERY_PATTERN.matcher(fieldString).matches() ||
            EVENT_SUBSCRIPTION_S7ANY_QUERY_PATTERN.matcher(fieldString).matches() ||
            EVENT_SUBSCRIPTION_DB_QUERY_PATTERN.matcher(fieldString).matches(); 

    }     
    
    public static S7SubscriptionField of(String fieldString) {
        {
            Matcher matcher = EVENT_SUBSCRIPTION_TYPE_PATTERN.matcher(fieldString);
            if (matcher.matches()){
                return new S7SubscriptionField(S7SubscriptionFieldType.EVENT_SUBSCRIPTION,
                            EventType.valueOf(fieldString));
            }
        }
        
        {
            //TODO: Actually only ALARM_S (SIG_1)
            Matcher matcher = EVENT_ALARM_ACK_PATTERN.matcher(fieldString);
            if (matcher.matches()){
                String[] arrIdAndSig;
                String EventIds = matcher.group(2);
                String[] arrStrEventId = EventIds.split(",");
                ArrayList<Integer> arrEventId = new ArrayList<>();
                for (String EventId:arrStrEventId){
                    EventId = EventId.replaceAll("16#", "");
                    arrIdAndSig =  EventId.split(";");
                    arrEventId.add(Integer.parseInt(arrIdAndSig[0], 16));
                    arrEventId.add(Integer.parseInt(arrIdAndSig[1], 16));                    
                }
                return new S7SubscriptionField(S7SubscriptionFieldType.ALARM_ACK,
                            arrEventId);
                
            }            
        }
        
        {
            //TODO: Support for ALARM_8            
            Matcher matcher = EVENT_ALARM_QUERY_PATTERN.matcher(fieldString); 
            TimeBase tb = null;
            short multi = 0;            
            if (matcher.matches()){

                return new S7SubscriptionField(S7SubscriptionFieldType.ALARM_QUERY,
                            AlarmType.ALARM_S);    
            }
        }
        
        {
            Matcher matcher = EVENT_SUBSCRIPTION_DB_QUERY_PATTERN.matcher(fieldString);
            TimeBase tb = null;
            short multi = 0;            
            if (matcher.matches()){
                String strAddress = matcher.group(9);
                strAddress = strAddress.replaceAll("\\[", ".0:BYTE[");
                String[] dbAddress = strAddress.split(",");       
                S7Field[] s7fields = new S7Field[dbAddress.length];
                int i=0;
                for (String address:dbAddress) {
                    s7fields[i] = S7Field.of(address);
                    i++;
                }                
                return new S7SubscriptionField(S7SubscriptionFieldType.CYCLIC_DB_SUBSCRIPTION,
                            s7fields,
                            tb,
                            multi);    
            }            
            
        }                
        
        {
            Matcher matcher = EVENT_SUBSCRIPTION_S7ANY_QUERY_PATTERN.matcher(fieldString);
            TimeBase tb = null;
            short multi = 0;
            if (matcher.matches()){                    
                S7Field[] myFields = null;
                String strAddress = matcher.group(9);  
                String[] fieldAddress = strAddress.split(","); 
                myFields = new S7Field[fieldAddress.length];            
                int i=0;
                for (String address:fieldAddress) {
                    myFields[i] = S7Field.of(address);
                    i++;
                }         
                return new S7SubscriptionField(S7SubscriptionFieldType.CYCLIC_SUBSCRIPTION,
                            myFields,
                tb,
                multi);                
            }
        }
        
        throw new PlcInvalidFieldException("Unable to parse address: " + fieldString);        
    }    
        
}

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
package org.apache.plc4x.java.s7.readwrite.field;

import java.time.Duration;
import java.util.ArrayList;
import java.util.Optional;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.model.PlcSubscriptionField;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.s7.readwrite.types.AlarmType;
import org.apache.plc4x.java.s7.readwrite.types.EventType;
import org.apache.plc4x.java.s7.readwrite.types.S7SubscriptionFieldType;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionField;

/**
 *
 * @author cgarcia
 */
public class S7SubscriptionField implements PlcField {

    //Event Subscription 
    private static final Pattern EVENT_SUBSCRIPTION_TYPE_PATTERN = 
        Pattern.compile("(^MODE)|(^SYS)|(^USR)|(^ALM_S)|(^ALM_8)");
    
    //Event ack
    private static final Pattern EVENT_ALARM_ACK_PATTERN = 
        Pattern.compile("(^ACK:)(((?:,{0,1})(16#[0-9a-fA-F]{8}))+)");
    
    //Query alarms from PLC.
    //TODO: Query SCAN 
    private static final Pattern EVENT_ALARM_QUERY_PATTERN = 
        Pattern.compile("(^QUERY:)((ALARM_S)|(ALARM_8))");    
       
    private final S7SubscriptionFieldType fieldtype;
    private final EventType eventtype;
    private final S7Field s7field;    
    private final ArrayList<Integer> ackalarms;
    private final AlarmType alarmquerytype;
    
    public S7SubscriptionField(S7SubscriptionFieldType fieldtype, EventType eventtype){
        this.fieldtype = fieldtype;
        this.eventtype = eventtype;
        this.s7field = null;
        this.ackalarms = null;
        this.alarmquerytype = null;
    }
       
    public S7SubscriptionField(S7SubscriptionFieldType fieldtype, ArrayList<Integer> ackalarms){
        this.fieldtype = fieldtype;
        this.eventtype = null;
        this.s7field   = null;
        this.ackalarms = ackalarms;
        this.alarmquerytype = null;
    }    

    public S7SubscriptionField(S7SubscriptionFieldType fieldtype, AlarmType alarmquerytype){
        this.fieldtype = fieldtype;
        this.eventtype = null;
        this.s7field = null;
        this.ackalarms = null;  
        this.alarmquerytype = alarmquerytype;        
    }     
    
    public S7SubscriptionField(S7SubscriptionFieldType fieldtype, S7Field s7field){
        this.fieldtype = fieldtype;
        this.eventtype = null;
        this.s7field = s7field;
        this.ackalarms = null;  
        this.alarmquerytype = null;        
    }     

    public S7SubscriptionFieldType getFieldtype() {
        return fieldtype;
    }

    public EventType getEventtype() {
        return eventtype;
    }

    public S7Field getS7field() {
        return s7field;
    }

    public ArrayList<Integer> getAckalarms() {
        return ackalarms;
    }

    public AlarmType getAlarmquerytype() {
        return alarmquerytype;
    }
           
    public static boolean matches(String fieldString) {
        return EVENT_SUBSCRIPTION_TYPE_PATTERN.matcher(fieldString).matches() ||
            EVENT_ALARM_ACK_PATTERN.matcher(fieldString).matches() ||
            EVENT_ALARM_QUERY_PATTERN.matcher(fieldString).matches() ||
                S7Field.matches(fieldString);
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
                String EventIds = matcher.group(2);
                String[] arrStrEventId = EventIds.split(",");
                ArrayList<Integer> arrEventId = new ArrayList<>();
                for (String EventId:arrStrEventId){
                    EventId = EventId.replaceAll("16#", "");
                    arrEventId.add(Integer.parseInt(EventId, 16));                    
                }
                return new S7SubscriptionField(S7SubscriptionFieldType.ALARM_ACK,
                            arrEventId);
                
            }            
        }
        
        {
            //TODO: Support for ALARM_8            
            Matcher matcher = EVENT_ALARM_QUERY_PATTERN.matcher(fieldString); 
            if (matcher.matches()){
            return new S7SubscriptionField(S7SubscriptionFieldType.ALARM_QUERY,
                            AlarmType.ALARM_S);    
            }
        }
        
        {
            if (S7Field.matches(fieldString)){
                S7Field s7field = S7Field.of(fieldString);
                switch(s7field.getDataType()) {
                    case BYTE:
                        
                    default:;
                    
                }
                return new S7SubscriptionField(S7SubscriptionFieldType.CYCLIC_SUBSCRIPTION,
                            s7field);                
            }
        }
        throw new PlcInvalidFieldException("Unable to parse address: " + fieldString);        
    }    
        
}

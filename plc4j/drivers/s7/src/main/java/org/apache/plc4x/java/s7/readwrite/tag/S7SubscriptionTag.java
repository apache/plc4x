/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.s7.readwrite.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.s7.readwrite.AlarmType;
import org.apache.plc4x.java.s7.readwrite.EventType;
import org.apache.plc4x.java.s7.readwrite.types.S7SubscriptionType;
import org.apache.plc4x.java.s7.readwrite.TimeBase;
        
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class S7SubscriptionTag implements PlcTag {

    //Event Subscription 
    private static final Pattern EVENT_SUBSCRIPTION_TYPE_PATTERN =
        //Pattern.compile("(^MODE)|(^SYS)|(^USR)|(^ALM_S)|(^ALM_8)");
        Pattern.compile("(^MODE)|(^SYS)|(^USR)|(^ALM)");
    //Event ack
    private static final Pattern EVENT_ALARM_ACK_PATTERN =
        Pattern.compile("(^ACK:)(((?:,{0,1})(16#[0-9a-fA-F]{8}))+)");

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
        Pattern.compile("(^CYC(\\((?<timeBase>((B01SEC)|(B1SEC)|(B10SEC))):(?<multiplier>[1-99])\\)):)(((?:,{0,1})((" + ADDRESS_PATTERN + ")|(" + DATA_BLOCK_ADDRESS_PATTERN + ")))+)");

    private static final Pattern EVENT_SUBSCRIPTION_DB_QUERY_PATTERN =
        Pattern.compile("(^CYC(\\((?<timeBase>((B01SEC)|(B1SEC)|(B10SEC))):(?<multiplier>[1-99])\\)):)(((?:,{0,1})(%DB(?<blockNumber>\\d{1,5}).DB(?<transferDBSizeCode>[B]?)(?<byteDBOffset>\\d{1,7})(\\[(?<numDBElements>\\d+)]))?)+)");

    private static final Pattern EVENT_CANCEL_JOB_QUERY_PATTERN =
        Pattern.compile("(^CANCEL:)((((?:,{0,1})(\\d+{1,3})))+)");    
    
    private static final String MEMORY_AREA = "memoryArea";
    private static final String TRANSFER_SIZE_CODE = "transferSizeCode";
    private static final String BYTE_OFFSET = "byteOffset";
    private static final String BIT_OFFSET = "bitOffset";
    private static final String DATA_TYPE = "dataType";
    private static final String NUM_ELEMENTS = "numElements";

    private static final String BLOCK_NUMBER = "blockNumber";
    private static final String TRANSFER_DB_SIZE_CODE = "transferSizeCode";
    private static final String BYTE_DB_OFFSET = "byteOffset";
    private static final String BIT_DB_OFFSET = "bitOffset";
    private static final String DATA_DB_TYPE = "dataType";
    private static final String NUM_DB_ELEMENTS = "numElements";

    private static final String STRING_LENGTH = "stringLength";
    private static final String TIME_BASE = "timeBase";
    private static final String TIME_BASE_MULTIPLIER = "multiplier";
    
    
//    private final String address;
    private final S7SubscriptionType subscriptionType;
    private final EventType eventtype;
    private final S7Tag[] s7tags;
    private final ArrayList<Integer> ackAlarms;
    private final AlarmType alarmQueryType;
    private final TimeBase timebase;
    private final short multiplier;    

    public S7SubscriptionTag(S7SubscriptionType subscriptionType, EventType eventtype) {
//        this.address = address;
        this.subscriptionType = subscriptionType;
        this.eventtype = eventtype;
        this.s7tags = null;
        this.ackAlarms = null;
        this.alarmQueryType = null;
        this.timebase = null;
        this.multiplier = 0;
    }

    public S7SubscriptionTag(S7SubscriptionType subscriptionType, ArrayList<Integer> ackAlarms) {
//        this.address = address;
        this.subscriptionType = subscriptionType;
        this.eventtype = null;
        this.s7tags = null;
        this.ackAlarms = ackAlarms;
        this.alarmQueryType = null;
        this.timebase = null;
        this.multiplier = 0;        
    }

    public S7SubscriptionTag(S7SubscriptionType subscriptionType, AlarmType alarmQueryType) {
//        this.address = address;
        this.subscriptionType = subscriptionType;
        this.eventtype = null;
        this.s7tags = null;
        this.ackAlarms = null;
        this.alarmQueryType = alarmQueryType;
        this.timebase = null;
        this.multiplier = 0;        
    }

//    public S7SubscriptionTag(S7SubscriptionType subscriptionType, S7Tag s7Tag) {
////        this.address = address;
//        this.subscriptionType = subscriptionType;
//        this.eventtype = null;
//        this.s7tags = s7Tag;
//        this.ackAlarms = null;
//        this.alarmQueryType = null;
//        this.timebase = null;
//        this.multiplier = 0;        
//    }
    
    public S7SubscriptionTag(S7SubscriptionType subscriptionType, S7Tag[] s7tags, TimeBase timebase, short multiplier) {
//        this.address = address;
        this.subscriptionType = subscriptionType;
        this.eventtype  = null;
        this.s7tags   = s7tags;
        this.ackAlarms   = null;  
        this.alarmQueryType = null;
        this.timebase   = timebase;
        this.multiplier = multiplier;
    }      

    @Override
    public String getAddressString() {
        return null;
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.RAW_BYTE_ARRAY;
    }

    @Override
    public List<ArrayInfo> getArrayInfo() {
        return Collections.emptyList();
    }

    public S7SubscriptionType getTagType() {
        return subscriptionType;
    }

    public EventType getEventType() {
        return eventtype;
    }

    public S7Tag[] getS7Tags() {
        return s7tags;
    }

    public ArrayList<Integer> getAckAlarms() {
        return ackAlarms;
    }

    public AlarmType getAlarmQueryType() {
        return alarmQueryType;
    }
    
    public TimeBase getTimeBase() {
        return timebase;
    }

    public short getMultiplier() {
        return multiplier;
    }        

    public static boolean matches(String tagString) {
        
        return EVENT_SUBSCRIPTION_TYPE_PATTERN.matcher(tagString).matches() ||
            EVENT_ALARM_ACK_PATTERN.matcher(tagString).matches() ||
            EVENT_ALARM_QUERY_PATTERN.matcher(tagString).matches() ||
            EVENT_SUBSCRIPTION_S7ANY_QUERY_PATTERN.matcher(tagString).matches() ||
            EVENT_SUBSCRIPTION_DB_QUERY_PATTERN.matcher(tagString).matches() ||
            EVENT_CANCEL_JOB_QUERY_PATTERN.matcher(tagString).matches();
        
//        return EVENT_SUBSCRIPTION_TYPE_PATTERN.matcher(tagString).matches() ||
//            EVENT_ALARM_ACK_PATTERN.matcher(tagString).matches() ||
//            EVENT_ALARM_QUERY_PATTERN.matcher(tagString).matches() ||
//            S7Tag.matches(tagString);
    }

    public static S7SubscriptionTag of(String tagString) {
        {
            Matcher matcher = EVENT_SUBSCRIPTION_TYPE_PATTERN.matcher(tagString);
            if (matcher.matches()) {
                return new S7SubscriptionTag(S7SubscriptionType.EVENT_SUBSCRIPTION,
                    EventType.valueOf(tagString));
            }
        }

        {
            //TODO: Actually only ALARM_S (SIG_1)
            Matcher matcher = EVENT_ALARM_ACK_PATTERN.matcher(tagString);
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
                return new S7SubscriptionTag(S7SubscriptionType.ALARM_ACK,
                            arrEventId);
                
            }            
        }

        {
            //TODO: Support for ALARM_8            
            Matcher matcher = EVENT_ALARM_QUERY_PATTERN.matcher(tagString); 
            TimeBase tb = null;
            short multi = 0;            
            if (matcher.matches()){

                return new S7SubscriptionTag(S7SubscriptionType.ALARM_QUERY,
                            AlarmType.ALARM_S);    
            }
        }

        {
            Matcher matcher = EVENT_SUBSCRIPTION_DB_QUERY_PATTERN.matcher(tagString);          
            if (matcher.matches()){
                TimeBase tb = TimeBase.valueOf(matcher.group(TIME_BASE));
                short multi = Short.parseShort(matcher.group(TIME_BASE_MULTIPLIER));                
                String strAddress = matcher.group(9);
                strAddress = strAddress.replaceAll("\\[", ".0:BYTE[");
                String[] dbAddress = strAddress.split(",");       
                S7Tag[] s7tags = new S7Tag[dbAddress.length];
                int i=0;
                for (String address:dbAddress) {
                    s7tags[i] = S7Tag.of(address);
                    i++;
                }                
                return new S7SubscriptionTag(S7SubscriptionType.CYCLIC_DB_SUBSCRIPTION,
                            s7tags,
                            tb,
                            multi);    
            }            
            
        }   
        
        {
            Matcher matcher = EVENT_SUBSCRIPTION_S7ANY_QUERY_PATTERN.matcher(tagString);
            if (matcher.matches()){     
                TimeBase tb = TimeBase.valueOf(matcher.group(TIME_BASE));
                short multi = Short.parseShort(matcher.group(TIME_BASE_MULTIPLIER));
                S7Tag[] myTags = null;
                String strAddress = matcher.group(9);  
                String[] fieldAddress = strAddress.split(","); 
                myTags = new S7Tag[fieldAddress.length];            
                int i=0;
                for (String address:fieldAddress) {
                    myTags[i] = S7Tag.of(address);
                    i++;
                }         
                return new S7SubscriptionTag(S7SubscriptionType.CYCLIC_SUBSCRIPTION,
                            myTags,
                tb,
                multi);                
            }
        }

         {
            Matcher matcher = EVENT_CANCEL_JOB_QUERY_PATTERN.matcher(tagString);
            if (matcher.matches()){
                String[] arrIdAndSig;
                String strJobIds = matcher.group(2);
                String[] arrStrEventId = strJobIds.split(",");
                ArrayList<Integer> arrJobId = new ArrayList<>();
                for (String jobId:arrStrEventId){
                    arrJobId.add(Integer.parseInt(jobId));                   
                }
                return new S7SubscriptionTag(S7SubscriptionType.CYCLIC_UNSUBSCRIPTION,
                            arrJobId);
                
            }            
        } 
        
        
        
        throw new PlcInvalidTagException("Unable to parse address: " + tagString);
    }

}

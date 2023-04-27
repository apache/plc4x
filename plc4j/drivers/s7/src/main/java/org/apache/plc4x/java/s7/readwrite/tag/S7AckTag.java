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
package org.apache.plc4x.java.s7.readwrite.tag;

import java.util.ArrayList;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.s7.readwrite.types.S7SubscriptionType;

public class S7AckTag implements PlcTag {
    
    private static final Pattern EVENT_ALARM_ACK_PATTERN = 
        Pattern.compile("(^ACK:)(((?:,{0,1})(16#[0-9a-fA-F]{8})(;([0-9a-fA-F]{2})))+)");
    
    private final ArrayList<Integer> ackAlarmIds;
    private final ArrayList<Integer> ackAlarmSigs;    
    
    public S7AckTag(ArrayList<Integer> ackAlarmIds, ArrayList<Integer> ackAlarmSigs){
        this.ackAlarmIds = ackAlarmIds;
        this.ackAlarmSigs = ackAlarmSigs;
    };    
    
    
    public ArrayList<Integer> getAlarmIds() {
        return ackAlarmIds;
    }

    public ArrayList<Integer> getAlarmSigs() {
        return ackAlarmSigs;
    }    
    
    public static boolean matches(String fieldString) {
        return EVENT_ALARM_ACK_PATTERN.matcher(fieldString).matches();
    }      
    
    public static S7AckTag of(String fieldString) {
        //TODO: Actually only ALARM_S (SIG_1)
        Matcher matcher = EVENT_ALARM_ACK_PATTERN.matcher(fieldString);
        if (matcher.matches()){
            String[] arrIdAndSig;
            String EventIds = matcher.group(2);
            String[] arrStrEventId = EventIds.split(",");
            ArrayList<Integer> arrAlarmIds = new ArrayList<>();
            ArrayList<Integer> arrAlarmSigs = new ArrayList<>();                
            for (String EventId:arrStrEventId){
                EventId = EventId.replaceAll("16#", "");
                arrIdAndSig =  EventId.split(";");
                arrAlarmIds.add(Integer.parseInt(arrIdAndSig[0], 16));
                arrAlarmSigs.add(Integer.parseInt(arrIdAndSig[1], 16));                    
            }
            return new S7AckTag(arrAlarmIds, arrAlarmSigs);                
        } else return null;                
    };    

    @Override
    public String getAddressString() {
        throw new UnsupportedOperationException("Not supported yet.");
    }
    
    
}

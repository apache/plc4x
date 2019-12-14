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

package org.apache.plc4x.java.s7.netty.model.types;

import java.util.HashMap;
import java.util.Map;
import org.apache.plc4x.java.api.model.PlcField;

/**
 *
 * @author cgarcia
 */
public enum PushEventType implements PlcField {
    
    ALARM_SQ("ALARM_SQ"),
    ALARM_S("ALARM_S"),
    ALARM_SC("ALARM_SC"),
    ALARM_DQ("ALARM_DQ"),
    ALARM_D("ALARM_D"),
    NOTIFY_8P("NOTIFY_8P"),
    ALARM("ALARM"),
    ALARM_8("ALARM_8"),
    ALARM_8P("ALARM_8P"),
    NOTIFY("NOTIFY");
    
    private static final Map<String, PushEventType> map;
    
    static {
        map = new HashMap<>();
        for (PushEventType s7event : PushEventType.values()) {
            map.put(s7event.event, s7event);
        }
    }
    
    private final String event;

    PushEventType(String event){
        this.event = event;
    }
    
    public String getEvent(){
        return event;
    }
    
    public static PushEventType valueOfEvent(String event) {
         return map.get(event);
    }
    
    
}

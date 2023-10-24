/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package org.apache.plc4x.merlot.api.impl;

import java.util.Hashtable;
import java.util.Map;
import java.util.UUID;
import org.osgi.service.dal.FunctionEvent;
import static org.osgi.service.dal.FunctionEvent.EVENT_CLASS;
import static org.osgi.service.dal.FunctionEvent.EVENT_PACKAGE;
import static org.osgi.service.dal.FunctionEvent.FUNCTION_UID;
import static org.osgi.service.dal.FunctionEvent.PROPERTY_NAME;
import static org.osgi.service.dal.FunctionEvent.PROPERTY_VALUE;


public class PlcFunctionEvent extends FunctionEvent {    
    
    public PlcFunctionEvent(PlcFunctionEventBuilder builder) {
        super(builder.topic, builder.metadata);
    }
      
    public static class PlcFunctionEventBuilder {
        private final Map<String, Object > metadata; 
        private String topic = null;

        public PlcFunctionEventBuilder() {
            metadata = new Hashtable<>();
        }      
        
        public PlcFunctionEventBuilder setTopic(String topic) {
            this.topic = topic;
            return this;
        }
        
        public PlcFunctionEventBuilder putEventClass(String evclass) {
            metadata.put(EVENT_CLASS, evclass);
            return this;
        }    
        
        public PlcFunctionEventBuilder putEventPackage(String evpackage) {
            metadata.put(EVENT_PACKAGE, evpackage);
            return this;
        }  
        
        public PlcFunctionEventBuilder putFunctionUid(UUID uid) {
            metadata.put(FUNCTION_UID, uid);
            return this;
        }   
        
        public PlcFunctionEventBuilder putPropertyName(String name) {
            metadata.put(PROPERTY_NAME, name);
            return this;
        }   

        public PlcFunctionEventBuilder putPropertyValue(Object o) {
            metadata.put(PROPERTY_VALUE, o);
            return this;
        } 
        
        public PlcFunctionEventBuilder putTopicPropertyChanged(Object o) {
            metadata.put(TOPIC_PROPERTY_CHANGED, o);
            return this;
        }         
        
        public PlcFunctionEventBuilder putMetaData(String k, Object o) {
            this.metadata.put(k, o);
            return this;
        } 
        
        public PlcFunctionEvent build(){
            PlcFunctionEvent femeta = new PlcFunctionEvent(this);
            validatePlcGroupObject(femeta);
            return femeta;
        }         
        
        private void validatePlcGroupObject(PlcFunctionEvent femeta) {
            //
        }            
        
    }    
    
}

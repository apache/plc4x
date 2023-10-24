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
import org.osgi.service.dal.FunctionData;
import static org.osgi.service.dal.FunctionData.DESCRIPTION;
import static org.osgi.service.dal.PropertyMetadata.DESCRIPTION;


public class PlcFunctionData extends FunctionData {

    public PlcFunctionData(PlcFunctionDataBuilder builder) {
        super(builder.metadata);
    }

    
    public static class PlcFunctionDataBuilder {
        private final Map<String, Object > metadata;
        private long timestamp = Long.MIN_VALUE;

        public PlcFunctionDataBuilder() {
            metadata = new Hashtable<>();
        }  
        
        public PlcFunctionDataBuilder putDescription(String desc) {
            this.metadata.put(DESCRIPTION, desc);
            return this;
        } 
        
        public PlcFunctionDataBuilder putTimeStamp(long timestamp) {
            this.timestamp = timestamp;
            return this;
        } 
        
        public PlcFunctionDataBuilder putMetaData(String k, Object o) {
            this.metadata.put(k, o);
            return this;
        }        
        
        public PlcFunctionData build(){
            PlcFunctionData fdmeta = new PlcFunctionData(this);
            validatePlcGroupObject(fdmeta);
            return fdmeta;
        }   
        
        private void validatePlcGroupObject(PlcFunctionData fdmeta) {
            //
        }           
        
    }
    
}

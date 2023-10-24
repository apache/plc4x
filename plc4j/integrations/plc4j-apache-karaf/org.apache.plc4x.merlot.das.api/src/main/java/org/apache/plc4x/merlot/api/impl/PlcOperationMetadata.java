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
package org.apache.plc4x.merlot.api.impl;

import java.util.ArrayList;
import java.util.Hashtable;
import java.util.List;
import java.util.Map;
import org.apache.plc4x.java.api.value.PlcValue;
import org.osgi.service.dal.OperationMetadata;
import static org.osgi.service.dal.OperationMetadata.DESCRIPTION;
import org.osgi.service.dal.PropertyMetadata;

/*
*
*/
public class PlcOperationMetadata implements OperationMetadata{

    private final Map<String, Object > metadata; 
    
    private PlcOperationMetadata valuemeta;
    private final PlcPropertyMetadata[] opmetas;
    
    public PlcOperationMetadata(PlcOperationMetaDataBuilder builder){
        this.metadata = builder.metadata;
        this.opmetas = builder.opmetas.toArray(new PlcPropertyMetadata[0]);
    }
        
    @Override
    public Map<String, ?> getMetadata() {
        return metadata;
    }

    @Override
    public PropertyMetadata getReturnValueMetadata() {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public PropertyMetadata[] getParametersMetadata() {
        return opmetas;
    }
    
    public static class PlcOperationMetaDataBuilder {
        
        private final Map<String, Object > metadata; 
        private final List<PlcPropertyMetadata> opmetas = new ArrayList<>();
        private PlcPropertyMetadata valuemeta;
    
        public PlcOperationMetaDataBuilder() {
            metadata = new Hashtable<>();            
        }
        
        public PlcOperationMetaDataBuilder putDescription(String desc) {
            this.metadata.put(DESCRIPTION, desc);
            return this;
        }  
        
        public PlcOperationMetaDataBuilder setReturnValueMetadata(PlcPropertyMetadata opmeta) {
            this.valuemeta = opmeta;
            return this;
        }        
        
        public PlcOperationMetaDataBuilder addOperationMetadata(PlcPropertyMetadata opmeta) {
            this.opmetas.add(opmeta);
            return this;
        }           
        
        public PlcOperationMetaDataBuilder putMetaData(String k, Object o) {
            this.metadata.put(k, o);
            return this;
        }
        
        public PlcOperationMetadata build(){
            PlcOperationMetadata opmeta = new PlcOperationMetadata(this);
            validatePlcGroupObject(opmeta);
            return opmeta;
        }
        
        private void validatePlcGroupObject(PlcOperationMetadata opmeta) {
            //
        }          
    }    
    
}

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

import java.util.Hashtable;
import java.util.Map;
import org.osgi.service.dal.FunctionData;
import org.osgi.service.dal.OperationMetadata;
import static org.osgi.service.dal.OperationMetadata.DESCRIPTION;
import org.osgi.service.dal.PropertyMetadata;
import static org.osgi.service.dal.PropertyMetadata.DESCRIPTION;

/*
*
*/
public class PlcPropertyMetadata implements PropertyMetadata{
    
    /**
     * Metadata key, which value represents the property long description. The
     * property value type is {@code java.lang.String}.
     */
    public static final String	DESC_LONG			= "long_description";    

    private final Map<String, Object > metadata;    
    
    public PlcPropertyMetadata(PlcPropertyMetaDataBuilder builder){
        this.metadata = builder.metadata;
    }
        
    @Override
    public Map<String, ?> getMetadata(String unit) {
        return this.metadata;
    }

    //TODO: Here you must consult the UNS/S88 service, which deploys the management of engineering units.
    @Override
    public FunctionData getStep(String unit) {
        return null;
    }

    //TODO: Here you must consult the UNS/S88 service, which deploys the management of engineering units.    
    @Override
    public FunctionData[] getEnumValues(String unit) {
        return null; 
    }

    //TODO: Here you must consult the UNS/S88 service, which deploys the management of engineering units.    
    @Override
    public FunctionData getMinValue(String unit) {
        return null;
    }

    //TODO: Here you must consult the UNS/S88 service, which deploys the management of engineering units.    
    @Override
    public FunctionData getMaxValue(String unit) {
        return null;
    }
    
    public static class PlcPropertyMetaDataBuilder {
        
        private final Map<String, Object > metadata;   
    
        public PlcPropertyMetaDataBuilder() {
            metadata = new Hashtable<>();
        }
        
        public PlcPropertyMetaDataBuilder putDescription(String desc) {
            this.metadata.put(DESCRIPTION, desc);
            return this;
        }  
        
        public PlcPropertyMetaDataBuilder putLongDescription(String desc) {
            this.metadata.put(DESC_LONG, desc);
            return this;
        }          
        
        public PlcPropertyMetaDataBuilder putAccess(int access) {
            this.metadata.put(ACCESS, access);
            return this;
        } 
        
        public PlcPropertyMetaDataBuilder putUnits(String... units) {
            this.metadata.put(UNITS, units);
            return this;
        }    
        
        public PlcPropertyMetaDataBuilder putMetaData(String k, Object o) {
            this.metadata.put(k, o);
            return this;
        }
        
        public PlcPropertyMetadata build(){
            PlcPropertyMetadata prometa = new PlcPropertyMetadata(this);
            validatePlcGroupObject(prometa);
            return prometa;
        }   
        
        private void validatePlcGroupObject(PlcPropertyMetadata prometa) {
            //
        }           
    }    
    
}

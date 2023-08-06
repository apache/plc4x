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
package org.apache.plc4x.merlot.modbus.sim.impl;

import org.apache.plc4x.merlot.modbus.dev.api.ModbusDevice;
import org.apache.plc4x.merlot.scheduler.api.JobContext;
import java.util.Random;


public class ModbusSimRandomImpl extends ModbusSimImpl {
    
    private boolean randombits = false;
    
    public ModbusSimRandomImpl(ModbusDevice mbdev) {
        super(mbdev);
        this.random = new Random();
    }

    @Override
    public void init() {
        super.init();
    }

    @Override
    public void execute(JobContext context) {
        Double value = 0.0;
        if (!started) return;
        if (tagF != null) {
            if (tagX != null) {
                x = getValue(mbdev,tagX,0);
            }               
            for (int i=0; i<tagF.range; i++) {
                switch(tagF.stType){
                    case pvBoolean: 
                        value = (random.nextBoolean()?1.0:0.0) * x ;
                        putValue(value, mbdev,tagF,i);
                        f = value;
                        break;
                    default: {                       
                        value = random.nextDouble() * x;                      
                        putValue(value, mbdev,tagF,i);
                        f = value;
                    }
                }
            }
            
        }        
    }
 
    
    
    
}

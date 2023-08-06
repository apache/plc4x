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
import org.apache.plc4x.merlot.modbus.sim.core.ModbusSimSignalType;
import org.apache.plc4x.merlot.scheduler.api.JobContext;


public class ModbusSimSignalImpl extends ModbusSimImpl {
    
    private Double t = 0.0;
    private Double fase = 0.0;
    private Double paso  = 0.0;    
    private ModbusSimSignalType signalType;
    
    public ModbusSimSignalImpl(ModbusDevice mbdev) {
        super(mbdev);
    }

    @Override
    public void setFunction(String strFunction) {
        this.strFunction = strFunction;
        signalType = ModbusSimSignalType.valueOfEvent(strFunction);
    }
    
    @Override
    public void execute(JobContext context) {
        
        Double result = 0.0;
        double frecuencia = 0.5; 
        double fase = 0.0;
        double magnitud = 1.0;

        if (!started) return;
        
        
        if (tagF != null) {
            f = getValue(mbdev,tagF,0);
        } else {
            return;
        }
        
        if (tagX != null) {
            frecuencia = getValue(mbdev,tagX,0);
            x = frecuencia;
        }
        if (tagY != null) {
            fase = getValue(mbdev,tagY,0);
            y = fase;
        }
        
        if (tagZ != null) {
            magnitud = getValue(mbdev,tagZ,0);
            z = magnitud;
        } 
        
        if (strFunction != null) {
            result = getSignalValue(signalType, paso, fase) * magnitud;
        }
        
        if (tagF != null) {
            putValue(result, mbdev,tagF,0);
            f = result;          
        }
        
        paso = paso + 2*Math.PI*frecuencia/10;         
        
    }
                   
    private Double getSignalValue(ModbusSimSignalType signal, Double t, Double fase){
        Double value = 0.0;
        switch(signal){
            case SQUARE:
                    value = Math.signum(Math.sin(t + fase));
                break;
            case SINE:
                    value = Math.sin(t + fase);
                break;
            case TRIANGULAR:
                    value = 2 * Math.abs(2 * ((t + fase) / (2 * Math.PI) - Math.floor((t + fase) / (2 * Math.PI) + 0.5))) - 1;
                break;
            case SAWTOOTH:
                    value =  2 * ((t + fase) / (2 * Math.PI) - Math.floor(0.5 + (t + fase) / (2 * Math.PI)));
                break;
        }
        return value;
    }
    
}

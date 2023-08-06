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

import org.apache.plc4x.merlot.modbus.sim.api.ModbusSim;
import org.apache.plc4x.merlot.modbus.sim.api.ModbusSimMBean;
import javax.management.NotCompliantMBeanException;
import javax.management.StandardMBean;


public class ModbusSimMBeanImpl extends StandardMBean implements ModbusSimMBean {

    private final ModbusSim ms; 
    
    public ModbusSimMBeanImpl(ModbusSim ms) throws NotCompliantMBeanException {
        super(ModbusSimMBean.class);
        this.ms = ms;
    }
    
    @Override
    public void init() {
        
    }

    @Override
    public void destroy() {
        
    }

    @Override
    public void start() {
        ms.start();
    }

    @Override
    public void stop() {
        ms.stop();
    }

    @Override
    public void setFunction(String strFunction) {
        ms.setFunction(strFunction);
    }

    @Override
    public String getFunction() {
        return ms.getFunction();
    }

    @Override
    public double getF() {
        return ms.getF();
    }

    @Override
    public void setF(double f) {
        ms.setF(f);
    }
    
    @Override
    public double getX() {
        return ms.getX();
    }

    @Override
    public void setX(double x) {
        ms.setX(x);
    }    
    
    @Override
    public double getY() {
        return ms.getY();
    }

    @Override
    public void setY(double y) {
        ms.setY(y);
    }    
    
    @Override
    public double getZ() {
        return ms.getZ();
    }
    
    @Override
    public void setZ(double z) {
        ms.setZ(z);
    }    
    
}

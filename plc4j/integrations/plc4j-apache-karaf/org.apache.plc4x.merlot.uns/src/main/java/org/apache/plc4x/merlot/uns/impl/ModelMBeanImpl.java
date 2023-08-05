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
package org.apache.plc4x.merlot.uns.impl;

import org.apache.plc4x.merlot.uns.api.Model;
import org.apache.plc4x.merlot.uns.api.ModelMBean;


public class ModelMBeanImpl implements ModelMBean {

    private final Model model;

    public ModelMBeanImpl(Model model) {
        this.model = model;
    }
                
    @Override
    public String getEnterprise() {
        return null;
    }

    @Override
    public String getPlant() {
        return null;
    }

    @Override
    public String getArea() {
        return null;
    }

    @Override
    public String getCell() {
        return null;
    }

    @Override
    public String getUnit() {
        return null;
    }

    @Override
    public String getDomain(Long id) {
        return null;
    }

    @Override
    public String getDomain(String pv) {
        return null;
    }

    @Override
    public String getStorageGroup() {
        return null;
    }
    
}

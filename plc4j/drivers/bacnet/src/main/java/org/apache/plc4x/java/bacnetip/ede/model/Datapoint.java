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
package org.apache.plc4x.java.bacnetip.ede.model;

import org.apache.plc4x.java.api.value.*;
import org.apache.plc4x.java.bacnetip.field.BacNetIpField;
import org.apache.plc4x.java.spi.values.PlcBOOL;
import org.apache.plc4x.java.spi.values.PlcDINT;
import org.apache.plc4x.java.spi.values.PlcLREAL;
import org.apache.plc4x.java.spi.values.PlcSTRING;

import java.util.HashMap;
import java.util.Map;

public class Datapoint {

    private final BacNetIpField address;
    private final String keyName;
    private final String objectName;
    private final String description;
    private final Double defaultValue;
    private final Double minValue;
    private final Double maxValue;
    private final Boolean commandable;
    private final Boolean supportsCov;
    private final Double hiLimit;
    private final Double lowLimit;
    private final String stateTextReference;
    private final Integer unitCode;
    private final Integer vendorSpecificAddress;
    private final Integer notificationClass;

    public Datapoint(BacNetIpField address, String keyName, String objectName, String description, Double defaultValue, Double minValue, Double maxValue, Boolean commandable, Boolean supportsCov, Double hiLimit, Double lowLimit, String stateTextReference, Integer unitCode, Integer vendorSpecificAddress, Integer notificationClass) {
        this.address = address;
        this.keyName = keyName;
        this.objectName = objectName;
        this.description = description;
        this.defaultValue = defaultValue;
        this.minValue = minValue;
        this.maxValue = maxValue;
        this.commandable = commandable;
        this.supportsCov = supportsCov;
        this.hiLimit = hiLimit;
        this.lowLimit = lowLimit;
        this.stateTextReference = stateTextReference;
        this.unitCode = unitCode;
        this.vendorSpecificAddress = vendorSpecificAddress;
        this.notificationClass = notificationClass;
    }

    public BacNetIpField getAddress() {
        return address;
    }

    public String getKeyName() {
        return keyName;
    }

    public String getObjectName() {
        return objectName;
    }

    public String getDescription() {
        return description;
    }

    public Double getDefaultValue() {
        return defaultValue;
    }

    public Double getMinValue() {
        return minValue;
    }

    public Double getMaxValue() {
        return maxValue;
    }

    public Boolean getCommandable() {
        return commandable;
    }

    public Boolean getSupportsCov() {
        return supportsCov;
    }

    public Double getHiLimit() {
        return hiLimit;
    }

    public Double getLowLimit() {
        return lowLimit;
    }

    public String getStateTextReference() {
        return stateTextReference;
    }

    public Integer getUnitCode() {
        return unitCode;
    }

    public Integer getVendorSpecificAddress() {
        return vendorSpecificAddress;
    }

    public Integer getNotificationClass() {
        return notificationClass;
    }

    public Map<String, PlcValue> toPlcValues() {
        Map<String, PlcValue> values = new HashMap<>();
        values.put("keyName", (keyName == null) ? null : new PlcSTRING(keyName));
        values.put("objectName", (objectName == null) ? null : new PlcSTRING(objectName));
        values.put("description", (description == null) ? null : new PlcSTRING(description));
        values.put("defaultValue", (defaultValue == null) ? null : new PlcLREAL(defaultValue));
        values.put("minValue", (minValue == null) ? null : new PlcLREAL(minValue));
        values.put("maxValue", (maxValue == null) ? null : new PlcLREAL(maxValue));
        values.put("commandable", (commandable == null) ? null : new PlcBOOL(commandable));
        values.put("supportsCov", (supportsCov == null) ? null : new PlcBOOL(supportsCov));
        values.put("hiLimit", (hiLimit == null) ? null : new PlcLREAL(hiLimit));
        values.put("lowLimit", (lowLimit == null) ? null : new PlcLREAL(lowLimit));
        values.put("stateTextReference", (stateTextReference == null) ? null : new PlcSTRING(stateTextReference));
        values.put("unitCode", (unitCode == null) ? null : new PlcDINT(unitCode));
        values.put("vendorSpecificAddress", (vendorSpecificAddress == null) ? null : new PlcDINT(vendorSpecificAddress));
        values.put("notificationClass", (notificationClass == null) ? null : new PlcDINT(notificationClass));
        return values;
    }

}

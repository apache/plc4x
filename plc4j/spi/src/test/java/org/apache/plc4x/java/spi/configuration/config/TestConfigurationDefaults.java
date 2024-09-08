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

package org.apache.plc4x.java.spi.configuration.config;

import org.apache.plc4x.java.spi.configuration.PlcConnectionConfiguration;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.*;

public class TestConfigurationDefaults implements PlcConnectionConfiguration {

    @ConfigurationParameter
    @BooleanDefaultValue(true)
    private boolean booleanField;

    @ConfigurationParameter
    @IntDefaultValue(42)
    private int integerField;

    @ConfigurationParameter
    @LongDefaultValue(232323232323232323L)
    private long longField;

    @ConfigurationParameter
    @FloatDefaultValue(3.1415927f)
    private float floatField;

    @ConfigurationParameter
    @DoubleDefaultValue(2.718281828459045d)
    private double doubleField;

    @ConfigurationParameter
    @StringDefaultValue("Hurz")
    private String stringField;

    public boolean isBooleanField() {
        return booleanField;
    }

    public int getIntegerField() {
        return integerField;
    }

    public long getLongField() {
        return longField;
    }

    public float getFloatField() {
        return floatField;
    }

    public double getDoubleField() {
        return doubleField;
    }

    public String getStringField() {
        return stringField;
    }

}

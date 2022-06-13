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
package org.apache.plc4x.java.bacnetip.ede.layouts;

public class EdeVersion2Layout implements EdeLayout {

    @Override
    public int getKeyNamePos() {
        return 0;
    }

    @Override
    public int getDeviceInstancePos() {
        return 1;
    }

    @Override
    public int getObjectNamePos() {
        return 2;
    }

    @Override
    public int getObjectTypePos() {
        return 3;
    }

    @Override
    public int getObjectInstancePos() {
        return 4;
    }

    @Override
    public int getDescriptionPos() {
        return 5;
    }

    @Override
    public int getDefaultValuePos() {
        return 6;
    }

    @Override
    public int getMinValuePos() {
        return 7;
    }

    @Override
    public int getMaxValuePos() {
        return 8;
    }

    @Override
    public int getCommandablePos() {
        return 9;
    }

    @Override
    public int getSupportsCovPos() {
        return 10;
    }

    @Override
    public int getHiLimitPos() {
        return 11;
    }

    @Override
    public int getLowLimitPos() {
        return 12;
    }

    @Override
    public int getStateTextReferencePos() {
        return 13;
    }

    @Override
    public int getUnitCodePos() {
        return 14;
    }

    @Override
    public int getVendorSpecificAddressPos() {
        return 15;
    }



}

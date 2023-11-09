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

package org.apache.plc4x.java.transport.serial;

import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.IntDefaultValue;

public class DefaultSerialTransportConfiguration implements SerialTransportConfiguration {

    @ConfigurationParameter("baud-rate")
    @IntDefaultValue(57600)
    private int baudRate;

    @ConfigurationParameter("num-data-bits")
    @IntDefaultValue(8)
    private int numDataBits;

    @ConfigurationParameter("num-stop-bits")
    @IntDefaultValue(1)
    private int numStopBits;

    @ConfigurationParameter("parity-bits")
    @IntDefaultValue(0)
    private int parityBits;

    @Override
    public int getBaudRate() {
        return baudRate;
    }

    public void setBaudRate(int baudRate) {
        this.baudRate = baudRate;
    }

    @Override
    public int getNumDataBits() {
        return numDataBits;
    }

    public void setNumDataBits(int numDataBits) {
        this.numDataBits = numDataBits;
    }

    @Override
    public int getNumStopBits() {
        return numStopBits;
    }

    public void setNumStopBits(int numStopBits) {
        this.numStopBits = numStopBits;
    }

    @Override
    public int getParityBits() {
        return parityBits;
    }

    public void setParityBits(int parityBits) {
        this.parityBits = parityBits;
    }

}

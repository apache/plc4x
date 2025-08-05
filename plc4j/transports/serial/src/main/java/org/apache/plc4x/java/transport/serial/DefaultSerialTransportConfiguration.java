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

package org.apache.plc4x.java.transport.serial;

import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.Description;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.IntDefaultValue;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.StringDefaultValue;

public class DefaultSerialTransportConfiguration implements SerialTransportConfiguration {

    @ConfigurationParameter("baud-rate")
    @IntDefaultValue(57600)
    @Description("Baud-rate the serial port is set to.\n" +
        "Typical values are:\n" +
        "- 9600\n" +
        "- 14400\n" +
        "- 19200\n" +
        "- 38400\n" +
        "- 57600\n" +
        "- 115200\n" +
        "- 128000\n" +
        "But can also be smaller or larger values."
    )
    private int baudRate;

    @ConfigurationParameter("num-data-bits")
    @IntDefaultValue(8)
    @Description("Number of data-bits used to send data.\n" +
        "Typical values are:\n" +
        "- 7\n" +
        "- 8\n" +
        "But can also be smaller or larger values."
    )
    private int numDataBits;

    @ConfigurationParameter("num-stop-bits")
    @IntDefaultValue(1)
    @Description("Number of stop-bits used to terminate data.\n" +
        "Typical values are:\n" +
        "- 1\n" +
        "- 2\n" +
        "(The theoretical 1.5 stop-bits setting is not supported)"
    )
    private int numStopBits;

    @ConfigurationParameter("parity")
    @StringDefaultValue("NO_PARITY")
    @Description("Number of bits used to calculate data parity.\n" +
        "This is used to detect errors in transmission.\n" +
        "Allowed values are:\n" +
        "- NO_PARITY\n" +
        "- ODD_PARITY\n" +
        "- EVEN_PARITY\n" +
        "- MARK_PARITY\n" +
        "- SPACE_PARITY")
    private SerialChannelOptionsParity parity;

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
    public SerialChannelOptionsParity getParity() {
        return parity;
    }

    public void setParity(SerialChannelOptionsParity parity) {
        this.parity = parity;
    }

}

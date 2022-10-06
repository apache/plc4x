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
package org.apache.plc4x.java.examples.pollloop;

import java.util.ArrayList;
import java.util.List;

public class Main {

    public static void main(String[] args) throws Exception {

        List<String> variables = new ArrayList<>();
        variables.add("%M89:REAL"); // currentSpeedInRpm	DriveVariables	Real	%MD89
        variables.add("%Q20:REAL"); // temperatureInCelsius	SensorVariables	Real	%QD20
        variables.add("%I58.0:BOOL"); // switchingStateOfCapSensor SensorVariable Bool %I58.0
        variables.add("%Q25:WORD"); // distanceInMm	SensorVariables	Word	%QW25
        variables.add("%Q82:REAL"); // currentDrivePercent	DriveVariables	Real	%QD82
        variables.add("%M86:INT");  // driveSetFreqInPercent	DriveVariables	Int	%MW86
        variables.add("%I66:WORD"); // rawValueOfTempSensor	SensorVariables	Word	%IW66
        variables.add("%I58:WORD"); // capSensorWord SensorVariables	Word	%IW58
        variables.add("%M0:BYTE"); // capSensorWord SensorVariables	Word	%IW58

        PollLoop pollLoop = new PollLoop( "s7://192.168.100.49/0/1",PollLoop.PLC4JTYPE_SIEMENS,
            variables, 1000);
        pollLoop.start();
    }

}

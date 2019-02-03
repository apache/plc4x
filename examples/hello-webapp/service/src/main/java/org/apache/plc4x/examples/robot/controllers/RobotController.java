/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */

package org.apache.plc4x.examples.robot.controllers;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import javax.annotation.PostConstruct;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

@RestController
@RequestMapping("api/robot")
public class RobotController {

    private static final Logger logger = LoggerFactory.getLogger(RobotController.class);

    private static final byte MOTOR_RIGHT_BACKWARDS = 0x01;
    private static final byte MOTOR_RIGHT_ON = 0x02;
    private static final byte MOTOR_LEFT_BACKWARDS = 0x04;
    private static final byte MOTOR_LEFT_ON = 0x08;

    @Value("${plc4x.connection-string}")
    private String connectionString;

    @Value("${plc4x.address-string}")
    private String addressString;

    private PlcConnection connection;

    @PostConstruct
    public void init() throws PlcConnectionException {
        connection = new PlcDriverManager().getConnection(connectionString);
    }

    @RequestMapping("move")
    public boolean move(@RequestParam(value="direction", defaultValue="stop") String direction) {
        logger.info("Move in direction: " + direction);
        byte state;
        switch (direction) {
            case "forward-right":
                state = MOTOR_LEFT_ON;
                break;
            case "forward":
                state = (byte) (MOTOR_LEFT_ON | MOTOR_RIGHT_ON);
                break;
            case "forward-left":
                state = MOTOR_RIGHT_ON;
                break;
            case "left":
                state = (byte) (MOTOR_LEFT_BACKWARDS | MOTOR_LEFT_ON | MOTOR_RIGHT_ON);
                break;
            case "right":
                state = (byte) (MOTOR_LEFT_ON | MOTOR_RIGHT_BACKWARDS | MOTOR_RIGHT_ON);
                break;
            case "backward-right":
                state = (byte) (MOTOR_LEFT_BACKWARDS | MOTOR_LEFT_ON);
                break;
            case "backward":
                state = (byte) (MOTOR_LEFT_BACKWARDS | MOTOR_LEFT_ON | MOTOR_RIGHT_BACKWARDS | MOTOR_RIGHT_ON);
                break;
            case "backward-left":
                state = (byte) (MOTOR_RIGHT_BACKWARDS | MOTOR_RIGHT_ON);
                break;
            default:
                state = 0;
                break;
        }
        try {
            PlcWriteRequest updateRequest = connection.writeRequestBuilder().addItem("state", addressString, state).build();
            updateRequest.execute().get(2000, TimeUnit.MILLISECONDS);
            return true;
        } catch (InterruptedException | ExecutionException | TimeoutException e) {
            logger.error("Caught Exception:", e);
            return false;
        }
    }

}

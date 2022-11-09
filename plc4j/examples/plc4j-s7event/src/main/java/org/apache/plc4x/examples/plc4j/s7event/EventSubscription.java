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
package org.apache.plc4x.examples.plc4j.s7event;

import java.util.Map;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.s7.events.S7AlarmEvent;
import org.apache.plc4x.java.s7.events.S7ModeEvent;
import org.apache.plc4x.java.s7.events.S7SysEvent;
import org.apache.plc4x.java.s7.events.S7UserEvent;

/**
 * Example for capturing events generated from a Siemens S7-300, S7-400 or VIPA PLC.
 * Support for mode events ("MODE"), system events ("SYS"), user events ("USR")
 * and alarms ("ALM").
 * Each consumer shows the fields and associated values of the "map" containing
 * the event parameters.
 */
public class EventSubscription {

    /**
     * @param args the command line arguments
     */
    public static void main(String[] args) throws Exception {
        try (PlcConnection connection = new PlcDriverManager().getConnection("s7://192.168.1.51?remote-rack=0&remote-slot=3&controller-type=S7_400")) {
            final PlcSubscriptionRequest.Builder subscription = connection.subscriptionRequestBuilder();

            subscription.addEventFieldAddress("myMODE", "MODE");
            subscription.addEventFieldAddress("mySYS", "SYS");
            subscription.addEventFieldAddress("myUSR", "USR");
            subscription.addEventFieldAddress("myALM", "ALM");

            final PlcSubscriptionRequest sub = subscription.build();
            final PlcSubscriptionResponse subresponse = sub.execute().get();

            //Si todo va bien con la subscripción puedo
            subresponse
                .getSubscriptionHandle("myMODE")
                .register(msg -> {
                    System.out.println("******** S7ModeEvent ********");
                    Map<String, Object> map = ((S7ModeEvent) msg).getMap();
                    map.forEach((x, y) -> {
                        System.out.println(x + " : " + y);
                    });
                    System.out.println("****************************");
                });

            subresponse
                .getSubscriptionHandle("mySYS")
                .register(msg -> {
                    System.out.println("******** S7SysEvent ********");
                    Map<String, Object> map = ((S7SysEvent) msg).getMap();
                    map.forEach((x, y) -> {
                        System.out.println(x + " : " + y);
                    });
                    System.out.println("****************************");
                });

            subresponse
                .getSubscriptionHandle("myUSR")
                .register(msg -> {
                    System.out.println("******** S7UserEvent *******");
                    Map<String, Object> map = ((S7UserEvent) msg).getMap();
                    map.forEach((x, y) -> {
                        System.out.println(x + " : " + y);
                    });
                    System.out.println("****************************");
                });

            subresponse
                .getSubscriptionHandle("myALM")
                .register(msg -> {
                    System.out.println("******** S7AlmEvent *********");
                    Map<String, Object> map = ((S7AlarmEvent) msg).getMap();
                    map.forEach((x, y) -> {
                        System.out.println(x + " : " + y);
                    });
                    System.out.println("****************************");
                });

            System.out.println("Waiting for events");

            Thread.sleep(120000);

            System.out.println("Bye...");

        }
    }

}

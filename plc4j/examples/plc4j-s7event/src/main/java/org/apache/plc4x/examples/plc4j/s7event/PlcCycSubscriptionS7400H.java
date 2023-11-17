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

import io.netty.util.concurrent.DefaultThreadFactory;
import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.s7.events.S7CyclicEvent;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Map;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.listener.ConnectionStateListener;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionResponse;
import org.apache.plc4x.java.s7.readwrite.connection.S7HDefaultNettyPlcConnection;
import org.apache.plc4x.java.s7.readwrite.protocol.S7HPlcConnection;


/**
 * Example for capturing events generated from a Siemens S7-300, S7-400 or VIPA PLC.
 * Support for mode events ("MODE"), system events ("SYS"), user events ("USR")
 * and alarms ("ALM").
 * Each consumer shows the tags and associated values of the "map" containing
 * the event parameters.
 */
public class PlcCycSubscriptionS7400H implements ConnectionStateListener {

    private static final Logger logger = LoggerFactory.getLogger(PlcCycSubscriptionS7400H.class);
    
    private S7HPlcConnection connection = null;
    private PlcSubscriptionResponse subresponse  = null;
    private boolean shutdown        = false;
    private boolean subscricted     = false; 
    private boolean unsubscricted   = false;
    /**
     * @param args the command line arguments
     */
    public static void main(String[] args) throws Exception {

        System.setProperty(org.slf4j.simple.SimpleLogger.DEFAULT_LOG_LEVEL_KEY, "Debug");        
        
        PlcCycSubscriptionS7400H myapp = new PlcCycSubscriptionS7400H();
        myapp.run(args);
      
                
    }
    
    public void makeconnection() {
        try {                
            connection =(S7HPlcConnection) new DefaultPlcDriverManager().getConnection("s7://10.10.1.80/10.10.1.81?remote-rack=0&remote-slot=3&remote-rack2=0&remote-slot=4&controller-type=S7_400&read-timeout=8&ping=true&ping-time=4&retry-time=5");
            connection.addEventListener(this);
            while(!connection.isConnected()){
                System.out.println("Conectando...");
            };
        } catch (PlcConnectionException ex) {
           System.out.println("RUNSCAN: " + ex.getMessage());
        }
    } 
    
    public void makesubscription() throws Exception {
        final PlcSubscriptionRequest.Builder subscription = connection.subscriptionRequestBuilder();
        //subscription.addEventTagAddress("myCYC", "CYC(B01SEC:5):%DB9002.DBB0[1]");
        subscription.addEventTagAddress("myCYC_01", "CYC(B01SEC:5):%MB190:BYTE");
        subscription.addEventTagAddress("myCYC_02", "CYC(B01SEC:5):%MW190:INT");            

        final PlcSubscriptionRequest sub = subscription.build();
        subresponse = sub.execute().get();

            //Si todo va bien con la subscripción puedo
//            subresponse
//                .getSubscriptionHandle("myMODE")
//                .register(msg -> {
//                    System.out.println("******** S7ModeEvent ********");
//                    Map<String, Object> map = ((S7ModeEvent) msg).getMap();
//                    map.forEach((x, y) -> {
//                        System.out.println(x + " : " + y);
//                    });
//                    System.out.println("****************************");
//                });
//
//            subresponse
//                .getSubscriptionHandle("mySYS")
//                .register(msg -> {
//                    System.out.println("******** S7SysEvent ********");
//                    Map<String, Object> map = ((S7SysEvent) msg).getMap();
//                    map.forEach((x, y) -> {
//                        if ("INFO1".equals(x)) {
//                            System.out.println(x + " : " + String.format("0x%04X", y));
//                        } else if ("INFO2".equals(x)) {
//                            System.out.println(x + " : " + String.format("0x%08X", y));    
//                        } else System.out.println(x + " : " + y);
//                    });
//                    System.out.println("****************************");
//                });
//
//            subresponse
//                .getSubscriptionHandle("myUSR")
//                .register(msg -> {
//                    System.out.println("******** S7UserEvent *******");
//                    Map<String, Object> map = ((S7UserEvent) msg).getMap();
//                    map.forEach((x, y) -> {
//                        System.out.println(x + " : " + y);
//                    });
//                    System.out.println("****************************");
//                });

         subresponse.getSubscriptionHandles().forEach( a -> {
             System.out.println(a.toString());
         });

        subresponse.getRequest().getTagNames().forEach(s -> {
            System.out.println("Tag " + s + " Tipo: " + subresponse.getRequest().getTag(s).getPlcValueType()); 
        });

        
             
        subresponse
            .getSubscriptionHandle("myCYC_01")
            .register(msg -> {
                if (null == msg) return;
                System.out.println("******** CYC Event *********");
                Map<String, Object> map = ((S7CyclicEvent) msg).getMap();
                map.forEach((x, y) -> {
                    if (x.startsWith("DATA_", 0)) {
                        System.out.println("Longitud de datos: " + ((byte[]) y).length);
                        System.out.println(x + ": " + Hex.encodeHexString((byte[]) y));
                    } else
                       System.out.println(x + " : " + y);
                });
                System.out.println("****************************");
            });
        
    }    

    @Override
    public void connected() {
        System.out.println("*************** CONECTADO *****************");
    }

    @Override
    public void disconnected() {
       System.out.println("*************** DESCONECTADO *****************");
       subscricted = false;        
    }
    
    
    public void makeunsubscription() throws Exception {
       System.out.println("*************** ELIMINANDO SUSCRIPCION *****************");        
        final PlcUnsubscriptionRequest.Builder unsubscription = connection.unsubscriptionRequestBuilder();
        
        unsubscription.addHandles(subresponse.getSubscriptionHandle("myCYC_01"));
        
        final PlcUnsubscriptionRequest res = unsubscription.build();
        
        res.execute().get();
        
    }    

    public void run (String[] args) throws Exception {
        try {
        
            makeconnection();
            
            System.out.println("Waiting for events");

            DefaultThreadFactory  tf = new
            DefaultThreadFactory("CYC", true);             

            Thread th01 = tf.newThread(() -> {
                      
                    while(!shutdown){
                        try {
                            
                            if ((connection.isConnected() == true) && (subscricted == false)) {
                                System.out.println("REALIZO LA SUSCRIPCION..." + connection.isConnected() + " : " + subscricted);
                                makesubscription();
                                subscricted = true;
                            } else if (connection.isConnected() == false) {
                                subscricted = false;    
                            } 
                            Thread.sleep(1000);
                        } catch (Exception ex){
                            System.out.println("Algo fallo aqui..." + ex.getMessage());
                            subscricted = false;
                        }
                    }
                }
            );
            
            th01.start();
            
            System.in.read();
            shutdown = true;
            
            System.out.println("Bye...");
            
            makeunsubscription();

            Thread.sleep(2000);
            
            connection.close();
            

        } catch (Exception ex) {
            System.out.println("Finaliza la ejecución...");
            ex.printStackTrace();
        };  
    }    
    

}

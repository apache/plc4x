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
import java.util.Scanner;
import java.util.concurrent.TimeUnit;

import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.s7.events.S7CyclicEvent;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.slf4j.simple.SimpleLogger;
import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.listener.ConnectionStateListener;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.s7.readwrite.protocol.S7HPlcConnection;


/**
 * Example for capturing events generated from a Siemens S7-300, S7-400 or VIPA PLC.
 * Support for cyclical (CYC) subscriptions.
 * Each consumer shows the tags and associated values of the "map" containing
 * the event parameters.
 */
public class CycSubscription implements ConnectionStateListener { //01

    private static final Logger logger = LoggerFactory.getLogger(CycSubscription.class); 
    private S7HPlcConnection connection = null;
    private PlcSubscriptionRequest.Builder subs_request_builder;
    private PlcSubscriptionRequest subs_request;
    
    private PlcSubscriptionResponse subs_response;
    
    private PlcConsumerRegistration csmr;
    
    private String jobId;
    
    /**
     * @param args the command line arguments
     */
    public static void main(String[] args) throws Exception {
        System.setProperty(SimpleLogger.DEFAULT_LOG_LEVEL_KEY, "debug");                   
            
        try  {
            CycSubscription cyc = new CycSubscription();
            
            cyc.OpenConnection();
            cyc.RegisterPlcTags();
            cyc.GetConnections();
            cyc.RegisterCYCHandles();
            
            System.out.println("Waiting for events...");
            Scanner scanner = new Scanner(System.in);
            System.out.println(scanner.nextLine());
            
            cyc.UnRegisterCYCHandles();
            cyc.CloseConnection();
            
            scanner.close();
                    
        } catch (Exception ex) {
            
        }
    }

    /*
    * Open the connection to the controller.
    * This connection specifies the monitoring parameters of the connection.
    * read-timeout = 8  ; Maximum wait time to see input data over the
    *                   ; TCP channel (sec.).
    * ping = true       ; Enable the PING function. (boolean)
    * ping-time = 4     ; Run the PING every N seconds.
    * retry-time = 5    ; Number of seconds to wait for connection monitoring. 
    *                   ; Default value zero (0) disables it.
    */
    public void OpenConnection() {
        try {
            connection = (S7HPlcConnection) new DefaultPlcDriverManager().
             getConnection("s7://10.10.1.33/10.10.1.34?remote-rack=0&remote-slot=3&remote-rack2=0&remote-slot=4&controller-type=S7_400&read-timeout=8&ping=true&ping-time=4&retry-time=5");  
            connection.addEventListener(this);           
        } catch (Exception ex) {
            ex.printStackTrace();
        }
    }
    
    /*
    * 
    */
    public void RegisterPlcTags(){
        subs_request_builder = connection.subscriptionRequestBuilder();

        subs_request_builder.addEventTagAddress("myCYC01", "CYC(B1SEC:2):%DB2.DBD2:REAL");
        subs_request_builder.addEventTagAddress("myCYC02", "CYC(B1SEC:2):%DB2.DBB3:BYTE"); 
        subs_request_builder.addEventTagAddress("myCYC03", "CYC(B1SEC:2):%MB190:BYTE");        

        subs_request = subs_request_builder.build();       
    }
    
    /*
    * Create the subscription Tags, only the first Tag is taken into account 
    * to define the request time.
    * Here the jobId is also captured, which is a value sent from the PLC. 
    * At the moment it is received instead of the Tagname assigned in 
    * the subscription.
    */
    public void GetConnections(){
        try {
            subs_response = subs_request.execute().get(); 
            jobId = (String) subs_response.getTagNames().toArray()[0];
        } catch (Exception  ex) {
            ex.printStackTrace();
        }       
    }
    
    /*
    * Closes the connection to the PLC.
    */
    public void CloseConnection() {          
        try {        
            connection.close();
        } catch (PlcConnectionException ex) {
            ex.printStackTrace();
        }
    }

    /*
    * Here I register my event consumer with the associated jobId. 
    */
    public void RegisterCYCHandles() {
        System.out.println("JOBID: " + jobId);
        csmr = subs_response
                .getSubscriptionHandle(jobId)
                .register( msg -> {
                    final S7CyclicEvent cycMsg = (S7CyclicEvent) msg;
                    System.out.println("******** <CYC Event 02> *********");
                    Map<String, Object> map = cycMsg.getMap();
                    map.keySet().forEach(key -> {System.out.println("Field available: " + key);});
                    byte[] data01 = (byte[]) map.get("myCYC01");
                    byte[] data02 = (byte[]) map.get("myCYC02");
                    System.out.println("myCYC01" + ": " + Hex.encodeHexString(data01));
                    System.out.println("myCYC02" + ": " + Hex.encodeHexString(data02));
                    System.out.println("myCYC01" + ": " + cycMsg.getFloat("myCYC01"));                    
                    System.out.println("myCYC02" + ": " + cycMsg.getByte("myCYC02"));
                    System.out.println("myCYC03" + ": " + cycMsg.getByte("myCYC03"));                    
                    System.out.println("****************************");
                });
    }
    
    /*
    * Secure unregistration process for event consumption.
    * 1. Unsubscribe from the event dispatcher.
    * 2. Cancels the subscription to the PLC, with the associated jobId. 
    */
    public void UnRegisterCYCHandles() {
        
        csmr.unregister();

        PlcSubscriptionResponse subresponse03 = null;  
       
        final PlcSubscriptionRequest.Builder subscription03 = connection.subscriptionRequestBuilder();
        subscription03.addEventTagAddress("CANCELTAG", "CANCEL:" + jobId);  
           
        final PlcSubscriptionRequest sub03 = subscription03.build();

        try {
            subresponse03 = sub03.execute().get(500, TimeUnit.MILLISECONDS);
        } catch (Exception ex) {
            ex.printStackTrace();
        }

        System.out.println("UnRegisterCYCHandles: " + subresponse03.getResponseCode(jobId));
      
    }   
    
    /*
    * For the S7H driver, up to three connections are foreseen, 
    * two real TCP and one virtual or embedded.
    * The connection of interest is the embedded one. 
    * It is the one that we must monitor until the internal 
    * machine changes its state to connected.
    */
    @Override
    public void connected() {
        System.out.println("Some connection return!");        
        if (connection.isConnected()) {
            System.out.println("Reconnecting return!"); 
            System.out.println("PASO 1!"); 
            RegisterPlcTags();
            System.out.println("PASO 2!");             
            GetConnections();
            System.out.println("PASO 3!");            
            RegisterCYCHandles();
            System.out.println("Salio de aqui!!!!.");
        }
    }

    /*
    * For the disconnection event we must unregister the event consumer.
    * For asynchronous events, subscriptions must be restarted.
    */
    @Override
    public void disconnected() {
        System.out.println("Trouble, trouble, we're offline! Status: " + connection.isConnected());
        csmr.unregister();
    }


}

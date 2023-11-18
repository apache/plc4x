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
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.s7.events.S7CyclicEvent;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Map;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.atomic.AtomicBoolean;
import org.apache.commons.lang3.time.StopWatch;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.listener.ConnectionStateListener;
import org.apache.plc4x.java.s7.readwrite.protocol.S7HPlcConnection;


/**
 * Cyclic subscription allows the acquisition of data from 
 * the controller in push mode, the PLC is responsible for sending 
 * the data to the client application. The minimum time base is 100 msec.
 * The data is sent as a stream of bytes so the client must 
 * maintain its consistency.This is the working mode of WinCC Scada.
 * The connection is supervised.
 * Plc:              SIMATIC S7-400
 * Model:            CPU 417
 * Part number:      6ES7 417-4XT05-0AB0
 * Firmware version: 5.1.0
 * CP1: 6GK7 443-1EX11-0XE0
 * CP2: 6GK7 443-1EX20-0XB0 
*/
public class PlcCycSubscriptionS7400H implements ConnectionStateListener {

    private static final Logger logger = LoggerFactory.getLogger(PlcCycSubscriptionS7400H.class);
    
    private S7HPlcConnection connection = null; 
    private AtomicBoolean isConnected = new AtomicBoolean(false);  
    private AtomicBoolean ShutDown = new AtomicBoolean(false);    
    private final DefaultThreadFactory  dtf = new DefaultThreadFactory("CYC", true);     
    /**
     * @param args the command line arguments
     */
    public static void main(String[] args) throws Exception {

        System.setProperty(org.slf4j.simple.SimpleLogger.DEFAULT_LOG_LEVEL_KEY, "Info");        
        
        PlcCycSubscriptionS7400H device = new PlcCycSubscriptionS7400H();
        device.run(args);                      
    }
    
    public void run (String[] args) throws Exception {
        
        logger.info("*****************************************************");
        logger.info("* Example of connection and read to a Simatic S7-400H");        
        logger.info("* Plc:              SIMATIC S7-400");
        logger.info("* Model:            CPU 417");
        logger.info("* Part number:      6ES7 417-4XT05-0AB0");
        logger.info("* Firmware version: 5.1.0");        
        logger.info("* CP1: 6GK7 443-1EX11-0XE0");   
        logger.info("* CP2: 6GK7 443-1EX20-0XB0");   
        logger.info("*");           
        logger.info("* Note: . ");  
        logger.info("*       . ");         
        logger.info("*****************************************************"); 
        

        
        logger.info("*****************************************************"); 
        logger.info("* 1. Once the connection is executed, it must "); 
        logger.info("*    suscrict the data contained in the address.");
        logger.info("*    URL to:s7://10.10.1.80/10.10.1.81?remote-rack=0&");
        logger.info("            remote-slot=3&remote-rack2=0&remote-slot=4&");
        logger.info("            controller-type=S7_400&read-timeout=8&");
        logger.info("            ping=true&ping-time=2&retry-time=3");
        logger.info("*    Press [ENTER]");        
        logger.info("*****************************************************"); 
        System.in.read();   
        
        OpenConnection("s7://10.10.1.80/10.10.1.81?remote-rack=0&"
                + "remote-slot=3&remote-rack2=0&remote-slot=4&"
                + "controller-type=S7_400&read-timeout=8&"
                + "ping=true&ping-time=2&retry-time=3"); //(01)        
        
        logger.info("*****************************************************"); 
        logger.info("* 2. In this step subscriptions are launched.");
        logger.info("     In the specific case of the S7-400,");
        logger.info("*    you can have up to 32 simultaneous subscriptions.");
        logger.info("*    You must be able to read the value from the PLC.");         
        logger.info("*    Press [ENTER]");        
        logger.info("*****************************************************"); 
        System.in.read();         

        MakeSubscription(); //(02)

        logger.info("*****************************************************"); 
        logger.info("* 3. Subscriptions are removed."); 
        logger.info("*    Depending on how you handle your tasks, ");
        logger.info("*    you should ensure that they are all completed.");         
        logger.info("*    Press [ENTER]");        
        logger.info("*****************************************************"); 
        System.in.read();         
        
        MakeUnsubscription(); //(03)

        logger.info("*****************************************************"); 
        logger.info("* 4. We close the connection and release resources.");        
        logger.info("*    Press [ENTER]");        
        logger.info("*****************************************************"); 
        System.in.read();         
   
        CloseConnection(); //(04)
             
    }       
    
    /***************************************************************************
    * Under normal conditions, the driver expects you to have the PLC 
    * connected to the network to start operations.
    * If a connection to the PLC cannot be established, an exception of type 
    * "PlcConnectionException" is generated, which must be handled by your 
    * application. In this example it waits for a connection to exist.
    * When you disconnect the "connection", the S7 driver will take care of the 
    * connection and reconnection process if necessary.
    * The internal wait time for the connection is one (01) second.
    ***************************************************************************/
    private void OpenConnection(String url) {
        int retrys = 0;
        StopWatch watch = new StopWatch();
        watch.start(); 
        while (null == connection) {
            try {        
                connection =(S7HPlcConnection) new DefaultPlcDriverManager().
                        getConnection(url);
                connection.addEventListener(this);
                while (!connection.isConnected());
                watch.stop();
                
                isConnected.set( connection.isConnected());
                
            logger.info("Time elapse for connection: {} in ms, with " +
                        "rettrys: {}",watch.getTime(), retrys);
            
            } catch (PlcConnectionException ex) {
                logger.info(ex.getMessage());
                 retrys++;
            }
        }        
    }
    
    /***************************************************************************
    * When the connection is closed, pending tasks and transactions are 
    * completed.
    * The rest of the work should be sent to the GC.
    ***************************************************************************/
    private void CloseConnection() {
        if (null == connection) return;
        isConnected.set(false);
        try {
            connection.close();
            connection = null; //GC do you job!.
            
        } catch (PlcConnectionException ex) {
            logger.info("CloseConnection: " + ex.getMessage());
        }
    }    
    
    /***************************************************************************
    * This method is called when there is a physical disconnection of the driver
    * Check the monitoring parameters given in the URL during connection.
    ***************************************************************************/ 
    public void MakeSubscription() throws Exception {
        
        Thread th01 = null;
        
        for (int i=1; i<=1; i++) {
        
            th01 = dtf.newThread(new SubscriptionRunnable(connection, i));
        
            th01.start();
            
        }

    }    

    /***************************************************************************
    * Depending on your design, you must close all subscribers
    * in an orderly manner.
    ***************************************************************************/         
    public void MakeUnsubscription() throws Exception {
        logger.info("Finish all subscriptions.");
        ShutDown.set(true);        
    }    
   
    /***************************************************************************
    * This method is called when the driver makes an internal TCP connection.
    * The first connection of the driver does not generate this event.
    * In the case of high availability systems, this signal should be used 
    * to restart subscriptions to events, alarms, etc. 
    ***************************************************************************/    
    @Override
    public void connected() {
        logger.info("*****************************************************");         
        logger.info("*************** Plc is connected. *******************");      
        logger.info("*****************************************************"); 
        isConnected.set(true);        
    }

    /***************************************************************************
    * This method is called when there is a physical disconnection of the driver
    * Check the monitoring parameters given in the URL during connection.
    ***************************************************************************/    
    @Override
    public void disconnected() {
        logger.info("*****************************************************");         
        logger.info("*************** Plc is disconnected. ****************");         
        logger.info("*****************************************************");         
        isConnected.set(false);
    }        
    
    /***************************************************************************
    * This object encapsulates the three steps required for cyclic subscription.
    * Try to handle all possible exceptions that are generated.
    ***************************************************************************/ 
    private class SubscriptionRunnable implements Runnable {
        private final Logger logger = LoggerFactory.getLogger(SubscriptionRunnable.class);          
        
        private final S7HPlcConnection myconnection; 
        private final int instance;
        private PlcSubscriptionRequest.Builder subscription;
        private PlcSubscriptionRequest sub;
        private PlcSubscriptionResponse subresponse;

        public SubscriptionRunnable(S7HPlcConnection connection, int instance) {
            this.myconnection = connection;
            this.instance = instance;
        }

        @Override
        public void run() {
            
            MakeMySubscription();
            
            MakeMyRegister();
            
            while (!ShutDown.get()) {
                try {
                Thread.sleep(100);
                } catch (Exception ex) {
                    logger.debug(ex.getMessage());
                }
            } ;
            
            MakeMyUnsubscription();
            
            logger.info("Bye!");            
        }
        
        /***********************************************************************
        * Registration of the different memory areas.
        * Only the first request sets the time period for the entire group.
        * The time bases are for the multiplier:
        * 
        * . B01SEC: Time base 0.1 Sec. (100 mSec.).
        * . B1SEC: Time base 1.0 Sec.
        * . B10SEC: Time base 10 Sec.
        ***********************************************************************/         
        private void MakeMySubscription() {
            subscription = myconnection.subscriptionRequestBuilder();
            subscription.addEventTagAddress(instance + "_myCYC_01", "CYC(B1SEC:5):%MB190:BYTE");
            subscription.addEventTagAddress(instance + "_myCYC_02", "CYC(B1SEC:5):%MW190:INT");
            
            sub = subscription.build();
            
            try {
                subresponse = sub.execute().get();
            } catch (InterruptedException ex) {
                logger.info(ex.getMessage());
            } catch (ExecutionException ex) {
                logger.info(ex.getMessage());
            }                       
        }
        
        /***********************************************************************
        * This object encapsulates the three steps required for cyclic subscription.
        * Try to handle all possible exceptions that are generated.
        ***********************************************************************/         
        private void MakeMyRegister() {
            
            subresponse
                .getSubscriptionHandle(instance + "_myCYC_01")
                .register(msg -> {
                    if (null == msg) return
                            ;
                    logger.info("******** {} CYC Event *********", instance);
                    Map<String, Object> map = ((S7CyclicEvent) msg).getMap();
                    map.forEach((x, y) -> {
                        if (x.startsWith("DATA_", 0)) {
                            logger.info("Longitud de datos: " + ((byte[]) y).length);
                            logger.info(x + ": " + Hex.encodeHexString((byte[]) y));
                        } else
                            logger.info(x + " : " + y);
                    });
                    logger.info("****************************");
                });;             
        }        
        
        /***********************************************************************
        * This object encapsulates the three steps required for cyclic subscription.
        * Try to handle all possible exceptions that are generated.
        ***********************************************************************/         
        private void MakeMyUnsubscription(){
            final PlcUnsubscriptionRequest.Builder unsubscription = myconnection.unsubscriptionRequestBuilder();
        
            unsubscription.addHandles(subresponse.getSubscriptionHandle(instance + "_myCYC_01"));
        
            final PlcUnsubscriptionRequest res = unsubscription.build();
        
            try {                        
                res.execute().get();
            } catch (InterruptedException ex) {
                logger.info(ex.getMessage());
            } catch (ExecutionException ex) {
                logger.info(ex.getMessage());
            }
        }
        
    }
    
    
}

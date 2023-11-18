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

import java.io.IOException;
import java.util.concurrent.atomic.AtomicBoolean;
import org.apache.commons.lang3.time.StopWatch;
import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.listener.ConnectionStateListener;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.s7.readwrite.protocol.S7HPlcConnection;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;


/*
* 
* Example of connection to a S7-400.
* The connection is supervised.
* Plc:              SIMATIC S7-400
* Model:            CPU 417
* Part number:      6ES7 417-4XT05-0AB0
* Firmware version: 5.1.0
* CP1: 6GK7 443-1EX11-0XE0
* CP2: 6GK7 443-1EX20-0XB0
*/
public class PlcReadDataS7400H implements ConnectionStateListener {
    
    private static final Logger logger = LoggerFactory.getLogger(PlcReadDataS7400H.class);   
    
    private S7HPlcConnection connection = null; 
    private AtomicBoolean isConnected = new AtomicBoolean(false);    

    /**
     * @param args the command line arguments
     */
    public static void main(String[] args) throws Exception {
        System.setProperty(org.slf4j.simple.SimpleLogger.DEFAULT_LOG_LEVEL_KEY, "Debug");         
        
        PlcReadDataS7400H device = new PlcReadDataS7400H();
        device.run();
    }
    
    
    public void run() throws IOException {
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
        
        OpenConnection("s7://10.10.1.80/10.10.1.81?remote-rack=0&"
                + "remote-slot=3&remote-rack2=0&remote-slot=4&"
                + "controller-type=S7_400&read-timeout=8&"
                + "ping=true&ping-time=2&retry-time=3"); //(01)
           
        logger.info("*****************************************************"); 
        logger.info("* 1. Once the connection is executed, it must read"); 
        logger.info("*    the data contained in the address.");
        logger.info("*    URL to:s7://10.10.1.80/10.10.1.81?remote-rack=0&");
        logger.info("            remote-slot=3&remote-rack2=0&remote-slot=4&");
        logger.info("            controller-type=S7_400&read-timeout=8&");
        logger.info("            ping=true&ping-time=2&retry-time=3");
        logger.info("*    Press [ENTER]");        
        logger.info("*****************************************************"); 
        System.in.read();
        
        Read(); //(01.1)
                 
        logger.info("*****************************************************"); 
        logger.info("* 2. Turn off/on PLC! This will cause the connection"); 
        logger.info("*    handlers to be lost. ");       
        logger.info("*    This simulates connection loss due to lack of ");
        logger.info("*    traffic, which is handled by OS. ");
        logger.info("*    When reading is executed, the driver must activate");
        logger.info("*    the reconnection process internally.");
        logger.info("*    In this operation the reading is lost, but ");
        logger.info("*    the reconnection process begins.");
        logger.info("*    Press [ENTER]");
        logger.info("*****************************************************"); 
        System.in.read();   
                
        Read(); //(02.1)        
        
        logger.info("*****************************************************"); 
        logger.info("* 3. The connection must be reestablished."); 
        logger.info("*    Remove primary connection.");         
        logger.info("*    Press [ENTER]");        
        logger.info("*****************************************************"); 
        System.in.read(); 
        
        Read(); //(03.1)     

        logger.info("*****************************************************"); 
        logger.info("* 4. Remove secondary connection.");        
        logger.info("*    Press [ENTER]");        
        logger.info("*****************************************************"); 
        System.in.read(); 
        
        Read(); //(03.1)         
        
        logger.info("*****************************************************"); 
        logger.info("* 5. Place primary connection."); 
        logger.info("*    Press [ENTER]");        
        logger.info("*****************************************************"); 
        System.in.read();  
        
        Read(); //(05.1)    
        
        logger.info("*****************************************************"); 
        logger.info("* 6. Place secondary connection."); 
        logger.info("*    Press [ENTER]");        
        logger.info("*****************************************************"); 
        System.in.read();  
        
        Read(); //(06.1)    
        
        logger.info("*****************************************************"); 
        logger.info("* 7. Once the connection is executed, it must read."); 
        logger.info("*    Press [ENTER]");        
        logger.info("*****************************************************");
        System.in.read();  
        
        Read(); //(07.1)    
        
        
        logger.info("*****************************************************"); 
        logger.info("* 8. And we close the connection."); 
        logger.info("*    Press [ENTER]");        
        logger.info("*****************************************************");        
           
        CloseConnection(); //(08.1)        
                
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
    * The reading process is standard. In case of an exception, 
    * the user must take the appropriate actions, but "do not close 
    * the connection":
    ***************************************************************************/    
    private void Read() {
        if (!isConnected.get()) return;
        try {
            final PlcReadRequest.Builder readrequest = connection.readRequestBuilder();  //(01)
            readrequest.addTagAddress("TEST", "%DB1000:4:INT"); //(02) 
            
            final PlcReadRequest rr = readrequest.build(); //(03)
            final PlcReadResponse response; //(04)            
            response = rr.execute().get(); //(05)
            
            if (response.getResponseCode("TEST") == PlcResponseCode.OK) { //(06)
                logger.info("Value: " + response.getString("TEST"));
            } else {
                logger.info("Problem reading...");                
            }              
        } catch (Exception ex) { //(07)
            logger.info("Read: " + ex.getMessage());
        };          
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
    
    
}

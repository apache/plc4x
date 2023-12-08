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
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicBoolean;
import org.apache.commons.lang3.time.StopWatch;
import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.listener.ConnectionStateListener;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
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
* 
* For the following program, the DB400 must be installed in the PLC, 
* with the following chains and the indicated lengths.
*
* DB400
* STRING001 - STRING[254]
* STRING002 - STRING[172]
* STRING003 - STRING[1]
* STRING004 - STRING[0]
* STRING005 - STRING[32]
* STRING006 - STRING[64]
* STRING007 - STRING[0]
* STRING008 - STRING[1]
* STRING009 - STRING[2]
* STRING010 - STRING[3]
* STRING011 - STRING[4]
* STRING012 - STRING[5]
* STRING013 - STRING[6]
* STRING014 - STRING[7]
* STRING015 - STRING[8]
*/
public class PlcReadWriteStrings implements ConnectionStateListener {

    private static final Logger logger = LoggerFactory.getLogger(PlcReadWriteStrings.class);   
    private static final long DELAY = 1000L;
    
    private static String TEST_STRING00 = "";    
    private static String TEST_STRING01 = "Y";    
    private static String TEST_STRING02 = "YZ";     
    private static String TEST_STRING08 = "01234567";
    private static String TEST_STRING254 = 
            "01234567890123456789012345678901234567890123456789"+    
            "01234567890123456789012345678901234567890123456789"+
            "01234567890123456789012345678901234567890123456789"+
            "01234567890123456789012345678901234567890123456789"+
            "01234567890123456789012345678901234567890123456789"+
            "0123";
    
    private S7HPlcConnection connection = null; 
    private AtomicBoolean isConnected = new AtomicBoolean(false);  
    
    private String[] tags = new String[]{"STRING007","STRING008","STRING009",
                                         "STRING010","STRING011","STRING012",
                                         "STRING013","STRING014","STRING015"};    
    
    /**
     * @param args the command line arguments
     */
    public static void main(String[] args) throws Exception {
        System.setProperty(org.slf4j.simple.SimpleLogger.DEFAULT_LOG_LEVEL_KEY, "Debug");         
        
        PlcReadWriteStrings device = new PlcReadWriteStrings();
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
                + "controller-type=S7_400&read-timeout=16&"                
                + "ping=true&ping-time=8&retry-time=3"); //(01)
           
        logger.info("*****************************************************"); 
        logger.info("* 1. Once the connection is executed, it must read"); 
        logger.info("*    the data contained in the address.");
        logger.info("*    URL to:s7://10.10.1.80/10.10.1.81?remote-rack=0&");
        logger.info("            remote-slot=3&remote-rack2=0&remote-slot=4&");
        logger.info("            controller-type=S7_400&read-timeout=8&");
        logger.info("            ping=true&ping-time=2&retry-time=3");
        logger.info("*    Press [ENTER]");        
        logger.info("*****************************************************"); 
        
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
        
        Write();
        
        logger.info("*****************************************************"); 
        logger.info("* Write null string.. Clean the DB in ST7/TIA.");     
        logger.info("* Press [ENTER]");        
        logger.info("*****************************************************");   
        System.in.read();          
        SafeWrite(TEST_STRING00);
        SafeRead(TEST_STRING00);
 
        
        logger.info("*****************************************************"); 
        logger.info("* Write one character string.. Clean the DB in ST7/TIA.");    
        logger.info("* Press [ENTER]");            
        logger.info("*****************************************************");  
        System.in.read();         
        SafeWrite(TEST_STRING01);    
        SafeRead(TEST_STRING01);        
  
                
        logger.info("*****************************************************"); 
        logger.info("* Write two character string.. Clean the DB in ST7/TIA.");    
        logger.info("* Press [ENTER]");            
        logger.info("*****************************************************");    
        System.in.read();           
        SafeWrite(TEST_STRING02);   
        SafeRead(TEST_STRING02);         

                
        logger.info("*****************************************************"); 
        logger.info("* Write eigth character string.. Clean the DB in ST7/TIA"); 
        logger.info("* Press [ENTER]");            
        logger.info("*****************************************************");  
        System.in.read();         
        SafeWrite(TEST_STRING08); 
        SafeRead(TEST_STRING08);         
  
                
        logger.info("*****************************************************"); 
        logger.info("* Write 254 character string.. Clean the DB in ST7/TIA");  
        logger.info("* Press [ENTER]");            
        logger.info("*****************************************************");  
        System.in.read();         
        SafeWrite(TEST_STRING254);    
        SafeRead(TEST_STRING254);          
  
                
        logger.info("*****************************************************"); 
        logger.info("* 8. And we close the connection."); 
        logger.info("*    Press [ENTER]");        
        logger.info("*****************************************************");        
        System.in.read();  
        
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
                //Avoid excessive CPU consumption
                try {
                    TimeUnit.MILLISECONDS.sleep(DELAY);
                } catch (InterruptedException ie) {
                    Thread.currentThread().interrupt();
                }                
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
            readrequest.addTagAddress("TAG01", "%DB400.DBX0.0:STRING"); //(02) 
            readrequest.addTagAddress("TAG06", "%DB400.DBX470.0:STRING"); //(02)             
            readrequest.addTagAddress("TAG02", "%MB190:BYTE"); //(02) 
            
            final PlcReadRequest rr = readrequest.build(); //(03)
            final PlcReadResponse response; //(04)            
            response = rr.execute().get(); //(05)
            
            if (response.getResponseCode("TAG01") == PlcResponseCode.OK) { //(06)
                logger.info("Value1: " + response.getString("TAG01"));
                logger.info("Value2: " + response.getString("TAG06"));                
                logger.info("Value3: " + response.getString("TAG02"));                
            } else {
                logger.info("Problem reading...");                
            }              
        } catch (Exception ex) { //(07)
            logger.info("Read: " + ex.getMessage());
        };          
    }    
    
    /***************************************************************************
    * The reading process is standard. In case of an exception, 
    * the user must take the appropriate actions, but "do not close 
    * the connection":
    ***************************************************************************/    
    private void Write() {
        if (!isConnected.get()) return;
        try {
            final PlcWriteRequest.Builder writeRequest = connection.writeRequestBuilder();  //(01)
            writeRequest.addTagAddress("TAG01", "%DB400.DBX0.0:STRING",TEST_STRING254 ); //(02) 
            writeRequest.addTagAddress("TAG06", "%DB400.DBX470.0:STRING",TEST_STRING254 ); //(02)           
            writeRequest.addTagAddress("TAG08", "%DB400.DBX436.0:STRING",TEST_STRING254 ); //(02)              
            
            final PlcWriteRequest wr = writeRequest.build(); //(03)
            final PlcWriteResponse response; //(04)            
            response = wr.execute().get(); //(05)
            
            if (response.getResponseCode("TAG01") == PlcResponseCode.OK) { //(06)
                logger.info("TAG01 Write sucefull...");
            } else {
                logger.info("TAG01 Problem reading...");                
            }  

            if (response.getResponseCode("TAG06") == PlcResponseCode.OK) { //(06)
                logger.info("TAG06 Write sucefull...");
            } else {
                logger.info("TAG06 Problem reading...");                
            }  
            
            if (response.getResponseCode("TAG08") == PlcResponseCode.OK) { //(06)
                logger.info("TAG07 Write sucefull...");
            } else {
                logger.info("TAG07 Problem reading...");                
            }              
            
        } catch (Exception ex) { //(07)
            logger.info("Read: " + ex.getMessage());
        };          
    }     
    
    /***************************************************************************
    * Write the String "str" to the different test memory areas.
    * No overlap of contents should be shown.
    ***************************************************************************/    
    private void SafeWrite(String str) {
        if (!isConnected.get()) return;
        try {
            final PlcWriteRequest.Builder writeRequest = connection.writeRequestBuilder();  //(01)
            writeRequest.addTagAddress(tags[0], "%DB400.DBX536.0:STRING",str); //(02) 
            writeRequest.addTagAddress(tags[1], "%DB400.DBX538.0:STRING",str); //(02)  
            writeRequest.addTagAddress(tags[2], "%DB400.DBX542.0:STRING",str); //(02) 
            writeRequest.addTagAddress(tags[3], "%DB400.DBX546.0:STRING",str); //(02)  
            writeRequest.addTagAddress(tags[4], "%DB400.DBX552.0:STRING",str); //(02) 
            writeRequest.addTagAddress(tags[5], "%DB400.DBX558.0:STRING",str); //(02)  
            writeRequest.addTagAddress(tags[6], "%DB400.DBX566.0:STRING",str); //(02) 
            writeRequest.addTagAddress(tags[7], "%DB400.DBX574.0:STRING",str); //(02) 
            writeRequest.addTagAddress(tags[8], "%DB400.DBX584.0:STRING",str); //(02)              
            
            final PlcWriteRequest wr = writeRequest.build(); //(03)
            final PlcWriteResponse response; //(04)            
            response = wr.execute().get(); //(05)
            
            for (String tag:tags) {
                if (response.getResponseCode(tag) == PlcResponseCode.OK) //(06)
                    logger.info(tag + " Write sucefull...");                
            }
            
        } catch (Exception ex) { //(07)
            logger.info("Read: " + ex.getMessage());
        };          
    }     
    
    /***************************************************************************
    * The reading process is standard. 
    * Each "value" obtained from the PLC is compared with "str".
    * If everything is fine, it should indicate a successful reading.
    ***************************************************************************/    
    private void SafeRead(String str) {
        int index = 0;
        int pos = 0;
        if (!isConnected.get()) return;
        try {
            final PlcReadRequest.Builder readRequest = connection.readRequestBuilder();  //(01)
            readRequest.addTagAddress(tags[0], "%DB400.DBX536.0:STRING"); //(02) 
            readRequest.addTagAddress(tags[1], "%DB400.DBX538.0:STRING"); //(02)  
            readRequest.addTagAddress(tags[2], "%DB400.DBX542.0:STRING"); //(02) 
            readRequest.addTagAddress(tags[3], "%DB400.DBX546.0:STRING"); //(02)  
            readRequest.addTagAddress(tags[4], "%DB400.DBX552.0:STRING"); //(02) 
            readRequest.addTagAddress(tags[5], "%DB400.DBX558.0:STRING"); //(02)  
            readRequest.addTagAddress(tags[6], "%DB400.DBX566.0:STRING"); //(02) 
            readRequest.addTagAddress(tags[7], "%DB400.DBX574.0:STRING"); //(02) 
            readRequest.addTagAddress(tags[8], "%DB400.DBX584.0:STRING"); //(02)  
            
            final PlcReadRequest rr = readRequest.build(); //(03)
            final PlcReadResponse response; //(04)            
            response = rr.execute().get(); //(05)
            
            
            for (String tag:tags) {
                if (response.getResponseCode(tag) == PlcResponseCode.OK){
                    pos = (index <= str.length())?index:str.length();
                    if (response.getString(tag).equals(str.substring(0, pos))){
                        logger.info(tag + " Read sucefull...");      
                    } else {
                        logger.info(tag + ": " + response.getString(tag) + " : " + str.substring(0, index));     
                    }
                }
                index++;
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

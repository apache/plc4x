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
import java.util.concurrent.ExecutionException;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.logging.Level;
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
* Example of connection to a Simatic S7-1200.
* The connection is supervised.
* Plc:              SIMATIC S7-1200
* Model:            CPU 1214C DC/DC/DC
* Part number:      6ES7 214-1AE30-0XB0
* Firmware version: 2.2
*/
public class PlcReadS71200 implements ConnectionStateListener {

    private static final Logger logger = LoggerFactory.getLogger(PlcReadS71200.class);   
    
    private S7HPlcConnection connection = null; 
    private AtomicBoolean isConnected = new AtomicBoolean(false);
    
    
    /**
     * @param args the command line arguments
     */
    public static void main(String[] args) throws Exception{
        System.setProperty(org.slf4j.simple.SimpleLogger.DEFAULT_LOG_LEVEL_KEY, "Debug"); 
        
        PlcReadS71200 device = new PlcReadS71200();
        device.run();
    }
    
    public void run() throws IOException {
        logger.info("*****************************************************");
        logger.info("* Example of connection and read to a Simatic S7-1200");        
        logger.info("* Plc:              SIMATIC S7-1200");
        logger.info("* Model:            CPU 1214C DC/DC/DC");
        logger.info("* Part number:      6ES7 214-1AE30-0XB0");
        logger.info("* Firmware version: 2.2");        
        logger.info("*");     
        logger.info("* Note: . All DBs must be non-optimized.");  
        logger.info("*       . PLC can be pinged.");         
        logger.info("*****************************************************"); 
        
        OpenConnection("s7://10.10.1.46"); //(01)
           
        logger.info("*****************************************************"); 
        logger.info("* 1. Once the connection is executed, it must read"); 
        logger.info("*    the data contained in the address.");
        logger.info("*    URL to: s7://10.10.1.46.");
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
        
        while (!isConnected.get());  //(04)

        logger.info("*****************************************************"); 
        logger.info("* 3. The connection must be reestablished."); 
        logger.info("*    You must be able to read the value from the PLC..");         
        logger.info("*    Press [ENTER]");        
        logger.info("*****************************************************"); 
        System.in.read(); 
        
        Read(); //(03.1)     
        
        logger.info("*****************************************************"); 
        logger.info("* 4. Now we close the connection and open it using"); 
        logger.info("*    other parameters.");    
        logger.info("*    The new connection is given by");            
        logger.info("*    URL to: s7://10.10.1.46?read-timeout=6&ping=true&ping-time=2");        
        logger.info("*    Press [ENTER]");        
        logger.info("*****************************************************");         
        System.in.read();     
        
        CloseConnection(); //(04.1)
        OpenConnection("s7://10.10.1.46?read-timeout=6&ping=true&ping-time=2"); //(04.2)  
        
        
        logger.info("*****************************************************"); 
        logger.info("* 5. Once the connection is executed, it must read."); 
        logger.info("*    Press [ENTER]");        
        logger.info("*****************************************************"); 
        System.in.read();  
        
        Read(); //(05.1)    
        
        logger.info("*****************************************************"); 
        logger.info("* 6. Turn off/on PLC! This will cause the connection"); 
        logger.info("*    handlers to be lost. ");
        logger.info("*    The driver will try to reconnect. ");        
        logger.info("*    Press [ENTER]");        
        logger.info("*****************************************************");  
        System.in.read();  
        
        while (!isConnected.get());          
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
    * When you disconnect the connection, the S7 driver will take care of the 
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
        try {
            final PlcReadRequest.Builder readrequest = connection.readRequestBuilder();  //(01)
            readrequest.addTagAddress("TEST", "%DB100:10:INT"); //(02) 
            
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

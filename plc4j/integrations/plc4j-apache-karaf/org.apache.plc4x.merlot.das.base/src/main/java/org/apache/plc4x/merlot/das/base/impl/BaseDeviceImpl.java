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

package org.apache.plc4x.merlot.das.base.impl;


import org.apache.plc4x.merlot.das.api.Driver;
import org.apache.plc4x.merlot.das.api.DriverCallback;
import org.apache.plc4x.merlot.das.api.DriverEvent;
import org.apache.plc4x.merlot.das.core.Merlot;
import org.apache.plc4x.merlot.das.base.api.BaseDevice;
import org.apache.plc4x.merlot.das.base.api.BaseDriver;
import java.time.LocalTime;
import java.time.format.DateTimeFormatter;
import java.util.ArrayList;
import java.util.List;
import java.util.function.Consumer;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.epics.pvdata.pv.ScalarType;
import static org.epics.pvdata.pv.ScalarType.pvBoolean;
import org.osgi.framework.BundleContext;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class BaseDeviceImpl implements BaseDevice, DriverCallback  {
	
    private static final Logger LOGGER = LoggerFactory.getLogger(BaseDeviceImpl.class);
    
    protected final BundleContext bc;
    protected boolean autostart = false;
    protected Driver driver = null;
    protected String url = null;
    
    public BaseDeviceImpl(BundleContext bc) {
        this.bc = bc;
    }
    
    @SuppressWarnings("unchecked")
    @Override
    public void init() throws Exception {    
    }

    @Override
    public void destroy() throws Exception {
    }

    @Override
    public void start() {
        autostart = true;
    }

    @Override
    public void stop() {
    }

    @Override
    public void attach(Driver driver) {
        this.driver = driver;
    }
            
    @Override
    public void putEvent(DriverEvent event) {  
        if (event.getCallback() == null) {
            event.setCallback(this);
        }
        driver.putEvent(event);
    }

    @Override
    public DriverEvent getEvent() {
        return driver.getEvent();
    }

    @Override
    public void noDriverFound() {
        LOGGER.info("no driver found...");
        driver = null;
    }

    @Override
    public void execute(DriverEvent cb) {
        LOGGER.info(cb.toString());
    }

    @Override
    public void doEvent(DriverEvent event, String tag) {
        throw new UnsupportedOperationException("Not supported yet."); 
    }

    @Override
    public void ReadRequest(String index, String item, DriverCallback cb) {
        PlcConnection plcConn = ((BaseDriver) driver).getPlcConnection();
        if (plcConn.isConnected()){
            PlcReadRequest.Builder requestBuilder = plcConn.readRequestBuilder();
            requestBuilder.addTagAddress(index, item);
            PlcReadRequest readRequest = requestBuilder.build(); //EX

            DriverEvent event = ((BaseDriver) driver).getEvent();
            event.setFunctionCode(Merlot.FUNCTION.FC_READ_MEMORY_BYTES);
            event.setPlcReadRequest(readRequest);
            event.setCallback((cb == null)?this:cb); 
            ((BaseDriver) driver).putEvent(event);
        } else {
            LOGGER.info("Device no connected");
        }
    }

    @Override
    public void WriteRequest(String scalar, String id, List<String> values, DriverCallback cb) {
        DateTimeFormatter dtf = DateTimeFormatter.ofPattern("HH:mm:ss:SSS");
        PlcConnection plcConn = ((BaseDriver) driver).getPlcConnection();
        if (plcConn.isConnected()){
            PlcWriteRequest.Builder requestBuilder = plcConn.writeRequestBuilder();
            List<Object> items = new ArrayList();
            switch(ScalarType.getScalarType(scalar)){
                case pvBoolean:for (String value:values)items.add(Boolean.parseBoolean(value));
                    break;                    
                case pvByte:
                case pvUByte:  for (String value:values)items.add(Byte.parseByte(value));                
                    break;
                case pvDouble: for (String value:values)items.add(Double.parseDouble(value));                         
                    break;
                case pvFloat:  for (String value:values)items.add(Float.parseFloat(value));     
                    break;
                case pvInt:
                case pvUInt:   for (String value:values)items.add(Integer.decode(value));    
                    break;
                case pvLong:
                case pvULong:  for (String value:values)items.add(Long.parseLong(value));
                    break;
                case pvShort:
                case pvUShort: for (String value:values)items.add(Short.parseShort(value));
                    break;
                case pvString: for (String value:values)items.add(value);
                    break;

            }
            requestBuilder.addTagAddress("Write@["+LocalTime.now().format(dtf)+"]", id, items.toArray());            
                     
            PlcWriteRequest writeRequest = requestBuilder.build(); //EX

            DriverEvent event = ((BaseDriver) driver).getEvent();
            event.setFunctionCode(Merlot.FUNCTION.FC_WRITE_DATA_BYTES);
            event.setPlcReadRequest(null);
            event.setPlcWriteRequest(writeRequest);
            event.setCallback((cb == null)?this:cb); 
            ((BaseDriver) driver).putEvent(event);
        } else {
            LOGGER.info("Device no connected");
        }        

    }
    
    @Override
    public void SubscriptionRequest(String... events) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public void UnsubscriptionRequest(String... events) {
        throw new UnsupportedOperationException("Not supported yet.");
    }   

    @Override
    public void ConsumerRegister(String event, Consumer<PlcSubscriptionEvent> consumer) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public void ConsumerUnRegister(String event) {
        throw new UnsupportedOperationException("Not supported yet.");
    }
    
 
    @Override
    public void setUrl(String url) {
        this.url = url;
    }

    @Override
    public String getUrl() {
        return url;
    }

}

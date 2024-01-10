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
package org.apache.plc4x.merlot.db.core;

import org.apache.plc4x.merlot.das.base.api.BaseDriver;
import org.apache.plc4x.merlot.db.api.DBControl;
import org.apache.plc4x.merlot.db.api.DBRecordFactory;
import org.apache.plc4x.merlot.scheduler.api.Job;
import org.apache.plc4x.merlot.scheduler.api.JobContext;
import java.util.ArrayList;
import java.util.Collections;
import java.util.Dictionary;
import java.util.Enumeration;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import org.apache.plc4x.java.api.PlcConnection;
import org.epics.pvdatabase.PVDatabase;
import org.epics.pvdatabase.PVRecord;
import org.osgi.framework.BundleContext;
import org.osgi.framework.Constants;
import org.osgi.framework.ServiceReference;
import org.osgi.service.cm.ConfigurationException;
import org.osgi.service.cm.ManagedServiceFactory;
import org.osgi.service.device.Driver;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;


public class DBRecordsManagedService implements ManagedServiceFactory, Job {
    private static final Logger LOGGER = LoggerFactory.getLogger(DBRecordsManagedService.class);  
    private final PVDatabase master;
    private static Map<String, Dictionary<String, ?>> waitingConfigs = null;    
    static final String PID = "org.apache.plc4x.merlot.db.records";   
    static final String FILE_PATH = "felix.fileinstall.filename";
 
    private String filter =  "(&(" + Constants.OBJECTCLASS + "=" + DBRecordFactory.class.getName() + ")"+
                           "(db.record.type=*))";

    private ServiceReference[] references = null;    
    private final BundleContext bundleContext; 
    
    private final DBControl dbControl;

    public DBRecordsManagedService(BundleContext bundleContext, PVDatabase master, DBControl dbControl) {
        this.bundleContext = bundleContext;
        this.master = master;
        this.dbControl = dbControl;
        waitingConfigs = Collections.synchronizedMap(new HashMap<String, Dictionary<String, ?>>());
    }
         
    @Override
    public String getName() {
        return PID;
    }

    @Override
    public void updated(String pid, Dictionary<String, ?> props) throws ConfigurationException {        
        String dataValue= null;
        String[] dataFields = null;
        DBRecordFactory recordFactory = null;
        PVRecord record = null;
        String device = null;  
        String strScalarType = null; 
        List<PVRecord> pvRecords = new ArrayList();
        String filename = (String) props.get("felix.fileinstall.filename");

        if (props.size() < 3){
            waitingConfigs.put(pid, props);  
            return;
        }
         
        if (filename != null){            
            int start = filename.lastIndexOf("-");
            int end = filename.indexOf(".cfg");
            if ((start > 0) && (end > 0)) {
                device = filename.substring(start+1, end);
            }
            if (device==null){
                LOGGER.debug("Bad configuration name found: " + filename);
                return;
            }
        } else {
                LOGGER.debug("Configuration file name not found.");
                return;
        }
        
        BaseDriver dummyDriver = getDriver(device);
        if (dummyDriver == null){
            LOGGER.debug("Device driver [" + device + "] is not deployed.");
            waitingConfigs.put(device, props); 
            return;
        }
        
        PlcConnection connPlc = dummyDriver.getPlcConnection();
        
        if (connPlc == null){
            LOGGER.debug("Device driver [" + device + "] native driver is not present.");
            waitingConfigs.put(device, props); 
            return;            
        }
        
        if(!connPlc.isConnected()){
            LOGGER.debug("Device driver [" + device + "] is not connected.");
            waitingConfigs.put(device, props); 
            return;            
        }

        if (props!=null) {
            Enumeration<String> keys = props.keys();
            for (Enumeration e = props.keys(); e.hasMoreElements();) {
                Object key = e.nextElement(); 
                dataValue = props.get(key).toString();                
                dataFields = dataValue.split(",");
                
                int start = dataFields[0].indexOf('[');
                int end = dataFields[0].indexOf(']');
                if ((start>0)){
                    strScalarType = dataFields[0].substring(0, start);
                } else {
                    strScalarType = dataFields[0];
                }
                
                System.out.println("Tipo: " + strScalarType);
                
                recordFactory = getRecordFactory(strScalarType);
                
                if (recordFactory != null){
                    PVRecord pvRecord = recordFactory.create(key.toString(), dataFields);

                    if (pvRecord != null){
                        pvRecords.add(pvRecord);
                    } else {
                        LOGGER.info("PVRecord '" + key +"' could not be created.");
                    }
                }
            }

            dbControl.attach(device, pvRecords);
            //Si todo OK, los agregoa a la base de datos
            pvRecords.forEach(pvr -> master.addRecord(pvr));

        }
    }
 
    @Override
    public void deleted(String pid) {
        LOGGER.info("Deleting config: " + pid);
    }
    
 
    @Override
    public void execute(JobContext arg0) {
        String pid = null;
        Dictionary<String, ?> props = null;
        Set<String> keys = new HashSet<>();
        keys.addAll(waitingConfigs.keySet());

        for (String key:keys){
            pid = key;
            props = waitingConfigs.remove(key);            
            try {
                updated(pid,props);
            } catch (ConfigurationException ex) {
                LOGGER.debug("Problem updating [" + key +"] from waiting list." );
            }
        }
    }    
    
    private DBRecordFactory getRecordFactory(String type){
        try {
            String strFilter = filter.replace("*", type);
            references = bundleContext.getServiceReferences((String) null, strFilter); 
            if (references != null){
               return (DBRecordFactory) bundleContext.getService(references[0]);
            } else {
                LOGGER.info("DBRecordFactory type: '" + type + "' don't exist.");
                return null;
            }
        } catch (Exception ex) {
            LOGGER.info("getRecordFactory: " + ex.toString());
        }
        return null;
    }
    
    private BaseDriver getDriver(String device){
        try{
            String filterdriver =  "(DRIVER_ID=" + device + ")"; 
            ServiceReference[] refdrvs = bundleContext.getAllServiceReferences(Driver.class.getName(), filterdriver);
            BaseDriver refdrv = (BaseDriver) bundleContext.getService(refdrvs[0]);
            if (refdrv==null) LOGGER.info("BasicDriver [" + device + "] don't found");
            return refdrv;            
        } catch (Exception ex){
            LOGGER.debug("getDriver: " + ex.toString());
        }
        return null;
    }    

    
}

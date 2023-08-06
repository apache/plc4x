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
package org.apache.plc4x.merlot.db.impl;

import org.apache.plc4x.merlot.db.api.DBCollector;
import org.apache.plc4x.merlot.db.api.DBControl;
import org.apache.plc4x.merlot.scheduler.api.Job;
import org.apache.plc4x.merlot.scheduler.api.JobContext;
import java.util.Collection;
import java.util.Hashtable;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;
import org.epics.pvdatabase.PVRecord;
import org.osgi.framework.BundleContext;
import org.osgi.framework.Constants;
import org.osgi.framework.ServiceReference;


public class DBControlImpl implements DBControl {

    
    BundleContext bundleContext;
    
    
    @Override
    public void attach(String device, Collection<PVRecord> pvRecords) {

        List<PVRecord> pvTemp = null;
        PVRecord pvRecord = null;
        
        Map<String, List<PVRecord>> byScanRate = pvRecords.stream()
        .collect(Collectors.groupingBy(record -> record.getPVStructure().getStringField("scan_rate").get()));        
        
        //System.out.println("Keys: " + byScanRate.keySet());
        for (String key:byScanRate.keySet()){
            DBCollector collector = registerCollector(device,key);
            collector.attach(byScanRate.get(key));
        }
        
    }

    @Override
    public void detach(PVRecord pvRecord) {
        System.out.println("detach PVRecord: "+ pvRecord.getPVStructure().getShortField("id"));
    }

    @Override
    public void execute(JobContext ctx) {

    }

    @Override
    public void setBundleContext(BundleContext bundleContext) {
        this.bundleContext = bundleContext;
    }

    private DBCollector registerCollector(String device, String scanrate){
        Hashtable properties = new Hashtable(); 
        properties.put("collector.das.device", device);
        properties.put("scheduler.name", device+":"+scanrate);
        properties.put("scheduler.period", scanrate);
        properties.put("scheduler.immediate", true);
        properties.put("scheduler.concurrent", false);
        
        DBCollector collector = new DBCollectorImpl(bundleContext, device);
        
        bundleContext.registerService(new String[]{DBCollector.class.getName(), Job.class.getName()}, collector, properties);         
        
        return collector;
    }
    
    private DBCollector getCollector(String driver, String device, String scanrate){
        try
        {
            String filter =  "(&(" + Constants.OBJECTCLASS + "=" + DBCollector.class.getName() + ")" +
                             "(&(scheduler.name=" + device+":"+scanrate + ")(scheduler.period=" + scanrate + ")))";
            
            ServiceReference[] references = bundleContext.getServiceReferences((String) null, filter);
            if (references != null){
                ServiceReference reference = references[0];
                return (DBCollector) reference.getBundle();
            }            
        } catch(Exception ex) {
            
        };
        return null;
    }
    
    
}

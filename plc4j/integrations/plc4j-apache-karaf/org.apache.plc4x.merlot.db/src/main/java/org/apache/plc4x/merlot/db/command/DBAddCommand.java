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
package org.apache.plc4x.merlot.db.command;

import org.apache.plc4x.merlot.db.api.DBControl;
import org.apache.plc4x.merlot.db.api.DBRecordFactory;
import org.apache.karaf.shell.api.action.Action;
import org.apache.karaf.shell.api.action.Argument;
import org.apache.karaf.shell.api.action.Command;
import org.apache.karaf.shell.api.action.Option;
import org.apache.karaf.shell.api.action.lifecycle.Reference;
import org.apache.karaf.shell.api.action.lifecycle.Service;
import org.epics.pvdata.pv.PVBoolean;
import org.epics.pvdata.pv.PVString;
import org.epics.pvdata.pv.PVStructure;
import org.epics.pvdatabase.PVDatabase;
import org.epics.pvdatabase.PVRecord;
import org.osgi.framework.BundleContext;
import org.osgi.framework.Constants;
import org.osgi.framework.ServiceReference;


@Command(scope = "db", name = "add", description = "Add a Record (tag) to database.")
@Service
public class DBAddCommand  implements Action {
  
    @Reference
    BundleContext bundleContext;
    
    @Reference
    PVDatabase master;
    
    @Reference
    DBControl dbControl;
    
    @Option(name = "-n", aliases = "--new", description = "New DBRecord.", required = false, multiValued = false)
    boolean driver;
    
    @Option(name = "-s", aliases = "--start", description = "Stop and delete the S7 device.", required = false, multiValued = false)
    boolean device;  
    
    @Option(name = "-a", aliases = "--array", description = "New DBRecord of the define type.", required = false, multiValued = false)
    boolean array = false;      
    
    @Argument(index = 0, name = "type", description = "PVType of the record.", required = true, multiValued = false)
    String type = null;   
    
    @Argument(index = 1, name = "name", description = "Short name of the record.", required = true, multiValued = false)
    String name = null;

    @Argument(index = 2, name = "id", description = "PLC address of the record.", required = true, multiValued = false)
    String id = null; 
    
    @Argument(index = 3, name = "scan", description = "Scan rate for the record.", required = true, multiValued = false)
    String scan = null;    
    
    @Argument(index = 4, name = "des", description = "Full description of the record.", required = false, multiValued = false)
    String descriptor = null;    

    
    public Object execute() throws Exception {
        
        boolean isArray = false;
        int lengthArray = 0;
        PVRecord record;
        
        String filter =  "(&(" + Constants.OBJECTCLASS + "=" + DBRecordFactory.class.getName() + ")"+
                           "(db.record.type=" + type + "))";
            
        ServiceReference[] references = bundleContext.getServiceReferences((String) null, filter);
        
        if (references != null){
            ServiceReference reference = references[0];
            if (type.equalsIgnoreCase((String)reference.getProperty("db.record.type"))){
                DBRecordFactory recordFactory = (DBRecordFactory) bundleContext.getService(reference);

                int start = id.indexOf('[')+1;
                int end = id.indexOf(']');
                if ((start>0) && (end!=-1)){
                    String strLength = id.substring(id.indexOf('[')+1, id.indexOf(']'));
                    lengthArray = Integer.parseInt(strLength);
                    isArray = true;
                }
                
                if (!isArray){
                    record = recordFactory.create(name);
                } else {
                    record = recordFactory.createArray(name,lengthArray);
                }

                PVStructure structure = record.getPVStructure();
                PVString pvDes = structure.getStringField("descriptor");
                pvDes.put(descriptor);
                PVString pvId = structure.getStringField("id");
                pvId.put(id);
                PVString pvScan = structure.getStringField("scan_rate");
                pvScan.put(scan);
                PVBoolean pvScanEnable = structure.getBooleanField("scan_enable");
                pvScanEnable.put(true);
                PVBoolean pvWriteEnable = structure.getBooleanField("write_enable");
                pvWriteEnable.put(true);                

                //TODO: Devolver "true" si se pudo agregar el record.
                //dbControl.attach(record);
                //TODO: Agregar solamente si se pudo agregar al driver
                //master.addRecord(record);  

                System.out.println("Record: \r\n" + record.toString());


            }
        }

        
        return null;
    }
}

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

import org.apache.plc4x.merlot.db.api.DBRPCReferenceListener;
import org.epics.pvaccess.server.rpc.RPCService;
import org.epics.pvdatabase.PVDatabase;
import org.epics.pvdatabase.PVRecord;
import org.osgi.framework.BundleContext;
import org.osgi.framework.ServiceReference;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;


public class DBRPCReferenceListenerImpl implements DBRPCReferenceListener {
    private static final Logger LOGGER = LoggerFactory.getLogger(DBRPCReferenceListenerImpl.class); 
    private final BundleContext bundleContext;
    private final PVDatabase master;

    public DBRPCReferenceListenerImpl(BundleContext bundleContext,
            PVDatabase master) {
        this.bundleContext = bundleContext;
        this.master = master;
    }
    
    
    @Override
    public void bind(ServiceReference reference) {
        try {
            if (reference == null) return;
            PVRecord rpcservice = (PVRecord) bundleContext.getService(reference);
            boolean isRecordAdd = master.addRecord(rpcservice);
            if (isRecordAdd){
                LOGGER.info("Service {} was Add.",reference.getProperty("rpcservice.id"));
            } else  {
                LOGGER.info("Service {} can't be Add.",reference.getProperty("rpcservice.id"));            
            }
        } catch (Exception ex){
            LOGGER.info(ex.getMessage());
        }        
    }

    //TODO: Check for compatible type
    @Override
    public void bind(RPCService service) {
        try {
            PVRecord record = (PVRecord) service;
                   
            boolean isRecordAdd = master.addRecord(record);
            if (isRecordAdd){
                LOGGER.info("Service {} was add.");
            } else  {
                LOGGER.info("Service {} can't be add.");            
            }
        } catch (Exception ex){
            LOGGER.info(ex.getMessage());
        }
    }

    @Override
    public void unbind(ServiceReference reference) {
        try {
            if (reference == null) return;
            PVRecord rpcservice = (PVRecord) bundleContext.getService(reference);
            boolean isRecordRemove = master.removeRecord(rpcservice);
            if (isRecordRemove){
                LOGGER.info("Service {} was remove.",reference.getProperty("rpcservice.id"));
            } else  {
                LOGGER.info("Service {} can't be remove.",reference.getProperty("rpcservice.id"));            
            }
        } catch (Exception ex){
            LOGGER.info(ex.getMessage());
        }        
    }
    
}

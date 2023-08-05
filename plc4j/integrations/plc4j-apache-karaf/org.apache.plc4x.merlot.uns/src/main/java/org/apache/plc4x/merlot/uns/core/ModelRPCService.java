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
package org.apache.plc4x.merlot.uns.core;

import org.apache.plc4x.merlot.uns.api.Model;
import java.util.UUID;
import org.epics.nt.NTScalar;
import org.epics.nt.NTScalarBuilder;
import org.epics.pvaccess.server.rpc.RPCRequestException;
import org.epics.pvaccess.server.rpc.RPCService;
import org.epics.pvaccess.server.rpc.Service;
import org.epics.pvdata.factory.FieldFactory;
import org.epics.pvdata.factory.PVDataFactory;
import org.epics.pvdata.pv.FieldCreate;
import org.epics.pvdata.pv.PVDataCreate;
import org.epics.pvdata.pv.PVLong;
import org.epics.pvdata.pv.PVScalarArray;
import org.epics.pvdata.pv.PVString;
import org.epics.pvdata.pv.PVStringArray;
import org.epics.pvdata.pv.PVStructure;
import org.epics.pvdata.pv.ScalarType;
import org.epics.pvdata.pv.Status;
import org.epics.pvdata.pv.Structure;
import org.epics.pvdatabase.PVDatabase;
import org.epics.pvdatabase.PVDatabaseFactory;
import org.epics.pvdatabase.PVRecord;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;


public class ModelRPCService extends PVRecord implements RPCService {

    private static final Logger LOGGER = LoggerFactory.getLogger(ModelRPCService.class);
    private static final String RPC_NAME = "model";

    private PVStructure pvTop;
    private Model model;
    private long request_counter = 0;

    public static ModelRPCService create(Model model) {
        FieldCreate fieldCreate = FieldFactory.getFieldCreate();
        PVDataCreate pvDataCreate = PVDataFactory.getPVDataCreate();
        NTScalarBuilder ntScalarBuilder = NTScalar.createBuilder();

        PVStructure pvTop = ntScalarBuilder
                .value(ScalarType.pvLong)
                .add("res", fieldCreate.createScalarArray(ScalarType.pvString))
                .addTimeStamp()
                .createPVStructure();

        ModelRPCService pvRecord = new ModelRPCService(RPC_NAME, pvTop, model);
        //PVDatabase master = PVDatabaseFactory.getMaster();
        //master.addRecord(pvRecord);        
        return pvRecord;
    }

    /**
     * Get the service
     */
    public Service getService(PVStructure pvRequest) {
        return this;
    }

    /**
     * Process a request from the client
     *
     * @param args The request from the client
     * @return The result.
     */
    public PVStructure request(PVStructure args) throws RPCRequestException {
        PVString pvOp = args.getSubField(PVString.class, "op");
        PVString pvQuery = args.getSubField(PVString.class, "query");
        if (pvOp == null) {
            throw new RPCRequestException(Status.StatusType.ERROR,
                    "PVString field with name 'op' expected.");
        }
        if (pvQuery == null) {
            throw new RPCRequestException(Status.StatusType.ERROR,
                    "PVString field with name 'query' expected.");
        }

        return execute(pvOp, pvQuery);

    }

    private ModelRPCService(String recordName, PVStructure pvStructure, Model model) {
        super(recordName, pvStructure);
        this.pvTop = pvStructure;
        this.model = model;
    }

    @Override
    public void process() {
        super.process(); //To change body of generated methods, choose Tools | Templates.
        request_counter++;
        PVLong pvCounter = pvTop.getSubField(PVLong.class, "value");
        pvCounter.put(request_counter);
    }

    private PVStructure execute(PVString op, PVString query) {
        lock();
        beginGroupPut();
        PVStringArray pvres = (PVStringArray) pvTop.getScalarArrayField("res", ScalarType.pvString);
        if ("getuuid".equals(op.get())) {
            if ("root".equals(query.get())) {
                String[] struuids = new String[1];
                struuids[0] = model.getEnterpriseUUID().toString();
                pvres.setCapacity(1);
                pvres.put(0, struuids.length, struuids, 0);
            } else {
                try {
                    UUID uuid = UUID.fromString(query.get());
                    UUID[] uuids = model.getAreas(uuid);

                    pvres.setCapacity(uuids.length);
                    String[] struuids = new String[uuids.length];
                    int i = 0;
                    for (UUID id : uuids) {
                        struuids[i] = id.toString();
                        i++;
                    }
                    pvres.put(0, struuids.length, struuids, 0);
                } catch (Exception ex) {

                }
            }
        } else if ("getid".equals(op.get())) {
            UUID uuid = model.getAreaUUID(query.get());
            pvres.setCapacity(1);
            String[] struuids = new String[1];
            if (uuid != null) {
                struuids[0] = uuid.toString();
            };
            pvres.put(0, struuids.length, struuids, 0);
        }
        if ("getdesc".equals(op.get())) {

        }
        process();
        endGroupPut();
        PVStructure pvResult = PVDataFactory.getPVDataCreate().createPVStructure(pvTop);
        unlock();
        return pvResult;
    }

}

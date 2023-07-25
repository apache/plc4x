/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package org.apache.plc4x.hop.transforms.plc4xevent;

import com.fasterxml.jackson.databind.ObjectMapper;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;
import java.util.concurrent.ConcurrentLinkedQueue;
import java.util.concurrent.locks.ReentrantLock;
import org.apache.hop.core.CheckResult;
import org.apache.hop.core.Const;
import org.apache.hop.core.ICheckResult;
import org.apache.hop.core.RowMetaAndData;
import org.apache.hop.core.exception.HopException;
import org.apache.hop.core.exception.HopPluginException;
import org.apache.hop.core.logging.LogLevel;
import org.apache.hop.core.row.IRowMeta;
import org.apache.hop.core.row.IValueMeta;
import org.apache.hop.core.row.RowDataUtil;
import org.apache.hop.core.row.RowMeta;
import org.apache.hop.core.row.value.ValueMetaFactory;
import org.apache.hop.core.util.StringUtil;
import org.apache.hop.i18n.BaseMessages;
import org.apache.hop.metadata.api.IHopMetadataProvider;
import org.apache.hop.pipeline.IExecutionStoppedListener;
import org.apache.hop.pipeline.Pipeline;
import org.apache.hop.pipeline.PipelineMeta;
import org.apache.hop.pipeline.engine.IPipelineEngine;
import org.apache.hop.pipeline.transform.BaseTransform;
import org.apache.hop.pipeline.transform.TransformMeta;
import org.apache.plc4x.hop.metadata.Plc4xConnection;
import org.apache.plc4x.hop.metadata.util.Plc4xLookup;
import org.apache.plc4x.hop.metadata.util.Plc4xWrapperConnection;
import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.listener.ConnectionStateListener;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.s7.events.S7Event;
import org.apache.plc4x.java.s7.events.S7ModeEvent;
import org.apache.plc4x.java.s7.events.S7SysEvent;
import org.apache.plc4x.java.s7.events.S7UserEvent;
import org.openide.util.Lookup;

/**
 * This transform receives an event from the S7 driver, of type MODE, SYS, USR
 * or ALM. Only one type can be processed at a time.
 *
 */
public class Plc4xEvent extends BaseTransform<Plc4xEventMeta, Plc4xEventData> implements ConnectionStateListener {

    public static String FIELD_MODE_EVENT = "MODE";
    public static String FIELD_USER_EVENT = "USR";
    public static String FIELD_SYS_EVENT = "SYS";
    public static String FIELD_ALARM_EVENT = "ALM";

    private static final Class<?> PKG = Plc4xEvent.class; // Needed by Translator

    private Plc4xConnection connmeta = null;
    private Plc4xWrapperConnection connwrapper = null;
    private PlcConsumerRegistration registerMode = null;
    private PlcConsumerRegistration registerUser = null;
    private PlcConsumerRegistration registerSys = null;
    private PlcConsumerRegistration registerAlarm = null;
    private PlcSubscriptionRequest subsRequest = null;
    private PlcSubscriptionResponse subresponse = null;

    private List<ICheckResult> remarks = new ArrayList<>(); // stores the errors...

    private Plc4xLookup lookup = Plc4xLookup.getDefault();
    private Lookup.Template template = null;
    private Lookup.Result<Plc4xWrapperConnection> result = null;


    /*
    * The transfer of events is done from the driver tasks. A delay can be added 
    * for the execution of this transformer.
     */
    private ObjectMapper mapper = new ObjectMapper();
    private ConcurrentLinkedQueue<S7Event> events = new ConcurrentLinkedQueue();
    private boolean stopBundle = false;
    private int index = 0;

    private static final ReentrantLock lock = new ReentrantLock();

    private static final String dummy = "dummy";

    public Plc4xEvent(TransformMeta transformMeta, Plc4xEventMeta meta, Plc4xEventData data, int copyNr, PipelineMeta pipelineMeta,
            Pipeline pipeline) {
        super(transformMeta, meta, data, copyNr, pipelineMeta, pipeline);
    }

    /*
    * Including Date and Time field for every row 
    *
    * @param meta Meta data from user dialog
    * @param remarks Error registers
    * @param origin transform instance name
     */
    public static final RowMetaAndData buildRow(Plc4xEventMeta meta,
            List<ICheckResult> remarks,
            String origin) throws HopPluginException {

        IRowMeta rowMeta = new RowMeta();
        Object[] rowData = RowDataUtil.allocateRowData(2);
        int index = 0;

        ArrayList<String> fields = new ArrayList<String>();

        if (meta.isModeEvent()) {
            for (S7ModeEvent.Fields field : S7ModeEvent.Fields.values()) {
                fields.add(field.name());
            }
        } else if (meta.isSysEvent()) {
            for (S7SysEvent.Fields field : S7SysEvent.Fields.values()) {
                fields.add(field.name());
            }
        } else if (meta.isUserEvent()) {
            for (S7UserEvent.Fields field : S7UserEvent.Fields.values()) {
                fields.add(field.name());
            }
        } else if (meta.isAlarmEvent()) {
            fields.add("ALARM");
        }

        for (String field : fields) {
            IValueMeta valueMeta
                    = ValueMetaFactory.createValueMeta(field, IValueMeta.TYPE_STRING); // build a  
            rowData[index] = StringUtil.EMPTY_STRING;
            // Now add value to the row!
            // This is in fact a copy from the fields row, but now with data.
            rowMeta.addValueMeta(valueMeta);
            index++;
        }

        return new RowMetaAndData(rowMeta, rowData);
    }

    /* 
  * 1. Block the other instances by means of a lock.  
  * 2. Try to locate an existing connection.
  * 3. If it doesn't exist, it tries to take control of the routine to 
  *    create an instance of PlcConnection and his wrapper.
  * 4. Register the connection wrapper for global access.
  * 5. If the connection to the PLC is made, then it creates the query 
  *    and executes it.
  *
     */
    @Override
    public boolean processRow() throws HopException {
        Object[] r = getRow(); // Get row from input rowset & set row busy!
        setLogLevel(LogLevel.DEBUG);

        if ((!meta.isNeverEnding() && data.rowsWritten >= data.rowLimit) && !isStopped()) {
            setOutputDone(); // signal end to receiver(s)
            return false;
        }

        // If we do not have the jobId, it is because we have not completed 
        // the CYC subscription.
        if (null == subsRequest) {
            RegisterPlcTags();
            GetSubscriptions();
            RegisterEventHandler();
        }

        if (!events.isEmpty()) {
            S7Event s7event = events.poll();
            index = 0;
            r = data.outputRowMeta.cloneRow(data.outputRowData);
            for (String name : data.outputRowMeta.getFieldNames()) {
                if (null != s7event.getMap().get(name)) {
                    if ("MAP".equals(name)) {
                        try {              
                            //jackson 2.15.x Does not support recursive maps. 
                            s7event.getMap().remove("MAP");
                            r[index++] = mapper.writer()
                                            .writeValueAsString(s7event.getMap());
                        } catch (Exception ex) {
                            logError(ex.getMessage());
                        }
                    } else {
                        r[index++] = s7event.getMap().get(name).toString();
                    };
                } else {
                    try {
                        //                r[index++] = mapper.writer()
                        //                                .writeValueAsString(s7event.getMap());
                        r[index++] = null;
                    } catch (Exception ex) {
                        logError(ex.getMessage());
                    }
                }
            }

            data.prevDate = data.rowDate;
            data.rowDate = new Date();

            putRow(data.outputRowMeta, r); // return your data
            data.rowsWritten++;
        }

        if (isStopped()) {
            setStopped(true);
            return false;
        }

        return true;
    }

    @Override
    public boolean init() {
        try {
            if (super.init()) {
                // Determine the number of rows to generate...
                data.rowLimit = Const.toLong(resolve(meta.getRowLimit()), -1L);

                if (data.rowLimit < 0L) { // Unable to parse
                    logError(BaseMessages.getString(PKG, "Plc4x.Read.Meta.Wrong.RowLimit.Number"));
                    return false; // fail
                }

                // Create a row (constants) with all the values in it...
                List<ICheckResult> remarks = new ArrayList<>(); // stores the errors...
                RowMetaAndData outputRow = buildRow(meta, remarks, getTransformName());
                if (!remarks.isEmpty()) {
                    for (int i = 0; i < remarks.size(); i++) {
                        CheckResult cr = (CheckResult) remarks.get(i);
                        logError(cr.getText());
                    }
                    return false;
                }

                data.outputRowData = outputRow.getData();
                data.outputRowMeta = outputRow.getRowMeta();

                mapper.findAndRegisterModules();

                getPlcConnection();

                return true;
            }

            return false;

        } catch (Exception ex) {
            setErrors(1L);
            logError("Error initializing transform", ex);
            return false;
        }
    }

    /*
    * Here, must perform the cleaning of any resource, main of the connection to 
    * the associated PLC.
     */
    @Override
    public void cleanup() {
        super.cleanup();
        logBasic("Cleanup. Release connection.");
        if (null != connwrapper) {
            if (null != registerMode) {
                registerMode.unregister();
            }
            if (null != registerUser) {
                registerUser.unregister();
            }
            if (null != registerSys) {
                registerSys.unregister();
            }
            if (null != registerAlarm) {
                registerAlarm.unregister();
            }
            connwrapper.release();
            if (connwrapper.refCnt() <= 0) {
                lookup.remove(connwrapper);
            }
        }
    }

    /*
    * Here, must perform the cleaning of any resource. 
    * 1. Check if we have reference to wrapper.
    * 2. Release de reference to object.
    * 3. The lastone remove the global reference to connection wrapper.
    * 4. Clear local references.
     */
    @Override
    public void dispose() {
        super.dispose();
        if (null != registerMode) {
            registerMode.unregister();
        }
        if (null != registerUser) {
            registerUser.unregister();
        }
        if (null != registerSys) {
            registerSys.unregister();
        }
        if (null != registerAlarm) {
            registerAlarm.unregister();
        }

        if (connwrapper != null) {
            logBasic("Dispose. Release connection: " + connwrapper.refCnt());
            connwrapper.release();
            if (connwrapper.refCnt() <= 0) {
                lookup.remove(connwrapper);
            }
            connwrapper = null;
            subsRequest = null;
            registerMode = null;
            registerUser = null;
            registerSys = null;
            registerAlarm = null;

        }
    }

    public void stopRunning() throws HopException {
        super.stopRunning();
        stopBundle = true;
    }

    private void getPlcConnection() {
        lock.lock(); //(01)
        try {

            IHopMetadataProvider metaprovider = getMetadataProvider();
            connmeta = metaprovider.getSerializer(Plc4xConnection.class).load(meta.getConnection());

            if (connwrapper == null) {
                template = new Lookup.Template<>(Plc4xWrapperConnection.class, meta.getConnection(), null);
                result = lookup.lookup(template);
                if (!result.allItems().isEmpty()) {
                    logBasic("Using connection: " + meta.getConnection());
                    connwrapper = (Plc4xWrapperConnection) result.allInstances().toArray()[0];
                    if (connwrapper != null) {
                        connwrapper.retain();
                    }
                }
            };

            if (connmeta == null) {
                logError(
                        BaseMessages.getString(
                                PKG,
                                "Plc4x.Read.Meta.Log.SetMetadata",
                                meta.getConnection()));
            }

            if ((connmeta != null) && (connwrapper == null)) {
                subsRequest = null;
                try {
                    final PlcConnection conn = new DefaultPlcDriverManager().getConnection(connmeta.getUrl()); //(03)
                    Thread.sleep(200);
                    if (conn.isConnected()) {
                        logBasic("Create new connection with url : " + connmeta.getUrl());
                        connwrapper = new Plc4xWrapperConnection(conn, meta.getConnection());
                        lookup.add(connwrapper);
                    }

                } catch (Exception ex) {
                    setErrors(1L);
                    logError("Unable to create connection to PLC. " + ex.getMessage());
                }
            }

        } catch (HopException ex) {
            logError("Unable to create connection to PLC. " + ex.getMessage());
        } finally {
            lock.unlock();
        }
    }

    /*
    * Registers the tags for the cyclical subscription.
    * In the first processing of the rows, a check of the tags is 
    * carried out in order that they are well formed, generating an exception 
    * if they are not.
     */
    public boolean RegisterPlcTags() {
        if ((connmeta != null) && (connwrapper != null)) {
            if (connwrapper.getConnection().isConnected()) {
                if (subsRequest == null) {
                    final PlcSubscriptionRequest.Builder subscription = connwrapper.getConnection().subscriptionRequestBuilder(); //(05)
                    if (meta.isModeEvent()) {
                        subscription.addEventTagAddress(FIELD_MODE_EVENT, FIELD_MODE_EVENT);
                    }
                    if (meta.isUserEvent()) {
                        subscription.addEventTagAddress(FIELD_USER_EVENT, FIELD_USER_EVENT);
                    }
                    if (meta.isSysEvent()) {
                        subscription.addEventTagAddress(FIELD_SYS_EVENT, FIELD_SYS_EVENT);
                    }
                    if (meta.isAlarmEvent()) {
                        subscription.addEventTagAddress(FIELD_ALARM_EVENT, FIELD_ALARM_EVENT);
                    }

                    subsRequest = subscription.build();

                }

            } else {
                setErrors(1L);
                logError("PLC is not connected.");
                setOutputDone();
                return false;
            }

        } else {
            setErrors(1L);
            logError("PLC connection don't exist.");
            setOutputDone();
            return false;
        }
        return true;
    }

    /*
    * This method performs the subscription to the events within the PLC.
    * 1. Take the waiting time as a reference to monitor the request.
    * 2. Captures the jobId associated with the subscription 
    *    which is assigned by the PLC.
     */
    public void GetSubscriptions() {
        try {
            subresponse = subsRequest.execute().get();
        } catch (Exception ex) {
            setErrors(1L);
            logError("Unable read from PLC. " + ex.getMessage());
        }
    }

    /*
    * Subscribe to incoming events by transferring them to
    * the local event buffer.
     */
    public void RegisterEventHandler() {

        if (meta.isModeEvent() && subresponse.getResponseCode(FIELD_MODE_EVENT) == PlcResponseCode.OK) {
            registerMode
                    = subresponse
                            .getSubscriptionHandle(FIELD_MODE_EVENT)
                            .register(msg -> {
                                events.add((S7Event) msg);
                            });
        }

        if (meta.isUserEvent() && subresponse.getResponseCode(FIELD_USER_EVENT) == PlcResponseCode.OK) {
            registerUser
                    = subresponse
                            .getSubscriptionHandle(FIELD_USER_EVENT)
                            .register(msg -> {
                                events.add((S7Event) msg);
                            });
        }

        if (meta.isSysEvent() && subresponse.getResponseCode(FIELD_SYS_EVENT) == PlcResponseCode.OK) {
            registerSys
                    = subresponse
                            .getSubscriptionHandle(FIELD_SYS_EVENT)
                            .register(msg -> {
                                events.add((S7Event) msg);
                            });
        }

        if (meta.isAlarmEvent() && subresponse.getResponseCode(FIELD_ALARM_EVENT) == PlcResponseCode.OK) {
            registerAlarm
                    = subresponse
                            .getSubscriptionHandle(FIELD_ALARM_EVENT)
                            .register(msg -> {
                                events.add((S7Event) msg);
                            });
        }
    }

    /*
    * Unsubscribes to the event generator from the driver.
    * TODO: It remains to verify the unsubscription in line. 
    *       At the moment only with the disconnection of the PLC.
     */
    public void UnRegisterCYCHandler() {
        if (null != registerMode) {
            registerMode.unregister();
        }
        if (null != registerUser) {
            registerUser.unregister();
        }
        if (null != registerSys) {
            registerSys.unregister();
        }
        if (null != registerAlarm) {
            registerAlarm.unregister();
        }
        registerMode = null;
        registerUser = null;
        registerSys = null;
        registerAlarm = null;
        subsRequest = null;
    }

    /*
    * When the driver makes a connection this method is called. 
    * In the case of the S7 driver, only the connection of the embedded channel
    * is served.
     */
    @Override
    public void connected() {
        if (connwrapper.getConnection().isConnected()) {
            RegisterPlcTags();
            GetSubscriptions();
            RegisterEventHandler();
        }
    }

    /*
    * This method is called from the driver to indicate that it has 
    * disconnected from the PLC.
    * If the PLC is disconnected, it automatically deletes all 
    * the status associated with this connection.
     */
    @Override
    public void disconnected() {
        logError("Driver disconnected!");
        UnRegisterCYCHandler();
    }

}

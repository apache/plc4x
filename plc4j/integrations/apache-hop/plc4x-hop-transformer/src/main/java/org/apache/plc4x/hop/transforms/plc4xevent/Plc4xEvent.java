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

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ConcurrentLinkedQueue;
import java.util.concurrent.LinkedBlockingQueue;
import java.util.concurrent.locks.ReentrantLock;
import java.util.logging.Level;
import java.util.logging.Logger;
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
import org.apache.hop.pipeline.Pipeline;
import org.apache.hop.pipeline.PipelineMeta;
import org.apache.hop.pipeline.transform.BaseTransform;
import org.apache.hop.pipeline.transform.ITransform;
import org.apache.hop.pipeline.transform.TransformMeta;
import org.apache.plc4x.hop.metadata.Plc4xConnection;
import org.apache.plc4x.hop.metadata.util.Plc4xLookup;
import org.apache.plc4x.hop.transforms.util.Plc4xGeneratorField;
import org.apache.plc4x.hop.metadata.util.Plc4xWrapperConnection;
import org.apache.plc4x.hop.transforms.plc4xinput.Plc4xRead;
import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.s7.events.S7AlarmEvent;
import org.apache.plc4x.java.s7.events.S7Event;
import org.apache.plc4x.java.s7.events.S7ModeEvent;
import org.apache.plc4x.java.s7.events.S7ModeEvent.Fields;
import org.apache.plc4x.java.s7.events.S7SysEvent;
import org.apache.plc4x.java.s7.events.S7UserEvent;
import org.apache.plc4x.java.s7.readwrite.ModeTransitionType;
import org.openide.util.Lookup;

/**
 * This transform receives an event from the S7 driver, of type MODE, SYS, 
 * USR or ALM. Only one type can be processed at a time.
 *
 */
public class Plc4xEvent extends BaseTransform<Plc4xEventMeta, Plc4xEventData> {

    public static String FIELD_MODE_EVENT = "MODE";
    public static String FIELD_USER_EVENT = "USR";
    public static String FIELD_SYS_EVENT = "SYS";
    public static String FIELD_ALARM_EVENT = "ALM";


    private static final Class<?> PKG = Plc4xEvent.class; // Needed by Translator

    private Plc4xConnection connmeta              = null;
    private Plc4xWrapperConnection connwrapper    = null;
    private PlcReadRequest readRequest            = null;
    private PlcConsumerRegistration registerMode  = null;
    private PlcConsumerRegistration registerUser  = null;
    private PlcConsumerRegistration registerSys   = null;
    private PlcConsumerRegistration registerAlarm   = null;
    private PlcSubscriptionRequest subsbuild      = null;
    private PlcSubscriptionResponse subresponse   = null;

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
                Pipeline pipeline ) {
    super( transformMeta, meta, data, copyNr, pipelineMeta, pipeline );
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
            for (S7ModeEvent.Fields field:S7ModeEvent.Fields.values()) {
                fields.add(field.name());
            }
        } else if (meta.isSysEvent()) {
            for (S7SysEvent.Fields field:S7SysEvent.Fields.values()) {
                fields.add(field.name());
            }        
        } else if (meta.isUserEvent()) {
            for (S7UserEvent.Fields field:S7UserEvent.Fields.values()) {
                fields.add(field.name());
            }        
        } else if (meta.isAlarmEvent()) {
            fields.add("ALARM");    
        }

        for (String field : fields) {
            IValueMeta valueMeta =
                ValueMetaFactory.createValueMeta(field, IValueMeta.TYPE_STRING); // build a  
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
  public  boolean processRow() throws HopException {
    Object[] r = getRow(); // Get row from input rowset & set row busy!
    setLogLevel(LogLevel.DEBUG);
    
    if ((!meta.isNeverEnding() && data.rowsWritten >= data.rowLimit) && !isStopped()) {   
      setOutputDone(); // signal end to receiver(s)
      return false;        
    }    
    
    lock.lock(); //(01)
    try {
        IHopMetadataProvider metaprovider = getMetadataProvider();
        connmeta = metaprovider.getSerializer(Plc4xConnection.class).load(meta.getConnection());
        if (connwrapper == null) {
            connwrapper = (Plc4xWrapperConnection) getPipeline().getExtensionDataMap().get(meta.getConnection()); //(02)
            if (connwrapper != null) connwrapper.retain();
        };

        if (connmeta == null){    
            logError(
                BaseMessages.getString(
                    PKG,
                    "Plc4x.Read.Meta.Log.SetMetadata",
                    meta.getConnection()));         
        }

        if ((connmeta != null) && (connwrapper == null)){
            readRequest = null;
            try{
                PlcConnection conn =  new DefaultPlcDriverManager().getConnection(connmeta.getUrl()); //(03)
                if (conn.isConnected()) {
                    connwrapper = new Plc4xWrapperConnection(conn, meta.getConnection());            
                    getPipeline().getExtensionDataMap().put(meta.getConnection(), connwrapper); //(04)
                }
            } catch (Exception ex){
                setErrors(1L);
                logError("Unable to create connection to PLC. " + ex.getMessage());
            }
        }
    } finally {
        lock.unlock();
    }
    
    if ((connmeta != null) && (connwrapper != null)){
        if (connwrapper.getConnection().isConnected()){
            if (readRequest == null){
                final PlcSubscriptionRequest.Builder subscription  = connwrapper.getConnection().subscriptionRequestBuilder(); //(05)
                if (meta.isModeEvent())  subscription.addEventTagAddress(FIELD_MODE_EVENT, FIELD_MODE_EVENT);
                if (meta.isUserEvent())  subscription.addEventTagAddress(FIELD_USER_EVENT, FIELD_USER_EVENT);
                if (meta.isSysEvent())   subscription.addEventTagAddress(FIELD_SYS_EVENT, FIELD_SYS_EVENT);
                if (meta.isAlarmEvent()) subscription.addEventTagAddress(FIELD_ALARM_EVENT, FIELD_ALARM_EVENT);
                
                subsbuild  = subscription.build();

            }            
            try {    
                 subresponse = subsbuild.execute().get();    
                 
                if (meta.isModeEvent() && subresponse.getResponseCode(FIELD_MODE_EVENT) == PlcResponseCode.OK) {
                    registerMode = 
                        subresponse
                        .getSubscriptionHandle(FIELD_MODE_EVENT)
                        .register(msg -> {
                            events.add((S7Event) msg);
                        });
                }    
                
                if (meta.isUserEvent() && subresponse.getResponseCode(FIELD_USER_EVENT) == PlcResponseCode.OK) {
                    registerUser = 
                        subresponse
                        .getSubscriptionHandle(FIELD_USER_EVENT)
                        .register(msg -> {
                            events.add((S7Event) msg);                            
                        });
                } 

                if (meta.isSysEvent() && subresponse.getResponseCode(FIELD_SYS_EVENT) == PlcResponseCode.OK) {
                    registerSys = 
                        subresponse
                        .getSubscriptionHandle(FIELD_SYS_EVENT)
                        .register(msg -> {
                            events.add((S7Event) msg);                            
                        });
                } 

                if (meta.isAlarmEvent() && subresponse.getResponseCode(FIELD_ALARM_EVENT) == PlcResponseCode.OK) {
                    registerAlarm = 
                        subresponse
                        .getSubscriptionHandle(FIELD_ALARM_EVENT)
                        .register(msg -> {
                            events.add((S7Event) msg);                            
                        });
                }                 

            } catch (Exception ex) {
                setErrors(1L);                
                logError("Unable read from PLC. " + ex.getMessage());
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

   while (events.size() == 0) {
        try {          
            Thread.sleep(100);
            if (stopBundle) {
                setOutputDone(); // signal end to receiver(s)
                return false; 
            }
        } catch (InterruptedException ex) {
            break;
        }
   } 
     
    S7Event s7event = events.poll();
    index = 0;
    r = data.outputRowMeta.cloneRow(data.outputRowData);     
    for (String name:data.outputRowMeta.getFieldNames()) {
        System.out.println(name + ": " + s7event.getMap().get(name));
        if (null != s7event.getMap().get(name)) {
            r[index++] = s7event.getMap().get(name).toString();
        } else {
            try {
                r[index++] = mapper.writer()
                                .writeValueAsString(s7event.getMap());
            } catch (Exception ex) {
                Logger.getLogger(Plc4xEvent.class.getName()).log(Level.SEVERE, null, ex);
            }
        }
    }
    

    data.prevDate = data.rowDate;
    data.rowDate = new Date();    

    putRow(data.outputRowMeta, r ); // return your data
    data.rowsWritten++;
    return true;
  }


    @Override
    public boolean init() {
        try {
            if(super.init()){     
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
            
        } catch (Exception ex){
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
        System.out.println("*************** CLEANUP *****************");          
        super.cleanup();
        logBasic("Cleanup. Release connection.");
        if (null != connwrapper) {
            if (null != registerMode )  registerMode.unregister();
            if (null != registerUser)   registerUser.unregister();
            if (null != registerSys)    registerSys.unregister();
            if (null != registerAlarm)  registerAlarm.unregister();    
            connwrapper.release();            
            if (connwrapper.refCnt() <= 0) 
                lookup.remove(connwrapper); 
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
        System.out.println("*************** DISPOSE *****************");          
        super.dispose();
        if (connwrapper != null) {
            logBasic("Dispose. Release connection: " + connwrapper.refCnt());            
            connwrapper.release();
            if (connwrapper.refCnt() <= 0) 
                lookup.remove(connwrapper);            
            connwrapper = null;
            readRequest     = null;
            registerMode    = null;
            registerUser    = null;
            registerSys     = null;
            registerAlarm   = null;

        }
    }

    @Override
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
                    System.out.println("Aqui encontro la conexion: " + meta.getConnection());
                    connwrapper = (Plc4xWrapperConnection) result.allInstances().toArray()[0];
                    if (connwrapper != null) connwrapper.retain();
                }
            };

            if (connmeta == null){    
                logError(
                    BaseMessages.getString(
                        PKG,
                        "Plc4x.Read.Meta.Log.SetMetadata",
                        meta.getConnection()));         
            }

            if ((connmeta != null) && (connwrapper == null)){
                readRequest = null;
                try{
                    System.out.println("Creo una nueva conexión...");
                    PlcConnection conn =  new DefaultPlcDriverManager().getConnection(connmeta.getUrl()); //(03)

                    if (conn.isConnected()) {
                        System.out.println("**** Agrego la segunda conexion. ****");
                        connwrapper = new Plc4xWrapperConnection(conn, meta.getConnection());            
                        lookup.add(connwrapper);
                    }

                } catch (Exception ex){
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
  
  
}

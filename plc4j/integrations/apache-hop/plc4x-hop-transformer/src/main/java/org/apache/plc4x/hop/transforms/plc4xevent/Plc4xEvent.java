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

import java.util.ArrayList;
import java.util.Date;
import java.util.List;
import java.util.Map;
import java.util.concurrent.locks.ReentrantLock;
import org.apache.hop.core.CheckResult;
import org.apache.hop.core.Const;
import org.apache.hop.core.ICheckResult;
import org.apache.hop.core.RowMetaAndData;
import org.apache.hop.core.exception.HopException;
import org.apache.hop.core.exception.HopPluginException;
import org.apache.hop.core.logging.LogLevel;
import org.apache.hop.core.row.IRowMeta;
import org.apache.hop.core.row.RowDataUtil;
import org.apache.hop.core.row.RowMeta;
import org.apache.hop.i18n.BaseMessages;
import org.apache.hop.metadata.api.IHopMetadataProvider;
import org.apache.hop.pipeline.Pipeline;
import org.apache.hop.pipeline.PipelineMeta;
import org.apache.hop.pipeline.transform.BaseTransform;
import org.apache.hop.pipeline.transform.ITransform;
import org.apache.hop.pipeline.transform.TransformMeta;
import org.apache.plc4x.hop.metadata.Plc4xConnection;
import org.apache.plc4x.hop.transforms.util.Plc4xWrapperConnection;
import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.s7.events.S7ModeEvent;
import org.apache.plc4x.java.s7.readwrite.ModeTransitionType;

/**
 * Transform That contains the basic skeleton needed to create your own plugin
 *
 */
public class Plc4xEvent extends BaseTransform<Plc4xEventMeta, Plc4xEventData> implements ITransform {

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
  public static final RowMetaAndData buildRow(
    Plc4xEventMeta meta, List<ICheckResult> remarks, String origin) throws HopPluginException {
    IRowMeta rowMeta = new RowMeta();
    Object[] rowData = RowDataUtil.allocateRowData(2);
    int index = 0;

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
                    connwrapper = new Plc4xWrapperConnection(conn);            
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
                            System.out.println("******** S7ModeEvent ********");
                            Map<String, Object> map = ((S7ModeEvent) msg).getMap();
                            map.forEach((x, y) -> { 
                                System.out.println(x + " : " + y);
                            });
                            short currentmode = (short) 
                                    map.get(S7ModeEvent.Fields.CURRENT_MODE.name());
                            System.out.println("CURRENT_MODE MSG: " + ModeTransitionType.enumForValue(currentmode).name());
                            System.out.println("****************************");
                        });
                }    
                
                if (meta.isUserEvent() && subresponse.getResponseCode(FIELD_USER_EVENT) == PlcResponseCode.OK) {
                    registerUser = 
                        subresponse
                        .getSubscriptionHandle(FIELD_USER_EVENT)
                        .register(msg -> {
                            System.out.println("******** S7ModeEvent ********");
                            Map<String, Object> map = ((S7ModeEvent) msg).getMap();
                            map.forEach((x, y) -> { 
                                System.out.println(x + " : " + y);
                            });
                            short currentmode = (short) 
                                    map.get(S7ModeEvent.Fields.CURRENT_MODE.name());
                            System.out.println("CURRENT_MODE MSG: " + ModeTransitionType.enumForValue(currentmode).name());
                            System.out.println("****************************");
                        });
                } 

                if (meta.isSysEvent() && subresponse.getResponseCode(FIELD_SYS_EVENT) == PlcResponseCode.OK) {
                    registerSys = 
                        subresponse
                        .getSubscriptionHandle(FIELD_SYS_EVENT)
                        .register(msg -> {
                            System.out.println("******** S7ModeEvent ********");
                            Map<String, Object> map = ((S7ModeEvent) msg).getMap();
                            map.forEach((x, y) -> { 
                                System.out.println(x + " : " + y);
                            });
                            short currentmode = (short) 
                                    map.get(S7ModeEvent.Fields.CURRENT_MODE.name());
                            System.out.println("CURRENT_MODE MSG: " + ModeTransitionType.enumForValue(currentmode).name());
                            System.out.println("****************************");
                        });
                } 

                if (meta.isAlarmEvent() && subresponse.getResponseCode(FIELD_ALARM_EVENT) == PlcResponseCode.OK) {
                    registerAlarm = 
                        subresponse
                        .getSubscriptionHandle(FIELD_ALARM_EVENT)
                        .register(msg -> {
                            System.out.println("******** S7ModeEvent ********");
                            Map<String, Object> map = ((S7ModeEvent) msg).getMap();
                            map.forEach((x, y) -> { 
                                System.out.println(x + " : " + y);
                            });
                            short currentmode = (short) 
                                    map.get(S7ModeEvent.Fields.CURRENT_MODE.name());
                            System.out.println("CURRENT_MODE MSG: " + ModeTransitionType.enumForValue(currentmode).name());
                            System.out.println("****************************");
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

    
    r = data.outputRowMeta.cloneRow(data.outputRowData); 
    logBasic("Tamano de los datos: " + r.length);
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
        super.cleanup();
        logBasic("Cleanup. Release connection.");
        if (connwrapper != null)
        connwrapper.release();     
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
        if (connwrapper != null) {
            logBasic("Dispose. Release connection: " + connwrapper.refCnt());            
            connwrapper.release();   
            if (!connwrapper.getConnection().isConnected()){           
                getPipeline().getExtensionDataMap().remove(meta.getConnection());
            }            
            connwrapper = null;
            readRequest = null;

        }
    }
 
  
  
  
}

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

import org.apache.plc4x.merlot.api.DriverEvent;
import org.apache.plc4x.merlot.api.core.Merlot;
import org.apache.plc4x.merlot.das.base.api.BaseDriver;
import org.apache.plc4x.merlot.db.api.DBCollector;
import org.apache.plc4x.merlot.scheduler.api.JobContext;
import java.util.Collection;
import java.util.Iterator;
import java.util.Set;
import java.util.Vector;
import java.util.concurrent.ConcurrentHashMap;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.epics.pvdata.copy.CreateRequest;
import org.epics.pvdata.copy.PVCopy;
import org.epics.pvdata.copy.PVCopyFactory;
import org.epics.pvdata.property.PVTimeStamp;
import org.epics.pvdata.property.PVTimeStampFactory;
import org.epics.pvdata.property.TimeStamp;
import org.epics.pvdata.property.TimeStampFactory;
import org.epics.pvdata.pv.BooleanArrayData;
import org.epics.pvdata.pv.ByteArrayData;
import org.epics.pvdata.pv.DoubleArrayData;
import org.epics.pvdata.pv.FloatArrayData;
import org.epics.pvdata.pv.IntArrayData;
import org.epics.pvdata.pv.LongArrayData;
import org.epics.pvdata.pv.PVBoolean;
import org.epics.pvdata.pv.PVBooleanArray;
import org.epics.pvdata.pv.PVByte;
import org.epics.pvdata.pv.PVByteArray;
import org.epics.pvdata.pv.PVDouble;
import org.epics.pvdata.pv.PVDoubleArray;
import org.epics.pvdata.pv.PVField;
import org.epics.pvdata.pv.PVFloat;
import org.epics.pvdata.pv.PVFloatArray;
import org.epics.pvdata.pv.PVInt;
import org.epics.pvdata.pv.PVIntArray;
import org.epics.pvdata.pv.PVLong;
import org.epics.pvdata.pv.PVLongArray;
import org.epics.pvdata.pv.PVScalar;
import org.epics.pvdata.pv.PVScalarArray;
import org.epics.pvdata.pv.PVShort;
import org.epics.pvdata.pv.PVShortArray;
import org.epics.pvdata.pv.PVString;
import org.epics.pvdata.pv.PVStringArray;
import org.epics.pvdata.pv.PVStructure;
import org.epics.pvdata.pv.PVUByte;
import org.epics.pvdata.pv.PVUInt;
import org.epics.pvdata.pv.PVULong;
import org.epics.pvdata.pv.PVUShort;
import org.epics.pvdata.pv.ShortArrayData;
import org.epics.pvdata.pv.StringArrayData;
import org.epics.pvdata.util.pvDataHelper.GetHelper;
import org.epics.pvdatabase.PVListener;
import org.epics.pvdatabase.PVRecord;
import org.epics.pvdatabase.PVRecordField;
import org.epics.pvdatabase.PVRecordStructure;
import org.osgi.framework.BundleContext;
import org.osgi.framework.ServiceReference;
import org.osgi.service.device.Driver;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;


//TODO: add the logger system, the write is INFO, the read maybe DEBUG
//TODO: Read/Write information for every tag for statical use.
public class DBCollectorImpl implements DBCollector, PVListener {
    private static final Logger LOGGER = LoggerFactory.getLogger(DBCollectorImpl.class);  
    private final BundleContext bundleContext;
    private final String device;
    private final BaseDriver devDriver;
    private final PlcConnection plcConn;
    private PlcReadRequest.Builder readRequestBuilder = null;
    private PlcWriteRequest.Builder writeRequestBuilder = null; 
    private PlcReadRequest readRequest = null;
    private PlcWriteRequest writeRequest = null;
    
    private final PVTimeStamp pvTimeStamp = PVTimeStampFactory.create();
    private final TimeStamp timeStamp = TimeStampFactory.create();
    private CreateRequest createRequest = CreateRequest.create();
    private String request = "field(timeStamp,value)";
    private PVStructure pvRequest = createRequest.createRequest(request);
    
    boolean run = false;
    private ConcurrentHashMap<String, PVRecord> items = new ConcurrentHashMap<String, PVRecord>();

    public DBCollectorImpl(BundleContext bundleContext, String device) {
        this.bundleContext = bundleContext;
        this.device = device;
        this.devDriver = getDriver(device);
        this.plcConn = devDriver.getPlcConnection();
        if (plcConn == null){LOGGER.info("DBCollectorImpl. Connection is  null");}
    }
    
    @Override
    public void execute(JobContext ctx) {
        if (run) {
            DriverEvent event = devDriver.getEvent();
            event.setFunctionCode(Merlot.FUNCTION.FC_READ_MEMORY_BYTES);
            event.setPlcReadRequest(readRequest);
            event.setCallback(this); 
            devDriver.putEvent(event);
        };
    }

    //TODO: Manejo por función solicitada
    //TODO: Replace item with Master database instance
    @Override
    public void execute(DriverEvent cb) {
        PVRecord pvRecord = null;
        PVField value = null;
        PVBoolean enable = null;
        boolean disable = false;
        if (cb.getPlcReadResponse() != null) {           
            PlcReadResponse drvResponse = cb.getPlcReadResponse();
            PlcReadRequest drvRequest = drvResponse.getRequest();
            
            Set<String> strFields = drvRequest.getTagNames();
            
            for(String strField:strFields) {

                pvRecord = items.get(strField);   
                value =  pvRecord.getPVStructure().getSubField("value");

                enable = pvRecord.getPVStructure().getBooleanField("scan_enable");
                
                disable = false;
                if (enable == null){
                    LOGGER.info("There is no \"enable\" field in the PVRecord. Could not process.");
                    disable = true;
                } else if (!enable.get()){
                    LOGGER.info("PVRecord[" + pvRecord.getRecordName()+"] is disabled.");
                    disable = true;
                }
                
                if (!disable) {
                    switch(value.getField().getType()) {
                        case scalar:
                                readScalar(drvResponse, pvRecord);
                            break;
                        case scalarArray:
                                readScalarArray(drvResponse, pvRecord);                    
                            break;
                        case structure:
                                executeStructure(drvResponse, pvRecord);                     
                            break;
                        case structureArray:
                                executeStructureArray(drvResponse, pvRecord);                     
                            break;
                        case union:
                                executeUnion(drvResponse, pvRecord);                     
                            break;
                        case unionArray:
                                executeUnionArray(drvResponse, pvRecord); 
                            break;                
                    }
                }
            }
        } else if (cb.getPlcWriteResponse() != null){
            LOGGER.info("Write ready...");
        }                
    }
    
    @Override
    public void attach(Collection<PVRecord> pvRecords) {
        run = false; 
    
        if (plcConn == null){LOGGER.info("Conexion null");}
        if (plcConn.isConnected()){      
            readRequestBuilder = this.plcConn.readRequestBuilder(); 
            for (PVRecord pvRecord:pvRecords){
                items.put(pvRecord.getRecordName(), pvRecord);
                readRequestBuilder.addTagAddress(pvRecord.getRecordName(), pvRecord.getPVStructure().getStringField("id").get());
                PVCopy pvCopy = PVCopyFactory.create(pvRecord.getPVRecordStructure().getPVStructure(), pvRequest, "");
                pvRecord.addListener(this, pvCopy);
            }
            readRequest = readRequestBuilder.build();
          } else {
            LOGGER.info("Driver not connected...");
        }
        run = true;
    }


    @Override
    public void detach(Collection<PVRecord> pvRecords) {
        throw new UnsupportedOperationException("Not supported yet.");
    }
    
    private BaseDriver getDriver(String device){
        try{
            String filterdriver =  "(DRIVER_ID=" + device + ")"; 
            ServiceReference[] refdrvs = bundleContext.getAllServiceReferences(Driver.class.getName(), filterdriver);
            BaseDriver refdrv = (BaseDriver) bundleContext.getService(refdrvs[0]);
            return refdrv;            
        } catch (Exception ex){
            LOGGER.info(ex.toString());
        }
        return null;
    }

    
    @Override
    public void dataPut(PVRecordField pvRecordField) {
    }

    @Override
    public void dataPut(PVRecordStructure pvStructure, PVRecordField pvRecordField) {   
    }

    @Override
    public void beginGroupPut(PVRecord pvRecord) {        
    }

    //TODO: Manejo de diferentes tipos    
    @Override
    public void endGroupPut(PVRecord pvRecord) {
        PVField pvField =  pvRecord.getPVStructure().getSubField("value");
        switch(pvField.getField().getType()){
            case scalar:
                    writeScalar(pvRecord);
                break;
            case scalarArray:
                    writeScalarArray(pvRecord);                    
                break;
            case structure:
;                     
                break;
            case structureArray:
                    
                break;
            case union:
                  
                break;
            case unionArray:
 
                break;                            
        }
    }

    @Override
    public void unlisten(PVRecord arg0) {
        LOGGER.info("Unlisten: ");
    }
    
    private void readScalar(PlcReadResponse response, PVRecord pvRecord){
        PVScalar value = (PVScalar) pvRecord.getPVStructure().getSubField("value");
        pvRecord.lock();
        try{
            switch(value.getScalar().getScalarType()){
                case pvBoolean:
                    ((PVBoolean) value).put(response.getBoolean(pvRecord.getRecordName()));
                    break;
                case pvByte:
                    ((PVByte) value).put(response.getByte(pvRecord.getRecordName()));
                    break;
                case pvShort:
                    ((PVShort) value).put(response.getShort(pvRecord.getRecordName()));
                    break;
                case pvInt:
                    ((PVInt) value).put(response.getInteger(pvRecord.getRecordName()));
                    break;
                case pvLong:
                    ((PVLong) value).put(response.getInteger(pvRecord.getRecordName()));                    
                    break;
                case pvUByte:
                    ((PVUByte) value).put(response.getByte(pvRecord.getRecordName()));                    
                    break;
                case pvUShort:
                    ((PVUShort) value).put(response.getShort(pvRecord.getRecordName()));                    
                    break;
                case pvUInt:
                    ((PVUInt) value).put(response.getInteger(pvRecord.getRecordName()));                    
                    break;
                case pvULong:
                    ((PVULong) value).put(response.getLong(pvRecord.getRecordName()));                    
                    break;
                case pvFloat:
                    ((PVFloat) value).put(response.getFloat(pvRecord.getRecordName()));                    
                    break;
                case pvDouble:
                    ((PVDouble) value).put(response.getDouble(pvRecord.getRecordName()));                    
                    break;
                case pvString:
                    ((PVString) value).put(response.getString(pvRecord.getRecordName()));                    
                    break;
                default:
                    LOGGER.info("Unsupported type.");
                    break;
            }
            pvRecord.process();
        } catch(Exception ex){
            LOGGER.info("executeScalar: " + ex.getMessage());
        } finally{
            pvRecord.unlock();
        }
    }
    
    private void readScalarArray(PlcReadResponse response, PVRecord pvRecord){
        PVScalarArray value = (PVScalarArray) pvRecord.getPVStructure().getSubField("value");        
        pvRecord.lock();
        try{
            switch(value.getScalarArray().getElementType()){
                case pvBoolean:{
                    Collection<Boolean> resObjects = response.getAllBooleans(pvRecord.getRecordName()); 
                    boolean[] arrBooleans = new boolean[resObjects.size()];
                    Iterator<Boolean> itr = resObjects.iterator();
                    Boolean objValue;
                    int i = 0;
                    while(itr.hasNext()){
                        objValue = itr.next();
                        if (objValue == null) break; 
                        arrBooleans[i] = objValue;
                        i++;
                    }
                    ((PVBooleanArray) value).put(0,resObjects.size(),arrBooleans,0);
                }
                    break;
                    
                case pvByte:
                case pvUByte:
                {
                    Collection<Byte> resObjects = response.getAllBytes(pvRecord.getRecordName()); 
                    byte[] arrBytes = new byte[resObjects.size()];
                    Iterator<Byte> itr = resObjects.iterator();
                    Byte objValue;
                    int i = 0;
                    while(itr.hasNext()){
                        objValue = itr.next();
                        if (objValue == null) break; 
                        arrBytes[i] = objValue;
                        i++;
                    }                    
                    ((PVByteArray) value).put(0,resObjects.size(),arrBytes,0);
                }
                    break;
                case pvShort: 
                case pvUShort:
                {
                    Collection<Short> resObjects = response.getAllShorts(pvRecord.getRecordName()); 
                    short[] arrShorts = new short[resObjects.size()];
                    Iterator<Short> itr = resObjects.iterator();
                    Short objValue;
                    int i = 0;
                    while(itr.hasNext()){
                        objValue = itr.next();
                        if (objValue == null) break; 
                        arrShorts[i] = objValue;
                        i++;
                    }     
                    int capacity = ((PVShortArray) value).getCapacity();
                    ((PVShortArray) value).put(0,capacity,arrShorts,0);
                }
                    break;
                case pvInt: 
                case pvUInt:                
                {
                    Collection<Integer> resObjects = response.getAllIntegers(pvRecord.getRecordName()); 
                    int[] arrIntegers = new int[resObjects.size()];
                    Iterator<Integer> itr = resObjects.iterator();
                    Integer objValue;
                    int i = 0;
                    while(itr.hasNext()){
                        objValue = itr.next();
                        if (objValue == null) break; 
                        arrIntegers[i] = objValue;
                        i++;
                    }     
                    int capacity = ((PVIntArray) value).getCapacity();
                    ((PVIntArray) value).put(0,capacity,arrIntegers,0);
                }
                    break;
                case pvLong: 
                case pvULong:                    
                {
                    Collection<Long> resObjects = response.getAllLongs(pvRecord.getRecordName()); 
                    long[] arrLongs = new long[resObjects.size()];
                    Iterator<Long> itr = resObjects.iterator();
                    Long objValue;
                    int i = 0;
                    while(itr.hasNext()){
                        objValue = itr.next();
                        if (objValue == null) break; 
                        arrLongs[i] = objValue;
                        i++;
                    }           
                    int capacity = ((PVLongArray) value).getCapacity();
                    ((PVLongArray) value).put(0,capacity,arrLongs,0);
                }               
                    break;
                case pvFloat:
                {
                    Collection<Float> resObjects = response.getAllFloats(pvRecord.getRecordName()); 
                    float[] arrFloats = new float[resObjects.size()];
                    Iterator<Float> itr = resObjects.iterator();
                    Float objValue;
                    int i = 0;
                    while(itr.hasNext()){
                        objValue = itr.next();
                        if (objValue == null) break;                     
                        arrFloats[i] = objValue;
                        i++;
                    }     
                    //System.out.println("Tamaño i: " + i);
                    int capacity = ((PVFloatArray) value).getCapacity();
                    ((PVFloatArray) value).put(0,capacity,arrFloats,0);
                }                                 
                    break;
                case pvDouble:
                {
                    Collection<Double> resObjects = response.getAllDoubles(pvRecord.getRecordName()); 
                    double[] arrDoubles = new double[resObjects.size()];
                    Iterator<Double> itr = resObjects.iterator();
                    Double objValue;
                    int i = 0;
                    while(itr.hasNext()){
                        objValue = itr.next();
                        if (objValue == null) break; 
                        arrDoubles[i] = objValue;
                        i++;
                    }         
                    int capacity = ((PVDoubleArray) value).getCapacity();
                    ((PVDoubleArray) value).put(0,resObjects.size(),arrDoubles,0);
                }                                         
                    break;
                case pvString:  
                {
                    Collection<String> resObjects = response.getAllStrings(pvRecord.getRecordName()); 
                    String[] arrStrings = new String[resObjects.size()];
                    Iterator<String> itr = resObjects.iterator();
                    String objValue;
                    int i = 0;
                    while(itr.hasNext()){
                        objValue = itr.next();
                        if (objValue == null) break; 
                        arrStrings[i] = objValue;
                        i++;
                    }    
                    int capacity = ((PVStringArray) value).getCapacity();
                    ((PVStringArray) value).put(0,resObjects.size(),arrStrings,0);
                }                      
                    break;
                default:
                        LOGGER.info("Unsupported type.");
                    break;
            }
            pvRecord.process();
        } catch (Exception ex){
            ex.printStackTrace();
            LOGGER.info(ex.getMessage());
        } finally{
            pvRecord.unlock();
        }        
    }
    
    private void executeStructure(PlcReadResponse response, PVRecord pvRecord){
        
    }
    
    private void executeStructureArray(PlcReadResponse response, PVRecord pvRecord){
        
    }   
    
    private void executeUnion(PlcReadResponse response, PVRecord pvRecord){
        
    }     
    
    private void executeUnionArray(PlcReadResponse response, PVRecord pvRecord){
        
    }     
    
    private void writeScalar(PVRecord pvRecord){
        try {    
            pvRecord.lock();            
            PVScalar valueScalar = (PVScalar) pvRecord.getPVStructure().getSubField("value");
            Object value = null;
            switch(valueScalar.getScalar().getScalarType()){
                case pvBoolean:
                    value = pvRecord.getPVStructure().getBooleanField("value").get();
                    break;                   
                case pvByte:
                case pvUByte:                     
                    value = pvRecord.getPVStructure().getByteField("value").get();
                    break;
                case pvShort:
                case pvUShort:                    
                    value = pvRecord.getPVStructure().getShortField("value").get();
                    break;
                case pvInt:
                case pvUInt:                    
                    value = pvRecord.getPVStructure().getIntField("value").get();
                    break;
                case pvLong:
                case pvULong:                    
                    value = pvRecord.getPVStructure().getLongField("value").get();
                    break;
                case pvFloat:
                    value = pvRecord.getPVStructure().getFloatField("value").get();                   
                    break;
                case pvDouble:
                    value = pvRecord.getPVStructure().getDoubleField("value").get();                   
                    break;
                case pvString:
                    value = pvRecord.getPVStructure().getStringField("value").get();                      
                    break;
                default:
                    LOGGER.info("Unsupported type.");
                    break;
            }
                                         
            writeRequestBuilder = plcConn.writeRequestBuilder();
            writeRequestBuilder.addTagAddress(pvRecord.getRecordName(), 
                    pvRecord.getPVStructure().getStringField("id").get(), 
                    value);
            writeRequest = writeRequestBuilder.build();
            
            DriverEvent event = devDriver.getEvent();
            event.setFunctionCode(Merlot.FUNCTION.FC_WRITE_DATA_BYTES);
            event.setPlcWriteRequest(writeRequest);
            event.setCallback(this); 
            devDriver.putEvent(event);  
        } catch (Exception ex){
            LOGGER.info("Exception ex: " + ex);
        }  finally{
            pvRecord.unlock();
        }  
        
    }  
    
    private void writeScalarArray(PVRecord pvRecord){
        try{
            pvRecord.lock();
            PVScalarArray valueScalarArray = (PVScalarArray) pvRecord.getPVStructure().getSubField("value");  
            Object value = null;
            switch(valueScalarArray.getScalarArray().getElementType()){
                case pvBoolean:
                    value = getBooleanVector((PVBooleanArray) valueScalarArray); 
                    break;                   
                case pvByte:
                case pvUByte:                     
                    value = getByteVector((PVByteArray) valueScalarArray); 
                    break;
                case pvShort:
                case pvUShort:                    
                    value = getShortVector((PVShortArray) valueScalarArray); 
                    break;
                case pvInt:
                case pvUInt:                    
                    value = getIntVector((PVIntArray) valueScalarArray);
                    break;
                case pvLong:
                case pvULong:                    
                    value = getLongVector((PVLongArray) valueScalarArray);
                    break;
                case pvFloat:
                    value = getFloatVector((PVFloatArray) valueScalarArray);                   
                    break;
                case pvDouble:
                    value = getDoubleVector((PVDoubleArray) valueScalarArray);                  
                    break;
                case pvString:
                    value = GetHelper.getStringVector((PVStringArray) valueScalarArray);                       
                    break;
                default:
                    LOGGER.info("Unsupported type.");
                    break;
            }
                                         
            writeRequestBuilder = plcConn.writeRequestBuilder();
            writeRequestBuilder.addTagAddress(pvRecord.getRecordName(), 
                    pvRecord.getPVStructure().getStringField("id").get(), 
                    ((Vector)value).toArray());
            writeRequest = writeRequestBuilder.build();
            
            DriverEvent event = devDriver.getEvent();
            event.setFunctionCode(Merlot.FUNCTION.FC_WRITE_DATA_BYTES);
            event.setPlcWriteRequest(writeRequest);
            event.setCallback(this); 
            devDriver.putEvent(event);              
        } catch (Exception ex){
            LOGGER.info("writeScalarArray: " + ex.toString());
        } finally {
            pvRecord.unlock();            
        }
    }
        
    public static Vector<Boolean> getBooleanVector( PVBooleanArray pv )
    {
        int len = pv.getLength();
        // double[] storage = new double[len];
        Vector<Boolean> ret = new Vector<>();
        BooleanArrayData data = new BooleanArrayData();
        int offset = 0;
        while(offset < len) {
            int num = pv.get(offset,(len-offset),data);
            for (int i=0; i<num; i++) ret.add(data.data[offset+i]);
            // System.arraycopy(data.data,data.offset,storage,offset,num);
            offset += num;
        }
        return ret;
    }        
    
    public static Vector<Byte> getByteVector( PVByteArray pv )
    {
        int len = pv.getLength();
        // double[] storage = new double[len];
        Vector<Byte> ret = new Vector<>();
        ByteArrayData data = new ByteArrayData();
        int offset = 0;
        while(offset < len) {
            int num = pv.get(offset,(len-offset),data);
            for (int i=0; i<num; i++) ret.add(data.data[offset+i]);
            // System.arraycopy(data.data,data.offset,storage,offset,num);
            offset += num;
        }
        return ret;
    }         
    
    public static Vector<Short> getShortVector( PVShortArray pv )
    {
        int len = pv.getLength();
        // double[] storage = new double[len];
        Vector<Short> ret = new Vector<>();
        ShortArrayData data = new ShortArrayData();
        int offset = 0;
        while(offset < len) {
            int num = pv.get(offset,(len-offset),data);
            for (int i=0; i<num; i++) ret.add(data.data[offset+i]);
            // System.arraycopy(data.data,data.offset,storage,offset,num);
            offset += num;
        }
        return ret;
    }    
    
    public static Vector<Integer> getIntVector( PVIntArray pv )
    {
        int len = pv.getLength();
        // double[] storage = new double[len];
        Vector<Integer> ret = new Vector<>();
        IntArrayData data = new IntArrayData();
        int offset = 0;
        while(offset < len) {
            int num = pv.get(offset,(len-offset),data);
            for (int i=0; i<num; i++) ret.add(data.data[offset+i]);
            // System.arraycopy(data.data,data.offset,storage,offset,num);
            offset += num;
        }
        return ret;
    }   
    
    public static Vector<Long> getLongVector( PVLongArray pv )
    {
        int len = pv.getLength();
        // double[] storage = new double[len];
        Vector<Long> ret = new Vector<>();
        LongArrayData data = new LongArrayData();
        int offset = 0;
        while(offset < len) {
            int num = pv.get(offset,(len-offset),data);
            for (int i=0; i<num; i++) ret.add(data.data[offset+i]);
            // System.arraycopy(data.data,data.offset,storage,offset,num);
            offset += num;
        }
        return ret;
    }        
        

    public static Vector<Float> getFloatVector( PVFloatArray pv )
    {
        int len = pv.getLength();
        // double[] storage = new double[len];
        Vector<Float> ret = new Vector<>();
        FloatArrayData data = new FloatArrayData();
        int offset = 0;
        while(offset < len) {
            int num = pv.get(offset,(len-offset),data);
            for (int i=0; i<num; i++) ret.add(data.data[offset+i]);
            // System.arraycopy(data.data,data.offset,storage,offset,num);
            offset += num;
        }
        return ret;
    }
    
    public static Vector<Double> getDoubleVector( PVDoubleArray pv )
    {
        int len = pv.getLength();
        // double[] storage = new double[len];
        Vector<Double> ret = new Vector<>();
        DoubleArrayData data = new DoubleArrayData();
        int offset = 0;
        while(offset < len) {
            int num = pv.get(offset,(len-offset),data);
            for (int i=0; i<num; i++) ret.add(data.data[offset+i]);
            // System.arraycopy(data.data,data.offset,storage,offset,num);
            offset += num;
        }
        return ret;
    } 
    
    public static Vector<String> getStringVector( PVStringArray pv )
    {
        int len = pv.getLength();
        // double[] storage = new double[len];
        Vector<String> ret = new Vector<>();
        StringArrayData data = new StringArrayData();
        int offset = 0;
        while(offset < len) {
            int num = pv.get(offset,(len-offset),data);
            for (int i=0; i<num; i++) ret.add(data.data[offset+i]);
            // System.arraycopy(data.data,data.offset,storage,offset,num);
            offset += num;
        }
        return ret;
    }      
    
}

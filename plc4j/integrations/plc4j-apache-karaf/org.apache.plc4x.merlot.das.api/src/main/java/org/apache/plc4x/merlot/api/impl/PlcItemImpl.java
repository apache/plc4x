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
package org.apache.plc4x.merlot.api.impl;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import java.time.LocalDate;
import java.time.ZoneOffset;
import java.util.Collection;
import java.util.Date;
import java.util.Hashtable;
import java.util.LinkedList;
import java.util.UUID;
import java.util.concurrent.locks.ReentrantLock;
import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.merlot.api.PlcItem;
import org.apache.plc4x.merlot.api.PlcItemListener;

/*
*
*/
public class PlcItemImpl implements PlcItem {
        
    private ReentrantLock lock = new ReentrantLock();    
    private String item_name;
    private String item_description;
    private String item_id;
    private UUID item_uid;
    
    private boolean item_enable = false;
   
    private int item_accessrigths = 0;
    
    private Boolean item_isarray     = false;    
    private Boolean item_disableoutput = false; 
    
    private final Hashtable<String, Object> myproperties;    
    
    private PlcTag item_plctag = null;
    private PlcValue item_plcvalue = null;
   
    private LinkedList<PlcItemListener> item_clients = null;
    private ByteBuf item_buffer = null;
    
    
    private long itemtransmit = 0;
    private long itemreceives = 0;    
    private long itemerrors = 0;


    private PlcReadResponse  plcresponse;   
    private PlcResponseCode dataquality;

    private Date lastreadtime;
    private Date lastwritetime;
    private Date lasterrortime;

 
    public PlcItemImpl(PlcItemBuilder builder) {
        this.myproperties = new Hashtable<>();        
        
        myproperties.put(PlcItem.ITEM_NAME, builder.item_name);        
        item_name = builder.item_name;
        
        myproperties.put(PlcItem.ITEM_DESCRIPTION, builder.item_description);         
        item_description = builder.item_description;
        
        myproperties.put(PlcItem.ITEM_UID, builder.item_uid.toString());           
        item_uid = builder.item_uid;        
        
        myproperties.put(PlcItem.ITEM_ID, builder.item_id);           
        item_id = builder.item_id;        

        
        item_enable = builder.item_enable;
        item_accessrigths = builder.item_accessrigths;
        
        item_isarray = builder.item_isarray;
        item_disableoutput = builder.item_disableoutput;
        
        item_buffer = Unpooled.buffer();
        item_clients = new LinkedList<>();
    }     

    @Override
    public void enable() {
        item_enable = true;
    }

    @Override
    public void disable() {
        item_enable = false;
    }
               
    @Override
    public UUID getItemUid() {
        return UUID.fromString((String) myproperties.get(PlcItem.ITEM_UID));
    }
       
    @Override
    public String getItemName() {
        return (String) myproperties.get(PlcItem.ITEM_NAME);
    }

    @Override
    public void setItemName(String itemname) {
        myproperties.put(PlcItem.ITEM_NAME, item_name);   
    }

    @Override
    public String getItemDescription() {
        return (String) myproperties.get(PlcItem.ITEM_DESCRIPTION);
    }

    @Override
    public void setItemDescription(String item_description) {
        myproperties.put(PlcItem.ITEM_DESCRIPTION, item_description);  
    }

    @Override
    public String getItemId() {
        return (String) myproperties.get(PlcItem.ITEM_ID);
    }

    @Override
    public void setItemId(String item_id) {
        myproperties.put(PlcItem.ITEM_ID, item_id);
    }

    @Override
    public PlcTag getItemPlcTag() {
        return item_plctag;
    }

    @Override
    public void setItemPlcTag(PlcTag itemplctag) {
        this.item_plctag = itemplctag;
    }
        
    @Override
    public Boolean getEnable() {
        return item_enable;
    }

    @Override
    public void setEnable(Boolean enable) {
        this.item_enable = enable;
                
    }    

    @Override
    public Boolean isEnable() {
        return  item_enable;
    }
                
    @Override
    public Boolean getIsArray() {
        return item_isarray;
    }

    @Override
    public void setIsArray(Boolean isArray) {
        this.item_isarray = isArray;
    }

    @Override
    public Boolean getIsDisableOutput() {
        return item_disableoutput;
    }

    @Override
    public void setIsDisableOutput(Boolean item_disableoutput) {
        this.item_disableoutput = item_disableoutput;
    }

    @Override
    public long getItemTransmits() {
        return itemtransmit;
    }

    @Override
    public long getItemReceives() {
        return itemreceives;
    }

    public long getItemErrors() {
        return itemerrors;
    }


    @Override
    public int getAccessRights() {
        return item_accessrigths;
    }

    @Override
    public void setAccessRights(int item_accessrigths) {
        this.item_accessrigths = item_accessrigths;
    }

    @Override
    public PlcResponseCode getDataQuality() {
        return dataquality;
    }

    @Override
    public void setPlcValue(PlcValue  plcvalue) {
        
        lock.lock();        
        try {
            this.item_plcvalue = plcvalue;
            item_buffer.clear();
            switch (item_plctag.getPlcValueType()){
                case BYTE:
                    item_plcvalue.getList().forEach(p -> {item_buffer.writeByte(p.getByte());});
                    break;
                case CHAR:                
                case STRING:
                    plcresponse.getAllStrings(item_uid.toString()).forEach(s -> {
                        item_buffer.writeBytes(s.getBytes());
                    });  
                    item_plcvalue.getList().forEach(p -> {item_buffer.writeBytes(p.getRaw());});                
                    break;
                case WORD:
                case USINT:
                case SINT:
                case UINT:
                case INT:
                case DINT:
                    item_plcvalue.getList().forEach(p -> {item_buffer.writeShort(p.getShort());});                 
                    break;
                case UDINT:
                case ULINT:
                case LINT:    
                    item_plcvalue.getList().forEach(p -> {item_buffer.writeLong(p.getLong());});                 
                    break;
                case BOOL:     
                    item_plcvalue.getList().forEach(p -> {item_buffer.writeBoolean(p.getBoolean());});                 
                    break;
                case REAL:
                case LREAL: 
                    item_plcvalue.getList().forEach(p -> {item_buffer.writeFloat(p.getFloat());});                  
                    break;
                case DATE_AND_TIME:  
                    item_plcvalue.getList().forEach(p -> {item_buffer.writeLong(p.getDateTime().toEpochSecond(ZoneOffset.UTC));});                  
                    break;
                case DATE: 
                    plcresponse.getAllDates(item_uid.toString()).forEach(dt -> {
                        item_buffer.writeLong(dt.toEpochDay());
                    }); 
                    item_plcvalue.getList().forEach(p -> {item_buffer.writeLong(p.getDateTime().toLocalDate().toEpochDay());});                 
                    break;
                case TIME:
                    break;
                case TIME_OF_DAY:
                    item_plcvalue.getList().forEach(p -> {item_buffer.writeLong(p.getTime().toEpochSecond(LocalDate.MAX, ZoneOffset.UTC));});                 
                    break;               
                default:
                    throw new NotImplementedException("The response type for datatype " + item_plcvalue.getPlcValueType() + " is not yet implemented");                
            }
        } catch (Exception ex) {
            
        } finally {
            lock.unlock();
        }
    }

    @Override
    public PlcValue getItemPlcValue() {
        lock.lock();
        PlcValue plcvalue;
        try {
            plcvalue = item_plcvalue;
        } finally {
            lock.unlock();
        }
        return plcvalue;
    }

    @Override
    public ByteBuf getItemByteBuf() {
        lock.lock();
        ByteBuf itembuffer;
        try {
            itembuffer = item_buffer.duplicate();
        } finally {
            lock.unlock();
        }        
        return itembuffer;
    }

    @Override
    public byte[] getItemBytes() {
        lock.lock();
        byte[] bytebuffer;
        try {
            bytebuffer = item_buffer.array();
        } finally {
            lock.unlock();
        }              
        return bytebuffer;
    }
                  
    @Override
    public void setDataQuality(PlcResponseCode dataquality) {
        this.dataquality = dataquality;
    }

    @Override
    public void addItemClient(PlcItemListener client) {
        if (!item_clients.contains(client)) {
            client.atach(this);
            item_clients.add(client);
        }
    }

    @Override
    public void removeItemClient(PlcItemListener client) {
        if (!item_clients.contains(client)) {
            item_clients.remove(client);            
            client.detach();
        }        
    }
            
    @Override
    public Hashtable<String, Object> getProperties() {
        return myproperties;
    }
               
    @Override
    public Date getLastReadDate() {
        return lastreadtime;
    }

    @Override
    public Date getLastWriteDate() {
        return lastwritetime;
    }

    @Override
    public Date getLastErrorDate() {
        return lasterrortime;
    }
    
    private void updateClients(){
        item_clients.forEach(c -> c.update());
    }

    public static class PlcItemBuilder {
        private final String item_name;
        private  UUID item_uid;    
        private String item_description;
        private String item_id;
        private Boolean item_enable         = false;   
        private int item_accessrigths       = 0;    
        private Boolean item_isarray        = false; 
        private Boolean item_disableoutput = false;         

        public PlcItemBuilder(String item_name) {
            this.item_name = item_name;
            this.item_uid = UUID.randomUUID();
        }

        public PlcItemBuilder setItemUid(UUID item_uid) {
            this.item_uid = item_uid;
            return this;
        }

        public PlcItemBuilder setItemDescription(String item_description) {
            this.item_description = item_description;
            return this;            
        }

        public PlcItemBuilder setItemId(String item_id) {
            this.item_id = item_id;
            return this;            
        }

        public PlcItemBuilder setItemEnable(boolean item_enable) {
            this.item_enable = item_enable;
            return this;            
        }

        public PlcItemBuilder setItemAccessrigths(int item_accessrigths) {
            this.item_accessrigths = item_accessrigths;
            return this;            
        }

        public PlcItemBuilder setItemIsarray(Boolean item_isarray) {
            this.item_isarray = item_isarray;
            return this;            
        }

        public PlcItemBuilder setItemDisableoutput(Boolean item_disableoutput) {
            this.item_disableoutput = item_disableoutput;
            return this;            
        }

        public PlcItem build() {
            PlcItem plcitem = new PlcItemImpl(this);
            validatePlcItemObject(plcitem);
            return plcitem;
        }
        
        private void validatePlcItemObject(PlcItem plcitem) {
            //
        }            
        
        
    }        
    
}

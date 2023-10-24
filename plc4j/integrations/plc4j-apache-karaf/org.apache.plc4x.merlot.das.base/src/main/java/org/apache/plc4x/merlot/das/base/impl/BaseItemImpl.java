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
package org.apache.plc4x.merlot.das.base.impl;

import io.netty.buffer.ByteBuf;
import java.util.Collection;
import java.util.Date;
import java.util.Hashtable;
import java.util.UUID;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.merlot.api.PlcItem;
import org.apache.plc4x.merlot.api.PlcItemListener;


public class BaseItemImpl implements PlcItem {

    private UUID uid;
    
    private String  itemname;
    private String  itemdescription;
    
    private String id;
    
    private boolean enable = false;
    
    private boolean isarray = false;
    private int arraystart = 0;
    
    private boolean disableoutput = true; 
    
    private long itemtransmits = 0;
    private long itemreceives = 0;
    private long itemerrors = 0;
    
    private int accessrights = 0;
    private int dataquality = 0; 

    private Date lastreaddate;
    private Date lastwritedate;
    private Date lasterrordate;

    @Override
    public void setPlcValue(PlcValue plcvalue) {
        throw new UnsupportedOperationException("Not supported yet."); // Generated from nbfs://nbhost/SystemFileSystem/Templates/Classes/Code/GeneratedMethodBody
    }

    @Override
    public PlcValue getItemPlcValue() {
        throw new UnsupportedOperationException("Not supported yet."); // Generated from nbfs://nbhost/SystemFileSystem/Templates/Classes/Code/GeneratedMethodBody
    }

    @Override
    public ByteBuf getItemByteBuf() {
        throw new UnsupportedOperationException("Not supported yet."); // Generated from nbfs://nbhost/SystemFileSystem/Templates/Classes/Code/GeneratedMethodBody
    }

    @Override
    public byte[] getItemBytes() {
        throw new UnsupportedOperationException("Not supported yet."); // Generated from nbfs://nbhost/SystemFileSystem/Templates/Classes/Code/GeneratedMethodBody
    }



    
    @Override
    public UUID getItemUid() {
        return uid; 
    }

    @Override
    public String getItemName() {
        return itemname; 
    }

    @Override
    public void setItemName(String itemname) {
        this.itemname = itemname;
    }

    @Override
    public String getItemDescription() {
        return itemdescription; 
    }

    @Override
    public void setItemDescription(String itemdescription) {
        this.itemdescription = itemdescription;
    }

    @Override
    public String getItemId() {
        return id; 
    }

    @Override
    public void setItemId(String itemid) {
        this.id = itemid; 
    }

    @Override
    public Boolean getEnable() {
        return enable; 
    }

    @Override
    public void setEnable(Boolean enable) {
        this.enable = enable; 
    }

    @Override
    public Boolean isEnable() {
        return enable;
    }

    
    @Override
    public Boolean getIsArray() {
        return isarray; 
    }

    @Override
    public void setIsArray(Boolean isarray) {
        this.isarray = isarray;
    }


    @Override
    public Boolean getIsDisableOutput() {
        return disableoutput; 
    }

    @Override
    public void setIsDisableOutput(Boolean isDisableOutput) {
        this.disableoutput = isDisableOutput;
    }

    @Override
    public long getItemTransmits() {
        return itemtransmits; 
    }

    @Override
    public long getItemReceives() {
        return itemreceives; 
    }

    @Override
    public long getItemErrors() {
        return itemerrors; 
    }

    @Override
    public int getAccessRights() {
        return accessrights; 
    }

    @Override
    public void setAccessRights(int accessrigths) {
        this.accessrights = accessrigths;
    }


    @Override
    public Date getLastReadDate() {
        return lastreaddate; 
    }

    @Override
    public Date getLastWriteDate() {
        return lastwritedate; 
    }

    @Override
    public Date getLastErrorDate() {
        return lasterrordate; 
    }

    @Override
    public PlcTag getItemPlcTag() {
        throw new UnsupportedOperationException("Not supported yet."); // Generated from nbfs://nbhost/SystemFileSystem/Templates/Classes/Code/GeneratedMethodBody
    }

    @Override
    public void setItemPlcTag(PlcTag itemplctag) {
        throw new UnsupportedOperationException("Not supported yet."); // Generated from nbfs://nbhost/SystemFileSystem/Templates/Classes/Code/GeneratedMethodBody
    }

    @Override
    public PlcResponseCode getDataQuality() {
        throw new UnsupportedOperationException("Not supported yet."); // Generated from nbfs://nbhost/SystemFileSystem/Templates/Classes/Code/GeneratedMethodBody
    }

    @Override
    public void setDataQuality(PlcResponseCode dataquality) {
        throw new UnsupportedOperationException("Not supported yet."); // Generated from nbfs://nbhost/SystemFileSystem/Templates/Classes/Code/GeneratedMethodBody
    }

    @Override
    public void addItemClient(PlcItemListener client) {
        throw new UnsupportedOperationException("Not supported yet."); // Generated from nbfs://nbhost/SystemFileSystem/Templates/Classes/Code/GeneratedMethodBody
    }

    @Override
    public void removeItemClient(PlcItemListener client) {
        throw new UnsupportedOperationException("Not supported yet."); // Generated from nbfs://nbhost/SystemFileSystem/Templates/Classes/Code/GeneratedMethodBody
    }

    @Override
    public void enable() {
        throw new UnsupportedOperationException("Not supported yet."); // Generated from nbfs://nbhost/SystemFileSystem/Templates/Classes/Code/GeneratedMethodBody
    }

    @Override
    public void disable() {
        throw new UnsupportedOperationException("Not supported yet."); // Generated from nbfs://nbhost/SystemFileSystem/Templates/Classes/Code/GeneratedMethodBody
    }

    @Override
    public Hashtable<String, Object> getProperties() {
        throw new UnsupportedOperationException("Not supported yet."); // Generated from nbfs://nbhost/SystemFileSystem/Templates/Classes/Code/GeneratedMethodBody
    }
    
}

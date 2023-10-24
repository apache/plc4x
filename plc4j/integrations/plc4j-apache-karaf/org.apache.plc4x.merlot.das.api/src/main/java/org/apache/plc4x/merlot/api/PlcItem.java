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
package org.apache.plc4x.merlot.api;

import io.netty.buffer.ByteBuf;
import java.util.Collection;
import java.util.Date;
import java.util.Hashtable;
import java.util.UUID; 
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;

public interface PlcItem {
    
    /*
    *
    */    
    public static final String	ITEM_UID            = "plc4x.item.uid";    
    
    public static final String	ITEM_NAME           = "plc4x.item.name";   
    
    public static final String	ITEM_DESCRIPTION    = "plc4x.item.desc"; 

    public static final String	ITEM_ID             = "plc4x.item.id";      
    
    /*
    *
    */
    public void enable();
    
    /*
    *
    */
    public void disable(); 
    
    
    /*
    *
    */    
    public UUID getItemUid();
    
    /*
    *
    */    
    public String getItemName();
    
    /*
    *
    */    
    public void setItemName(String itemname);    
    
    /*
    *
    */    
    public String getItemDescription();
    
    /*
    *
    */    
    public void setItemDescription(String itemdescription);    
    
    /*
    *
    */    
    public String getItemId();
    
    /*
    *
    */    
    public void setItemId(String itemid);  
    
    /*
    *
    */    
    public PlcTag getItemPlcTag();
    
    /*
    *
    */    
    public void setItemPlcTag(PlcTag itemplctag);    
    
    /*
    *
    */    
    public Boolean getEnable();
    
    /*
    *
    */    
    public void setEnable(Boolean enable);
    
        /*
    *
    */    
    public Boolean isEnable();
    
    /*
    *
    */    
    public Boolean getIsArray();
    
    /*
    *
    */    
    public void setIsArray(Boolean isArray);    
      
    /*
    *
    */    
    public Boolean getIsDisableOutput();
    
    /*
    *
    */    
    public void setIsDisableOutput(Boolean isDisableOutput);    
    
    /*
    *
    */    
    public long getItemTransmits();
    
    /*
    *
    */    
    public long getItemReceives();
    
    /*
    *
    */    
    public long getItemErrors();
    
    /*
    *
    */    
    public int getAccessRights();
    
    /*
    *
    */    
    public void setAccessRights(int accessrigths);    
    
    /*
    *
    */    
    public PlcResponseCode getDataQuality();
    
    /*
    *
    */    
    public void setDataQuality(PlcResponseCode dataquality);

    /*
    *
    */    
    public void setPlcValue(PlcValue plcvalue);
       
    /*
    *
    */    
    public PlcValue getItemPlcValue();
    
    /*
    *
    */    
    public ByteBuf getItemByteBuf();
    
    /*
    *
    */    
    public byte[] getItemBytes();
    
    /*
    *
    */    
    public void addItemClient(PlcItemListener client);
    
    /*
    *
    */    
    public void removeItemClient(PlcItemListener client);    
    
    /*
    *
    */
    public Hashtable<String, Object> getProperties();
    
    /*
    *
    */    
    public Date getLastReadDate();
    
    /*
    *
    */    
    public Date getLastWriteDate();
    
    /*
    *
    */    
    public Date getLastErrorDate();
    
}

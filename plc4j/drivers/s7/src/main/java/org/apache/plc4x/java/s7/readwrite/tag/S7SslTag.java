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
package org.apache.plc4x.java.s7.readwrite.tag;

import java.util.Collections;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcValueType;

public class S7SslTag implements PlcTag {
    
   //SZL_ID=0xYYYY;INDEX=0xZZZZ
    private static final Pattern SSL_ADDRESS_PATTERN =
        Pattern.compile("^SSL_ID=(?<sslId>16#[0-9a-fA-F]{4});INDEX=(?<index>16#[0-9a-fA-F]{4})");  
    
    private static final String SSL_ID = "sslId";    
    private static final String INDEX = "index";    
    
    private final int szlId;
    private final int index;

    public S7SslTag(int sslId, int index) {
        this.szlId = sslId;
        this.index = index;
    }

    @Override
    public String getAddressString() {
        return String.format("SSL_ID=%s;INDEX=16#%d", szlId, index);
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.RAW_BYTE_ARRAY;
    }

    @Override
    public List<ArrayInfo> getArrayInfo() {
        return Collections.emptyList();
    }

    public int getSslId() {
        return szlId;
    }

    public int getIndex() {
        return index;
    }
          
    public static boolean matches(String tagString) {
        return SSL_ADDRESS_PATTERN.matcher(tagString).matches();
    }   
    
    public static S7SslTag of(String tagString) {
        Matcher matcher = SSL_ADDRESS_PATTERN.matcher(tagString);
        if (matcher.matches()){
            String strSxlId = matcher.group(SSL_ID);
            String strIndex = matcher.group(INDEX);
            strSxlId = strSxlId.replaceAll("16#", "");
            strIndex = strIndex.replaceAll("16#", "");
            return new S7SslTag(Integer.parseInt(strSxlId, 16),Integer.parseInt(strIndex, 16));
        }  
        throw new PlcInvalidTagException("Unable to parse address: " + tagString);
    }       

}

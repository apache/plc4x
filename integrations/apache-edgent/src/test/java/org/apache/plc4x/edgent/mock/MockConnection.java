/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.edgent.mock;

import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.CompletableFuture;

import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.connection.AbstractPlcConnection;
import org.apache.plc4x.java.api.connection.PlcReader;
import org.apache.plc4x.java.api.connection.PlcWriter;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.exceptions.PlcIoException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.model.Address;

public class MockConnection extends AbstractPlcConnection implements PlcReader, PlcWriter {

    private final String url;
    private final PlcAuthentication authentication;
    private final Map<Address,Object> dataValueMap = new HashMap<>();
    private long curReadCnt;
    private int readExceptionTriggerCount;
    private String readExceptionMsg;
    private long curWriteCnt;
    private int writeExceptionTriggerCount;
    private String writeExceptionMsg;
  
    public MockConnection(String url) {
      this(url, null);
    }
  
    public MockConnection(String url, PlcAuthentication authentication) {
      this.url = url;
      this.authentication = authentication;
    }
  
    public PlcAuthentication getAuthentication() {
        return authentication;
    }
    
    public String getUrl() {
      return url;
    }

    @Override
    public void connect() throws PlcConnectionException {

    }

    @Override
    public void close() throws Exception {

    }

    @Override
    public Address parseAddress(String addressString) throws PlcException {
      return new MockAddress(addressString);
    }

    @SuppressWarnings("unchecked")
    @Override
    public <T> CompletableFuture<PlcReadResponse<T>> read(PlcReadRequest<T> readRequest) {
      curReadCnt++;
      if (readExceptionTriggerCount > 0 && curReadCnt == readExceptionTriggerCount) {
        curReadCnt = 0;
        CompletableFuture<PlcReadResponse<T>> cf = new CompletableFuture<>();
        cf.completeExceptionally(new PlcIoException(readExceptionMsg));
        return cf;
      }
      PlcReadResponse<T> response = readRequest.createResponse((T)getDataValue(readRequest.getAddress()));
      CompletableFuture<PlcReadResponse<T>> cf = CompletableFuture.completedFuture(response);
      return cf;
    }

    @Override
    public <T> CompletableFuture<PlcWriteResponse<T>> write(PlcWriteRequest<T> writeRequest) {
      curWriteCnt++;
      if (writeExceptionTriggerCount > 0 && curWriteCnt == writeExceptionTriggerCount) {
        curWriteCnt = 0;
        CompletableFuture<PlcWriteResponse<T>> cf = new CompletableFuture<>();
        cf.completeExceptionally(new PlcIoException(writeExceptionMsg));
        return cf;
      }
      setDataValue(writeRequest.getAddress(), writeRequest.getValue());
      CompletableFuture<PlcWriteResponse<T>> cf = CompletableFuture.completedFuture(writeRequest.createResponse());
      return cf;
    }

    public void setDataValue(Address address, Object o) {
      dataValueMap.put(address, o);
    }
    public Object getDataValue(Address address) {
      return dataValueMap.get(address);
    }
    public Map<Address,Object> getAllDataValues() {
      return dataValueMap;
    }
    public void clearAllDataValues() {
      dataValueMap.clear();
    }
    public void setReadException(int readTriggerCount, String msg) {
      readExceptionTriggerCount = readTriggerCount;
      readExceptionMsg = msg;
      curReadCnt = 0;
    }
    public void setWriteException(int writeTriggerCount, String msg) {
      writeExceptionTriggerCount = writeTriggerCount;
      writeExceptionMsg = msg;
      curWriteCnt = 0;
    }
}

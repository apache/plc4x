/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.examples.integration.iotdb;

import java.util.Collections;
import org.apache.iotdb.rpc.IoTDBConnectionException;
import org.apache.iotdb.rpc.StatementExecutionException;
import org.apache.iotdb.rpc.TSStatusCode;
import org.apache.iotdb.session.pool.SessionPool;
import org.apache.iotdb.tsfile.file.metadata.enums.CompressionType;
import org.apache.iotdb.tsfile.file.metadata.enums.TSDataType;
import org.apache.iotdb.tsfile.file.metadata.enums.TSEncoding;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class IoTDBWriterWithSession implements  IIoTDBWriter {
    private static Logger LOGGER = LoggerFactory.getLogger(IoTDBWriterWithSession.class);

    //sessionPool is more user friendly than session
    SessionPool sessionPool;

    public IoTDBWriterWithSession(String iotdbIpPort, String user, String password) {
        sessionPool = new SessionPool(iotdbIpPort.split(":")[0], Integer.valueOf(iotdbIpPort.split(":")[1]), user, password, 1);
    }

    @Override
    public void initStorageGroup(String storageGroup) {
        try {
            sessionPool.setStorageGroup(storageGroup);
        } catch (IoTDBConnectionException e) {
            LOGGER.error("Error", e);
        } catch (StatementExecutionException e) {
            if (e.getStatusCode() != TSStatusCode.PATH_ALREADY_EXIST_ERROR.getStatusCode()) { // 300 means the storage group exist already.
                LOGGER.error("Error", e);
            }
        }
    }

    @Override
    public void writeData(String deviceId, String field, long timestamp, Integer value) {
        try {
            //when using the default configuration of IoTDB, then int value will be considered as float.
            //so we use float here.
            //change the data type by modify `integer_string_infer_type` parameter.

            //if you create time-series manually, REMEMBER TO MODIFY THE TSDATATYPE.FLOAT AS WHAT YOU REALLY NEED.
            sessionPool.insertRecord(deviceId, timestamp, Collections.singletonList(field), Collections.singletonList(
                TSDataType.FLOAT), Collections.singletonList(value + 0.0f));
        } catch (IoTDBConnectionException | StatementExecutionException e) {
            LOGGER.error("Error storing connection.", e);
        }
    }

    @Override
    public void close() {
        sessionPool.close();
    }

    @Override
    public void createTimeseries(String timeseries, String dataType) {
        try {
            sessionPool.createTimeseries(timeseries, TSDataType.valueOf(dataType), TSEncoding.RLE, CompressionType.SNAPPY);
        } catch (IoTDBConnectionException e) {
            LOGGER.error("Error", e);
        } catch (StatementExecutionException e) {
            if (e.getStatusCode() != TSStatusCode.PATH_ALREADY_EXIST_ERROR.getStatusCode()) { // 300 means the time series exist already.
                LOGGER.error("Error", e);
            }
        }
    }

}

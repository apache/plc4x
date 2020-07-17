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

import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.SQLException;
import java.sql.Statement;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class IoTDBWriterWithJDBC implements IIoTDBWriter{
    private static Logger LOGGER = LoggerFactory.getLogger(IoTDBWriterWithJDBC.class);

    //IoTDB JDBC connection
    Connection connection;

    //IoTDB JDBC Statement
    Statement statement;

    public IoTDBWriterWithJDBC(String ipPort, String username, String password)
        throws ClassNotFoundException, SQLException {
        // Get IoTDB connection
        Class.forName("org.apache.iotdb.jdbc.IoTDBDriver");
        connection = DriverManager.getConnection("jdbc:iotdb://" + ipPort+"/",
            username, password);
        statement = connection.createStatement();
    }

    @Override
    public void initStorageGroup(String storageGroup) {
        try {
            statement.execute("SET STORAGE GROUP TO root." + storageGroup);
        } catch (SQLException e) {
            //from v0.9.0, you can use the error code to check whether the sg exists.
            LOGGER.error(e.getMessage());
        }
    }

    @Override
    public void writeData(String deviceId, String field, long timestamp, Integer value) {
        //please modify this method if you want to write multiple fields once.
        try {
            String sql = String.format("insert into %s (timestamp, %s) values (%d, %s)",
                deviceId, field, timestamp, value + "");
            statement.execute(sql);
        } catch (SQLException e) {
            LOGGER.error("Error storing data.", e);
        }
    }

    @Override
    public void close() {
        try {
            connection.close();
        } catch (SQLException e) {
            LOGGER.error("Error closing connection.", e);
        }
    }

    @Override
    public void createTimeseries(String timeseries, String dataType) {
        try {
            statement.execute("create timeseries " + timeseries + " with dataType=" + dataType +",encoding=RLE");
        } catch (SQLException e) {
            //from v0.9.0, you can use the error code to check whether the sg exists.
            LOGGER.error(e.getMessage());
        }
    }

}

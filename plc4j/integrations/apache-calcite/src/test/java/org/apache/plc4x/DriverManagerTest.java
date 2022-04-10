/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x;

import org.apache.calcite.jdbc.CalciteConnection;
import org.apache.calcite.jdbc.Driver;
import org.apache.plc4x.java.scraper.config.ScraperConfiguration;
import org.apache.plc4x.java.scraper.config.ScraperConfigurationClassicImpl;
import org.apache.plc4x.java.scraper.exception.ScraperException;
import org.assertj.core.api.WithAssertions;
import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.sql.Connection;
import java.sql.ResultSet;
import java.sql.ResultSetMetaData;
import java.sql.SQLException;
import java.util.Properties;

public class DriverManagerTest implements WithAssertions {

    @Test
    void query() throws SQLException, IOException, ScraperException {
        Driver driver = new Driver();
        Connection connection = driver.connect("jdbc:calcite:asdf;lex=MYSQL_ANSI", new Properties());

        CalciteConnection calciteConnection = connection.unwrap(CalciteConnection.class);
        calciteConnection.getRootSchema().add("plc4x", new Plc4xSchema(ScraperConfiguration.fromFile("src/test/resources/example.yml", ScraperConfigurationClassicImpl.class), 100));

        ResultSet rs = connection.prepareStatement("SELECT * FROM \"plc4x\".\"job1\"").executeQuery();
        validateResult(rs);

        connection.close();
    }

    @Test
    void query2() throws IOException, SQLException {
        Driver driver = new Driver();
        Connection connection = driver.connect("jdbc:calcite:model=src/test/resources/model.json", new Properties());

        ResultSet rs = connection.prepareStatement("SELECT * FROM \"PLC4X-TABLES\".\"job1\"").executeQuery();

        validateResult(rs);

        connection.close();
    }

    private void validateResult(ResultSet rs) throws SQLException {
        // Assert columns
        ResultSetMetaData metadata = rs.getMetaData();
        assertThat(metadata.getColumnCount()).isEqualTo(4);
        // Column names
        assertThat(metadata.getColumnName(1)).isEqualTo("timestamp");
        assertThat(metadata.getColumnName(2)).isEqualTo("source");
        assertThat(metadata.getColumnName(3)).isEqualTo("test");
        assertThat(metadata.getColumnName(4)).isEqualTo("test2");
        // Column types
        assertThat(metadata.getColumnTypeName(1)).isEqualTo("TIMESTAMP");
        assertThat(metadata.getColumnTypeName(2)).isEqualTo("VARCHAR");
        assertThat(metadata.getColumnTypeName(3)).isEqualTo("INTEGER");
        assertThat(metadata.getColumnTypeName(4)).isEqualTo("VARCHAR");

        int rowCount = 0;
        while (rs.next()) {
            rowCount++;
        }

        assertThat(rowCount).isEqualTo(100);
    }

}

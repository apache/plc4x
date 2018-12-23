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
package org.apache.plc4x;

import org.apache.calcite.jdbc.CalciteConnection;
import org.apache.calcite.jdbc.Driver;
import org.apache.plc4x.java.scraper.config.ScraperConfiguration;
import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.sql.Connection;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.util.Properties;

public class DriverManagerTest {

    @Test
    void instanciateJdbcConnection() throws SQLException, IOException {
        Driver driver = new Driver();
        Connection connection = driver.connect("jdbc:calcite://asdf;config=abc;lex=MYSQL_ANSI", new Properties());

        CalciteConnection calciteConnection = connection.unwrap(CalciteConnection.class);
        calciteConnection.getRootSchema().add("plc4x", new Plc4xSchema(ScraperConfiguration.fromFile("src/test/resources/example.yml"), 10));

        // ResultSet rs = connection.prepareStatement("SELECT STREAM \"test\", \"test\" * 2, \"test2\" FROM \"plc4x\".\"job1\"").executeQuery();
        ResultSet rs = connection.prepareStatement("SELECT STREAM * FROM \"plc4x\".\"job1\" WHERE source = 'test'").executeQuery();

        // Print the header
        int count = rs.getMetaData().getColumnCount();
        for (int i = 1; i <= count; i++) {
            System.out.print(rs.getMetaData().getColumnLabel(i) + "(" + rs.getMetaData().getColumnTypeName(i) + ")" + "\t");
        }
        System.out.println("");

        while (rs.next()) {
            for (int i = 1; i <= count; i++) {
                System.out.print(rs.getString(i) + "\t");
            }
            System.out.println("");
        }

        connection.close();
    }

    @Test
    void instantiateDirect() throws IOException, SQLException {
        Driver driver = new Driver();
        Connection connection = driver.connect("jdbc:calcite:model=src/test/resources/model.yml;lex=MYSQL_ANSI", new Properties());

        // ResultSet rs = connection.prepareStatement("SELECT STREAM \"test\", \"test\" * 2, \"test2\" FROM \"plc4x\".\"job1\"").executeQuery();
        ResultSet rs = connection.prepareStatement("SELECT * FROM \"plc4x-tables\".\"job1\"").executeQuery();

        // Print the header
        int count = rs.getMetaData().getColumnCount();
        for (int i = 1; i <= count; i++) {
            System.out.print(rs.getMetaData().getColumnLabel(i) + "(" + rs.getMetaData().getColumnTypeName(i) + ")" + "\t");
        }
        System.out.println("");

        int row = 1;

        while (rs.next()) {
            System.out.print(row++ + "\t");
            for (int i = 1; i <= count; i++) {
                System.out.print(rs.getString(i) + "\t");
            }
            System.out.println("");
        }

        connection.close();
    }

}

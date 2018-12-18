package org.apache.plc4x;

import org.apache.calcite.jdbc.CalciteConnection;
import org.apache.calcite.jdbc.Driver;
import org.apache.calcite.schema.Schema;
import org.apache.plc4x.java.scraper.config.ScraperConfiguration;
import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.sql.Connection;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.util.Collections;
import java.util.Properties;

public class DriverManagerTest {

    @Test
    void instanciateJdbcConnection() throws SQLException, IOException {
        Driver driver = new Driver();
        Connection connection = driver.connect("jdbc:calcite://asdf;config=abc", new Properties());

        CalciteConnection calciteConnection = connection.unwrap(CalciteConnection.class);
        calciteConnection.getRootSchema().add("plc4x", new Plc4xSchema(ScraperConfiguration.fromFile("src/test/resources/example.yml")));

        ResultSet rs = connection.prepareStatement("SELECT STREAM * FROM \"plc4x\".\"job1\"").executeQuery();

        while (rs.next()) {
            System.out.print("Spalte 1: " + rs.getString(1) + "\t");
            System.out.println("Spalte 2: " + rs.getString(2));
        }

        connection.close();
    }
}

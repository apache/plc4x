package org.apache.plc4x.java.api.authentication;

import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class PlcUsernamePasswordAuthenticationTest {

    @Test
    @Tag("fast")
    void authenication() {
        PlcUsernamePasswordAuthentication authenication = new PlcUsernamePasswordAuthentication("user", "password");

        assertTrue(authenication.getUsername().equals("user"), "Unexpected user name");
        assertTrue(authenication.getPassword().equals("password"), "Unexpected password");
    }

}
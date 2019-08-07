package org.apache.plc4x.java.df1;

import org.junit.Test;

import java.util.regex.Matcher;

import static org.apache.plc4x.java.df1.DF1PlcDriver.DF1_URI_PATTERN;
import static org.apache.plc4x.java.df1.DF1PlcDriver.SERIAL_PATTERN;
import static org.junit.Assert.*;

public class DF1PlcDriverTest {

    @Test
    public void matchExpression() {
        Matcher matcher = SERIAL_PATTERN.matcher("serial:///COM4");

        assertTrue(matcher.matches());
    }

    @Test
    public void matchExpression2() {
        Matcher matcher = DF1_URI_PATTERN.matcher("df1:serial:///COM4");

        assertTrue(matcher.matches());
    }
}
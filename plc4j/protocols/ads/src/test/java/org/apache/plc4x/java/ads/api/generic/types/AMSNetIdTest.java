package org.apache.plc4x.java.ads.api.generic.types;

import org.junit.Test;

import static org.hamcrest.Matchers.is;
import static org.junit.Assert.*;

public class AMSNetIdTest {

    @Test
    public void netIdBytes() {
        // note bytes in reverse order
        AMSNetId netid = AMSNetId.of((byte)0x01, (byte)0x02, (byte)0x03, (byte)0x04, (byte)0x05, (byte)0x06);
        assertThat(netid.toString(), is("1.2.3.4.5.6"));
    }

    @Test
    public void netIdString() {
        // note bytes in reverse order
        AMSNetId netid = AMSNetId.of("1.2.3.4.5.6");
        assertThat(netid.toString(), is("1.2.3.4.5.6"));
    }
    
    @Test(expected = IllegalArgumentException.class)
    public void netIdTooShort() {
        // note bytes in reverse order
        AMSNetId netid = AMSNetId.of("1.2.3.4");
    }

    @Test(expected = IllegalArgumentException.class)
    public void netIdStringTooLong() {
        // note bytes in reverse order
        AMSNetId netid = AMSNetId.of("1.2.3.4.5.6.7.8");
    }

    @Test(expected = IllegalArgumentException.class)
    public void netIdStringWrongSeperator() {
        // note bytes in reverse order
        AMSNetId netid = AMSNetId.of("1:2:3:4:5:6");
    }

}
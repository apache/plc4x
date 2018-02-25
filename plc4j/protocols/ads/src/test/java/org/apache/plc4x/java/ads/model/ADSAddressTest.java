package org.apache.plc4x.java.ads.model;

import org.junit.After;
import org.junit.Before;
import org.junit.Test;

import static org.hamcrest.CoreMatchers.is;
import static org.junit.Assert.*;

public class ADSAddressTest {

    @Test
    public void of() {
        ADSAddress address = ADSAddress.of("1/10");
        assertThat(address.getIndexGroup(), is(1L));
        assertThat(address.getIndexOffset(), is(10L));
    }

    @Test(expected  = IllegalArgumentException.class)
    public void stringInAddress() {
        ADSAddress address = ADSAddress.of("group/offset");
    }

    @Test(expected  = IllegalArgumentException.class)
    public void singleNumberAddress() {
        ADSAddress address = ADSAddress.of("10");
    }

    @Test(expected  = IllegalArgumentException.class)
    public void wrongSeperator() {
        ADSAddress address = ADSAddress.of("1:10");
    }

    @Test
    public void getGroupAndOffset() {
        ADSAddress address = ADSAddress.of(2L, 20L);
        assertThat(address.getIndexGroup(), is(2L));
        assertThat(address.getIndexOffset(), is(20L));
    }
}
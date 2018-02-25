package org.apache.plc4x.java.ads.api.generic.types;

import org.junit.Test;

import java.util.EnumSet;

import static org.hamcrest.Matchers.is;
import static org.hamcrest.Matchers.isEmptyString;
import static org.hamcrest.Matchers.not;
import static org.junit.Assert.*;

public class StateTest {

    // Not the best unit tests but here for coverage

    @Test
    public void stateBitFields() {
        int bitMask = State.StateMask.RESPONSE.getMask() | State.StateMask.NO_RETURN.getMask() | State.StateMask.ADS_COMMAND.getMask()
                        | State.StateMask.SYSTEM_COMMAND.getMask() | State.StateMask.HIGH_PRIORITY_COMMAND.getMask() | State.StateMask.TIMESTAMP_ADDED.getMask()
                        | State.StateMask.UDP_COMMAND.getMask() | State.StateMask.INIT_COMMAND.getMask() | State.StateMask.BROADCAST.getMask();
        State state = State.of(bitMask);

        assertThat(state.toString(), not((isEmptyString())));
    }

    @Test
    public void equals() {
        State a = State.of((byte) 0x1, (byte) 0x2);
        State b = State.of((byte) 0x1, (byte) 0x2);
        State c = State.of((byte) 0x1, (byte) 0x4);
        byte array[] = {(byte) 0x1, (byte) 0x2};

        assertThat(a.equals(a), is(true));
        assertThat(a.equals(b), is(true));
        assertThat(a.equals(c), is(false));
        assertThat(a.equals(1), is(false));
        assertThat(a.equals((byte) 1), is(false));
        assertThat(a.equals(array), is(false));
    }

}
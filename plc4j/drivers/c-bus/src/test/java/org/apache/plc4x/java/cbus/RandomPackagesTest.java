package org.apache.plc4x.java.cbus;

import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.cbus.readwrite.CBusCommand;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;

@Disabled("non of those work yet")
public class RandomPackagesTest {

    static final String BACKSLASH = "5C";

    // 4.2.9.1
    @Test
    void pointToPointCommand1() throws Exception {
        byte[] bytes = Hex.decodeHex(BACKSLASH + "0603002102D4");
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        CBusCommand cBusCommand = CBusCommand.staticParse(readBufferByteBased, false);
        assertThat(cBusCommand)
            .isNotNull();
    }

    // 4.2.9.1
    @Test
    void pointToPointCommand2() throws Exception {
        byte[] bytes = Hex.decodeHex(BACKSLASH + "06420903210289");
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        CBusCommand cBusCommand = CBusCommand.staticParse(readBufferByteBased, false);
        assertThat(cBusCommand)
            .isNotNull();
    }

    // 4.2.9.2
    @Test
    void pointToMultiPointCommand1() throws Exception {
        byte[] bytes = Hex.decodeHex(BACKSLASH + "0538000108BA");
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        CBusCommand cBusCommand = CBusCommand.staticParse(readBufferByteBased, false);
        assertThat(cBusCommand)
            .isNotNull();
    }

    // 4.2.9.2
    @Test
    void pointToMultiPointCommand2() throws Exception {
        byte[] bytes = Hex.decodeHex(BACKSLASH + "05FF007A38004A");
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        CBusCommand cBusCommand = CBusCommand.staticParse(readBufferByteBased, false);
        assertThat(cBusCommand)
            .isNotNull();
    }

    // 4.2.9.3
    @Test
    void pointToPointToMultiPointCommand2() throws Exception {
        byte[] bytes = Hex.decodeHex(BACKSLASH + "03420938010871");
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        CBusCommand cBusCommand = CBusCommand.staticParse(readBufferByteBased, false);
        assertThat(cBusCommand)
            .isNotNull();
    }

    // 4.3.3.1
    @Test
    void calReply1() throws Exception {
        byte[] bytes = Hex.decodeHex(BACKSLASH + "0605002102");
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        CBusCommand cBusCommand = CBusCommand.staticParse(readBufferByteBased, false);
        assertThat(cBusCommand)
            .isNotNull();
    }

    // 4.3.3.1
    @Test
    void calReply2() throws Exception {
        byte[] bytes = Hex.decodeHex(BACKSLASH + "86059300890231E22E363620207F");
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        CBusCommand cBusCommand = CBusCommand.staticParse(readBufferByteBased, false);
        assertThat(cBusCommand)
            .isNotNull();
    }

    // 4.3.3.2
    @Test
    void monitoredSal() throws Exception {
        byte[] bytes = Hex.decodeHex(BACKSLASH + "0503380079083F");
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        CBusCommand cBusCommand = CBusCommand.staticParse(readBufferByteBased, false);
        assertThat(cBusCommand)
            .isNotNull();
    }


}

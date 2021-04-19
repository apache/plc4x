package org.apache.plc4x.java.canopen.transport;

import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.generation.MessageIO;
import org.apache.plc4x.java.spi.transport.Transport;

public interface CANTransport extends Transport {

    MessageIO<CANOpenFrame, CANOpenFrame> getMessageIO(Configuration cfg);

    Class<CANOpenFrame> getMessageType();
}

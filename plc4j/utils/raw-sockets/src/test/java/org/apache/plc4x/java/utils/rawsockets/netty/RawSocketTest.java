package org.apache.plc4x.java.utils.rawsockets.netty;

import java.util.Arrays;

public class RawSocketTest {

    public static void main(String[] args) throws Exception {
        // Protocol number 1 = ICMP (Ping)
        RawSocket rawSocket = new RawSocket(1);

        // Simply print the result to the console
        rawSocket.addListener(new RawSocketListener() {
            @Override
            public void packetReceived(byte[] rawData) {
                System.out.println("Got response:");
                System.out.println(Arrays.toString(rawData));
            }
        });

        // Connect to the remote address
        // This doesn't really connect in the usual sense, it
        // does the ARP MAC address lookup and sets up the listener
        // to accept packets sent from that mac address to the
        // current machines with the given IP protocol id.
        rawSocket.connect("10.10.56.1");

        // Simple ICMP (Ping packet)
        byte[] rawData = new byte[] {
            // Type (ICMP Ping Request) & Code (just 0)
            (byte) 0x08, (byte) 0x00,
            // Checksum
            (byte) 0xe3, (byte) 0xe5,
            // Identifier
            (byte) 0x00, (byte) 0x01,
            // Sequence Number
            (byte) 0x00, (byte) 0x00,
            // Payload (Just random data that was used to fit to the checksum)
            (byte) 0x00, (byte) 0x01, (byte) 0x02, (byte) 0x03, (byte) 0x04, (byte) 0x05, (byte) 0x06, (byte) 0x07, (byte) 0x08, (byte) 0x09};

        // Write the raw packet to the remote host.
        rawSocket.write(rawData);
    }

}

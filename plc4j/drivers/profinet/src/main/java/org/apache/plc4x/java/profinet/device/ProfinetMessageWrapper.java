package org.apache.plc4x.java.profinet.device;

import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.profinet.readwrite.*;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBufferByteBased;
import org.apache.plc4x.java.utils.rawsockets.netty.RawSocketChannel;

import java.io.IOException;
import java.net.DatagramPacket;

public class ProfinetMessageWrapper {

    public static void sendMessage(ProfinetCallable callable, ProfinetDevice context) throws RuntimeException {
        try {
            DceRpc_Packet packet = callable.create();
            // Serialize it to a byte-payload
            WriteBufferByteBased writeBuffer = new WriteBufferByteBased(packet.getLengthInBytes());
            packet.serialize(writeBuffer);
            // Create a udp packet.
            DatagramPacket connectRequestPacket = new DatagramPacket(writeBuffer.getData(), writeBuffer.getData().length);
            connectRequestPacket.setAddress(context.getIpAddress());
            connectRequestPacket.setPort(context.getPort());

            // Send it.
            context.getUdpSocket().send(connectRequestPacket);

            // Receive the response.
            byte[] resultBuffer = new byte[packet.getLengthInBytes()];
            DatagramPacket connectResponsePacket = new DatagramPacket(resultBuffer, resultBuffer.length);
            context.getUdpSocket().receive(connectResponsePacket);
            ReadBufferByteBased readBuffer = new ReadBufferByteBased(resultBuffer);
            final DceRpc_Packet dceRpc_packet = DceRpc_Packet.staticParse(readBuffer);
            callable.handle(dceRpc_packet);
        } catch (SerializationException e) {
            throw new RuntimeException(e);
        } catch (IOException e) {
            throw new RuntimeException(e);
        } catch (ParseException e) {
            throw new RuntimeException(e);
        } catch (PlcException e) {
            throw new RuntimeException(e);
        }

    }
}

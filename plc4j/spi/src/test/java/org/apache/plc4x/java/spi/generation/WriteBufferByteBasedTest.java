package org.apache.plc4x.java.spi.generation;

import org.apache.commons.io.HexDump;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;

import java.io.ByteArrayOutputStream;
import java.math.BigInteger;

import static org.junit.jupiter.api.Assertions.*;

class WriteBufferByteBasedTest {

    @Nested
    class WriteBigInteger {
        @Nested
        class BigEndian {

            @Test
            void zero() throws Exception {
                WriteBufferByteBased SUT = new WriteBufferByteBased(1, ByteOrder.LITTLE_ENDIAN);
                SUT.writeBigInteger(8, BigInteger.ZERO);
                byte[] data = SUT.getData();
                System.out.println(toHex(data));
                // TODO: check right representation
                assertArrayEquals(new byte[]{0b0000_0000}, data);
                assertEquals(BigInteger.ZERO, new BigInteger(data));
            }

            @Test
            void one() throws Exception {
                WriteBufferByteBased SUT = new WriteBufferByteBased(1, ByteOrder.LITTLE_ENDIAN);
                SUT.writeBigInteger(8, BigInteger.ONE);
                byte[] data = SUT.getData();
                System.out.println(toHex(data));
                // TODO: check right representation
                assertArrayEquals(new byte[]{0b0000_0001}, data);
                assertEquals(BigInteger.ZERO, new BigInteger(data));
            }

            @Test
            void minusOne() throws Exception {
                WriteBufferByteBased SUT = new WriteBufferByteBased(8, ByteOrder.LITTLE_ENDIAN);
                SUT.writeBigInteger(8, BigInteger.ZERO.subtract(BigInteger.ONE));
                byte[] data = SUT.getData();
                System.out.println(toHex(data));
                // TODO: check right representation
                assertArrayEquals(new byte[]{0b0000_0001}, data);
                assertEquals(BigInteger.ZERO, new BigInteger(data));
            }

            @Test
            void minus255() throws Exception {
                WriteBufferByteBased SUT = new WriteBufferByteBased(8, ByteOrder.LITTLE_ENDIAN);
                SUT.writeBigInteger(8, BigInteger.valueOf(-255L));
                byte[] data = SUT.getData();
                System.out.println(toHex(data));
                // TODO: check right representation
                assertArrayEquals(new byte[]{(byte) 0b1000_0000, 0b0000_0001}, data);
                assertEquals(BigInteger.valueOf(-255L), new BigInteger(data));
            }

        }

        @Nested
        class LittleEndian {

            @Test
            void writeBigInteger_LE() throws Exception {
                WriteBufferByteBased SUT_LE = new WriteBufferByteBased(8012, ByteOrder.LITTLE_ENDIAN);
                SUT_LE.writeBigInteger(1, BigInteger.ZERO);
            }
        }
    }

    @Nested
    class WriteUnsignedBigInteger {
        @Nested
        class BigEndian {

            @Test
            void zero() throws Exception {
                WriteBufferByteBased SUT = new WriteBufferByteBased(1, ByteOrder.LITTLE_ENDIAN);
                SUT.writeUnsignedBigInteger(8, BigInteger.ZERO);
                byte[] data = SUT.getData();
                System.out.println(toHex(data));
                // TODO: check right representation
                assertArrayEquals(new byte[]{0b0000_0000}, data);
                assertEquals(BigInteger.ZERO, new BigInteger(data));
            }

            @Test
            void one() throws Exception {
                WriteBufferByteBased SUT = new WriteBufferByteBased(1, ByteOrder.LITTLE_ENDIAN);
                SUT.writeUnsignedBigInteger(8, BigInteger.ONE);
                byte[] data = SUT.getData();
                System.out.println(toHex(data));
                // TODO: check right representation
                assertArrayEquals(new byte[]{0b0000_0001}, data);
                assertEquals(BigInteger.ZERO, new BigInteger(data));
            }

            @Test
            void minusOne() throws Exception {
                WriteBufferByteBased SUT = new WriteBufferByteBased(8, ByteOrder.LITTLE_ENDIAN);
                SUT.writeUnsignedBigInteger(8, BigInteger.ZERO.subtract(BigInteger.ONE));
                byte[] data = SUT.getData();
                System.out.println(toHex(data));
                // TODO: check right representation
                assertArrayEquals(new byte[]{0b0000_0001}, data);
                assertEquals(BigInteger.ZERO, new BigInteger(data));
            }

            @Test
            void minus255() throws Exception {
                WriteBufferByteBased SUT = new WriteBufferByteBased(8, ByteOrder.LITTLE_ENDIAN);
                SUT.writeUnsignedBigInteger(8, BigInteger.valueOf(-255L));
                byte[] data = SUT.getData();
                System.out.println(toHex(data));
                // TODO: check right representation
                assertArrayEquals(new byte[]{(byte) 0b1000_0000, 0b0000_0001}, data);
                assertEquals(BigInteger.valueOf(-255L), new BigInteger(data));
            }
        }

        @Nested
        class LittleEndian {

            @Test
            void writeBigInteger_LE() throws Exception {
                WriteBufferByteBased SUT_LE = new WriteBufferByteBased(8012, ByteOrder.LITTLE_ENDIAN);
                SUT_LE.writeBigInteger(1, BigInteger.ZERO);
            }
        }
    }

    public static String toHex(byte[] bytes) throws Exception {
        try (ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream()) {
            HexDump.dump(bytes, 0, byteArrayOutputStream, 0);
            return byteArrayOutputStream.toString();
        }
    }
}
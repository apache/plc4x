package org.apache.plc4x.java.spi.utils;

public class ByteOrderUtils {
    /**
     * Byte order is 'A B C D'
     * @param num int number
     * @return int number in big endian format
     */
    public static int BigEndian(int num){
        return num;
    }
    /**
     * Byte order is 'D C B A'
     * @param num int number
     * @return int number in little endian format
     */
    public static int LittleEndian(int num){
        return Integer.reverseBytes(num);
    }
    /**
     * Byte order is 'B A D C'
     * @param num int number
     * @return int number in big endian swap format
     */
    public static int BigEndianWordSwap(int num){
        return (num << 16)|(num >>> 16);
    }
    /**
     * Byte order is 'C D A B'
     * @param num int number
     * @return int number in little endian swap format
     */
    public static int LittleEndianWordSwap(int num){
        return ((num&0xff00)>>>8)|
                ((num<<8)&0xff00)|
                ((num<<8)&0xff000000)|
                ((num &0xff000000)>>>8);
    }



    /**
     * Byte order is 'A B C D E F G H'
     * @param num long number
     * @return long number in big endian format
     */
    public static long BigEndian(long num){
        return num;
    }

    /**
     * Byte order is 'H G F E D C B A'
     * @param num long number
     * @return long number in little endian format
     */
    public static long LittleEndian(long num){
        return Long.reverseBytes(num);
    }

    /**
     * Byte order is 'B A D C F E H G'
     * @param num long number
     * @return long number in big endian format
     */
    public static long BigEndianWordSwap(long num){

        return (num & 0x00ff00ff00ff00ffL) << 8 | (num >>> 8) & 0x00ff00ff00ff00ffL;
    }

    /**
     * Byte order is 'G H E F C D A B'
     * @param num long number
     * @return long number in little endian format
     */
    public static long LittleEndianWordSwap(long num){

        return (num & 0xffff000000000000L) >>> 48 |
                (num & 0x0000ffff00000000L) >>> 16 |
                (num & 0x00000000ffff0000L) << 16 |
                (num & 0x000000000000ffffL) << 48;
    }


}

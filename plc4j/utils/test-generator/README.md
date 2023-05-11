# Test-Generator

It is the purpose of the test-generator, to take a given pcap capture file and to produce a ParserSerializer testsuite document from it.

As the tool can only generate XML for things implemented in PLC4X this might sound counter-intuitive. 
This however is only a first step: By converting, we now have each packet in the PCAP recording in a form, that we can use to build a DriverTestsuite type of test. 

## Usage

In order to generate the XML:

    ParserSerializerTestsuiteGenerator -o "read-write" -p "s7" -t "CyclicExchange" org.apache.plc4x.java.s7.readwrite.TPKTPacket /Users/christoferdutz/Downloads/cyc-cotp.pcapng CyclicExchange.xml

The `-o` parameter refers to the `outputFlavor` element in the output document.

The `-p` parameter refers to the `protocolName` element in the output document

The `-t` parameter refers to the root element in the output of each packet

The class name following refers to the class name of the PLC4J type used to parse this message type.

After that comes the path to the `pcap` file

Last comes the output filename of the generated document.

NOTES:

- Be sure that the driver module for the driver you want to generate a testsuite for is available in the classpath, or the generator will not be able to load the base packet type.
- On Mac I needed to define the `jna.library.path` to where the `libpcap` is found, or it will use the Mac default, which is not compatible: `-Djna.library.path=/usr/local/Cellar/libpcap/1.10.4/lib`
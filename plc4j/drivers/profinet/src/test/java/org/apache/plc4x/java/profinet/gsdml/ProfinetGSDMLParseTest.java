package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.databind.jsontype.BasicPolymorphicTypeValidator;
import com.fasterxml.jackson.databind.jsontype.PolymorphicTypeValidator;
import com.fasterxml.jackson.dataformat.xml.XmlMapper;
import org.junit.jupiter.api.*;

import java.io.File;
import java.io.IOException;

public class ProfinetGSDMLParseTest {

    @Test
    public void readGsdmlFile() throws IOException {
        XmlMapper xmlMapper = new XmlMapper();

        ProfinetISO15745Profile value = xmlMapper.readValue(new File("src/test/resources/gsdml.xml"), ProfinetISO15745Profile.class);
        System.out.println(8);
    }

}

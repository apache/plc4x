package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.databind.jsontype.BasicPolymorphicTypeValidator;
import com.fasterxml.jackson.databind.jsontype.PolymorphicTypeValidator;
import com.fasterxml.jackson.dataformat.xml.XmlMapper;
import io.vavr.control.Option;
import org.apache.commons.io.filefilter.FalseFileFilter;
import org.junit.jupiter.api.*;

import java.io.File;
import java.io.IOException;

import static org.junit.jupiter.api.Assertions.assertEquals;

@TestInstance(TestInstance.Lifecycle.PER_CLASS)
public class ProfinetGSDMLParseTest {

    private ProfinetISO15745Profile gsdml = null;

    @BeforeAll
    public void setUp() {
        try {
            XmlMapper xmlMapper = new XmlMapper();
            this.gsdml = xmlMapper.readValue(new File("src/test/resources/gsdml.xml"), ProfinetISO15745Profile.class);
        } catch(IOException e) {
            assert false;
        }
    }

    @Test
    public void readGsdmlFile()  {
        assertEquals(this.gsdml.getProfileBody().getDeviceIdentity().getVendorName().getValue(), "Apache PLC4X");
    }

    @Test
    public void readGsdmlFileStartupMode()  {
        ProfinetInterfaceSubmoduleItem interfaceModule = (ProfinetInterfaceSubmoduleItem) this.gsdml.getProfileBody().getApplicationProcess().getDeviceAccessPointList().get(0).getSystemDefinedSubmoduleList().get(0);
        assertEquals(interfaceModule.getApplicationRelations().getStartupMode(), "Advanced");
    }


}

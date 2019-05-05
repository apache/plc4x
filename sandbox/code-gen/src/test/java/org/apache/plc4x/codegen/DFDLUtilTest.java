package org.apache.plc4x.codegen;

import org.junit.Test;

import java.io.File;

public class DFDLUtilTest {

    @Test
    public void loadDFDL() {
        final File schema = new File("src/test/resources/protocol.dfdl.xsd");
        final File outDir = new File("/tmp/");

        System.out.println(schema.getAbsolutePath().toString());

        final DFDLUtil util = new DFDLUtil();

        util.transform(schema, outDir);
    }
}
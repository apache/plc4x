/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */

package org.apache.plc4x.protocols.s7;

import org.apache.commons.codec.binary.Hex;
import org.apache.daffodil.japi.Compiler;
import org.apache.daffodil.japi.*;
import org.apache.daffodil.japi.infoset.JDOMInfosetOutputter;
import org.apache.daffodil.japi.io.InputSourceDataInputStream;
import org.jdom2.Document;
import org.jdom2.output.Format;
import org.jdom2.output.XMLOutputter;

import java.io.ByteArrayInputStream;
import java.net.URI;
import java.net.URL;
import java.util.List;

public class ManualTest {

    public static void main(String[] args) throws Exception {
        Compiler c = Daffodil.compiler();
        c.setValidateDFDLSchemas(true);
        URL schemaUrl = ManualTest.class.getClassLoader().getResource("org/apache/plc4x/protocols/s7/protocol.dfdl.xsd");
        if (schemaUrl != null) {
            URI schemaUri = schemaUrl.toURI();
            ProcessorFactory pf = c.compileSource(schemaUri);
            DataProcessor dp = pf.onPath("/");
            logDiagnosticInformation(dp);

            /*byte[] packet = Hex.decodeHex(
                    "0300004002f080320300" +
                    "0000010002002b000004" +
                    "07ff0300010000ff04002000000000ff0400080000ff0500100000ff0400080000ff0400100000ff04000800");*/
            byte[] packet = Hex.decodeHex(
                "0300004002f080320300" +
                "0000010002002b000004" +
                "07ff0300010000ff04002000000000ff0400080000ff0500100000ff0400080000ff0400100000ff04000800");


            // After having enough bytes available, process the current package.
            JDOMInfosetOutputter outputter = new JDOMInfosetOutputter();
            ParseResult byteMessage = dp.parse(
                new InputSourceDataInputStream(new ByteArrayInputStream(packet)), outputter);
            if (byteMessage.isError()) {
                logDiagnosticInformation(byteMessage);
                return;
            }

            // Get the resulting XML document from the parser.
            Document message = outputter.getResult();

            System.out.println(new XMLOutputter(Format.getPrettyFormat()).outputString(message));
        }
    }

    protected static void logDiagnosticInformation(WithDiagnostics withDiagnostics) {
        List<Diagnostic> diags = withDiagnostics.getDiagnostics();
        for (Diagnostic d : diags) {
            System.err.println(d.getSomeMessage());
        }
    }


}

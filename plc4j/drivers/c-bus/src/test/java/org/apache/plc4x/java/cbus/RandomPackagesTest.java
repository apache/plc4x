/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.cbus;

import org.apache.plc4x.java.cbus.readwrite.*;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.apache.plc4x.java.spi.generation.WriteBufferByteBased;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;

import java.nio.charset.StandardCharsets;

import static org.apache.plc4x.java.cbus.Util.assertMessageMatches;
import static org.assertj.core.api.Assertions.assertThat;

public class RandomPackagesTest {

    public static final CBusOptions C_BUS_OPTIONS_WITH_SRCHK = new CBusOptions(false, false, false, false, false, false, false, false, true);
    RequestContext requestContext;
    CBusOptions cBusOptions;

    @BeforeEach
    void setUp() {
        requestContext = new RequestContext(false, false, false);
        cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, false);
    }

    @Disabled
    @Test
    void whatEverThisIs() throws Exception {
        byte[] bytes = "\\3436303230303231303167\r".getBytes(StandardCharsets.UTF_8);
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
        assertThat(msg).isNotNull();
        System.out.println(msg);
        assertMessageMatches(bytes, msg);
    }

    @Test
    void deviceManagementInstruction() throws Exception {
        byte[] bytes = "@1A2001\r".getBytes(StandardCharsets.UTF_8);
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
        assertThat(msg).isNotNull();
        System.out.println(msg);
        assertMessageMatches(bytes, msg);
    }

    @Test
    void setLight() throws Exception {
        byte[] bytes = "\\0538000100g\r".getBytes(StandardCharsets.UTF_8);
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
        assertThat(msg).isNotNull();
        System.out.println(msg);
        CBusMessageToServer msgToServer = (CBusMessageToServer) msg;
        RequestCommand requestCommand = (RequestCommand) msgToServer.getRequest();
        CBusCommand cbusCommand = requestCommand.getCbusCommand();
        System.out.println(cbusCommand);
        assertMessageMatches(bytes, msg);
    }

    @Test
    void identifyResponse() throws Exception {
        byte[] bytes = "g.890150435F434E49454421\r\n".getBytes(StandardCharsets.UTF_8);
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        // We know we send an identify command so we set the cal flag
        requestContext = new RequestContext(true, false, false);
        cBusOptions = C_BUS_OPTIONS_WITH_SRCHK;
        CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
        assertThat(msg).isNotNull();
        System.out.println(msg);
        CBusMessageToClient messageToClient = (CBusMessageToClient) msg;
        ReplyOrConfirmationConfirmation confirmationReply = (ReplyOrConfirmationConfirmation) messageToClient.getReply();
        ReplyOrConfirmationReply normalReply = (ReplyOrConfirmationReply) confirmationReply.getEmbeddedReply();
        EncodedReplyCALReply encodedReplyCALReply = (EncodedReplyCALReply) ((ReplyEncodedReply) normalReply.getReply()).getEncodedReply();
        System.out.println(encodedReplyCALReply.getCalReply());
        assertMessageMatches(bytes, msg);
    }

    @Test
    void someResponse() throws Exception {
        byte[] bytes = "nl.8220025C\r\n".getBytes(StandardCharsets.UTF_8);
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        requestContext = new RequestContext(true, false, false);
        cBusOptions = C_BUS_OPTIONS_WITH_SRCHK;
        CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
        assertThat(msg).isNotNull();
        System.out.println(msg);
        CBusMessageToClient messageToClient = (CBusMessageToClient) msg;
        ReplyOrConfirmationConfirmation confirmationReply = (ReplyOrConfirmationConfirmation) messageToClient.getReply();
        ReplyOrConfirmationReply normalReply = (ReplyOrConfirmationReply) confirmationReply.getEmbeddedReply();
        EncodedReplyCALReply encodedReplyCALReply = (EncodedReplyCALReply) ((ReplyEncodedReply) normalReply.getReply()).getEncodedReply();
        System.out.println(encodedReplyCALReply.getCalReply());
        assertMessageMatches(bytes, msg);
    }

    @Test
    void someOtherResponse() throws Exception {
        byte[] bytes = "\\0538000100g\r".getBytes(StandardCharsets.UTF_8);
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
        assertThat(msg).isNotNull();
        System.out.println(msg);
        CBusMessageToServer messageToServer = (CBusMessageToServer) msg;
        RequestCommand requestCommand = (RequestCommand) messageToServer.getRequest();
        System.out.println(requestCommand.getCbusCommand());
        assertMessageMatches(bytes, msg);
    }


    @Test
    void identifyRequest2() throws Exception {
        byte[] bytes = "21021A2102i\r".getBytes(StandardCharsets.UTF_8);
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
        assertThat(msg).isNotNull();
        System.out.println(msg);
        CALData calData = ((RequestObsolete) ((CBusMessageToServer) msg).getRequest()).getCalData();
        System.out.println(calData);
        assertMessageMatches(bytes, msg);
    }

    @Test
    void identifyResponse2() throws Exception {
        byte[] bytes = "i.8902352E342E3030202010\r\n".getBytes(StandardCharsets.UTF_8);
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        // We know we send an identify command so we set the cal flag
        requestContext = new RequestContext(true, false, true);
        cBusOptions = C_BUS_OPTIONS_WITH_SRCHK;
        CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
        assertThat(msg).isNotNull();
        System.out.println(msg);
        CBusMessageToClient messageToClient = (CBusMessageToClient) msg;
        ReplyOrConfirmationConfirmation confirmationReply = (ReplyOrConfirmationConfirmation) messageToClient.getReply();
        ReplyOrConfirmationReply normalReply = (ReplyOrConfirmationReply) confirmationReply.getEmbeddedReply();
        EncodedReplyCALReply encodedReplyCALReply = (EncodedReplyCALReply) ((ReplyEncodedReply) normalReply.getReply()).getEncodedReply();
        System.out.println(encodedReplyCALReply.getCalReply());
        assertMessageMatches(bytes, msg);
    }

    @Test
    void recall() throws Exception {
        byte[] bytes = "@1A2001\r".getBytes(StandardCharsets.UTF_8);
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
        assertThat(msg).isNotNull();
        System.out.println(msg);
        CBusMessageToServer messageToServer = (CBusMessageToServer) msg;
        RequestDirectCommandAccess requestDirectCommandAccess = (RequestDirectCommandAccess) messageToServer.getRequest();
        CALData calData = ((RequestDirectCommandAccess) ((CBusMessageToServer) msg).getRequest()).getCalData();
        System.out.println(calData);

        WriteBufferByteBased writeBuffer = new WriteBufferByteBased(bytes.length);
        msg.serialize(writeBuffer);
        assertThat(writeBuffer.getBytes()).isEqualTo(bytes);
        assertMessageMatches(bytes, msg);
    }

    @Test
    void identifyTypeReply() throws Exception {
        byte[] bytes = "h.890150435F434E49454421\r\n".getBytes(StandardCharsets.UTF_8);
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        requestContext = new RequestContext(true, false, true);
        cBusOptions = C_BUS_OPTIONS_WITH_SRCHK;
        CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
        assertThat(msg).isNotNull();
        System.out.println(msg);
        CBusMessageToClient messageToClient = (CBusMessageToClient) msg;
        ReplyOrConfirmationConfirmation confirmationReply = (ReplyOrConfirmationConfirmation) messageToClient.getReply();
        ReplyOrConfirmationReply normalReply = (ReplyOrConfirmationReply) confirmationReply.getEmbeddedReply();
        EncodedReplyCALReply encodedReplyCALReply = (EncodedReplyCALReply) ((ReplyEncodedReply) normalReply.getReply()).getEncodedReply();
        System.out.println(encodedReplyCALReply.getCalReply());

        assertMessageMatches(bytes, msg);
    }

    @Test
    void write30to9755() throws Exception {
        byte[] bytes = "A3309755s\r".getBytes(StandardCharsets.UTF_8);
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        cBusOptions = C_BUS_OPTIONS_WITH_SRCHK;
        CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
        assertThat(msg).isNotNull();
        System.out.println(msg);
        System.out.println(((RequestObsolete) ((CBusMessageToServer) msg).getRequest()).getCalData());

        assertMessageMatches(bytes, msg);
    }

    @Test
    void strangeNotYetParsableCommandResponse() throws Exception {
        byte[] bytes = "s.860202003230977D\r\n".getBytes(StandardCharsets.UTF_8);
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        requestContext = new RequestContext(true, false, false);
        cBusOptions = C_BUS_OPTIONS_WITH_SRCHK;
        CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
        assertThat(msg).isNotNull();
        System.out.println(msg);
        CBusMessageToClient messageToClient = (CBusMessageToClient) msg;
        ReplyOrConfirmationConfirmation confirmationReply = (ReplyOrConfirmationConfirmation) messageToClient.getReply();
        ReplyOrConfirmationReply normalReply = (ReplyOrConfirmationReply) confirmationReply.getEmbeddedReply();
        EncodedReplyCALReply encodedReplyCALReply = (EncodedReplyCALReply) ((ReplyEncodedReply) normalReply.getReply()).getEncodedReply();
        System.out.println(encodedReplyCALReply.getCalReply());

        assertMessageMatches(bytes, msg);
    }

    @Test
    void statusRequestBinaryState() throws Exception {
        byte[] bytes = "\\05FF00FAFF00v\r".getBytes(StandardCharsets.UTF_8);
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, false);
        CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
        assertThat(msg).isNotNull();
        System.out.println(msg);
        CBusMessageToServer messageToServer = (CBusMessageToServer) msg;
        RequestCommand requestCommand = (RequestCommand) messageToServer.getRequest();
        CBusCommand cbusCommand = requestCommand.getCbusCommand();
        System.out.println(cbusCommand);
        assertMessageMatches(bytes, msg);
    }

    @Disabled
    @Test
    void wat() throws Exception {
        byte[] bytes = "D8FF0024000002000000000000000008000000000000000000\r\n".getBytes(StandardCharsets.UTF_8);
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        requestContext = new RequestContext(true, false, false);
        CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
        assertThat(msg).isNotNull();
        System.out.println(msg);
        CBusMessageToClient messageToClient = (CBusMessageToClient) msg;
        ReplyOrConfirmationReply reply = (ReplyOrConfirmationReply) messageToClient.getReply();
        EncodedReplyStandardFormatStatusReply encodedReplyCALReply = (EncodedReplyStandardFormatStatusReply) ((ReplyEncodedReply) reply.getReply()).getEncodedReply();
        System.out.println(encodedReplyCALReply.getReply());
        assertMessageMatches(bytes, msg);
    }

    @Test
    void WriteCommand() throws Exception {
        byte[] bytes = "\\46310900A400410600r\r".getBytes(StandardCharsets.UTF_8);
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, false);
        CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
        assertThat(msg).isNotNull();
        System.out.println(msg);
        System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
        assertMessageMatches(bytes, msg);
    }

    @Test
    void statusReply() throws Exception {
        byte[] bytes = "D8FF5800000000000000000000000000000000000000000000D1\r\n".getBytes(StandardCharsets.UTF_8);
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        cBusOptions = C_BUS_OPTIONS_WITH_SRCHK;
        requestContext = new RequestContext(false, true, false);
        CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
        assertThat(msg).isNotNull();
        System.out.println(msg);
        Reply normalReply = ((ReplyOrConfirmationReply) ((CBusMessageToClient) msg).getReply()).getReply();
        EncodedReplyStandardFormatStatusReply encodedReplyStandardFormatStatusReply = (EncodedReplyStandardFormatStatusReply) ((ReplyEncodedReply) normalReply).getEncodedReply();
        System.out.println(encodedReplyStandardFormatStatusReply.getReply());

        assertMessageMatches(bytes, msg);
    }

    @Test
    void identifyUnitSummary() throws Exception {
        byte[] bytes = "2110\r".getBytes(StandardCharsets.UTF_8);
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        cBusOptions = C_BUS_OPTIONS_WITH_SRCHK;
        CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
        assertThat(msg).isNotNull();
        System.out.println(msg);
        System.out.println(((RequestObsolete) ((CBusMessageToServer) msg).getRequest()).getCalData());

        assertMessageMatches(bytes, msg);
    }

    @Test
    void identifyUnitSummaryResponse() throws Exception {
        byte[] bytes = "o.8510020000FF6A\r\n".getBytes(StandardCharsets.UTF_8);
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        cBusOptions = C_BUS_OPTIONS_WITH_SRCHK;
        requestContext = new RequestContext(true, false, true);
        CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
        assertThat(msg).isNotNull();
        System.out.println(msg);
        System.out.println(((ReplyEncodedReply) ((ReplyOrConfirmationReply) ((ReplyOrConfirmationConfirmation) ((CBusMessageToClient) msg).getReply()).getEmbeddedReply()).getReply()).getEncodedReply());

        assertMessageMatches(bytes, msg);
    }

    @Test
    void hvacAndCoolingSAL() throws Exception {
        byte[] bytes = "0531AC0079042F0401430316000011\r\n".getBytes(StandardCharsets.UTF_8);
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        cBusOptions = C_BUS_OPTIONS_WITH_SRCHK;
        requestContext = new RequestContext(false, false, false);
        CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
        assertThat(msg).isNotNull();
        System.out.println(msg);
        System.out.println(((ReplyEncodedReply) ((ReplyOrConfirmationReply) ((CBusMessageToClient) msg).getReply()).getReply()).getEncodedReply());

        assertMessageMatches(bytes, msg);
    }

    @Test
    void calIdentifyReplyAndAnotherCal() throws Exception {
        byte[] bytes = "h.860102008902312E362E30302020832138FFAE\r\n".getBytes(StandardCharsets.UTF_8);
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        cBusOptions = C_BUS_OPTIONS_WITH_SRCHK;
        requestContext = new RequestContext(false, false, true);
        CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
        assertThat(msg).isNotNull();
        System.out.println(msg);
        System.out.println(((ReplyEncodedReply) ((ReplyOrConfirmationReply) ((ReplyOrConfirmationConfirmation) ((CBusMessageToClient) msg).getReply()).getEmbeddedReply()).getReply()).getEncodedReply());

        assertMessageMatches(bytes, msg);
    }

    @Test
    void routedAcknowledge() throws Exception {
        byte[] bytes = "r.8631020100320041D3\r\n".getBytes(StandardCharsets.UTF_8);
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        cBusOptions = C_BUS_OPTIONS_WITH_SRCHK;
        requestContext = new RequestContext(false, false, true);
        CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
        assertThat(msg).isNotNull();
        System.out.println(msg);
        System.out.println(((ReplyEncodedReply) ((ReplyOrConfirmationReply) ((ReplyOrConfirmationConfirmation) ((CBusMessageToClient) msg).getReply()).getEmbeddedReply()).getReply()).getEncodedReply());

        assertMessageMatches(bytes, msg);
    }

    @Test
    void gavValuesCurrentReply() throws Exception {
        byte[] bytes = "w.860C02008A08000000C8000000000012\r\n".getBytes(StandardCharsets.UTF_8);
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        cBusOptions = C_BUS_OPTIONS_WITH_SRCHK;
        requestContext = new RequestContext(false, false, true);
        CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
        assertThat(msg).isNotNull();
        System.out.println(msg);
        System.out.println(((ReplyEncodedReply) ((ReplyOrConfirmationReply) ((ReplyOrConfirmationConfirmation) ((CBusMessageToClient) msg).getReply()).getEmbeddedReply()).getReply()).getEncodedReply());

        assertMessageMatches(bytes, msg);
    }

    @Test
    void SetHvacLevel() throws Exception {
        byte[] bytes = "0531AC0036040108FF0000DC\r\n".getBytes(StandardCharsets.UTF_8);
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        requestContext = new RequestContext(false, false, false);
        CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
        assertThat(msg).isNotNull();
        System.out.println(msg);
        System.out.println(((ReplyEncodedReply) ((ReplyOrConfirmationReply) ((CBusMessageToClient) msg).getReply()).getReply()).getEncodedReply());

        // TODO: apparently the set the first bit of AuxiliaryLevel to true wich is not valid according to the documentation
        //assertMessageMatches(bytes, msg);
    }
    @Test
    void salHvac() throws Exception {
        byte[] bytes = "0531AC0036040142037F001F\r\n".getBytes(StandardCharsets.UTF_8);
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        cBusOptions = C_BUS_OPTIONS_WITH_SRCHK;
        requestContext = new RequestContext(false, false, false);
        CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
        assertThat(msg).isNotNull();
        System.out.println(msg);
        System.out.println(((ReplyEncodedReply) ((ReplyOrConfirmationReply) ((CBusMessageToClient) msg).getReply()).getReply()).getEncodedReply());

        assertMessageMatches(bytes, msg);
    }

    @Disabled("Not clear yet what this is")
    @Test
    void closestFitIsAStatusRequestButWeDonTHaveAnyBytesBeforeThat() throws Exception {
        byte[] bytes = "FAFF00r\r".getBytes(StandardCharsets.UTF_8);
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        cBusOptions = C_BUS_OPTIONS_WITH_SRCHK;
        CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
        assertThat(msg).isNotNull();
        System.out.println(msg);

        assertMessageMatches(bytes, msg);
    }

    @Test
    void ownSal() throws Exception {
        byte[] bytes = "003809AF10\r\n".getBytes(StandardCharsets.UTF_8);
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        cBusOptions = C_BUS_OPTIONS_WITH_SRCHK;
        CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
        assertThat(msg).isNotNull();
        System.out.println(msg);
        System.out.println(((ReplyEncodedReply) ((ReplyOrConfirmationReply) ((CBusMessageToClient) msg).getReply()).getReply()).getEncodedReply());

        assertMessageMatches(bytes, msg);
    }
}

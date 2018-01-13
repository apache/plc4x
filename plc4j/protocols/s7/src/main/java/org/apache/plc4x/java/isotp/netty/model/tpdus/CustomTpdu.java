package org.apache.plc4x.java.isotp.netty.model.tpdus;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.isotp.netty.model.params.Parameter;
import org.apache.plc4x.java.isotp.netty.model.types.TpduCode;

import java.util.List;

public class CustomTpdu extends Tpdu  {

    public CustomTpdu(Byte tpduCode, List<Parameter> parameters, ByteBuf userData) {
        super(TpduCode.valueOf(tpduCode), parameters, userData);
    }

}

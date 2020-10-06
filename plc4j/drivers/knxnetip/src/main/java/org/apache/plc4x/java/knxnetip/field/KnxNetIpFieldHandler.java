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
package org.apache.plc4x.java.knxnetip.field;

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.api.value.PlcValues;
import org.apache.plc4x.java.spi.connection.DefaultPlcFieldHandler;

import java.math.BigDecimal;
import java.math.BigInteger;

public class KnxNetIpFieldHandler extends DefaultPlcFieldHandler {

    @Override
    public PlcField createField(String fieldQuery) {
        if (KnxNetIpField.matches(fieldQuery)) {
            return KnxNetIpField.of(fieldQuery);
        }
        throw new PlcInvalidFieldException(fieldQuery);
    }

    @Override
    public PlcValue encodeBoolean(PlcField field, Object[] values) {
        if(values.length == 1) {
            return PlcValues.of((Boolean) values[0]);
        }
        return PlcValues.of((Boolean[]) values);
    }

    @Override
    public PlcValue encodeByte(PlcField field, Object[] values) {
        if(values.length == 1) {
            return PlcValues.of((Byte) values[0]);
        }
        return PlcValues.of((Byte[]) values);
    }

    @Override
    public PlcValue encodeShort(PlcField field, Object[] values) {
        if(values.length == 1) {
            return PlcValues.of((Short) values[0]);
        }
        return PlcValues.of((Short[]) values);
    }

    @Override
    public PlcValue encodeInteger(PlcField field, Object[] values) {
        if(values.length == 1) {
            return PlcValues.of((Integer) values[0]);
        }
        return PlcValues.of((Integer[]) values);
    }

    @Override
    public PlcValue encodeBigInteger(PlcField field, Object[] values) {
        if(values.length == 1) {
            return PlcValues.of((BigInteger) values[0]);
        }
        return PlcValues.of((BigInteger[]) values);
    }

    @Override
    public PlcValue encodeLong(PlcField field, Object[] values) {
        if(values.length == 1) {
            return PlcValues.of((Long) values[0]);
        }
        return PlcValues.of((Long[]) values);
    }

    @Override
    public PlcValue encodeFloat(PlcField field, Object[] values) {
        if(values.length == 1) {
            return PlcValues.of((Float) values[0]);
        }
        return PlcValues.of((Float[]) values);
    }

    @Override
    public PlcValue encodeBigDecimal(PlcField field, Object[] values) {
        if(values.length == 1) {
            return PlcValues.of((BigDecimal) values[0]);
        }
        return PlcValues.of((BigDecimal[]) values);
    }

    @Override
    public PlcValue encodeDouble(PlcField field, Object[] values) {
        if(values.length == 1) {
            return PlcValues.of((Double) values[0]);
        }
        return PlcValues.of((Double[]) values);
    }

    @Override
    public PlcValue encodeString(PlcField field, Object[] values) {
        if(values.length == 1) {
            return PlcValues.of((String) values[0]);
        }
        return PlcValues.of((String[]) values);
    }

}

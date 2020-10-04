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
package org.apache.plc4x.java.opcua.protocol;


import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.value.*;
import org.apache.plc4x.java.spi.connection.DefaultPlcFieldHandler;

import java.math.BigInteger;
import java.util.ArrayList;

/**
 */
public class OpcuaPlcFieldHandler extends DefaultPlcFieldHandler {

    @Override
    public PlcField createField(String fieldQuery) {
        if (OpcuaField.matches(fieldQuery)) {
            return OpcuaField.of(fieldQuery);
        }
        throw new PlcInvalidFieldException(fieldQuery);
    }

    @Override
    public PlcValue encodeString(PlcField field, Object[] values) {
        OpcuaField opcField = (OpcuaField) field;
        ArrayList<String> resultSet = new ArrayList<>();
        for (Object item : values) {
            resultSet.add(item.toString());
        }
        if(resultSet.size() == 1) {
            return new PlcString(resultSet.get(0));
        } else {
            return new PlcList(resultSet);
        }
    }

    @Override
    public PlcValue encodeBoolean(PlcField field, Object[] values) {
        return internalEncode(field, values, "BOOL");
    }

    @Override
    public PlcValue encodeByte(PlcField field, Object[] values) {
        return internalEncode(field, values, "SINT");
    }

    @Override
    public PlcValue encodeShort(PlcField field, Object[] values) {
        return internalEncode(field, values, "INT");
    }

    @Override
    public PlcValue encodeInteger(PlcField field, Object[] values) {
        return internalEncode(field, values, "DINT");
    }

    @Override
    public PlcValue encodeBigInteger(PlcField field, Object[] values) {
        OpcuaField opcField = (OpcuaField) field;
        ArrayList<BigInteger> resultSet = new ArrayList<>();
        for (Object item : values) {
            resultSet.add((BigInteger) item);
        }
        if(resultSet.size() == 1) {
            return new PlcBigInteger(resultSet.get(0));
        } else {
            return new PlcList(resultSet);
        }
    }

    @Override
    public PlcValue encodeLong(PlcField field, Object[] values) {
        return internalEncode(field, values, "LINT");
    }

    @Override
    public PlcValue encodeFloat(PlcField field, Object[] values) {
        return internalEncode(field, values, "REAL");
    }


    @Override
    public PlcValue encodeDouble(PlcField field, Object[] values) {
        OpcuaField opcField = (OpcuaField) field;
        ArrayList<Double> resultSet = new ArrayList<>();
        for (Object item : values) {
            resultSet.add((Double) item);
        }
        if(resultSet.size() == 1) {
            return new PlcLREAL(resultSet.get(0));
        } else {
            return new PlcList(resultSet);
        }
    }

    private PlcValue internalEncode(PlcField field, Object[] values, String datatype) {
        OpcuaField opcField = (OpcuaField) field;
        try {
            switch (datatype) {
                //Implement Custom PlcValue types here
                default:
                    return PlcValues.of(values, Class.forName(PlcValues.class.getPackage().getName() + ".Plc" + datatype));
            }
        } catch (ClassNotFoundException e) {
            throw new PlcRuntimeException("Invalid encoder for type " + datatype + e);
        }
    }

}

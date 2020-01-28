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
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.value.*;
import org.apache.plc4x.java.spi.connection.DefaultPlcFieldHandler;

import java.math.BigInteger;
import java.util.ArrayList;

/**
 * Created by Matthias Milan Strljic on 10.05.2019
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
        OpcuaField opcField = (OpcuaField) field;
        ArrayList<Boolean> resultSet = new ArrayList<>();
        for (Object item : values) {
            resultSet.add((Boolean) item);
        }
        if(resultSet.size() == 1) {
            return new PlcBoolean(resultSet.get(0));
        } else {
            return new PlcList(resultSet);
        }
    }

    @Override
    public PlcValue encodeByte(PlcField field, Object[] values) {
        OpcuaField opcField = (OpcuaField) field;
        ArrayList<Byte> resultSet = new ArrayList<>();
        for (Object item : values) {
            resultSet.add((Byte) item);
        }
        if(resultSet.size() == 1) {
            return new PlcInteger(resultSet.get(0));
        } else {
            return new PlcList(resultSet);
        }
    }

    @Override
    public PlcValue encodeShort(PlcField field, Object[] values) {
        OpcuaField opcField = (OpcuaField) field;
        ArrayList<Short> resultSet = new ArrayList<>();
        for (Object item : values) {
            resultSet.add((Short) item);
        }
        if(resultSet.size() == 1) {
            return new PlcInteger(resultSet.get(0));
        } else {
            return new PlcList(resultSet);
        }
    }

    @Override
    public PlcValue encodeInteger(PlcField field, Object[] values) {
        OpcuaField opcField = (OpcuaField) field;
        ArrayList<Integer> resultSet = new ArrayList<>();
        for (Object item : values) {
            resultSet.add((Integer) item);
        }
        if(resultSet.size() == 1) {
            return new PlcInteger(resultSet.get(0));
        } else {
            return new PlcList(resultSet);
        }
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
        OpcuaField opcField = (OpcuaField) field;
        ArrayList<Long> resultSet = new ArrayList<>();
        for (Object item : values) {
            resultSet.add((Long) item);
        }
        if(resultSet.size() == 1) {
            return new PlcLong(resultSet.get(0));
        } else {
            return new PlcList(resultSet);
        }
    }

    @Override
    public PlcValue encodeFloat(PlcField field, Object[] values) {
        OpcuaField opcField = (OpcuaField) field;
        ArrayList<Float> resultSet = new ArrayList<>();
        for (Object item : values) {
            resultSet.add((Float) item);
        }
        if(resultSet.size() == 1) {
            return new PlcFloat(resultSet.get(0));
        } else {
            return new PlcList(resultSet);
        }
    }


    @Override
    public PlcValue encodeDouble(PlcField field, Object[] values) {
        OpcuaField opcField = (OpcuaField) field;
        ArrayList<Double> resultSet = new ArrayList<>();
        for (Object item : values) {
            resultSet.add((Double) item);
        }
        if(resultSet.size() == 1) {
            return new PlcDouble(resultSet.get(0));
        } else {
            return new PlcList(resultSet);
        }
    }

}

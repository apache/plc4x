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
package org.apache.plc4x.language.python;

import org.apache.commons.lang3.StringUtils;
import org.apache.commons.lang3.math.NumberUtils;
import org.apache.commons.text.WordUtils;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.terms.DefaultStringLiteral;
import org.apache.plc4x.plugins.codegenerator.protocol.freemarker.BaseFreemarkerLanguageTemplateHelper;
import org.apache.plc4x.plugins.codegenerator.protocol.freemarker.FreemarkerException;
import org.apache.plc4x.plugins.codegenerator.protocol.freemarker.Tracer;
import org.apache.plc4x.plugins.codegenerator.types.definitions.*;
import org.apache.plc4x.plugins.codegenerator.types.fields.*;
import org.apache.plc4x.plugins.codegenerator.types.references.*;
import org.apache.plc4x.plugins.codegenerator.types.terms.*;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.*;
import java.util.function.Function;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

@SuppressWarnings({"unused", "WeakerAccess"})
public class PythonLanguageTemplateHelper extends BaseFreemarkerLanguageTemplateHelper {

    private static final Logger LOGGER = LoggerFactory.getLogger(PythonLanguageTemplateHelper.class);

    private final Map<String, String> options;

    public SortedSet<String> requiredImports = new TreeSet<>();

    public SortedSet<String> requiredImportsForDataIo = new TreeSet<>();

    public PythonLanguageTemplateHelper(TypeDefinition thisType, String protocolName, String flavorName, Map<String, TypeDefinition> types,
                                      Map<String, String> options) {
        super(thisType, protocolName, flavorName, types);
        this.options = options;
    }

    public String packageName() {
        return packageName(protocolName, "python", flavorName);
    }

    public String packageName(String protocolName, String languageName, String languageFlavorName) {
        return Optional.ofNullable(options.get("package")).orElseGet(() ->
                String.join("", protocolName.split("-")) + "." +
                String.join("", languageFlavorName.split("-")));
    }

    public String packageName(String languageFlavorName) {
        return String.join("", languageFlavorName.split("\\-"));
    }

    @Override
    public String getLanguageTypeNameForField(Field field) {
        // If the referenced type is a DataIo type, the value is of type PlcValue.
        if (field.isPropertyField()) {
            PropertyField propertyField = field.asPropertyField().orElseThrow(IllegalStateException::new);
            if (propertyField.getType().isComplexTypeReference()) {
                ComplexTypeReference complexTypeReference = propertyField.getType().asComplexTypeReference().orElseThrow(IllegalStateException::new);
                final TypeDefinition typeDefinition = getTypeDefinitions().get(complexTypeReference.getName());
                if (typeDefinition instanceof DataIoTypeDefinition) {
                    emitRequiredImport(complexTypeReference.getName());
                    return "PlcValue";
                }
            }
        }
        return getLanguageTypeNameForTypeReference(((TypedField) field).getType(), !field.isOptionalField());
    }

    public String getNonPrimitiveLanguageTypeNameForField(TypedField field) {
        return getLanguageTypeNameForTypeReference(field.getType(), false);
    }

    public String getLanguageTypeNameForSpecType(TypeReference typeReference) {
        return getLanguageTypeNameForTypeReference(typeReference, true);
    }

    @Override
    public String getLanguageTypeNameForTypeReference(TypeReference typeReference) {
        return getLanguageTypeNameForTypeReference(typeReference, false);
    }

    public String getLanguageTypeNameForTypeReference(TypeReference typeReference, boolean allowPrimitive) {
        if (typeReference == null) {
            // TODO: shouldn't this be an error case
            return "";
        }
        if (typeReference.isArrayTypeReference()) {
            final ArrayTypeReference arrayTypeReference = (ArrayTypeReference) typeReference;
            TypeReference elementTypeReference = arrayTypeReference.getElementTypeReference();
            emitRequiredImport("from typing import List");
            return "List[" + getLanguageTypeNameForTypeReference(elementTypeReference) + "]";
        }
        if (typeReference.isNonSimpleTypeReference()) {
            emitRequiredImport("from plc4py.protocols." + protocolName + "." + flavorName.replace("-", "") + "." + typeReference.asNonSimpleTypeReference().orElseThrow().getName() + " import " + typeReference.asNonSimpleTypeReference().orElseThrow().getName());
            return typeReference.asNonSimpleTypeReference().orElseThrow().getName();
        }
        SimpleTypeReference simpleTypeReference = typeReference.asSimpleTypeReference().orElseThrow();
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                emitRequiredImport("from ctypes import c_bool");
                return "c_bool";
            case BYTE:
                emitRequiredImport("from ctypes import c_byte");
                return "c_byte";
            case UINT:
                IntegerTypeReference unsignedIntegerTypeReference = simpleTypeReference.asIntegerTypeReference().orElseThrow();
                if (unsignedIntegerTypeReference.getSizeInBits() <= 8) {
                    emitRequiredImport("from ctypes import c_uint8");
                    return "c_uint8";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 16) {
                    emitRequiredImport("from ctypes import c_uint16");
                    return "c_uint16";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 32) {
                    emitRequiredImport("from ctypes import c_uint32");
                    return "c_uint32";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 64) {
                    emitRequiredImport("from ctypes import c_uint64");
                    return "c_uint64";
                }
                emitRequiredImport("from ctypes import c_longlong");
                return "c_longlong";
            case INT:
                IntegerTypeReference integerTypeReference = simpleTypeReference.asIntegerTypeReference().orElseThrow();
                if (integerTypeReference.getSizeInBits() <= 8) {
                    emitRequiredImport("from ctypes import c_int8");
                    return "c_int8";
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    emitRequiredImport("from ctypes import c_int16");
                    return "c_int16";
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    emitRequiredImport("from ctypes import c_int32");
                    return "c_int32";
                }
                if (integerTypeReference.getSizeInBits() <= 64) {
                    emitRequiredImport("from ctypes import c_int64");
                    return "c_int64";
                }
                emitRequiredImport("from ctypes import c_longlong");
                return "c_longlong";
            case FLOAT:
            case UFLOAT:
                FloatTypeReference floatTypeReference = simpleTypeReference.asFloatTypeReference().orElseThrow();
                int sizeInBits = floatTypeReference.getSizeInBits();
                if (sizeInBits <= 32) {
                    emitRequiredImport("from ctypes import c_float");
                    return "c_float";
                }
                if (sizeInBits <= 64) {
                    emitRequiredImport("from ctypes import c_double");
                    return "c_double";
                }
                emitRequiredImport("from ctypes import c_longdouble");
                return "c_longdouble";
            case STRING:
            case VSTRING:
                return "str";
            case TIME:
            case DATE:
            case DATETIME:
                return "time.Time";
            default:
                throw new RuntimeException("Unsupported simple type");
        }
    }

    public String getPlcValueTypeForTypeReference(TypeReference typeReference) {
        if (!(typeReference instanceof SimpleTypeReference)) {
            return "PlcStruct";
        }
        SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return "PlcBOOL";
            case BYTE:
                return "PlcSINT";
            case UINT:
                IntegerTypeReference unsignedIntegerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (unsignedIntegerTypeReference.getSizeInBits() <= 4) {
                    return "PlcUSINT";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 8) {
                    return "PlcUINT";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 16) {
                    return "PlcUDINT";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 32) {
                    return "PlcULINT";
                }
            case INT:
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (integerTypeReference.getSizeInBits() <= 8) {
                    return "PlcSINT";
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    return "PlcINT";
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    return "PlcDINT";
                }
                if (integerTypeReference.getSizeInBits() <= 64) {
                    return "PlcLINT";
                }

            case FLOAT:
            case UFLOAT:
                FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                int sizeInBits = floatTypeReference.getSizeInBits();
                if (sizeInBits <= 32) {
                    return "PlcREAL";
                }
                if (sizeInBits <= 64) {
                    return "PlcLREAL";
                }
            case STRING:
            case VSTRING:
                return "PlcSTRING";
            case TIME:
            case DATE:
            case DATETIME:
                return "PlcTIME";
        }
        throw new RuntimeException("Unsupported simple type");
    }

    @Override
    public String getNullValueForTypeReference(TypeReference typeReference) {
        if (typeReference instanceof SimpleTypeReference) {
            SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
            switch (simpleTypeReference.getBaseType()) {
                case BIT:
                    return "False";
                case BYTE:
                    return "0";
                case UINT:
                    IntegerTypeReference unsignedIntegerTypeReference = (IntegerTypeReference) simpleTypeReference;
                    if (unsignedIntegerTypeReference.getSizeInBits() <= 16) {
                        return "0";
                    }
                    if (unsignedIntegerTypeReference.getSizeInBits() <= 32) {
                        return "0l";
                    }
                    return "None";
                case INT:
                    IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                    if (integerTypeReference.getSizeInBits() <= 32) {
                        return "0";
                    }
                    if (integerTypeReference.getSizeInBits() <= 64) {
                        return "0l";
                    }
                    return "None";
                case FLOAT:
                    FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                    int sizeInBits = floatTypeReference.getSizeInBits();
                    if (sizeInBits <= 32) {
                        return "0.0f";
                    }
                    if (sizeInBits <= 64) {
                        return "0.0";
                    }
                    return "None";
                case STRING:
                case VSTRING:
                    return "None";
            }
            throw new FreemarkerException("Unmapped base-type" + simpleTypeReference.getBaseType());
        } else {
            return "None";
        }
    }

    public int getNumBits(SimpleTypeReference simpleTypeReference) {
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return 1;
            case BYTE:
                return 8;
            case UINT:
            case INT:
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                return integerTypeReference.getSizeInBits();
            case FLOAT:
                FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                return floatTypeReference.getSizeInBits();
            case STRING:
                StringTypeReference stringTypeReference = (StringTypeReference) simpleTypeReference;
                return stringTypeReference.getSizeInBits();
            case VSTRING:
                throw new IllegalArgumentException("getSizeInBits doesn't work for 'vstring' fields");
            default:
                return 0;
        }
    }

    @Deprecated
    @Override
    public String getReadBufferReadMethodCall(SimpleTypeReference simpleTypeReference, String valueString, TypedField field) {
        return "/*TODO: migrate me*/" + getReadBufferReadMethodCall("", simpleTypeReference, valueString, field);
    }

    @Deprecated
    public String getReadBufferReadMethodCall(String logicalName, SimpleTypeReference simpleTypeReference, String valueString, TypedField field) {
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                String bitType = "Bit";
                return "/*TODO: migrate me*/" + "readBuffer.read" + bitType + "(\"" + logicalName + "\")";
            case BYTE:
                String byteType = "Byte";
                return "/*TODO: migrate me*/" + "readBuffer.read" + byteType + "(\"" + logicalName + "\")";
            case UINT:
                String unsignedIntegerType;
                IntegerTypeReference unsignedIntegerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (unsignedIntegerTypeReference.getSizeInBits() <= 4) {
                    unsignedIntegerType = "UnsignedByte";
                } else if (unsignedIntegerTypeReference.getSizeInBits() <= 8) {
                    unsignedIntegerType = "UnsignedShort";
                } else if (unsignedIntegerTypeReference.getSizeInBits() <= 16) {
                    unsignedIntegerType = "UnsignedInt";
                } else if (unsignedIntegerTypeReference.getSizeInBits() <= 32) {
                    unsignedIntegerType = "UnsignedLong";
                } else {
                    unsignedIntegerType = "UnsignedBigInteger";
                }
                return "/*TODO: migrate me*/" + "readBuffer.read" + unsignedIntegerType + "(\"" + logicalName + "\", " + simpleTypeReference.getSizeInBits() + ")";
            case INT:
                String integerType;
                if (simpleTypeReference.getSizeInBits() <= 8) {
                    integerType = "SignedByte";
                } else if (simpleTypeReference.getSizeInBits() <= 16) {
                    integerType = "Short";
                } else if (simpleTypeReference.getSizeInBits() <= 32) {
                    integerType = "Int";
                } else if (simpleTypeReference.getSizeInBits() <= 64) {
                    integerType = "Long";
                } else {
                    integerType = "BigInteger";
                }
                return "/*TODO: migrate me*/" + "readBuffer.read" + integerType + "(\"" + logicalName + "\", " + simpleTypeReference.getSizeInBits() + ")";
            case FLOAT:
                String floatType = (simpleTypeReference.getSizeInBits() <= 32) ? "Float" : "Double";
                return "/*TODO: migrate me*/" + "readBuffer.read" + floatType + "(\"" + logicalName + "\", " + simpleTypeReference.getSizeInBits() + ")";
            case STRING:
            case VSTRING:
                String stringType = "String";
                final Term encodingTerm = field.getEncoding().orElse(new DefaultStringLiteral("UTF-8"));
                if (!(encodingTerm instanceof StringLiteral)) {
                    throw new RuntimeException("Encoding must be a quoted string value");
                }
                String encoding = ((StringLiteral) encodingTerm).getValue();
                String length = Integer.toString(simpleTypeReference.getSizeInBits());
                if (simpleTypeReference.getBaseType() == SimpleTypeReference.SimpleBaseType.VSTRING) {
                    VstringTypeReference vstringTypeReference = (VstringTypeReference) simpleTypeReference;
                    length = toParseExpression(field, INT_TYPE_REFERENCE, vstringTypeReference.getLengthExpression(), null);
                }
                return "/*TODO: migrate me*/" + "readBuffer.read" + stringType + "(\"" + logicalName + "\", " + length + ", \"" +
                    encoding + "\")";
        }
        return "/*TODO: migrate me*/" + "";
    }

    public String getDataReaderCall(TypeReference typeReference) {
        return getDataReaderCall(typeReference, "enumForValue");
    }

    public String getDataReaderCall(TypeReference typeReference, String resolverMethod) {
        if (typeReference.isEnumTypeReference()) {
            final String languageTypeName = getLanguageTypeNameForTypeReference(typeReference);
            final SimpleTypeReference enumBaseTypeReference = getEnumBaseTypeReference(typeReference);
            return "DataReaderEnumDefault(" + languageTypeName + "." + resolverMethod + ", " + getDataReaderCall(enumBaseTypeReference) + ")";
        } else if (typeReference.isArrayTypeReference()) {
            final ArrayTypeReference arrayTypeReference = typeReference.asArrayTypeReference().orElseThrow();
            return getDataReaderCall(arrayTypeReference.getElementTypeReference(), resolverMethod);
        } else if (typeReference.isSimpleTypeReference()) {
            SimpleTypeReference simpleTypeReference = typeReference.asSimpleTypeReference().orElseThrow(IllegalStateException::new);
            return getDataReaderCall(simpleTypeReference);
        } else if (typeReference.isComplexTypeReference()) {
            StringBuilder paramsString = new StringBuilder();
            ComplexTypeReference complexTypeReference = typeReference.asComplexTypeReference().orElseThrow(IllegalStateException::new);
            ComplexTypeDefinition typeDefinition = complexTypeReference.getTypeDefinition();
            String parserCallString = getLanguageTypeNameForTypeReference(typeReference);
            // In case of DataIo we actually need to use the type name and not what above returns.
            // (In this case the mspec type name and the result type name differ)
            if (typeReference.isDataIoTypeReference()) {
                parserCallString = typeReference.asDataIoTypeReference().orElseThrow().getName();
            }
            if (typeDefinition.isDiscriminatedChildTypeDefinition()) {
                parserCallString = "(" + getLanguageTypeNameForTypeReference(typeReference) + ") " + typeDefinition.getParentType().orElseThrow().getName();
            }
            List<Term> paramTerms = complexTypeReference.getParams().orElse(Collections.emptyList());
            for (int i = 0; i < paramTerms.size(); i++) {
                Term paramTerm = paramTerms.get(i);
                final TypeReference argumentType = getArgumentType(complexTypeReference, i);
                paramsString
                    .append(", ")
                    .append(getLanguageTypeNameForTypeReference(argumentType, true))
                    .append("(")
                    .append(toParseExpression(null, argumentType, paramTerm, null))
                    .append(")");
            }
            return "DataReaderComplexDefault(" + parserCallString + ".staticParse(readBuffer" + paramsString + "), readBuffer)";
        } else {
            throw new IllegalStateException("What is this type? " + typeReference);
        }
    }

    public String getDataReaderCall(SimpleTypeReference simpleTypeReference) {
        final int sizeInBits = simpleTypeReference.getSizeInBits();
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return "readBoolean(readBuffer)";
            case BYTE:
                return "readByte(readBuffer, " + sizeInBits + ")";
            case UINT:
                if (sizeInBits <= 4) return "readUnsignedByte(readBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 8) return "readUnsignedShort(readBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 16) return "readUnsignedInt(readBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 32) return "readUnsignedLong(readBuffer, " + sizeInBits + ")";
                return "readUnsignedBigInteger(readBuffer, " + sizeInBits + ")";
            case INT:
                if (sizeInBits <= 8) return "readSignedByte(readBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 16) return "readSignedShort(readBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 32) return "readSignedInt(readBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 64) return "readSignedLong(readBuffer, " + sizeInBits + ")";
                return "readSignedBigInteger(readBuffer, " + sizeInBits + ")";
            case FLOAT:
                if (sizeInBits <= 32) return "readFloat(readBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 64) return "readDouble(readBuffer, " + sizeInBits + ")";
                return "readBigDecimal(readBuffer, " + sizeInBits + ")";
            case STRING:
                return "readString(readBuffer, " + sizeInBits + ")";
            case VSTRING:
                VstringTypeReference vstringTypeReference = (VstringTypeReference) simpleTypeReference;
                return "readString(readBuffer, " + toParseExpression(null, INT_TYPE_REFERENCE, vstringTypeReference.getLengthExpression(), null) + ")";
            case TIME:
                return "readTime(readBuffer)";
            case DATE:
                return "readDate(readBuffer)";
            case DATETIME:
                return "readDateTime(readBuffer)";
            default:
                throw new UnsupportedOperationException("Unsupported type " + simpleTypeReference.getBaseType());
        }
    }

    public String getDataWriterCall(TypeReference typeReference, String fieldName) {
        if (typeReference.isSimpleTypeReference()) {
            SimpleTypeReference simpleTypeReference = typeReference.asSimpleTypeReference().orElseThrow(IllegalStateException::new);
            return getDataWriterCall(simpleTypeReference);
        } else if (typeReference.isArrayTypeReference()) {
            final ArrayTypeReference arrayTypeReference = typeReference.asArrayTypeReference().orElseThrow();
            return getDataWriterCall(arrayTypeReference.getElementTypeReference(), fieldName);
        } else if (typeReference.isComplexTypeReference()) {
            return "DataWriterComplexDefault(writeBuffer)";
        } else {
            throw new IllegalStateException("What is this type? " + typeReference);
        }
    }

    public String getEnumDataWriterCall(EnumTypeReference typeReference, String fieldName, String attributeName) {
        if (!typeReference.isEnumTypeReference()) {
            throw new IllegalArgumentException("this method should only be called for enum types");
        }
        final String languageTypeName = getLanguageTypeNameForTypeReference(typeReference);
        SimpleTypeReference outputTypeReference;
        if ("value".equals(attributeName)) {
            outputTypeReference = getEnumBaseTypeReference(typeReference);
        } else {
            outputTypeReference = getEnumFieldSimpleTypeReference(typeReference.asNonSimpleTypeReference().orElseThrow(), attributeName);
        }
        return "DataWriterEnumDefault(" + languageTypeName + ".get" + StringUtils.capitalize(attributeName) + ", " + languageTypeName + ".name, " + getDataWriterCall(outputTypeReference, fieldName) + ")";
    }

    public String getDataWriterCall(SimpleTypeReference simpleTypeReference) {
        final int sizeInBits = simpleTypeReference.getSizeInBits();
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return "writeBoolean(writeBuffer)";
            case BYTE:
                return "writeByte(writeBuffer, " + sizeInBits + ")";
            case UINT:
                if (sizeInBits <= 4) return "writeUnsignedByte(writeBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 8) return "writeUnsignedShort(writeBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 16) return "writeUnsignedInt(writeBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 32) return "writeUnsignedLong(writeBuffer, " + sizeInBits + ")";
                return "writeUnsignedBigInteger(writeBuffer, " + sizeInBits + ")";
            case INT:
                if (sizeInBits <= 8) return "writeSignedByte(writeBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 16) return "writeSignedShort(writeBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 32) return "writeSignedInt(writeBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 64) return "writeSignedLong(writeBuffer, " + sizeInBits + ")";
                return "writeSignedBigInteger(writeBuffer, " + sizeInBits + ")";
            case FLOAT:
                if (sizeInBits <= 32) return "writeFloat(writeBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 64) return "writeDouble(writeBuffer, " + sizeInBits + ")";
                return "writeBigDecimal(writeBuffer, " + sizeInBits + ")";
            case STRING:
                return "writeString(writeBuffer, " + sizeInBits + ")";
            case VSTRING:
                VstringTypeReference vstringTypeReference = (VstringTypeReference) simpleTypeReference;
                return "writeString(writeBuffer, " + toParseExpression(null, INT_TYPE_REFERENCE, vstringTypeReference.getLengthExpression(), null) + ")";
            case TIME:
                return "writeTime(writeBuffer)";
            case DATE:
                return "writeDate(writeBuffer)";
            case DATETIME:
                return "writeDateTime(writeBuffer)";
            default:
                throw new UnsupportedOperationException("Unsupported type " + simpleTypeReference.getBaseType());
        }
    }

    @Deprecated
    @Override
    public String getWriteBufferWriteMethodCall(SimpleTypeReference simpleTypeReference, String fieldName, TypedField field) {
        return "/*TODO: migrate me*/" + getWriteBufferWriteMethodCall("", simpleTypeReference, fieldName, field);
    }

    @Deprecated
    public String getWriteBufferWriteMethodCall(String logicalName, SimpleTypeReference simpleTypeReference, String fieldName, TypedField field, String... writerArgs) {
        String writerArgsString = "";
        if (writerArgs.length > 0) {
            writerArgsString += ", " + StringUtils.join(writerArgs, ", ");
        }
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return "/*TODO: migrate me*/" + "writeBuffer.writeBit(\"" + logicalName + "\", (boolean) " + fieldName + "" + writerArgsString + ")";
            case BYTE:
                ByteTypeReference byteTypeReference = (ByteTypeReference) simpleTypeReference;
                return "/*TODO: migrate me*/" + "writeBuffer.writeByte(\"" + logicalName + "\", ((Number) " + fieldName + ").byteValue()" + writerArgsString + ")";
            case UINT:
                IntegerTypeReference unsignedIntegerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (unsignedIntegerTypeReference.getSizeInBits() <= 4) {
                    return "/*TODO: migrate me*/" + "writeBuffer.writeUnsignedByte(\"" + logicalName + "\", " + unsignedIntegerTypeReference.getSizeInBits() + ", ((Number) " + fieldName + ").byteValue()" + writerArgsString + ")";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 8) {
                    return "/*TODO: migrate me*/" + "writeBuffer.writeUnsignedShort(\"" + logicalName + "\", " + unsignedIntegerTypeReference.getSizeInBits() + ", ((Number) " + fieldName + ").shortValue()" + writerArgsString + ")";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 16) {
                    return "/*TODO: migrate me*/" + "writeBuffer.writeUnsignedInt(\"" + logicalName + "\", " + unsignedIntegerTypeReference.getSizeInBits() + ", ((Number) " + fieldName + ").intValue()" + writerArgsString + ")";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 32) {
                    return "/*TODO: migrate me*/" + "writeBuffer.writeUnsignedLong(\"" + logicalName + "\", " + unsignedIntegerTypeReference.getSizeInBits() + ", ((Number) " + fieldName + ").longValue()" + writerArgsString + ")";
                }
                return "/*TODO: migrate me*/" + "writeBuffer.writeUnsignedBigInteger(\"" + logicalName + "\", " + unsignedIntegerTypeReference.getSizeInBits() + ", (BigInteger) " + fieldName + "" + writerArgsString + ")";
            case INT:
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (integerTypeReference.getSizeInBits() <= 8) {
                    return "/*TODO: migrate me*/" + "writeBuffer.writeSignedByte(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ", ((Number) " + fieldName + ").byteValue()" + writerArgsString + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    return "/*TODO: migrate me*/" + "writeBuffer.writeShort(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ", ((Number) " + fieldName + ").shortValue()" + writerArgsString + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    return "/*TODO: migrate me*/" + "writeBuffer.writeInt(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ", ((Number) " + fieldName + ").intValue()" + writerArgsString + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 64) {
                    return "/*TODO: migrate me*/" + "writeBuffer.writeLong(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ", ((Number) " + fieldName + ").longValue()" + writerArgsString + ")";
                }
                return "/*TODO: migrate me*/" + "writeBuffer.writeBigInteger(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ", BigInteger.valueOf( " + fieldName + ")" + writerArgsString + ")";
            case FLOAT:
            case UFLOAT:
                FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                if (floatTypeReference.getSizeInBits() <= 32) {
                    return "/*TODO: migrate me*/" + "writeBuffer.writeFloat(\"" + logicalName + "\", " + floatTypeReference.getSizeInBits() + "," + fieldName + "" + writerArgsString + ")";
                } else if (floatTypeReference.getSizeInBits() <= 64) {
                    return "/*TODO: migrate me*/" + "writeBuffer.writeDouble(\"" + logicalName + "\", " + floatTypeReference.getSizeInBits() + "," + fieldName + "" + writerArgsString + ")";
                } else {
                    throw new RuntimeException("Unsupported float type");
                }
            case STRING:
            case VSTRING:
                final Term encodingTerm = field.getEncoding().orElse(new DefaultStringLiteral("UTF-8"));
                if (!(encodingTerm instanceof StringLiteral)) {
                    throw new RuntimeException("Encoding must be a quoted string value");
                }
                String encoding = ((StringLiteral) encodingTerm).getValue();
                String length = Integer.toString(simpleTypeReference.getSizeInBits());
                if (simpleTypeReference.getBaseType() == SimpleTypeReference.SimpleBaseType.VSTRING) {
                    VstringTypeReference vstringTypeReference = (VstringTypeReference) simpleTypeReference;
                    length = toSerializationExpression(field, INT_TYPE_REFERENCE, vstringTypeReference.getLengthExpression(), thisType.getParserArguments().orElse(Collections.emptyList()));
                }
                return "/*TODO: migrate me*/" + "writeBuffer.writeString(\"" + logicalName + "\", " + length + ", \"" +
                    encoding + "\", (String) " + fieldName + "" + writerArgsString + ")";
        }
        throw new FreemarkerException("Unmapped basetype" + simpleTypeReference.getBaseType());
    }

    public String getReservedValue(ReservedField reservedField) {
        final String languageTypeName = getLanguageTypeNameForTypeReference(reservedField.getType(), true);
        if ("BigInteger".equals(languageTypeName)) {
            return "BigInteger.valueOf(" + reservedField.getReferenceValue() + ")";
        } else {
            return "(" + languageTypeName + ") " + reservedField.getReferenceValue();
        }
    }

    /**
     * @param field           this generally only is needed in order to access field attributes such as encoding etc.
     * @param resultType      the type the resulting expression should have
     * @param term            the term representing the expression
     * @param parserArguments any parser arguments, which could be referenced in expressions (Needed for getting the type)
     * @return Java code which does the things defined in 'term'
     */
    public String toParseExpression(Field field, TypeReference resultType, Term term, List<Argument> parserArguments) {
        Tracer tracer = Tracer.start("toParseExpression");
        return tracer + toExpression(field, resultType, term, variableLiteral -> tracer.dive("variableExpressionGenerator") + toVariableParseExpression(field, resultType, variableLiteral, parserArguments));
    }

    /**
     * @param field               this generally only is needed in order to access field attributes such as encoding etc.
     * @param resultType          the type the resulting expression should have
     * @param term                the term representing the expression
     * @param serializerArguments any serializer arguments, which could be referenced in expressions (Needed for getting the type)
     * @return Java code which does the things defined in 'term'
     */
    public String toSerializationExpression(Field field, TypeReference resultType, Term term, List<Argument> serializerArguments) {
        Tracer tracer = Tracer.start("toSerializationExpression");
        return tracer + toExpression(field, resultType, term, variableLiteral -> tracer.dive("variableExpressionGenerator") + toVariableSerializationExpression(field, resultType, variableLiteral, serializerArguments));
    }

    private String toExpression(Field field, TypeReference resultType, Term term, Function<VariableLiteral, String> variableExpressionGenerator) {
        Tracer tracer = Tracer.start("toExpression");
        if (term == null) {
            return tracer + "";
        }
        if (term instanceof Literal) {
            return toLiteralTermExpression(field, resultType, (Literal) term, variableExpressionGenerator, tracer);
        } else if (term instanceof UnaryTerm) {
            return toUnaryTermExpression(field, resultType, (UnaryTerm) term, variableExpressionGenerator, tracer);
        } else if (term instanceof BinaryTerm) {
            return toBinaryTermExpression(field, resultType, (BinaryTerm) term, variableExpressionGenerator, tracer);
        } else if (term instanceof TernaryTerm) {
            return toTernaryTermExpression(field, resultType, (TernaryTerm) term, variableExpressionGenerator, tracer);
        } else {
            throw new RuntimeException("Unsupported Term type " + term.getClass().getName() + ". Actual type " + resultType);
        }
    }

    private String toLiteralTermExpression(Field field, TypeReference resultType, Literal literal, Function<VariableLiteral, String> variableExpressionGenerator, Tracer tracer) {
        tracer = tracer.dive("literal term instanceOf");
        if (literal instanceof NullLiteral) {
            tracer = tracer.dive("null literal instanceOf");
            return tracer + "null";
        } else if (literal instanceof BooleanLiteral) {
            tracer = tracer.dive("boolean literal instanceOf");
            String bool = Boolean.toString(((BooleanLiteral) literal).getValue());
            return tracer + bool.substring(0,1).toUpperCase() + bool.substring(1);
        } else if (literal instanceof NumericLiteral) {
            tracer = tracer.dive("numeric literal instanceOf");
            final String numberString = ((NumericLiteral) literal).getNumber().toString();
            if (resultType.isIntegerTypeReference()) {
                final IntegerTypeReference integerTypeReference = resultType.asIntegerTypeReference().orElseThrow(RuntimeException::new);
                if (integerTypeReference.getBaseType() == SimpleTypeReference.SimpleBaseType.UINT && integerTypeReference.getSizeInBits() >= 32) {
                    tracer = tracer.dive("uint >= 32bit");
                    return tracer + numberString + "L";
                } else if (integerTypeReference.getBaseType() == SimpleTypeReference.SimpleBaseType.INT && integerTypeReference.getSizeInBits() > 32) {
                    tracer = tracer.dive("int > 32bit");
                    return tracer + numberString + "L";
                }
            } else if (resultType.isFloatTypeReference()) {
                final FloatTypeReference floatTypeReference = resultType.asFloatTypeReference().orElseThrow(RuntimeException::new);
                if (floatTypeReference.getSizeInBits() <= 32) {
                    tracer = tracer.dive("float < 32bit");
                    return tracer + numberString + "F";
                }
            }
            return tracer + numberString;
        } else if (literal instanceof HexadecimalLiteral) {
            tracer = tracer.dive("hexadecimal literal instanceOf");
            final String hexString = ((HexadecimalLiteral) literal).getHexString();
            if (resultType.isIntegerTypeReference()) {
                final IntegerTypeReference integerTypeReference = resultType.asIntegerTypeReference().orElseThrow(RuntimeException::new);
                if (integerTypeReference.getBaseType() == SimpleTypeReference.SimpleBaseType.UINT && integerTypeReference.getSizeInBits() >= 32) {
                    tracer = tracer.dive("uint >= 32bit");
                    return tracer + hexString + "L";
                } else if (integerTypeReference.getBaseType() == SimpleTypeReference.SimpleBaseType.INT && integerTypeReference.getSizeInBits() > 32) {
                    tracer = tracer.dive("int > 32bit");
                    return tracer + hexString + "L";
                }
            }
            return tracer + hexString;
        } else if (literal instanceof StringLiteral) {
            tracer = tracer.dive("string literal instanceOf");
            return tracer + "\"" + ((StringLiteral) literal).getValue() + "\"";
        } else if (literal instanceof VariableLiteral) {
            tracer = tracer.dive("variable literal instanceOf");
            VariableLiteral variableLiteral = (VariableLiteral) literal;
            if ("curPos".equals(((VariableLiteral) literal).getName())) {
                return "(positionAware.getPos() - startPos)";
            }
            // If this literal references an Enum type, then we have to output it differently.
            if (getTypeDefinitions().get(variableLiteral.getName()) instanceof EnumTypeDefinition) {
                tracer = tracer.dive("enum definition instanceOf");
                VariableLiteral enumDefinitionChild = variableLiteral.getChild()
                    .orElseThrow(() -> new RuntimeException("enum definitions should have childs"));
                return tracer + variableLiteral.getName() + "." + enumDefinitionChild.getName() +
                    enumDefinitionChild.getChild().map(child -> "." + toVariableExpressionRest(field, resultType, child)).orElse("");
            } else {
                return tracer + variableExpressionGenerator.apply(variableLiteral);
            }
        } else {
            throw new RuntimeException("Unsupported Literal type " + literal.getClass().getName());
        }
    }

    private String toUnaryTermExpression(Field field, TypeReference resultType, UnaryTerm unaryTerm, Function<VariableLiteral, String> variableExpressionGenerator, Tracer tracer) {
        tracer = tracer.dive("unary term instanceOf");
        Term a = unaryTerm.getA();
        switch (unaryTerm.getOperation()) {
            case "!":
                tracer = tracer.dive("case !");
                if ((resultType != getAnyTypeReference()) && !resultType.isBooleanTypeReference()) {
                    throw new IllegalArgumentException("'!(...)' expression requires boolean type. Actual type " + resultType);
                }
                return tracer + "!(" + toExpression(field, resultType, a, variableExpressionGenerator) + ")";
            case "-":
                tracer = tracer.dive("case -");
                if ((resultType != getAnyTypeReference()) && !resultType.isIntegerTypeReference() && !resultType.isFloatTypeReference()) {
                    throw new IllegalArgumentException("'-(...)' expression requires integer or floating-point type. Actual type " + resultType);
                }
                return tracer + "-(" + toExpression(field, resultType, a, variableExpressionGenerator) + ")";
            case "()":
                tracer = tracer.dive("case ()");
                return tracer + "(" + toExpression(field, resultType, a, variableExpressionGenerator) + ")";
            default:
                throw new RuntimeException("Unsupported unary operation type " + unaryTerm.getOperation() + ". Actual type " + resultType);
        }
    }

    private String toBinaryTermExpression(Field field, TypeReference resultType, BinaryTerm binaryTerm, Function<VariableLiteral, String> variableExpressionGenerator, Tracer tracer) {
        tracer = tracer.dive("binary term instanceOf");
        Term a = binaryTerm.getA();
        Term b = binaryTerm.getB();
        String operation = binaryTerm.getOperation();
        switch (operation) {
            case "^": {
                tracer = tracer.dive(operation);
                if ((resultType != getAnyTypeReference()) && !resultType.isIntegerTypeReference() && !resultType.isFloatTypeReference()) {
                    throw new IllegalArgumentException("'A^B' expression requires numeric result type. Actual type " + resultType);
                }
                return tracer + "Math.pow((" + toExpression(field, resultType, a, variableExpressionGenerator) + "), (" + toExpression(field, resultType, b, variableExpressionGenerator) + "))";
            }
            case "*":
            case "/":
            case "%":
            case "+":
            case "-": {
                tracer = tracer.dive(operation);
                if ((resultType != getAnyTypeReference()) && !resultType.isIntegerTypeReference() && !resultType.isFloatTypeReference()) {
                    throw new IllegalArgumentException("'A" + operation + "B' expression requires numeric result type. Actual type " + resultType);
                }
                return tracer + "(" + toExpression(field, resultType, a, variableExpressionGenerator) + ") " + operation + " (" + toExpression(field, resultType, b, variableExpressionGenerator) + ")";
            }
            case ">>":
            case "<<": {
                tracer = tracer.dive(operation);
                return tracer + "(" + toExpression(field, resultType, a, variableExpressionGenerator) + ") " + operation + " (" + toExpression(field, INT_TYPE_REFERENCE, b, variableExpressionGenerator) + ")";
            }
            case ">=":
            case "<=":
            case ">":
            case "<":
            case "==":
            case "!=":
                if ((resultType != getAnyTypeReference()) && !resultType.isBooleanTypeReference()) {
                    throw new IllegalArgumentException("'A" + operation + "B' expression requires boolean result type. Actual type " + resultType);
                }
                // TODO: Try to infer the types of the arguments in this case
                return tracer + "(" + toExpression(field, ANY_TYPE_REFERENCE, a, variableExpressionGenerator) + ") " + operation + " (" + toExpression(field, ANY_TYPE_REFERENCE, b, variableExpressionGenerator) + ")";
            case "&&":
            case "||":
                if ((resultType != getAnyTypeReference()) && !resultType.isBooleanTypeReference()) {
                    throw new IllegalArgumentException("'A" + operation + "B' expression requires boolean result type. Actual type " + resultType);
                }
                return tracer + "(" + toExpression(field, resultType, a, variableExpressionGenerator) + ") " + operation + " (" + toExpression(field, resultType, b, variableExpressionGenerator) + ")";
            case "&":
            case "|":
                if ((resultType != getAnyTypeReference()) && !resultType.isIntegerTypeReference() && !resultType.isByteTypeReference()) {
                    throw new IllegalArgumentException("'A" + operation + "B' expression requires byte or integer result type. Actual type " + resultType);
                }
                return tracer + "(" + toExpression(field, resultType, a, variableExpressionGenerator) + ") " + operation + " (" + toExpression(field, resultType, b, variableExpressionGenerator) + ")";
            default:
                throw new IllegalArgumentException("Unsupported ternary operation type " + operation);
        }
    }

    private String toTernaryTermExpression(Field field, TypeReference resultType, TernaryTerm ternaryTerm, Function<VariableLiteral, String> variableExpressionGenerator, Tracer tracer) {
        tracer = tracer.dive("ternary term instanceOf");
        if ("if".equals(ternaryTerm.getOperation())) {
            Term a = ternaryTerm.getA();
            Term b = ternaryTerm.getB();
            Term c = ternaryTerm.getC();
            return tracer +
                "(" +
                "(" + toExpression(field, BOOL_TYPE_REFERENCE, a, variableExpressionGenerator) + ") ? " +
                toExpression(field, resultType, b, variableExpressionGenerator) + " : " +
                toExpression(field, resultType, c, variableExpressionGenerator) + "" +
                ")";
        } else {
            throw new IllegalArgumentException("Unsupported ternary operation type " + ternaryTerm.getOperation() + ". Actual type " + resultType);
        }
    }

    public String toVariableEnumAccessExpression(VariableLiteral variableLiteral) {
        return variableLiteral.getName();
    }

    private String toVariableParseExpression(Field field, TypeReference resultType, VariableLiteral variableLiteral, List<Argument> parserArguments) {
        Tracer tracer = Tracer.start("toVariableParseExpression");
        // CAST expressions are special as we need to add a ".class" to the second parameter in Java.
        if ("CAST".equals(variableLiteral.getName())) {
            return toCastVariableParseExpression(field, resultType, variableLiteral, parserArguments, tracer);
        }
        // Special handling for ByteOrder enums (Built in enums)
        else if ("BIG_ENDIAN".equals(variableLiteral.getName())) {
            return "ByteOrder.BIG_ENDIAN";
        } else if ("LITTLE_ENDIAN".equals(variableLiteral.getName())) {
            return "ByteOrder.LITTLE_ENDIAN";
        }
        // If we're referencing an implicit field, we need to handle that differently.
        else if (isVariableLiteralImplicitField(variableLiteral)) { // If we are accessing implicit fields, we need to rely on a local variable instead.
            return toImplicitVariableParseExpression(field, resultType, variableLiteral, tracer);
        }
        // Call a static function in the drivers StaticHelper
        else if ("STATIC_CALL".equals(variableLiteral.getName())) {
            return toStaticCallParseExpression(field, resultType, variableLiteral, parserArguments, tracer);
        }
        // Call a built-in global static function
        else if (variableLiteral.getName().equals(variableLiteral.getName().toUpperCase())) { // All uppercase names are not fields, but utility methods.
            return toFunctionCallParseExpression(field, resultType, variableLiteral, parserArguments, tracer);
        }
        // The synthetic checksumRawData is a local field and should not be accessed as bean property.
        boolean isParserArg = "readBuffer".equals(variableLiteral.getName());
        boolean isTypeArg = "_type".equals(variableLiteral.getName());
        if (!isParserArg && !isTypeArg && parserArguments != null) {
            for (Argument serializerArgument : parserArguments) {
                if (serializerArgument.getName().equals(variableLiteral.getName())) {
                    isParserArg = true;
                    break;
                }
            }
        }
        if (isParserArg) {
            tracer = tracer.dive("parser arg");
            return tracer + variableLiteral.getName() + variableLiteral.getChild().map(child -> "." + toVariableExpressionRest(field, resultType, child)).orElse("");
        } else if (isTypeArg) {
            tracer = tracer.dive("type arg");
            String part = variableLiteral.getChild().map(VariableLiteral::getName).orElse("");
            switch (part) {
                case "name":
                    return tracer + "\"" + field.getTypeName() + "\"";
                case "length":
                    return tracer + "\"" + ((SimpleTypeReference) field).getSizeInBits() + "\"";
                case "encoding":
                    String encoding = ((StringLiteral) field.getEncoding().orElse(new DefaultStringLiteral("UTF-8"))).getValue();
                    return tracer + "\"" + encoding + "\"";
                default:
                    return tracer + "";
            }
        } else {
            String indexAddon = "";
            if (variableLiteral.getIndex().isPresent()) {
                indexAddon = ".get(" + variableLiteral.getIndex().orElseThrow() + ")";
            }
            return tracer + variableLiteral.getName() + indexAddon + variableLiteral.getChild().map(child -> "." + toVariableExpressionRest(field, resultType, child)).orElse("");
        }
    }

    private String toCastVariableParseExpression(Field field, TypeReference resultType, VariableLiteral variableLiteral, List<Argument> parserArguments, Tracer tracer) {
        tracer = tracer.dive("CAST");
        List<Term> arguments = variableLiteral.getArgs().orElseThrow(() -> new RuntimeException("A Cast expression needs arguments"));
        if (arguments.size() != 2) {
            throw new RuntimeException("A CAST expression expects exactly two arguments.");
        }
        VariableLiteral firstArgument = arguments.get(0).asLiteral()
            .orElseThrow(() -> new RuntimeException("First argument should be a literal"))
            .asVariableLiteral()
            .orElseThrow(() -> new RuntimeException("First argument should be a Variable literal"));
        StringLiteral typeArgument = arguments.get(1).asLiteral().orElseThrow(() -> new RuntimeException("Second argument should be a String literal"))
            .asStringLiteral()
            .orElseThrow(() -> new RuntimeException("Second argument should be a String literal"));
        String sb = "CAST" + "(" +
            toVariableParseExpression(field, ANY_TYPE_REFERENCE, firstArgument, parserArguments) +
            ", " +
            typeArgument.getValue() + ".class)";
        return tracer + sb + variableLiteral.getChild().map(child -> "." + toVariableExpressionRest(field, resultType, child)).orElse("");
    }

    private String toImplicitVariableParseExpression(Field field, TypeReference resultType, VariableLiteral variableLiteral, Tracer tracer) {
        tracer = tracer.dive("implicit");
        return tracer + variableLiteral.getName();
    }

    private String toStaticCallParseExpression(Field field, TypeReference resultType, VariableLiteral variableLiteral, List<Argument> parserArguments, Tracer tracer) {
        tracer = tracer.dive("STATIC_CALL");
        List<Term> arguments = variableLiteral.getArgs().orElseThrow(() -> new RuntimeException("A STATIC_CALL expression needs arguments"));
        if (arguments.size() < 1) {
            throw new RuntimeException("A STATIC_CALL expression expects at least one argument.");
        }
        // TODO: make it as static import with a emitImport so if a static call is present a "utils" package must be present in the import
        StringBuilder sb = new StringBuilder();
        sb.append(packageName()).append(".utils.StaticHelper.");
        // Get the class and method name
        String methodName = arguments.get(0).asLiteral()
            .orElseThrow(() -> new RuntimeException("First argument should be a literal"))
            .asStringLiteral()
            .orElseThrow(() -> new RuntimeException("Expecting the first argument of a 'STATIC_CALL' to be a StringLiteral")).
            getValue();
        sb.append(methodName).append("(");
        for (int i = 1; i < arguments.size(); i++) {
            Term arg = arguments.get(i);
            if (i > 1) {
                sb.append(", ");
            }
            sb.append(toParseExpression(field, ANY_TYPE_REFERENCE, arg, parserArguments));
           /*if (arg instanceof VariableLiteral) {
                VariableLiteral variableLiteralArg = (VariableLiteral) arg;
                // "readBuffer" is the default name of the reader argument which is always available.
                boolean isParserArg = "readBuffer".equals(variableLiteralArg.getName());
                boolean isTypeArg = "_type".equals(variableLiteralArg.getName());
                if (!isParserArg && !isTypeArg && parserArguments != null) {
                    for (Argument parserArgument : parserArguments) {
                        if (parserArgument.getName().equals(variableLiteralArg.getName())) {
                            isParserArg = true;
                            break;
                        }
                    }
                }
                if (isParserArg) {
                    sb.append(variableLiteralArg.getName()).append(variableLiteralArg.getChild().map(child -> "." + toVariableExpressionRest(child)).orElse(""));
                } else if (isTypeArg) {// We have to manually evaluate the type information at code-generation time.
                    String part = variableLiteralArg.getChild().map(VariableLiteral::getName).orElse("");
                    switch (part) {
                        case "name":
                            sb.append("\"").append(field.getTypeName()).append("\"");
                            break;
                        case "length":
                            sb.append("\"").append(((SimpleTypeReference) field).getSizeInBits()).append("\"");
                            break;
                        case "encoding":
                            String encoding = ((StringLiteral) field.getEncoding().orElse(new DefaultStringLiteral("UTF-8"))).getValue();
                            sb.append("\"").append(encoding).append("\"");
                            break;
                    }
                } else {
                    sb.append(toVariableParseExpression(field, variableLiteralArg, null));
                }
            } else if (arg instanceof StringLiteral) {
                sb.append(((StringLiteral) arg).getValue());
            }*/
        }
        sb.append(")");
        if (variableLiteral.getIndex().isPresent()) {
            // TODO: If this is a byte typed field, this needs to be an array accessor instead.
            sb.append(".get(").append(variableLiteral.getIndex().orElseThrow()).append(")");
        }
        return tracer + sb.toString();
    }

    private String toFunctionCallParseExpression(Field field, TypeReference resultType, VariableLiteral variableLiteral, List<Argument> parserArguments, Tracer tracer) {
        tracer = tracer.dive("FunctionCall");
        StringBuilder sb = new StringBuilder(variableLiteral.getName());
        if (variableLiteral.getArgs().isPresent()) {
            sb.append("(");
            boolean firstArg = true;
            for (Term arg : variableLiteral.getArgs().get()) {
                if (!firstArg) {
                    sb.append(", ");
                }
                // TODO: Try to infer the type of the argument ...
                sb.append(toParseExpression(field, ANY_TYPE_REFERENCE, arg, parserArguments));
                firstArg = false;
            }
            sb.append(")");
        }
        if (variableLiteral.getIndex().isPresent()) {
            // TODO: If this is a byte typed field, this needs to be an array accessor instead.
            sb.append(".get(").append(variableLiteral.getIndex().orElseThrow()).append(")");
        }
        return tracer + sb.toString() + variableLiteral.getChild().map(child -> "." + toVariableExpressionRest(field, resultType, child)).orElse("");
    }

    private String toVariableSerializationExpression(Field field, TypeReference resultType, VariableLiteral variableLiteral, List<Argument> serialzerArguments) {
        Tracer tracer = Tracer.start("variable serialization expression");
        if ("STATIC_CALL".equals(variableLiteral.getName())) {
            return toStaticCallSerializationExpression(field, resultType, variableLiteral, serialzerArguments, tracer);
        }
        // All uppercase names are not fields, but utility methods.
        else if (variableLiteral.getName().equals(variableLiteral.getName().toUpperCase())) {
            return toGlobalFunctionCallSerializationExpression(field, resultType, variableLiteral, serialzerArguments, tracer);
        } else if (isVariableLiteralImplicitField(variableLiteral)) { // If we are accessing implicit fields, we need to rely on a local variable instead.
            tracer = tracer.dive("implicit field");
            final ImplicitField referencedImplicitField = getReferencedImplicitField(variableLiteral);
            return tracer + toSerializationExpression(referencedImplicitField, referencedImplicitField.getType(), getReferencedImplicitField(variableLiteral).getSerializeExpression(), serialzerArguments);
        } else if (isVariableLiteralVirtualField(variableLiteral)) {
            tracer = tracer.dive("virtual field");
            return tracer + toVariableExpressionRest(field, resultType, variableLiteral);
        }
        // The synthetic checksumRawData is a local field and should not be accessed as bean property.
        boolean isSerializerArg = "writeBuffer".equals(variableLiteral.getName()) || "checksumRawData".equals(variableLiteral.getName()) || "_value".equals(variableLiteral.getName()) || "element".equals(variableLiteral.getName()) || "size".equals(variableLiteral.getName());
        boolean isTypeArg = "_type".equals(variableLiteral.getName());
        if (!isSerializerArg && !isTypeArg && serialzerArguments != null) {
            for (Argument serializerArgument : serialzerArguments) {
                if (serializerArgument.getName().equals(variableLiteral.getName())) {
                    isSerializerArg = true;
                    break;
                }
            }
        }
        if (isSerializerArg) {
            tracer = tracer.dive("serializer arg");
            return tracer + variableLiteral.getName() + variableLiteral.getChild().map(child -> "." + toVariableExpressionRest(field, resultType, child)).orElse("");
        } else if (isTypeArg) {
            tracer = tracer.dive("type arg");
            String part = variableLiteral.getChild().map(VariableLiteral::getName).orElse("");
            switch (part) {
                case "name":
                    return tracer + "\"" + field.getTypeName() + "\"";
                case "length":
                    return tracer + "\"" + ((SimpleTypeReference) field).getSizeInBits() + "\"";
                case "encoding":
                    String encoding = ((StringLiteral) field.getEncoding().orElse(new DefaultStringLiteral("UTF-8"))).getValue();
                    return tracer + "\"" + encoding + "\"";
                default:
                    return tracer + "";
            }
        } else {
            return tracer + toVariableExpressionRest(field, resultType, variableLiteral);
        }
    }

    private String toGlobalFunctionCallSerializationExpression(Field field, TypeReference resultType, VariableLiteral variableLiteral, List<Argument> serialzerArguments, Tracer tracer) {
        tracer = tracer.dive("GLOBAL_FUNCTION_CALL");
        StringBuilder sb = new StringBuilder(variableLiteral.getName());
        if (variableLiteral.getArgs().isPresent()) {
            sb.append("(");
            boolean firstArg = true;
            for (Term arg : variableLiteral.getArgs().get()) {
                if (!firstArg) {
                    sb.append(", ");
                }
                sb.append(toSerializationExpression(field, ANY_TYPE_REFERENCE, arg, serialzerArguments));
                firstArg = false;
                /*if (arg instanceof VariableLiteral) {
                    VariableLiteral va = (VariableLiteral) arg;
                    boolean isSerializerArg = "readBuffer".equals(va.getName()) || "writeBuffer".equals(va.getName());
                    boolean isTypeArg = "_type".equals(va.getName());
                    if (!isSerializerArg && !isTypeArg && serialzerArguments != null) {
                        for (Argument serializerArgument : serialzerArguments) {
                            if (serializerArgument.getName().equals(va.getName())) {
                                isSerializerArg = true;
                                break;
                            }
                        }
                    }
                    if (isSerializerArg) {
                        sb.append(va.getName()).append(va.getChild().map(child -> "." + toVariableExpressionRest(child)).orElse(""));
                    } else if (isTypeArg) {
                        String part = va.getChild().map(VariableLiteral::getName).orElse("");
                        switch (part) {
                            case "name":
                                sb.append("\"").append(field.getTypeName()).append("\"");
                                break;
                            case "length":
                                sb.append("\"").append(((SimpleTypeReference) field).getSizeInBits()).append("\"");
                                break;
                            case "encoding":
                                String encoding = ((StringLiteral) field.getEncoding().orElse(new DefaultStringLiteral("UTF-8"))).getValue();
                                sb.append("\"").append(encoding).append("\"");
                                break;
                        }
                    } else {
                        sb.append(toVariableSerializationExpression(field, va, serialzerArguments));
                    }
                } else if (arg instanceof StringLiteral) {
                    sb.append(((StringLiteral) arg).getValue());
                }*/
            }
            sb.append(")");
        }
        return tracer + sb.toString();
    }

    private String toStaticCallSerializationExpression(Field field, TypeReference resultType, VariableLiteral variableLiteral, List<Argument> serialzerArguments, Tracer tracer) {
        tracer = tracer.dive("STATIC_CALL");
        StringBuilder sb = new StringBuilder();
        List<Term> arguments = variableLiteral.getArgs().orElseThrow(() -> new RuntimeException("A STATIC_CALL expression needs arguments"));
        if (arguments.size() < 1) {
            throw new RuntimeException("A STATIC_CALL expression expects at least one argument.");
        }
        // TODO: make it as static import with a emitImport so if a static call is present a "utils" package must be present in the import
        sb.append(packageName()).append(".utils.StaticHelper.");
        // Get the class and method name
        String methodName = arguments.get(0).asLiteral()
            .orElseThrow(() -> new RuntimeException("First argument should be a literal"))
            .asStringLiteral()
            .orElseThrow(() -> new RuntimeException("Expecting the first argument of a 'STATIC_CALL' to be a StringLiteral")).
            getValue();
        //methodName = methodName.substring(1, methodName.length() - 1);
        sb.append(methodName).append("(");
        for (int i = 1; i < arguments.size(); i++) {
            Term arg = arguments.get(i);
            if (i > 1) {
                sb.append(", ");
            }
            sb.append(toSerializationExpression(field, ANY_TYPE_REFERENCE, arg, serialzerArguments));
            /*if (arg instanceof VariableLiteral) {
                VariableLiteral va = (VariableLiteral) arg;
                // "readBuffer" and "_value" are always available in every parser.
                boolean isSerializerArg = "readBuffer".equals(va.getName()) || "writeBuffer".equals(va.getName()) || "_value".equals(va.getName()) || "element".equals(va.getName());
                boolean isTypeArg = "_type".equals(va.getName());
                if (!isSerializerArg && !isTypeArg && serialzerArguments != null) {
                    for (Argument serializerArgument : serialzerArguments) {
                        if (serializerArgument.getName().equals(va.getName())) {
                            isSerializerArg = true;
                            break;
                        }
                    }
                }
                if (isSerializerArg) {
                    sb.append(va.getName()).append(va.getChild().map(child -> "." + toVariableExpressionRest(child)).orElse(""));
                } else if (isTypeArg) {
                    String part = va.getChild().map(VariableLiteral::getName).orElse("");
                    switch (part) {
                        case "name":
                            sb.append("\"").append(field.getTypeName()).append("\"");
                            break;
                        case "length":
                            sb.append("\"").append(((SimpleTypeReference) field).getSizeInBits()).append("\"");
                            break;
                        case "encoding":
                            String encoding = ((StringLiteral) field.getEncoding().orElse(new DefaultStringLiteral("UTF-8"))).getValue();
                            sb.append("\"").append(encoding).append("\"");
                            break;
                    }
                } else {
                    sb.append(toVariableSerializationExpression(field, va, serialzerArguments));
                }
            } else if (arg instanceof StringLiteral) {
                sb.append(((StringLiteral) arg).getValue());
            }*/
        }
        sb.append(")");
        return tracer + sb.toString();
    }

    private String toVariableExpressionRest(Field field, TypeReference resultType, VariableLiteral variableLiteral) {
        Tracer tracer = Tracer.start("variable expression rest");
        // length is kind of a keyword in mspec, so we shouldn't be naming variables length. if we ask for the length of a object we can just return length().
        // This way we can get the length of a string when serializing
        String variableLiteralName = variableLiteral.getName();
        if (variableLiteralName.equals("length")) {
            tracer = tracer.dive("length");
            return tracer + variableLiteralName + "()" + ((variableLiteral.getIndex().isPresent() ? ".get(" + variableLiteral.getIndex().orElseThrow() + ")" : "") +
                variableLiteral.getChild().map(child -> "." + toVariableExpressionRest(field, resultType, child)).orElse(""));
        }
        return tracer + "self.get" + WordUtils.capitalize(variableLiteralName) + "()" + ((variableLiteral.getIndex().isPresent() ? ".get(" + variableLiteral.getIndex().orElseThrow() + ")" : "") +
            variableLiteral.getChild().map(child -> "." + toVariableExpressionRest(field, resultType, child)).orElse(""));
    }

    public String getSizeInBits(ComplexTypeDefinition complexTypeDefinition, List<Argument> parserArguments) {
        int sizeInBits = 0;
        StringBuilder sb = new StringBuilder();
        for (Field field : complexTypeDefinition.getFields()) {
            if (field instanceof ArrayField) {
                ArrayField arrayField = (ArrayField) field;
                final SimpleTypeReference type = (SimpleTypeReference) arrayField.getType();
                switch (arrayField.getLoopType()) {
                    case COUNT:
                        sb.append("(").append(toSerializationExpression(null, INT_TYPE_REFERENCE, arrayField.getLoopExpression(), parserArguments)).append(" * ").append(type.getSizeInBits()).append(") + ");
                        break;
                    case LENGTH:
                        sb.append("(").append(toSerializationExpression(null, INT_TYPE_REFERENCE, arrayField.getLoopExpression(), parserArguments)).append(" * 8) + ");
                        break;
                    case TERMINATED:
                        // No terminated.
                        break;
                }
            } else if (field instanceof TypedField) {
                TypedField typedField = (TypedField) field;
                final TypeReference type = typedField.getType();
                if (field instanceof ManualField) {
                    ManualField manualField = (ManualField) field;
                    sb.append("(").append(toSerializationExpression(null, INT_TYPE_REFERENCE, manualField.getLengthExpression(), parserArguments)).append(") + ");
                } else if (type instanceof SimpleTypeReference) {
                    SimpleTypeReference simpleTypeReference = (SimpleTypeReference) type;
                    if (simpleTypeReference instanceof VstringTypeReference) {
                        sb.append(toSerializationExpression(null, INT_TYPE_REFERENCE, ((VstringTypeReference) simpleTypeReference).getLengthExpression(), parserArguments)).append(" + ");
                    } else {
                        sizeInBits += simpleTypeReference.getSizeInBits();
                    }
                }
            }
        }
        return sb.toString() + sizeInBits;
    }

    public String escapeValue(TypeReference typeReference, String valueString) {
        if (valueString == null) {
            return null;
        }
        if (typeReference instanceof SimpleTypeReference) {
            SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
            switch (simpleTypeReference.getBaseType()) {
                case UINT:
                case INT:
                    // If it's a one character string and is numeric, output it as char.
                    if (!NumberUtils.isParsable(valueString) && (valueString.length() == 1)) {
                        return "'" + valueString + "'";
                    }
                    break;
                case STRING:
                case VSTRING:
                    return "\"" + valueString + "\"";
            }
        }
        return valueString;
    }

    public String getFieldOptions(TypedField field, List<Argument> parserArguments) {
        StringBuilder sb = new StringBuilder();
        final Optional<Term> encodingOptional = field.getEncoding();
        if (encodingOptional.isPresent()) {
            final String encoding = toParseExpression(field, field.getType(), encodingOptional.get(), parserArguments);
            sb.append(", WithOption.WithEncoding(").append(encoding).append(")");
        }
        final Optional<Term> byteOrderOptional = field.getByteOrder();
        if (byteOrderOptional.isPresent()) {
            final String byteOrder = toParseExpression(field, field.getType(), byteOrderOptional.get(), parserArguments);
            sb.append(", WithOption.WithByteOrder(").append(byteOrder).append(")");
        }
        return sb.toString();
    }

    public boolean isBigIntegerSource(Term term) {
        boolean isBigInteger = term.asLiteral()
            .flatMap(LiteralConversions::asVariableLiteral)
            .flatMap(VariableLiteral::getChild)
            .map(Term.class::cast)
            .map(this::isBigIntegerSource)
            .orElse(false);
        return isBigInteger || term.asLiteral()
            .flatMap(LiteralConversions::asVariableLiteral)
            .map(VariableLiteral::getTypeReference)
            .flatMap(TypeReferenceConversions::asIntegerTypeReference)
            .map(integerTypeReference -> integerTypeReference.getSizeInBits() >= 64)
            .orElse(false);
    }

    public boolean needsLongMarker(Optional<SimpleTypeReference> baseTypeReference) {
        return baseTypeReference.isPresent() && baseTypeReference.get().isIntegerTypeReference() && baseTypeReference.get().asIntegerTypeReference().orElseThrow().getSizeInBits() >= 32;
    }

    public void emitRequiredImport(String requiredImport) {
        LOGGER.debug("emitting import '\"{}\"'", requiredImport);
        requiredImports.add(requiredImport);
    }

    public void emitRequiredImport(String alias, String requiredImport) {
        LOGGER.debug("emitting import '{} \"{}'\"", alias, requiredImport);
        requiredImports.add(alias + ' ' + '"' + requiredImport + '"');
    }

    public Set<String> getRequiredImports() {
        return requiredImports;
    }

    public void emitDataIoRequiredImport(String requiredImport) {
        LOGGER.debug("emitting io import '\"{}\"'", requiredImport);
        requiredImportsForDataIo.add(requiredImport);
    }

    public void emitDataIoRequiredImport(String alias, String requiredImport) {
        LOGGER.debug("emitting data io import '{} \"{}'\"", alias, requiredImport);
        requiredImportsForDataIo.add(alias + ' ' + '"' + requiredImport + '"');
    }

    public Set<String> getRequiredImportsForDataIo() {
        return requiredImportsForDataIo;
    }

}

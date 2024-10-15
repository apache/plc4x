/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */

package org.apache.plc4x.language.cs;

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

import java.util.*;
import java.util.function.Function;

@SuppressWarnings({"unused", "WeakerAccess"})
public class CsLanguageTemplateHelper extends BaseFreemarkerLanguageTemplateHelper {

    private final Map<String, String> options;

    public CsLanguageTemplateHelper(TypeDefinition thisType, String protocolName, String flavorName, Map<String, TypeDefinition> types,
                                    Map<String, String> externalTypes, Map<String, String> options) {
        super(thisType, protocolName, flavorName, types);
        this.options = options;
    }

    public String fileName(String protocolName, String languageName, String languageFlavorName) {
        return "drivers." + String.join("", protocolName.split("-")) + "." +
            String.join("", languageFlavorName.split("-"));
    }

    public String packageName() {
        return packageName(protocolName, "cs", flavorName);
    }

    public String packageName(String protocolName, String languageName, String languageFlavorName) {
        return Optional.ofNullable((String) options.get("package")).orElseGet(() ->
            "org.apache.plc4x." + String.join("", languageName.split("-")) + "." +
                String.join("", protocolName.split("-")) + "." +
                String.join("", languageFlavorName.split("-")));
    }

    @Override
    public String getLanguageTypeNameForField(Field field) {
        // If the referenced type is a DataIo type, the value is of type PlcValue.
        if (field.isPropertyField()) {
            PropertyField propertyField = field.asPropertyField().orElseThrow(IllegalStateException::new);
            if (propertyField.getType().isComplexTypeReference()) {
                NonSimpleTypeReference nonSimpleTypeReference = propertyField.getType().asNonSimpleTypeReference().orElseThrow(IllegalStateException::new);
                final TypeDefinition typeDefinition = getTypeDefinitions().get(nonSimpleTypeReference.getName());
                if (typeDefinition instanceof DataIoTypeDefinition) {
                    return "PlcValue";
                }
            }
        }
        return getLanguageTypeNameForTypeReference(((TypedField) field).getType());
    }

    @Override
    public String getLanguageTypeNameForTypeReference(TypeReference typeReference) {
        Objects.requireNonNull(typeReference);
        if (typeReference instanceof ArrayTypeReference) {
            final ArrayTypeReference arrayTypeReference = (ArrayTypeReference) typeReference;
            return getLanguageTypeNameForTypeReference(arrayTypeReference.getElementTypeReference()) + "[]";
        }
        // DataIo data-types always have properties of type PlcValue
        if (typeReference.isDataIoTypeReference()) {
            return "PlcValue";
        }
        if (typeReference.isNonSimpleTypeReference()) {
            return typeReference.asNonSimpleTypeReference().orElseThrow().getName();
        }
        SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return "bool";
            case BYTE:
                return "byte";
            case UINT:
                IntegerTypeReference unsignedIntegerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (unsignedIntegerTypeReference.getSizeInBits() <= 8) {
                    return "byte";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 16) {
                    return "ushort";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 32) {
                    return "uint";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 64) {
                    return "ulong";
                }
                throw new FreemarkerException("Unsupported simple type");
            case INT:
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (integerTypeReference.getSizeInBits() <= 8) {
                    return "sbyte";
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    return "short";
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    return "int";
                }
                if (integerTypeReference.getSizeInBits() <= 64) {
                    return "long";
                }
                throw new FreemarkerException("Unsupported simple type");
            case FLOAT:
            case UFLOAT:
                FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                int sizeInBits = floatTypeReference.getSizeInBits();
                if (sizeInBits <= 32) {
                    return "float";
                }
                if (sizeInBits <= 64) {
                    return "double";
                }
                throw new FreemarkerException("Unsupported simple type");
            case STRING:
            case VSTRING:
                return "string";
            case TIME:
                return "time";
            case DATE:
                return "date";
            case DATETIME:
                return "datetime2";

        }
        throw new FreemarkerException("Unsupported simple type");
    }

    public String getPlcValueTypeForTypeReference(TypeReference typeReference) {
        if (!(typeReference instanceof SimpleTypeReference)) {
            return "PlcStruct";
        }
        SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
        int sizeInBits = simpleTypeReference.getSizeInBits();
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return "PlcBOOL";
            case BYTE:
                return "PlcBYTE";
            case UINT:
                if (sizeInBits <= 8) {
                    return "PlcUSINT";
                }
                if (sizeInBits <= 16) {
                    return "PlcUINT";
                }
                if (sizeInBits <= 32) {
                    return "PlcUDINT";
                }
                if (sizeInBits <= 64) {
                    return "PlcULINT";
                }
                throw new FreemarkerException("Unsupported UINT with bit length " + sizeInBits);
            case INT:
                if (sizeInBits <= 8) {
                    return "PlcSINT";
                }
                if (sizeInBits <= 16) {
                    return "PlcINT";
                }
                if (sizeInBits <= 32) {
                    return "PlcDINT";
                }
                if (sizeInBits <= 64) {
                    return "PlcLINT";
                }
                throw new FreemarkerException("Unsupported INT with bit length " + sizeInBits);
            case FLOAT:
            case UFLOAT:
                if (sizeInBits <= 32) {
                    return "PlcREAL";
                }
                if (sizeInBits <= 64) {
                    return "PlcLREAL";
                }
                throw new FreemarkerException("Unsupported REAL with bit length " + sizeInBits);
            case STRING:
            case VSTRING:
                return "PlcSTRING";
            case TIME:
            case DATE:
            case DATETIME:
                return "PlcTIME";
        }
        throw new FreemarkerException("Unsupported simple type");
    }

    @Override
    public String getNullValueForTypeReference(TypeReference typeReference) {
        if (typeReference instanceof SimpleTypeReference) {
            SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
            switch (simpleTypeReference.getBaseType()) {
                case BIT:
                    return "false";
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
                    return "null";
                case INT:
                    IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                    if (integerTypeReference.getSizeInBits() <= 32) {
                        return "0";
                    }
                    if (integerTypeReference.getSizeInBits() <= 64) {
                        return "0l";
                    }
                    return "null";
                case FLOAT:
                    FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                    int sizeInBits = floatTypeReference.getSizeInBits();
                    if (sizeInBits <= 32) {
                        return "0.0f";
                    }
                    if (sizeInBits <= 64) {
                        return "0.0";
                    }
                    return "null";
                case STRING:
                case VSTRING:
                    return "null";
            }
            throw new FreemarkerException("Unmapped base-type" + simpleTypeReference.getBaseType());
        } else {
            return "null";
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
        return getReadBufferReadMethodCall("", simpleTypeReference, valueString, field);
    }

    @Deprecated
    public String getReadBufferReadMethodCall(String logicalName, SimpleTypeReference simpleTypeReference, String valueString, TypedField field) {
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return "readBuffer.ReadBit(\"" + logicalName + "\")";
            case BYTE:
                return "readBuffer.ReadByte(\"" + logicalName + "\", 8)";
            case UINT:
                String unsignedIntegerType;
                IntegerTypeReference unsignedIntegerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (unsignedIntegerTypeReference.getSizeInBits() <= 8) {
                    unsignedIntegerType = "Byte";
                } else if (unsignedIntegerTypeReference.getSizeInBits() <= 16) {
                    unsignedIntegerType = "Ushort";
                } else if (unsignedIntegerTypeReference.getSizeInBits() <= 32) {
                    unsignedIntegerType = "Uint";
                } else if (unsignedIntegerTypeReference.getSizeInBits() <= 64) {
                    unsignedIntegerType = "Ulong";
                } else {
                    throw new FreemarkerException("Unsupported type");
                }
                return "readBuffer.Read" + unsignedIntegerType + "(\"" + logicalName + "\", " + simpleTypeReference.getSizeInBits() + ")";
            case INT:
                String integerType;
                if (simpleTypeReference.getSizeInBits() <= 8) {
                    integerType = "Sbyte";
                } else if (simpleTypeReference.getSizeInBits() <= 16) {
                    integerType = "Short";
                } else if (simpleTypeReference.getSizeInBits() <= 32) {
                    integerType = "Int";
                } else if (simpleTypeReference.getSizeInBits() <= 64) {
                    integerType = "Long";
                } else {
                    throw new FreemarkerException("Unsupported type");
                }
                return "readBuffer.Read" + integerType + "(\"" + logicalName + "\", " + simpleTypeReference.getSizeInBits() + ")";
            case FLOAT:
                String floatType = (simpleTypeReference.getSizeInBits() <= 32) ? "Float" : "Double";
                return "readBuffer.Read" + floatType + "(\"" + logicalName + "\", " + simpleTypeReference.getSizeInBits() + ")";
            case STRING:
            case VSTRING:
                String stringType = "String";
                final Term encodingTerm = field.getEncoding().orElse(new DefaultStringLiteral("UTF-8"));
                if (!(encodingTerm instanceof StringLiteral)) {
                    throw new FreemarkerException("Encoding must be a quoted string value");
                }
                String encoding = ((StringLiteral) encodingTerm).getValue();
                String length = Integer.toString(simpleTypeReference.getSizeInBits());
                if (simpleTypeReference.getBaseType() == SimpleTypeReference.SimpleBaseType.VSTRING) {
                    VstringTypeReference vstringTypeReference = (VstringTypeReference) simpleTypeReference;
                    length = toParseExpression(field, INT_TYPE_REFERENCE, vstringTypeReference.getLengthExpression(), null);
                }
                return "readBuffer.Read" + stringType + "(\"" + logicalName + "\", " + length + ", System.Text.Encoding.GetEncoding(\"" +
                    encoding + "\"))";
        }
        return "";
    }

    public String getDataReaderCall(TypeReference typeReference) {
        return getDataReaderCall(typeReference, "enumForValue");
    }

    public String getDataReaderCall(TypeReference typeReference, String resolverMethod) {
        if (typeReference.isEnumTypeReference()) {
            final String languageTypeName = getLanguageTypeNameForTypeReference(typeReference);
            final SimpleTypeReference enumBaseTypeReference = getEnumBaseTypeReference(typeReference);
            return "new DataReaderEnumDefault<>(" + languageTypeName + "::" + resolverMethod + ", " + getDataReaderCall(enumBaseTypeReference) + ")";
        } else if (typeReference.isArrayTypeReference()) {
            final ArrayTypeReference arrayTypeReference = typeReference.asArrayTypeReference().orElseThrow();
            return getDataReaderCall(arrayTypeReference.getElementTypeReference(), resolverMethod);
        } else if (typeReference.isSimpleTypeReference()) {
            SimpleTypeReference simpleTypeReference = typeReference.asSimpleTypeReference().orElseThrow(IllegalStateException::new);
            return getDataReaderCall(simpleTypeReference);
        } else if (typeReference.isNonSimpleTypeReference()) {
            StringBuilder paramsString = new StringBuilder();
            NonSimpleTypeReference nonSimpleTypeReference = typeReference.asNonSimpleTypeReference().orElseThrow(IllegalStateException::new);
            ComplexTypeDefinition typeDefinition = typeReference.asNonSimpleTypeReference().orElseThrow().getTypeDefinition().asComplexTypeDefinition().orElseThrow();
            String parserCallString = getLanguageTypeNameForTypeReference(typeReference);
            if (typeDefinition.isDiscriminatedChildTypeDefinition()) {
                parserCallString = "(" + getLanguageTypeNameForTypeReference(typeReference) + ") " + typeDefinition.getParentType().orElseThrow().getName();
            }
            List<Term> paramTerms = nonSimpleTypeReference.getParams().orElse(Collections.emptyList());
            for (int i = 0; i < paramTerms.size(); i++) {
                Term paramTerm = paramTerms.get(i);
                final TypeReference argumentType = getArgumentType(nonSimpleTypeReference, i);
                paramsString
                    .append(", (")
                    .append(getLanguageTypeNameForTypeReference(argumentType))
                    .append(") (")
                    .append(toParseExpression(null, argumentType, paramTerm, null))
                    .append(")");
            }
            return "readComplex(() -> " + parserCallString + "IO.staticParse(readBuffer" + paramsString + "), readBuffer)";
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
        } else if (typeReference.isComplexTypeReference()) {
            return "new DataWriterComplexDefault<>(writeBuffer)";
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
            outputTypeReference = getEnumFieldSimpleTypeReference(typeReference, attributeName);
        }
        return "new DataWriterEnumDefault<>(" + languageTypeName + "::get" + StringUtils.capitalize(attributeName) + ", " + languageTypeName + "::name, " + getDataWriterCall(outputTypeReference, fieldName) + ")";
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
                return "writeDateTime(readBuffer)";
            default:
                throw new UnsupportedOperationException("Unsupported type " + simpleTypeReference.getBaseType());
        }
    }

    @Deprecated
    @Override
    public String getWriteBufferWriteMethodCall(SimpleTypeReference simpleTypeReference, String fieldName, TypedField field) {
        return getWriteBufferWriteMethodCall("", simpleTypeReference, fieldName, field);
    }

    @Deprecated
    public String getWriteBufferWriteMethodCall(String logicalName, SimpleTypeReference simpleTypeReference, String fieldName, TypedField field, String... writerArgs) {
        String writerArgsString = "";
        if (writerArgs.length > 0) {
            writerArgsString += ", " + StringUtils.join(writerArgs, ", ");
        }
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return "writeBuffer.WriteBit(\"" + logicalName + "\", " + fieldName + "" + writerArgsString + ")";
            case BYTE:
                ByteTypeReference byteTypeReference = (ByteTypeReference) simpleTypeReference;
                return "writeBuffer.WriteByte(\"" + logicalName + "\", " + fieldName + writerArgsString + ", 8)";
            case UINT:
                IntegerTypeReference unsignedIntegerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (unsignedIntegerTypeReference.getSizeInBits() <= 8) {
                    return "writeBuffer.WriteByte(\"" + logicalName + "\", " + unsignedIntegerTypeReference.getSizeInBits() + ", (byte) " + fieldName + writerArgsString + ")";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 16) {
                    return "writeBuffer.WriteUshort(\"" + logicalName + "\", " + unsignedIntegerTypeReference.getSizeInBits() + ", (ushort) " + fieldName + writerArgsString + ")";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 32) {
                    return "writeBuffer.WriteUint(\"" + logicalName + "\", " + unsignedIntegerTypeReference.getSizeInBits() + ", (uint) " + fieldName + writerArgsString + ")";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 64) {
                    return "writeBuffer.WriteUlong(\"" + logicalName + "\", " + unsignedIntegerTypeReference.getSizeInBits() + ", (ulong) " + fieldName + writerArgsString + ")";
                }
                throw new FreemarkerException("Unsupported uint type");
            case INT:
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (integerTypeReference.getSizeInBits() <= 8) {
                    return "writeBuffer.WriteSbyte(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ", (sbyte) " + fieldName + writerArgsString + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    return "writeBuffer.WriteShort(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ", (short) " + fieldName + writerArgsString + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    return "writeBuffer.WriteInt(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ", (int) " + fieldName + writerArgsString + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 64) {
                    return "writeBuffer.WriteLong(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ", (long) " + fieldName + writerArgsString + ")";
                }
                throw new FreemarkerException("Unsupported int type");
            case FLOAT:
            case UFLOAT:
                FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                if (floatTypeReference.getSizeInBits() <= 32) {
                    return "writeBuffer.WriteFloat(\"" + logicalName + "\", " + floatTypeReference.getSizeInBits() + "," + fieldName + writerArgsString + ")";
                } else if (floatTypeReference.getSizeInBits() <= 64) {
                    return "writeBuffer.WriteDouble(\"" + logicalName + "\", " + floatTypeReference.getSizeInBits() + "," + fieldName + writerArgsString + ")";
                } else {
                    throw new FreemarkerException("Unsupported float type");
                }
            case STRING:
            case VSTRING:
                final Term encodingTerm = field.getEncoding().orElse(new DefaultStringLiteral("UTF-8"));
                if (!(encodingTerm instanceof StringLiteral)) {
                    throw new FreemarkerException("Encoding must be a quoted string value");
                }
                String encoding = ((StringLiteral) encodingTerm).getValue();
                String length = Integer.toString(simpleTypeReference.getSizeInBits());
                if (simpleTypeReference.getBaseType() == SimpleTypeReference.SimpleBaseType.VSTRING) {
                    VstringTypeReference vstringTypeReference = (VstringTypeReference) simpleTypeReference;
                    length = toSerializationExpression(field, INT_TYPE_REFERENCE, vstringTypeReference.getLengthExpression(), thisType.getParserArguments().orElse(Collections.emptyList()));
                }
                return "writeBuffer.WriteString(\"" + logicalName + "\", " + length + ", \"" +
                    encoding + "\", (string) " + fieldName + "" + writerArgsString + ")";
        }
        throw new FreemarkerException("Unmapped basetype" + simpleTypeReference.getBaseType());
    }

    public String getReservedValue(ReservedField reservedField) {
        final String languageTypeName = getLanguageTypeNameForTypeReference(reservedField.getType());
        if ("BigInteger".equals(languageTypeName)) {
            return "BigInteger.valueOf(" + reservedField.getReferenceValue() + ")";
        } else {
            return "" + reservedField.getReferenceValue();
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
            throw new FreemarkerException("Unsupported Term type " + term.getClass().getName());
        }
    }

    private String toLiteralTermExpression(Field field, TypeReference resultType, Literal literal, Function<VariableLiteral, String> variableExpressionGenerator, Tracer tracer) {
        tracer = tracer.dive("literal term instanceOf");
        if (literal instanceof NullLiteral) {
            tracer = tracer.dive("null literal instanceOf");
            return tracer + "null";
        } else if (literal instanceof BooleanLiteral) {
            tracer = tracer.dive("boolean literal instanceOf");
            return tracer + Boolean.toString(((BooleanLiteral) literal).getValue());
        } else if (literal instanceof NumericLiteral) {
            tracer = tracer.dive("numeric literal instanceOf");
            final String numberString = ((NumericLiteral) literal).getNumber().toString();
            if (resultType.isIntegerTypeReference()) {
                final IntegerTypeReference integerTypeReference = resultType.asIntegerTypeReference().orElseThrow(FreemarkerException::new);
                if (integerTypeReference.getBaseType() == SimpleTypeReference.SimpleBaseType.UINT && integerTypeReference.getSizeInBits() >= 32) {
                    tracer = tracer.dive("uint >= 32bit");
                    return tracer + numberString + "L";
                } else if (integerTypeReference.getBaseType() == SimpleTypeReference.SimpleBaseType.INT && integerTypeReference.getSizeInBits() > 32) {
                    tracer = tracer.dive("int > 32bit");
                    return tracer + numberString + "L";
                }
            } else if (resultType.isFloatTypeReference()) {
                final FloatTypeReference floatTypeReference = resultType.asFloatTypeReference().orElseThrow(FreemarkerException::new);
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
                final IntegerTypeReference integerTypeReference = resultType.asIntegerTypeReference().orElseThrow(FreemarkerException::new);
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
                return "(readBuffer.getPos() - startPos)";
            }
            // If this literal references an Enum type, then we have to output it differently.
            if (getTypeDefinitions().get(variableLiteral.getName()) instanceof EnumTypeDefinition) {
                tracer = tracer.dive("enum definition instanceOf");
                VariableLiteral enumDefinitionChild = variableLiteral.getChild()
                    .orElseThrow(() -> new FreemarkerException("enum definitions should have childs"));
                return tracer + variableLiteral.getName() + "." + enumDefinitionChild.getName() +
                    enumDefinitionChild.getChild().map(child -> "." + toVariableExpressionRest(field, resultType, child)).orElse("");
            } else {
                return tracer + variableExpressionGenerator.apply(variableLiteral);
            }
        } else {
            throw new FreemarkerException("Unsupported Literal type " + literal.getClass().getName());
        }
    }

    private String toUnaryTermExpression(Field field, TypeReference resultType, UnaryTerm unaryTerm, Function<VariableLiteral, String> variableExpressionGenerator, Tracer tracer) {
        tracer = tracer.dive("unary term instanceOf");
        Term a = unaryTerm.getA();
        switch (unaryTerm.getOperation()) {
            case "!":
                tracer = tracer.dive("case !");
                if ((resultType != getAnyTypeReference()) && !resultType.isBooleanTypeReference()) {
                    throw new IllegalArgumentException("'!(...)' expression requires boolean type");
                }
                return tracer + "!(" + toExpression(field, resultType, a, variableExpressionGenerator) + ")";
            case "-":
                tracer = tracer.dive("case -");
                if ((resultType != getAnyTypeReference()) && !resultType.isIntegerTypeReference() && !resultType.isFloatTypeReference()) {
                    throw new IllegalArgumentException("'-(...)' expression requires integer or floating-point type");
                }
                return tracer + "-(" + toExpression(field, resultType, a, variableExpressionGenerator) + ")";
            case "()":
                tracer = tracer.dive("case ()");
                return tracer + "(" + toExpression(field, resultType, a, variableExpressionGenerator) + ")";
            default:
                throw new FreemarkerException("Unsupported unary operation type " + unaryTerm.getOperation());
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
                    throw new IllegalArgumentException("'A^B' expression requires numeric result type");
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
                    throw new IllegalArgumentException("'A" + operation + "B' expression requires numeric result type");
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
            case "!=": {
                if ((resultType != getAnyTypeReference()) && !resultType.isBooleanTypeReference()) {
                    throw new IllegalArgumentException("'A" + operation + "B' expression requires boolean result type");
                }
                // TODO: Try to infer the types of the arguments in this case
                return tracer + "(" + toExpression(field, ANY_TYPE_REFERENCE, a, variableExpressionGenerator) + ") " + operation + " (" + toExpression(field, ANY_TYPE_REFERENCE, b, variableExpressionGenerator) + ")";
            }
            case "&&":
            case "||": {
                if ((resultType != getAnyTypeReference()) && !resultType.isBooleanTypeReference()) {
                    throw new IllegalArgumentException("'A" + operation + "B' expression requires boolean result type");
                }
                return tracer + "(" + toExpression(field, resultType, a, variableExpressionGenerator) + ") " + operation + " (" + toExpression(field, resultType, b, variableExpressionGenerator) + ")";
            }
            case "&":
            case "|": {
                throw new IllegalArgumentException("Implement this some day ...");
            }
            default: {
                throw new IllegalArgumentException("Unsupported ternary operation type " + operation);
            }
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
            throw new IllegalArgumentException("Unsupported ternary operation type " + ternaryTerm.getOperation());
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
            return tracer + variableLiteral.getName() + variableLiteral.getChild().map(child -> "." + toVariableExpressionRest(field, resultType, child)).orElse("");
        }
    }

    private String toCastVariableParseExpression(Field field, TypeReference resultType, VariableLiteral variableLiteral, List<Argument> parserArguments, Tracer tracer) {
        tracer = tracer.dive("CAST");
        List<Term> arguments = variableLiteral.getArgs().orElseThrow(() -> new FreemarkerException("A Cast expression needs arguments"));
        if (arguments.size() != 2) {
            throw new FreemarkerException("A CAST expression expects exactly two arguments.");
        }
        VariableLiteral firstArgument = arguments.get(0).asLiteral()
            .orElseThrow(() -> new FreemarkerException("First argument should be a literal"))
            .asVariableLiteral()
            .orElseThrow(() -> new FreemarkerException("First argument should be a Variable literal"));
        StringLiteral typeArgument = arguments.get(1).asLiteral().orElseThrow(() -> new FreemarkerException("Second argument should be a String literal"))
            .asStringLiteral()
            .orElseThrow(() -> new FreemarkerException("Second argument should be a String literal"));
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
        List<Term> arguments = variableLiteral.getArgs().orElseThrow(() -> new FreemarkerException("A STATIC_CALL expression needs arguments"));
        if (arguments.size() < 1) {
            throw new FreemarkerException("A STATIC_CALL expression expects at least one argument.");
        }
        // TODO: make it as static import with a emitImport so if a static call is present a "utils" package must be present in the import
        StringBuilder sb = new StringBuilder();
        sb.append(packageName()).append(".utils.StaticHelper.");
        // Get the class and method name
        String methodName = arguments.get(0).asLiteral()
            .orElseThrow(() -> new FreemarkerException("First argument should be a literal"))
            .asStringLiteral()
            .orElseThrow(() -> new FreemarkerException("Expecting the first argument of a 'STATIC_CALL' to be a StringLiteral")).
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
        List<Term> arguments = variableLiteral.getArgs().orElseThrow(() -> new FreemarkerException("A STATIC_CALL expression needs arguments"));
        if (arguments.size() < 1) {
            throw new FreemarkerException("A STATIC_CALL expression expects at least one argument.");
        }
        // TODO: make it as static import with a emitImport so if a static call is present a "utils" package must be present in the import
        sb.append(packageName()).append(".utils.StaticHelper.");
        // Get the class and method name
        String methodName = arguments.get(0).asLiteral()
            .orElseThrow(() -> new FreemarkerException("First argument should be a literal"))
            .asStringLiteral()
            .orElseThrow(() -> new FreemarkerException("Expecting the first argument of a 'STATIC_CALL' to be a StringLiteral")).
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
        return tracer + "get" + WordUtils.capitalize(variableLiteralName) + "()" + ((variableLiteral.getIndex().isPresent() ? ".get(" + variableLiteral.getIndex().orElseThrow() + ")" : "") +
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
        if (typeReference.isSimpleTypeReference()) {
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
        } else if (typeReference.isEnumTypeReference()) {
            return "model." + typeReference.asEnumTypeReference().orElseThrow().getName() + "." + valueString;
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

}

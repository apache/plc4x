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
package org.apache.plc4x.language.go;

import org.apache.commons.lang3.StringUtils;
import org.apache.commons.lang3.math.NumberUtils;
import org.apache.commons.text.CaseUtils;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.definitions.DefaultArgument;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.references.DefaultBooleanTypeReference;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.references.DefaultByteOrderTypeReference;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.references.DefaultFloatTypeReference;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.references.DefaultIntegerTypeReference;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.terms.DefaultStringLiteral;
import org.apache.plc4x.plugins.codegenerator.protocol.freemarker.BaseFreemarkerLanguageTemplateHelper;
import org.apache.plc4x.plugins.codegenerator.protocol.freemarker.FreemarkerException;
import org.apache.plc4x.plugins.codegenerator.protocol.freemarker.Tracer;
import org.apache.plc4x.plugins.codegenerator.types.definitions.*;
import org.apache.plc4x.plugins.codegenerator.types.enums.EnumValue;
import org.apache.plc4x.plugins.codegenerator.types.fields.*;
import org.apache.plc4x.plugins.codegenerator.types.references.*;
import org.apache.plc4x.plugins.codegenerator.types.terms.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;

@SuppressWarnings({"unused", "WeakerAccess"})
public class GoLanguageTemplateHelper extends BaseFreemarkerLanguageTemplateHelper {

    private static final Logger LOGGER = LoggerFactory.getLogger(GoLanguageTemplateHelper.class);

    private final Map<String, String> options;

    // TODO: we could condense it to one import set as these can be emitted per template and are not hardcoded anymore

    public final SortedSet<String> requiredImports = new TreeSet<>();

    public final SortedSet<String> requiredImportsForDataIo = new TreeSet<>();

    public GoLanguageTemplateHelper(TypeDefinition thisType, String protocolName, String flavorName, Map<String, TypeDefinition> types,
                                    Map<String, String> options) {
        super(thisType, protocolName, flavorName, types);
        this.options = options;
    }

    public String fileName(String protocolName, String languageName, String languageFlavorName) {
        return String.join("", protocolName.split("\\-")) + "." +
            String.join("", languageFlavorName.split("\\-"));
    }

    public String getSanitizedPackageName() {
        String sanitizedName = getProtocolName().replaceAll("-", "");
        sanitizedName = sanitizedName.replaceAll("\\.", "/");
        sanitizedName = sanitizedName.toLowerCase();
        return sanitizedName;
    }

    // TODO: check if protocol name can be enforced to only contain valid chars
    public String getSanitizedProtocolName() {
        String sanitizedName = getProtocolName().replaceAll("-", "");
        sanitizedName = CaseUtils.toCamelCase(sanitizedName, false, '.');
        return sanitizedName;
    }

    public String packageName(String languageFlavorName) {
        return String.join("", languageFlavorName.split("\\-"));
    }

    @Override
    public String getLanguageTypeNameForField(Field field) {
        boolean optional = field instanceof OptionalField;
        // If the referenced type is a DataIo type, the value is of type PlcValue.
        if (field instanceof PropertyField) {
            PropertyField propertyField = (PropertyField) field;
            if (propertyField.getType() instanceof ComplexTypeReference) {
                ComplexTypeReference complexTypeReference = (ComplexTypeReference) propertyField.getType();
                final TypeDefinition typeDefinition = complexTypeReference.getTypeDefinition();
                if (typeDefinition instanceof DataIoTypeDefinition) {
                    return "PlcValue";
                }
            }
        }
        TypedField typedField = field.asTypedField().orElseThrow();
        String encoding = null;
        Optional<Term> encodingAttribute = field.getAttribute("encoding");
        if (encodingAttribute.isPresent()) {
            encoding = encodingAttribute.get().toString();
        }
        return getLanguageTypeNameForTypeReference(typedField.getType(), encoding);
    }

    public boolean isComplex(Field field) {
        return field instanceof PropertyField && ((PropertyField) field).getType() instanceof NonSimpleTypeReference;
    }

    @Override
    public String getLanguageTypeNameForTypeReference(TypeReference typeReference) {
        return getLanguageTypeNameForTypeReference(typeReference, null);
    }

    public String getLanguageTypeNameForTypeReference(TypeReference typeReference, String encoding) {
        if (typeReference == null) {
            // TODO: shouldn't this be an error case
            return "";
        }
        if (typeReference.isArrayTypeReference()) {
            final ArrayTypeReference arrayTypeReference = (ArrayTypeReference) typeReference;
            TypeReference elementTypeReference = arrayTypeReference.getElementTypeReference();
            return "[]" + getLanguageTypeNameForTypeReference(elementTypeReference);
        }
        if (typeReference.isNonSimpleTypeReference()) {
            return typeReference.asNonSimpleTypeReference().orElseThrow().getName();
        }
        if (typeReference instanceof ByteOrderTypeReference) {
            return "binary.byteOrder";
        }
        SimpleTypeReference simpleTypeReference = typeReference.asSimpleTypeReference().orElseThrow();
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return "bool";
            case BYTE:
                return "byte";
            case UINT:
                IntegerTypeReference unsignedIntegerTypeReference = simpleTypeReference.asIntegerTypeReference().orElseThrow();
                if (unsignedIntegerTypeReference.getSizeInBits() <= 8) {
                    return "uint8";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 16) {
                    return "uint16";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 32) {
                    return "uint32";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 64) {
                    return "uint64";
                }
                emitRequiredImport("math/big");
                return "*big.Int";
            case INT:
                IntegerTypeReference integerTypeReference = simpleTypeReference.asIntegerTypeReference().orElseThrow();
                if (integerTypeReference.getSizeInBits() <= 8) {
                    return "int8";
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    return "int16";
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    return "int32";
                }
                if (integerTypeReference.getSizeInBits() <= 64) {
                    return "int64";
                }
                emitRequiredImport("math/big");
                return "*big.Int";
            case FLOAT:
            case UFLOAT:
                FloatTypeReference floatTypeReference = simpleTypeReference.asFloatTypeReference().orElseThrow();
                int sizeInBits = floatTypeReference.getSizeInBits();
                if (sizeInBits <= 32) {
                    return "float32";
                }
                if (sizeInBits <= 64) {
                    return "float64";
                }
                emitRequiredImport("math/big");
                return "*big.Float";
            case STRING:
            case VSTRING:
                return "string";
            case TIME:
            case DATE:
            case DATETIME:
                emitRequiredImport("time");
                return "time.Time";
            default:
                throw new FreemarkerException("Unsupported simple type");
        }
    }

    public String getPlcValueTypeForTypeReference(TypeReference typeReference) {
        if (typeReference == null) {
            // TODO: shouldn't this be an error case
            return "";
        }
        if (typeReference.isNonSimpleTypeReference()) {
            return ((NonSimpleTypeReference) typeReference).getName();
        }
        SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return "values.NewPlcBOOL";
            case BYTE:
                return "values.NewPlcBYTE";
            case UINT:
                IntegerTypeReference unsignedIntegerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (unsignedIntegerTypeReference.getSizeInBits() <= 8) {
                    return "values.NewPlcUSINT";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 16) {
                    return "values.NewPlcUINT";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 32) {
                    return "values.NewPlcUDINT";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 64) {
                    return "values.NewPlcULINT";
                }
                return "values.NewPlcBINT";
            case INT:
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (integerTypeReference.getSizeInBits() <= 8) {
                    return "values.NewPlcSINT";
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    return "values.NewPlcINT";
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    return "values.NewPlcDINT";
                }
                if (integerTypeReference.getSizeInBits() <= 64) {
                    return "values.NewPlcLINT";
                }
                return "values.NewPlcBINT";
            case FLOAT:
            case UFLOAT:
                FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                int sizeInBits = floatTypeReference.getSizeInBits();
                if (sizeInBits <= 32) {
                    return "values.NewPlcREAL";
                }
                if (sizeInBits <= 64) {
                    return "values.NewPlcLREAL";
                }
                return "values.NewPlcBREAL";
            case STRING:
            case VSTRING:
                return "values.NewPlcSTRING";
            case TIME:
                return "values.NewPlcTIME";
            case DATE:
                return "values.NewPlcDATE";
            case DATETIME:
                return "values.NewPlcDATE_AND_TIME";
            default:
                throw new FreemarkerException("Unsupported simple type" + simpleTypeReference.getBaseType());
        }
    }

    @Override
    public String getNullValueForTypeReference(TypeReference typeReference) {
        if (typeReference instanceof SimpleTypeReference) {
            SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
            switch (simpleTypeReference.getBaseType()) {
                case BIT:
                    return "false";
                case BYTE:
                case UINT:
                case INT:
                    return "0";
                case FLOAT:
                    return "0.0";
                case STRING:
                case VSTRING:
                    return "\"\"";
            }
        } else if (typeReference.isEnumTypeReference()) {
            return "0";
        }
        return "nil";
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
            case VSTRING:
                StringTypeReference stringTypeReference = (StringTypeReference) simpleTypeReference;
                return stringTypeReference.getSizeInBits();
            default:
                return 0;
        }
    }

    public boolean needsPointerAccess(PropertyField field) {
        boolean isAnTypeOfOptional = "optional".equals(field.getTypeName());
        return isAnTypeOfOptional && needPointerAccess(field.getType());
    }

    public boolean needPointerAccess(TypeReference typeReference) {
        boolean isNotAnComplexTypeReference = !typeReference.isComplexTypeReference();
        boolean arrayTypeIsNotAnComplexTypeReference = !(typeReference.isArrayTypeReference() && typeReference.asArrayTypeReference().orElseThrow().getElementTypeReference().isComplexTypeReference());
        return isNotAnComplexTypeReference && arrayTypeIsNotAnComplexTypeReference;
    }

    @Deprecated
    public String getSpecialReadBufferReadMethodCall(String logicalName, SimpleTypeReference simpleTypeReference, TypedField field) {
        return "/*TODO: migrate me*/" + getReadBufferReadMethodCall(logicalName, simpleTypeReference, null, field);
    }

    @Deprecated
    public String getReadBufferReadMethodCall(String logicalName, SimpleTypeReference simpleTypeReference, TypedField field) {
        return "/*TODO: migrate me*/" + getReadBufferReadMethodCall(logicalName, simpleTypeReference, null, field);
    }

    @Override
    @Deprecated
    public String getReadBufferReadMethodCall(SimpleTypeReference simpleTypeReference, String valueString, TypedField field) {
        return "/*TODO: migrate me*/" + getReadBufferReadMethodCall("", simpleTypeReference, valueString, field);
    }

    @Deprecated
    public String getReadBufferReadMethodCall(String logicalName, SimpleTypeReference simpleTypeReference, String valueString, TypedField field) {
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return "/*TODO: migrate me*/" + "readBuffer.ReadBit(\"" + logicalName + "\")";
            case BYTE:
                return "/*TODO: migrate me*/" + "readBuffer.ReadByte(\"" + logicalName + "\")";
            case UINT:
                IntegerTypeReference unsignedIntegerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (unsignedIntegerTypeReference.getSizeInBits() <= 8) {
                    return "/*TODO: migrate me*/" + "readBuffer.ReadUint8(\"" + logicalName + "\", " + unsignedIntegerTypeReference.getSizeInBits() + ")";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 16) {
                    return "/*TODO: migrate me*/" + "readBuffer.ReadUint16(\"" + logicalName + "\", " + unsignedIntegerTypeReference.getSizeInBits() + ")";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 32) {
                    return "/*TODO: migrate me*/" + "readBuffer.ReadUint32(\"" + logicalName + "\", " + unsignedIntegerTypeReference.getSizeInBits() + ")";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 64) {
                    return "/*TODO: migrate me*/" + "readBuffer.ReadUint64(\"" + logicalName + "\", " + unsignedIntegerTypeReference.getSizeInBits() + ")";
                }
                return "/*TODO: migrate me*/" + "readBuffer.ReadBigInt(\"" + logicalName + "\", " + unsignedIntegerTypeReference.getSizeInBits() + ")";
            case INT:
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (integerTypeReference.getSizeInBits() <= 8) {
                    return "/*TODO: migrate me*/" + "readBuffer.ReadInt8(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    return "/*TODO: migrate me*/" + "readBuffer.ReadInt16(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    return "/*TODO: migrate me*/" + "readBuffer.ReadInt32(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 64) {
                    return "/*TODO: migrate me*/" + "readBuffer.ReadInt64(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ")";
                }
                return "/*TODO: migrate me*/" + "readBuffer.ReadBigInt(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ")";
            case FLOAT:
                FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                if (floatTypeReference.getSizeInBits() <= 32) {
                    return "/*TODO: migrate me*/" + "readBuffer.ReadFloat32(\"" + logicalName + "\", " + floatTypeReference.getSizeInBits() + ")";
                }
                if (floatTypeReference.getSizeInBits() <= 64) {
                    return "/*TODO: migrate me*/" + "readBuffer.ReadFloat64(\"" + logicalName + "\", " + floatTypeReference.getSizeInBits() + ")";
                }
                return "/*TODO: migrate me*/" + "readBuffer.ReadBigFloat(\"" + logicalName + "\", " + floatTypeReference.getSizeInBits() + ")";
            case STRING: {
                String encoding = "UTF-8";
                if (field != null) {
                    final Term encodingTerm = field.getEncoding().orElse(new DefaultStringLiteral(encoding));
                    encoding = encodingTerm.asLiteral()
                        .orElseThrow(() -> new FreemarkerException("Encoding must be a literal"))
                        .asStringLiteral()
                        .orElseThrow(() -> new FreemarkerException("Encoding must be a quoted string value")).getValue();
                }
                String length = Integer.toString(simpleTypeReference.getSizeInBits());
                return "/*TODO: migrate me*/" + "readBuffer.ReadString(\"" + logicalName + "\", uint32(" + length + "), utils.WithEncoding(\"" + encoding + "\"))";
            }
            case VSTRING: {
                String encoding = "UTF-8";
                VstringTypeReference vstringTypeReference = (VstringTypeReference) simpleTypeReference;
                if (field != null) {
                    final Term encodingTerm = field.getEncoding().orElse(new DefaultStringLiteral(encoding));
                    encoding = encodingTerm.asLiteral()
                        .orElseThrow(() -> new FreemarkerException("Encoding must be a literal"))
                        .asStringLiteral()
                        .orElseThrow(() -> new FreemarkerException("Encoding must be a quoted string value")).getValue();
                }
                String lengthExpression = toExpression(field, null, vstringTypeReference.getLengthExpression(), null, null, false, false);
                if (vstringTypeReference.getLengthExpression().isTernaryTerm()) {
                    lengthExpression = "(" + lengthExpression + ").(uint32)";
                } else {
                    lengthExpression = "uint32(" + lengthExpression + ")";
                }
                return "/*TODO: migrate me*/" + "readBuffer.ReadString(\"" + logicalName + "\", " + lengthExpression + ", utils.WithEncoding(\"" + encoding + "\"))";
            }
            case TIME:
            case DATE:
            case DATETIME:
                emitRequiredImport("time");
                return "/*TODO: migrate me*/" + "func() (time.Time, error) {raw, err := readBuffer.ReadUint32(\"" + logicalName + "\", 32);return time.UnixMilli(int64(raw)), err;}()";
            default:
                throw new FreemarkerException("Unsupported base type " + simpleTypeReference.getBaseType());
        }
    }

    public String getDataReaderCall(TypeReference typeReference) {
        emitDataReaderRequiredImports();
        return getDataReaderCall(typeReference, "ByValue");
    }

    public String getDataReaderCall(TypeReference typeReference, String resolverMethod) {
        emitDataReaderRequiredImports();
        if (typeReference.isEnumTypeReference()) {
            final String languageTypeName = getLanguageTypeNameForTypeReference(typeReference);
            final SimpleTypeReference enumBaseTypeReference = getEnumBaseTypeReference(typeReference);
            return "ReadEnum(" + languageTypeName + resolverMethod + ", " + getDataReaderCall(enumBaseTypeReference) + ")";
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
            String typeName = getLanguageTypeNameForTypeReference(typeReference);
            String parserCallString = typeName;
            // In case of DataIo we actually need to use the type name and not what above returns.
            // (In this case the mspec type name and the result type name differ)
            if (typeReference.isDataIoTypeReference()) {
                parserCallString = typeReference.asDataIoTypeReference().orElseThrow().getName();
            }
            if (typeDefinition.isDiscriminatedChildTypeDefinition()) {
                parserCallString = typeDefinition.getParentType().orElseThrow().getName();
            }
            List<Term> paramTerms = complexTypeReference.getParams().orElse(Collections.emptyList());
            for (int i = 0; i < paramTerms.size(); i++) {
                Term paramTerm = paramTerms.get(i);
                final TypeReference argumentType = getArgumentType(complexTypeReference, i);
                paramsString
                    .append(", (")
                    .append(getLanguageTypeNameForTypeReference(argumentType))
                    .append(") (")
                    .append(toParseExpression(null, argumentType, paramTerm, null))
                    .append(")");
            }
            String paramsStringString = paramsString.toString();
            if (StringUtils.isNotBlank(paramsStringString) || typeDefinition.isDiscriminatedChildTypeDefinition()) { // In this case we need to spell the function out
                String genericTypeParam = "";
                if (typeDefinition.isDiscriminatedParentTypeDefinition() || typeDefinition.isDiscriminatedChildTypeDefinition()) {
                    genericTypeParam = "[" + typeName + "]";
                }
                return "ReadComplex[" + typeName + "](" + parserCallString + "ParseWithBufferProducer" + genericTypeParam + "(" + StringUtils.substring(paramsStringString, 2) + "), readBuffer)";
            } else {
                return "ReadComplex[" + typeName + "](" + parserCallString + "ParseWithBuffer, readBuffer)";
            }
        } else {
            throw new IllegalStateException("What is this type? " + typeReference);
        }
    }

    public String getDataReaderCall(SimpleTypeReference simpleTypeReference) {
        emitDataReaderRequiredImports();
        final int sizeInBits = simpleTypeReference.getSizeInBits();
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return "ReadBoolean(readBuffer)";
            case BYTE:
                return "ReadByte(readBuffer, " + sizeInBits + ")";
            case UINT:
                if (sizeInBits <= 8) return "ReadUnsignedByte(readBuffer, uint8(" + sizeInBits + "))";
                if (sizeInBits <= 16) return "ReadUnsignedShort(readBuffer, uint8(" + sizeInBits + "))";
                if (sizeInBits <= 32) return "ReadUnsignedInt(readBuffer, uint8(" + sizeInBits + "))";
                if (sizeInBits <= 64) return "ReadUnsignedLong(readBuffer, uint8(" + sizeInBits + "))";
                return "ReadUnsignedBigInteger(readBuffer, uint8(" + sizeInBits + "))";
            case INT:
                if (sizeInBits <= 8) return "ReadSignedByte(readBuffer, uint8(" + sizeInBits + "))";
                if (sizeInBits <= 16) return "ReadSignedShort(readBuffer, uint8(" + sizeInBits + "))";
                if (sizeInBits <= 32) return "ReadSignedInt(readBuffer, uint8(" + sizeInBits + "))";
                if (sizeInBits <= 64) return "ReadSignedLong(readBuffer, uint8(" + sizeInBits + "))";
                return "ReadSignedBigInteger(readBuffer, uint8(" + sizeInBits + "))";
            case FLOAT:
                if (sizeInBits <= 32) return "ReadFloat(readBuffer, uint8(" + sizeInBits + "))";
                if (sizeInBits <= 64) return "ReadDouble(readBuffer, uint8(" + sizeInBits + "))";
                return "ReadBigDecimal(readBuffer, uint8(" + sizeInBits + "))";
            case STRING:
                return "ReadString(readBuffer, uint32(" + sizeInBits + "))";
            case VSTRING:
                VstringTypeReference vstringTypeReference = (VstringTypeReference) simpleTypeReference;
                return "ReadString(readBuffer, uint32(" + toParseExpression(null, INT_TYPE_REFERENCE, vstringTypeReference.getLengthExpression(), null) + "))";
            case TIME:
                return "ReadTime(readBuffer)";
            case DATE:
                return "ReadDate(readBuffer)";
            case DATETIME:
                return "ReadDateTime(readBuffer)";
            default:
                throw new UnsupportedOperationException("Unsupported type " + simpleTypeReference.getBaseType());
        }
    }

    public String getDataWriterCall(TypeReference typeReference, String fieldName) {
        emitDataReaderRequiredImports();
        if (typeReference.isSimpleTypeReference()) {
            SimpleTypeReference simpleTypeReference = typeReference.asSimpleTypeReference().orElseThrow(IllegalStateException::new);
            return getDataWriterCall(simpleTypeReference);
        } else if (typeReference.isArrayTypeReference()) {
            final ArrayTypeReference arrayTypeReference = typeReference.asArrayTypeReference().orElseThrow();
            return getDataWriterCall(arrayTypeReference.getElementTypeReference(), fieldName);
        } else if (typeReference.isComplexTypeReference()) {
            return "WriteComplex[" + getLanguageTypeNameForTypeReference(typeReference) + "](writeBuffer)";
        } else {
            throw new IllegalStateException("What is this type? " + typeReference);
        }
    }

    public String getEnumDataWriterCall(EnumTypeReference typeReference, String fieldName, String attributeName) {
        emitDataReaderRequiredImports();
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
        return "WriteEnum[" + languageTypeName + "," + getLanguageTypeNameForTypeReference(outputTypeReference) + "](" + languageTypeName + ".Get" + StringUtils.capitalize(attributeName) + ", " + languageTypeName + ".PLC4XEnumName, " + getDataWriterCall(outputTypeReference, fieldName) + ")";
    }

    public String getDataWriterCall(SimpleTypeReference simpleTypeReference) {
        emitDataReaderRequiredImports();
        final int sizeInBits = simpleTypeReference.getSizeInBits();
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return "WriteBoolean(writeBuffer)";
            case BYTE:
                return "WriteByte(writeBuffer, " + sizeInBits + ")";
            case UINT:
                if (sizeInBits <= 8) return "WriteUnsignedByte(writeBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 16) return "WriteUnsignedShort(writeBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 32) return "WriteUnsignedInt(writeBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 64) return "WriteUnsignedLong(writeBuffer, " + sizeInBits + ")";
                return "WriteUnsignedBigInteger(writeBuffer, " + sizeInBits + ")";
            case INT:
                if (sizeInBits <= 8) return "WriteSignedByte(writeBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 16) return "WriteSignedShort(writeBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 32) return "WriteSignedInt(writeBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 64) return "WriteSignedLong(writeBuffer, " + sizeInBits + ")";
                return "WriteSignedBigInteger(writeBuffer, " + sizeInBits + ")";
            case FLOAT:
                if (sizeInBits <= 32) return "WriteFloat(writeBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 64) return "WriteDouble(writeBuffer, " + sizeInBits + ")";
                return "WriteBigDecimal(writeBuffer, " + sizeInBits + ")";
            case STRING:
                return "WriteString(writeBuffer, " + sizeInBits + ")";
            case VSTRING:
                VstringTypeReference vstringTypeReference = (VstringTypeReference) simpleTypeReference;
                return "WriteString(writeBuffer, int32(" + toSerializationExpression(null, INT_TYPE_REFERENCE, vstringTypeReference.getLengthExpression(), null) + "))";
            case TIME:
                return "WriteTime(writeBuffer)";
            case DATE:
                return "WriteDate(writeBuffer)";
            case DATETIME:
                return "WriteDateTime(writeBuffer)";
            default:
                throw new UnsupportedOperationException("Unsupported type " + simpleTypeReference.getBaseType());
        }
    }

    @Override
    @Deprecated
    public String getWriteBufferWriteMethodCall(SimpleTypeReference simpleTypeReference, String fieldName, TypedField field) {
        // Fallback if somewhere the method gets called without a name
        String logicalName = fieldName.replaceAll("[\"()*]", "").replaceFirst("_", "");
        return "/*TODO: migrate me*/" + getWriteBufferWriteMethodCall(logicalName, simpleTypeReference, fieldName, field);
    }

    @Deprecated
    public String getWriteBufferWriteMethodCall(String logicalName, SimpleTypeReference simpleTypeReference, Term valueTerm, TypedField field, String... writerArgs) {
        if (valueTerm instanceof BooleanLiteral) {
            return "/*TODO: migrate me*/" + getWriteBufferWriteMethodCall(logicalName, simpleTypeReference, Boolean.toString(((BooleanLiteral) valueTerm).getValue()), field, writerArgs);
        }
        if (valueTerm instanceof NumericLiteral) {
            return "/*TODO: migrate me*/" + getWriteBufferWriteMethodCall(logicalName, simpleTypeReference, ((NumericLiteral) valueTerm).getNumber().toString(), field, writerArgs);
        }
        if (valueTerm instanceof HexadecimalLiteral) {
            return "/*TODO: migrate me*/" + getWriteBufferWriteMethodCall(logicalName, simpleTypeReference, ((HexadecimalLiteral) valueTerm).getHexString(), field, writerArgs);
        }
        if (valueTerm instanceof StringLiteral) {
            return "/*TODO: migrate me*/" + getWriteBufferWriteMethodCall(logicalName, simpleTypeReference, "\"" + ((StringLiteral) valueTerm).getValue() + "\"", field, writerArgs);
        }
        throw new FreemarkerException("Outputting " + valueTerm.toString() + " not implemented yet. Please continue defining other types in the GoLanguageHelper.getWriteBufferWriteMethodCall.");
    }

    @Deprecated
    public String getWriteBufferWriteMethodCall(String logicalName, SimpleTypeReference simpleTypeReference, String fieldName, TypedField field, String... writerArgs) {
        String writerArgsString = "";
        if (writerArgs.length > 0) {
            writerArgsString += ", " + StringUtils.join(writerArgs, ", ");
        }
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return "/*TODO: migrate me*/" + "writeBuffer.WriteBit(\"" + logicalName + "\", " + fieldName + writerArgsString + ")";
            case BYTE:
                return "/*TODO: migrate me*/" + "writeBuffer.WriteByte(\"" + logicalName + "\", " + fieldName + writerArgsString + ")";
            case UINT:
                IntegerTypeReference unsignedIntegerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (unsignedIntegerTypeReference.getSizeInBits() <= 8) {
                    return "/*TODO: migrate me*/" + "writeBuffer.WriteUint8(\"" + logicalName + "\", " + unsignedIntegerTypeReference.getSizeInBits() + ", uint8(" + fieldName + ")" + writerArgsString + ")";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 16) {
                    return "/*TODO: migrate me*/" + "writeBuffer.WriteUint16(\"" + logicalName + "\", " + unsignedIntegerTypeReference.getSizeInBits() + ", uint16(" + fieldName + ")" + writerArgsString + ")";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 32) {
                    return "/*TODO: migrate me*/" + "writeBuffer.WriteUint32(\"" + logicalName + "\", " + unsignedIntegerTypeReference.getSizeInBits() + ", uint32(" + fieldName + ")" + writerArgsString + ")";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 64) {
                    return "/*TODO: migrate me*/" + "writeBuffer.WriteUint64(\"" + logicalName + "\", " + unsignedIntegerTypeReference.getSizeInBits() + ", uint64(" + fieldName + ")" + writerArgsString + ")";
                }
                return "/*TODO: migrate me*/" + "writeBuffer.WriteBigInt(\"" + logicalName + "\", " + unsignedIntegerTypeReference.getSizeInBits() + ", " + fieldName + writerArgsString + ")";
            case INT:
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (integerTypeReference.getSizeInBits() <= 8) {
                    return "/*TODO: migrate me*/" + "writeBuffer.WriteInt8(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ", int8(" + fieldName + ")" + writerArgsString + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    return "/*TODO: migrate me*/" + "writeBuffer.WriteInt16(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ", int16(" + fieldName + ")" + writerArgsString + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    return "/*TODO: migrate me*/" + "writeBuffer.WriteInt32(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ", int32(" + fieldName + ")" + writerArgsString + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 64) {
                    return "/*TODO: migrate me*/" + "writeBuffer.WriteInt64(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ", int64(" + fieldName + ")" + writerArgsString + ")";
                }
                return "/*TODO: migrate me*/" + "writeBuffer.WriteBigInt(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ", " + fieldName + writerArgsString + ")";
            case FLOAT:
            case UFLOAT:
                FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                if (floatTypeReference.getSizeInBits() <= 32) {
                    return "/*TODO: migrate me*/" + "writeBuffer.WriteFloat32(\"" + logicalName + "\", " + floatTypeReference.getSizeInBits() + ", " + fieldName + writerArgsString + ")";
                }
                if (floatTypeReference.getSizeInBits() <= 64) {
                    return "/*TODO: migrate me*/" + "writeBuffer.WriteFloat64(\"" + logicalName + "\", " + floatTypeReference.getSizeInBits() + ", " + fieldName + writerArgsString + ")";
                }
                return "/*TODO: migrate me*/" + "writeBuffer.WriteBigFloat(\"" + logicalName + "\", " + floatTypeReference.getSizeInBits() + ", " + fieldName + writerArgsString + ")";
            case STRING: {
                StringTypeReference stringTypeReference = (StringTypeReference) simpleTypeReference;
                String encoding = "UTF-8";
                if (field != null) {
                    final Term encodingTerm = field.getEncoding().orElse(new DefaultStringLiteral("UTF-8"));
                    encoding = encodingTerm.asLiteral()
                        .orElseThrow(() -> new FreemarkerException("Encoding must be a literal"))
                        .asStringLiteral()
                        .orElseThrow(() -> new FreemarkerException("Encoding must be a quoted string value")).getValue();
                }
                String length = Integer.toString(simpleTypeReference.getSizeInBits());
                return "/*TODO: migrate me*/" + "writeBuffer.WriteString(\"" + logicalName + "\", uint32(" + length + "), " + fieldName + writerArgsString + ", utils.WithEncoding(\"" + encoding + ")\"))";
            }
            case VSTRING: {
                VstringTypeReference vstringTypeReference = (VstringTypeReference) simpleTypeReference;
                String encoding = "UTF-8";
                if (field != null) {
                    final Term encodingTerm = field.getEncoding().orElse(new DefaultStringLiteral("UTF-8"));
                    encoding = encodingTerm.asLiteral()
                        .orElseThrow(() -> new FreemarkerException("Encoding must be a literal"))
                        .asStringLiteral()
                        .orElseThrow(() -> new FreemarkerException("Encoding must be a quoted string value")).getValue();
                }
                String lengthExpression = toExpression(field, null, vstringTypeReference.getLengthExpression(), null, Collections.singletonList(new DefaultArgument("stringLength", new DefaultIntegerTypeReference(SimpleTypeReference.SimpleBaseType.INT, 32))), true, false);
                if (vstringTypeReference.getLengthExpression().isTernaryTerm()) {
                    lengthExpression = "(" + lengthExpression + ").(uint32)";
                } else {
                    lengthExpression = "uint32(" + lengthExpression + ")";
                }
                String length = Integer.toString(simpleTypeReference.getSizeInBits());
                return "/*TODO: migrate me*/" + "writeBuffer.WriteString(\"" + logicalName + "\", " + lengthExpression + ", " + fieldName + writerArgsString + ", utils.WithEncoding(\"" + encoding + ")\"))";
            }
            case DATE:
            case TIME:
            case DATETIME:
                return "/*TODO: migrate me*/" + "writeBuffer.WriteUint32(\"" + logicalName + "\", uint32(" + fieldName + ")" + writerArgsString + ")";
            default:
                throw new FreemarkerException("Unsupported base type " + simpleTypeReference.getBaseType());
        }
    }

    public String getReservedValue(ReservedField reservedField) {
        final String languageTypeName = getLanguageTypeNameForTypeReference(reservedField.getType());
        switch (languageTypeName) {
            case "*big.Int":
                emitRequiredImport("math/big");
                return "big.NewInt(" + reservedField.getReferenceValue() + ")";
            case "*big.Float":
                emitRequiredImport("math/big");
                return "*big.Float(" + reservedField.getReferenceValue() + ")";
            default:
                return languageTypeName + "(" + reservedField.getReferenceValue() + ")";
        }
    }

    public String toTypeSafeCompare(ReservedField reservedField) {
        final String languageTypeName = getLanguageTypeNameForTypeReference(reservedField.getType());
        switch (languageTypeName) {
            case "*big.Int":
            case "*big.Float":
                emitRequiredImport("math/big");
                return "reserved.Cmp(" + getReservedValue(reservedField) + ") != 0";
            default:
                return "reserved != " + getReservedValue(reservedField);
        }
    }

    public String toParseExpression(Field field, TypeReference resultType, Term term, List<Argument> parserArguments) {
        Tracer tracer = Tracer.start("toParseExpression");
        return tracer + toTypedParseExpression(field, resultType, term, parserArguments);
    }

    public String toParseExpression(Field field, TypeReference resultType, Term term, List<Argument> parserArguments, boolean suppressPointerAccess) {
        Tracer tracer = Tracer.start("toParseExpression");
        return tracer + toTypedParseExpression(field, resultType, term, parserArguments, suppressPointerAccess);
    }

    public String toSerializationExpression(Field field, TypeReference resultType, Term term, List<Argument> serializerArguments) {
        Tracer tracer = Tracer.start("toSerializationExpression");
        return tracer + toTypedSerializationExpression(field, resultType, term, serializerArguments);
    }

    public String toBooleanParseExpression(Field field, Term term, List<Argument> parserArguments) {
        Tracer tracer = Tracer.start("toBooleanParseExpression");
        return tracer + toTypedParseExpression(field, new DefaultBooleanTypeReference(), term, parserArguments);
    }

    public String toBooleanSerializationExpression(Field field, Term term, List<Argument> serializerArguments) {
        Tracer tracer = Tracer.start("toBooleanSerializationExpression");
        return tracer + toTypedSerializationExpression(field, new DefaultBooleanTypeReference(), term, serializerArguments);
    }

    public String toIntegerParseExpression(Field field, int sizeInBits, Term term, List<Argument> parserArguments) {
        Tracer tracer = Tracer.start("toIntegerParseExpression");
        return tracer + toTypedParseExpression(field, new DefaultIntegerTypeReference(SimpleTypeReference.SimpleBaseType.UINT, sizeInBits), term, parserArguments);
    }

    public String toIntegerSerializationExpression(Field field, int sizeInBits, Term term, List<Argument> serializerArguments) {
        Tracer tracer = Tracer.start("toIntegerSerializationExpression");
        return tracer + toTypedSerializationExpression(field, new DefaultIntegerTypeReference(SimpleTypeReference.SimpleBaseType.UINT, sizeInBits), term, serializerArguments);
    }

    public String toTypedParseExpression(Field field, TypeReference fieldType, Term term, List<Argument> parserArguments) {
        Tracer tracer = Tracer.start("toTypedParseExpression");
        return tracer + toExpression(field, fieldType, term, parserArguments, null, false, fieldType != null && fieldType.isComplexTypeReference());
    }

    public String toTypedParseExpression(Field field, TypeReference fieldType, Term term, List<Argument> parserArguments, boolean suppressPointerAccess) {
        Tracer tracer = Tracer.start("toTypedParseExpression");
        return tracer + toExpression(field, fieldType, term, parserArguments, null, false, suppressPointerAccess);
    }

    public String toTypedSerializationExpression(Field field, TypeReference fieldType, Term term, List<Argument> serializerArguments) {
        Tracer tracer = Tracer.start("toTypedSerializationExpression");
        return tracer + toExpression(field, fieldType, term, null, serializerArguments, true, false);
    }

    String getCastExpressionForTypeReference(TypeReference typeReference) {
        Tracer tracer = Tracer.start("castExpression");
        if (typeReference instanceof SimpleTypeReference) {
            return tracer.dive("simpleTypeRef") + getLanguageTypeNameForTypeReference(typeReference);
        } else if (typeReference instanceof ByteOrderTypeReference) {
            return tracer.dive("byteOrderTypeRef") + "binary.ByteOrder";
        } else if (typeReference != null) {
            return tracer.dive("anyTypeRef") + "Cast" + getLanguageTypeNameForTypeReference(typeReference);
        } else {
            return tracer.dive("noTypeRef") + "";
        }
    }

    private String toExpression(Field field, TypeReference fieldType, Term term, List<Argument> parserArguments, List<Argument> serializerArguments, boolean serialize, boolean suppressPointerAccess) {
        Tracer tracer = Tracer.start("toExpression(suppressPointerAccess=" + suppressPointerAccess + ")");
        if (term == null) {
            return "";
        }
        if (term instanceof Literal) {
            return toLiteralTermExpression(field, fieldType, term, parserArguments, serializerArguments, serialize, suppressPointerAccess, tracer);
        } else if (term instanceof UnaryTerm) {
            return toUnaryTermExpression(field, fieldType, (UnaryTerm) term, parserArguments, serializerArguments, serialize, tracer);
        } else if (term instanceof BinaryTerm) {
            return toBinaryTermExpression(field, fieldType, (BinaryTerm) term, parserArguments, serializerArguments, serialize, tracer);
        } else if (term instanceof TernaryTerm) {
            return toTernaryTermExpression(field, fieldType, (TernaryTerm) term, parserArguments, serializerArguments, serialize, tracer);
        } else {
            throw new FreemarkerException("Unsupported Term type " + term.getClass().getName());
        }
    }

    private String toTernaryTermExpression(Field field, TypeReference fieldType, TernaryTerm ternaryTerm, List<Argument> parserArguments, List<Argument> serializerArguments, boolean serialize, Tracer tracer) {
        tracer = tracer.dive("ternary term instanceOf");
        if ("if".equals(ternaryTerm.getOperation())) {
            Term a = ternaryTerm.getA();
            Term b = ternaryTerm.getB();
            Term c = ternaryTerm.getC();
            String castExpressionForTypeReference = getCastExpressionForTypeReference(fieldType);
            String inlineIf = "utils.InlineIf(" + toExpression(field, new DefaultBooleanTypeReference(), a, parserArguments, serializerArguments, serialize, false) + ", " +
                "func() any {return " + castExpressionForTypeReference + "(" + toExpression(field, fieldType, b, parserArguments, serializerArguments, serialize, false) + ")}, " +
                "func() any {return " + castExpressionForTypeReference + "(" + toExpression(field, fieldType, c, parserArguments, serializerArguments, serialize, false) + ")})";
            if (fieldType != null) {
                if (fieldType instanceof ByteOrderTypeReference) {
                    return tracer.dive("byteordertypereference") + "(" + inlineIf + ").(binary.ByteOrder)";
                }
                if (fieldType.isNonSimpleTypeReference()) {
                    return tracer.dive("nonsimpletypereference") + castExpressionForTypeReference + "(" + inlineIf + ")";
                }
                return tracer + inlineIf + ".(" + castExpressionForTypeReference + ")";
            }
            return tracer + inlineIf;
        } else {
            throw new FreemarkerException("Unsupported ternary operation type " + ternaryTerm.getOperation());
        }
    }

    private String toBinaryTermExpression(Field field, TypeReference fieldType, BinaryTerm binaryTerm, List<Argument> parserArguments, List<Argument> serializerArguments, boolean serialize, Tracer tracer) {
        tracer = tracer.dive("binary term instanceOf");
        Term a = binaryTerm.getA();
        Term b = binaryTerm.getB();
        String operation = binaryTerm.getOperation();
        String castExpressionForTypeReference = getCastExpressionForTypeReference(fieldType);
        switch (operation) {
            case "^":
                tracer = tracer.dive("^");
                emitRequiredImport("math");
                return tracer + "Math.pow(" +
                    castExpressionForTypeReference + "(" + toExpression(field, fieldType, a, parserArguments, serializerArguments, serialize, false) + "), " +
                    castExpressionForTypeReference + "(" + toExpression(field, fieldType, b, parserArguments, serializerArguments, serialize, false) + "))";
            // If we start casting for comparisons, equals or non equals, really messy things happen.
            case "==":
            case "!=":
            case ">":
            case "<":
            case ">=":
            case "<=":
                tracer = tracer.dive("compare");
                // For every access of optional elements we need pointer access ...
                // Except for doing a nil or not-nil check :-(
                // So in case of such a check, we need to suppress the pointer-access.
                boolean suppressPointerAccessOverride = (operation.equals("==") || operation.equals("!=")) && ((a instanceof NullLiteral) || (b instanceof NullLiteral));
                String aExpression = toExpression(field, null, a, parserArguments, serializerArguments, serialize, suppressPointerAccessOverride);
                String bExpression = toExpression(field, null, b, parserArguments, serializerArguments, serialize, suppressPointerAccessOverride);
                return tracer + "bool((" + aExpression + ") " + operation + " (" + bExpression + "))";
            case ">>":
            case "<<":
            case "|":
            case "&":
                tracer = tracer.dive("bitwise");
                // We don't want casts here
                return tracer +
                    toExpression(field, fieldType, a, parserArguments, serializerArguments, serialize, false) +
                    operation + " " +
                    toExpression(field, fieldType, b, parserArguments, serializerArguments, serialize, false);
            default:
                tracer = tracer.dive("default");
                if (fieldType instanceof StringTypeReference) {
                    tracer = tracer.dive("string type reference");
                    return tracer + toExpression(field, fieldType, a, parserArguments, serializerArguments, serialize, false) +
                        operation + " " +
                        toExpression(field, fieldType, b, parserArguments, serializerArguments, serialize, false);
                }
                return tracer +
                    castExpressionForTypeReference + "(" + toExpression(field, fieldType, a, parserArguments, serializerArguments, serialize, false) + ") " +
                    operation + " " +
                    castExpressionForTypeReference + "(" + toExpression(field, fieldType, b, parserArguments, serializerArguments, serialize, false) + ")";
        }
    }

    private String toUnaryTermExpression(Field field, TypeReference fieldType, UnaryTerm unaryTerm, List<Argument> parserArguments, List<Argument> serializerArguments, boolean serialize, Tracer tracer) {
        tracer = tracer.dive("unary term instanceOf");
        Term a = unaryTerm.getA();
        switch (unaryTerm.getOperation()) {
            case "!":
                tracer = tracer.dive("case !");
                return tracer + "!(" + toExpression(field, fieldType, a, parserArguments, serializerArguments, serialize, false) + ")";
            case "-":
                tracer = tracer.dive("case -");
                return tracer + "-(" + toExpression(field, fieldType, a, parserArguments, serializerArguments, serialize, false) + ")";
            case "()":
                tracer = tracer.dive("case ()");
                return tracer + "(" + toExpression(field, fieldType, a, parserArguments, serializerArguments, serialize, false) + ")";
            default:
                throw new FreemarkerException("Unsupported unary operation type " + unaryTerm.getOperation());
        }
    }

    private String toLiteralTermExpression(Field field, TypeReference fieldType, Term term, List<Argument> parserArguments, List<Argument> serializerArguments, boolean serialize, boolean suppressPointerAccess, Tracer tracer) {
        tracer = tracer.dive("literal term instanceOf");
        if (term instanceof NullLiteral) {
            tracer = tracer.dive("null literal instanceOf");
            return tracer + "nil";
        } else if (term instanceof BooleanLiteral) {
            tracer = tracer.dive("boolean literal instanceOf");
            return tracer + getCastExpressionForTypeReference(fieldType) + "(" + ((BooleanLiteral) term).getValue() + ")";
        } else if (term instanceof NumericLiteral) {
            tracer = tracer.dive("numeric literal instanceOf");
            if (getCastExpressionForTypeReference(fieldType).equals("string")) {
                tracer = tracer.dive("type reference string");
                return tracer + "(" + ((NumericLiteral) term).getNumber().toString() + ")";
            } else {
                return tracer + getCastExpressionForTypeReference(fieldType) + "(" + ((NumericLiteral) term).getNumber().toString() + ")";
            }
        } else if (term instanceof HexadecimalLiteral) {
            tracer = tracer.dive("hexadecimal literal instanceOf");
            return tracer + ((HexadecimalLiteral) term).getHexString();
        } else if (term instanceof StringLiteral) {
            tracer = tracer.dive("string literal instanceOf");
            return tracer + "\"" + ((StringLiteral) term).getValue() + "\"";
        } else if (term instanceof VariableLiteral) {
            tracer = tracer.dive("variable literal instanceOf");
            VariableLiteral variableLiteral = (VariableLiteral) term;
            if ("curPos".equals(((VariableLiteral) term).getName())) {
                return "(positionAware.GetPos() - startPos)";
            } else if ("BIG_ENDIAN".equals(((VariableLiteral) term).getName()) && (fieldType instanceof ByteOrderTypeReference)) {
                return "binary.BigEndian";
            } else if ("LITTLE_ENDIAN".equals(((VariableLiteral) term).getName()) && (fieldType instanceof ByteOrderTypeReference)) {
                return "binary.LittleEndian";
            }
            return tracer + toVariableExpression(field, fieldType, (VariableLiteral) term, parserArguments, serializerArguments, serialize, suppressPointerAccess);
        } else {
            throw new FreemarkerException("Unsupported Literal type " + term.getClass().getName());
        }
    }

    private String toVariableExpression(Field field, TypeReference typeReference, VariableLiteral variableLiteral, List<Argument> parserArguments, List<Argument> serializerArguments, boolean serialize, boolean suppressPointerAccess) {
        return toVariableExpression(field, typeReference, variableLiteral, parserArguments, serializerArguments, serialize, suppressPointerAccess, false);
    }

    private String toVariableExpression(Field field, TypeReference typeReference, VariableLiteral variableLiteral, List<Argument> parserArguments, List<Argument> serializerArguments, boolean serialize, boolean suppressPointerAccess, boolean isChild) {
        Tracer tracer = Tracer.start("toVariableExpression(serialize=" + serialize + ")");
        String variableLiteralName = variableLiteral.getName();
        boolean isEnumTypeReference = typeReference != null && typeReference.isEnumTypeReference();
        if ("lengthInBytes".equals(variableLiteralName)) {
            return toLengthInBytesVariableExpression(typeReference, serialize, tracer);
        } else if ("lengthInBits".equals(variableLiteralName)) {
            return toLengthInBitsVariableExpression(typeReference, serialize, tracer);
        } else if ("_value".equals(variableLiteralName)) {
            return toValueVariableExpression(field, typeReference, variableLiteral, parserArguments, serializerArguments, serialize, suppressPointerAccess, tracer);
        }
        if ("_lastItem".equals(variableLiteralName)) {
            return toLastItemVariableExpression(typeReference, serialize, tracer);
        }
        if ("length".equals(variableLiteral.getChild().map(VariableLiteral::getName).orElse(""))) {
            return toLengthVariableExpression(field, variableLiteral, serialize, tracer);
        }
        // If this literal references an Enum type, then we have to output it differently.
        else if (getTypeDefinitions().get(variableLiteralName) instanceof EnumTypeDefinition) {
            return toEnumVariableExpression(field, typeReference, variableLiteral, parserArguments, serializerArguments, suppressPointerAccess, tracer);
        }
        // If we are accessing enum constants, these also need to be output differently.
        else if (thisType.isComplexTypeDefinition()
            && thisType.asComplexTypeDefinition()
            .orElseThrow(IllegalAccessError::new)
            .getPropertyFieldByName(variableLiteralName)
            .filter(EnumField.class::isInstance)
            .isPresent()
            && (variableLiteral.getChild().isPresent())
        ) {
            return toConstantVariableExpression(field, typeReference, variableLiteral, parserArguments, serializerArguments, suppressPointerAccess, tracer);
        }
        // If we are accessing optional fields, (we might need to use pointer-access).
        else if (!serialize && thisType.isComplexTypeDefinition()
            && thisType.asComplexTypeDefinition()
            .orElseThrow(IllegalStateException::new)
            .getPropertyFieldByName(variableLiteralName)
            .filter(OptionalField.class::isInstance)
            .isPresent()
        ) {
            tracer = tracer.dive("non serialize optional fields");
            return toOptionalVariableExpression(field, typeReference, variableLiteral, parserArguments, serializerArguments, suppressPointerAccess, tracer);
        }
        // If we are accessing optional fields, (we might need to use pointer-access).
        else if (thisType.isComplexTypeDefinition()
            && thisType.asComplexTypeDefinition()
            .orElseThrow(IllegalStateException::new)
            .getPropertyFieldByName(variableLiteralName)
            .filter(OptionalField.class::isInstance)
            .isPresent()
        ) {
            tracer = tracer.dive("optional fields");
            OptionalField optionalField = thisType.asComplexTypeDefinition().orElseThrow().getPropertyFieldByName(variableLiteralName).orElseThrow().asOptionalField().orElseThrow();
            return tracer + "(" + (suppressPointerAccess || optionalField.getType().isComplexTypeReference() ? "" : "*") + "m.Get" + capitalize(variableLiteral.getName()) + "())" +
                variableLiteral.getChild().map(child -> "." + capitalize(toVariableExpression(field, typeReference, child, parserArguments, serializerArguments, false, suppressPointerAccess, true))).orElse("");
        }
        // If we are accessing implicit fields, we need to rely on local variable instead.
        //else if (isVariableLiteralImplicitField(vl)) {
        //    tracer = tracer.dive("implicit");
        //    return tracer + (serialize ? vl.getName() : vl.getName()) + ((vl.getChild() != null) ?
        //        "." + capitalize(toVariableExpression(typeReference, vl.getChild(), parserArguments, serializerArguments, false, suppressPointerAccess)) : "");
        //}
        // If we are accessing implicit fields, we need to rely on a local variable instead.

        // CAST expressions are special as we need to add a ".class" to the second parameter in Java.
        else if ("CAST".equals(variableLiteralName)) {
            return toCastVariableExpression(field, typeReference, variableLiteral, parserArguments, serializerArguments, serialize, suppressPointerAccess, tracer);
        } else if ("STATIC_CALL".equals(variableLiteralName)) {
            return toStaticCallVariableExpression(field, typeReference, variableLiteral, parserArguments, serializerArguments, serialize, suppressPointerAccess, tracer);
        } else if (!isEnumTypeReference && "COUNT".equals(variableLiteralName)) {
            return toCountVariableExpression(field, typeReference, variableLiteral, parserArguments, serializerArguments, serialize, suppressPointerAccess, tracer);
        } else if (!isEnumTypeReference && "ARRAY_SIZE_IN_BYTES".equals(variableLiteralName)) {
            return toArraySizeInBytesVariableExpression(field, typeReference, variableLiteral, parserArguments, serializerArguments, suppressPointerAccess, tracer);
        } else if ("CEIL".equals(variableLiteralName)) {
            return toCeilVariableExpression(field, variableLiteral, parserArguments, serializerArguments, serialize, suppressPointerAccess, tracer);
        } else if ("STR_LEN".equals(variableLiteralName)) {
            return toStrLenVariableExpression(field, typeReference, variableLiteral, parserArguments, serializerArguments, serialize, suppressPointerAccess, tracer);
        }
        // All uppercase names are not fields, but utility methods.
        // TODO: It seems we also run into this, in case of using enum constants in type-switches.
        else if (variableLiteralName.equals(variableLiteralName.toUpperCase())) {
            tracer = tracer.dive("utility");
            return toUppercaseVariableExpression(field, typeReference, variableLiteral, parserArguments, serializerArguments, serialize, suppressPointerAccess, tracer);
        }
        // If the current property references a discriminator value, we have to serialize it differently.
        else if (thisType.isComplexTypeDefinition() && thisType.asComplexTypeDefinition()
            .orElseThrow(IllegalStateException::new)
            .getPropertyFieldFromThisOrParentByName(variableLiteralName)
            .filter(DiscriminatorField.class::isInstance)
            .isPresent()) {
            tracer = tracer.dive("discriminator value");
            // TODO: Should this return something?
        }
        // If the variable has a child element and we're able to find a type for this ... get the type.
        else if ((variableLiteral.getChild().isPresent()) && ((ComplexTypeDefinition) thisType).getTypeReferenceForProperty(variableLiteralName).isPresent()) {
            tracer = tracer.dive("child element");
            final Optional<NonSimpleTypeReference> typeReferenceForProperty = ((ComplexTypeDefinition) thisType).getTypeReferenceForProperty(variableLiteralName)
                .flatMap(TypeReferenceConversions::asNonSimpleTypeReference);
            if (typeReferenceForProperty.isPresent()) {
                tracer = tracer.dive("complex");
                final NonSimpleTypeReference nonSimpleTypeReference = typeReferenceForProperty.get();
                TypeDefinition typeDefinition = nonSimpleTypeReference.getTypeDefinition();
                if (typeDefinition instanceof ComplexTypeDefinition) {
                    tracer = tracer.dive("complex");
                    ComplexTypeDefinition complexTypeDefinition = (ComplexTypeDefinition) typeDefinition;
                    String childProperty = variableLiteral.getChild()
                        .orElseThrow(() -> new FreemarkerException("complex needs a child"))
                        .getName();
                    final Optional<Field> matchingDiscriminatorField = complexTypeDefinition.getFields().stream()
                        .filter(curField -> (curField instanceof DiscriminatorField) && ((DiscriminatorField) curField).getName().equals(childProperty))
                        .findFirst();
                    if (matchingDiscriminatorField.isPresent()) {
                        return tracer + "Cast" + getLanguageTypeNameForTypeReference(nonSimpleTypeReference) + "(" + variableLiteralName + ").Get" + capitalize(childProperty) + "()";
                    }
                    // TODO: is this really meant to fall through?
                    tracer = tracer.dive("we fell through the complex complex");
                } else if (typeDefinition instanceof EnumTypeDefinition) {
                    tracer = tracer.dive("enum");
                    String variableAccess = variableLiteralName;
                    if (isChild) {
                        variableAccess = "Get" + capitalize(variableLiteralName) + "()";
                    }
                    return tracer + (serialize ? "m.Get" + capitalize(variableLiteralName) + "()" : variableAccess) +
                        "." + capitalize(variableLiteral.getChild().orElseThrow(() -> new FreemarkerException("enum needs a child")).getName()) + "()";
                }
            }
            // TODO: is this really meant to fall through?
            tracer = tracer.dive("we fell through the child complete");
        } else if (isVariableLiteralImplicitField(variableLiteral)) {
            tracer = tracer.dive("implicit");
            if (serialize) {
                tracer = tracer.dive("serialize");
                final ImplicitField referencedImplicitField = getReferencedImplicitField(variableLiteral);
                return tracer + toSerializationExpression(referencedImplicitField, referencedImplicitField.getType(), getReferencedImplicitField(variableLiteral).getSerializeExpression(), serializerArguments);
            } else {
                return tracer + variableLiteralName;
                //return toParseExpression(getReferencedImplicitField(vl), getReferencedImplicitField(vl).getSerializeExpression(), serializerArguments);
            }
        }

        // This is a special case for DataIo string types, which need to access the stringLength
        if ((serializerArguments != null) && serializerArguments.stream()
            .anyMatch(argument -> argument.getName().equals(variableLiteralName)) && "stringLength".equals(variableLiteralName)) {
            tracer = tracer.dive("serialization argument");
            return tracer + variableLiteralName +
                variableLiteral.getChild()
                    .map(child -> "." + capitalize(toVariableExpression(field, typeReference, child, parserArguments, serializerArguments, false, suppressPointerAccess, true)))
                    .orElse("");
        } else if ((serializerArguments != null) && serializerArguments.stream()
            .anyMatch(argument -> argument.getName().equals(variableLiteralName))) {
            tracer = tracer.dive("serialization argument");
            return tracer + "m.Get" + capitalize(variableLiteralName) +"()"+
                variableLiteral.getChild()
                    .map(child -> "." + capitalize(toVariableExpression(field, typeReference, child, parserArguments, serializerArguments, false, suppressPointerAccess, true)))
                    .orElse("");
        }
        if ((parserArguments != null) && parserArguments.stream()
            .anyMatch(argument -> argument.getName().equals(variableLiteralName))) {
            tracer = tracer.dive("parser argument");
            return tracer + variableLiteralName +
                variableLiteral.getChild()
                    .map(child -> "." + capitalize(toVariableExpression(field, typeReference, child, parserArguments, serializerArguments, false, suppressPointerAccess, true)))
                    .orElse("");
        }
        String indexCall = "";
        if (variableLiteral.getIndex().isPresent()) {
            tracer = tracer.dive("indexCall");
            // We have a index call
            indexCall = "[" + variableLiteral.getIndex().orElseThrow() + "]";
        }
        tracer = tracer.dive("else");
        Tracer tracer2 = tracer;
        String variableAccess = variableLiteralName;
        if (isChild) {
            variableAccess = "Get" + capitalize(variableAccess) + "()";
        }
        return tracer + (serialize ? "m.Get" + capitalize(variableLiteralName) + "()" : variableAccess) + indexCall +
            variableLiteral.getChild()
                .map(child -> "." + capitalize(toVariableExpression(field, typeReference, child, parserArguments, serializerArguments, false, suppressPointerAccess, true)))
                .orElse("");
    }

    private String toUppercaseVariableExpression(Field field, TypeReference typeReference, VariableLiteral variableLiteral, List<Argument> parserArguments, List<Argument> serializerArguments, boolean serialize, boolean suppressPointerAccess, Tracer tracer) {
        tracer = tracer.dive("toUppercaseVariableExpression");
        StringBuilder sb = new StringBuilder("Get" + capitalize(variableLiteral.getName()) + "()");
        if (variableLiteral.getArgs().isPresent()) {
            sb.append("(");
            boolean firstArg = true;
            for (Term arg : variableLiteral.getArgs().get()) {
                if (!firstArg) {
                    sb.append(", ");
                }
                sb.append(toExpression(field, typeReference, arg, parserArguments, serializerArguments, serialize, suppressPointerAccess));
                firstArg = false;
            }
            sb.append(")");
        }
        if (variableLiteral.getIndex().isPresent()) {
            sb.append("[").append(variableLiteral.getIndex().orElseThrow()).append("]");
        }
        return tracer + sb.toString() + variableLiteral.getChild()
            .map(child -> "." + toVariableExpression(field, typeReference, child, parserArguments, serializerArguments, false, suppressPointerAccess, true))
            .orElse("");
    }

    private String toCeilVariableExpression(Field field, VariableLiteral variableLiteral, List<Argument> parserArguments, List<Argument> serializerArguments, boolean serialize, boolean suppressPointerAccess, Tracer tracer) {
        tracer = tracer.dive("ceil");
        Term va = variableLiteral.getArgs().orElseThrow(() -> new FreemarkerException("CEIL needs at least one arg"))
            .stream().findFirst().orElseThrow(IllegalStateException::new);
        // The Ceil function expects 64 bit floating point values.
        TypeReference tr = new DefaultFloatTypeReference(SimpleTypeReference.SimpleBaseType.FLOAT, 64);
        emitRequiredImport("math");
        return tracer + "math.Ceil(" + toExpression(field, tr, va, parserArguments, serializerArguments, serialize, suppressPointerAccess) + ")";
    }

    private String toArraySizeInBytesVariableExpression(Field field, TypeReference typeReference, VariableLiteral variableLiteral, List<Argument> parserArguments, List<Argument> serializerArguments, boolean suppressPointerAccess, Tracer tracer) {
        tracer = tracer.dive("array size in bytes");
        VariableLiteral va = variableLiteral.getArgs()
            .orElseThrow(() -> new FreemarkerException("ARRAY_SIZE_IN_BYTES needs at least one arg"))
            .stream().findFirst().orElseThrow(IllegalStateException::new)
            .asLiteral()
            .orElseThrow(() -> new FreemarkerException("ARRAY_SIZE_IN_BYTES needs a literal"))
            .asVariableLiteral()
            .orElseThrow(() -> new FreemarkerException("ARRAY_SIZE_IN_BYTES needs a variable literal"));
        // "io" and "m" are always available in every parser.
        boolean isSerializerArg = "readBuffer".equals(va.getName()) || "writeBuffer".equals(va.getName()) || "m".equals(va.getName()) || "element".equals(va.getName());
        if (!isSerializerArg && serializerArguments != null) {
            for (Argument serializerArgument : serializerArguments) {
                if (serializerArgument.getName().equals(va.getName())) {
                    isSerializerArg = true;
                    break;
                }
            }
        }
        StringBuilder sb = new StringBuilder();
        if (isSerializerArg) {
            sb.append(va.getName()).append(va.getChild().map(child -> "." + toVariableExpression(field, typeReference, child, parserArguments, serializerArguments, true, suppressPointerAccess, true)).orElse(""));
        } else {
            sb.append(toVariableExpression(field, typeReference, va, parserArguments, serializerArguments, true, suppressPointerAccess));
        }
        return tracer + getCastExpressionForTypeReference(typeReference) + "(" + va.getName() + "ArraySizeInBytes(" + sb + "))";
    }

    private String toCountVariableExpression(Field field, TypeReference typeReference, VariableLiteral variableLiteral, List<Argument> parserArguments, List<Argument> serializerArguments, boolean serialize, boolean suppressPointerAccess, Tracer tracer) {
        tracer = tracer.dive("count");
        VariableLiteral countLiteral = variableLiteral.getArgs()
            .orElseThrow(() -> new FreemarkerException("Count needs at least one arg"))
            .get(0)
            .asLiteral()
            .orElseThrow(() -> new FreemarkerException("Count needs a literal"))
            .asVariableLiteral()
            .orElseThrow(() -> new FreemarkerException("Count needs a variable literal"));
        return tracer + (typeReference instanceof SimpleTypeReference ? getCastExpressionForTypeReference(typeReference) : "") + "(len(" +
            toVariableExpression(field, typeReference, countLiteral, parserArguments, serializerArguments, serialize, suppressPointerAccess) +
            "))";
    }

    private String toStrLenVariableExpression(Field field, TypeReference typeReference, VariableLiteral variableLiteral, List<Argument> parserArguments, List<Argument> serializerArguments, boolean serialize, boolean suppressPointerAccess, Tracer tracer) {
        tracer = tracer.dive("str-len");
        VariableLiteral countLiteral = variableLiteral.getArgs()
            .orElseThrow(() -> new FreemarkerException("Str-len needs at least one arg"))
            .get(0)
            .asLiteral()
            .orElseThrow(() -> new FreemarkerException("Str-len needs a literal"))
            .asVariableLiteral()
            .orElseThrow(() -> new FreemarkerException("Str-len needs a variable literal"));
        return tracer + (typeReference instanceof SimpleTypeReference ? getCastExpressionForTypeReference(typeReference) : "") + "(len(" +
            toVariableExpression(field, typeReference, countLiteral, parserArguments, serializerArguments, serialize, suppressPointerAccess) +
            "))";
    }

    private String toStaticCallVariableExpression(Field field, TypeReference typeReference, VariableLiteral variableLiteral, List<Argument> parserArguments, List<Argument> serializerArguments, boolean serialize, boolean suppressPointerAccess, Tracer tracer) {
        tracer = tracer.dive("STATIC_CALL");
        StringBuilder sb = new StringBuilder();
        List<Term> arguments = variableLiteral.getArgs().orElseThrow(() -> new FreemarkerException("A STATIC_CALL expression needs arguments"));
        if (arguments.size() < 1) {
            throw new FreemarkerException("A STATIC_CALL expression expects at least one argument.");
        }
        // Get the class and method name
        String staticCall = arguments.get(0).asLiteral()
            .orElseThrow(() -> new FreemarkerException("First argument should be a literal"))
            .asStringLiteral()
            .orElseThrow(() -> new FreemarkerException("Expecting the first argument of a 'STATIC_CALL' to be a StringLiteral")).
            getValue();
        sb.append(capitalize(staticCall)).append("(").append("ctx, ");
        for (int i = 1; i < arguments.size(); i++) {
            Term arg = arguments.get(i);
            if (i > 1) {
                sb.append(", ");
            }
            if (arg instanceof UnaryTerm) {
                arg = ((UnaryTerm) arg).getA();
            }
            if (arg instanceof VariableLiteral) {
                tracer = tracer.dive("VariableLiteral nr." + i);
                VariableLiteral va = (VariableLiteral) arg;
                // "io" is the default name of the reader argument which is always available.
                boolean isParserArg = "readBuffer".equals(va.getName()) || "writeBuffer".equals(va.getName()) || ((thisType instanceof DataIoTypeDefinition) && "_value".equals(va.getName()));
                boolean isTypeArg = "_type".equals(va.getName());
                if (!isParserArg && !isTypeArg && parserArguments != null) {
                    for (Argument parserArgument : parserArguments) {
                        if (parserArgument.getName().equals(va.getName())) {
                            isParserArg = true;
                            break;
                        }
                    }
                }
                if (isParserArg) {
                    tracer = tracer.dive("isParserArg");
                    if (va.getName().equals("_value")) {
                        tracer = tracer.dive("is _value");
                        sb.append(va.getName().substring(1)).append(va.getChild().map(child -> "." + toVariableExpression(field, typeReference, child, parserArguments, serializerArguments, false, suppressPointerAccess, true)).orElse(""));
                    } else {
                        sb.append(va.getName()).append((va.getChild().isPresent()) ?
                            ".Get" + capitalize(toVariableExpression(field, typeReference, va.getChild().orElseThrow(IllegalStateException::new), parserArguments, serializerArguments, false, suppressPointerAccess)) + "()" : "");
                    }
                }
                // We have to manually evaluate the type information at code-generation time.
                else if (isTypeArg) {
                    String part = va.getChild().map(VariableLiteral::getName).orElse("");
                    switch (part) {
                        case "name":
//                                sb.append("\"").append(field.getTypeName()).append("\"");
                            break;
                        case "length":
                            sb.append("\"").append(((SimpleTypeReference) typeReference).getSizeInBits()).append("\"");
                            break;
                        case "encoding":
                            final Term encodingTerm = field.getEncoding().orElse(new DefaultStringLiteral("UTF-8"));
                            if (!(encodingTerm instanceof StringLiteral)) {
                                throw new FreemarkerException("Encoding must be a quoted string value");
                            }
                            String encoding = ((StringLiteral) encodingTerm).getValue();
                            sb.append("\"").append(encoding).append("\"");
                            break;
                    }
                } else {
                    sb.append(toVariableExpression(field, typeReference, va, parserArguments, serializerArguments, serialize, suppressPointerAccess));
                }
            } else if (arg instanceof StringLiteral) {
                tracer = tracer.dive("StringLiteral");
                sb.append(((StringLiteral) arg).getValue());
            } else if (arg instanceof BooleanLiteral) {
                tracer = tracer.dive("BooleanLiteral");
                sb.append(((BooleanLiteral) arg).getValue());
            } else if (arg instanceof NumericLiteral) {
                tracer = tracer.dive("NumericLiteral");
                sb.append(((NumericLiteral) arg).getNumber());
            } else if (arg instanceof BinaryTerm) {
                tracer = tracer.dive("BinaryTerm");
                sb.append(toBinaryTermExpression(field, typeReference, (BinaryTerm) arg, parserArguments, serializerArguments, serialize, tracer));
            } else {
                throw new FreemarkerException(arg.getClass().getName());
            }
        }
        sb.append(")");
        return tracer + sb.toString();
    }

    private String toCastVariableExpression(Field field, TypeReference typeReference, VariableLiteral variableLiteral, List<Argument> parserArguments, List<Argument> serializerArguments, boolean serialize, boolean suppressPointerAccess, Tracer tracer) {
        tracer = tracer.dive("CAST");
        List<Term> arguments = variableLiteral.getArgs().orElseThrow(() -> new FreemarkerException("A Cast expression needs arguments"));
        if (arguments.size() != 2) {
            throw new FreemarkerException("A CAST expression expects exactly two arguments.");
        }
        VariableLiteral firstArgument = arguments.get(0).asLiteral()
            .orElseThrow(() -> new FreemarkerException("First argument should be a literal"))
            .asVariableLiteral()
            .orElseThrow(() -> new FreemarkerException("First argument should be a Variable literal"));
        StringLiteral typeLiteral = arguments.get(1).asLiteral()
            .orElseThrow(() -> new FreemarkerException("Second argument should be a String literal"))
            .asStringLiteral()
            .orElseThrow(() -> new FreemarkerException("Second argument should be a String literal"));
        final TypeDefinition typeDefinition = getTypeDefinitions().get(typeLiteral.getValue());
        StringBuilder sb = new StringBuilder();
        if (typeDefinition.isComplexTypeDefinition()) {
            sb.append("Cast");
        }
        sb.append(typeLiteral.getValue());
        sb.append("(").append(toVariableExpression(field, typeReference, firstArgument, parserArguments, serializerArguments, serialize, suppressPointerAccess)).append(")");
        return tracer + sb.toString() + variableLiteral.getChild().map(child -> "." + capitalize(toVariableExpression(field, typeReference, child, parserArguments, serializerArguments, false, suppressPointerAccess, true))).orElse("");
    }

    private String toOptionalVariableExpression(Field field, TypeReference typeReference, VariableLiteral variableLiteral, List<Argument> parserArguments, List<Argument> serializerArguments, boolean suppressPointerAccess, Tracer tracer) {
        tracer = tracer.dive("optional fields");
        return tracer + "(" + (suppressPointerAccess || (typeReference != null && typeReference.isComplexTypeReference()) ? "" : "*") + variableLiteral.getName() + ")" +
            variableLiteral.getChild().map(child -> "." + capitalize(toVariableExpression(field, typeReference, child, parserArguments, serializerArguments, false, suppressPointerAccess, true))).orElse("");
    }

    private String toConstantVariableExpression(Field field, TypeReference typeReference, VariableLiteral variableLiteral, List<Argument> parserArguments, List<Argument> serializerArguments, boolean suppressPointerAccess, Tracer tracer) {
        tracer = tracer.dive("enum constant");
        VariableLiteral child = variableLiteral.getChild().orElseThrow(IllegalStateException::new);
        return tracer + variableLiteral.getName() + "." + capitalize(child.getName()) + "()" +
            child.getChild().map(childChild -> "." + toVariableExpression(field, typeReference, childChild, parserArguments, serializerArguments, false, suppressPointerAccess, true)).orElse("");
    }

    private String toEnumVariableExpression(Field field, TypeReference typeReference, VariableLiteral variableLiteral, List<Argument> parserArguments, List<Argument> serializerArguments, boolean suppressPointerAccess, Tracer tracer) {
        tracer = tracer.dive("enum");
        VariableLiteral child = variableLiteral.getChild().orElseThrow(() -> new FreemarkerException("Enum should have a child"));
        return tracer + variableLiteral.getName() + "_" + child.getName() +
            child.getChild().map(childChild -> "." + toVariableExpression(field, typeReference, childChild, parserArguments, serializerArguments, false, suppressPointerAccess, true)).orElse("");
    }

    private String toLastItemVariableExpression(TypeReference typeReference, boolean serialize, Tracer tracer) {
        tracer = tracer.dive("lastItem");
        return tracer + "utils.GetLastItemFromContext(ctx)";
    }

    private String toLengthVariableExpression(Field field, VariableLiteral variableLiteral, boolean serialize, Tracer tracer) {
        tracer = tracer.dive("length");
        return tracer + (serialize ? ("len(m.Get" + capitalize(variableLiteral.getName()) + "())") : ("(" + variableLiteral.getName() + ")"));
    }

    private String toValueVariableExpression(Field field, TypeReference typeReference, VariableLiteral variableLiteral, List<Argument> parserArguments, List<Argument> serializerArguments, boolean serialize, boolean suppressPointerAccess, Tracer tracer) {
        final Tracer tracer2 = tracer.dive("_value");
        return variableLiteral.getChild()
            .map(child -> tracer2.dive("withChild") + "m." + toUppercaseVariableExpression(field, typeReference, child, parserArguments, serializerArguments, serialize, suppressPointerAccess, tracer2))
            .orElse(tracer2 + "m");
    }

    private String toLengthInBitsVariableExpression(TypeReference typeReference, boolean serialize, Tracer tracer) {
        tracer = tracer.dive("lengthInBits");
        return tracer + (serialize ? getCastExpressionForTypeReference(typeReference) + "(m.Get" : "Get") + "LengthInBits" + (serialize ? "(ctx))" : "(ctx)");
    }

    private String toLengthInBytesVariableExpression(TypeReference typeReference, boolean serialize, Tracer tracer) {
        tracer = tracer.dive("lengthInBytes");
        return tracer + (serialize ? getCastExpressionForTypeReference(typeReference) + "(m.Get" : "Get") + "LengthInBytes" + (serialize ? "(ctx))" : "(ctx)");
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
                        sb.append("(").append(toTypedSerializationExpression(field, type, arrayField.getLoopExpression(), parserArguments)).append(" * ").append(type.getSizeInBits()).append(") + ");
                        break;
                    case LENGTH:
                        sb.append("(").append(toTypedSerializationExpression(field, type, arrayField.getLoopExpression(), parserArguments)).append(" * 8) + ");
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
                    sb.append("(").append(toSerializationExpression(manualField, getIntTypeReference(), manualField.getLengthExpression(), parserArguments)).append(") + ");
                } else if (type instanceof SimpleTypeReference) {
                    SimpleTypeReference simpleTypeReference = (SimpleTypeReference) type;
                    sizeInBits += simpleTypeReference.getSizeInBits();
                } else {
                    throw new IllegalStateException("No ComplexTypeReference supported");
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

    public String escapeEnumValue(TypeReference typeReference, String valueString) {
        // Currently the only case in which here complex type references are used are when referencing enum constants.
        if (typeReference != null && typeReference.isNonSimpleTypeReference()) {
            // C doesn't like NULL values for enums, so we have to return something else (we'll treat -1 as NULL)
            if ("null".equals(valueString) || valueString == null) {
                return "0";
            }
            if (valueString.contains(".")) {
                String typeName = valueString.substring(0, valueString.indexOf('.'));
                String constantName = valueString.substring(valueString.indexOf('.') + 1);
                return typeName + "_" + constantName;
            }
            return valueString;
        } else {
            return escapeValue(typeReference, valueString);
        }
    }

    public Collection<EnumValue> getUniqueEnumValues(List<EnumValue> enumValues) {
        Map<String, EnumValue> filteredEnumValues = new TreeMap<>();
        for (EnumValue enumValue : enumValues) {
            if (!filteredEnumValues.containsKey(enumValue.getValue())) {
                filteredEnumValues.put(enumValue.getValue(), enumValue);
            }
        }
        return filteredEnumValues.values();
    }

    public List<DiscriminatedComplexTypeDefinition> getUniqueSwitchCases(List<DiscriminatedComplexTypeDefinition> allSwitchCases) {
        Map<String, DiscriminatedComplexTypeDefinition> switchCases = new LinkedHashMap<>();
        for (DiscriminatedComplexTypeDefinition switchCase : allSwitchCases) {
            if (!switchCases.containsKey(switchCase.getName())) {
                switchCases.put(switchCase.getName(), switchCase);
            }
        }
        return new ArrayList<>(switchCases.values());
    }

    public void emitRequiredImportRaw(String requiredImport) {
        LOGGER.debug("emitting import '{}'", requiredImport);
        requiredImports.add(requiredImport);
    }

    public void emitRequiredImport(String requiredImport) {
        LOGGER.debug("emitting import '\"{}\"'", requiredImport);
        requiredImports.add('"' + requiredImport + '"');
    }

    public void emitRequiredImport(String alias, String requiredImport) {
        LOGGER.debug("emitting import '{} \"{}'\"", alias, requiredImport);
        requiredImports.add(alias + ' ' + '"' + requiredImport + '"');
    }

    public Set<String> getRequiredImports() {
        return requiredImports;
    }

    public void emitDataReaderRequiredImports() {
        requiredImports.add(". \"github.com/apache/plc4x/plc4go/spi/codegen/fields\"");
        requiredImports.add(". \"github.com/apache/plc4x/plc4go/spi/codegen/io\"");
    }

    public void emitCodegenRequiredImports() {
        requiredImports.add("\"github.com/apache/plc4x/plc4go/spi/codegen\"");
    }

    public void emitDataIoRequiredImport(String requiredImport) {
        LOGGER.debug("emitting io import '\"{}\"'", requiredImport);
        requiredImportsForDataIo.add('"' + requiredImport + '"');
    }

    public void emitDataIoRequiredImport(String alias, String requiredImport) {
        LOGGER.debug("emitting data io import '{} \"{}'\"", alias, requiredImport);
        requiredImportsForDataIo.add(alias + ' ' + '"' + requiredImport + '"');
    }

    public Set<String> getRequiredImportsForDataIo() {
        return requiredImportsForDataIo;
    }

    public String getVariableName(Field field) {
        if (!(field instanceof NamedField)) {
            return "_";
        }
        NamedField namedField = (NamedField) field;

        String name = null;
        for (Field curField : ((ComplexTypeDefinition) thisType).getFields()) {
            if (curField == field) {
                name = namedField.getName();
            } else if (name != null) {
                if (curField instanceof ArrayField) {
                    ArrayField arrayField = (ArrayField) curField;
                    if (arrayField.getLoopExpression().contains(name)) {
                        return name;
                    }
                } else if (curField instanceof ChecksumField) {
                    ChecksumField checksumField = (ChecksumField) curField;
                    if (checksumField.getChecksumExpression().contains(name)) {
                        return name;
                    }
                } else if (curField instanceof ImplicitField) {
                    ImplicitField implicitField = (ImplicitField) curField;
                    if (implicitField.getSerializeExpression().contains(name)) {
                        return name;
                    }
                } else if (curField instanceof ManualArrayField) {
                    ManualArrayField manualArrayField = (ManualArrayField) curField;
                    if (manualArrayField.getLengthExpression().contains(name)) {
                        return name;
                    }
                    if (manualArrayField.getLoopExpression().contains(name)) {
                        return name;
                    }
                    if (manualArrayField.getParseExpression().contains(name)) {
                        return name;
                    }
                    if (manualArrayField.getSerializeExpression().contains(name)) {
                        return name;
                    }
                } else if (curField instanceof ManualField) {
                    ManualField manualField = (ManualField) curField;
                    if (manualField.getLengthExpression().contains(name)) {
                        return name;
                    }
                    if (manualField.getParseExpression().contains(name)) {
                        return name;
                    }
                    if (manualField.getSerializeExpression().contains(name)) {
                        return name;
                    }
                } else if (curField instanceof OptionalField) {
                    OptionalField optionalField = (OptionalField) curField;
                    if (optionalField.getConditionExpression().isPresent() && optionalField.getConditionExpression().orElseThrow(IllegalStateException::new).contains(name)) {
                        return name;
                    }
                } else if (curField instanceof SwitchField) {
                    SwitchField switchField = (SwitchField) curField;
                    for (Term discriminatorExpression : switchField.getDiscriminatorExpressions()) {
                        if (discriminatorExpression.contains(name)) {
                            return name;
                        }
                    }
                    for (DiscriminatedComplexTypeDefinition curCase : switchField.getCases()) {
                        for (Argument parserArgument : curCase.getParserArguments().orElse(Collections.emptyList())) {
                            if (parserArgument.getName().equals(name)) {
                                return name;
                            }
                        }
                    }
                } else if (curField instanceof VirtualField) {
                    VirtualField virtualField = (VirtualField) curField;
                    if (virtualField.getValueExpression().contains(name)) {
                        return name;
                    }
                }
                List<Term> params = field.asTypedField()
                    .map(typedField -> typedField.getType().asNonSimpleTypeReference()
                        .map(NonSimpleTypeReference::getParams)
                        .map(terms -> terms.orElse(Collections.emptyList()))
                        .orElse(Collections.emptyList())
                    )
                    .orElse(Collections.emptyList());
                for (Term param : params) {
                    if (param.contains(name)) {
                        return name;
                    }
                }
            }
        }

        return "_";
    }

    public boolean needsVariable(Field field, String variableName, boolean serialization) {
        if (!serialization) {
            if (field instanceof ArrayField) {
                ArrayField arrayField = (ArrayField) field;
                if (arrayField.getLoopExpression().contains(variableName)) {
                    return true;
                }
            }
        }
        if (field instanceof VirtualField) {
            VirtualField virtualField = (VirtualField) field;
            if (virtualField.getValueExpression().contains(variableName)) {
                return true;
            }
        }
        if (field instanceof PaddingField) {
            PaddingField paddingField = (PaddingField) field;
            if (paddingField.getPaddingCondition().contains(variableName)) {
                return true;
            }
            if (paddingField.getPaddingValue().contains(variableName)) {
                return true;
            }
        }
        return field.asTypedField()
            .map(typedField -> typedField.getType().asNonSimpleTypeReference()
                .map(nonSimpleTypeReference -> nonSimpleTypeReference.getParams()
                    .map(params -> params.stream()
                        .anyMatch(param -> param.contains(variableName))
                    )
                    .orElse(false)
                )
                .orElse(false)
            )
            .orElse(false);
    }

    /**
     * Right now only the ARRAY_SIZE_IN_BYTES requires helpers to be generated.
     * Also right now only the Modbus protocol requires this and here the referenced
     * properties are all also members of the current complex type,
     * so we'll simplify things here for now.
     *
     * @param functionName name of the
     * @return something
     */
    public Map<String, String> requiresHelperFunctions(String functionName) {
        Map<String, String> result = new HashMap<>();
        boolean usesFunction = false;
        // As the ARRAY_SIZE_IN_BYTES only applies to ArrayFields, search for these
        for (Field curField : ((ComplexTypeDefinition) thisType).getFields()) {
            if (curField instanceof ArrayField) {
                ArrayField arrayField = (ArrayField) curField;
                if (arrayField.getLoopExpression().contains(functionName)) {
                    usesFunction = true;
                }
                result.put(arrayField.getName(), getLanguageTypeNameForField(arrayField));
            } else if (curField instanceof ImplicitField) {
                ImplicitField implicitField = (ImplicitField) curField;
                if (implicitField.getSerializeExpression().contains(functionName)) {
                    usesFunction = true;
                }
            }
        }
        if (!usesFunction) {
            return Collections.emptyMap();
        }
        return result;
    }

    public boolean requiresCurPos() {
        if (thisType instanceof ComplexTypeDefinition) {
            ComplexTypeDefinition complexTypeDefinition = (ComplexTypeDefinition) this.thisType;
            for (Field curField : complexTypeDefinition.getFields()) {
                if (requiresVariable(curField, "curPos")) {
                    return true;
                }
            }
        }
        return false;
    }

    public boolean requiresStartPos() {
        if (thisType instanceof ComplexTypeDefinition) {
            ComplexTypeDefinition complexTypeDefinition = (ComplexTypeDefinition) this.thisType;
            for (Field curField : complexTypeDefinition.getFields()) {
                if (requiresVariable(curField, "startPos")) {
                    return true;
                }
            }
        }
        return false;
    }

    public boolean requiresVariable(Field curField, String variable) {
        if (curField.isArrayField()) {
            ArrayField arrayField = (ArrayField) curField;
            if (arrayField.getLoopExpression().contains(variable)) {
                return true;
            }
        } else if (curField.isOptionalField()) {
            OptionalField optionalField = (OptionalField) curField;
            if (optionalField.getConditionExpression().isPresent() && optionalField.getConditionExpression().orElseThrow(IllegalStateException::new).contains(variable)) {
                return true;
            }
        }
        return curField.asTypedField()
            .map(typedField -> typedField.getType().asNonSimpleTypeReference()
                .map(nonSimpleTypeReference -> nonSimpleTypeReference.getParams()
                    .map(params -> params.stream()
                        .anyMatch(param -> param.contains(variable))
                    )
                    .orElse(false)
                )
                .orElse(false)
            )
            .orElse(false);
    }

    public Term findTerm(Term baseTerm, String name) {
        if (baseTerm instanceof VariableLiteral) {
            VariableLiteral variableLiteral = (VariableLiteral) baseTerm;
            if (variableLiteral.getName().equals(name)) {
                return variableLiteral;
            }
            if (variableLiteral.getChild().isPresent()) {
                Term found = findTerm(variableLiteral.getChild().get(), name);
                if (found != null) {
                    return found;
                }
            }
            for (Term arg : variableLiteral.getArgs().orElse(Collections.emptyList())) {
                Term found = findTerm(arg, name);
                if (found != null) {
                    return found;
                }
            }
        } else if (baseTerm instanceof UnaryTerm) {
            UnaryTerm unaryTerm = (UnaryTerm) baseTerm;
            return findTerm(unaryTerm.getA(), name);
        } else if (baseTerm instanceof BinaryTerm) {
            BinaryTerm binaryTerm = (BinaryTerm) baseTerm;
            Term found = findTerm(binaryTerm.getA(), name);
            if (found != null) {
                return found;
            }
            found = findTerm(binaryTerm.getB(), name);
            return found;
        } else if (baseTerm instanceof TernaryTerm) {
            TernaryTerm ternaryTerm = (TernaryTerm) baseTerm;
            Term found = findTerm(ternaryTerm.getA(), name);
            if (found != null) {
                return found;
            }
            found = findTerm(ternaryTerm.getB(), name);
            if (found != null) {
                return found;
            }
            found = findTerm(ternaryTerm.getC(), name);
            return found;
        }
        return null;
    }

    public String getEnumExpression(String expression) {
        String enumName = expression.substring(0, expression.indexOf('.'));
        String enumConstant = expression.substring(expression.indexOf('.') + 1);
        return enumName + "_" + enumConstant;
    }

    public boolean needsReferenceForParserArgument(String propertyName, TypeReference argumentType) {
        return argumentType.asComplexTypeReference()
            .map(complexTypeReference -> {
                // Check if this is a local field.
                // FIXME: shouldn't this look onto the argumentType? this is nowhere used...
                return thisType.asComplexTypeDefinition()
                    .map(
                        complexTypeDefinition -> complexTypeDefinition.getPropertyFieldByName(propertyName)
                            .map(TypedField.class::cast)
                            .map(TypedField::getType)
                            .filter(NonSimpleTypeReference.class::isInstance)
                            .map(NonSimpleTypeReference.class::cast)
                            .map(NonSimpleTypeReference::getTypeDefinition)
                            .map(typeDefinition -> !(typeDefinition instanceof EnumTypeDefinition))
                            .orElse(false)
                    )
                    .orElse(false);
            })
            .orElse(false);
    }

    public String capitalize(String str) {
        Tracer dummyTracer = Tracer.start("");
        String extractedTrace = dummyTracer.extractTraces(str);
        String cleanedString = dummyTracer.removeTraces(str);
        return extractedTrace + StringUtils.capitalize(cleanedString);
    }

    public String getEndiannessOptions(boolean read, boolean separatorPrefix) {
        return getEndiannessOptions(read, separatorPrefix, Collections.emptyList());
    }

    public String getEndiannessOptions(boolean read, boolean separatorPrefix, List<Argument> parserArguments) {
        Optional<Term> byteOrder = thisType.getAttribute("byteOrder");
        if (byteOrder.isPresent()) {
            emitRequiredImport("encoding/binary");
            if (read) {
                return (separatorPrefix ? ", " : "") + "utils.WithByteOrderForReadBufferByteBased(" +
                    toParseExpression(null, new DefaultByteOrderTypeReference(), byteOrder.orElseThrow(), parserArguments) +
                    ")";
            } else {
                return (separatorPrefix ? ", " : "") + "utils.WithByteOrderForByteBasedBuffer(" +
                    toSerializationExpression(null, new DefaultByteOrderTypeReference(), byteOrder.orElseThrow(), parserArguments) +
                    ")";
            }
        }
        return "";
    }

    public String getFieldOptions(TypedField field, List<Argument> parserArguments) {
        StringBuilder sb = new StringBuilder();
        field.getEncoding().ifPresent(term -> {
            emitCodegenRequiredImports();
            final String encoding = toParseExpression(field, field.getType(), term, parserArguments);
            sb.append(", codegen.WithEncoding(").append(encoding).append(")");
        });

        field.getByteOrder().ifPresent(term -> {
            emitCodegenRequiredImports();
            emitRequiredImport("encoding/binary");
            String byteOrder = "binary.BigEndian";
            switch (term.stringRepresentation()) {
                case "BIG_ENDIAN":
                    byteOrder = "binary.BigEndian";
                    break;
                case "LITTLE_ENDIAN":
                    byteOrder = "binary.LittleEndian";
                    break;
                default:
                    throw new RuntimeException("unmapped bytes order " + term.stringRepresentation());
            }
            sb.append(", codegen.WithByteOrder(").append(byteOrder).append(")");
        });

        field.getAttribute("nullBytesHex").ifPresent(term -> {
            emitCodegenRequiredImports();
            final String nullBytesHex = toParseExpression(field, field.getType(), term, parserArguments);
            sb.append(", codegen.WithNullBytesHex(\"").append(nullBytesHex).append("\")");
        });
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

    public boolean isGeneratePropertiesForParserArguments() {
        return options.getOrDefault("generate-properties-for-parser-arguments", "false").equals("true");
    }

    public boolean isGeneratePropertiesForReservedFields() {
        return options.getOrDefault("generate-properties-for-reserved-fields", "false").equals("true");
    }

}

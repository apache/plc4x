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

package org.apache.plc4x.language.cs;

import org.apache.commons.lang3.StringUtils;
import org.apache.commons.lang3.math.NumberUtils;
import org.apache.plc4x.plugins.codegenerator.protocol.freemarker.BaseFreemarkerLanguageTemplateHelper;
import org.apache.plc4x.plugins.codegenerator.types.definitions.*;
import org.apache.plc4x.plugins.codegenerator.types.enums.EnumValue;
import org.apache.plc4x.plugins.codegenerator.types.fields.*;
import org.apache.plc4x.plugins.codegenerator.types.references.*;
import org.apache.plc4x.plugins.codegenerator.types.terms.*;

import java.util.*;

@SuppressWarnings({"unused", "WeakerAccess"})
public class CsLanguageTemplateHelper extends BaseFreemarkerLanguageTemplateHelper {

    public CsLanguageTemplateHelper(TypeDefinition thisType, String protocolName, String flavorName, Map<String, TypeDefinition> types) {
        super(thisType, protocolName, flavorName, types);
    }

    public String fileName(String protocolName, String languageName, String languageFlavorName) {
        return String.join("", protocolName.split("\\-")) + "." +
            String.join("", languageFlavorName.split("\\-"));
    }

    public String packageName(String languageFlavorName) {
        return String.join("", languageFlavorName.split("\\-"));
    }

    @Override
    public String getLanguageTypeNameForField(Field field) {
        boolean optional = field instanceof OptionalField;
        // If the referenced type is a DataIo type, the value is of type PlcValue.
        if(field instanceof PropertyField) {
            PropertyField propertyField = (PropertyField) field;
            if(propertyField.getType() instanceof ComplexTypeReference) {
                ComplexTypeReference complexTypeReference = (ComplexTypeReference) propertyField.getType();
                final TypeDefinition typeDefinition = getTypeDefinitions().get(complexTypeReference.getName());
                if(typeDefinition instanceof DataIoTypeDefinition) {
                    return "PlcValue";
                }
            }
        }
        return getLanguageTypeNameForTypeReference(((TypedField) field).getType());
    }

    @Override
    public String getLanguageTypeNameForTypeReference(TypeReference typeReference) {
        if(typeReference instanceof SimpleTypeReference) {
            SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
            switch (simpleTypeReference.getBaseType()) {
                case BIT: {
                    return "bool";
                }
                case UINT: {
                    IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                    if (integerTypeReference.getSizeInBits() <= 8) {
                        return "byte";
                    }
                    if (integerTypeReference.getSizeInBits() <= 16) {
                        return "ushort";
                    }
                    if (integerTypeReference.getSizeInBits() <= 32) {
                        return "uint";
                    }
                    if (integerTypeReference.getSizeInBits() <= 64) {
                        return "ulong";
                    }
                    throw new RuntimeException("Unsupported simple type");
                }
                case INT: {
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
                    throw new RuntimeException("Unsupported simple type");
                }
                case FLOAT:
                case UFLOAT: {
                    FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                    int sizeInBits = ((floatTypeReference.getBaseType() == SimpleTypeReference.SimpleBaseType.FLOAT) ? 1 : 0) +
                        floatTypeReference.getExponent() + floatTypeReference.getMantissa();
                    if (sizeInBits <= 32) {
                        return "float";
                    }
                    if (sizeInBits <= 64) {
                        return "double";
                    }
                    throw new RuntimeException("Unsupported simple type");
                }
                case STRING: {
                    return "string";
                }
                case TIME: {
                    return "time";
                }
                case DATE: {
                    return "date";
                }
                case DATETIME: {
                    return "datetime2";
                }
            }
            throw new RuntimeException("Unsupported simple type");
        } else {
            return (typeReference != null) ? ((ComplexTypeReference) typeReference).getName() : "";
        }
    }

    public String getPlcValueTypeForTypeReference(TypeReference typeReference) {
        if(typeReference instanceof SimpleTypeReference) {
            SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
            switch (simpleTypeReference.getBaseType()) {
                case BIT: {
                    return "PlcBOOL";
                }
                case UINT: {
                    IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                    if (integerTypeReference.getSizeInBits() <= 8) {
                        return "PlcUSINT";
                    }
                    if (integerTypeReference.getSizeInBits() <= 16) {
                        return "PlcUINT";
                    }
                    if (integerTypeReference.getSizeInBits() <= 32) {
                        return "PlcUDINT";
                    }
                    if (integerTypeReference.getSizeInBits() <= 64) {
                        return "PlcULINT";
                    }
                    throw new RuntimeException("Unsupported simple type");
                }
                case INT: {
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
                    throw new RuntimeException("Unsupported simple type");
                }
                case FLOAT:
                case UFLOAT: {
                    FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                    int sizeInBits = ((floatTypeReference.getBaseType() == SimpleTypeReference.SimpleBaseType.FLOAT) ? 1 : 0) +
                        floatTypeReference.getExponent() + floatTypeReference.getMantissa();
                    if (sizeInBits <= 32) {
                        return "PlcREAL";
                    }
                    if (sizeInBits <= 64) {
                        return "PlcLREAL";
                    }
                    throw new RuntimeException("Unsupported simple type");
                }
                case STRING: {
                    return "PlcSTRING";
                }
                case TIME: {
                    return "PlcTIME";
                }
                case DATE: {
                    return "PlcDATE";
                }
                case DATETIME: {
                    return "PlcDATE_AND_TIME";
                }
            }
            throw new RuntimeException("Unsupported simple type");
        } else {
            return (typeReference != null) ? ((ComplexTypeReference) typeReference).getName() : "";
        }
    }

    @Override
    public String getNullValueForTypeReference(TypeReference typeReference) {
        if(typeReference instanceof SimpleTypeReference) {
            SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
            switch (simpleTypeReference.getBaseType()) {
                case BIT: {
                    return "false";
                }
                case UINT:
                case INT: {
                    return "0";
                }
                case FLOAT: {
                    return "0.0";
                }
                case STRING: {
                    return "\"\"";
                }
            }
        } else if(typeReference instanceof ComplexTypeReference) {
            return "0";
        }
        return "null";
    }

    public int getNumBits(SimpleTypeReference simpleTypeReference) {
        switch (simpleTypeReference.getBaseType()) {
            case BIT: {
                return 1;
            }
            case UINT:
            case INT: {
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                return integerTypeReference.getSizeInBits();
            }
            case FLOAT: {
                FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                return floatTypeReference.getSizeInBits();
            }
            case STRING: {
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                return integerTypeReference.getSizeInBits();
            }
            default: {
                return 0;
            }
        }
    }

    public String getReadBufferReadMethodCall(SimpleTypeReference simpleTypeReference, String valueString) {
        switch (simpleTypeReference.getBaseType()) {
            case BIT: {
                return "io.ReadBit()";
            }
            case UINT: {
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (integerTypeReference.getSizeInBits() <= 8) {
                    return "io.ReadByte(" + integerTypeReference.getSizeInBits() + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    return "io.ReadUshort(" + integerTypeReference.getSizeInBits() + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    return "io.ReadUint(" + integerTypeReference.getSizeInBits() + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 64) {
                    return "io.ReadUlong(" + integerTypeReference.getSizeInBits() + ")";
                }
            }
            case INT: {
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (integerTypeReference.getSizeInBits() <= 8) {
                    return "io.ReadSbyte(" + integerTypeReference.getSizeInBits() + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    return "io.ReadShort(" + integerTypeReference.getSizeInBits() + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    return "io.ReadInt(" + integerTypeReference.getSizeInBits() + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 64) {
                    return "io.ReadLong(" + integerTypeReference.getSizeInBits() + ")";
                }
            }
            case FLOAT: {
                FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                if (floatTypeReference.getSizeInBits() <= 32) {
                    return "io.ReadFloat(true, " + floatTypeReference.getExponent() + ", " + floatTypeReference.getMantissa() + ")";
                }
                if (floatTypeReference.getSizeInBits() <= 64) {
                    return "io.ReadDouble(true, " + floatTypeReference.getExponent() + ", " + floatTypeReference.getMantissa() + ")";
                }
            }
            case STRING: {
                StringTypeReference stringTypeReference = (StringTypeReference) simpleTypeReference;
                return "io.ReadString(" + stringTypeReference.getSizeInBits() + ", Encoding.GetEncoding(\"" +
                    stringTypeReference.getEncoding().substring(1, stringTypeReference.getEncoding().length() - 1) + "\"))";
            }
        }
        return "Hurz";
    }

    @Override
    public String getWriteBufferWriteMethodCall(SimpleTypeReference simpleTypeReference, String fieldName) {
        switch (simpleTypeReference.getBaseType()) {
            case BIT: {
                return "io.WriteBit(" + fieldName + ")";
            }
            case UINT: {
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (integerTypeReference.getSizeInBits() <= 8) {
                    return "io.WriteByte(" + integerTypeReference.getSizeInBits() + ", " + fieldName + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    return "io.WriteUshort(" + integerTypeReference.getSizeInBits() + ", " + fieldName + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    return "io.WriteUint(" + integerTypeReference.getSizeInBits() + ", " + fieldName + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 64) {
                    return "io.WriteUlong(" + integerTypeReference.getSizeInBits() + ", " + fieldName + ")";
                }
            }
            case INT: {
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (integerTypeReference.getSizeInBits() <= 8) {
                    return "io.WriteSbyte(" + integerTypeReference.getSizeInBits() + ", " + fieldName + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    return "io.WriteShort(" + integerTypeReference.getSizeInBits() + ", " + fieldName + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    return "io.WriteInt(" + integerTypeReference.getSizeInBits() + ", " + fieldName + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 64) {
                    return "io.WriteLong(" + integerTypeReference.getSizeInBits() + ", " + fieldName + ")";
                }
            }
            case FLOAT:
            case UFLOAT: {
                FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                if (floatTypeReference.getSizeInBits() <= 32) {
                    return "io.WriteFloat(" + floatTypeReference.getSizeInBits() + ", " + fieldName + ")";
                }
                if (floatTypeReference.getSizeInBits() <= 64) {
                    return "io.WriteDouble(" + floatTypeReference.getSizeInBits() + ", " + fieldName + ")";
                }
            }
            case STRING: {
                StringTypeReference stringTypeReference = (StringTypeReference) simpleTypeReference;
                String encoding = ((stringTypeReference.getEncoding() != null) && (stringTypeReference.getEncoding().length() > 2)) ?
                    stringTypeReference.getEncoding().substring(1, stringTypeReference.getEncoding().length() - 1) : "UTF-8";
                return "io.WriteString(" + stringTypeReference.getSizeInBits() + ", \"" +
                    encoding + "\", " + fieldName + ")";
            }
        }
        return "Hurz";
    }

    public String getReservedValue(ReservedField reservedField) {
        final String languageTypeName = getLanguageTypeNameForTypeReference(reservedField.getType());
        return reservedField.getReferenceValue().toString();
    }

    public String toParseExpression(TypedField field, Term term, Argument[] parserArguments) {
        return toTypedParseExpression((field != null) ? field.getType() : null, term, parserArguments);
    }

    public String toSerializationExpression(TypedField field, Term term, Argument[] serializerArguments) {
        return toTypedSerializationExpression((field != null) ? field.getType() : null, term, serializerArguments);
    }

    public String toBooleanParseExpression(Term term, Argument[] parserArguments) {
        return toTypedParseExpression(new DefaultBooleanTypeReference(), term, parserArguments);
    }

    public String toBooleanSerializationExpression(Term term, Argument[] serializerArguments) {
        return toTypedSerializationExpression(new DefaultBooleanTypeReference(), term, serializerArguments);
    }

    public String toIntegerParseExpression(int sizeInBits, Term term, Argument[] parserArguments) {
        return toTypedParseExpression(new DefaultIntegerTypeReference(SimpleTypeReference.SimpleBaseType.UINT, sizeInBits), term, parserArguments);
    }

    public String toIntegerSerializationExpression(int sizeInBits, Term term, Argument[] serializerArguments) {
        return toTypedSerializationExpression(new DefaultIntegerTypeReference(SimpleTypeReference.SimpleBaseType.UINT, sizeInBits), term, serializerArguments);
    }

    public String toTypedParseExpression(TypeReference fieldType, Term term, Argument[] parserArguments) {
        return toExpression(fieldType, term, parserArguments, null, false, false);
    }

    public String toTypedSerializationExpression(TypeReference fieldType, Term term, Argument[] serializerArguments) {
        return toExpression(fieldType, term, null, serializerArguments, true, false);
    }

    String getCastExpressionForTypeReference(TypeReference typeReference) {
        if(typeReference instanceof SimpleTypeReference) {
            return getLanguageTypeNameForTypeReference(typeReference);
        } else if(typeReference != null) {
            return "Cast" + getLanguageTypeNameForTypeReference(typeReference);
        } else {
            return "";
        }
    }

    private String toExpression(TypeReference fieldType, Term term, Argument[] parserArguments, Argument[] serializerArguments, boolean serialize, boolean suppressPointerAccess) {
        if(term == null) {
            return "";
        }
        if(term instanceof Literal) {
            if(term instanceof NullLiteral) {
                return "nil";
            } else if(term instanceof BooleanLiteral) {
                return getCastExpressionForTypeReference(fieldType) + "(" + ((BooleanLiteral) term).getValue() + ")";
            } else if(term instanceof NumericLiteral) {
                return getCastExpressionForTypeReference(fieldType) + "(" + ((NumericLiteral) term).getNumber().toString() + ")";
            } else if(term instanceof StringLiteral) {
                return "\"" + ((StringLiteral) term).getValue() + "\"";
            } else if(term instanceof VariableLiteral) {
                return toVariableExpression(fieldType, (VariableLiteral) term, parserArguments, serializerArguments, serialize, suppressPointerAccess);
            } else {
                throw new RuntimeException("Unsupported Literal type " + term.getClass().getName());
            }
        } else if (term instanceof UnaryTerm) {
            UnaryTerm ut = (UnaryTerm) term;
            Term a = ut.getA();
            switch(ut.getOperation()) {
                case "!":
                    return "!(" + toExpression(fieldType, a, parserArguments, serializerArguments, serialize, false) + ")";
                case "-":
                    return "-(" + getCastExpressionForTypeReference(fieldType) + "(" + toExpression(fieldType, a, parserArguments, serializerArguments, serialize, false) + "))";
                case "()":
                    return getCastExpressionForTypeReference(fieldType) + "(" + toExpression(fieldType, a, parserArguments, serializerArguments, serialize, false) + ")";
                default:
                    throw new RuntimeException("Unsupported unary operation type " + ut.getOperation());
            }
        } else if (term instanceof BinaryTerm) {
            BinaryTerm bt = (BinaryTerm) term;
            Term a = bt.getA();
            Term b = bt.getB();
            String operation = bt.getOperation();
            switch (operation) {
                case "^":
                    return "Math.pow(" +
                        getCastExpressionForTypeReference(fieldType) + "(" + toExpression(fieldType, a, parserArguments, serializerArguments, serialize, false) + "), " +
                        getCastExpressionForTypeReference(fieldType) + "(" + toExpression(fieldType, b, parserArguments, serializerArguments, serialize, false) + "))";
                // If we start casting for comparisons, equals or non equals, really messy things happen.
                case "==":
                case "!=":
                case ">":
                case "<":
                case ">=":
                case "<=":
                    // For every access of optional elements we need pointer access ...
                    // Except for doing a nil or not-nil check :-(
                    // So in case of such a check, we need to suppress the pointer-access.
                    boolean suppressPointerAccessOverride = (operation.equals("==") || operation.equals("!=")) && ((a instanceof NullLiteral) || (b instanceof NullLiteral));
                    return "bool((" + toExpression(null, a, parserArguments, serializerArguments, serialize, suppressPointerAccessOverride) + ") " +
                        operation +
                        " (" + toExpression(null, b, parserArguments, serializerArguments, serialize, suppressPointerAccessOverride) + "))";
                default:
                    return getCastExpressionForTypeReference(fieldType) + "(" + toExpression(fieldType, a, parserArguments, serializerArguments, serialize, false) + ") " +
                        operation + " " +
                        getCastExpressionForTypeReference(fieldType) +"(" + toExpression(fieldType, b, parserArguments, serializerArguments, serialize, false) + ")";
            }
        } else if (term instanceof TernaryTerm) {
            TernaryTerm tt = (TernaryTerm) term;
            if ("if".equals(tt.getOperation())) {
                Term a = tt.getA();
                Term b = tt.getB();
                Term c = tt.getC();
                // TODO: This is not quite correct with the cast to uint16
                return "utils.InlineIf(" + toExpression(new DefaultBooleanTypeReference(), a, parserArguments, serializerArguments, serialize, false) + ", " +
                    "uint16(" + toExpression(fieldType, b, parserArguments, serializerArguments, serialize, false) + "), " +
                    "uint16(" + toExpression(fieldType, c, parserArguments, serializerArguments, serialize, false) + "))";
            } else {
                throw new RuntimeException("Unsupported ternary operation type " + tt.getOperation());
            }
        } else {
            throw new RuntimeException("Unsupported Term type " + term.getClass().getName());
        }
    }

    private String toVariableExpression(TypeReference typeReference, VariableLiteral vl, Argument[] parserArguments, Argument[] serializerArguments, boolean serialize, boolean suppressPointerAccess) {
        if ("lengthInBytes".equals(vl.getName())) {
            return (serialize ? getCastExpressionForTypeReference(typeReference) + "(m." : "") + "LengthInBytes()" + (serialize ? ")" : "");
        } else if ("lengthInBits".equals(vl.getName())) {
            return (serialize ? getCastExpressionForTypeReference(typeReference) + "(m." : "") + "LengthInBits()" + (serialize ? ")" : "");
        }
        // If this literal references an Enum type, then we have to output it differently.
        else if (getTypeDefinitions().get(vl.getName()) instanceof EnumTypeDefinition) {
            return vl.getName() + "_" + vl.getChild().getName() +
                ((vl.getChild().getChild() != null) ?
                    "." + toVariableExpression(typeReference, vl.getChild().getChild(), parserArguments, serializerArguments, false, suppressPointerAccess) : "");
        }
        // If we are accessing enum constants, these also need to be output differently.
        else if ((getFieldForNameFromCurrent(vl.getName()) instanceof EnumField) && (vl.getChild() != null)) {
            return vl.getName() + "." + StringUtils.capitalize(vl.getChild().getName()) + "()" +
                ((vl.getChild().getChild() != null) ?
                    "." + toVariableExpression(typeReference, vl.getChild().getChild(), parserArguments, serializerArguments, false, suppressPointerAccess) : "");
        }
        // If we are accessing optional fields of simple type, we need to use pointer-access.
        else if (!serialize && (getFieldForNameFromCurrent(vl.getName()) instanceof OptionalField) &&
            (((OptionalField) getFieldForNameFromCurrent(vl.getName())).getType() instanceof SimpleTypeReference)) {
            return "(*" + vl.getName() + ")" +
                ((vl.getChild() != null) ?
                    "." + toVariableExpression(typeReference, vl.getChild(), parserArguments, serializerArguments, serialize, suppressPointerAccess) : "");
        }
        // CAST expressions are special as we need to add a ".class" to the second parameter in Java.
        else if ("CAST".equals(vl.getName())) {
            if ((vl.getArgs() == null) || (vl.getArgs().size() != 2)) {
                throw new RuntimeException("A CAST expression expects exactly two arguments.");
            }
            VariableLiteral typeLiteral = (VariableLiteral) vl.getArgs().get(1);
            final TypeDefinition typeDefinition = getTypeDefinitions().get(typeLiteral.getName());
            TypeReference type = typeDefinition.getTypeReference();
            StringBuilder sb = new StringBuilder();
            if (type instanceof ComplexTypeReference) {
                sb.append("Cast");
            }
            sb.append(typeLiteral.getName());
            sb.append("(").append(toVariableExpression(typeReference, (VariableLiteral) vl.getArgs().get(0), parserArguments, serializerArguments, serialize, suppressPointerAccess)).append(")");
            return sb.toString() + ((vl.getChild() != null) ? "." + StringUtils.capitalize(toVariableExpression(typeReference, vl.getChild(), parserArguments, serializerArguments, false, suppressPointerAccess)) : "");
        } else if ("STATIC_CALL".equals(vl.getName())) {
            StringBuilder sb = new StringBuilder();
            if (!(vl.getArgs().get(0) instanceof StringLiteral)) {
                throw new RuntimeException("Expecting the first argument of a 'STATIC_CALL' to be a StringLiteral");
            }
            // Get the class and method name
            String staticCall = ((StringLiteral) vl.getArgs().get(0)).getValue();
            // Cut off the double-quotes
            staticCall = staticCall.substring(1, staticCall.length() - 1);
            // Remove all the previous parts prior to the Class name (Which starts with an uppercase letter)
            while(staticCall.contains(".") && !StringUtils.isAllUpperCase(staticCall.substring(0,1))) {
                staticCall = staticCall.substring(staticCall.indexOf(".") + 1);
            }
            String className = staticCall.substring(0, staticCall.indexOf("."));
            String methodName = staticCall.substring(staticCall.indexOf(".") + 1);
            sb.append(className).append(StringUtils.capitalize(methodName)).append("(");
            for (int i = 1; i < vl.getArgs().size(); i++) {
                Term arg = vl.getArgs().get(i);
                if (i > 1) {
                    sb.append(", ");
                }
                if (arg instanceof VariableLiteral) {
                    VariableLiteral va = (VariableLiteral) arg;
                    // "io" is the default name of the reader argument which is always available.
                    boolean isParserArg = "io".equals(va.getName()) || ((getThisTypeDefinition() instanceof DataIoTypeDefinition) && "_value".equals(va.getName()));
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
                        if(va.getName().equals("_value")) {
                            sb.append(va.getName().substring(1) + ((va.getChild() != null) ?
                                "." + toVariableExpression(typeReference, vl.getChild(), parserArguments, serializerArguments, false, suppressPointerAccess) : ""));
                        } else {
                            sb.append(va.getName() + ((va.getChild() != null) ?
                                "." + toVariableExpression(typeReference, vl.getChild(), parserArguments, serializerArguments, false, suppressPointerAccess) : ""));
                        }
                    }
                    // We have to manually evaluate the type information at code-generation time.
                    else if (isTypeArg) {
                        String part = va.getChild().getName();
                        switch (part) {
                            case "name":
//                                sb.append("\"").append(field.getTypeName()).append("\"");
                                break;
                            case "length":
                                sb.append("\"").append(((SimpleTypeReference) typeReference).getSizeInBits()).append("\"");
                                break;
                            case "encoding":
                                String encoding = ((StringTypeReference) typeReference).getEncoding();
                                // Cut off the single quotes.
                                encoding = encoding.substring(1, encoding.length() - 1);
                                sb.append("\"").append(encoding).append("\"");
                                break;
                        }
                    } else {
                        sb.append(toVariableExpression(typeReference, va, parserArguments, serializerArguments, serialize, suppressPointerAccess));
                    }
                } else if (arg instanceof StringLiteral) {
                    sb.append(((StringLiteral) arg).getValue());
                }
            }
            sb.append(")");
            return sb.toString();
        } else if ("COUNT".equals(vl.getName())) {
            return (typeReference instanceof SimpleTypeReference ? getCastExpressionForTypeReference(typeReference) : "") + "(len(" +
                toVariableExpression(typeReference, (VariableLiteral) vl.getArgs().get(0), parserArguments, serializerArguments, serialize, suppressPointerAccess) +
                "))";
        } else if ("ARRAY_SIZE_IN_BYTES".equals(vl.getName())) {
            VariableLiteral va = (VariableLiteral) vl.getArgs().get(0);
            // "io" and "m" are always available in every parser.
            boolean isSerializerArg = "io".equals(va.getName()) || "m".equals(va.getName()) || "element".equals(va.getName());
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
                sb.append(va.getName()).append(((va.getChild() != null) ? "." + toVariableExpression(typeReference, va.getChild(), parserArguments, serializerArguments, true, suppressPointerAccess) : ""));
            } else {
                sb.append(toVariableExpression(typeReference, va, parserArguments, serializerArguments, true, suppressPointerAccess));
            }
            return getCastExpressionForTypeReference(typeReference) + "(" + ((VariableLiteral) vl.getArgs().get(0)).getName() + "ArraySizeInBytes(" + sb.toString() + "))";
        }
        else if("CEIL".equals(vl.getName())) {
            Term va = vl.getArgs().get(0);
            // The Ceil function expects 64 bit floating point values.
            TypeReference tr = new DefaultFloatTypeReference(SimpleTypeReference.SimpleBaseType.FLOAT, 11, 52);
            return "math.Ceil(" + toExpression(tr, va, parserArguments, serializerArguments, serialize, suppressPointerAccess) + ")";
        }
        // All uppercase names are not fields, but utility methods.
        else if (vl.getName().equals(vl.getName().toUpperCase())) {
            StringBuilder sb = new StringBuilder(vl.getName());
            if (vl.getArgs() != null) {
                sb.append("(");
                boolean firstArg = true;
                for (Term arg : vl.getArgs()) {
                    if (!firstArg) {
                        sb.append(", ");
                    }
                    sb.append(toExpression(typeReference, arg, parserArguments, serializerArguments, serialize, suppressPointerAccess));
                    firstArg = false;
                }
                sb.append(")");
            }
            if (vl.getIndex() != VariableLiteral.NO_INDEX) {
                sb.append("[").append(vl.getIndex()).append("]");
            }
            return sb.toString() + ((vl.getChild() != null) ?
                "." + toVariableExpression(typeReference, vl.getChild(), parserArguments, serializerArguments, false, suppressPointerAccess) : "");
        }
        // If the current property references a discriminator value, we have to serialize it differently.
        else if ((getFieldForNameFromCurrentOrParent(vl.getName()) != null) && (getFieldForNameFromCurrentOrParent(vl.getName()) instanceof DiscriminatorField)) {
            final DiscriminatorField discriminatorField = (DiscriminatorField) getFieldForNameFromCurrentOrParent(vl.getName());
            System.out.println(discriminatorField);
        }
        // If the current property references a parserArguments property and that is a discriminator property, we also have to serialize it differently..
        else if ((vl.getChild() != null) && (getTypeReferenceForProperty(((ComplexTypeDefinition) getThisTypeDefinition()), vl.getName()) != null)) {
            final Optional<TypeReference> typeReferenceForProperty = getTypeReferenceForProperty(((ComplexTypeDefinition) getThisTypeDefinition()), vl.getName());
            if(typeReferenceForProperty.isPresent() && typeReferenceForProperty.get() instanceof ComplexTypeReference) {
                final TypeReference complexTypeReference = typeReferenceForProperty.get();
                TypeDefinition typeDefinition = getTypeDefinitionForTypeReference(complexTypeReference);
                if(typeDefinition instanceof ComplexTypeDefinition) {
                    ComplexTypeDefinition complexTypeDefinition = (ComplexTypeDefinition) typeDefinition;
                    String childProperty = vl.getChild().getName();
                    final Optional<Field> matchingDiscriminatorField = complexTypeDefinition.getFields().stream().filter(field -> (field instanceof DiscriminatorField) && ((DiscriminatorField) field).getName().equals(childProperty)).findFirst();
                    if(matchingDiscriminatorField.isPresent()) {
                        return "Cast" + getLanguageTypeNameForTypeReference(complexTypeReference) + "(" + vl.getName() + ")." + StringUtils.capitalize(childProperty) + "()";
                    }
                }
            }
        }
        // If the current term references a serialization argument, handle it differently (don't prefix it with "m.")
        else if((serializerArguments != null) && Arrays.stream(serializerArguments).anyMatch(argument -> argument.getName().equals(vl.getName()))) {
            return vl.getName() + ((vl.getChild() != null) ?
                "." + toVariableExpression(typeReference, vl.getChild(), parserArguments, serializerArguments, false, suppressPointerAccess) : "");
        }
        return (serialize ? "m." + StringUtils.capitalize(vl.getName()) : vl.getName()) + ((vl.getChild() != null) ?
            "." + toVariableExpression(typeReference, vl.getChild(), parserArguments, serializerArguments, false, suppressPointerAccess) : "");
    }

    public String getSizeInBits(ComplexTypeDefinition complexTypeDefinition, Argument[] parserArguments) {
        int sizeInBits = 0;
        StringBuilder sb = new StringBuilder("");
        for (Field field : complexTypeDefinition.getFields()) {
            if(field instanceof ArrayField) {
                ArrayField arrayField = (ArrayField) field;
                final SimpleTypeReference type = (SimpleTypeReference) arrayField.getType();
                switch (arrayField.getLoopType()) {
                    case COUNT:
                        sb.append("(").append(toTypedSerializationExpression(type, arrayField.getLoopExpression(), parserArguments)).append(" * ").append(type.getSizeInBits()).append(") + ");
                        break;
                    case LENGTH:
                        sb.append("(").append(toTypedSerializationExpression(type, arrayField.getLoopExpression(), parserArguments)).append(" * 8) + ");
                        break;
                    case TERMINATED:
                        // No terminated.
                        break;
                }
            } else if(field instanceof TypedField) {
                TypedField typedField = (TypedField) field;
                final TypeReference type = typedField.getType();
                if(field instanceof ManualField) {
                    ManualField manualField = (ManualField) field;
                    sb.append("(").append(toSerializationExpression(manualField, manualField.getLengthExpression(), parserArguments)).append(") + ");
                }
                else if(type instanceof SimpleTypeReference) {
                    SimpleTypeReference simpleTypeReference = (SimpleTypeReference) type;
                    sizeInBits += simpleTypeReference.getSizeInBits();
                } else {
                    // No ComplexTypeReference supported.
                }
            }
        }
        return sb.toString() + sizeInBits;
    }

    public String escapeValue(TypeReference typeReference, String valueString) {
        if(valueString == null) {
            return null;
        }
        if(typeReference instanceof SimpleTypeReference) {
            SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
            switch (simpleTypeReference.getBaseType()) {
                case UINT:
                case INT:
                    // If it's a one character string and is numeric, output it as char.
                    if(!NumberUtils.isParsable(valueString) && (valueString.length() == 1)) {
                        return "'" + valueString + "'";
                    }
                    break;
                case STRING:
                    return "\"" + valueString + "\"";
            }
        }
        return valueString;
    }

    public String escapeEnumValue(TypeReference typeReference, String valueString) {
        // Currently the only case in which here complex type references are used are when referencing enum constants.
        if (typeReference instanceof ComplexTypeReference) {
            // C doesn't like NULL values for enums, so we have to return something else (we'll treat -1 as NULL)
            if ("null".equals(valueString)) {
                return "0";
            }
            String typeName = valueString.substring(0, valueString.indexOf('.'));
            String constantName = valueString.substring(valueString.indexOf('.') + 1);
            return typeName + "_" + constantName;
        } else {
            return escapeValue(typeReference, valueString);
        }
    }

    public Collection<EnumValue> getUniqueEnumValues(EnumValue[] enumValues) {
        Map<String, EnumValue> filteredEnumValues = new TreeMap<>();
        for (EnumValue enumValue : enumValues) {
            if (!filteredEnumValues.containsKey(enumValue.getValue())) {
                filteredEnumValues.put(enumValue.getValue(), enumValue);
            }
        }
        return filteredEnumValues.values();
    }

    public List<String> getRequiredImports() {
        ComplexTypeDefinition complexTypeDefinition = (ComplexTypeDefinition) getThisTypeDefinition();
        List<String> imports = new ArrayList<>();

        if(complexTypeDefinition.getAllPropertyFields().stream().anyMatch(field -> isArrayField(field) && getLanguageTypeNameForField(field).equals("int8"))) {
            imports.add("\"encoding/base64\"");
        }

        imports.add("\"encoding/xml\"");

        // For "Fields with complex type", constant, typeSwitch,  fields: "errors"
        if(!complexTypeDefinition.getFields().isEmpty()) {
            imports.add("\"errors\"");
        }

        imports.add("\"io\"");

        // At least one reserved field or simple field with complex type
        if(complexTypeDefinition.getFields().stream().anyMatch(field ->
            (field instanceof ReservedField))) {
            imports.add("log \"github.com/sirupsen/logrus\"");
        }

        // For CEIL functions: "math"
/*        if(complexTypeDefinition.getFields().stream().anyMatch(field ->
            FieldUtils.contains(field, "CEIL"))) {
            imports.add("\"math\"");
        }*/

        imports.add("\"github.com/apache/plc4x/plc4go/internal/plc4go/utils\"");

        // For Constant field: "strconv"
        if(complexTypeDefinition.getFields().stream().anyMatch(field ->
            (field instanceof ConstField))/* || complexTypeDefinition.getAllPropertyFields().stream().anyMatch(
                propertyField -> isSimpleField(propertyField))*/) {
            imports.add("\"strconv\"");
        }

        if(isDiscriminatedParentTypeDefinition()) {
            imports.add("\"reflect\"");
            imports.add("\"strings\"");
        }

        return imports;
    }

    public List<String> getRequiredImportsForDataIo() {
        DataIoTypeDefinition dataIo = (DataIoTypeDefinition) getThisTypeDefinition();

        List<String> imports = new ArrayList<>();

        imports.add("\"errors\"");
        imports.add("\"github.com/apache/plc4x/plc4go/internal/plc4go/model/values\"");
        imports.add("\"github.com/apache/plc4x/plc4go/internal/plc4go/utils\"");
        imports.add("api \"github.com/apache/plc4x/plc4go/pkg/plc4go/values\"");

        if(dataIo.getSwitchField().getCases().stream().anyMatch(typeCase ->
            (typeCase.getName().equals("TIME_OF_DAY") && hasFieldsWithNames(typeCase.getFields(), "hour", "minutes", "seconds")) ||
                (typeCase.getName().equals("DATE") && hasFieldsWithNames(typeCase.getFields(), "year", "month", "day")) ||
                (typeCase.getName().equals("DATE_AND_TIME") && hasFieldsWithNames(typeCase.getFields(), "year", "month", "day", "hour", "minutes", "seconds")))) {
            imports.add("\"time\"");
        }
        return imports;
    }

    public String getVariableName(Field field) {
        if(!(field instanceof NamedField)) {
            return "_";
        }
        NamedField namedField = (NamedField) field;

        String name = null;
        for (Field curField : ((ComplexTypeDefinition) getThisTypeDefinition()).getFields()) {
            if(curField == field) {
                name = namedField.getName();
            } else if(name != null) {
                if(curField instanceof ArrayField) {
                    ArrayField arrayField = (ArrayField) curField;
                    if(arrayField.getLoopExpression().contains(name)) {
                        return name;
                    }
                } else if(curField instanceof ChecksumField) {
                    ChecksumField checksumField = (ChecksumField) curField;
                    if(checksumField.getChecksumExpression().contains(name)) {
                        return name;
                    }
                } else if(curField instanceof ImplicitField) {
                    ImplicitField implicitField = (ImplicitField) curField;
                    if(implicitField.getSerializeExpression().contains(name)) {
                        return name;
                    }
                } else if(curField instanceof ManualArrayField) {
                    ManualArrayField manualArrayField = (ManualArrayField) curField;
                    if(manualArrayField.getLengthExpression().contains(name)) {
                        return name;
                    }
                    if(manualArrayField.getLoopExpression().contains(name)) {
                        return name;
                    }
                    if(manualArrayField.getParseExpression().contains(name)) {
                        return name;
                    }
                    if(manualArrayField.getSerializeExpression().contains(name)) {
                        return name;
                    }
                } else if(curField instanceof ManualField) {
                    ManualField manualField = (ManualField) curField;
                    if(manualField.getLengthExpression().contains(name)) {
                        return name;
                    }
                    if(manualField.getParseExpression().contains(name)) {
                        return name;
                    }
                    if(manualField.getSerializeExpression().contains(name)) {
                        return name;
                    }
                } else if(curField instanceof OptionalField) {
                    OptionalField optionalField = (OptionalField) curField;
                    if(optionalField.getConditionExpression().contains(name)) {
                        return name;
                    }
                } else if(curField instanceof SwitchField) {
                    SwitchField switchField = (SwitchField) curField;
                    for (Term discriminatorExpression : switchField.getDiscriminatorExpressions()) {
                        if(discriminatorExpression.contains(name)) {
                            return name;
                        }
                    }
                    for (DiscriminatedComplexTypeDefinition curCase : switchField.getCases()) {
                        for (Argument parserArgument : curCase.getParserArguments()) {
                            if(parserArgument.getName().equals(name)) {
                                return name;
                            }
                        }
                    }
                } else if(curField instanceof VirtualField) {
                    VirtualField virtualField = (VirtualField) curField;
                    if(virtualField.getValueExpression().contains(name)) {
                        return name;
                    }
                }
                if(curField.getParams() != null) {
                    for (Term param : curField.getParams()) {
                        if(param.contains(name)) {
                            return name;
                        }
                    }
                }
            }
        }

        return "_";
    }

    public boolean needsVariable(ArrayField field, String variableName, boolean serialization) {
        if(!serialization) {
            if (field.getLoopExpression().contains(variableName)) {
                return true;
            }
        }
        if((field.getParams() != null) && (field.getParams().length > 0)){
            for (Term param : field.getParams()) {
                if(param.contains(variableName)) {
                    return true;
                }
            }
        }
        return false;
    }

    /**
     * Right now only the ARRAY_SIZE_IN_BYTES requires helpers to be generated.
     * Also right now only the Modbus protocol requires this and here the referenced
     * properties are all also members of the current complex type,
     * so we'll simplify things here for now.
     *
     * @param functionName name of the
     * @return
     */
    public Map<String, String> requiresHelperFunctions(String functionName) {
        Map<String, String> result = new HashMap<>();
        boolean usesFunction = false;
        // As the ARRAY_SIZE_IN_BYTES only applies to ArrayFields, search for these
        for (Field curField : ((ComplexTypeDefinition) getThisTypeDefinition()).getFields()) {
            if(curField instanceof ArrayField) {
                ArrayField arrayField = (ArrayField) curField;
                if(arrayField.getLoopExpression().contains(functionName)) {
                    usesFunction = true;
                }
                result.put(arrayField.getName(), getLanguageTypeNameForField(arrayField));
            } else if(curField instanceof ImplicitField) {
                ImplicitField implicitField = (ImplicitField) curField;
                if(implicitField.getSerializeExpression().contains(functionName)) {
                    usesFunction = true;
                }
            }
        }
        if(usesFunction) {
            return result;
        } else {
            return Collections.emptyMap();
        }
    }

    public boolean requiresStartPosAndCurPos() {
        if(getThisTypeDefinition() instanceof ComplexTypeDefinition) {
            for (Field curField : ((ComplexTypeDefinition) getThisTypeDefinition()).getFields()) {
                if (requiresVariable(curField, "curPos")) {
                    return true;
                }
            }
        }
        return false;
    }

    public boolean requiresVariable(Field curField, String variable) {
        if(curField instanceof ArrayField) {
            ArrayField arrayField = (ArrayField) curField;
            if(arrayField.getLoopExpression().contains(variable)) {
                return true;
            }
        } else if(curField instanceof OptionalField) {
            OptionalField optionalField = (OptionalField) curField;
            if(optionalField.getConditionExpression().contains(variable)) {
                return true;
            }
        }
        if(curField.getParams() != null) {
            for (Term paramTerm : curField.getParams()) {
                if (paramTerm.contains(variable)) {
                    return true;
                }
            }
        }
        return false;
    }

    public Term findTerm(Term baseTerm, String name) {
        if(baseTerm instanceof VariableLiteral) {
            VariableLiteral variableLiteral = (VariableLiteral) baseTerm;
            if(variableLiteral.getName().equals(name)) {
                return variableLiteral;
            }
            if(variableLiteral.getChild() != null) {
                Term found = findTerm(variableLiteral.getChild(), name);
                if(found != null) {
                    return found;
                }
            }
            for (Term arg : variableLiteral.getArgs()) {
                Term found = findTerm(arg, name);
                if(found != null) {
                    return found;
                }
            }
        } else if(baseTerm instanceof UnaryTerm) {
            UnaryTerm unaryTerm = (UnaryTerm) baseTerm;
            Term found = findTerm(unaryTerm.getA(), name);
            if(found != null) {
                return found;
            }
        } else if(baseTerm instanceof BinaryTerm) {
            BinaryTerm binaryTerm = (BinaryTerm) baseTerm;
            Term found = findTerm(binaryTerm.getA(), name);
            if(found != null) {
                return found;
            }
            found = findTerm(binaryTerm.getB(), name);
            if(found != null) {
                return found;
            }
        } else if(baseTerm instanceof TernaryTerm) {
            TernaryTerm ternaryTerm = (TernaryTerm) baseTerm;
            Term found = findTerm(ternaryTerm.getA(), name);
            if(found != null) {
                return found;
            }
            found = findTerm(ternaryTerm.getB(), name);
            if(found != null) {
                return found;
            }
            found = findTerm(ternaryTerm.getC(), name);
            if(found != null) {
                return found;
            }
        }
        return null;
    }

}

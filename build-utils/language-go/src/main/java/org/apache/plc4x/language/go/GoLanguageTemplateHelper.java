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

package org.apache.plc4x.language.go;

import org.apache.commons.lang3.math.NumberUtils;
import org.apache.plc4x.plugins.codegenerator.protocol.freemarker.BaseFreemarkerLanguageTemplateHelper;
import org.apache.plc4x.plugins.codegenerator.types.definitions.*;
import org.apache.plc4x.plugins.codegenerator.types.enums.EnumValue;
import org.apache.plc4x.plugins.codegenerator.types.fields.*;
import org.apache.plc4x.plugins.codegenerator.types.references.*;
import org.apache.plc4x.plugins.codegenerator.types.terms.*;

import java.util.*;
import java.util.function.Function;

@SuppressWarnings({"unused", "WeakerAccess"})
public class GoLanguageTemplateHelper extends BaseFreemarkerLanguageTemplateHelper {

    public GoLanguageTemplateHelper(TypeDefinition thisType, String protocolName, String flavorName, Map<String, TypeDefinition> types) {
        super(thisType, protocolName, flavorName, types);
    }

    public String fileName(String protocolName, String languageName, String languageFlavorName) {
        return "plc4go." + String.join("", protocolName.split("\\-")) + "." +
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
                        return "uint8";
                    }
                    if (integerTypeReference.getSizeInBits() <= 16) {
                        return "uint16";
                    }
                    if (integerTypeReference.getSizeInBits() <= 32) {
                        return "uint32";
                    }
                    if (integerTypeReference.getSizeInBits() <= 64) {
                        return "uint64";
                    }
                    throw new RuntimeException("Unsupported simple type");
                }
                case INT: {
                    IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
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
                    throw new RuntimeException("Unsupported simple type");
                }
                case FLOAT:
                case UFLOAT: {
                    FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                    int sizeInBits = ((floatTypeReference.getBaseType() == SimpleTypeReference.SimpleBaseType.FLOAT) ? 1 : 0) +
                        floatTypeReference.getExponent() + floatTypeReference.getMantissa();
                    if (sizeInBits <= 32) {
                        return "float32";
                    }
                    if (sizeInBits <= 64) {
                        return "float64";
                    }
                    throw new RuntimeException("Unsupported simple type");
                }
                case STRING: {
                    return "string";
                }
                case TIME: {
                    return "Time";
                }
                case DATE: {
                    return "Date";
                }
                case DATETIME: {
                    return "Date";
                }
            }
            throw new RuntimeException("Unsupported simple type");
        } else {
            return ((ComplexTypeReference) typeReference).getName();
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
        return "nil";
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
                    return "io.ReadUint8(" + integerTypeReference.getSizeInBits() + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    return "io.ReadUint16(" + integerTypeReference.getSizeInBits() + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    return "io.ReadUint32(" + integerTypeReference.getSizeInBits() + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 64) {
                    return "io.ReadUint64(" + integerTypeReference.getSizeInBits() + ")";
                }
            }
            case INT: {
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (integerTypeReference.getSizeInBits() <= 8) {
                    return "io.ReadInt8(" + integerTypeReference.getSizeInBits() + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    return "io.ReadInt16(" + integerTypeReference.getSizeInBits() + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    return "io.ReadInt32(" + integerTypeReference.getSizeInBits() + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 64) {
                    return "io.ReadInt64(" + integerTypeReference.getSizeInBits() + ")";
                }
            }
            case FLOAT: {
                FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                if (floatTypeReference.getSizeInBits() <= 32) {
                    return "io.ReadFloat32(" + floatTypeReference.getSizeInBits() + ")";
                }
                if (floatTypeReference.getSizeInBits() <= 64) {
                    return "io.ReadFloat64(" + floatTypeReference.getSizeInBits() + ")";
                }
            }
            case STRING: {
                StringTypeReference stringTypeReference = (StringTypeReference) simpleTypeReference;
                return "io.ReadString(" + stringTypeReference.getSizeInBits() + ", \"" +
                    stringTypeReference.getEncoding() + "\")";
            }
        }
        return "Hurz";
    }

    @Override
    public String getWriteBufferWriteMethodCall(SimpleTypeReference simpleTypeReference, String fieldName) {
        switch (simpleTypeReference.getBaseType()) {
            case BIT: {
                return "io.WriteBit((bool) " + fieldName + ")";
            }
            case UINT: {
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (integerTypeReference.getSizeInBits() <= 8) {
                    return "io.WriteUint8(" + integerTypeReference.getSizeInBits() + ", " + fieldName + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    return "io.WriteUint16(" + integerTypeReference.getSizeInBits() + ", " + fieldName + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    return "io.WriteUint32(" + integerTypeReference.getSizeInBits() + ", " + fieldName + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 64) {
                    return "io.WriteUint64(" + integerTypeReference.getSizeInBits() + ", " + fieldName + ")";
                }
            }
            case INT: {
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (integerTypeReference.getSizeInBits() <= 8) {
                    return "io.WriteInt8(" + integerTypeReference.getSizeInBits() + ", " + fieldName + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    return "io.WriteInt16(" + integerTypeReference.getSizeInBits() + ", " + fieldName + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    return "io.WriteInt32(" + integerTypeReference.getSizeInBits() + ", " + fieldName + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 64) {
                    return "io.WriteInt64(" + integerTypeReference.getSizeInBits() + ", " + fieldName + ")";
                }
            }
            case FLOAT:
            case UFLOAT: {
                FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                if (floatTypeReference.getSizeInBits() <= 32) {
                    return "io.WriteFloat32(" + floatTypeReference.getSizeInBits() + ", " + fieldName + ")";
                }
                if (floatTypeReference.getSizeInBits() <= 64) {
                    return "io.WriteFloat64(" + floatTypeReference.getSizeInBits() + ", " + fieldName + ")";
                }
            }
            case STRING: {
                StringTypeReference stringTypeReference = (StringTypeReference) simpleTypeReference;
                return "io.WriteString(" + stringTypeReference.getSizeInBits() + ", \"" +
                    stringTypeReference.getEncoding() + "\", " + fieldName + ")";
            }
        }
        return "Hurz";
    }

    public String getReservedValue(ReservedField reservedField) {
        final String languageTypeName = getLanguageTypeNameForTypeReference(reservedField.getType());
        return languageTypeName + "(" + reservedField.getReferenceValue() + ")";
    }

    public String toParseExpression(TypedField field, Term term, Argument[] parserArguments) {
        return toExpression(field, term, term1 -> toVariableParseExpression(field, term1, parserArguments));
    }

    public String toSerializationExpression(TypedField field, Term term, Argument[] parserArguments) {
        return toExpression(field, term, term1 -> toVariableSerializationExpression(field, term1, parserArguments));
    }

    private String toExpression(TypedField field, Term term, Function<Term, String> variableExpressionGenerator) {
        if(term == null) {
            return "";
        }
        if(term instanceof Literal) {
            if(term instanceof NullLiteral) {
                return "nil";
            } else if(term instanceof BooleanLiteral) {
                return Boolean.toString(((BooleanLiteral) term).getValue());
            } else if(term instanceof NumericLiteral) {
                return ((NumericLiteral) term).getNumber().toString();
            } else if(term instanceof StringLiteral) {
                return "\"" + ((StringLiteral) term).getValue() + "\"";
            } else if(term instanceof VariableLiteral) {
                VariableLiteral variableLiteral = (VariableLiteral) term;
                // If this literal references an Enum type, then we have to output it differently.
                if(getTypeDefinitions().get(variableLiteral.getName()) instanceof EnumTypeDefinition) {
                    return variableLiteral.getName() + "_" + variableLiteral.getChild().getName() +
                        ((variableLiteral.getChild().getChild() != null) ?
                            "." + toVariableExpressionRest(variableLiteral.getChild().getChild()) : "");
                } else {
                    return variableExpressionGenerator.apply(term);
                }
            } else {
                throw new RuntimeException("Unsupported Literal type " + term.getClass().getName());
            }
        } else if (term instanceof UnaryTerm) {
            UnaryTerm ut = (UnaryTerm) term;
            Term a = ut.getA();
            switch(ut.getOperation()) {
                case "!":
                    return "!(" + toExpression(field, a, variableExpressionGenerator) + ")";
                case "-":
                    return "-(" + toExpression(field, a, variableExpressionGenerator) + ")";
                case "()":
                    return "(" + toExpression(field, a, variableExpressionGenerator) + ")";
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
                    return "Math.pow((" + toExpression(field, a, variableExpressionGenerator) + "), (" + toExpression(field, b, variableExpressionGenerator) + "))";
                default:
                    return "(" + toExpression(field, a, variableExpressionGenerator) + ") " + operation + " (" + toExpression(field, b, variableExpressionGenerator) + ")";
            }
        } else if (term instanceof TernaryTerm) {
            TernaryTerm tt = (TernaryTerm) term;
            if("if".equals(tt.getOperation())) {
                Term a = tt.getA();
                Term b = tt.getB();
                Term c = tt.getC();
                return "spi.InlineIf((" +  toExpression(field, a, variableExpressionGenerator) + "), uint16(" + toExpression(field, b, variableExpressionGenerator) + "), uint16(" + toExpression(field, c, variableExpressionGenerator) + "))";
            } else {
                throw new RuntimeException("Unsupported ternary operation type " + tt.getOperation());
            }
        } else {
            throw new RuntimeException("Unsupported Term type " + term.getClass().getName());
        }
    }

    private String toVariableParseExpression(TypedField field, Term term, Argument[] parserArguments) {
        VariableLiteral vl = (VariableLiteral) term;
        // CAST expressions are special as we need to add a ".class" to the second parameter in Java.
        if("CAST".equals(vl.getName())) {
            if((vl.getArgs() == null) || (vl.getArgs().size() != 2)) {
                throw new RuntimeException("A CAST expression expects exactly two arguments.");
            }
            VariableLiteral type = (VariableLiteral) vl.getArgs().get(1);
            StringBuilder sb = new StringBuilder(type.getName());
            sb.append("(").append(toVariableParseExpression(field, vl.getArgs().get(0), parserArguments)).append(")");
            return sb.toString() + ((vl.getChild() != null) ? "." + toVariableExpressionRest(vl.getChild()) : "");
        }
        else if("STATIC_CALL".equals(vl.getName())) {
            StringBuilder sb = new StringBuilder();
            if (!(vl.getArgs().get(0) instanceof StringLiteral)) {
                throw new RuntimeException("Expecting the first argument of a 'STATIC_CALL' to be a StringLiteral");
            }
            // Get the class and method name
            String methodName = ((StringLiteral) vl.getArgs().get(0)).getValue();
            // Cut off the double-quptes
            methodName = methodName.substring(1, methodName.length() - 1);
            sb.append(methodName).append("(");
            for (int i = 1; i < vl.getArgs().size(); i++) {
                Term arg = vl.getArgs().get(i);
                if (i > 1) {
                    sb.append(", ");
                }
                if (arg instanceof VariableLiteral) {
                    VariableLiteral va = (VariableLiteral) arg;
                    // "io" is the default name of the reader argument which is always available.
                    boolean isParserArg = "io".equals(va.getName());
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
                        sb.append(va.getName() + ((va.getChild() != null) ? "." + toVariableExpressionRest(va.getChild()) : ""));
                    }
                    // We have to manually evaluate the type information at code-generation time.
                    else if (isTypeArg) {
                        String part = va.getChild().getName();
                        switch (part) {
                            case "name":
                                sb.append("\"").append(field.getTypeName()).append("\"");
                                break;
                            case "length":
                                sb.append("\"").append(((SimpleTypeReference) field).getSizeInBits()).append("\"");
                                break;
                            case "encoding":
                                String encoding = ((StringTypeReference) field.getType()).getEncoding();
                                // Cut off the single quotes.
                                encoding = encoding.substring(1, encoding.length() - 1);
                                sb.append("\"").append(encoding).append("\"");
                                break;
                        }
                    } else {
                        sb.append(toVariableParseExpression(field, va, null));
                    }
                } else if (arg instanceof StringLiteral) {
                    sb.append(((StringLiteral) arg).getValue());
                }
            }
            sb.append(")");
            return sb.toString();
        }
        else if("COUNT".equals(vl.getName())) {
            return "uint8(len(" + vl.getName() + ((vl.getChild() != null) ? "." + toVariableExpressionRest(vl.getChild()) : "") + "))";
        }
        // All uppercase names are not fields, but utility methods.
        else if(vl.getName().equals(vl.getName().toUpperCase())) {
            StringBuilder sb = new StringBuilder(vl.getName());
            if(vl.getArgs() != null) {
                sb.append("(");
                boolean firstArg = true;
                for(Term arg : vl.getArgs()) {
                    if(!firstArg) {
                        sb.append(", ");
                    }
                    sb.append(toParseExpression(field, arg, parserArguments));
                    firstArg = false;
                }
                sb.append(")");
            }
            if(vl.getIndex() != VariableLiteral.NO_INDEX) {
                sb.append("[").append(vl.getIndex()).append("]");
            }
            return sb.toString() + ((vl.getChild() != null) ? "." + toVariableExpressionRest(vl.getChild()) : "");
        }
        return vl.getName() + ((vl.getChild() != null) ? "." + toVariableExpressionRest(vl.getChild()) : "");
    }

    private String toVariableSerializationExpression(TypedField field, Term term, Argument[] serialzerArguments) {
        VariableLiteral vl = (VariableLiteral) term;
        if("STATIC_CALL".equals(vl.getName())) {
            StringBuilder sb = new StringBuilder();
            if(!(vl.getArgs().get(0) instanceof StringLiteral)) {
                throw new RuntimeException("Expecting the first argument of a 'STATIC_CALL' to be a StringLiteral");
            }
            String methodName = ((StringLiteral) vl.getArgs().get(0)).getValue();
            methodName = methodName.substring(1, methodName.length() - 1);
            sb.append(methodName).append("(");
            for(int i = 1; i < vl.getArgs().size(); i++) {
                Term arg = vl.getArgs().get(i);
                if(i > 1) {
                    sb.append(", ");
                }
                if(arg instanceof VariableLiteral) {
                    VariableLiteral va = (VariableLiteral) arg;
                    // "io" and "m" are always available in every parser.
                    boolean isSerializerArg = "io".equals(va.getName()) || "m".equals(va.getName()) || "element".equals(va.getName());
                    boolean isTypeArg = "_type".equals(va.getName());
                    if(!isSerializerArg && !isTypeArg && serialzerArguments != null) {
                        for (Argument serializerArgument : serialzerArguments) {
                            if (serializerArgument.getName().equals(va.getName())) {
                                isSerializerArg = true;
                                break;
                            }
                        }
                    }
                    if(isSerializerArg) {
                        sb.append(va.getName() + ((va.getChild() != null) ? "." + toVariableExpressionRest(va.getChild()) : ""));
                    } else if(isTypeArg) {
                        String part = va.getChild().getName();
                        switch (part) {
                            case "name":
                                sb.append("\"").append(field.getTypeName()).append("\"");
                                break;
                            case "length":
                                sb.append("\"").append(((SimpleTypeReference) field).getSizeInBits()).append("\"");
                                break;
                            case "encoding":
                                String encoding = ((StringTypeReference) field.getType()).getEncoding();
                                // Cut off the single quotes.
                                encoding = encoding.substring(1, encoding.length() - 1);
                                sb.append("\"").append(encoding).append("\"");
                                break;
                        }
                    } else {
                        sb.append(toVariableSerializationExpression(field, va, null));
                    }
                } else if(arg instanceof StringLiteral) {
                    sb.append(((StringLiteral) arg).getValue());
                }
            }
            sb.append(")");
            return sb.toString();
        }
        else if("COUNT".equals(vl.getName())) {
            VariableLiteral va = (VariableLiteral) vl.getArgs().get(0);
            // "io" and "m" are always available in every parser.
            boolean isSerializerArg = "io".equals(va.getName()) || "m".equals(va.getName()) || "element".equals(va.getName());
            if(!isSerializerArg && serialzerArguments != null) {
                for (Argument serializerArgument : serialzerArguments) {
                    if (serializerArgument.getName().equals(va.getName())) {
                        isSerializerArg = true;
                        break;
                    }
                }
            }
            StringBuilder sb = new StringBuilder();
            if(isSerializerArg) {
                sb.append(va.getName() + ((va.getChild() != null) ? "." + toVariableExpressionRest(va.getChild()) : ""));
            } else {
                sb.append(toVariableSerializationExpression(field, va, null));
            }
            return getLanguageTypeNameForField(field) + "(len(" + sb.toString() + "))";
        }
        else if("ARRAY_SIZE_IN_BYTES".equals(vl.getName())) {
            VariableLiteral va = (VariableLiteral) vl.getArgs().get(0);
            // "io" and "m" are always available in every parser.
            boolean isSerializerArg = "io".equals(va.getName()) || "m".equals(va.getName()) || "element".equals(va.getName());
            if(!isSerializerArg && serialzerArguments != null) {
                for (Argument serializerArgument : serialzerArguments) {
                    if (serializerArgument.getName().equals(va.getName())) {
                        isSerializerArg = true;
                        break;
                    }
                }
            }
            StringBuilder sb = new StringBuilder();
            if(isSerializerArg) {
                sb.append(va.getName()).append(((va.getChild() != null) ? "." + toVariableExpressionRest(va.getChild()) : ""));
            } else {
                sb.append(toVariableSerializationExpression(field, va, null));
            }
            return getLanguageTypeNameForField(field) + "(" + ((VariableLiteral) vl.getArgs().get(0)).getName() + "ArraySizeInBytes(" + sb.toString() + "))";
        }
        // The synthetic checksumRawData is a local field and should not be accessed as bean property.
        boolean isSerializerArg = "checksumRawData".equals(vl.getName()) || "_value".equals(vl.getName()) || "element".equals(vl.getName());
        boolean isTypeArg = "_type".equals(vl.getName());
        if(!isSerializerArg && !isTypeArg && serialzerArguments != null) {
            for (Argument serializerArgument : serialzerArguments) {
                if (serializerArgument.getName().equals(vl.getName())) {
                    isSerializerArg = true;
                    break;
                }
            }
        }
        if(isSerializerArg) {
            return vl.getName() + ((vl.getChild() != null) ? "." + toVariableExpressionRest(vl.getChild()) : "");
        } else if(isTypeArg) {
            String part = vl.getChild().getName();
            switch (part) {
                case "name":
                    return"\"" + field.getTypeName() + "\"";
                case "length":
                    return"\"" + ((SimpleTypeReference) field).getSizeInBits() + "\"";
                case "encoding":
                    String encoding = ((StringTypeReference) field.getType()).getEncoding();
                    // Cut off the single quotes.
                    encoding = encoding.substring(1, encoding.length() - 1);
                    return"\"" + encoding + "\"";
                default:
                    return "";
            }
        } else {
            return "m." + toVariableExpressionRest(vl);
        }
    }

    private String toVariableExpressionRest(VariableLiteral vl) {
        if("lengthInBytes".equals(vl.getName())) {
            return "LengthInBytes()";
        } else if("lengthInBits".equals(vl.getName())) {
            return "LengthInBits()";
        }
        return vl.getName() + ((vl.isIndexed() ? "[" + vl.getIndex() + "]" : "") +
            ((vl.getChild() != null) ? "." + toVariableExpressionRest(vl.getChild()) : ""));
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
                        sb.append("(").append(toSerializationExpression(null, arrayField.getLoopExpression(), parserArguments)).append(" * ").append(type.getSizeInBits()).append(") + ");
                        break;
                    case LENGTH:
                        sb.append("(").append(toSerializationExpression(null, arrayField.getLoopExpression(), parserArguments)).append(" * 8) + ");
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
                    sb.append("(").append(toSerializationExpression(null, manualField.getLengthExpression(), parserArguments)).append(") + ");
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
                return "-1";
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
        List<String> imports = new ArrayList<>();
        // At least one reserved field or simple field with complex type
        if(((ComplexTypeDefinition) getThisTypeDefinition()).getFields().stream().anyMatch(field ->
            (field instanceof ReservedField))) {
            imports.add("log \"github.com/sirupsen/logrus\"");
        }

        // For "Fields with complex type", constant, typeSwitch,  fields: "errors"
        if(((ComplexTypeDefinition) getThisTypeDefinition()).getFields().stream().anyMatch(field ->
            (field instanceof ConstField) || (field instanceof SwitchField) ||
                ((field instanceof TypedField) && ((TypedField) field).getType() instanceof ComplexTypeReference))) {
            imports.add("\"errors\"");
        }

        // "Fields with complex type": "reflect"
        if(((ComplexTypeDefinition) getThisTypeDefinition()).getFields().stream().anyMatch(field ->
            !(field instanceof EnumField) &&
            ((field instanceof TypedField) && ((TypedField) field).getType() instanceof ComplexTypeReference))) {
            imports.add("\"reflect\"");
        }

        // For Constant field: "strconv"
        if(((ComplexTypeDefinition) getThisTypeDefinition()).getFields().stream().anyMatch(field ->
            (field instanceof ConstField))) {
            imports.add("\"strconv\"");
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
        for (Field curField : ((ComplexTypeDefinition) getThisTypeDefinition()).getFields()) {
            if(requiresVariable(curField, "curPos")) {
                return true;
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

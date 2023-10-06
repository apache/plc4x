<#--
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
-->
<#-- Prevent freemarker from escaping stuff -->
<#outputformat "undefined">
<#-- Declare the name and type of variables passed in to the template -->
<#-- @ftlvariable name="languageName" type="java.lang.String" -->
<#-- @ftlvariable name="protocolName" type="java.lang.String" -->
<#-- @ftlvariable name="outputFlavor" type="java.lang.String" -->
<#-- @ftlvariable name="helper" type="org.apache.plc4x.language.java.JavaLanguageTemplateHelper" -->
<#-- @ftlvariable name="tracer" type="org.apache.plc4x.plugins.codegenerator.protocol.freemarker.Tracer" -->
<#-- @ftlvariable name="type" type="org.apache.plc4x.plugins.codegenerator.types.definitions.ComplexTypeDefinition" -->
${helper.packageName(protocolName, languageName, outputFlavor)?replace(".", "/")}/${type.name}.py
#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#

<#macro importSectionWithContentBelow>
    <#local sectionContent><#nested></#local>
    <#if helper.getRequiredImports()?has_content>
        <#list helper.getRequiredImports() as import>
${import}
        </#list>
    </#if>
    ${sectionContent}
</#macro>
<#macro emitImport import>${helper.emitRequiredImport(import)}</#macro>
<@importSectionWithContentBelow>


<#-- TODO: the code below implies that parserArguments will be null if not present... not pretty  -->
<#if type.parserArguments.isPresent()><#assign parserArguments=type.parserArguments.orElseThrow()></#if>
class ${type.name}:
    <@emitImport import="from loguru import logging as log" />

    <@emitImport import="from abc import staticmethod" />
    @staticmethod
    def static_parse(read_buffer: ReadBuffer<#if type.parserArguments.isPresent()>, <#list type.parserArguments.orElseThrow() as parserArgument>${helper.camelCaseToSnakeCase(parserArgument.name)}: ${helper.getLanguageTypeNameForTypeReference(parserArgument.type, false)}<#sep>, </#sep></#list></#if>):
        <#assign defaultCaseOutput=false>
        <#assign dataIoTypeDefinition=type.asDataIoTypeDefinition().orElseThrow()>
        <#list dataIoTypeDefinition.switchField.orElseThrow().cases as case>
                <#if case.discriminatorValueTerms?has_content>
        if <@compress single_line=true>
                    <#list case.discriminatorValueTerms as discriminatorValueTerm>
                        <#assign discriminatorExpression=dataIoTypeDefinition.switchField.orElseThrow().discriminatorExpressions[discriminatorValueTerm?index].asLiteral().orElseThrow().asVariableLiteral().orElseThrow()>
                        <#assign discriminatorType=helper.getDiscriminatorTypes()[discriminatorExpression.name]>
                        EvaluationHelper.equals(
                        ${helper.camelCaseToSnakeCase(helper.toParseExpression(dataIoTypeDefinition.switchField.orElseThrow(), discriminatorType, discriminatorExpression, parserArguments))},
                        <#if discriminatorType.isEnumTypeReference()>
                        ${helper.getLanguageTypeNameForTypeReference(discriminatorType)}.${helper.toParseExpression(dataIoTypeDefinition.switchField.orElseThrow(), discriminatorType, discriminatorValueTerm, parserArguments)}
                        <#else>
                        ${helper.camelCaseToSnakeCase(helper.toParseExpression(dataIoTypeDefinition.switchField.orElseThrow(), discriminatorType, discriminatorValueTerm, parserArguments))}
                        </#if>
                    )
                        <#sep> and </#sep>
                    </#list>
                    :</@compress>  # ${case.name}
                <#else>
                    <#assign defaultCaseOutput=true>
                </#if>
            <#assign valueDefined=false>
            <#list case.fields as field>
                <#switch field.typeName>
                    <#case "array">
                        <#assign arrayField = field.asArrayField().orElseThrow()>
                        <#assign elementTypeReference=arrayField.type.elementTypeReference>
            # Array field (${helper.camelCaseToSnakeCase(arrayField.name)})
            <#-- Only update curPos if the length expression uses it -->
                        <#if arrayField.loopExpression.contains("curPos")>
            cur_pos = read_buffer.get_pos() - start_pos
                        </#if>
            <#-- If this is a count array, we can directly initialize an array with the given size -->
                        <#if field.isCountArrayField()>
            # Count array
            if ${helper.toParseExpression(arrayField, helper.intTypeReference, arrayField.loopExpression,parserArguments)} > Integer.MAX_VALUE:
                raise ParseException("Array count of " + (${helper.toParseExpression(arrayField, helper.intTypeReference, arrayField.loopExpression,parserArguments)}) + " exceeds the maximum allowed count of " + Integer.MAX_VALUE)

            item_count: int = int(${helper.toParseExpression(arrayField, helper.intTypeReference, arrayField.loopExpression,parserArguments)})
            ${helper.camelCaseToSnakeCase(arrayField.name)}: List[PlcValue] = []
            for cur_item in range(item_count):
                ${helper.camelCaseToSnakeCase(arrayField.name)}.append(${helper.getPlcValueTypeForTypeReference(elementTypeReference)}(${helper.getLanguageTypeNameForTypeReference(elementTypeReference, false)}(<#if elementTypeReference.isSimpleTypeReference()>${helper.getReadBufferReadMethodCall(elementTypeReference.asSimpleTypeReference().orElseThrow(), "", arrayField)})<#else>${elementTypeReference.asComplexTypeReference().orElseThrow().name}IO.static_parse(read_buffer<#if elementTypeReference.params.isPresent()>, <#list elementTypeReference.params.orElseThrow() as parserArgument>(${helper.getLanguageTypeNameForTypeReference(helper.getArgumentType(elementTypeReference, parserArgument?index), true)}) (${helper.toParseExpression(arrayField, elementTypeReference, parserArgument,parserArguments)})<#sep>, </#sep></#list></#if>)</#if>))

            <#-- In all other cases do we have to work with a list, that is later converted to an array -->
                        <#else>
            <#-- For a length array, we read data till the read position of the buffer reaches a given position -->
                            <#if arrayField.isLengthArrayField()>
            # Length array
            _${helper.camelCaseToSnakeCase(arrayField.name)}_length: long = ${helper.toParseExpression(arrayField, helper.intTypeReference, arrayField.loopExpression,parserArguments)}
            ${helper.camelCaseToSnakeCase(arrayField.name)}_end_pos: long = read_buffer.get_pos() + _${helper.camelCaseToSnakeCase(arrayField.name)}_length
            value: List[PlcValue]= []
            while read_buffer.get_pos() < ${helper.camelCaseToSnakeCase(arrayField.name)}_end_pos):
                value.append(<@compress single_line=true>
                <#if elementTypeReference.isSimpleTypeReference()>
                    ${helper.getPlcValueTypeForTypeReference(elementTypeReference)}(${helper.getReadBufferReadMethodCall(elementTypeReference.asSimpleTypeReference().orElseThrow(), "", arrayField)})
                <#else>${elementTypeReference.asNonSimpleTypeReference().orElseThrow().name}IO.static_parse(readBuffer
                    <#if elementTypeReference.params.isPresent()>,
                        <#list elementTypeReference.params.orElseThrow() as parserArgument>
                            ${helper.getLanguageTypeNameForTypeReference(helper.getArgumentType(elementTypeReference, parserArgument?index), true)}(${helper.toParseExpression(arrayField, elementTypeReference, parserArgument,parserArguments)})
                            <#sep>, </#sep>
                        </#list>
                    </#if>
                    )
                </#if>
                )
                </@compress>

            <#-- A terminated array keeps on reading data as long as the termination expression evaluates to false -->
                            <#elseif arrayField.isTerminatedArrayField()>
            # Terminated array
            ${arrayField.name}: ${helper.getNonPrimitiveLanguageTypeNameForField(arrayField)} = new LinkedList<>()
            while not bool(${helper.camelCaseToSnakeCase(helper.toParseExpression(arrayField, helper.boolTypeReference, arrayField.loopExpression,parserArguments))})):
                ${helper.camelCaseToSnakeCase(arrayField.name)}.append(<#if elementTypeReference.isSimpleTypeReference()>${helper.getReadBufferReadMethodCall(elementTypeReference.asSimpleTypeReference().orElseThrow(), "", arrayField)}<#else>${elementTypeReference.asComplexTypeReference().orElseThrow().name}IO.static_parse(readBuffer<#if arrayField.params.isPresent()>, <#list arrayField.params.orElseThrow() as parserArgument>(${helper.getLanguageTypeNameForTypeReference(helper.getArgumentType(elementTypeReference, parserArgument?index), true)}) (${helper.toParseExpression(arrayField, parserArgument,parserArguments)})<#sep>, </#sep></#list></#if>)</#if>)

                <#-- After parsing, update the current position, but only if it's needed -->
                                <#if arrayField.loopExpression.contains("curPos")>
                cur_pos = read_buffer.get_pos() - start_pos
                                </#if>

                            </#if>
                        </#if>
                        <#if arrayField.name == "value">
                            <#assign valueDefined=true>
                        </#if>
                    <#break>
                    <#case "const">
                        <#assign constField=field.asConstField().orElseThrow()>

            # Const Field (${helper.camelCaseToSnakeCase(constField.name)})
            ${helper.camelCaseToSnakeCase(constField.name)}: ${helper.getNonPrimitiveLanguageTypeNameForField(constField)} = ${helper.getReadBufferReadMethodCall(constField.type.asSimpleTypeReference().orElseThrow(), "", constField)}
            if ${helper.camelCaseToSnakeCase(constField.name)} != ${dataIoTypeDefinition.name}.${constField.name?upper_case}):
                raise ParseException("Expected constant value " + ${dataIoTypeDefinition.name}.${constField.name?upper_case} + " but got " + ${helper.camelCaseToSnakeCase(constField.name)})

                        <#if constField.name == "value">
                            <#assign valueDefined=true>
                        </#if>
                    <#break>
                    <#case "enum">
                        <#assign enumField=field.asEnumField().orElseThrow()>

            # Enum field (${helper.camelCaseToSnakeCase(enumField.name)})
            ${helper.camelCaseToSnakeCase(enumField.name)}: ${helper.getNonPrimitiveLanguageTypeNameForField(enumField)} = ${helper.getNonPrimitiveLanguageTypeNameForField(enumField)}.enum_for_value(${helper.getReadBufferReadMethodCall(helper.getEnumBaseTypeReference(enumField.type.asSimpleTypeReference().orElseThrow()), "", enumField)})
                        <#if enumField.name == "value">
                            <#assign valueDefined=true>
                        </#if>
                    <#break>
                    <#case "manual">
                        <#assign manualField=field.asManualField().orElseThrow()>

            # Manual Field (${helper.camelCaseToSnakeCase(manualField.name)})
            ${helper.camelCaseToSnakeCase(manualField.name)}: ${helper.getLanguageTypeNameForField(field)} = (${helper.getLanguageTypeNameForField(manualField)}) (${helper.toParseExpression(manualField, manualField.type, manualField.parseExpression,parserArguments)})
                        <#if manualField.name == "value">
                            <#assign valueDefined=true>
                        </#if>
                    <#break>
                    <#case "reserved">
                        <#assign reservedField=field.asReservedField().orElseThrow()>

            # Reserved Field (Compartmentalized so the "reserved" variable can't leak)
            reserved: ${helper.getLanguageTypeNameForField(field)} = ${helper.getReadBufferReadMethodCall(reservedField.type.asSimpleTypeReference().orElseThrow(), "", reservedField)}
            if reserved != ${helper.getReservedValue(reservedField)}:
                log.info("Expected constant value " + str(${reservedField.referenceValue}) + " but got " + reserved + " for reserved field.")
                    <#break>
                    <#case "simple">
                        <#assign simpleField=field.asSimpleField().orElseThrow()>

                        <#if helper.isEnumField(simpleField)>
            # Enum field (${simpleField.name})
            ${helper.camelCaseToSnakeCase(simpleField.name)}: ${helper.getNonPrimitiveLanguageTypeNameForField(simpleField)} = ${helper.getNonPrimitiveLanguageTypeNameForField(simpleField)}(${helper.getReadBufferReadMethodCall(helper.getEnumBaseTypeReference(simpleField.type), "", simpleField)})
                        <#else>
            # Simple Field (${simpleField.name})
            ${helper.camelCaseToSnakeCase(simpleField.name)}: ${helper.getNonPrimitiveLanguageTypeNameForField(simpleField)} = <#if simpleField.type.isSimpleTypeReference()>${helper.getReadBufferReadMethodCall(simpleField.type.asSimpleTypeReference().orElseThrow(), "", simpleField)}<#else>${simpleField.type.asComplexTypeReference().orElseThrow().name}IO.staticParse(readBuffer<#if simpleField.params.isPresent()>, <#list field.params.orElseThrow() as parserArgument>${helper.getLanguageTypeNameForTypeReference(helper.getArgumentType(simpleField.type, parserArgument?index), true)}(${helper.toParseExpression(simpleField, simpleField.type, parserArgument,parserArguments)})<#sep>, </#sep></#list></#if>)</#if>
                        </#if>
                        <#if case.name == "Struct" ||
                            ((case.name == "DATE_AND_TIME") && ((simpleField.name == "year") || (simpleField.name == "month") || (simpleField.name == "day") || (simpleField.name == "hour") || (simpleField.name == "minutes") || (simpleField.name == "seconds"))) ||
                            ((case.name == "DATE_AND_TIME") && (simpleField.name == "secondsSinceEpoch")) ||
                            ((case.name == "DATE") && ((simpleField.name == "year") || (simpleField.name == "month") || (simpleField.name == "day"))) ||
                            ((case.name == "TIME_OF_DAY") && ((simpleField.name == "hour") || (simpleField.name == "minutes") || (simpleField.name == "seconds"))) ||
                        simpleField.name == "value">
                            <#assign valueDefined=true>
                        </#if>
                    <#break>
                </#switch>
            </#list>
            <#if case.name == "Struct">

                <#-- In this case we need to wrap each field in a PlcValue that matches it's natural type -->
            _map: Dict[str, PlcValue] = {}
                <#list case.fields as field>
                    <#if field.isArrayField()>
                        <#assign field=field.asArrayField().orElseThrow()>
            _map["${field.name}"] = PlcList(${field.name})
                    <#elseif field.isPropertyField()>
                        <#assign field=field.asPropertyField().orElseThrow()>
                        <#switch helper.getLanguageTypeNameForTypeReference(field.type)>
                            <#case "Boolean">
            _map["${field.name}"] = PlcBOOL(${field.name})
                                <#break>
                            <#case "Byte">
            _map["${field.name}"] = PlcSINT(${field.name})
                                <#break>
                            <#case "Short">
            _map["${field.name}"] = PlcINT(${field.name})
                                <#break>
                            <#case "Integer">
            _map["${field.name}"] = PlcDINT(${field.name})
                                <#break>
                            <#case "Long">
            _map["${field.name}"] = PlcLINT(${field.name})
                                <#break>
                            <#case "BigInteger">
            _map["${field.name}"] = PlcBigInteger(${field.name})
                                <#break>
                            <#case "Float">
            _map["${field.name}"] = PlcREAL(${field.name})
                                <#break>
                            <#case "Double">
            _map["${field.name}"] = PlcLREAL(${field.name})
                                <#break>
                            <#case "BigDecimal">
            _map["${field.name}"] = PlcBigDecimal(${field.name})
                                <#break>
                            <#case "String">
            _map["${field.name}"] = PlcSTRING(${field.name})
                                <#break>
                            <#case "LocalTime">
            _map["${field.name}"] =  PlcTIME_OF_DAY(${field.name})
                                <#break>
                            <#case "LocalDate">
            _map["${field.name}"] =  PlcDATE(${field.name})
                                <#break>
                            <#case "LocalDateTime">
            _map["${field.name}"] = PlcDATE_AND_TIME(${field.name})
                                <#break>
                        </#switch>
                    </#if>
                </#list>
                <#assign valueDefined=true>
            </#if>

            <#if valueDefined>
                <#switch case.name>
                    <#case "TIME">
            return PlcTIME(value)
                    <#break>
                    <#case "DATE">
                    <#if helper.hasFieldsWithNames(case.fields, "year", "month", "day")>
            value: LocalDate = LocalDate(int(year), (month == 0) ? 1 : int(month), (day == 0) ? 1 : int(day))
                    </#if>
            return PlcDATE(value)
                    <#break>
                    <#case "TIME_OF_DAY">
                    <#if helper.hasFieldsWithNames(case.fields, "hour", "minutes", "seconds", "nanos")>
            value: LocalTime = LocalTime(int(hour), int(minutes), int(seconds), int(nanos))
                    <#elseif helper.hasFieldsWithNames(case.fields, "hour", "minutes", "seconds")>
            value: LocalTime = LocalTime(int(hour), int(minutes), int(seconds))
                    </#if>
            return PlcTIME_OF_DAY(value)
                    <#break>
                    <#case "DATE_AND_TIME">
                    <#if helper.hasFieldsWithNames(case.fields, "year", "month", "day", "hour", "minutes", "seconds", "nanos")>
            value: LocalDateTime = LocalDateTime(int(year), (month == 0) ? 1 : int(month), (day == 0) ? 1 : int(day), int(hour), int(minutes), int(seconds), int(nanos))
                    <#elseif helper.hasFieldsWithNames(case.fields, "year", "month", "day", "hour", "minutes", "seconds")>
            value: LocalDateTime = LocalDateTime(int(year), (month == 0) ? 1 : int(month), (day == 0) ? 1 : int(day), int(hour), int(minutes), int(seconds))
                    <#elseif helper.hasFieldsWithNames(case.fields, "secondsSinceEpoch")>
            value: LocalDateTime = LocalDateTime.ofEpochSecond(secondsSinceEpoch, 0, ZoneOffset.UTC)
                    </#if>
            return PlcDATE_AND_TIME(value)
                    <#break>
                    <#case "Struct">
            return PlcStruct(_map)
                    <#break>
                    <#case "List">
            return PlcList(value)
                    <#break>
                    <#default>
            return Plc${case.name}(value)
                </#switch>
            </#if>
        <#sep></#sep></#list>
        <#if !defaultCaseOutput>
        return None
        </#if>

<#if outputFlavor != "passive">
    <@emitImport import="from abc import staticmethod" />
    @staticmethod
    def static_serialize(writeBuffer: WriteBuffer, _value: PlcValue<#if type.parserArguments.isPresent()>, <#list type.parserArguments.orElseThrow() as parserArgument>${parserArgument.name}: ${helper.getLanguageTypeNameForTypeReference(parserArgument.type, false)}<#sep>, </#sep></#list></#if>) -> None:
        static_serialize(writeBuffer, _value<#if type.parserArguments.isPresent()>, <#list type.parserArguments.orElseThrow() as parserArgument>${parserArgument.name}<#sep>, </#sep></#list></#if>, ByteOrder.BIG_ENDIAN)

    <@emitImport import="from abc import staticmethod" />
    @staticmethod
    def static_serialize(writeBuffer: WriteBuffer, _value: PlcValue<#if type.parserArguments.isPresent()>, <#list type.parserArguments.orElseThrow() as parserArgument>${parserArgument.name}: ${helper.getLanguageTypeNameForTypeReference(parserArgument.type, false)}<#sep>, </#sep></#list></#if>, byteOrder: ByteOrder) -> None:
        <#assign defaultCaseOutput=false>
        <#assign dataIoTypeDefinition=type.asDataIoTypeDefinition().orElseThrow()>
        <#list dataIoTypeDefinition.switchField.orElseThrow().cases as case>
            <#if case.discriminatorValueTerms?has_content>
        if <@compress single_line=true>
                <#list case.discriminatorValueTerms as discriminatorValueTerm>
                    <#assign discriminatorExpression=dataIoTypeDefinition.switchField.orElseThrow().discriminatorExpressions[discriminatorValueTerm?index].asLiteral().orElseThrow().asVariableLiteral().orElseThrow()>
                    <#assign discriminatorType=helper.getDiscriminatorTypes()[discriminatorExpression.name]>
                    EvaluationHelper.equals(
                    ${helper.toParseExpression(dataIoTypeDefinition.switchField.orElseThrow(), discriminatorType, discriminatorExpression, parserArguments)},
                    <#if discriminatorType.isEnumTypeReference()>
                    ${helper.getLanguageTypeNameForTypeReference(discriminatorType)}.${helper.toParseExpression(dataIoTypeDefinition.switchField.orElseThrow(), discriminatorType, discriminatorValueTerm, parserArguments)}
                    <#else>
                    ${helper.toParseExpression(dataIoTypeDefinition.switchField.orElseThrow(), discriminatorType, discriminatorValueTerm, parserArguments)}
                    </#if>)
                    <#sep> and </#sep>
                </#list>
                : # ${case.name}</@compress>
            <#else>
                <#assign defaultCaseOutput=true>
            </#if>

            <#list case.fields as field>
                <#switch field.typeName>
                    <#case "array">
                        <#assign arrayField=field.asArrayField().orElseThrow()>
                        <#assign elementTypeReference=arrayField.type.elementTypeReference>
            values: PlcList = _value

                        <#if case.name == "Struct">
            for val in values.getStruct().get("${arrayField.name}").getList():
                <#if elementTypeReference.isByteBased()>
                value: ${helper.getLanguageTypeNameForField(arrayField)} = val.getRaw()
                writeBuffer.writeByteArray("", value)
                <#else>
                value: ${helper.getLanguageTypeNameForField(arrayField)} = val.get${helper.getLanguageTypeNameForField(arrayField)?cap_first}()
                ${helper.getWriteBufferWriteMethodCall(elementTypeReference.asSimpleTypeReference().orElseThrow(), "value", arrayField)}
                </#if>
                        <#else>
            for val in values.getList():
                <#if elementTypeReference.isByteBased()>
                value: list[byte] = val.getRaw()
                writeBuffer.writeByteArray("", value)
                <#else>
                value: ${helper.getLanguageTypeNameForTypeReference(elementTypeReference)} = val.get${helper.getLanguageTypeNameForTypeReference(elementTypeReference)?cap_first}()
                ${helper.getWriteBufferWriteMethodCall(elementTypeReference.asSimpleTypeReference().orElseThrow(), "(" + arrayField.name + ")", arrayField)}
                </#if>
                        </#if>

                    <#if case.name == "BOOL">
            while writeBuffer.getPos() < len(writeBuffer.getData()):
                writeBuffer.writeBit(false)
                    </#if>
                    <#break>
                    <#case "const">
                        <#assign constField=field.asConstField().orElseThrow()>
            # Const Field (${constField.name})
            ${helper.getWriteBufferWriteMethodCall(constField.type.asSimpleTypeReference().orElseThrow(), constField.referenceValue, constField)}
                    <#break>
                    <#case "enum">
                        <#assign enumField=field.asEnumField().orElseThrow()>
            # Enum field (${enumField.name})
            ${enumField.name}: ${helper.getLanguageTypeNameForField(field)} = _value.get${enumField.name?cap_first}()
            ${helper.getWriteBufferWriteMethodCall(helper.getEnumBaseTypeReference(field.asTypedField().orElseThrow().type), "(" + enumField.name + ".getValue())", enumField)}
                    <#break>
                    <#case "manual">
                        <#assign manualField=field.asManualField().orElseThrow()>
            # Manual Field (${manualField.name})
            ${helper.toSerializationExpression(manualField, manualField.type, manualField.serializeExpression, type.parserArguments.orElse(null))}
                    <#break>
                    <#case "reserved">
                        <#assign reservedField=field.asReservedField().orElseThrow()>
            # Reserved Field
            ${helper.getWriteBufferWriteMethodCall(reservedField.type.asSimpleTypeReference().orElseThrow(), helper.getReservedValue(reservedField), reservedField)}
                    <#break>
                    <#case "simple">
                        <#assign simpleField=field.asSimpleField().orElseThrow()>
            # Simple Field (${simpleField.name})
                        <#if case.name == "Struct">
            ${simpleField.name}: ${helper.getLanguageTypeNameForField(simpleField)} = _value.getStruct().get("${simpleField.name}").get${helper.getLanguageTypeNameForField(simpleField)?cap_first}()
                        <#else>
                            <#if simpleField.name == "value">
            ${simpleField.name}: ${helper.getLanguageTypeNameForField(simpleField)} = _value.get${helper.getLanguageTypeNameForField(simpleField)?cap_first}()
                            <#elseif simpleField.name == "secondsSinceEpoch">
            ${simpleField.name}: ${helper.getLanguageTypeNameForField(simpleField)} = _value.get${helper.getLanguageTypeNameForField(simpleField)?cap_first}()
                            <#else>
                                <#-- Just for now -->
            ${simpleField.name}: ${helper.getLanguageTypeNameForField(simpleField)} = ${helper.getNullValueForTypeReference(simpleField.type)}
                            </#if>
                        </#if>
                        <#if simpleField.type.isSimpleTypeReference()>
            ${helper.getWriteBufferWriteMethodCall(simpleField.type.asSimpleTypeReference().orElseThrow(), "(" + simpleField.name + ")", simpleField)}
                        <#else>
            ${simpleField.type.asComplexTypeReference().orElseThrow().name}IO.staticSerialize(writeBuffer, ${simpleField.name})
                        </#if>
                    <#break>
                </#switch>
            </#list>
        <#sep></#sep></#list>
</#if>

    <@emitImport import="from abc import staticmethod" />
    <@emitImport import="import math" />
    @staticmethod
    def get_length_in_bytes(_value: PlcValue<#if type.parserArguments.isPresent()>, <#list type.parserArguments.orElseThrow() as parserArgument>${parserArgument.name}: ${helper.getLanguageTypeNameForTypeReference(parserArgument.type, false)}<#sep>, </#sep></#list></#if>) -> int:
        return int(math.ceil(float(getLengthInBits(_value<#if type.parserArguments.isPresent()>, <#list type.parserArguments.orElseThrow() as parserArgument>${parserArgument.name}<#sep>, </#sep></#list></#if>)) / 8.0))


    <@emitImport import="from abc import staticmethod" />
    @staticmethod
    def get_length_in_bits(_value: PlcValue<#if type.parserArguments.isPresent()>, <#list type.parserArguments.orElseThrow() as parserArgument>${parserArgument.name}: ${helper.getLanguageTypeNameForTypeReference(parserArgument.type, false)}<#sep>, </#sep></#list></#if>) -> int:
        sizeInBits: int = 0
        <#assign defaultCaseOutput=false>
        <#assign dataIoTypeDefinition=type.asDataIoTypeDefinition().orElseThrow()>
        <#list dataIoTypeDefinition.switchField.orElseThrow().cases as case>
        <#if case.discriminatorValueTerms?has_content>
        if <@compress single_line=true>
            <#list case.discriminatorValueTerms as discriminatorValueTerm>
            <#assign discriminatorExpression=dataIoTypeDefinition.switchField.orElseThrow().discriminatorExpressions[discriminatorValueTerm?index].asLiteral().orElseThrow().asVariableLiteral().orElseThrow()>
            <#assign discriminatorType=helper.getDiscriminatorTypes()[discriminatorExpression.name]>
                EvaluationHelper.equals(
                    ${helper.toParseExpression(dataIoTypeDefinition.switchField.orElseThrow(), discriminatorType, discriminatorExpression, parserArguments)},
                    <#if discriminatorType.isEnumTypeReference()>
                    ${helper.getLanguageTypeNameForTypeReference(discriminatorType)}.${helper.toParseExpression(dataIoTypeDefinition.switchField.orElseThrow(), discriminatorType, discriminatorValueTerm, parserArguments)}
                    <#else>
                    ${helper.toParseExpression(dataIoTypeDefinition.switchField.orElseThrow(), discriminatorType, discriminatorValueTerm, parserArguments)}
                    </#if>
                )
            <#sep> and </#sep>
            </#list>
        </@compress>: # ${case.name}
        <#else>
            <#assign defaultCaseOutput=true>
        </#if>
            <#list case.fields as field>
            <#switch field.typeName>
            <#case "array">
                <#assign arrayField=field.asArrayField().orElseThrow()>
                <#assign elementTypeReference=arrayField.type.elementTypeReference>
                values: PlcList = _value
                <#if case.name == "Struct">
                # TODO: Finish this!
                <#elseif elementTypeReference.isComplexTypeReference()>
                # TODO: Finish this!
                <#else>
                sizeInBits += values.getList().size() * ${elementTypeReference.asSimpleTypeReference().orElseThrow().sizeInBits}
                </#if>
                <#break>
            <#case "const">
                <#assign constField=field.asConstField().orElseThrow()>
                # Const Field (${constField.name})
                <#-- const fields are only simple type -->
                sizeInBits += ${constField.type.asSimpleTypeReference().orElseThrow().sizeInBits}
                <#break>
            <#case "enum">
                <#assign enumField=field.asEnumField().orElseThrow()>
                # Enum field (${enumField.name})
                sizeInBits += ${helper.getEnumFieldSimpleTypeReference(enumField.type, enumField.fieldName).sizeInBits}
                <#break>
            <#case "manual">
                <#assign manualField=field.asManualField().orElseThrow()>
                # Manual Field (${manualField.name})
                sizeInBits += ${helper.toSerializationExpression(manualField, helper.intTypeReference, manualField.lengthExpression, type.parserArguments.orElse(null))}
                <#break>
            <#case "reserved">
                <#assign reservedField=field.asReservedField().orElseThrow()>
                # Reserved Field
                sizeInBits += ${reservedField.type.asSimpleTypeReference().orElseThrow().sizeInBits}
                <#break>
            <#case "simple">
                <#assign simpleField=field.asSimpleField().orElseThrow()>
                # Simple Field (${simpleField.name})
                sizeInBits += ${simpleField.type.asSimpleTypeReference().orElseThrow().sizeInBits}
                <#break>
            </#switch>
            </#list>
        <#sep></#sep></#list>
        return sizeInBits


</@importSectionWithContentBelow>

</#outputformat>

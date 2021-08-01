grammar MSpec;
/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

file
 : complexTypeDefinition* EOF
 ;

complexTypeDefinition
 : LBRACKET complexType RBRACKET
 ;

complexType
 : 'type' name=idExpression (LBRACKET params=argumentList RBRACKET)? fieldDefinition*
 | 'discriminatedType' name=idExpression (LBRACKET params=argumentList RBRACKET)? fieldDefinition+
 | 'enum' (type=typeReference)? name=idExpression (LBRACKET params=argumentList RBRACKET)? enumValues=enumValueDefinition+
 | 'dataIo' name=idExpression (LBRACKET params=argumentList RBRACKET)? dataIoTypeSwitch=dataIoDefinition
 ;

fieldDefinition
 : LBRACKET field (LBRACKET params=multipleExpressions RBRACKET)? RBRACKET
 ;

dataIoDefinition
 : LBRACKET typeSwitchField (LBRACKET params=multipleExpressions RBRACKET)? RBRACKET
 ;

field
 : abstractField
 | arrayField
 | checksumField
 | constField
 | discriminatorField
 | enumField
 | implicitField
 | manualArrayField
 | manualField
 | optionalField
 | paddingField
 | reservedField
 | simpleField
 | typeSwitchField
 | unknownField
 | virtualField
 ;

abstractField
 : 'abstract' type=typeReference name=idExpression
 ;

arrayField
 : 'array' type=typeReference name=idExpression loopType=ARRAY_LOOP_TYPE loopExpression=expression
 ;

checksumField
 : 'checksum' type=dataType name=idExpression checksumExpression=expression
 ;

constField
 : 'const' type=dataType name=idExpression expected=expression
 ;

discriminatorField
 : 'discriminator' type=typeReference name=idExpression
 ;

enumField
 : 'enum' type=typeReference name=idExpression (fieldName=idExpression)?
 ;

implicitField
 : 'implicit' type=dataType name=idExpression serializeExpression=expression
 ;

manualArrayField
 : 'manualArray' type=typeReference name=idExpression loopType=ARRAY_LOOP_TYPE loopExpression=expression parseExpression=expression serializeExpression=expression lengthExpression=expression
 ;

manualField
 : 'manual' type=typeReference name=idExpression parseExpression=expression serializeExpression=expression lengthExpression=expression
 ;

optionalField
 : 'optional' type=typeReference name=idExpression condition=expression
 ;

paddingField
 : 'padding' type=dataType name=idExpression paddingValue=expression paddingCondition=expression
 ;

reservedField
 : 'reserved' type=dataType expected=expression
 ;

simpleField
 : 'simple' type=typeReference name=idExpression
 ;

typeSwitchField
 : 'typeSwitch' discriminators=multipleExpressions caseStatement*
 ;

unknownField
 : 'unknown' type=dataType
 ;

virtualField
 : 'virtual' type=typeReference name=idExpression valueExpression=expression
 ;

enumValueDefinition
 : LBRACKET (valueExpression=expression)? name=IDENTIFIER_LITERAL (LBRACKET constantValueExpressions=multipleExpressions RBRACKET)? RBRACKET
 ;

typeReference
 : complexTypeReference=IDENTIFIER_LITERAL
 | simpleTypeReference=dataType
 ;

caseStatement
 : LBRACKET (discriminatorValues=multipleExpressions)? name=IDENTIFIER_LITERAL (LBRACKET params=argumentList RBRACKET)? fieldDefinition* RBRACKET
 ;

dataType
 : base='bit'
 | base='byte'
 | base='int' size=INTEGER_LITERAL
 | base='uint' size=INTEGER_LITERAL
 | base='float' exponent=INTEGER_LITERAL '.' mantissa=INTEGER_LITERAL
 | base='ufloat' exponent=INTEGER_LITERAL '.' mantissa=INTEGER_LITERAL
 | base='string' length=expression (encoding=idExpression)?
 | base='time'
 | base='date'
 | base='dateTime'
 ;

argument
 : type=typeReference name=idExpression
 ;

argumentList
 : argument (',' argument)*
 ;

expression
 : TICK expr=innerExpression TICK
 ;

multipleExpressions
 : expression (',' expression)*
 ;

innerExpression
 : BOOLEAN_LITERAL
 // Explicitly allow the loop type keywords in expressions
 | ARRAY_LOOP_TYPE
 | IDENTIFIER_LITERAL ('(' (innerExpression (',' innerExpression)* )? ')')? ('[' innerExpression ']')?
 | innerExpression '.' innerExpression // Field Reference or method call
 | innerExpression '[' + INTEGER_LITERAL + ']' // Array index
 | innerExpression BinaryOperator innerExpression  // Addition
 | innerExpression '?' innerExpression ':' innerExpression
 | '(' innerExpression ')'
 | '"' innerExpression '"'
 | '!' innerExpression
 | HEX_LITERAL
 | INTEGER_LITERAL
 | STRING_LITERAL
 ;

idExpression
 : TICK id=IDENTIFIER_LITERAL TICK
 // Explicitly allow the loop type keywords in id-expressions
 | TICK id=ARRAY_LOOP_TYPE TICK
 ;

TICK : '\'';
LBRACKET : '[';
RBRACKET : ']';
LCBRACKET : '{';
RCBRACKET : '}';

BinaryOperator
 : '+'
 | '-'
 | '/'
 | '*'
 | '^'
 | '=='
 | '!='
 | '>>'
 | '<<'
 | '>'
 | '<'
 | '>='
 | '<='
 | '&&'
 | '||'
 | '&'
 | '|'
 | '%'
 ;

ARRAY_LOOP_TYPE
 : 'count'
 | 'length'
 | 'terminated'
 ;

// Integer literals

INTEGER_LITERAL
 : INTEGER_CHARACTERS
 ;

fragment
INTEGER_CHARACTERS
 : INTEGER_CHARACTER+
 ;

fragment
INTEGER_CHARACTER
 : [0-9]
 ;

// Hexadecimal literals

HEX_LITERAL
 : '0' [xX] HEX_CHARACTERS
 ;

fragment
HEX_CHARACTERS
 : HEX_CHARACTER+
 ;

fragment
HEX_CHARACTER
 : [0-9a-fA-F]
 ;

// Boolean literals

BOOLEAN_LITERAL
 : 'true'
 | 'false'
 ;

// String literals

STRING_LITERAL
 : '"' STRING_CHARACTERS? '"'
 ;

// As we're generating property names and class names from these,
// we have to put more restrictions on them.

IDENTIFIER_LITERAL
 : [A-Za-z0-9_-]+
 ;

fragment
STRING_CHARACTERS
 : STRING_CHARACTER+
 ;

fragment
STRING_CHARACTER
 : ~["\\\r\n]
 ;

// Stuff we just want to ignore

LINE_COMMENT
 : '//' ~[\r\n]* -> skip
 ;

BLOCK_COMMENT
 : '/*' .*? '*/' -> skip
 ;

WS
 : [ \t\r\n\u000C]+ -> skip
 ;

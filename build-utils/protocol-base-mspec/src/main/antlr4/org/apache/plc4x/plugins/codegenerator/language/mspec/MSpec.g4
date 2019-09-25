grammar MSpec;
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

file
 : complexTypeDefinition* EOF
 ;

complexTypeDefinition
 : COMMENT
 | LBRACKET complexType RBRACKET
 ;

complexType
 : 'type' name=idExpression (LBRACKET params=argumentList RBRACKET)? fieldDefinition+
 | 'discriminatedType' name=idExpression (LBRACKET params=argumentList RBRACKET)? fieldDefinition+
 | 'enum' type=typeReference name=idExpression (LBRACKET params=argumentList RBRACKET)? enumValues=enumValueDefinition+
 ;


fieldDefinition
 : LBRACKET field (LBRACKET params=multipleExpressions RBRACKET)? RBRACKET
 ;

field
 : arrayField
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
 | virtualField
 ;

arrayField
 : 'array' type=typeReference name=idExpression loopType=arrayType loopExpression=expression
 ;

checksumField
 : 'checksum' type=dataType name=idExpression checksumExpression=expression
 ;

constField
 : 'const' type=dataType name=idExpression expected=expression
 ;

discriminatorField
 : 'discriminator' type=dataType name=idExpression
 ;

enumField
 : 'enum' type=typeReference name=idExpression
 ;

implicitField
 : 'implicit' type=dataType name=idExpression serializationExpression=expression
 ;

manualArrayField
 : 'manualArray' type=typeReference name=idExpression loopType=arrayType loopExpression=expression deserializationExpression=expression serializationExpression=expression lengthExpression=expression
 ;

manualField
 : 'manual' type=typeReference name=idExpression deserializationExpression=expression serializationExpression=expression lengthExpression=expression
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

virtualField
 : 'virtual' type=typeReference name=idExpression valueExpression=expression
 ;

enumValueDefinition
 : LBRACKET valueExpression=expression name=IDENTIFIER (LBRACKET constantValueExpressions=multipleExpressions RBRACKET)? RBRACKET
 ;


typeReference
 : complexTypeReference=IDENTIFIER
 | simpleTypeReference=dataType
 ;

caseStatement
 : LBRACKET (discriminatorValues=multipleExpressions)? name=IDENTIFIER (LBRACKET params=argumentList RBRACKET)? fieldDefinition* RBRACKET
 ;

dataType
 : base='bit'
 | base='int' size=INTEGER_LITERAL
 | base='uint' size=INTEGER_LITERAL
 | base='float' size=INTEGER_LITERAL
 | base='string'
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
 : 'A' | 'B' | 'C' | 'D' | 'E' | 'F'
 | HEX_LITERAL
 | INTEGER_LITERAL
 | (IDENTIFIER | arrayType) ('(' (innerExpression (',' innerExpression)* )? ')')? ('[' innerExpression ']')?
 | innerExpression '.' innerExpression // Field Reference or method call
 | innerExpression '[' + INTEGER_LITERAL + ']' // Array index
 | innerExpression BinaryOperator innerExpression  // Addition
 | innerExpression '?' innerExpression ':' innerExpression
 | '(' innerExpression ')'
 | '"' innerExpression '"'
 | '!' innerExpression
 ;

COMMENT
 : K_COMMENT [a-zA-Z0-9,.'":;()/ \t\r\n\u000C-]*
 | '//' [a-zA-Z0-9,.'":;()/ \t-]*
 ;

INTEGER_LITERAL
 : [0-9]+
 ;

HEX_LITERAL
 : HexNumeral
 ;

fragment HexNumeral
 : '0' [xX] HexDigit HexDigit?;

fragment HexDigit
 : [0-9a-fA-F]
;

arrayType
 : 'count'
 | 'length'
 | 'terminated'
 ;

idExpression
 : TICK id=IDENTIFIER TICK
 ;

fragment K_COMMENT : '<--';

TICK : '\'';
TIMES : 'x';
LBRACKET : '[';
RBRACKET : ']';
LCBRACKET : '{';
RCBRACKET : '}';

BinaryOperator
 : '+'
 | '-'
 | '/'
 | '*'
 | '=='
 | '!='
 | '>'
 | '<'
 | '>='
 | '<='
 | '&&'
 | '||'
 | '%'
 ;

ZERO : '0';
HEX_VALUE : [0-9A-F];

IDENTIFIER
 : [A-Za-z0-9_-]+
 ;

WS  :  [ \t\r\n\u000C]+ -> skip
;
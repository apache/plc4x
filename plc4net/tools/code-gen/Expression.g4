// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
// Based on https://github.com/bkiers/tiny-language-antlr4 (Unlicense).
//

grammar Expression;

expressionString
 : expression EOF
 ;

expression
 : '-' expression                                       #unaryMinusExpression
 | '!' expression                                       #notExpression
 | <assoc=right> expression '^' expression              #powerExpression
 | expression op=( '*' | '/' | '%' ) expression         #multExpression
 | expression op=( '+' | '-' ) expression               #addExpression
 | expression op=( '>>' | '<<' ) expression             #bitShiftExpression
 | expression op=( '>=' | '<=' | '>' | '<' ) expression #compExpression
 | expression op=( '==' | '!=' ) expression             #eqExpression
 | expression '&&' expression                           #andExpression
 | expression '&' expression                            #bitAndExpression
 | expression '||' expression                           #orExpression
 | expression '|' expression                            #bitOrExpression
 | expression '?' expression ':' expression             #ifExpression
 | Number                                               #numberExpression
 | HexExpression                                        #hexExpression
 | Bool                                                 #boolExpression
 | Null                                                 #nullExpression
 | identifierSegment                                    #identifierExpression
 | String indexes?                                      #stringExpression
 | '(' expression ')' indexes?                          #expressionExpression
 ;

identifierSegment
 : name=Identifier args=identifierSegmentArguments? index=identifierSegmentIndexes? ('.' rest=identifierSegmentRest)?
 ;

identifierSegmentArguments
 : arguments
 ;

identifierSegmentIndexes
 : indexes
 ;

identifierSegmentRest
 : identifierSegment
 ;

arguments
 : '(' (expression (',' expression)*)? ')'
 ;

indexes
 : ( '[' expression ']' )+
 ;

Null     : 'null';

Bool
 : 'true'
 | 'false'
 ;

Number
 : Int ( '.' Digit* )?
 ;

Identifier
 : [a-zA-Z_] [a-zA-Z_0-9]*
 ;

String
 : ["] ( ~["\r\n\\] | '\\' ~[\r\n] )* ["]
 | ['] ( ~['\r\n\\] | '\\' ~[\r\n] )* [']
 ;

Space
 : [ \t\r\n\u000C] -> skip
 ;

fragment Int
 : [1-9] Digit*
 | '0'
 ;

fragment Digit
 : [0-9]
 ;

HexExpression
  : '0' [xX] HexCharacters
  ;

fragment HexCharacters
  : HexCharacter+
  ;

fragment HexCharacter
  : [0-9a-fA-F]
  ;
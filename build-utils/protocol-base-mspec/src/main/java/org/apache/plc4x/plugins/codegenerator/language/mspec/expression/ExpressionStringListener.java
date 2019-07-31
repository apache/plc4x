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

package org.apache.plc4x.plugins.codegenerator.language.mspec.expression;

import org.apache.plc4x.plugins.codegenerator.types.terms.*;

import java.util.LinkedList;
import java.util.List;
import java.util.Stack;

public class ExpressionStringListener extends ExpressionBaseListener {

    private Stack<List<Term>> parserContexts;

    private Term root;

    public Term getRoot() {
        return root;
    }

    @Override
    public void enterExpressionString(ExpressionParser.ExpressionStringContext ctx) {
        parserContexts = new Stack<>();
        parserContexts.push(new LinkedList<>());
    }

    @Override
    public void exitExpressionString(ExpressionParser.ExpressionStringContext ctx) {
        List<Term> roots = parserContexts.pop();
        if(roots.isEmpty()) {
            throw new RuntimeException("Empty Expression not supported.");
        }
        if(roots.size() != 1) {
            throw new RuntimeException("Expression can only contain one root term.");
        }
        root = roots.get(0);
    }

    /////////////////////////////////////////////////////////////////////////////////////////
    // Literals
    /////////////////////////////////////////////////////////////////////////////////////////

    @Override
    public void exitNullExpression(ExpressionParser.NullExpressionContext ctx) {
        parserContexts.peek().add(new NullLiteral());
    }

    @Override
    public void exitBoolExpression(ExpressionParser.BoolExpressionContext ctx) {
        parserContexts.peek().add(new BooleanLiteral(Boolean.valueOf(ctx.getText())));
    }

    @Override
    public void exitNumberExpression(ExpressionParser.NumberExpressionContext ctx) {
        String strValue = ctx.Number().getText();
        if(strValue.contains(".")) {
            parserContexts.peek().add(new NumericLiteral(Double.valueOf(strValue)));
        } else {
            parserContexts.peek().add(new NumericLiteral(Long.valueOf(strValue)));
        }
    }

    @Override
    public void exitStringExpression(ExpressionParser.StringExpressionContext ctx) {
        parserContexts.peek().add(new StringLiteral(ctx.getText()));
    }

    @Override
    public void enterIdentifierExpression(ExpressionParser.IdentifierExpressionContext ctx) {
        parserContexts.push(new LinkedList<>());
    }

    @Override
    public void exitIdentifierExpression(ExpressionParser.IdentifierExpressionContext ctx) {
        List<Term> args = parserContexts.pop();
        parserContexts.peek().add(getVariableLiteral(ctx.identifierSegment(), args));
    }

    /////////////////////////////////////////////////////////////////////////////////////////
    // Unary Terms
    /////////////////////////////////////////////////////////////////////////////////////////

    @Override
    public void enterNotExpression(ExpressionParser.NotExpressionContext ctx) {
        parserContexts.push(new LinkedList<>());
    }

    @Override
    public void exitNotExpression(ExpressionParser.NotExpressionContext ctx) {
        UnaryTerm ut = getUnaryTerm("!", parserContexts.pop());
        parserContexts.peek().add(ut);
    }

    @Override
    public void enterUnaryMinusExpression(ExpressionParser.UnaryMinusExpressionContext ctx) {
        parserContexts.push(new LinkedList<>());
    }

    @Override
    public void exitUnaryMinusExpression(ExpressionParser.UnaryMinusExpressionContext ctx) {
        UnaryTerm ut = getUnaryTerm("-", parserContexts.pop());
        parserContexts.peek().add(ut);
    }

    @Override
    public void enterExpressionExpression(ExpressionParser.ExpressionExpressionContext ctx) {
        parserContexts.push(new LinkedList<>());
    }

    @Override
    public void exitExpressionExpression(ExpressionParser.ExpressionExpressionContext ctx) {
        UnaryTerm ut = getUnaryTerm("()", parserContexts.pop());
        parserContexts.peek().add(ut);
    }

    /////////////////////////////////////////////////////////////////////////////////////////
    // Binary Terms
    /////////////////////////////////////////////////////////////////////////////////////////

    @Override
    public void enterOrExpression(ExpressionParser.OrExpressionContext ctx) {
        parserContexts.push(new LinkedList<>());
    }

    @Override
    public void exitOrExpression(ExpressionParser.OrExpressionContext ctx) {
        BinaryTerm bt = getBinaryTerm("||", parserContexts.pop());
        parserContexts.peek().add(bt);
    }

    @Override
    public void enterPowerExpression(ExpressionParser.PowerExpressionContext ctx) {
        parserContexts.push(new LinkedList<>());
    }

    @Override
    public void exitPowerExpression(ExpressionParser.PowerExpressionContext ctx) {
        BinaryTerm bt = getBinaryTerm("^", parserContexts.pop());
        parserContexts.peek().add(bt);
    }

    @Override
    public void enterEqExpression(ExpressionParser.EqExpressionContext ctx) {
        parserContexts.push(new LinkedList<>());
    }

    @Override
    public void exitEqExpression(ExpressionParser.EqExpressionContext ctx) {
        BinaryTerm bt = getBinaryTerm(ctx.op.getText(), parserContexts.pop());
        parserContexts.peek().add(bt);
    }

    @Override
    public void enterAndExpression(ExpressionParser.AndExpressionContext ctx) {
        parserContexts.push(new LinkedList<>());
    }

    @Override
    public void exitAndExpression(ExpressionParser.AndExpressionContext ctx) {
        BinaryTerm bt = getBinaryTerm("&&", parserContexts.pop());
        parserContexts.peek().add(bt);
    }

    @Override
    public void enterAddExpression(ExpressionParser.AddExpressionContext ctx) {
        parserContexts.push(new LinkedList<>());
    }

    @Override
    public void exitAddExpression(ExpressionParser.AddExpressionContext ctx) {
        BinaryTerm bt = getBinaryTerm(ctx.op.getText(), parserContexts.pop());
        parserContexts.peek().add(bt);
    }

    @Override
    public void enterCompExpression(ExpressionParser.CompExpressionContext ctx) {
        parserContexts.push(new LinkedList<>());
    }

    @Override
    public void exitCompExpression(ExpressionParser.CompExpressionContext ctx) {
        BinaryTerm bt = getBinaryTerm(ctx.op.getText(), parserContexts.pop());
        parserContexts.peek().add(bt);
    }

    @Override
    public void enterMultExpression(ExpressionParser.MultExpressionContext ctx) {
        parserContexts.push(new LinkedList<>());
    }

    @Override
    public void exitMultExpression(ExpressionParser.MultExpressionContext ctx) {
        BinaryTerm bt = getBinaryTerm(ctx.op.getText(), parserContexts.pop());
        parserContexts.peek().add(bt);
    }

    /////////////////////////////////////////////////////////////////////////////////////////
    // Ternary Terms
    /////////////////////////////////////////////////////////////////////////////////////////

    @Override
    public void enterIfExpression(ExpressionParser.IfExpressionContext ctx) {
        parserContexts.push(new LinkedList<>());
    }

    @Override
    public void exitIfExpression(ExpressionParser.IfExpressionContext ctx) {
        TernaryTerm tt = getTernaryTerm("if", parserContexts.pop());
        parserContexts.peek().add(tt);
    }

    /////////////////////////////////////////////////////////////////////////////////////////
    // Helpers
    /////////////////////////////////////////////////////////////////////////////////////////

    private VariableLiteral getVariableLiteral(ExpressionParser.IdentifierSegmentContext ctx, List<Term> args) {
        String name = ctx.name.getText();
        int index = (ctx.index != null) ? Integer.valueOf(ctx.index.getText().substring(1, ctx.index.getText().length() - 1)) : VariableLiteral.NO_INDEX;
        VariableLiteral child = (ctx.rest != null) ? getVariableLiteral(ctx.rest, args) : null;
        return new VariableLiteral(name, args, index, child);
    }

    private UnaryTerm getUnaryTerm(String op, List<Term> terms) {
        if(terms.size() != 1) {
            throw new RuntimeException(op + " should be a unary operation");
        }
        Term a = terms.get(0);
        return new UnaryTerm(a, op);
    }

    private BinaryTerm getBinaryTerm(String op, List<Term> terms) {
        if(terms.size() != 2) {
            throw new RuntimeException(op + " should be a binary operation");
        }
        Term a = terms.get(0);
        Term b = terms.get(1);
        return new BinaryTerm(a, b, op);
    }

    private TernaryTerm getTernaryTerm(String op, List<Term> terms) {
        if(terms.size() != 3) {
            throw new RuntimeException(op + " should be a ternary operation");
        }
        Term a = terms.get(0);
        Term b = terms.get(1);
        Term c = terms.get(1);
        return new TernaryTerm(a, b, c, op);
    }

}

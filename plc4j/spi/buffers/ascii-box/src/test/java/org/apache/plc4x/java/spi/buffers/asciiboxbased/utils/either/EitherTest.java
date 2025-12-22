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
package org.apache.plc4x.java.spi.buffers.asciiboxbased.utils.either;

import org.junit.jupiter.api.Test;

import java.util.Optional;

import static org.junit.jupiter.api.Assertions.*;

class EitherTest {

    @Test
    void testLeftCreation() {
        Either<String, Integer> left = Either.left("error");

        assertTrue(left.isLeft());
        assertFalse(left.isRight());
        assertEquals("error", left.getLeftOrThrow());
    }

    @Test
    void testRightCreation() {
        Either<String, Integer> right = Either.right(42);

        assertFalse(right.isLeft());
        assertTrue(right.isRight());
        assertEquals(42, right.getRightOrThrow());
    }

    @Test
    void testLeftGetRightThrows() {
        Either<String, Integer> left = Either.left("error");

        assertThrows(IllegalStateException.class, left::getRightOrThrow);
    }

    @Test
    void testRightGetLeftThrows() {
        Either<String, Integer> right = Either.right(42);

        assertThrows(IllegalStateException.class, right::getLeftOrThrow);
    }

    @Test
    void testLeftOpt() {
        Either<String, Integer> left = Either.left("error");
        Either<String, Integer> right = Either.right(42);

        assertEquals(Optional.of("error"), left.leftOpt());
        assertEquals(Optional.empty(), right.leftOpt());
    }

    @Test
    void testRightOpt() {
        Either<String, Integer> left = Either.left("error");
        Either<String, Integer> right = Either.right(42);

        assertEquals(Optional.empty(), left.rightOpt());
        assertEquals(Optional.of(42), right.rightOpt());
    }

    @Test
    void testFold() {
        Either<String, Integer> left = Either.left("error");
        Either<String, Integer> right = Either.right(42);

        String leftResult = left.fold(l -> "Left: " + l, r -> "Right: " + r);
        String rightResult = right.fold(l -> "Left: " + l, r -> "Right: " + r);

        assertEquals("Left: error", leftResult);
        assertEquals("Right: 42", rightResult);
    }

    @Test
    void testMapOnRight() {
        Either<String, Integer> right = Either.right(42);
        Either<String, Integer> mapped = right.map(x -> x * 2);

        assertTrue(mapped.isRight());
        assertEquals(84, mapped.getRightOrThrow());
    }

    @Test
    void testMapOnLeft() {
        Either<String, Integer> left = Either.left("error");
        Either<String, Integer> mapped = left.map(x -> x * 2);

        assertTrue(mapped.isLeft());
        assertEquals("error", mapped.getLeftOrThrow());
    }

    @Test
    void testMapLeftOnLeft() {
        Either<String, Integer> left = Either.left("error");
        Either<String, Integer> mapped = left.mapLeft(s -> s.toUpperCase());

        assertTrue(mapped.isLeft());
        assertEquals("ERROR", mapped.getLeftOrThrow());
    }

    @Test
    void testMapLeftOnRight() {
        Either<String, Integer> right = Either.right(42);
        Either<String, Integer> mapped = right.mapLeft(s -> s.toUpperCase());

        assertTrue(mapped.isRight());
        assertEquals(42, mapped.getRightOrThrow());
    }

    @Test
    void testFlatMapOnRight() {
        Either<String, Integer> right = Either.right(42);
        Either<String, Integer> result = right.flatMap(x -> Either.right(x * 2));

        assertTrue(result.isRight());
        assertEquals(84, result.getRightOrThrow());
    }

    @Test
    void testFlatMapOnLeft() {
        Either<String, Integer> left = Either.left("error");
        Either<String, Integer> result = left.flatMap(x -> Either.right(x * 2));

        assertTrue(result.isLeft());
        assertEquals("error", result.getLeftOrThrow());
    }

    @Test
    void testSwap() {
        Either<String, Integer> left = Either.left("error");
        Either<String, Integer> right = Either.right(42);

        Either<Integer, String> swappedLeft = left.swap();
        Either<Integer, String> swappedRight = right.swap();

        assertTrue(swappedLeft.isRight());
        assertEquals("error", swappedLeft.getRightOrThrow());

        assertTrue(swappedRight.isLeft());
        assertEquals(42, swappedRight.getLeftOrThrow());
    }

    @Test
    void testToOptional() {
        Either<String, Integer> left = Either.left("error");
        Either<String, Integer> right = Either.right(42);

        assertEquals(Optional.empty(), left.toOptional());
        assertEquals(Optional.of(42), right.toOptional());
    }

    @Test
    void testGetOrElse() {
        Either<String, Integer> left = Either.left("error");
        Either<String, Integer> right = Either.right(42);

        assertEquals(0, left.getOrElse(0));
        assertEquals(42, right.getOrElse(0));
    }

    @Test
    void testGetOrElseGet() {
        Either<String, Integer> left = Either.left("error");
        Either<String, Integer> right = Either.right(42);

        assertEquals(5, left.getOrElseGet(s -> s.length()));
        assertEquals(42, right.getOrElseGet(s -> s.length()));
    }

    @Test
    void testLeftEquals() {
        Either<String, Integer> left1 = Either.left("error");
        Either<String, Integer> left2 = Either.left("error");
        Either<String, Integer> left3 = Either.left("different");
        Either<String, Integer> right = Either.right(42);

        assertEquals(left1, left2);
        assertNotEquals(left1, left3);
        assertNotEquals(left1, right);
        assertNotEquals(left1, null);
        assertNotEquals(left1, "error");
    }

    @Test
    void testRightEquals() {
        Either<String, Integer> right1 = Either.right(42);
        Either<String, Integer> right2 = Either.right(42);
        Either<String, Integer> right3 = Either.right(100);
        Either<String, Integer> left = Either.left("error");

        assertEquals(right1, right2);
        assertNotEquals(right1, right3);
        assertNotEquals(right1, left);
        assertNotEquals(right1, null);
        assertNotEquals(right1, 42);
    }

    @Test
    void testHashCode() {
        Either<String, Integer> left1 = Either.left("error");
        Either<String, Integer> left2 = Either.left("error");
        Either<String, Integer> right1 = Either.right(42);
        Either<String, Integer> right2 = Either.right(42);

        assertEquals(left1.hashCode(), left2.hashCode());
        assertEquals(right1.hashCode(), right2.hashCode());
    }

    @Test
    void testToString() {
        Either<String, Integer> left = Either.left("error");
        Either<String, Integer> right = Either.right(42);

        assertEquals("Left(error)", left.toString());
        assertEquals("Right(42)", right.toString());
    }

    @Test
    void testResultSuccess() {
        Either<String, Integer> result = Either.Result.success(42);

        assertTrue(result.isRight());
        assertEquals(42, result.getRightOrThrow());
    }

    @Test
    void testResultFailure() {
        Either<String, Integer> result = Either.Result.failure("error");

        assertTrue(result.isLeft());
        assertEquals("error", result.getLeftOrThrow());
    }

    @Test
    void testLeftEqualsItself() {
        Either<String, Integer> left = Either.left("error");
        assertEquals(left, left);
    }

    @Test
    void testRightEqualsItself() {
        Either<String, Integer> right = Either.right(42);
        assertEquals(right, right);
    }
}

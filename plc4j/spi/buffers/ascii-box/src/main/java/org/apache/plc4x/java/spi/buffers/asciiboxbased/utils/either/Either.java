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

import java.io.Serializable;
import java.util.Objects;
import java.util.Optional;
import java.util.function.Function;

public abstract class Either<L, R> implements Serializable {

    private Either() {
    }

    // -- Factories --------------------------------------------------------------

    public static <L, R> Either<L, R> left(L value) {
        return new Left<>(value);
    }

    public static <L, R> Either<L, R> right(R value) {
        return new Right<>(value);
    }

    // -- Basic API --------------------------------------------------------------

    public abstract boolean isLeft();

    public abstract boolean isRight();

    public abstract L getLeftOrThrow();

    public abstract R getRightOrThrow();

    public final Optional<L> leftOpt() {
        return isLeft() ? Optional.of(getLeftOrThrow()) : Optional.empty();
    }

    public final Optional<R> rightOpt() {
        return isRight() ? Optional.of(getRightOrThrow()) : Optional.empty();
    }

    /**
     * Fold: apply leftFn if Left, rightFn if Right.
     */
    public abstract <T> T fold(Function<? super L, ? extends T> leftFn,
                               Function<? super R, ? extends T> rightFn);

    /**
     * Map the Right value. Left passes through unchanged.
     */
    public abstract <R2> Either<L, R2> map(Function<? super R, ? extends R2> fn);

    /**
     * Map the Left value. Right passes through unchanged.
     */
    public abstract <L2> Either<L2, R> mapLeft(Function<? super L, ? extends L2> fn);

    /**
     * Flat-map the Right value. Left passes through unchanged.
     */
    public abstract <R2> Either<L, R2> flatMap(Function<? super R, ? extends Either<L, R2>> fn);

    /**
     * Swap Left/Right.
     */
    public abstract Either<R, L> swap();

    public final Optional<R> toOptional() {
        return rightOpt();
    }

    public final R getOrElse(R fallback) {
        return isRight() ? getRightOrThrow() : fallback;
    }

    public final R getOrElseGet(Function<? super L, ? extends R> fallback) {
        return isRight() ? getRightOrThrow() : fallback.apply(getLeftOrThrow());
    }

    // -- Implementations --------------------------------------------------------

    private static final class Left<L, R> extends Either<L, R> {
        private final L value;

        private Left(L value) {
            this.value = value;
        }

        @Override
        public boolean isLeft() {
            return true;
        }

        @Override
        public boolean isRight() {
            return false;
        }

        @Override
        public L getLeftOrThrow() {
            return value;
        }

        @Override
        public R getRightOrThrow() {
            throw new IllegalStateException("Not a Right: Left(" + value + ")");
        }

        @Override
        public <T> T fold(Function<? super L, ? extends T> leftFn,
                          Function<? super R, ? extends T> rightFn) {
            return leftFn.apply(value);
        }

        @Override
        public <R2> Either<L, R2> map(Function<? super R, ? extends R2> fn) {
            return left(value);
        }

        @Override
        public <L2> Either<L2, R> mapLeft(Function<? super L, ? extends L2> fn) {
            return left(fn.apply(value));
        }

        @Override
        public <R2> Either<L, R2> flatMap(Function<? super R, ? extends Either<L, R2>> fn) {
            return left(value);
        }

        @Override
        public Either<R, L> swap() {
            return right(value);
        }

        @Override
        public boolean equals(Object o) {
            return (this == o) || (o instanceof Left && Objects.equals(value, ((Left<?, ?>) o).value));
        }

        @Override
        public int hashCode() {
            return Objects.hashCode(value);
        }

        @Override
        public String toString() {
            return "Left(" + value + ")";
        }
    }

    private static final class Right<L, R> extends Either<L, R> {
        private final R value;

        private Right(R value) {
            this.value = value;
        }

        @Override
        public boolean isLeft() {
            return false;
        }

        @Override
        public boolean isRight() {
            return true;
        }

        @Override
        public L getLeftOrThrow() {
            throw new IllegalStateException("Not a Left: Right(" + value + ")");
        }

        @Override
        public R getRightOrThrow() {
            return value;
        }

        @Override
        public <T> T fold(Function<? super L, ? extends T> leftFn,
                          Function<? super R, ? extends T> rightFn) {
            return rightFn.apply(value);
        }

        @Override
        public <R2> Either<L, R2> map(Function<? super R, ? extends R2> fn) {
            return right(fn.apply(value));
        }

        @Override
        public <L2> Either<L2, R> mapLeft(Function<? super L, ? extends L2> fn) {
            return right(value);
        }

        @Override
        public <R2> Either<L, R2> flatMap(Function<? super R, ? extends Either<L, R2>> fn) {
            return fn.apply(value);
        }

        @Override
        public Either<R, L> swap() {
            return left(value);
        }

        @Override
        public boolean equals(Object o) {
            return (this == o) || (o instanceof Right && Objects.equals(value, ((Right<?, ?>) o).value));
        }

        @Override
        public int hashCode() {
            return Objects.hashCode(value);
        }

        @Override
        public String toString() {
            return "Right(" + value + ")";
        }
    }

    // -- Optional ergonomic aliases (Result style) ------------------------------

    public static final class Result<T, E> {
        private Result() {
        }

        public static <T, E> Either<E, T> success(T value) {
            return Either.right(value);
        }

        public static <T, E> Either<E, T> failure(E error) {
            return Either.left(error);
        }
    }
}
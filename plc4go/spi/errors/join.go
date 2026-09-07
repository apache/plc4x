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

package errors

import stderrors "errors"

// Join returns an error that wraps the given errors. Delegates to stdlib errors.Join.
func Join(errs ...error) error { return stderrors.Join(errs...) }

// ErrUnsupported indicates that a requested operation cannot be performed,
// because it is unsupported. Re-exported from the stdlib so callers only need
// to import this package. Equal (==) to stdlib errors.ErrUnsupported.
var ErrUnsupported = stderrors.ErrUnsupported

// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

// Copyright 2022 Google LLC.
// Licensed under the Apache License, Version 2.0 (the "License").

exception Error {
  // The reason of the error. This is a constant value that identifies the
  // proximate cause of the error. Error reasons are unique within a particular
  // domain of errors. This should be at most 63 characters and match a
  // regular expression of `[A-Z][A-Z0-9_]+[A-Z0-9]`, which represents
  // UPPER_SNAKE_CASE.
  1: string reason,
  
  // A developer-facing human-readable error message in English. It should
  // both explain the error and offer an actionable resolution to it.
  2: string message,

  // Additional error information that the client code can use to handle
  // the error, such as retry info or a help link.
  3: map<string, string> details
}

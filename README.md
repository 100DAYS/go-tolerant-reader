# go-tolerant-reader
Go Implementation of a Tolerant Reader as described by Martin Fowler

This repo is in miserable condition. I just picked it out of a project and hardly did any work on it.
It has no design patterns and some switch-case with multi-level nesting. At some point i might want to refactor but for now, the code is just straight-forward.

Idea is that JSON input is first Unmarshalled into map[string]interface{}.

Then the resulting structure is passed to the tolerant reader along with a target struct.

The target struct should have some tags describing where in the given source data structure a certain information is expected.
Tolerant-Reader will try to find this information and convert it to the data type in the respective struct field if possible.

A usage example can be seen in the test file.

We use this to parse json encoded messages thrown on a kafka message bus by some php applications.

# cmplint

[![Go Reference](https://pkg.go.dev/badge/fillmore-labs.com/cmplint.svg)](https://pkg.go.dev/fillmore-labs.com/cmplint)
[![Test](https://github.com/fillmore-labs/cmplint/actions/workflows/test.yml/badge.svg?branch=dev)](https://github.com/fillmore-labs/cmplint/actions/workflows/test.yml)
[![CodeQL](https://github.com/fillmore-labs/cmplint/actions/workflows/github-code-scanning/codeql/badge.svg?branch=dev)](https://github.com/fillmore-labs/cmplint/actions/workflows/github-code-scanning/codeql)
[![Coverage](https://codecov.io/gh/fillmore-labs/cmplint/branch/dev/graph/badge.svg?token=J5SNKW3NJ0)](https://codecov.io/gh/fillmore-labs/cmplint)
[![Go Report Card](https://goreportcard.com/badge/fillmore-labs.com/cmplint)](https://goreportcard.com/report/fillmore-labs.com/cmplint)
[![License](https://img.shields.io/github/license/fillmore-labs/cmplint)](https://www.apache.org/licenses/LICENSE-2.0)

`cmplint` is a Go linter that detects comparisons like `ptr == &MyStruct{}` — code that is almost always incorrect. Each
composite literal or `new()` call allocates a unique address, so these comparisons typically evaluate to `false` or have
undefined behavior when zero-sized types are involved.

## Installation

### Go

```console
go install fillmore-labs.com/cmplint@latest
```

### Homebrew

```console
brew install fillmore-labs/tap/cmplint
```

### Eget

[Eget](https://github.com/zyedidia/eget) downloads prebuilt binaries. After installing it:

```console
eget fillmore-labs/cmplint
```

## Usage

```console
cmplint ./...
```

## The Problem

Comparing pointers to newly allocated values is a source of subtle bugs in Go. Consider:

```go
import (
  "time"

  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

  d := &metav1.Duration{30 * time.Second}
  if (d == &metav1.Duration{30 * time.Second}) { // This will always be false!
    // This code never executes
  }
```

The reason: every `&T{}` or `new(T)` expression creates a fresh allocation with a
[unique address](https://go.dev/ref/spec#Composite_literals). So `ptr == &MyStruct{}` is almost always `false`,
regardless of what `ptr` points to.

For zero-sized types, the behavior is [undefined](https://go.dev/ref/spec#Address_operators):

> _"Pointers to distinct zero-size variables may or may not be equal."_

### Examples of Problematic Code

Here are examples that `cmplint` will flag:

#### Direct Pointer Comparisons

```go
import (
  "github.com/operator-framework/api/pkg/operators/v1alpha1"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Checking if a operator update strategy matches expected values
func validateUpdateStrategy(spec *v1alpha1.CatalogSourceSpec) {
  expectedTime := 30 * time.Second

  // This comparison will always be false - &metav1.Duration{} creates a unique address.
  if (spec.UpdateStrategy.Interval != &metav1.Duration{Duration: expectedTime}) {
    // ...
  }

  // Correct approach: Dereference the pointer and compare values (after a nil check).
  if spec.UpdateStrategy.Interval == nil || spec.UpdateStrategy.Interval.Duration != expectedTime {
    // ...
  }
}
```

#### Error Handling with `errors.Is`

```go
func connectToDatabase() {
  db, err := dbConnect()
  // This will always be false - &url.Error{} creates a unique address.
  if errors.Is(err, &url.Error{}) {
    log.Fatal("Cannot connect to DB")
  }

  // Correct approach:
  var urlErr *url.Error
  if errors.As(err, &urlErr) {
    log.Fatal("Error connecting to DB:", urlErr)
  }
  // ...
}

func unmarshalEvent(msg []byte) {
  var es []cloudevents.Event
  err := json.Unmarshal(msg, &es)
  // This comparison will always be false:
  if errors.Is(err, &json.UnmarshalTypeError{}) {
    //...
  }

  // Correct approach:
  var typeErr *json.UnmarshalTypeError
  if errors.As(err, &typeErr) {
    //...
  }
}
```

## Special Cases

### `errors.Is` and Similar Functions

`cmplint` suppresses diagnostics when error types implement `Unwrap` or `Is` methods, since these enable legitimate use
with [`errors.Is`](https://pkg.go.dev/errors#Is). Specifically, diagnostics are suppressed when:

- **The error type has an `Unwrap() error` method**, as `errors.Is` traverses the error tree.

<details><summary><b><code>Unwrap() error</code> tree example.</b></summary>

```go
type wrappedError struct{ Cause error }

func (e *wrappedError) Error() string { return "wrapped: " + e.Cause.Error() }
func (e *wrappedError) Unwrap() error { return e.Cause } // This suppresses the diagnostic.

  // No warning for this code:
  if errors.Is(&wrappedError{os.ErrNoDeadline}, os.ErrNoDeadline) { // Valid due to "Unwrap" method.
    // ...
  }
```

</details>

- **The error type has an `Is(error) bool` method**, as custom comparison logic is executed.

<details><summary><b>Custom <code>Is(error) bool</code> method example.</b></summary>

When the static type of an error is just the `error` interface, the analyzer cannot know its dynamic type. Therefore,
the diagnostic is also suppressed when the _target_ has an `Is(error) bool` method:

```go
type customError struct{ Code int }

func (i *customError) Error() string { return fmt.Sprintf("custom error %d", i.Code) }

func (i *customError) Is(err error) bool { // This suppresses the diagnostic.
  _, ok := err.(*customError)
  return ok
}

  err = func() error {
    return &customError{100}
  }()

  // No warning for this code:
  if errors.Is(err, &customError{200}) { // Valid due to custom "Is" method.
    // ...
  }
```

</details>

#### Rare False Positives

The applied heuristic can lead to false positives in rare cases. For example, if one error type's `Is` method is
designed to compare against a different error type, `cmplint` may flag valid code. This pattern is uncommon and
potentially confusing.

<details><summary>This workaround improves clarity and suppresses the linting error.</summary>

```go
type errorA struct{ Code int }

func (e *errorA) Error() string { return fmt.Sprintf("error a %d", e.Code) }

type errorB struct{ Code int }

func (e *errorB) Error() string { return fmt.Sprintf("error b %d", e.Code) }

func (e *errorB) Is(err error) bool {
  if err, ok := err.(*errorA); ok { // errorB knows how to check against errorA.
    return e.Code == err.Code
  }

  return false
}

  err := func() error {
    return &errorB{100}
  }()

  // This valid code gets flagged:
  if errors.Is(err, &errorA{100}) { // Flagged, but technically correct.
    // ...
  }

  // Document to clarify intent and assign to an identifier to suppress the warning:
  target := &errorA{100} // errorB's "Is" method should match.
  if errors.Is(err, target) {
    // ...
  }
```

</details>

## Diagnostics

- **“Result of comparison with address of new variable of type "..." is always false”**

  This indicates a comparison like `ptr == &MyStruct{}` that will never be true. Consider these fixes:
  - _Compare values instead:_

    ```go
      *ptr == MyStruct{}
    ```

  - _Use `errors.As` for errors:_

    ```go
      var target *MyError
      if errors.As(err, &target) {
        // ...
      }
    ```

  - _Check dynamic type (for interface types):_

    ```go
      if v, ok := v.(*MyStruct); ok {
        if v.SomeField == expected { /* ... */ }
      }
    ```

  - _Pre-declare the target:_

    ```go
      var sentinel = &MyStruct{}
      // ...
      if ptr == sentinel { /* ... */ }
    ```

- **“Result of comparison with address of new variable of type "..." is false or undefined”**

  This diagnostic appears for zero-sized types where the comparison behavior is undefined:

  ```go
  type Skip struct{}

  func (e *Skip) Error() string { return "host hook execution skipped." }

  func (r renderRunner) RunHostHook(ctx context.Context, hook *hostHook) {
    if err := hook.run(ctx /*, ... */); errors.Is(err, &Skip{}) { // Undefined behavior.
      // ...
    }
  }
  ```

  or

  ```go
      defer func() {
        err := recover()

        if err, ok := err.(error); ok &&
          errors.Is(err, &runtime.PanicNilError{}) { // Undefined behavior.
          log.Print("panic called with nil argument")
        }
      }()

      panic(nil)
  ```

  While this might work due to Go runtime optimizations, its logic is unsound. Use `errors.As` instead:

  ```go
    var panicErr *runtime.PanicNilError
    if errors.As(err, &panicErr) {
      log.Println("panic error")
    }
  ```

## Integration

### `go vet`

```shell
go vet -vettool=$(which cmplint) ./...
```

### `errortype`

[`errortype`](https://github.com/fillmore-labs/errortype) is a comprehensive error handling linter that includes
`cmplint`'s checks. If you're already using `errortype`, you don't need `cmplint` separately — though running `cmplint`
alone might be faster on large codebases.

### `golangci-lint` Module Plugin

Add a `.custom-gcl.yaml` file to your project root:

```yaml
---
version: v2.9.0
plugins:
  - module: fillmore-labs.com/cmplint
    import: fillmore-labs.com/cmplint/gclplugin
    version: v0.0.6
```

See also the `golangci-lint`
[module plugin system](https://golangci-lint.run/docs/plugins/module-plugins/#the-automatic-way) documentation.

### CI/CD Integration

Add `cmplint` to your CI pipeline:

```yaml
# GitHub Actions example
jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6
      - uses: actions/setup-go@v6
        with:
          go-version: 1.25
      - name: Run cmplint
        run: go run fillmore-labs.com/cmplint@latest ./...
```

## Links

- [Equality of Pointers to Zero-Sized Types](https://blog.fillmore-labs.com/posts/zerosized-1/).

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.

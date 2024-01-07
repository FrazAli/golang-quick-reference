# About this document

This document is a quick reference for Go programming language.

The objective is to get a quick overview of basic syntax with examples, useful tools and tips.

**This is not a complete reference.** For a comprehensive reference, read the official [Go language specification](https://go.dev/ref/spec).

**The document structure is not organized in any particular order.**

# Introduction

Go is a compiled language. The key highlights of the language are: ease of use, modern concurrency features and performance.

Checkout [A tour of Go](https://go.dev/tour/list)

Learn more about Go [program initialization and execution](https://go.dev/ref/spec#Program_initialization_and_execution).

## Hello world

Similar to C language, a Go program requires an entry point function named `main`.

```go
import "fmt"

func main() {
  fmt.Println("Hello World!");
}
```

## The 'init' function

The `init()` function in Go allows to run code before the `main()` function.

```go
import "fmt"

func init() {
  fmt.Println("init is called before main!");
}
```

Read More: [https://go.dev/doc/effective_go#init](https://go.dev/doc/effective_go#init)

# Concurrency

Here is a brief overview of concurrency features in Go language. To learn more about it, read More: [https://go.dev/doc/effective_go#concurrency](https://go.dev/doc/effective_go#concurrency)

## Goroutines

- A goroutine is a lightweight thread.
- A goroutine can be started with `go` keyword.

```go
// Example 1
go func() {
  fmt.Println("Hello World! from goroutine");
}()

// Example 2
go list.Sort()  // run list.Sort concurrently; don't wait for it.
```

## Channels

- Channels are the preferred way to communicate/synchronize between goroutines.
- Channels can be created with `make` keyword.
- By default channels are bidirectional (`chan` keyword), but can be declared as unidirectional e.g. A send only unidirectional channel: `chan <- int`

```go
// uni-directional send-only channel
cs := make(chan<- int, 2)
cs <- 42
fmt.Println(<-co)  // Will cause error, since co is send-only

// uni-directional receive-only channel
cr := make(<-chan int, 2)
cr <- 42  // Will cause error, since cp is receive-only
```

- Since channels are types, they can be assigned to each other depending on the type of the channel.
  - When assigning channels to each-other,
    - channels can be assigned from bidirectional to unidirectional, not vice-versa.
    - unidirectional channels cannot be assigned to other unidirectional channels.

```go
c := make(chan int)
cr := make(<- chan int) // receive-only
cs := make(chan<- int)  // send-only

// uni-directional to bi-directional channel
c = cs  // ❌ Will cause error, since cs is send-only
c = cr  // ❌ Will cause error, since cs is receive-only

// bi-directional to uni-directional channel can be assigned
cs = c  // ✅
cr = c  // ✅
```

- By defaults channels are unbuffered (i.e. buffer size is zero).
- Unbuffered channels are blocking until the value on the channel is received.
- Beffered channels are non-blocking until the buffer is full.
- **Tip:** Use unbuffered channels unless there are specific reasons to use a buffered channel. The purpose of channels is to 'synchronize' and therefor a buffered channel is rarely used.

```go
// unbuffered channel of integers, default size zero
ci := make(chan int)
// unbuffered channel of integers, explicit size set to zero
cj := make(chan int, 0)
// buffered channel of pointers to Files with size 100
cf := make(chan *os.File, 100)

// putting values on the channel
ci <- 42

// receiving values from the channel
fmt.Println(<-ci)
v := <-ci


```

### Examples

- Signal using a channel

```go
c := make(chan int)  // Allocate a channel.
// Start the sort in a goroutine; when it completes, signal on the channel.
go func() {
  list.Sort()
  c <- 1  // Send a signal; value does not matter.
}()
doSomethingForAWhile()
<-c   // Wait for sort to finish; discard sent value.
```

- Send and receive using a channel

```go
package main

import (
  "fmt"
)

func main() {
  c := make(chan int)

  go send(c)

  for i := 0; i < 100; i++ {
    fmt.Println(<-c)
  }
}

// bi-directional to send-only channel as parameter
func send(c chan<- int) {
  for i := 0; i < 100; i++ {
    c <- i
  }
}
```

### Using `range` with Channels

- The `range` keyword can be used to iterate over a channel.
- The channel must be closed to stop the iteration.

```go
package main

import (
  "fmt"
)

func main() {
  c := make(chan int)

  go send(c)

  for v := range c {
    fmt.Println(v)
  }
}

func send(c chan<- int) {
  for i := 0; i < 100; i++ {
    c <- i
  }
  close(c) // Needed for range to work
}
```

### Using `select` with Channels

- The `select` keyword can be used to wait for multiple channels to receive a value.

```go
package main

import (
  "fmt"
)

func main() {
  even := make(chan int)
  odd := make(chan int)
  quit := make(chan bool)

  go send(even, odd, quit)
  receive(even, odd, quit)

  fmt.Println("Exiting")
}

func send(even, odd chan<- int, quit chan<- bool) {
  for i := 0; i < 100; i++ {
    if i%2 == 0 {
      even <- i
    } else {
      odd <- i
    }
  }
  close(quit)
}

func receive(even, odd <-chan int, quit <-chan bool) {
  for {
    select {
    case v := <-even:
      fmt.Println("Even: ", v)
    case v := <-odd:
      fmt.Println("Odd: ", v)
    // comma, ok idiom can be used with channels
    case i, ok := <-quit:
      if !ok {
        fmt.Println("Channel closed")
        return
      } else {
        fmt.Println("Received: ", i)
      }
  }
  }
}
```

### Patterns - Pipelines and Cancellation

**Tip:** It is best to follow well known concurrency patterns e.g. 'Fan-in' and 'Fan-out', and using 'Context' to cancel any recursively launched go routines.

- Read more about useful concurrency patterns: [https://go.dev/blog/pipelines](https://go.dev/blog/pipelines)
- Also read about how to use context to cancel go routines: [https://pkg.go.dev/context](https://pkg.go.dev/context)

## Wait Groups

- A `WaitGroup` is used to wait for a group of goroutines to finish.

```go
import (
  "fmt"
  "runtime"
  "sync"
)

var wg sync.WaitGroup

func main () {
  wg.Add(1) // Increment the WaitGroup coutner
  go func() {
    fmt.Println("Working from a goroutine!")
    wg.Done() // Decrements the WaitGroup Counter
  }()

  fmt.Println("Waiting to finish from the main")
  wg.Wait() // Waits for WaitGroup count to reach zero
}
```

See a more elaborate example: [https://pkg.go.dev/sync#example-WaitGroup](https://pkg.go.dev/sync#example-WaitGroup)

## Race Conditions

- Race conditions occur when multiple goroutines access shared data.
- The preferred way to avoid race conditions, is to not use shared memory, instead use Go channels to synchronize go routines.
- In case shared data is unavoidable, race conditions can be resolved by using [Mutexes](#mutexes) or by using the atomic operations through the standard library package [`sync/atomic`](https://pkg.go.dev/sync/atomic).

### Mutexes

- Mutexes are used to synchronize access to shared resources.
- See the example below to see how race conditions can be resolved using a mutex.
- **Tip:** Use `-race` flag to enable race detector in the command line. For example `go run -race main.go`

```go
package main

import (
  "fmt"
  "runtime"
  "sync"
)

var mu sync.Mutex

func main() {
  fmt.Println("CPUs: ", runtime.NumCPU())
  fmt.Println("Go routines: ", runtime.NumGoroutine())

  counter := 0

  const goroutines = 100
  var wg sync.WaitGroup
  wg.Add(goroutines)

  var mutex sync.Mutex

  for i := 0; i < goroutines; i++ {
    go func() {
      mutex.Lock()
      v := counter
      runtime.Gosched() // Yield the CPU - similar to time.Sleep()
      v++
      counter = v
      mutex.Unlock()
      wg.Done()
    }()
    fmt.Println("Go routines: ", runtime.NumGoroutine())
  }
  wg.Wait()
  fmt.Println("Go routines: ", runtime.NumGoroutine())
  fmt.Println("Counter: ", counter)
}
```

## Atomic Operations

The example above can be re-written to use atomic operations instead of mutexes.

```go
package main

import (
  "fmt"
  "runtime"
  "sync"
  "sync/atomic"
)

func main() {
  fmt.Println("CPUs: ", runtime.NumCPU())
  fmt.Println("Go routines: ", runtime.NumGoroutine())

  var counter int64 = 0 // int64 allows to use Atomic functions

  const goroutines = 100
  var wg sync.WaitGroup
  wg.Add(goroutines)

  for i := 0; i < goroutines; i++ {
    go func() {
      atomic.AddInt64(&counter, 1) // atomically increment
      runtime.Gosched()
      // atomically read
      fmt.Println("Counter: ", atomic.LoadInt64(&counter))
      wg.Done()
    }()
    fmt.Println("Go routines: ", runtime.NumGoroutine())
  }
  wg.Wait()
  fmt.Println("Go routines: ", runtime.NumGoroutine())
  fmt.Println("Counter: ", counter)
}
```

# Tooling

## Packages

- All code in go is organized in .go source files that can be grouped into packages. This is specified at the top of the source files e.g.

```go
package main
```

This indicates that this source file belongs to a package named 'main'.

- A package can be imported with `import` statement.

```go
import "fmt"
```

- Import external packages as dependencies with `go get` command.

```bash
go get github.com/user/repo/path/package@latest
```

- User `go mod tidy` to cleanup up the dependencies.

### Exported Names

- The name of all exported symbols in a package start with a **capital Letter** e.g. fmt.Printf

## Modules

- The go packages are organized in modules.
- The modules help manage dependencies and specify a global address to the module.

```bash
go mod init github.com/user/repo/path/modname
```

This command initializes a module by specifying the global address to the repository where the module source code will be hosted. The command generates a `go.mod` file in the current directory with module specification e.g.

```text
module github.com/user/repo/path/modname

go 1.21.4
```

## Build and Run

- Run a go program with `go run` command.

```bash
go run main.go
```

- Build a go program with `go build` command.

```bash
go build main.go
```

- Build for a specific platform with `GOOS=linux GOARCH=amd64 go build` command, e.g.

```bash
GOOS=linux go build
```

- Check Go environment variables

```bash
go env
```

- Go installs binaries in `GOPATH/bin`.

## Format

- Format a go file with `go fmt` command.
- The `go fmt` command is an alias to `gofmt -w -l`.

```bash
go fmt main.go
```

## Lint

- ~~Read more about the `golint` tool on the official github repo.~~[go get -u golang.org/x/lint/golint](go get -u golang.org/x/lint/golint)
- The `golint` tool is deprecated, instead, use `go vet` or use [`staticcheck`](https://staticcheck.dev/docs/getting-started/) tool.
- Install the `staticcheck` tool with `go install honnef.co/go/tools/cmd/staticcheck@latest`
- Lint a go file with `staticcheck` command.

```bash
staticcheck main.go
```

# Comments

Single line comments start with `//` at the beginning of a line. Multi-line comments are started with `/*` and ended with `*/`.

# Variables

- Go is a statically typed language.
- Variables can be declared with `var` keyword.
- By convention variable names are `camelCase`.

```go
var variableName = "value"
```

- The type is optional and it is inferred from the value at **compile time** if an initial value is provided.
- A shorthand for declaring a variable is `:=`, it cannot be used outside functions.

```go
variableName := "value"
a, b, d, _, s := 0, 1, 2.4, false, "Hello!"
```

- Type **must** be specified if variable is declared without an initial value.
- Unlike most other programming languages, variable type is specified after the variable name.

```go
var countItems int
```

Uninitialized variables are assigned the [zero values](#zero-values) automatically.

**Convention**: Use the shorthand (:=) variable declaration and assignment, the only exception is when no explicit value other than zero can be used to initialize the variable.

## Combined variable declaration

- You can declare multiple variables in one line, type specification at the end applies to all preceding variables.

```go
var a, b, c int
```

## Zero values

- **`0, 0.0`** for Numeric types
- **`false`** for Boolean types
- **`""`** for String types
- **`nil`** for Reference types

# Constants

Similar to most other programming languages constants are declared with `const` keyword.

```go
const Pi float64 = 3.14159265358979323846
const zero = 0.0         // untyped floating-point constant
const (
  size int64 = 1024
  eof        = -1  // untyped integer constant
)
// untyped integer and string constants
const a, b, c = 3, 4, "foo"  // a = 3, b = 4, c = "foo",
const u, v float32 = 0, 3    // u = 0.0, v = 3.0
```

## Iota

Iota is a "[constant generator](https://go.dev/ref/spec#Iota)" used in `const` enumerations. Its value starts at zero and is incremented automatically for each const specification, see examples below:

```go
const (
  c0 = iota  // c0 == 0
  c1 = iota  // c1 == 1
  c2 = iota  // c2 == 2
)
```

- iota only needs to be specifies in the first const specification.

```go
const (
  _ = iota   // iota = 0, ignored
  a          // a == 1  (iota == 1, implicit)
  b          // b == 2  (iota == 2, implicit)
  c          // c == 3  (iota == 3, implicit)
)
```

- `iota` is reset to zero for each new const statement and specification.

```go
const x = iota  // x == 0
const y = iota  // y == 0
```

- Multiple uses of iota in the same const specification have the same value:

```go
const (
  // bit0 == 1, mask0 == 0  (iota == 0)
  bit0, mask0 = 1 << iota, 1 << iota - 1
  bit1, mask1                  // bit1 == 2, mask1 == 1 (iota == 1)
  _, _                         // unused                (iota == 2)
  bit3, mask3                  // bit3 == 8, mask3 == 7 (iota == 3)
)
```

This last example exploits the implicit repetition of the last non-empty expression list.

```go
type ByteSize int

const (
  _           = iota // ignore first value by assigning to blank identifier
  KB ByteSize = 1 << (10 * iota)
  MB
  GB
  TB
  PB
  EB
)
```

Source: https://go.dev/ref/spec#Iota

# Pointers

- Similar to C/C++, you can declare pointers with `*` operator.
- Use address operator `&` to get the address of a variable.

```go
variableName := "value"
var pointerToVariable *string = &variableName

fmt.Printf("value '%v' of type %T at address %v\n",
  *pointerToVariable, *pointerToVariable, &variableName)
```

- Unlike C/C++, pointer arithmetic is not allowed.

# Printing

Print functions are part of the Go core package `fmt`. A package must be imported before any functions from it can be used. The print functions are similar to other programming languages, here are few examples:

```go
import "fmt"

func main() {
  x := 3
  y := 5

  // spaces are automatically added to parameters
  fmt.Println("The sum of", x, "and", y, "is", x+y);
  // Same output with a formatted print function
  fmt.Printf("The sum of %v and %v is %v\n", x, y, x+y);
  // Print the type of a variable
  fmt.Printf("The type of variable x is %T\n", x);
  // Raw string literals
  fmt.Println(`This is \n raw string without
    and is printed as it is`);
}
```

# Primitive Types

## Boolean:

```text
bool
```

## Numeric:

```text
uint8       unsigned  8-bit integers (0 to 255)
uint16      unsigned 16-bit integers (0 to 65535)
uint32      unsigned 32-bit integers (0 to 4294967295)
uint64      unsigned 64-bit integers (0 to 18446744073709551615)

int8        signed  8-bit integers (-128 to 127)
int16       signed 16-bit integers (-32768 to 32767)
int32       signed 32-bit integers (-2147483648 to 2147483647)
int64       signed 64-bit integers (-9223372036854775808 to 9223372036854775807)

float32     IEEE-754 32-bit floating-point numbers
float64     IEEE-754 64-bit floating-point numbers

complex64   complex numbers with float32 real and imaginary parts
complex128  complex numbers with float64 real and imaginary parts

byte        alias for uint8
rune        alias for int32, used to represent  a single unicode
```

> "Explicit conversions are required when different numeric types are mixed in an expression or assignment. For instance, int32 and int are not the same type even though they may have the same size on a particular architecture."

Source: [Go Reference: Numeric types](https://go.dev/ref/spec#Numeric_types)

The `int`, `uint`, and `uintptr` types are usually 32 bits wide on 32 bit systems and 64 bits wide on 64 bit systems.

**Convention**: Use `int` type unless there are specific reasons to use a signed or unsigned integer type.

## Strings:

- Strings are enclosed in double quotes.
- Strings are immutable sequence of unicode characters.
  - In comparison `[]byte` is a mutable array of bytes.
  - Go provides functions to convert between `[]byte` and `string`.

```text
string      sequence of characters
```

## Other

```text
any               any type
comparable        any type that supports all comparison operators
uintptr           unsafe pointer type (not typically used)
```

### Examples

```go
str := "Hello"
bytes := []byte(str)

bytes := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f}
str2 := string(bytes)
```

# Aggregate Types

## Arrays (Fixed length)

- Arrays are fixed sized and are not commonly used in Go.
- More commonly [slices](#slices) (dynamically sized) are preferred for storing elements of the same type.

### Syntax

```text
ArrayType   = "[" ArrayLength "]" ElementType
ArrayLength = Expression
ElementType = Type
```

### Examples

```text
[32]byte
[2*N] struct { x, y int32 }
[1000]*float64
[3][5]int
[2][2][2]float64  // same as [2]([2]([2]float64))

// use initialized values to set the array size
names := [...]string{"John", "Jonas", "Martin"}
```

## Slices (Dynamic length)

- Slices are built on top of arrays but are dynamically sized.
  See the source code for slice type: [src/runtime/slice.go](https://cs.opensource.google/go/go/+/refs/tags/go1.21.5:src/runtime/slice.go;drc=bdc6ae579aa86d21183c612c8c37916f397afaa8;l=16)
- use `append` builtin function to add elements to the slice
- Slices have a length (elements in the slice) and a capacity (elements in the underlying array). Read More: [https://go.dev/blog/slices-intro](https://go.dev/blog/slices-intro)
- As the slice grows, the capacity increases automatically. Behind the scenes, new larger array is allocated by copying data from old array to the new array and discarding the old array.
- Use `copy` builtin function to copy elements from one slice to another

### Syntax

```go
// Declaration
SliceType = "[" "]" ElementType

// Allocation
make([]T, length, capacity)
new([capacity]T)[0:length]
```

### Examples

```go
someSlice = []int

someSlice = make([]int, 50, 100)
fmt.Println(len(someSlice)) // length: 50
fmt.Println(cap(someSlice)) // capacity: 100

// OR less commonly
someSlice = new([100]int)[0:50]

// slice with initial values
names := []string{"John", "Jonas", "Martin"}

// append values to slice
names = append(names, "Fraz", "Garrett")

// Create a slice from another slice
lastTwo = names[3:] // ["Fraz", "Garrett"], storage is shared

// Deleting an element from slice, deletes element at index 3
names = append(names[0:3], names[4:]...) // '...' expands the slice values

// deep copy slices
newNames := make([]string, len(names))
copy(newNames, names)

// multi-dimensional slices
lists := [][]string{{"a", "b"}, {"c", "d"}} // [["a", "b"], ["c", "d"]]
```

## maps

- Maps are a collection of key-value pairs.
- Keys are unique and values are of the same type.
- Maps are unordered.
- Maps can also be created with `make` keyword.
- Use `delete` builtin function to delete elements from the map
- Map returns zero-value instead of an error if a key is not found, use 'comma, ok' idiom to check if an element exists in the map

```go
map[KeyType]ValueType
make(map[KeyType]ValueType)
```

### Examples

```go
myMap = map[string]int{}
someMap := make(map[string]int)

// assign values
myMap["Jake"] = 32 // map[string]int{"Jake": 32}

// delete an element from the map
delete(myMap, "Jake")

// check if an element exists in the map
v, ok := myMap["Jake"] // ok = false
```

## Structs

- Structs are a collection of values of different types.
- Structs can be declared as types with `struct` keyword.
- Structs can also be anonymous.
- Stucts can be composed by embedding other structs. The fields of embedded structs are promoted to the parent struct (This is similar to inheritance in other programming languages).

### Examples

```go
type Person struct {
  name string
  age int
}

p := Person{name: "John", age: 32}
fmt.Println(p)
fmt.Println(p.name, p.age) // access struct fields with dot operator

// Anonymous struct
person := struct {
  name string
  age int
}{
  name: "John",
  age: 32,
}

// Composition - embedded structs
type Employee struct {
  Person
  title string
}

emp := Employee{
  Person: p
  title: "Software Engineer",
}

/* Embedded struct fields / functions are promoted to the enclosing struct.
cessing Person fields directly from Employee.

fmt.Println("Employee:", emp.name, emp.age, emp.title) // Employee: John 32 Software Engineer

```

# Generic Types

- Generic types allow to write code that can work with any type.
- Generic types are denoted with the capital letter `T` by convention.
- Generic type can be declared to be any type with `any` keyword.
- Generic types can be used to restrict or group other types.

## Examples

```go
// Generic functions include a type specification in square brackets
func generic[T any](x T) T {}

// Generic type T that can be an `int` or `float64`
func generic[T int | float64](x T) T {}

// Typically interfaces are used to group types that can be used in generic functions
type myNumber interface {
  ~int | ~float64
}
func generic[T myNumber](x T) T {}
```

## Type Constraints - The `~` token / operator

- The `~` token is used to specify that the type specification applies to the underlying types. For example, `~string` means that a set of all types who have string as their underlying types.
- The `~` token is useful when types are redefined or when type aliases are used.
- By convention, use the builtin `constraints` package for various pre-defined constraint types (Complex, Integer, Float, Ordered, Signed, Unsigned). Read more about the [constraints package](https://pkg.go.dev/golang.org/x/exp/constraints).

### Examples\*\*

```go
// An alias of int type
type BigInt int

// An alias of float64
type BigFloat float64

// nyNum interface accepts any of the types: int, float64, BigInt, BigFloat
type myNum interface {
  ~int | ~float64
}
```

# Operators

## Logical Operators

```text
&&    conditional AND    p && q  is  "if p then q else false"
||    conditional OR     p || q  is  "if p then true else q"
!     NOT                !p      is  "not p"
```

## ... 'Operator'

`...` means different things depending on where it is used.

- Array size from initializer list

```go
// Declares an array where size is determined from initializer list
arr := [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
```

- Variadic parameters

```go
// Variadic parameters i.e. any number of arguments as a slice
func Println(a ...string) (int, error)
```

- Expand and pass slice

```go
// Expand a slice i.e. ["a", "b"] -> "a", "b"
names = []string{"John", "Jonas", "Martin"}
names = append(names, []string{"Fraz", "Garrett"}...)
```

- Recursive directory pattern

```bash
# Looks for tests in all sub-directories recursively
go test ./...
```

# Control Flow

## Conditionals

### If

```go
if true {
  fmt.Println("true")
} else if false {
  fmt.Println("false")
} else {
  fmt.Println("Not true or false")
}
```

### if with inline statements

```go
if i := 0; i < 5 {
  fmt.Println("less than 5")
}
```

### Switch

`fallthrough` is used to execute the next case if the current one matches. The example below prints both "zero" and "one".

```go
switch i := 0; i {
case 0:
  fmt.Println("zero")
  fallthrough
case 1:
  fmt.Println("one")
default:
  fmt.Println("other")


switch {
case x < 42:
  fmt.Println("less than 42")
case x == 42:
  fmt.Println("42")
case x > 42:
  fmt.Println("greater than 42")
default:
  fmt.Println("other")

```

## Loops

Unlike most other languages, in go `for` keyword is used for all types of loops.

```go
// Like a C for
for i := 0; i < 5; i++ {
  fmt.Println(i)
}

// Like a C while
for true { }

// Like a C for(;;) - use 'break' to exit
for { }

// range loop
for key, value := range myMap {
  fmt.Println(key, value)
}

// multi-variable for loop - Reverse 'a'
for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
  a[i], a[j] = a[j], a[i]
}
```

## Functions

### Syntax

```text
func (receiver) identifier(parameters) (returns) { body }
```

### Examples

```go
// multiple returns
func add(x, y int) (string, int) {
  return "sum", x + y
}

// variadic function
func sum(nums ...any) int {
  total := 0
  for _, num := range nums { // 'nums is slice type
    total += num.(int)
  }
  return total
}

// recursion
func factorial(n int) int {
  if n == 0 {
    return 1
  }
  return n * factorial(n-1)
}
```

### Anonymous functions

#### Syntax

```text
func(parameters) (returns) { body }
func(parameters) (returns) { body }(arguments)
```

- Similar to javascript anonymous functions can be used as 'Immediately executed function expressions'

#### Examples

```go
// Anonymous functions
func(x, y int) int {
  return x + y
}

// Anonymous functions inline call
func(s string) {
  fmt.Println("Hello", s)
}(" World!")
```

#### Function Expressions

- Similar to other languages, functions can be assigned to variables or passed around as arguments to other functions.
- Anonymous functions can be returned as a value.
- Anonymous functions can be passed as **callback**.

```go
add := func(x, y int) int {
  return x + y
}

ret := add(1, 2)

// Anonymous functions can be returned as a value
func adder() func(int, int) int {
  return func(x, y int) int {
    return x + y
  }
}

f := adder()
f(1, 2)

// Anonymous function  as callback
func doOp(a int, b int, f func(int, int) int) int {
  return f(a, b)
}

ret = doOp(1, 2, add)
```

##### Closure

```go
func inc() func() int {
  x := 0
  return func(x int) int {
    x++
    return x
  }
}

fmt.Println(adder()) // 1
fmt.Println(adder()) // 2
fmt.Println(adder()) // 3
```

### Parameters

- Similar to [combined variable declaration](#combined-variable-declaration), parameter types can be combined **if the parameters share the same type**, in such a case, the type of the last parameter applies to parameters before it. Compare the two equivalent examples below.

```go
// `int` type applies to both parameters
func add(a, b int) int {
  reutrn a + b
}
// Same as above
func add(a int, b int) int {
  reutrn a + b
}

```

### Return values

- Functions can return multiple values:

```go
func swap(x, y int) (int, int) {
  return y, x
}
```

#### Named Returns

```go
func split(sum int) (x, y int) {
  x = sum * 4 / 9
  y = sum - x
  return
}
```

### Deferred Function Calls

- Functions calls can be deferred with `defer` keyword.
- Multiple deferred function calls are executed in last-in-first-out order.

```go
func main() {
  // prints 4, 3, 2, 1, 0
  for i := 0; i < 5; i++ {
    defer fmt.Println(i)
  }
}
```

- Deferred functions calls are evaluated (parameters parsed) when the `defer` statement is executed, its only the execution that is 'deferred'
  - This behavior allow for useful features such as tracing.

```go
func trace(s string) string {
  fmt.Println("entering:", s)
  return s
}

func un(s string) {
  fmt.Println("leaving:", s)
}

func a() {
  defer un(trace("a"))
  fmt.Println("in a")
}

func b() {
  defer un(trace("b"))
  fmt.Println("in b")
  a()
}

func main() {
  b()
}

/*
ts

ring: b

ring: a

ing: a
ing: b

```

      Source: [Effective Go: defer](https://go.dev/doc/effective_go#defer)

## Methods

- Methods are functions attached to a type through the 'receiver' parameter before function name.

### Syntax

```text
func (someType) methodName(parameterList) (returnList)
```

### Examples

```go
type person struct {
  name string
  age  int
}

func (p person) toString() string {
  return fmt.Sprintf("%v (%v years)", p.name, p.age)
}
```

# Interfaces

- Interfaces declare behaviors (i.e. functions) that a type must implement.
- Many types can implement an interface i.e. Through interfaces 'polymorphism' is possible in Go.
- Types can implicitly implement an interface e.g. Any type that implements a function with signature `String() string` implicitly implements the builtin [`fmt.Stringer`](https://pkg.go.dev/fmt#Stringer) interface
  - Builtin interfaces allows to leverage the builtin functions (e.g. fmt.Println()) to be used with the custom types.

## Syntax

```text
type interfaceName interface {
methodName(parameterList) (returnList)
}
```

## Examples

```go
type person struct {
  name string
  age  int
}

type employee struct {
  person
}

// implicitly implements fmt.Stringer infterface
func (p person) String() string {
  return fmt.Sprintf("%v (%v years)", p.name, p.age)
}

func (p person) speak() string {
  return fmt.Sprintf("I am %v and I am %v years old", p.name, p.age)
}

func (e employee) speak() string {
  return fmt.Sprintf("I am %v", p.name)
}

type human interface {
  speak()
}

func saySomething(h human) {
  h.speak()
}

func main() {
  p := person{name: "John", age: 32}
  fmt.Println(p) // Uses the String function above
  log.Println(p) // Uses the 'Stringer' interface to log with timestamp

  // Test polymorphism of values through an interface 'human'
  e := employee{person{name: "Jenny", age: 28}}
  saySomething(p)
  saySomething(e)
}
```

# Error Handling

- There is no exception handling in Go, instead any errors are returned as an error type, and should be handled through comma-error idiom.
- **Tip:** Use comma-error idiom to check for errors as much as possible.
- The `errors` package provides a set of functions to work with errors.
- New errors can be created with the `errors.New` function.
  - The `fmt.Errorf` provides a more convenient way to create errors with parameters.
- The `log` package allows to write the errors with timestamp to stdout or to a file.
- Critical errors can be handled with the `panic` function which recursively stops all go routines.
- In case of a critical error that requires immediate termination, use the `log.Fatalf` function to log the error and exit, it internally uses `os.Exit` to exit with an error.

## Examples

```go
import (
  "errors"
  "fmt"
)

func main() {
  _, err := strconv.Atoi("a")
  if err != nil {
    fmt.Errorf("Error converting string to integer: %w", err)
  }
}
```

# Useful Applications using the Standard Library

## Random number

The `math/rand` package can be used to generate random numbers.
For example you can generate a pseudo random number between 1 and 100:

```go
package main

import (
  "fmt"
  "math/rand"
)

func main() {
  fmt.Println("Random number (0 - 99): ", rand.Intn(100))
}
```

## The `Writer` Interface

- The `io.Writer` interface is implemented by any value that has a method with the following signature:

```go
Write(p []byte) (n int, err error)
```

- For example, the [`File`](https://pkg.go.dev/os#File) type in Go implements the `io.Writer` interface.
- Similarly, `fmt.Println` function which wraps `fmt.Fprintln` function, that implements the `io.Writer` interface. It uses os.Stdout as the File type for writing to stdout.
- The example below shows how the `io.Writer` interface enables a consistent and powerful way to implement and use Interfaces in Go.

```go
import (
  "fmt"
  "os"
)

func main() {
  fmt.Println("Hello, World!")
  fmt.Fprintln(os.Stdout, "Hello, World!")
  io.WriteString(os.Stdout, "Hello, World!")
}
```

## Buffer IO

- `bytes` package in the standard library provides functions for working with byte slices and byte buffers.

### Examples

```go
type person struct {
  first string
  last string
}

func (p person) writeOut(w io.Writer) error {
  _, err := w.Write([]byte(p.first))
  return err
}

func main() {
  b := bytes.NewBufferString("Hello")
  fmt.Println(b.String()) // Hello
  b.WriteString(" World!")
  fmt.Println(b.String())  // Hello World!
  b.reset()
  b.WriteString("Wow!")
  fmt.Println(b.String())  // Wow!
  b.Write([]byte(" Cool!"))
  fmt.Println(b.String())  // Wow! Cool!

  // Write buffer to a file using the io.Writer interface
  f, _ := os.Create("test.txt")
  if err != nil {
    log.Fatalf("error %s", err)
  }
  defer f.close()

  p := person{first: "Morgan", last: "freeman"}

  var b bytes.Buffer
  p.writeOut(f)  // check and handle errors
  p.writeOut(&b) // check and handle errors
  fmt.Println(b.String())
}
```

## Sorting

- `sort` package in the standard library provides functions for sorting slices and maps.

### Examples

```go
import (
  "fmt"
  "sort"
)

func main() {
  s := []int{5, 2, 6, 3, 1, 4}
  names := []string[]{"John", "Jonas", "Martin", "Fraz", "Garrett"}

  sort.Ints(s)
  fmt.Println(s)  // [1 2 3 4 5 6]

  sort.String(names)
  fmt.Println(names) // ["Fraz" "Garrett" "John" "Jonas" "Martin"]
}
```

### Custom Sorting

- To sort a slice of custom types, implement the `sort.Interface` interface.
- The example below shows how to sort a slice of `person` structs by the age of the person.

```go
type person struct {
  name string
  age  string
}

func (p person) String() string {
  return fmt.Sprintf("%s (%s)", p.name, p.age)
}

func (p person) Len() int {
  return len(p)
}

func (p person) Swap(i, j int) {
  p[i], p[j] = p[j], p[i]
}

func (p person) Less(i, j int) bool {
  return p[i].age < p[j].age
}

func main() {
  persons := []person{
    {"Fraz", "35"},
    {"John", "32"},
    {"Garrett", "36"},
    {"Jonas", "33"},
    {"Martin", "34"},
  }

  sort.Sort(persons)
  fmt.Println(persons)
}
```

For more examples see the [sort package](https://pkg.go.dev/sort) in the standard library.

## Password Encryption

- `crypto/bcrypt` package in the standard library provides functions for password encryption and decryption.

### Example

```go
import (
  "crypto/bcrypt"
  "fmt"
)

func main() {
  password := "password123"
  hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
  fmt.Println(string(hash))

  // password from client
  err := bcrypt.CompareHashAndPassword(hash, []byte(password))
  if err != nil {
    fmt.Println(err)
    return
  }
  // Password matches
  fmt.Println("Success!")
}
```

## Marshalling and Unmarshalling JSON

- `encoding/json` package in the standard library provides functions for working with JSON.
- Marshalling means converting a Go struct to a JSON string.
- Unmarshalling means converting a JSON string to a Go struct.
- **Tip:** Use [json-to-go](https://mholt.github.io/json-to-go/) to convert JSON to Go structs.

### Example

```go
package main

import (
  "encoding/json"
  "fmt"
)

type person struct {
  first string
  last  string
}

func main() {
  p := person{first: "Morgan", last: "Freeman"}
  b, _ := json.Marshal(p)
  fmt.Println(string(b))

  // Marshalling JSON
  person_json, err := json.Marshal(person)
  if err != nil {
    log.Fatalf("Error marshalling person: %s", err)
  }
  fmt.Println(string(person_json))

  // Unmarshalling JSON
  var p1 person
  // without string literals: "{\"first\": "Morgan", \"last\": \"Freeman\"}"
  json.Unmarshal([]byte(`{"first": "Morgan", "last": "Freeman"}`), &p1)
  fmt.Println(p1)
}
```

## Print Runtime System Information

- `runtime` package in the standard library provides functions for printing some useful system information.

### Example

```go
import (
  "fmt"
  "runtime"
)

func main() {
  fmt.Println("OS:", runtime.GOOS)
  fmt.Println("ARCH:", runtime.GOARCH)
  fmt.Println("CPUs:", runtime.NumCPU())
  fmt.Println("Goroutines:", runtime.NumGoroutine())
  fmt.Println("GOMAXPROCS:", runtime.GOMAXPROCS(0))
  fmt.Println("Version:", runtime.Version())

  /* Example output
linux
: amd64
: 8
utines: 1
XPROCS: 8
ion: go1.21.4

}
```

# Testing

- **Tip:** Remember to BET i.e. Benchmark, Example and Test the code.
- Quick command reference:

```bash
go test
go test -bench .
go test -bench . -benchmem
go test -cover
go test -coverprofile coverage.out
go tool cover -html=coverage.out
```

## Unit tests

A unit test is a function that tests a single function, method, or struct.

- Test file names must end with `_test.go`
- Test files must be in the same package
- Test functions must start with `Test` followed by the function name that should start with a capital letter. e.g. `func TestSomeFunction(t *testing.T)`

### Examples

#### Basic unit test

```go
// main.go
package main

func Add(x, y int) int {
  return x + y
}

// main_test.go
package main

import "testing"

func TestAdd(t *testing.T) {
  result := Add(1, 2)
  if result != 3 {
    t.Errorf("Incorrect sum, want 3, got %d", Add(1, 2))
  }
}
```

#### Unit test example with mocked database using interface

```go
package main

// User represents a user with an id and first name
type User struct {
  ID int
  First string
}

// MockDatastore is a temporary service that stores retrievable data.
type MockDatastore struct {
  Users map[int]User
}

func (md MockDatastore) GetUser(id int) (User, error) {
  if user, ok := md.Users[id]; !ok {
    return User{}, fmt.Errorf("User %d not found", id)
  }
  return user, nil
}

func (md MockDatastore) SaveUser(u User) error {
  _, ok := md.Users[u.ID]
  if ok {
    return fmt.Errorf("User %d already exists", u.ID)
  }
  md.Users[u.ID] = u
  return nil
}

// Datastore defines an interface for storing retrievable data.
// Any TYPE that implements this interface (has these two methods) is also of TYPE `Datastore`.
// This means any value of TYPE `MockDatastore` is also of TYPE `Datastore`.
// This means we could have a value of TYPE `*sql.DB` and it can also be of TYPE `Datastore`
// This means we can write functions to take TYPE `Datastore` and use either of these:
// -- `MockDatastore`
// -- `*sql.DB`
type Datastore interface {
  GetUser(int) (User, error)
  SaveUser(u User) error
}

// Service represents a service that stores retrievable data.
// We will attach methods to `Service` so that we can use either of these:
// -- `MockDatastore`
// -- `*sql.DB`
type Service struct {
  ds Datastore
}

func (s Service) GetUser(id int) (User, error) {
  return s.ds.GetUser(id)
}

func (s Service) SaveUser(u User) error {
  return s.ds.SaveUser(u)
}

func main() {
  db := MockDatastore{
    Users: make(map[int]User),
  }

  srvc := Service{
    ds := db,
  }

  u1 := User{
    ID: 1,
    First: "Morgan",
  }

  err := srvc.SaveUser(u1)
  if err != nil {
    log.Fatalf("Error saving user: %s", err)
  }

  u1Returned, err := srvc.GetUser(u1.ID)
  if err != nil {
    log.Fatalf("Error getting user: %s", err)
  }

  fmt.Println(u1)
  fmt.Println(u1Returned)
}
```

The test file

```go
package main

import (
  "testing"
)

func TestGetUser(t *testing.T) {
  md := MockDatastore{
    Users: map[int]User{
      1: {ID: 1, First: "Morgan"},
    },
  }
  s := &Service{
    ds: md,
  }

  u, err := s.GetUser(1)
  if err != nil {
    t.Errorf("Error getting user: %s", err)
  }

  if u.First != "Morgan" {
    t.Errorf("Expected Morgan, got %s", u.First)
  }
}
```

## Example Tests

- Read more about example tests: [https://go.dev/blog/examples](https://go.dev/blog/examples)
- Unit tests can also be written in a way that they can be used as documentation examples.
- The examples on godoc.org are written as example tests.
- Example test names start with `Example`.
- Example tests use the print functions together with comments to indicate expected output.
- Example tests can be run similar to other unit tests, in addition the Example tests automatically show up as Examples on the generated HTML documentation.

### Example

```go
// main.go
package main

func Add(x, y int) int {
  return x + y
}
```

```go
// main_test.go
import "fmt"

func ExampleAdd(t *testing.T) {
  fmt.Println(Add(1, 2))
  // Output: 3
}
```

## Running tests

From the package root, run the tests with:

- Run all tests in current package: `go test`

## Coverage

- For test coverage, use -cover: `go test -cover`
- To generate coverage report, use -coverprofile: `go test -coverprofile=coverage.out`
- To display coverage report, use `go tool cover -html=coverage.out`
  - _This will generate an HTML report in `coverage.html`, with covered lines in green and uncovered lines in red._
- Learn more `go tool cover -h`

## Benchmark

- Benchmark is part of the go `testing` package i.e. `testing.B` or `testing.T`.
- Benchmark tests are added to the test module but with the prefix `Benchmark`.
- Benchmark a go program with `go test -bench=.`

```bash
go test -bench=. -benchmem
```

- read more about the benchmark flags: `go help testflag`

### Example

```go
// main_test.go, see main.go above for reference
package main

import (
  "testing"
)

func BenchmarkAdd(b *testing.B) {
  for i := 0; i < b.N; i++ {
    Add(3, 4)
  }
}
```

# Documentation

## The `go doc` command

- Read the help doc: `go doc -help`
- You can access documentation with `go doc` in the terminal.

```bash
go doc fmt  # Display documentation of the 'fmt' package
```

- You can view the documentation of your current package with `go doc` in the package root directory.

## Useful Links

- [godoc.org](https://godoc.org) - Documentation for third-party and standard library packages.
- [go.dev/doc](https://go.dev/doc) - Short tutorials and documentation for standard library packages.

## Package Documentation

### Generation

- Documentation for a package is generated with `godoc` command (`godoc` is not the same as `go doc`) in the terminal, the inline code comments are the main source of documentation in this case. See [Comment Syntax](#comment-syntax) for more details.
  - Install `godoc` package with `go install golang.org/x/tools/cmd/godoc@latest`
  - A local documentation server can be started with `godoc -http=:8080`

### Comment Recommendations

- Do not use URLs in comments, instead, include reference to any README that may include links / URLs
- Use short and clear sentences over long lines.
- Comments should provide context to the code that is not clear from the code itself.
- For detailed documentation, use a separate file `doc.go` in the package root.
  - For example, `doc.go` file for the `fmt` package can be found at: [`fmt/doc.go`](https://cs.opensource.google/go/go/+/refs/tags/go1.21.5:src/fmt/doc.go)

### Comment Syntax

See builtin go modules for syntax of code comments that allows documentation to be generated with `go doc` command. Here is a brief overview:

#### Package

A brief one-liner describing the purpose of the package.

```go
// Package math provides mathematical constants and functions.
package math
```

#### Function / Method

A brief one-liner describing the purpose of the function / method.

```go
// Add returns the sum of x and y.
func Add(x, y int) int {
  return x + y
}
```

#### Type

A brief one-liner describing the purpose of the type.

```go
// Person represents a person with a name and age.
type Person struct {
  name string
  age int
}
```

#### Constant / Variable

A brief one-liner describing the constant / variable.

```go
// Pi is the ratio of a circle's circumference to its diameter.
const Pi = 3.14
```

## Publishing godoc documentation to godoc.org

Once the documentation is generated, and pushed to github, it can be published to godoc.org with a very simple search.

- Copy the github url of the pakcage and paste it in the search bar on godoc.org
- The documentation should now be in the search indexed on godoc.org and can be searched using the fuzzy search.

# Interface Segregation in Go - Idiomatic Pattern

## Question
Why do we have the same interfaces defined in multiple packages?

In my code, I see `http` and `message` have the same interfaces:
```go
type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, ticketID string) error
}
```

## Answer

Having the same interfaces defined in both packages is actually **good practice** in Go for these reasons:

### 1. **Dependency Inversion**
Each package (`http` and `message`) defines only the interface it needs, not importing the actual implementations. This means:
- `project/http/handler.go` doesn't depend on the concrete spreadsheets or receipts packages
- `project/message/handlers.go` doesn't depend on the concrete implementations either
- Both depend on **abstractions**, not concretions

### 2. **Decoupling**
If you created a shared interface package, you'd introduce coupling:
```go
// ❌ BAD: Creates unnecessary dependencies
project/
  interfaces/      # shared package
    services.go    # both http and message import this
```

Now `http` and `message` both depend on a third package. Changes to one consumer affect the other.

### 3. **Interface Ownership**
In Go, **"Accept interfaces, return structs"** means the **consumer** owns the interface:
- The `http` package defines what **it** needs from a SpreadsheetsAPI
- The `message` package defines what **it** needs from a SpreadsheetsAPI
- They happen to be identical now, but could diverge in the future

### 4. **Future Flexibility**
If tomorrow the `http` handler needs a different method:
```go
// http/handler.go
type SpreadsheetsAPI interface {
    AppendRow(ctx context.Context, sheetName string, row []string) error
    GetRow(ctx context.Context, sheetName string, rowID int) ([]string, error) // new!
}
```

You can change it without affecting the `message` package at all.

## Practical Example

The concrete implementation (wherever it lives) satisfies **both** interfaces implicitly:

```go
// Concrete type satisfies both interfaces automatically
type GoogleSheets struct { /*...*/ }

func (g *GoogleSheets) AppendRow(ctx context.Context, sheetName string, row []string) error {
    // implementation
}
```

This is Go's **implicit interface satisfaction** at work—no explicit "implements" declaration needed.

## Key Principle

**"It's a good practice to keep interfaces close to the code that uses them."**

This follows the Go proverb: *"The bigger the interface, the weaker the abstraction."*

Each consumer defines the minimal interface it needs, making the code:
- More testable (easier to mock)
- More maintainable (changes are localized)
- More flexible (interfaces can evolve independently)

---

**TL;DR**: Duplicate interfaces are intentional. Each package defines the contract it needs, keeping packages decoupled and maintainable. This is idiomatic Go design!

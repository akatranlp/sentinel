//go:generate go tool go-enum --marshal
package components

// ENUM(error, warning, pending)
type ToastVariant string

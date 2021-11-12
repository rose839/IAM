package errors

import (
	"fmt"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("do some setup")
	m.Run()
	fmt.Println("do some cleanup")
}

func ExampleNew() {
	err := New("hello")
	fmt.Println(err)
	// output: hello
}

func ExampleWithMessage() {
	cause := New("whoops")
	err := WithMessage(cause, "oh noes")
	fmt.Println(err)

	// Output: oh noes
}

func ExampleWithMessage_printf() {
	cause := New("whoops")
	err := WithMessage(cause, "oh noes")
	fmt.Printf("%#+v", err)

	// Output: oh noes
}

func ExampleWrap() {
	cause := New("whoops")
	err := Wrap(cause, "oh noes")
	fmt.Println(err)

	// Output: oh noes
}

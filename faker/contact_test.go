package faker

import (
	"fmt"
	"testing"
)

func ExampleFaker_Contact() {
	Global.Seed(11)
	contact := Global.Contact()
	fmt.Println(contact.Phone)
	fmt.Println(contact.Email)
	// Output: 6136459948
	// carolecarroll@bosco.com
}

func BenchmarkContact(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Global.Contact()
	}
}

func ExampleFaker_Phone() {
	Global.Seed(11)
	fmt.Println(Global.Phone())
	// Output: 6136459948
}

func BenchmarkPhone(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Global.Phone()
	}
}

func ExampleFaker_PhoneFormatted() {
	Global.Seed(11)
	fmt.Println(Global.PhoneFormatted())
	// Output: 136-459-9489
}

func BenchmarkPhoneFormatted(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Global.PhoneFormatted()
	}
}

func ExampleFaker_Email() {
	Global.Seed(11)
	fmt.Println(Global.Email())
	// Output: markusmoen@pagac.net
}

func BenchmarkEmail(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Global.Email()
	}
}

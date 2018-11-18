package faker

import (
	"fmt"
	"testing"
)

func TestFaker_TemplateCustom(t *testing.T) {

}

func TestFaker_Template(t *testing.T) {
	for _, test := range []struct {
		template, should string
	}{
		{"", ""},
		{"{", "{"},
		{"{{", "{{"},
		{"{}", "{}"},
		{"{notfound}", "{notfound}"},
		{"Alfa", "Alfa"},
		{"Al{fa", "Al{fa"},
		{"{name}", "Jeromy Schmeler"},
		{"{name}}", "Jeromy Schmeler}"},
		{"{{name}}", "{Jeromy Schmeler}"},
		{"{{name}", "{Jeromy Schmeler"},
		{"{name}!", "Jeromy Schmeler!"},
		{"X{name}X", "XJeromy SchmelerX"},
		{"X{name}X{name}X", "XJeromy SchmelerXKim SteuberX"},
		{"{name", "{name"},
		{"name}", "name}"},
		{"name{", "name{"},
	} {
		fastFaker := NewFastFaker()
		fastFaker.Seed(42)
		got := fastFaker.Template(test.template)
		if got == test.should {
			continue
		}

		t.Errorf("fot template '%s', got '%s' should '%s'",
			test.template, got, test.should)
	}
}

func TestFaker_TemplateJson(t *testing.T) {
	for _, test := range []struct {
		template, should string
	}{
		{`{name:"{name}"}`, `{name:"Jeromy Schmeler"}`},
		{`["{name}", "{name}"]`, `["Jeromy Schmeler", "Kim Steuber"]`},

		{`{name:"{name}", age: {digit}}`, `{name:"Jeromy Schmeler", age: 8}`},
	} {
		fastFaker := NewFastFaker()
		fastFaker.Seed(42)
		got := fastFaker.Template(test.template)

		if got == test.should {
			continue
		}

		t.Errorf("fot template '%s', got '%s' should '%s'",
			test.template, got, test.should)
	}
}

func ExampleFaker_Template() {
	template := "Hello {name}!"

	fastFaker := NewFastFaker() // not concurrent safe, see NewSafeFaker()
	fastFaker.Seed(42)          //for each seed value will generate a different result

	fmt.Printf("%s\n", fastFaker.Template(template))
	fmt.Printf("%s\n", fastFaker.Template(template))

	// Output: Hello Jeromy Schmeler!
	//Hello Kim Steuber!
}

func ExampleFaker_TemplateJSON() {
	template := `{name:"{name}", age: {digit}}`

	fastFaker := NewFastFaker() // not concurrent safe, see NewSafeFaker()
	fastFaker.Seed(42)          //for each seed value will generate a different result

	fmt.Printf("%s\n", fastFaker.Template(template))
	// Output:{name:"Jeromy Schmeler", age: 8}
}

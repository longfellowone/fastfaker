package faker

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

// Template replaces all the found variables into the template with the actual
// results from the Faker.*() function.
//		fastFaker.Template("Hello {name}!") //Hello Jeromy Schmeler!
// For a list of all the variables see ./templates_variables.md
// All the "#" will be replaced by a digit and all the "?" by a ASCII letter.
// To use custom delimiters (instead of {}) see TemplateCustom()
func (f *Faker) Template(pattern string) string {

	//to allow simple patterns like phone numbers ##-###-###-###
	//and be backward compatibility with Generate()
	pattern = f.Numerify(pattern)
	pattern = f.Lexify(pattern)

	result, _ := f.TemplateCustom(pattern, "{", "}")
	return result
}

type keyPos struct {
	start, end   int
	variableFunc fakerer
}

// TemplateAllowedDelimiters the runes that are allowed as variable delimiters
// in Custom Templates. Must be ASCII (1 byte size) and not interfere with the regex expressions.
const TemplateAllowedDelimiters = "{}%#~<>-:@`"

func (f *Faker) TemplateCustom(template, delimStart, delimEnd string) (string, error) {

	//edge case, the template is only a variable "name"
	if delimStart == "" || delimEnd == "" {
		variableFunc, exists := templateVariables[template]
		if exists {
			return variableFunc(f), nil
		}
		return template, nil
	}

	for _, r := range delimStart + delimEnd {
		if strings.ContainsRune(TemplateAllowedDelimiters, r) == false {
			return "", fmt.Errorf("delimiters error, supported ones are: '%s'", TemplateAllowedDelimiters)
		}
	}

	//To allow more TemplateAllowedDelimiters we must add escape characters to each rune
	//for example for a delimiter "|||" we must transform it to "\|\|\|"

	//better Panic than sorry
	var pattern = regexp.MustCompile(`(` + delimStart + `[a-zA-Z0-9_-]+` + delimEnd + `)`)

	templateAsByte := []byte(template)
	indexes := pattern.FindAllIndex(templateAsByte, -1)

	//filter templateVariables, we will find all the locations in the template
	//of all variables.
	var toReplace []keyPos
	for _, match := range indexes {
		start := match[0]
		end := match[1]

		variableWithDelim := string(templateAsByte[start:end])
		sizeDelimStart := len(delimStart)
		sizeDelimEnd := len(delimEnd)

		//remove the delimiters, eg: {name} => name
		variable := strings.ToLower(variableWithDelim[sizeDelimStart : len(variableWithDelim)-sizeDelimEnd])

		variableFunc, exists := templateVariables[variable]
		if exists == false {
			//the variable does not exists
			continue
		}
		toReplace = append(toReplace, keyPos{start, end, variableFunc})
	}

	if len(indexes) == 0 {
		return template, nil
	}

	//cannot use strings.Builder to keep 1.0 compatibility
	buff := bytes.Buffer{}
	buff.Grow(len(templateAsByte) + 3*len(indexes)) //at least the input with some MagicNumber 3

	var lastEnd int
	//we go trough each byte and replace the variable with a value
	for _, posToReplace := range toReplace {
		buff.Write(templateAsByte[lastEnd:posToReplace.start])
		buff.WriteString(posToReplace.variableFunc(f))
		lastEnd = posToReplace.end
	}
	buff.Write(templateAsByte[lastEnd:])
	return buff.String(), nil
}

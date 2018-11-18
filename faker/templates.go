package faker

import (
	"bytes"
	"log"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// Template replaces all the found variables into the template with the actual
// results from the Faker.*() function.
//		fastFaker.Template("Hello {name}!") //Hello Jeromy Schmeler!
// For a list of all the variables see ./templates_variables.md
// To use custom delimiters (instead of {}) see TemplateCustom()
func (f *Faker) Template(pattern string) string {
	result, err := f.TemplateCustom(pattern, "{", "}")
	if err != nil {
		log.Panic("TemplateCustom is not working with default {}")
	}
	return result
}

type keyPos struct {
	start, end   int
	variableFunc fakerer
}

func (f *Faker) TemplateCustom(template, delimStart, delimEnd string) (string, error) {
	var pattern, err = regexp.Compile(`(\` + delimStart + `[a-zA-Z0-9_-]+\` + delimEnd + `)`)
	if err != nil {
		return "", errors.New("wrong delimiter characters, try {} or %%")
	}

	templateAsByte := []byte(template)
	indexes := pattern.FindAllIndex(templateAsByte, -1)

	//filter templateVariables
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

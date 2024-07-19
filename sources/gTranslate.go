package sources

import (
	"fmt"

	"github.com/bregydoc/gtranslate"
)

var result string

// Gtranslate translates text from one language to another using the gtranslate package.
func Gtrasnlate(text string, from string, to string) (string, error) {
	// Call gtranslate TranslateWithParams function to translate the text
	translated, err := gtranslate.TranslateWithParams(
		text,
		gtranslate.TranslationParams{
			From: from,
			To:   to,
		},
	)

	// Check if an error occurred during translation
	if err != nil {
		// Print the error message to standard output
		fmt.Println("Error:", err)
		// Return an empty string and the error
		return "", err
	}

	// If translation was successful, print the translated text
	fmt.Println(translated)

	// Set the result to the translated text
	result := translated

	// Return the translated text and nil (no) error
	return result, nil
}

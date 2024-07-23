package sources

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type TranslationResponse struct {
	ResponseData struct {
		TranslatedText string `json:"translatedText"`
	} `json:"responseData"`
}

func MyMemory(text string, from string, to string) (string, error) {
	lang := from + "|" + to

	// Build the API request URL
	apiURL := fmt.Sprintf("https://api.mymemory.translated.net/get?q=%s&langpair=%s", url.QueryEscape(text), lang)

	// Make the API request
	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("Error:", err)
		return result, err
	}
	defer resp.Body.Close()

	// Parse the JSON response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return result, err
	}

	// Unmarshal the JSON response into a struct
	var translationResponse TranslationResponse
	err = json.Unmarshal(body, &translationResponse)
	if err != nil {
		fmt.Println("Error:", err)
		return result, err
	}

	result := translationResponse.ResponseData.TranslatedText
	fmt.Println(result)

	return result, nil
}

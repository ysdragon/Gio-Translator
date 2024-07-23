package sources

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/corpix/uarand"
)

var (
	vqd4  string
	agent = uarand.GetRandom()
)

func DuckDuckGoAiTranslate(text string, from string, to string) (string, error) {

	statusURL := "https://duckduckgo.com/duckchat/v1/status"
	headers := http.Header{
		"Accept":             {"text/event-stream"},
		"Content-Type":       {"application/json"},
		"Cookie":             {"dcm=1"},
		"Origin":             {"https://duckduckgo.com"},
		"Referer":            {"https://duckduckgo.com/"},
		"Sec-Ch-Ua":          {`"Chromium";v="125", "Not.A/Brand";v="24"`},
		"Sec-Ch-Ua-Mobile":   {"?0"},
		"Sec-Ch-Ua-Platform": {`"Windows"`},
		"User-Agent":         {agent},
	}

	req, err := http.NewRequest("GET", statusURL, nil)
	if err != nil {
		return "Error", err
	}
	req.Header = headers
	req.Header.Set("x-vqd-accept", "1")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "Error", err
	}
	defer resp.Body.Close()

	vqd4 = resp.Header.Get("x-vqd-4")

	chatURL := "https://duckduckgo.com/duckchat/v1/chat"

	headers = http.Header{
		"Accept":             {"text/event-stream"},
		"Content-Type":       {"application/json"},
		"Cookie":             {"dcm=1"},
		"Origin":             {"https://duckduckgo.com"},
		"Referer":            {"https://duckduckgo.com/"},
		"Sec-Ch-Ua":          {`"Chromium";v="125", "Not.A/Brand";v="24"`},
		"Sec-Ch-Ua-Mobile":   {"?0"},
		"Sec-Ch-Ua-Platform": {`"Windows"`},
		"User-Agent":         {agent},
		"X-Vqd-4":            {vqd4},
	}

	data := fmt.Sprintf(`{"model":"claude-3-haiku-20240307","messages":[{"role":"user","content":"%s"}]}`, "Act as a professional translator and translate this text from "+from+" to "+to+" : "+text)
	req, err = http.NewRequest("POST", chatURL, strings.NewReader(data))
	if err != nil {
		return "Error", err
	}
	req.Header = headers

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return "Error", err
	}
	defer resp.Body.Close()

	contentStream := ""
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Error", err
	}

	lines := strings.Split(string(body), "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "data: ") && !strings.Contains(line, "[DONE]") {
			var data map[string]string
			json.Unmarshal([]byte(line[6:]), &data)
			message := data["message"]
			if message != "" {
				contentStream += message
			}
		}
	}

	return contentStream, err
}

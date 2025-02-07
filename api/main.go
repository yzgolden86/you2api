package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type YouChatResponse struct {
	YouChatToken string `json:"youChatToken"`
}

type OpenAIStreamResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Delta        Delta  `json:"delta"`
	Index        int    `json:"index"`
	FinishReason string `json:"finish_reason"`
}

type Delta struct {
	Content string `json:"content"`
}

type OpenAIRequest struct {
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
	Model    string    `json:"model"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []OpenAIChoice `json:"choices"`
}

type OpenAIChoice struct {
	Message      Message `json:"message"`
	Index        int     `json:"index"`
	FinishReason string  `json:"finish_reason"`
}

type ModelResponse struct {
	Object string        `json:"object"`
	Data   []ModelDetail `json:"data"`
}

type ModelDetail struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

var modelMap = map[string]string{
	"deepseek-reasoner":  "deepseek_r1",
	"deepseek-chat":      "deepseek_v3",
	"o3-mini-high":       "openai_o3_mini_high",
	"o3-mini-medium":     "openai_o3_mini_medium",
	"o1":                 "openai_o1",
	"o1-mini":            "openai_o1_mini",
	"o1-preview":         "openai_o1_preview",
	"gpt-4o":             "gpt_4o",
	"gpt-4o-mini":        "gpt_4o_mini",
	"gpt-4-turbo":        "gpt_4_turbo",
	"gpt-3.5-turbo":      "gpt_3.5",
	"claude-3-opus":      "claude_3_opus",
	"claude-3-sonnet":    "claude_3_sonnet",
	"claude-3.5-sonnet":  "claude_3_5_sonnet",
	"claude-3.5-haiku":   "claude_3_5_haiku",
	"gemini-1.5-pro":     "gemini_1_5_pro",
	"gemini-1.5-flash":   "gemini_1_5_flash",
	"llama-3.2-90b":      "llama3_2_90b",
	"llama-3.1-405b":     "llama3_1_405b",
	"mistral-large-2":    "mistral_large_2",
	"qwen-2.5-72b":       "qwen2p5_72b",
	"qwen-2.5-coder-32b": "qwen2p5_coder_32b",
	"command-r-plus":     "command_r_plus",
}

func getReverseModelMap() map[string]string {
	reverse := make(map[string]string, len(modelMap))
	for k, v := range modelMap {
		reverse[v] = k
	}
	return reverse
}

func mapModelName(openAIModel string) string {
	if mappedModel, exists := modelMap[openAIModel]; exists {
		return mappedModel
	}
	return "deepseek_v3"
}

func reverseMapModelName(youModel string) string {
	reverseMap := getReverseModelMap()
	if mappedModel, exists := reverseMap[youModel]; exists {
		return mappedModel
	}
	return "deepseek-chat"
}

var originalModel string

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/v1/chat/completions" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "You2Api Service Running...",
			"message": "MoLoveSze...",
		})
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Missing or invalid authorization header", http.StatusUnauthorized)
		return
	}
	dsToken := strings.TrimPrefix(authHeader, "Bearer ")

	var openAIReq OpenAIRequest
	if err := json.NewDecoder(r.Body).Decode(&openAIReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	originalModel = openAIReq.Model
	lastMessage := openAIReq.Messages[len(openAIReq.Messages)-1].Content
	var chatHistory []map[string]interface{}
	for _, msg := range openAIReq.Messages {
		chatMsg := map[string]interface{}{
			"question": msg.Content,
			"answer":   "",
		}
		if msg.Role == "assistant" {
			chatMsg["question"] = ""
			chatMsg["answer"] = msg.Content
		}
		chatHistory = append(chatHistory, chatMsg)
	}

	chatHistoryJSON, _ := json.Marshal(chatHistory)

	youReq, _ := http.NewRequest("GET", "https://you.com/api/streamingSearch", nil)

	q := youReq.URL.Query()
	q.Add("q", lastMessage)
	q.Add("page", "1")
	q.Add("count", "10")
	q.Add("safeSearch", "Moderate")
	q.Add("mkt", "zh-HK")
	q.Add("enable_worklow_generation_ux", "true")
	q.Add("domain", "youchat")
	q.Add("use_personalization_extraction", "true")
	q.Add("pastChatLength", fmt.Sprintf("%d", len(chatHistory)-1))
	q.Add("selectedChatMode", "custom")
	q.Add("selectedAiModel", mapModelName(openAIReq.Model))
	q.Add("enable_agent_clarification_questions", "true")
	q.Add("use_nested_youchat_updates", "true")
	q.Add("chat", string(chatHistoryJSON))
	youReq.URL.RawQuery = q.Encode()

	youReq.Header = http.Header{
		"sec-ch-ua-platform":         {"Windows"},
		"Cache-Control":              {"no-cache"},
		"sec-ch-ua":                  {`"Not(A:Brand";v="99", "Microsoft Edge";v="133", "Chromium";v="133"`},
		"sec-ch-ua-bitness":          {"64"},
		"sec-ch-ua-model":            {""},
		"sec-ch-ua-mobile":           {"?0"},
		"sec-ch-ua-arch":             {"x86"},
		"sec-ch-ua-full-version":     {"133.0.3065.39"},
		"Accept":                     {"text/event-stream"},
		"User-Agent":                 {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36 Edg/133.0.0.0"},
		"sec-ch-ua-platform-version": {"19.0.0"},
		"Sec-Fetch-Site":             {"same-origin"},
		"Sec-Fetch-Mode":             {"cors"},
		"Sec-Fetch-Dest":             {"empty"},
		"Host":                       {"you.com"},
	}

	cookies := getCookies(dsToken)
	var cookieStrings []string
	for name, value := range cookies {
		cookieStrings = append(cookieStrings, fmt.Sprintf("%s=%s", name, value))
	}
	youReq.Header.Add("Cookie", strings.Join(cookieStrings, ";"))

	if !openAIReq.Stream {
		handleNonStreamingResponse(w, youReq)
		return
	}

	handleStreamingResponse(w, youReq)
}

func getCookies(dsToken string) map[string]string {
	return map[string]string{
		"guest_has_seen_legal_disclaimer": "true",
		"youchat_personalization":         "true",
		"DS":                              dsToken,
		"you_subscription":                "youpro_standard_year",
		"youpro_subscription":             "true",
		"ai_model":                        "deepseek_r1",
		"youchat_smart_learn":             "true",
	}
}

func handleNonStreamingResponse(w http.ResponseWriter, youReq *http.Request) {
	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	resp, err := client.Do(youReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var fullResponse strings.Builder
	scanner := bufio.NewScanner(resp.Body)

	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "event: youChatToken") {
			scanner.Scan()
			data := scanner.Text()
			if !strings.HasPrefix(data, "data: ") {
				continue
			}
			var token YouChatResponse
			if err := json.Unmarshal([]byte(strings.TrimPrefix(data, "data: ")), &token); err != nil {
				continue
			}
			fullResponse.WriteString(token.YouChatToken)
		}
	}

	if scanner.Err() != nil {
		http.Error(w, "Error reading response", http.StatusInternalServerError)
		return
	}

	openAIResp := OpenAIResponse{
		ID:      "chatcmpl-" + fmt.Sprintf("%d", time.Now().Unix()),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   reverseMapModelName(mapModelName(originalModel)),
		Choices: []OpenAIChoice{
			{
				Message: Message{
					Role:    "assistant",
					Content: fullResponse.String(),
				},
				Index:        0,
				FinishReason: "stop",
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(openAIResp); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func handleStreamingResponse(w http.ResponseWriter, youReq *http.Request) {
	client := &http.Client{}
	resp, err := client.Do(youReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "event: youChatToken") {
			scanner.Scan()
			data := scanner.Text()

			var token YouChatResponse
			json.Unmarshal([]byte(strings.TrimPrefix(data, "data: ")), &token)

			openAIResp := OpenAIStreamResponse{
				ID:      "chatcmpl-" + fmt.Sprintf("%d", time.Now().Unix()),
				Object:  "chat.completion.chunk",
				Created: time.Now().Unix(),
				Model:   reverseMapModelName(mapModelName(originalModel)),
				Choices: []Choice{
					{
						Delta: Delta{
							Content: token.YouChatToken,
						},
						Index:        0,
						FinishReason: "",
					},
				},
			}

			respBytes, _ := json.Marshal(openAIResp)
			fmt.Fprintf(w, "data: %s\n\n", string(respBytes))
			w.(http.Flusher).Flush()
		}
	}
}

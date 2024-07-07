package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/richardwilkes/unison"
	"golang.design/x/clipboard"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func RunClipboardWatcher(markdown *unison.Markdown) {
	if err := clipboard.Init(); err != nil {
		panic(err)
	}

	ctx := context.Background()
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY is not set")
	}

	file, err := os.OpenFile("kotoba.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockNone,
		},
	}

	writeChan := make(chan genai.Part, 100)
	ch := clipboard.Watch(ctx, clipboard.FmtText)
	go writeToFile(file, writeChan)

	for content := range ch {
		data := invert(string(content))
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		part := genai.Text(
			fmt.Sprintf(
				"\n---------------------------------\n%s\n---------------------------------\n%s\n",
				timestamp,
				data,
			),
		)
		writeChan <- part
		fmt.Println(part)
		stream := model.GenerateContentStream(
			ctx,
			genai.Text(
				fmt.Sprintf(
					"Dịch sang tiếng Việt và giải thích ngữ pháp, viết cách đọc bằng hiragana, giải nghĩa tất cả các chữ Kanji xuất hiện trong đoạn văn kèm theo phiên âm furigana và từ Hán Nôm tương ứng:\n %s",
					data,
				),
			),
		)
		for {
			resp, err := stream.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			printResponse(resp, writeChan, markdown)
		}
	}
}

func writeToFile(file *os.File, writeChan chan genai.Part) {
	for part := range writeChan {
		fmt.Fprint(file, part)
	}
}

func printResponse(
	resp *genai.GenerateContentResponse,
	writeChan chan genai.Part,
	markdown *unison.Markdown,
) {
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				writeChan <- part
				fmt.Print(part)
				partText, ok := part.(genai.Text)
				if !ok {
					continue
				}
				currentText := markdown.ContentBytes()
				if len(currentText)+len(partText) > 10000 {
					currentText = currentText[(len(currentText) + len(partText) - 1000):]
				}
				currentText = append(currentText, partText...)
				markdown.SetContent(string(currentText), 0)
			}
		}
	}
}

func invert(s string) string {
	lines := strings.Split(s, "\n")
	s = ""
	for i := len(lines) - 1; i >= 0; i-- {
		s += lines[i]
		if i != 0 {
			s += "\n"
		}
	}
	return s
}

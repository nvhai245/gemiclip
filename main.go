package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"golang.design/x/clipboard"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func main() {
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

	for data := range ch {
		fmt.Println()
		fmt.Println("---------------------------------")
		fmt.Println()
		fmt.Println("---------------------------------")
		fmt.Println(string(data))
		fmt.Println()
		stream := model.GenerateContentStream(
			ctx,
			genai.Text(
				fmt.Sprintf(
					"Dịch sang tiếng Việt và giải thích ngữ pháp, giải nghĩa các chữ Kanji có trong đoạn văn kèm theo phiên âm furigana, và cung cấp từ Hán Việt tương ứng:\n %s",
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
			printResponse(resp, writeChan)
		}
	}
}

func writeToFile(file *os.File, writeChan chan genai.Part) {
	for part := range writeChan {
		fmt.Fprint(file, part)
	}
}

func printResponse(resp *genai.GenerateContentResponse, writeChan chan genai.Part) {
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				writeChan <- part
				fmt.Print(part)
			}
		}
	}
}

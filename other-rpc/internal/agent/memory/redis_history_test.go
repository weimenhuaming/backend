package memory

import (
	"testing"

	"github.com/tmc/langchaingo/llms"
)

func TestStoredMessageRoundTrip(t *testing.T) {
	cases := []llms.ChatMessage{
		llms.HumanChatMessage{Content: "你好"},
		llms.AIChatMessage{Content: "你好，有什么可以帮你？"},
		llms.SystemChatMessage{Content: "system prompt"},
	}

	for _, original := range cases {
		restored := fromStored(toStored(original))
		if restored.GetType() != original.GetType() {
			t.Fatalf("type mismatch: got %s want %s", restored.GetType(), original.GetType())
		}
		if restored.GetContent() != original.GetContent() {
			t.Fatalf("content mismatch: got %q want %q", restored.GetContent(), original.GetContent())
		}
	}
}

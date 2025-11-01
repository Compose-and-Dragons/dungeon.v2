---
marp: true
theme: default
paginate: true
---
<style>
.dodgerblue {
  color: dodgerblue;
}
.indianred {
  color: indianred;
}
.seagreen {
  color: seagreen;
}
</style>
### Golang 🩵 Generative AI: <span class="indianred">OpenAI API</span>

```golang
client := openai.NewClient(
    option.WithBaseURL("http://localhost:12434/engines/v1/"),
    option.WithAPIKey("I💙DockerModelRunner"),
)

messages := []openai.ChatCompletionMessageParamUnion{
    openai.SystemMessage(`Your name is Elara, Weaver of the Arcane`),
    openai.UserMessage("Tell me something about you"),
}

param := openai.ChatCompletionNewParams{
    Messages:    messages,
    Model:       "ai/qwen2.5:latest",
    Temperature: openai.Opt(0.7),    
}

completion, err := client.Chat.Completions.New(ctx, param)
```

<!--

-->

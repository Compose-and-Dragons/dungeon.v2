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
### Golang 🩵 Generative AI: **GenKit**

```golang
g := genkit.Init(ctx, genkit.WithPlugins(&openai.OpenAI{
    APIKey: "I💙DockerModelRunner",
    Opts: []option.RequestOption{
        option.WithBaseURL("http://localhost:12434/engines/v1/"),
    },
}))

resp, err := genkit.Generate(ctx, g,
    ai.WithModelName("openai/ai/qwen2.5:latest"),
    ai.WithSystem("Your name is Elara, Weaver of the Arcane."),
    ai.WithPrompt("Tell me something about you"),
    ai.WithConfig(map[string]any{"temperature": 0.7}),
)

fmt.Println(resp.Text())
```
<!--

-->

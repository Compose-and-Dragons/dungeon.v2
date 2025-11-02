---
marp: true
html: true
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

## 🛠️🪛⚙️🤖 Fine-Tuning of a LLM

- Taking <span class="dodgerblue">**a pre-trained AI model**</span> and training it further on <span class="indianred">**specialized data**</span> 
- <span class="dodgerblue">**Adapting**</span> general knowledge to <span class="indianred">**specific contexts**</span> (medical, legal, etc.)
- Fine-tuning parameters with <span class="seagreen">**targeted examples**</span> without starting from scratch
```json
[
  {
    "prompt": "What is your name and title?",
    "response": "My name is Queen Pédauque, also known as the Goose-Footed Queen."
  },
]
```
<!--
Fine-tuning is the technique of taking a pre-trained AI model (like GPT, Mistral, or LLaMA) and training it further on a specialized dataset tailored to a specific task or domain. This allows the model, which already has general knowledge, to adapt precisely to a particular context (e.g., medical vocabulary, legal corpus, specific writing style) without starting from scratch. The process fine-tunes the model's parameters using relevant examples, making it more effective for targeted applications.
-->


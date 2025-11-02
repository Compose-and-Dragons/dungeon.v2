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
</style>
# ⛓️‍💥 Constraints of Small Language Models

- Smaller context windows (ex: `4K`, `8K` tokens) 
  > *Maximum amount of text that can be processed at once (per request)*
- Fewer parameters (ex: <span class="dodgerblue">`qwen2.5:`**`0.5b`**</span>)
  - *Less Capacity to learn/understand complex patterns*
  - *But Less Resources needed (memory, compute)*
- Less training data
- Bad at function calling (🤬 with MCP)

<!--
Smaller context windows: Maximum amount of text that can be processed at once (per request)
- Taille maximale de texte traitable en une fois
- Limite la quantité d'information que le modèle peut "voir" en une fois.
Fewer parameters: Determines the model's capacity to learn complex patterns and the resources needed (memory, compute).
-->


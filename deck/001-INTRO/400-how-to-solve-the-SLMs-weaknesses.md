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
# ⛓️‍💥 Solving the <span class="indianred">Weaknesses</span> of the SLMs

| Weaknesses              | Solutions                 | Methods          |
|-------------------------|---------------------------|------------------|
| Less training data 1    | provide more data.        | **Fine tuning**  |
| Less training data 2    | provide data on demand    | **RAG**, **MCP** |
| Fewer parameters        | externalize some tasks.   | **MCP** (A2A)    |
| Smaller context windows | smaller but accurate data | **RAG**          |
| Bad at function calling | SLMs natively optimized for MCP| 🕵️‍♂️🔎🔎🔎          |

> - *RAG: Retrieval-Augmented Generation*
> - *MCP: Model Context Protocol*, *A2A: Agent to Agent*

<!--

-->
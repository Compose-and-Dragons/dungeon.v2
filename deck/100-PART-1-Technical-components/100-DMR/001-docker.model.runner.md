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
# 🐳🧠 Docker Model Runner (**DMR**)

#### What?
- Docker **plugin** `==` <span class="dodgerblue">**LLM Engine**</span> (based on `llama.cpp`)
> and soon with **vLLM** support 🎉
#### Features
* Pull and push **models** to and from Docker Hub (or any Container Registry)
* Pull **models** directly from 🤗 **Hugging Face** 
* Package **GGUF** model files as **OCI Artifacts** and publish them
* Run and interact with AI models directly
  - From the **command line** or the **Docker Desktop GUI**
  - <span class="dodgerblue">**OpenAI**</span> compatible APIs


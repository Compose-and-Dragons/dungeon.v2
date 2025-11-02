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

## 🛠️🪛⚙️🤖 Fine-Tuning of a LLM - <span class="indianred">Tools</span>

| Task                                      | Solutions                   | 
|-------------------------------------------|-----------------------------|
| 📝 **Specialized** Dataset + **Cleaning** | **Claude** AI                   | 
| ⚙️ Launch Fine-tuning **Training**        | 🦥 **Unsloth** + Docker **Offload** 1️⃣ |
| 🤔 **Evaluate** on **New** Examples       | Docker **Model Runner**         |
| 🚀 **Deploy** the Fine-tuned Model        | Docker **Hub** 2️⃣ |

</br></br></br>
> 1️⃣ Or AMD + NVIDIA GPUs with **Unsloth**
> 2️⃣ Or private **Model Registry** (OCI compliant)
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
</style>
# AI Agent?
> Augmented 🧠 LLM with 🛠️ Tools

<div class="mermaid">
flowchart LR
    In((In)) --> LLM[GenAI app + LLM]
    LLM --> Out((Out))
    
    LLM -.-> |Query/Results| Retrieval[LLM Data/RAG/Context]
    Retrieval -.-> LLM
    
    LLM -.-> |Call/Response| Tools(Tools/MCP)
    Tools -.-> |Feedback| LLM


    Tools -->|Action| Environment((Environment))
    Environment -->|Feedback| Tools
    
    
    LLM -.-> |Read/Write| Memory[Memory]
    Memory -.-> LLM
    
    style In fill:#FFDDDD,stroke:#DD8888
    style Out fill:#FFDDDD,stroke:#DD8888
    style LLM fill:#DDFFDD,stroke:#88DD88
    style Retrieval fill:#DDDDFF,stroke:#8888DD
    style Tools fill:#DDDDFF,stroke:#8888DD
    style Memory fill:#DDDDFF,stroke:#8888DD

    classDef env fill:#FFDDDD,stroke:#DD8888
    class Environment env

</div>

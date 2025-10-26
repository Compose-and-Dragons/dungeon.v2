# Generate dataset for NPC agents to fine-tune language models

## Description

This program automatically generates training datasets for Language Models (LLMs) to create Non-Player Characters (NPCs) with consistent personalities in a role-playing game.

### How it works

The program analyzes a document describing an NPC's background and personality, splits it into logical sections, then uses an LLM to generate prompt/response pairs that capture the character's essence. This data can then be used to fine-tune a language model specific to the character.

### Main steps

1. **Configuration**: Loads parameters from environment variables (model URL, temperature, etc.)
2. **Context loading**: Reads the markdown file containing the character description
3. **RAG splitting**: Divides the document into chunks by sections for optimal processing
4. **Generation**: For each chunk, asks the LLM to generate N prompt/response pairs
5. **Validation**: Retries up to 3 times in case of JSON parsing failure
6. **Save**: Exports all entries to a training JSON file

### Environment variables

- `MODEL_RUNNER_BASE_URL`: LLM engine API URL
- `CHAT_MODEL`: Model identifier to use
- `MODEL_TEMPERATURE`: Model temperature (creativity)
- `MODEL_TOP_P`: Top-P sampling
- `SYSTEM_INSTRUCTIONS`: System instructions for the agent
- `CONTEXT_PATH`: Path to the character description file
- `DATASET_ENTRIES_PER_CHUNK`: Number of entries to generate per chunk
- `NPC_NAME`: Character name
- `TEMPLATE`: Template for prompt generation

## Architecture diagrams

### Process overview

```mermaid
flowchart TD
    A[Start] --> B[Load Configuration]
    B --> C[Read character file]
    C --> D[Split into RAG chunks]
    D --> E{For each chunk}
    E --> F[Interpolate template]
    F --> G[Call LLM JSON Completion]
    G --> H{Success?}
    H -->|No| I{Retry < 3?}
    I -->|Yes| G
    I -->|No| J[Skip chunk]
    H -->|Yes| K[Parse JSON]
    K --> L[Add to dataset]
    J --> M{More chunks?}
    L --> M
    M -->|Yes| E
    M -->|No| N[Save training_data.json]
    N --> O[End]
```

### Component architecture

```mermaid
graph LR
    A[main.go] --> B[NPCAgent]
    A --> C[helpers]
    A --> D[rag]
    B --> E[LLM API]
    C --> F[Env Variables]
    C --> G[File I/O]
    D --> H[Markdown Parser]

    style A fill:#f9f,stroke:#333
    style E fill:#bbf,stroke:#333
    style F fill:#bfb,stroke:#333
```

### Data flow

```mermaid
sequenceDiagram
    participant M as Main
    participant A as NPCAgent
    participant R as RAG
    participant L as LLM
    participant F as FileSystem

    M->>F: Read character.md
    F-->>M: Markdown content
    M->>R: SplitMarkdownBySections()
    R-->>M: List of chunks

    loop For each chunk
        M->>M: Interpolate template
        M->>A: JsonCompletion(chunk)
        A->>L: Generation request
        L-->>A: JSON prompt/response
        A-->>M: DataSetEntry[]
        M->>M: Append to trainingData
    end

    M->>F: Write training_data.json
    F-->>M: Success
```

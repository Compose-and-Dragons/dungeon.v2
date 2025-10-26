# Dataset Generation Strategy for Fine-Tuning Small Language Models

## Executive Summary

This document describes a **chunked iterative generation strategy** designed specifically for fine-tuning small language models (SLMs) with limited context windows. The approach breaks down large character descriptions into manageable chunks and uses an iterative refinement process to generate diverse, high-quality training data.

**Key Innovation**: Instead of relying on large models like GPT-4 or Claude to generate thousands of examples in one pass, this strategy uses small models (like jan-nano-128k or qwen2.5:0.5b) to generate datasets incrementally, making the process cost-effective and accessible.

---

## Problem Statement

### Challenges with Small Models

When fine-tuning small language models for NPC (Non-Player Character) generation in games, we face several challenges:

1. **Limited Context Window**: Small models typically have context windows of 8K-16K tokens, preventing them from processing large character sheets in one pass
2. **Cost Constraints**: Using large commercial models (GPT-4, Claude) for dataset generation is expensive
3. **Quality Requirements**: Need 1000+ diverse, high-quality examples for effective fine-tuning
4. **Consistency**: Must maintain character consistency across all generated examples
5. **JSON Reliability**: Small models often struggle with consistent JSON output format

---

## Solution Architecture

### High-Level Strategy

```mermaid
graph TB
    A[Large Character Document] --> B[Manual Chunking by Section]
    B --> C[Chunk 1: Profile]
    B --> D[Chunk 2: Behavior]
    B --> E[Chunk 3: Appearance]
    B --> F[Chunk N: ...]

    C --> G[Iterative Generation Process]
    D --> G
    E --> G
    F --> G

    G --> H[Individual JSON Files]
    H --> I[Merge All Files]
    I --> J[Final training_data.json]

    style A fill:#e1f5ff
    style G fill:#fff4e1
    style J fill:#e8f5e9
```

### Core Components

The system consists of three main components:

1. **Document Preparation**: Manual chunking of character sheets into focused sections
2. **Iterative Generation Engine**: Multi-pass generation with refinement
3. **Merge Pipeline**: Aggregation of all generated datasets

---

## Detailed Process Flow

### 1. Document Preparation Phase

```mermaid
flowchart LR
    A[Character Sheet<br/>queen_pedauque_character_sheet.md] --> B{Manual Analysis}
    B --> C[Identify Logical Sections]
    C --> D[Profile]
    C --> E[Behavior]
    C --> F[Appearance]
    C --> G[Magic Powers]
    C --> H[History]
    C --> I[Special Rules]

    D --> J[Create Chunk Files]
    E --> J
    F --> J
    G --> J
    H --> J
    I --> J

    J --> K[01.profile.md]
    J --> L[02.behavior.md]
    J --> M[03.appearance.md]
    J --> N[...]

    style A fill:#bbdefb
    style J fill:#c8e6c9
```

Each chunk file contains **4 sections separated by `----------` delimiter**:

1. **Context**: Brief system instructions for the character
2. **Document**: The actual content to generate examples from
3. **Settings**: JSON configuration for generation parameters
4. **Prompt Template**: Template for generating questions

**Example Chunk Structure**:
```markdown
# Queen Pédauque - Legendary Sorceress NPC System Instructions
You are Queen Pédauque, a legendary sorceress NPC...
----------
## Character Profile
- **Name**: Queen Pédauque (also known as the Goose-Footed Queen)
- **Race**: Human with Ancient Fae Heritage
...
----------
{
    "nb_dataset_entries": 20,
    "nb_iterations": 5
}
----------
From this document related to {{.NameOfTheNPC}}:
{{.Chunk}}
Generate {{.NbEntriesPerChunk}} dataset entries...
```

---

### 2. Iterative Generation Engine

```mermaid
sequenceDiagram
    participant M as Main Program
    participant C as Chunk File
    participant T as Template Engine
    participant A as NPCAgent
    participant L as Small LLM
    participant F as File System

    M->>C: Read chunk file
    C-->>M: Context, Document, Settings, Template

    Note over M: Parse Settings JSON<br/>(nb_dataset_entries, nb_iterations)

    M->>T: Interpolate template with data
    T-->>M: Formatted prompt

    rect rgb(220, 240, 255)
        Note over M,L: ITERATION 1
        M->>A: JsonCompletion(prompt)
        A->>L: Generate 5 examples
        L-->>A: JSON array
        A-->>M: 5 DataSetEntry objects
        M->>M: Append to trainingData
    end

    rect rgb(255, 240, 220)
        Note over M,L: ITERATIONS 2-5
        loop 4 more times
            M->>M: Modify prompt with previous response
            M->>A: JsonCompletion(enhanced prompt)
            A->>L: Generate 5 NEW examples
            L-->>A: JSON array
            A-->>M: 5 DataSetEntry objects
            M->>M: Append to trainingData
        end
    end

    M->>F: Write 1.queen_pedauque_training_data.json
    F-->>M: Success
```

---

### 3. Generation Loop with Retry Logic

```mermaid
flowchart TD
    A[Start Generation for Chunk] --> B[Read Chunk Settings]
    B --> C{Parse Settings}
    C -->|Success| D[nb_dataset_entries = 5<br/>nb_iterations = 5]
    C -->|Failure| Z[Fatal Error]

    D --> E[Iteration Counter = 1]
    E --> F[Interpolate Prompt Template]
    F --> G[Call LLM JsonCompletion]

    G --> H{JSON Valid?}
    H -->|No| I{Retry < 3?}
    I -->|Yes| J[Retry Counter++]
    J --> G
    I -->|No| K[Skip This Iteration]

    H -->|Yes| L[Parse JSON to DataSetEntry array]
    L --> M[Append 5 entries to trainingData]
    M --> N{Iteration < 5?}

    N -->|Yes| O[Iteration Counter++]
    O --> P[Modify Prompt:<br/>Include previous response<br/>Request DIFFERENT questions]
    P --> G

    N -->|No| Q[Total Entries = 25]
    Q --> R[Marshal to JSON with Indent]
    R --> S[Write file:<br/>fileIndex.trainDataFile]
    S --> T[fileIndex++]
    T --> U{More Chunks?}

    U -->|Yes| A
    U -->|No| V[Merge All Files]
    V --> W[End]

    K --> N

    style G fill:#fff4e1
    style H fill:#ffebee
    style L fill:#e8f5e9
    style V fill:#e1f5ff
```

**Key Points**:
- Each chunk generates `nb_dataset_entries × nb_iterations` examples (e.g., 5 × 5 = 25)
- Each iteration explicitly asks for **DIFFERENT** questions to ensure diversity
- Retry mechanism (up to 3 attempts) handles JSON parsing failures
- Failed iterations are skipped rather than crashing the entire process

---

### 4. Iterative Refinement Strategy

The system uses a clever prompt modification technique to ensure diversity:

```mermaid
graph TD
    A[Iteration 1:<br/>Base Prompt] --> B[LLM generates<br/>5 examples]
    B --> C[Save Response 1]

    C --> D[Iteration 2:<br/>Enhanced Prompt]
    D --> E["Prompt includes:<br/>• Previous response<br/>• Explicit request for NEW questions<br/>• Same document context"]
    E --> F[LLM generates<br/>5 DIFFERENT examples]
    F --> G[Save Response 2]

    G --> H[Iteration 3-5:<br/>Continue pattern]
    H --> I[Total: 25 diverse examples<br/>from single chunk]

    style A fill:#e3f2fd
    style D fill:#fff3e0
    style H fill:#fce4ec
    style I fill:#e8f5e9
```

**Example Prompt Evolution**:

**Iteration 1**:
```
From this document related to Queen Pédauque:
[Document content...]
Generate 5 dataset entries...
```

**Iteration 2**:
```
Here is the previous response you gave:
[Previous JSON array...]

Now, please generate NEW more dataset entries for the same document
but with DIFFERENT prompts:
[Same document content...]
Generate 5 dataset entries...
```

This approach:
- ✅ Maintains context awareness
- ✅ Encourages diversity by showing what was already generated
- ✅ Keeps prompts focused and manageable for small models

---

## File Organization and Merging

### Individual File Generation

```mermaid
graph LR
    A[Chunk 01.profile.md] --> B[1.queen_pedauque_training_data.json]
    C[Chunk 02.behavior.md] --> D[2.queen_pedauque_training_data.json]
    E[Chunk 03.interactions.md] --> F[3.queen_pedauque_training_data.json]
    G[Chunk 18.cultural-significance.md] --> H[18.queen_pedauque_training_data.json]

    style B fill:#e8f5e9
    style D fill:#e8f5e9
    style F fill:#e8f5e9
    style H fill:#e8f5e9
```

### Merge Process

```mermaid
flowchart TD
    A[Start Merge] --> B[Get all files in directory]
    B --> C{For each file}
    C --> D{Ends with<br/>trainDataFile?}
    D -->|No| E[Skip file]
    D -->|Yes| F[Read JSON file]
    F --> G{Valid JSON<br/>array?}
    G -->|No| H[Log warning, continue]
    G -->|Yes| I[Parse to DataSetEntry array]
    I --> J[Append all entries to mergedData]

    E --> K{More files?}
    H --> K
    J --> K
    K -->|Yes| C
    K -->|No| L[Marshal merged array to JSON]
    L --> M[Write training_trainDataFile]
    M --> N[Report total entries]
    N --> O[End]

    style F fill:#e3f2fd
    style I fill:#fff3e0
    style M fill:#e8f5e9
```

**Result**:
- Input: `1.queen_pedauque_training_data.json`, `2.queen_pedauque_training_data.json`, ... `18.queen_pedauque_training_data.json`
- Output: `training_queen_pedauque_training_data.json` (single merged file with 450+ examples)

---

## Configuration Management

### Docker Compose Configuration

The system uses Docker Compose for reproducible environments:

```yaml
environment:
  MAX_GENERATION_RETRIES: 3
  NPC_NAME: Queen Pedauque
  SYSTEM_INSTRUCTIONS: |
    You are a helpful AI assistant...
  MODEL_CONTEXT_DIRECTORY: /app/data/chunks
  TRAINING_DATA_DIRECTORY: /app/data
  TRAINING_DATA_FILE: queen_pedauque_training_data.json
  MODEL_TEMPERATURE: 1.0
  MODEL_TOP_P: 0.9

models:
  chat-model:
    model: hf.co/menlo/jan-nano-128k-gguf:q4_k_m
    context_size: 16384
```

### Parameter Tuning Guide

```mermaid
graph TD
    A[Configuration Parameters] --> B[Generation Parameters]
    A --> C[Model Parameters]
    A --> D[Retry Parameters]

    B --> B1[nb_dataset_entries<br/>Per-iteration count]
    B --> B2[nb_iterations<br/>Refinement passes]

    C --> C1[MODEL_TEMPERATURE<br/>Creativity level]
    C --> C2[MODEL_TOP_P<br/>Sampling diversity]
    C --> C3[context_size<br/>Token limit]

    D --> D1[MAX_GENERATION_RETRIES<br/>Failure tolerance]

    style B fill:#e3f2fd
    style C fill:#fff3e0
    style D fill:#ffebee
```

**Recommended Settings**:

| Parameter | Recommended Value | Rationale |
|-----------|------------------|-----------|
| `nb_dataset_entries` | 5 | Manageable for small models, balances quality vs quantity |
| `nb_iterations` | 5 | Provides diversity without excessive redundancy |
| `MODEL_TEMPERATURE` | 1.0 | Higher creativity for varied questions |
| `MODEL_TOP_P` | 0.9 | Balanced sampling for coherent but diverse outputs |
| `MAX_GENERATION_RETRIES` | 3 | Tolerates occasional JSON failures without infinite loops |

---

## Data Quality Assurance

### Quality Control Mechanisms

```mermaid
flowchart LR
    A[Raw Generation] --> B[JSON Schema Validation]
    B --> C{Valid Structure?}
    C -->|No| D[Retry up to 3 times]
    C -->|Yes| E[Content Validation]

    E --> F{Contains required fields?}
    F -->|No| D
    F -->|Yes| G[Diversity Check]

    G --> H[Compare with previous iterations]
    H --> I{Too similar?}
    I -->|Yes| J[Prompt emphasizes NEW questions]
    I -->|No| K[Accept entry]

    D --> L{Retries exhausted?}
    L -->|Yes| M[Skip iteration]
    L -->|No| B

    K --> N[Add to training dataset]
    M --> O[Log warning]

    style B fill:#e3f2fd
    style E fill:#fff3e0
    style G fill:#f3e5f5
    style K fill:#e8f5e9
    style M fill:#ffebee
```

### Output Format Validation

Each generated entry must conform to:

```json
{
  "prompt": "question string (non-empty)",
  "response": "answer string (non-empty, based only on provided document)"
}
```

**Validation Rules**:
1. Must be valid JSON array
2. Each element must have exactly 2 fields: `prompt` and `response`
3. Both fields must be non-empty strings
4. Responses must be grounded in the source document (no hallucinations)
5. Prompts must be diverse across iterations

---

## Scalability and Performance

### Horizontal Scaling Strategy

```mermaid
graph TB
    A[Character Sheet] --> B[Manual Chunking]
    B --> C1[Chunk Group 1<br/>Chunks 1-6]
    B --> C2[Chunk Group 2<br/>Chunks 7-12]
    B --> C3[Chunk Group 3<br/>Chunks 13-18]

    C1 --> D1[Worker Container 1]
    C2 --> D2[Worker Container 2]
    C3 --> D3[Worker Container 3]

    D1 --> E1[Files 1-6]
    D2 --> E2[Files 7-12]
    D3 --> E3[Files 13-18]

    E1 --> F[Final Merge Process]
    E2 --> F
    E3 --> F

    F --> G[training_data.json]

    style C1 fill:#e3f2fd
    style C2 fill:#fff3e0
    style C3 fill:#f3e5f5
    style G fill:#e8f5e9
```

**Performance Characteristics**:
- **Sequential Processing**: One chunk at a time by default
- **Parallelization Potential**: Chunks can be processed in parallel containers
- **Time per Chunk**: ~2-5 minutes for 25 examples (depends on model speed)
- **Total Time**: 18 chunks × 3 minutes = ~54 minutes (sequential)
- **Parallel Time**: ~6 minutes with 3 workers

---

## Best Practices and Recommendations

### Chunking Strategy

```mermaid
mindmap
  root((Effective<br/>Chunking))
    Thematic Coherence
      Single topic per chunk
      Related concepts together
      Clear section boundaries
    Size Management
      500-1500 tokens per chunk
      Fits in context window
      Enough content for diversity
    Logical Ordering
      Basic info first
      Complex concepts later
      Dependencies respected
    Configuration
      Adjust nb_entries by complexity
      More iterations for rich sections
      Fewer for simple facts
```

### Document Preparation Checklist

- [ ] **Identify Natural Sections**: Use existing headings and structure
- [ ] **Create Focused Chunks**: Each chunk should cover one aspect (profile, appearance, powers, etc.)
- [ ] **Include Context**: Every chunk gets the same character introduction for consistency
- [ ] **Configure Settings**: Adjust `nb_dataset_entries` and `nb_iterations` based on section richness
- [ ] **Write Clear Prompts**: Templates should guide the model to generate specific question types

### Example Chunking Decision Tree

```mermaid
flowchart TD
    A[Character Section] --> B{Content Type?}
    B -->|Simple Facts| C[Profile, Name, Age]
    B -->|Descriptive| D[Appearance, Clothing]
    B -->|Complex Rules| E[Magic Powers, Abilities]
    B -->|Narrative| F[History, Background]

    C --> G[Low Iteration Count<br/>nb_iterations: 3-4]
    D --> H[Medium Iteration Count<br/>nb_iterations: 4-5]
    E --> I[High Iteration Count<br/>nb_iterations: 5-7]
    F --> I

    G --> J[Simple Questions]
    H --> K[Descriptive Questions]
    I --> L[Complex Questions]

    style C fill:#c8e6c9
    style D fill:#fff9c4
    style E fill:#ffccbc
    style F fill:#ffccbc
```

---

## Troubleshooting Guide

### Common Issues and Solutions

```mermaid
flowchart TD
    A[Problem Detected] --> B{Issue Type?}

    B -->|JSON Parse Failures| C[Increase MAX_GENERATION_RETRIES]
    B -->|Low Diversity| D[Increase nb_iterations<br/>Strengthen diversity prompts]
    B -->|Hallucinations| E[Emphasize document-only grounding<br/>Add validation examples]
    B -->|Repetitive Examples| F[Review prompt templates<br/>Add explicit variation requests]
    B -->|Context Overflow| G[Reduce chunk size<br/>Split complex sections]

    C --> H[Monitor retry logs]
    D --> I[Analyze prompt effectiveness]
    E --> J[Add grounding examples to prompts]
    F --> K[Enhance iteration prompts]
    G --> L[Recalculate chunk boundaries]

    H --> M{Solved?}
    I --> M
    J --> M
    K --> M
    L --> M

    M -->|No| N[Adjust model parameters<br/>or switch model]
    M -->|Yes| O[Document solution]

    style C fill:#ffebee
    style D fill:#fff3e0
    style E fill:#e3f2fd
    style F fill:#f3e5f5
    style G fill:#e8f5e9
```

### Debugging Workflow

1. **Check Logs**: Monitor console output for retry attempts and failures
2. **Inspect Partial Outputs**: Review individual chunk files for quality
3. **Validate JSON Structure**: Use JSON validators on generated files
4. **Review Prompt Templates**: Ensure templates are clear and unambiguous
5. **Test with Single Chunk**: Isolate problematic chunks for focused debugging

---

## Success Metrics

### Quality Indicators

```mermaid
graph LR
    A[Dataset Quality] --> B[Completeness]
    A --> C[Diversity]
    A --> D[Accuracy]
    A --> E[Format Compliance]

    B --> B1[All chunks processed]
    B --> B2[Target entry count met]

    C --> C1[Unique prompts]
    C --> C2[Varied question types]

    D --> D1[Responses match document]
    D --> D2[No hallucinations]

    E --> E1[Valid JSON structure]
    E --> E2[Correct field names]

    style A fill:#e8f5e9
    style B fill:#e3f2fd
    style C fill:#fff3e0
    style D fill:#f3e5f5
    style E fill:#ffebee
```

**Target Metrics**:
- **Total Examples**: 1000+ entries
- **Chunk Coverage**: 100% of character aspects
- **JSON Validity**: >95% on first attempt
- **Diversity Score**: <10% duplicate prompts
- **Accuracy**: 100% responses grounded in source document

---

## Conclusion

This chunked iterative generation strategy enables **cost-effective, high-quality dataset creation using small language models**. By breaking down complex character descriptions, iterating for diversity, and maintaining strict quality controls, the system produces training data suitable for fine-tuning NPC models.

**Key Advantages**:
- ✅ Works with small, resource-efficient models
- ✅ Produces diverse, high-quality examples
- ✅ Handles JSON generation failures gracefully
- ✅ Scales horizontally for faster processing
- ✅ Maintains character consistency across all examples

**Ideal Use Cases**:
- Fine-tuning small models for game NPCs
- Creating character-specific chatbots
- Training models with limited GPU resources
- Educational projects for learning LLM fine-tuning
- Cost-sensitive production environments

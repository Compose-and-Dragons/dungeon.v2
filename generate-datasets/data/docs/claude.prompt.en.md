I need to fine-tune a small LLM (qwen2.5:0.5b-instruct) using data in the following format:

```json
[
  {
    "prompt": "question",
    "response": "answer"
  },
  {
    "prompt": "question2",
    "response": "answer2"
  },
  // etc...
]
```


- This data will be saved in a training_data.json file.
- The data in this file will be generated (in English) from the information contained in the attached queen_pedauque_character_sheet.md document.
- There must be enough examples (ideally 1000+) for the fine-tuning to be effective.
- Each example must be unique and not repeated.
- Each example must be relevant to the character described in queen_pedauque_character_sheet.md.
- Each example must include varied questions covering different aspects of the character.

We will proceed in several steps (the training_data.json file will be generated in multiple passes):
Note: for each step, always maintain a global view of Queen Pédauque's character context.

### First Step:

Generate 100 examples of question-answer pairs in the JSON format described above, using only the information available in the queen_pedauque_character_sheet.md document about Queen Pédauque character, specifically using information from the beginning of the document up to and including the "Character Profile" section.

The LLM should be able to answer questions such as:
- What is your name?
- Who are you?
- What is your class?
- What is your role?
- etc...

The LLM should speak in the first person as if it were the character.

### Second Step:

Use the information from the "Behavioral Instructions" and "Interaction Guidelines" sections of the queen_pedauque_character_sheet.md document to generate 100 additional examples of question-answer pairs in the JSON format described above. These examples should reflect the behavior and interaction directives specified in these sections.

### Third Step:

Use the information from the "Special Triggers" section of the queen_pedauque_character_sheet.md document to generate 100 additional examples of question-answer pairs in the JSON format described above. These examples must include questions that trigger the special responses mentioned in this section. It is important to ensure the LLM has a specific reaction to these triggers: "Chocolatine" and "Pain au chocolat".

### Fourth Step:

In the main section "Queen Pédauque - Background and Personality", use the subsections "Kind", "Name and Title" and "Age" to generate 100 additional examples of question-answer pairs in the JSON format described above.

### Fifth Step:

In the main section "Queen Pédauque - Background and Personality", use the "Family" subsection to generate 100 additional examples of question-answer pairs in the JSON format described above about the main character's family. Questions should be varied and cover different aspects of the character's family.

### Sixth Step:

In the main section "Queen Pédauque - Background and Personality", use the "Occupation" subsection to generate 100 additional examples of question-answer pairs in the JSON format described above about the main character's occupations.

### Seventh Step:

In the main section "Queen Pédauque - Background and Personality", use the "Physical Appearance" subsection to generate 100 additional examples of question-answer pairs in the JSON format described above about the main character's physical appearance.
Questions should be varied and cover different aspects of physical appearance.

### Eighth Step:

In the main section "Queen Pédauque - Background and Personality", use the "Clothing" subsection to generate 100 additional examples of question-answer pairs in the JSON format described above about how the main character dresses.
Questions should be varied and cover different aspects of clothing.

### Ninth Step:

In the main section "Queen Pédauque - Background and Personality", use the "Food Preferences" subsection to generate 100 additional examples of question-answer pairs in the JSON format described above about the main character's food preferences.
Questions should be varied and cover different aspects of their food preferences.

### Tenth Step:

In the main section "Queen Pédauque - Background and Personality", use the "Background Story" subsection along with the "The Invention of the Chocolatine" and "Present Day" subsections to generate 100 additional examples of question-answer pairs in the JSON format described above about the main character's history.
Questions should be varied and cover different aspects of their history.

### Eleventh Step:

In the main section "Queen Pédauque - Background and Personality", use the "Personality and Character Traits" subsection along with the associated subsections "Strengths", "Weaknesses", "Distinctive Traits" to generate 300 additional examples of question-answer pairs in the JSON format described above about the main character's personality.
Questions should be varied and cover different aspects of their personality and character.

### Twelfth Step:

In the main section "Queen Pédauque - Background and Personality", use the "Magical Powers" subsection along with the associated subsections "Culinary Sorcery", "Transformation Magic", "Defensive & Combat Magic", "Blessing & Curse Magic" and "Regional Magic" to generate 300 additional examples of question-answer pairs in the JSON format described above about the main character's magical powers.
Questions should be varied and cover different aspects of their magical powers.

### Thirteenth Step:

In the main section "Queen Pédauque - Background and Personality", use the "Quote" subsection to generate 100 additional examples of question-answer pairs in the JSON format described above about the main character's favorite quote.

### Fourteenth Step:

In the main section "Queen Pédauque - Background and Personality", use the "Secret Keyword" subsection to generate 100 additional examples of question-answer pairs in the JSON format described above about the main character's secret word.

### Fifteenth Step:

Use the entire main section "The Sacred Encyclopedia of Chocolatine" to generate 500 additional examples of question-answer pairs in the JSON format described above. These examples should reflect a global understanding of what the chocolatine is, integrating information from all subsections of the main section.

Make sure the final training_data.json file is a valid JSON file containing all the question-answer pairs generated during the different steps.

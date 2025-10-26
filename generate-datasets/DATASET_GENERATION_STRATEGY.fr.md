# Stratégie de Génération de Données pour le Fine-Tuning de Petits Modèles de Langage

## Résumé Exécutif

Ce document décrit une **stratégie de génération itérative par fragments** conçue spécifiquement pour le fine-tuning de petits modèles de langage (SLM) avec des fenêtres de contexte limitées. L'approche divise les grandes descriptions de personnages en fragments gérables et utilise un processus de raffinement itératif pour générer des données d'entraînement diversifiées et de haute qualité.

**Innovation Clé** : Au lieu de s'appuyer sur de grands modèles comme GPT-4 ou Claude pour générer des milliers d'exemples en une seule passe, cette stratégie utilise de petits modèles (comme jan-nano-128k ou qwen2.5:0.5b) pour générer des datasets de manière incrémentale, rendant le processus économique et accessible.

---

## Énoncé du Problème

### Défis avec les Petits Modèles

Lors du fine-tuning de petits modèles de langage pour la génération de PNJ (Personnages Non-Joueurs) dans les jeux, nous faisons face à plusieurs défis :

1. **Fenêtre de Contexte Limitée** : Les petits modèles ont généralement des fenêtres de contexte de 8K-16K tokens, les empêchant de traiter de grandes fiches de personnages en une seule passe
2. **Contraintes de Coût** : Utiliser de grands modèles commerciaux (GPT-4, Claude) pour la génération de datasets est coûteux
3. **Exigences de Qualité** : Besoin de 1000+ exemples diversifiés et de haute qualité pour un fine-tuning efficace
4. **Cohérence** : Doit maintenir la cohérence du personnage à travers tous les exemples générés
5. **Fiabilité JSON** : Les petits modèles ont souvent du mal avec un format de sortie JSON cohérent

---

## Architecture de la Solution

### Stratégie de Haut Niveau

```mermaid
graph TB
    A[Grand Document Personnage] --> B[Découpage Manuel par Section]
    B --> C[Fragment 1: Profil]
    B --> D[Fragment 2: Comportement]
    B --> E[Fragment 3: Apparence]
    B --> F[Fragment N: ...]

    C --> G[Processus de Génération Itératif]
    D --> G
    E --> G
    F --> G

    G --> H[Fichiers JSON Individuels]
    H --> I[Fusion de Tous les Fichiers]
    I --> J[training_data.json Final]

    style A fill:#e1f5ff
    style G fill:#fff4e1
    style J fill:#e8f5e9
```

### Composants Principaux

Le système se compose de trois composants principaux :

1. **Préparation du Document** : Découpage manuel des fiches de personnage en sections ciblées
2. **Moteur de Génération Itérative** : Génération multi-passes avec raffinement
3. **Pipeline de Fusion** : Agrégation de tous les datasets générés

---

## Flux de Processus Détaillé

### 1. Phase de Préparation du Document

```mermaid
flowchart LR
    A[Fiche Personnage<br/>queen_pedauque_character_sheet.md] --> B{Analyse Manuelle}
    B --> C[Identifier les Sections Logiques]
    C --> D[Profil]
    C --> E[Comportement]
    C --> F[Apparence]
    C --> G[Pouvoirs Magiques]
    C --> H[Histoire]
    C --> I[Règles Spéciales]

    D --> J[Créer Fichiers Fragments]
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

Chaque fichier fragment contient **4 sections séparées par le délimiteur `----------`** :

1. **Contexte** : Instructions système brèves pour le personnage
2. **Document** : Le contenu réel à partir duquel générer des exemples
3. **Paramètres** : Configuration JSON pour les paramètres de génération
4. **Template de Prompt** : Modèle pour générer les questions

**Exemple de Structure de Fragment** :
```markdown
# Queen Pédauque - Instructions Système PNJ Sorcière Légendaire
Vous êtes Queen Pédauque, une sorcière légendaire PNJ...
----------
## Profil du Personnage
- **Nom** : Queen Pédauque (aussi connue comme la Reine aux Pieds d'Oie)
- **Race** : Humaine avec Héritage Féérique Ancien
...
----------
{
    "nb_dataset_entries": 20,
    "nb_iterations": 5
}
----------
À partir de ce document lié à {{.NameOfTheNPC}} :
{{.Chunk}}
Générer {{.NbEntriesPerChunk}} entrées de dataset...
```

---

### 2. Moteur de Génération Itérative

```mermaid
sequenceDiagram
    participant M as Programme Principal
    participant C as Fichier Fragment
    participant T as Moteur Template
    participant A as NPCAgent
    participant L as Petit LLM
    participant F as Système Fichiers

    M->>C: Lire fichier fragment
    C-->>M: Contexte, Document, Paramètres, Template

    Note over M: Parser JSON Paramètres<br/>(nb_dataset_entries, nb_iterations)

    M->>T: Interpoler template avec données
    T-->>M: Prompt formaté

    rect rgb(220, 240, 255)
        Note over M,L: ITÉRATION 1
        M->>A: JsonCompletion(prompt)
        A->>L: Générer 5 exemples
        L-->>A: Tableau JSON
        A-->>M: 5 objets DataSetEntry
        M->>M: Ajouter à trainingData
    end

    rect rgb(255, 240, 220)
        Note over M,L: ITÉRATIONS 2-5
        loop 4 fois de plus
            M->>M: Modifier prompt avec réponse précédente
            M->>A: JsonCompletion(prompt amélioré)
            A->>L: Générer 5 NOUVEAUX exemples
            L-->>A: Tableau JSON
            A-->>M: 5 objets DataSetEntry
            M->>M: Ajouter à trainingData
        end
    end

    M->>F: Écrire 1.queen_pedauque_training_data.json
    F-->>M: Succès
```

---

### 3. Boucle de Génération avec Logique de Retry

```mermaid
flowchart TD
    A[Démarrer Génération pour Fragment] --> B[Lire Paramètres Fragment]
    B --> C{Parser Paramètres}
    C -->|Succès| D[nb_dataset_entries = 5<br/>nb_iterations = 5]
    C -->|Échec| Z[Erreur Fatale]

    D --> E[Compteur Itération = 1]
    E --> F[Interpoler Template Prompt]
    F --> G[Appeler LLM JsonCompletion]

    G --> H{JSON Valide?}
    H -->|Non| I{Retry < 3?}
    I -->|Oui| J[Compteur Retry++]
    J --> G
    I -->|Non| K[Passer Cette Itération]

    H -->|Oui| L[Parser JSON vers tableau DataSetEntry]
    L --> M[Ajouter 5 entrées à trainingData]
    M --> N{Itération < 5?}

    N -->|Oui| O[Compteur Itération++]
    O --> P[Modifier Prompt:<br/>Inclure réponse précédente<br/>Demander questions DIFFÉRENTES]
    P --> G

    N -->|Non| Q[Total Entrées = 25]
    Q --> R[Marshaller vers JSON avec Indentation]
    R --> S[Écrire fichier:<br/>fileIndex.trainDataFile]
    S --> T[fileIndex++]
    T --> U{Plus de Fragments?}

    U -->|Oui| A
    U -->|Non| V[Fusionner Tous les Fichiers]
    V --> W[Fin]

    K --> N

    style G fill:#fff4e1
    style H fill:#ffebee
    style L fill:#e8f5e9
    style V fill:#e1f5ff
```

**Points Clés** :
- Chaque fragment génère `nb_dataset_entries × nb_iterations` exemples (ex: 5 × 5 = 25)
- Chaque itération demande explicitement des questions **DIFFÉRENTES** pour assurer la diversité
- Mécanisme de retry (jusqu'à 3 tentatives) gère les échecs de parsing JSON
- Les itérations échouées sont sautées plutôt que de crasher tout le processus

---

### 4. Stratégie de Raffinement Itératif

Le système utilise une technique astucieuse de modification de prompt pour assurer la diversité :

```mermaid
graph TD
    A[Itération 1:<br/>Prompt de Base] --> B[LLM génère<br/>5 exemples]
    B --> C[Sauvegarder Réponse 1]

    C --> D[Itération 2:<br/>Prompt Amélioré]
    D --> E["Prompt inclut:<br/>• Réponse précédente<br/>• Demande explicite de NOUVELLES questions<br/>• Même contexte document"]
    E --> F[LLM génère<br/>5 exemples DIFFÉRENTS]
    F --> G[Sauvegarder Réponse 2]

    G --> H[Itération 3-5:<br/>Continuer le pattern]
    H --> I[Total: 25 exemples divers<br/>d'un seul fragment]

    style A fill:#e3f2fd
    style D fill:#fff3e0
    style H fill:#fce4ec
    style I fill:#e8f5e9
```

**Exemple d'Évolution de Prompt** :

**Itération 1** :
```
À partir de ce document lié à Queen Pédauque :
[Contenu document...]
Générer 5 entrées de dataset...
```

**Itération 2** :
```
Voici la réponse précédente que vous avez donnée :
[Tableau JSON précédent...]

Maintenant, veuillez générer de NOUVELLES entrées de dataset pour le même document
mais avec des prompts DIFFÉRENTS :
[Même contenu document...]
Générer 5 entrées de dataset...
```

Cette approche :
- ✅ Maintient la conscience du contexte
- ✅ Encourage la diversité en montrant ce qui a déjà été généré
- ✅ Garde les prompts ciblés et gérables pour les petits modèles

---

## Organisation des Fichiers et Fusion

### Génération de Fichiers Individuels

```mermaid
graph LR
    A[Fragment 01.profile.md] --> B[1.queen_pedauque_training_data.json]
    C[Fragment 02.behavior.md] --> D[2.queen_pedauque_training_data.json]
    E[Fragment 03.interactions.md] --> F[3.queen_pedauque_training_data.json]
    G[Fragment 18.cultural-significance.md] --> H[18.queen_pedauque_training_data.json]

    style B fill:#e8f5e9
    style D fill:#e8f5e9
    style F fill:#e8f5e9
    style H fill:#e8f5e9
```

### Processus de Fusion

```mermaid
flowchart TD
    A[Démarrer Fusion] --> B[Obtenir tous les fichiers du répertoire]
    B --> C{Pour chaque fichier}
    C --> D{Se termine par<br/>trainDataFile?}
    D -->|Non| E[Passer fichier]
    D -->|Oui| F[Lire fichier JSON]
    F --> G{Tableau JSON<br/>valide?}
    G -->|Non| H[Logger avertissement, continuer]
    G -->|Oui| I[Parser vers tableau DataSetEntry]
    I --> J[Ajouter toutes les entrées à mergedData]

    E --> K{Plus de fichiers?}
    H --> K
    J --> K
    K -->|Oui| C
    K -->|Non| L[Marshaller tableau fusionné vers JSON]
    L --> M[Écrire training_trainDataFile]
    M --> N[Rapporter total entrées]
    N --> O[Fin]

    style F fill:#e3f2fd
    style I fill:#fff3e0
    style M fill:#e8f5e9
```

**Résultat** :
- Entrée : `1.queen_pedauque_training_data.json`, `2.queen_pedauque_training_data.json`, ... `18.queen_pedauque_training_data.json`
- Sortie : `training_queen_pedauque_training_data.json` (fichier unique fusionné avec 450+ exemples)

---

## Gestion de la Configuration

### Configuration Docker Compose

Le système utilise Docker Compose pour des environnements reproductibles :

```yaml
environment:
  MAX_GENERATION_RETRIES: 3
  NPC_NAME: Queen Pedauque
  SYSTEM_INSTRUCTIONS: |
    Vous êtes un assistant IA utile...
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

### Guide d'Ajustement des Paramètres

```mermaid
graph TD
    A[Paramètres de Configuration] --> B[Paramètres de Génération]
    A --> C[Paramètres de Modèle]
    A --> D[Paramètres de Retry]

    B --> B1[nb_dataset_entries<br/>Compte par itération]
    B --> B2[nb_iterations<br/>Passes de raffinement]

    C --> C1[MODEL_TEMPERATURE<br/>Niveau de créativité]
    C --> C2[MODEL_TOP_P<br/>Diversité d'échantillonnage]
    C --> C3[context_size<br/>Limite de tokens]

    D --> D1[MAX_GENERATION_RETRIES<br/>Tolérance aux échecs]

    style B fill:#e3f2fd
    style C fill:#fff3e0
    style D fill:#ffebee
```

**Paramètres Recommandés** :

| Paramètre | Valeur Recommandée | Justification |
|-----------|-------------------|---------------|
| `nb_dataset_entries` | 5 | Gérable pour petits modèles, équilibre qualité vs quantité |
| `nb_iterations` | 5 | Fournit diversité sans redondance excessive |
| `MODEL_TEMPERATURE` | 1.0 | Créativité élevée pour questions variées |
| `MODEL_TOP_P` | 0.9 | Échantillonnage équilibré pour sorties cohérentes mais diverses |
| `MAX_GENERATION_RETRIES` | 3 | Tolère échecs JSON occasionnels sans boucles infinies |

---

## Assurance Qualité des Données

### Mécanismes de Contrôle Qualité

```mermaid
flowchart LR
    A[Génération Brute] --> B[Validation Schéma JSON]
    B --> C{Structure Valide?}
    C -->|Non| D[Réessayer jusqu'à 3 fois]
    C -->|Oui| E[Validation Contenu]

    E --> F{Contient champs requis?}
    F -->|Non| D
    F -->|Oui| G[Vérification Diversité]

    G --> H[Comparer avec itérations précédentes]
    H --> I{Trop similaire?}
    I -->|Oui| J[Prompt insiste sur NOUVELLES questions]
    I -->|Non| K[Accepter entrée]

    D --> L{Tentatives épuisées?}
    L -->|Oui| M[Passer itération]
    L -->|Non| B

    K --> N[Ajouter au dataset d'entraînement]
    M --> O[Logger avertissement]

    style B fill:#e3f2fd
    style E fill:#fff3e0
    style G fill:#f3e5f5
    style K fill:#e8f5e9
    style M fill:#ffebee
```

### Validation du Format de Sortie

Chaque entrée générée doit se conformer à :

```json
{
  "prompt": "chaîne question (non vide)",
  "response": "chaîne réponse (non vide, basée uniquement sur document fourni)"
}
```

**Règles de Validation** :
1. Doit être un tableau JSON valide
2. Chaque élément doit avoir exactement 2 champs : `prompt` et `response`
3. Les deux champs doivent être des chaînes non vides
4. Les réponses doivent être ancrées dans le document source (pas d'hallucinations)
5. Les prompts doivent être diversifiés à travers les itérations

---

## Scalabilité et Performance

### Stratégie de Scaling Horizontal

```mermaid
graph TB
    A[Fiche Personnage] --> B[Découpage Manuel]
    B --> C1[Groupe Fragments 1<br/>Fragments 1-6]
    B --> C2[Groupe Fragments 2<br/>Fragments 7-12]
    B --> C3[Groupe Fragments 3<br/>Fragments 13-18]

    C1 --> D1[Conteneur Worker 1]
    C2 --> D2[Conteneur Worker 2]
    C3 --> D3[Conteneur Worker 3]

    D1 --> E1[Fichiers 1-6]
    D2 --> E2[Fichiers 7-12]
    D3 --> E3[Fichiers 13-18]

    E1 --> F[Processus Fusion Finale]
    E2 --> F
    E3 --> F

    F --> G[training_data.json]

    style C1 fill:#e3f2fd
    style C2 fill:#fff3e0
    style C3 fill:#f3e5f5
    style G fill:#e8f5e9
```

**Caractéristiques de Performance** :
- **Traitement Séquentiel** : Un fragment à la fois par défaut
- **Potentiel de Parallélisation** : Les fragments peuvent être traités en conteneurs parallèles
- **Temps par Fragment** : ~2-5 minutes pour 25 exemples (dépend de la vitesse du modèle)
- **Temps Total** : 18 fragments × 3 minutes = ~54 minutes (séquentiel)
- **Temps Parallèle** : ~6 minutes avec 3 workers

---

## Meilleures Pratiques et Recommandations

### Stratégie de Découpage

```mermaid
mindmap
  root((Découpage<br/>Efficace))
    Cohérence Thématique
      Un seul sujet par fragment
      Concepts liés ensemble
      Limites de section claires
    Gestion Taille
      500-1500 tokens par fragment
      Tient dans fenêtre contexte
      Assez de contenu pour diversité
    Ordre Logique
      Info de base d'abord
      Concepts complexes après
      Dépendances respectées
    Configuration
      Ajuster nb_entries par complexité
      Plus d'itérations pour sections riches
      Moins pour faits simples
```

### Checklist Préparation Document

- [ ] **Identifier Sections Naturelles** : Utiliser en-têtes et structure existants
- [ ] **Créer Fragments Ciblés** : Chaque fragment devrait couvrir un aspect (profil, apparence, pouvoirs, etc.)
- [ ] **Inclure Contexte** : Chaque fragment reçoit la même introduction de personnage pour cohérence
- [ ] **Configurer Paramètres** : Ajuster `nb_dataset_entries` et `nb_iterations` selon richesse section
- [ ] **Écrire Prompts Clairs** : Les templates devraient guider le modèle à générer types de questions spécifiques

### Arbre de Décision de Découpage

```mermaid
flowchart TD
    A[Section Personnage] --> B{Type de Contenu?}
    B -->|Faits Simples| C[Profil, Nom, Âge]
    B -->|Descriptif| D[Apparence, Vêtements]
    B -->|Règles Complexes| E[Pouvoirs Magiques, Capacités]
    B -->|Narratif| F[Histoire, Background]

    C --> G[Faible Compte Itération<br/>nb_iterations: 3-4]
    D --> H[Moyen Compte Itération<br/>nb_iterations: 4-5]
    E --> I[Haut Compte Itération<br/>nb_iterations: 5-7]
    F --> I

    G --> J[Questions Simples]
    H --> K[Questions Descriptives]
    I --> L[Questions Complexes]

    style C fill:#c8e6c9
    style D fill:#fff9c4
    style E fill:#ffccbc
    style F fill:#ffccbc
```

---

## Guide de Dépannage

### Problèmes Courants et Solutions

```mermaid
flowchart TD
    A[Problème Détecté] --> B{Type de Problème?}

    B -->|Échecs Parse JSON| C[Augmenter MAX_GENERATION_RETRIES]
    B -->|Faible Diversité| D[Augmenter nb_iterations<br/>Renforcer prompts diversité]
    B -->|Hallucinations| E[Insister sur ancrage document<br/>Ajouter exemples validation]
    B -->|Exemples Répétitifs| F[Revoir templates prompts<br/>Ajouter demandes variation explicites]
    B -->|Dépassement Contexte| G[Réduire taille fragments<br/>Diviser sections complexes]

    C --> H[Monitorer logs retry]
    D --> I[Analyser efficacité prompts]
    E --> J[Ajouter exemples ancrage aux prompts]
    F --> K[Améliorer prompts itération]
    G --> L[Recalculer limites fragments]

    H --> M{Résolu?}
    I --> M
    J --> M
    K --> M
    L --> M

    M -->|Non| N[Ajuster paramètres modèle<br/>ou changer modèle]
    M -->|Oui| O[Documenter solution]

    style C fill:#ffebee
    style D fill:#fff3e0
    style E fill:#e3f2fd
    style F fill:#f3e5f5
    style G fill:#e8f5e9
```

### Workflow de Débogage

1. **Vérifier Logs** : Monitorer sortie console pour tentatives retry et échecs
2. **Inspecter Sorties Partielles** : Revoir fichiers fragments individuels pour qualité
3. **Valider Structure JSON** : Utiliser validateurs JSON sur fichiers générés
4. **Revoir Templates Prompts** : S'assurer que les templates sont clairs et non ambigus
5. **Tester avec Fragment Unique** : Isoler fragments problématiques pour débogage ciblé

---

## Métriques de Succès

### Indicateurs de Qualité

```mermaid
graph LR
    A[Qualité Dataset] --> B[Complétude]
    A --> C[Diversité]
    A --> D[Exactitude]
    A --> E[Conformité Format]

    B --> B1[Tous fragments traités]
    B --> B2[Compte cible entrées atteint]

    C --> C1[Prompts uniques]
    C --> C2[Types questions variés]

    D --> D1[Réponses correspondent document]
    D --> D2[Pas d'hallucinations]

    E --> E1[Structure JSON valide]
    E --> E2[Noms champs corrects]

    style A fill:#e8f5e9
    style B fill:#e3f2fd
    style C fill:#fff3e0
    style D fill:#f3e5f5
    style E fill:#ffebee
```

**Métriques Cibles** :
- **Total Exemples** : 1000+ entrées
- **Couverture Fragments** : 100% des aspects personnage
- **Validité JSON** : >95% à la première tentative
- **Score Diversité** : <10% de prompts dupliqués
- **Exactitude** : 100% réponses ancrées dans document source

---

## Conclusion

Cette stratégie de génération itérative par fragments permet la **création de datasets économique et de haute qualité en utilisant de petits modèles de langage**. En décomposant les descriptions de personnages complexes, en itérant pour la diversité, et en maintenant des contrôles qualité stricts, le système produit des données d'entraînement adaptées au fine-tuning de modèles PNJ.

**Avantages Clés** :
- ✅ Fonctionne avec des modèles petits et économes en ressources
- ✅ Produit des exemples diversifiés et de haute qualité
- ✅ Gère les échecs de génération JSON avec élégance
- ✅ Scale horizontalement pour traitement plus rapide
- ✅ Maintient cohérence personnage à travers tous les exemples

**Cas d'Usage Idéaux** :
- Fine-tuning de petits modèles pour PNJ de jeux
- Création de chatbots spécifiques à un personnage
- Entraînement de modèles avec ressources GPU limitées
- Projets éducatifs pour apprendre le fine-tuning LLM
- Environnements de production sensibles aux coûts

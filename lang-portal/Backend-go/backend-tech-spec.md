Here’s your corrected and refined version with improved spelling, grammar, and clarity:

---

# **Backend Server Technical Specifications**

## **Business Goal**

A language learning school wants to build a prototype learning portal that will serve three primary functions:

- **Vocabulary Inventory:** A repository of possible vocabulary that can be learned.
- **Learning Record Store (LRS):** A system to track correct and incorrect vocabulary practice scores.
- **Unified Launchpad:** A centralized interface to access various learning applications.

## **Technical Requirements**

- The backend will be built using **Go**.
- The database will be **SQLite3**.
- The API will be built using **Gin**.
- **Mage** will be used as a task runner for Go.
- The API will always return **JSON**.
- **No authentication or authorization** will be implemented.
- **All interactions** will be treated as a single-user system.

## **Directory Structure**

```text
backend_go/
├── cmd/
│   └── server/
├── internal/
│   ├── models/     # Data structures and database operations
│   ├── handlers/   # HTTP handlers organized by feature (dashboard, words, groups, etc.)
│   └── service/    # Business logic
├── db/
│   ├── migrations/
│   └── seeds/      # For initial data population
├── magefile.go
├── go.mod
└── words.db
```

## **Database Schema**

The system will use a single **SQLite** database named `words.db`, located in the root of the `backend_go` project folder.

### **Tables and Schema**

- **`words`** – Stores vocabulary words.
  - `id` (integer, primary key)
  - `japanese` (string)
  - `romaji` (string)
  - `english` (string)
  - `parts` (JSON)

- **`words_groups`** – Join table for words and groups (many-to-many relationship).
  - `id` (integer, primary key)
  - `word_id` (integer, foreign key)
  - `group_id` (integer, foreign key)

- **`groups`** – Thematic groups of words.
  - `id` (integer, primary key)
  - `name` (string)

- **`study_sessions`** – Records of study sessions, grouping `word_review_items`.
  - `id` (integer, primary key)
  - `group_id` (integer, foreign key)
  - `created_at` (datetime)
  - `study_activity_id` (integer, foreign key)

- **`study_activities`** – Specific study activities linking a study session to a group.
  - `id` (integer, primary key)
  - `study_session_id` (integer, foreign key)
  - `group_id` (integer, foreign key)
  - `created_at` (datetime)

- **`word_review_items`** – Records of word practice, tracking correctness.
  - `word_id` (integer, foreign key)
  - `study_session_id` (integer, foreign key)
  - `correct` (boolean)
  - `created_at` (datetime)

---

## **API Endpoints**

### **Dashboard API**
- `GET /api/dashboard/last_study_session`  
  _Retrieves details of the last study session._

- `GET /api/dashboard/study_progress`  
  _Returns overall study progress._

- `GET /api/dashboard/quick_stats`  
  _Provides a summary of key statistics._

### **Study Activities API**
- `GET /api/study_activities/:id`  
  _Retrieves details of a specific study activity._

- `GET /api/study_activities/:id/study_sessions`  
  _Returns all study sessions related to a specific activity._

- `POST /api/study_activities`  
  _Creates a new study activity._  
  **Required Parameters:**  
  - `group_id` (integer) – The ID of the group associated with the activity.  
  - `study_activity_id` (integer) – The ID of the study activity.

### **Words API**
- `GET /api/words`  
  _Retrieves a paginated list of all words (100 items per page)._

- `GET /api/words/:id`  
  _Retrieves details of a specific word._

### **Word Groups API**
- `GET /api/groups`  
  _Retrieves a paginated list of word groups (100 items per page)._

- `GET /api/groups/:id`  
  _Retrieves details of a specific word group._

- `GET /api/groups/:id/words`  
  _Returns a list of words associated with a specific group._

- `GET /api/groups/:id/study_sessions`  
  _Returns study sessions related to a specific group._

### **Study Sessions API**
- `GET /api/study_sessions`  
  _Retrieves a paginated list of all study sessions (100 items per page)._

- `GET /api/study_sessions/:id`  
  _Retrieves details of a specific study session._

- `GET /api/study_sessions/:id/words`  
  _Returns all words reviewed in a specific study session._

- `POST /api/study_sessions/:id/words/:word_id/review`  
  _Records a word review result._  
  **Required Parameter:**  
  - `correct` (boolean) – Indicates whether the word was answered correctly.

### **Reset API**
- `POST /api/reset_history`  
  _Deletes all study sessions and word review items._

- `POST /api/full_reset`  
  _Drops all tables and recreates them with seed data._
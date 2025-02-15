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

# **API Endpoints Documentation**

## **Dashboard Endpoints**

### **GET /api/dashboard/last_study_session**  
**Description:** Retrieves details of the last study session.  

**Response:**
```json
{
  "activity_id": "string",
  "activity_name": "string",
  "timestamp": "2024-03-21T15:30:00Z",
  "stats": {
    "correct_answers": 15,
    "incorrect_answers": 5
  },
  "group_id": "string",
  "group_name": "string"
}
```

---

### **GET /api/dashboard/study_progress**  
**Description:** Retrieves the overall study progress.  

**Response:**
```json
{
  "total_words_studied": 3,
  "total_available_words": 124,
  "mastery_percentage": 45.5
}
```

---

### **GET /api/dashboard/quick_stats**  
**Description:** Provides a summary of key statistics.  

**Response:**
```json
{
  "success_rate": 80.5,
  "total_study_sessions": 4,
  "active_groups": 3,
  "study_streak_days": 4
}
```

---

## **Study Activities Endpoints**

### **GET /api/study_activities**  
**Description:** Retrieves a list of all available study activities.  

**Response:**
```json
{
  "activities": [
    {
      "id": "string",
      "name": "string",
      "thumbnail_url": "string",
      "description": "string"
    }
  ]
}
```

---

### **GET /api/study_activities/:id**  
**Description:** Retrieves details of a specific study activity.  

**Response:**
```json
{
  "id": "string",
  "name": "string",
  "thumbnail_url": "string",
  "description": "string",
  "launch_url": "string"
}
```

---

### **GET /api/study_activities/:id/study_sessions**  
**Description:** Retrieves all study sessions associated with a specific study activity.  

**Response:**
```json
{
  "items": [
    {
      "id": "string",
      "activity_name": "string",
      "group_name": "string",
      "start_time": "2024-03-21T15:30:00Z",
      "end_time": "2024-03-21T16:00:00Z",
      "review_items_count": 20
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_items": 100,
    "items_per_page": 20
  }
}
```

---

### **POST /api/study_activities**  
**Description:** Creates a new study activity.  

**Request Body:**
```json
{
  "activity_id": "string",
  "group_id": "string"
}
```

**Response:**
```json
{
  "study_session_id": "string",
  "activity_url": "string"
}
```

---

## **Words Endpoints**

### **GET /api/words**  
**Description:** Retrieves a paginated list of all words (100 items per page).  

**Response:**
```json
{
  "items": [
    {
      "id": "string",
      "japanese": "string",
      "romaji": "string",
      "english": "string",
      "correct_count": 10,
      "wrong_count": 2
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_items": 500,
    "items_per_page": 100
  }
}
```

---

### **GET /api/words/:id**  
**Description:** Retrieves details of a specific word.  

**Response:**
```json
{
  "id": "string",
  "japanese": "string",
  "romaji": "string",
  "english": "string",
  "stats": {
    "correct_count": 10,
    "wrong_count": 2
  },
  "groups": [
    {
      "id": "string",
      "name": "string"
    }
  ]
}
```

---

## **Groups Endpoints**

### **GET /api/groups**  
**Description:** Retrieves a paginated list of word groups (100 items per page).  

**Response:**
```json
{
  "items": [
    {
      "id": "string",
      "name": "string",
      "word_count": 50
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 3,
    "total_items": 50,
    "items_per_page": 20
  }
}
```

---

### **GET /api/groups/:id**  
**Description:** Retrieves details of a specific word group.  

**Response:**
```json
{
  "id": "string",
  "name": "string",
  "statistics": {
    "total_words": 50
  }
}
```

---

### **GET /api/groups/:id/words**  
**Description:** Retrieves a paginated list of words associated with a specific group.  

**Response:**
```json
{
  "items": [
    {
      "id": "string",
      "japanese": "string",
      "romaji": "string",
      "english": "string",
      "correct_count": 10,
      "wrong_count": 2
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_items": 50,
    "items_per_page": 10
  }
}
```

---

### **GET /api/groups/:id/study_sessions**  
**Description:** Retrieves all study sessions associated with a specific group.  

**Response:**
```json
{
  "study_sessions": [
    {
      "id": "string",
      "activity_name": "string",
      "start_time": "2024-03-21T15:30:00Z",
      "end_time": "2024-03-21T16:00:00Z",
      "review_items_count": 20
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 3,
    "total_items": 30,
    "items_per_page": 10
  }
}
```

---

## **Study Sessions Endpoints**

### **GET /api/study_sessions**  
**Description:** Retrieves a paginated list of all study sessions.  

**Response:**
```json
{
  "study_sessions": [
    {
      "id": "string",
      "activity_name": "string",
      "group_name": "string",
      "start_time": "2024-03-21T15:30:00Z",
      "end_time": "2024-03-21T16:00:00Z",
      "review_items_count": 20
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_items": 100,
    "items_per_page": 20
  }
}
```

---

### **GET /api/study_sessions/:id**  
**Description:** Retrieves details of a specific study session.  

**Response:**
```json
{
  "id": "string",
  "activity_name": "string",
  "group_name": "string",
  "start_time": "2024-03-21T15:30:00Z",
  "end_time": "2024-03-21T16:00:00Z",
  "review_items_count": 20
}
```

---

## **Settings Endpoints**

### **POST /api/reset_history**  
**Description:** Deletes all study sessions and word review items.  

**Response:**
```json
{
  "success": true,
  "message": "Study history has been reset successfully"
}
```

---

### **POST /api/full_reset**  
**Description:** Drops all tables and reseeds the system with initial data.  

**Response:**
```json
{
  "success": true,
  "message": "System has been fully reset and seeded with initial data"
}
```

---

### **POST /api/study_sessions/:id/words/:word_id/review**
**Description:** Records the review result for a single word within a study session.

**URL Parameters:**
- `id` (study_session_id) - integer
- `word_id` - integer

**Request Body:**
```json
{
  "correct": boolean
}
```

**Response:**
```json
{
  "success": true,
  "word_id": 1,
  "study_session_id": 123,
  "correct": true,
  "created_at": "2025-02-08T17:33:07-05:00"

```

---

## Tasks

Lets list out possible tasks we need for our lang portal.

### Initialize Database
This task will initialize the sqlite database called `words.db`. Located in the project dir of `backend-go`

### Migrate Database
TThis task will run a series of migrations sql files on the database called `words.db`

Migrations live in the `migrations` folder.
The migration files will be run in order of their file name.
The file names should looks like this:

```sql
0001_init.sql
0002_create_words_table.sql
```

### Seed Data
This task will import json files and transform them into target data for our database.

All seed files live in the `seeds` folder.
All seed files must be loaded.

In our task we should have DSL to specific each seed file and its expected group word name.

```json
[
  {
    "kanji": "払う",
    "romaji": "harau",
    "english": "to pay",
  },
  ...
]
```
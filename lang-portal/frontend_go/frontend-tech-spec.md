
# **Frontend Technical Specifications**

## **Business Goal**

A language learning school wants to build a prototype learning portal that will serve three primary functions:

- **Vocabulary Inventory:** A repository of possible vocabulary that can be learned.
- **Learning Record Store (LRS):** A system to track correct and incorrect vocabulary practice scores.
- **Unified Launchpad:** A centralized interface to access various learning applications.

You have been tasked with creating the frontend API of the application.
The fractional CTO has made strong recommendation that you settle on a frontend stack that is being commonly adopted and optimized for AI prototyping services.

The frontend application should powered by my backend API.(prompt fot tech spec)

Technical Restrictions:
The technical stack should be the following:
Typescript (statically typed javascript)
Tailwind CSS (css framework)
Vite.js (frontend tool)
ShadCN (UI components)

## **Pages**

### **Dashboard** (`/dashboard`)

#### **Purpose**
This page provides a summary of learning progress and serves as the default landing page when users visit the web app.

#### **Components**
- **Last Study Session**
  - Displays the last study activity used.
  - Shows the timestamp of the last activity.
  - Summarizes the number of correct vs. incorrect answers from the last activity.
  - Provides a link to the corresponding group.

- **Study Progress**
  - **Total words studied** (e.g., `3/124`)
    - Displays the total number of words studied across all study sessions out of the total available words in the database.
  - **Mastery progress** (e.g., `0% - 100%`)

- **Quick Stats**
  - **Success rate** (e.g., `80%`)
  - **Total study sessions** (e.g., `4`)
  - **Total active groups** (e.g., `3`)
  - **Study streak** (e.g., `4 days`)

- **Start Studying Button**
  - Redirects to the study activities page.

#### **Required API Endpoints**
- `GET /api/dashboard/last_study_session`
- `GET /api/dashboard/study_progress`
- `GET /api/dashboard/quick_stats`

---

### **Study Activities Index** (`/study_activities`)

#### **Purpose**
This page displays a collection of study activities, each represented by a thumbnail and name, allowing users to launch or view study activities.

#### **Components**
- **Study Activity Card**
  - Displays a thumbnail of the study activity.
  - Shows the name of the study activity.
  - Includes a "Launch" button to start the study activity.
  - Includes a "View" button to access past study sessions for the selected activity.

#### **Required API Endpoints**
- `GET /api/study_activities`

---

### **Study Activity Show** (`/study_activities/:id`)

#### **Purpose**
This page provides details about a specific study activity and its past study sessions.

#### **Components**
- **Study Activity Details**
  - Name of the study activity.
  - Thumbnail image.
  - Description of the study activity.
  - Launch button.

- **Study Activities Paginated List**
  - `id`
  - Activity name
  - Group name
  - Start time
  - End time (inferred from the last `word_review_item` submitted)
  - Number of review items

#### **Required API Endpoints**
- `GET /api/study_activities/:id`
- `GET /api/study_activities/:id/study_sessions`

---

### **Study Activity Launch** (`/study_activities/:id/launch`)

#### **Purpose**
This page is responsible for launching a study activity.

#### **Components**
- **Study Activity Name**
- **Launch Form**
  - Select field for group selection.
  - "Launch Now" button.

#### **Behavior**
- After submission, a new tab opens with the study activity based on the URL stored in the database.
- After submission, the page redirects to the study session show page.

#### **Required API Endpoints**
- `POST /api/study_activities`

---

### **Words Index** (`/words`)

#### **Purpose**
This page displays all words available in the database.

#### **Components**
- **Paginated Word List**
  - **Columns**
    - Japanese
    - Romaji
    - English
    - Correct Count
    - Wrong Count
  - Pagination: 100 items per page.
  - Clicking on a Japanese word redirects to the word details page.

#### **Required API Endpoints**
- `GET /api/words`

---

### **Word Show** (`/words/:id`)

#### **Purpose**
This page displays information about a specific word.

#### **Components**
- Japanese
- Romaji
- English
- **Study Statistics**
  - Correct Count
  - Wrong Count
- **Word Groups**
  - Displayed as a series of tags (e.g., pills).
  - Clicking a group name redirects to the corresponding group details page.

#### **Required API Endpoints**
- `GET /api/words/:id`

---

### **Word Groups Index** (`/groups`)

#### **Purpose**
This page lists all word groups available in the database.

#### **Components**
- **Paginated Group List**
  - **Columns**
    - Group Name
    - Word Count
  - Clicking a group name redirects to the group details page.

#### **Required API Endpoints**
- `GET /api/groups`

---

### **Group Show** (`/groups/:id`)

#### **Purpose**
This page provides details about a specific word group.

#### **Components**
- **Group Name**
- **Group Statistics**
  - Total Word Count
- **Words in Group** (Paginated List of Words)
  - Uses the same component as the Words Index page.
- **Study Sessions** (Paginated List of Study Sessions)
  - Uses the same component as the Study Sessions Index page.

#### **Required API Endpoints**
- `GET /api/groups/:id` (retrieves group name and statistics)
- `GET /api/groups/:id/words`
- `GET /api/groups/:id/study_sessions`

---

## **Study Sessions Index** (`/study_sessions`)

#### **Purpose**
This page displays a list of all study sessions recorded in the database.

#### **Components**
- **Paginated Study Session List**
  - **Columns**
    - `id`
    - Activity Name
    - Group Name
    - Start Time
    - End Time
    - Number of Review Items
  - Clicking a study session ID redirects to the study session details page.

#### **Required API Endpoints**
- `GET /api/study_sessions`

---

### **Study Session Show** (`/study_sessions/:id`)

#### **Purpose**
This page provides details about a specific study session.

#### **Components**
- **Study Session Details**
  - Activity Name
  - Group Name
  - Start Time
  - End Time
  - Number of Review Items
- **Words Review Items** (Paginated List of Words)
  - Uses the same component as the Words Index page.

#### **Required API Endpoints**
- `GET /api/study_sessions/:id`
- `GET /api/study_sessions/:id/words`

---

### **Settings Page** (`/settings`)

#### **Purpose**
This page allows users to configure the study portal.

#### **Components**
- **Theme Selection**
  - Options: Light, Dark, System Default
- **Reset History Button**
  - Deletes all study sessions and word review items.
- **Full Reset Button**
  - Drops all tables and recreates them with seed data.

#### **Required API Endpoints**
- `POST /api/reset_history`
- `POST /api/full_reset`
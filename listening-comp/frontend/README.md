```markdown
# **Listening Learning App**  

## **Objective**  
The goal is to develop a **Language Listening Comprehension App** that helps users improve their listening skills through AI-generated content. The app will pull **practice listening comprehension examples** from YouTube language learning tests and generate similar exercises.  

---

## **Approach**  

1️⃣ **Business Goal**  
   - You are an **Applied AI Engineer** tasked with building this **AI-driven language learning app**.  
   - The app will **extract YouTube content** to generate **customized listening exercises**.  

2️⃣ **Technical Uncertainty**  
   - ❓ **No Japanese knowledge** required—AI handles comprehension.  
   - ❓ **Vector store with SQLite3** for storing and retrieving documents.  
   - ❓ **TTS (Text-to-Speech) may be limited** in some languages.  
   - ❓ **ASR (Automatic Speech Recognition) quality may vary** by language.  
   - ❓ **Extracting accurate YouTube transcripts** is a challenge.  

---

## **Technical Requirements**  

✅ **Speech-to-Text (ASR)** – Convert spoken words into text.  
   - **Amazon Transcribe, OpenWhisper** as potential ASR solutions.  

✅ **YouTube Transcript API** – Retrieve and process transcripts from YouTube videos.  

✅ **LLM + Tool Use “Agent”** – AI models to process and generate learning content.  

✅ **SQLite3 - Knowledge Base** – Store structured text for efficient retrieval.  

✅ **Text-to-Speech (TTS)** – Convert generated text into speech.  
   - **Amazon Polly or equivalent AI-driven TTS solutions**.  

✅ **AI Coding Assistant** – Support for development.  
   - **Amazon Developer Q, Windsurf, Cursor, GitHub Copilot**.  

✅ **Frontend Framework** – User interface for accessing content.  
   - **Streamlit** for rapid frontend development.  

✅ **Guardrails** – Implement **content moderation & accuracy checks**.  

---

## **Running the Application**  

### **Run Frontend**  
```sh
streamlit run frontend/main.py
```

### **Run Backend**  
```sh
cd backend
pip install -r requirements.txt
cd ..
python backend/main.py
```

---

## **Next Steps**  
1️⃣ Validate **YouTube transcript extraction quality**.  
2️⃣ Test **ASR and TTS performance across multiple languages**.  
3️⃣ Implement **frontend interface for interaction**.  
4️⃣ Optimize **LLM processing for language exercises**.  
5️⃣ Integrate **SQLite3 for storing structured learning data**.  

🚀 **Listening Learning App development in progress!** 🚀  
```
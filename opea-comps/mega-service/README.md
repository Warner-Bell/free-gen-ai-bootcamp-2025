# **OPEA Mega Service**

**Created Mega Service**


**Install Requirements:**

Installing OPEA comps using a requirements file in the mega-service directory, and the below command.

```
pip install -r requirements.txt
```

**Run Mega-Service**

```
python app.py
```

**Test Code**



```sh
curl -X POST http://localhost:9000/v1/example-service \
  -H "Content-Type: application/json" \
  -d '{
    "model": "llama3.2:1b",
    "messages": "Yooo, wut up fam?"
  }' \
  -o response.json
  (-o output.json)
```


```sh
  curl -X POST http://localhost:9000/v1/example-service \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [
      {
        "role": "user",
        "content": "Aye, this a vet message"
      }
    ],
    "model": "test-model",
    "max_tokens": 100,
    "temperature": 0.7
  }'
```

##**Parked MS Implementation**
We spent an hour working through the MS implementation, we learned chat Q&A can work with Ollama, we need to figure out what the correct request handle is. we paused the research and TS to revist at a later time.
**prompt**
@workspace I have a transcript of a JLPT listening practice test. I want to extract all the
questions out. I need to extract it in the following structure per question:

introduction:
conversation:
question:

I want to use Amazon Bedrock with the Amazon Nova micro model.

I have other code that can be used to figure out how to write the code over in the
@get_transcript.py file. write the implementation for the data structure in the same style. Output the complete code here.

**Response Snippet**
```
I'll help you create a solution using Amazon Bedrock with the Claude model to extract JLPT listening questions from transcripts. Here's the implementation that follows a similar style to your existing code:

```
```
This implementation:

Uses a JLPTQuestion dataclass to structure the question data

Creates a JLPTTranscriptParser class that handles:

Interaction with Amazon Bedrock

Parsing transcripts into structured questions

Saving the parsed questions to a file

Follows error handling patterns similar to your existing code

Maintains similar file organization structure

To use this code:

Ensure you have the required AWS credentials and permissions set up

Create a "questions" directory for the output files

Place your transcript files in the "transcripts" directory

Run the script with appropriate file paths

The code uses Amazon Bedrock to analyze the transcript and extract the questions in the required format. It then structures this data and saves it in a readable format.

Note: You'll need to have the boto3 library installed and proper AWS credentials configured in your environment to use Amazon Bedrock.

```

**Result**



`I dont see reference to the nova micromodel?`

`like this, and dont returen json, use text instead; from typing import Optional, Dict, List
import boto3
import os

Model ID
#MODEL_ID = "amazon.nova-micro-v1:0"
MODEL_ID = "amazon.nova-lite-v1:0"

class TranscriptStructurer:`

#We should implement a cleanup strategy for the vector store files to prevent unnecessary buildup.

# Regular processing
python backend/vector_store.py

# Clean up old files
python backend/vector_store.py --cleanup

# Complete reset
python backend/vector_store.py --reset

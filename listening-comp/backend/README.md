# ** Listening-Comp-App**

**Build Steps**
Pulled fresh code from the repo, installed dependencies and verified app funtionality.

Using Q developer to assist with code generation.

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

Q behave similarly to Andrew's experience with winsurf  ie mising nova model id and trying to implement Titan models.

after checking with bedrock realized I do have access to Nova but could not enable it for some reason so had to use claud v2 it was very slow, and the /dev update were also slow so the process was becoming tedious.

I decided to try higher versions of claude and settled on v3.5 haiku, it made the question processing go much much faster.


`I dont see reference to the nova micromodel?`

`like this, and dont return json, use text instead; from typing import Optional, Dict, List
import boto3
import os

Model ID
#MODEL_ID = "anthropic.claude-3-haiku-20240307-v1:0"
MODEL_ID =  "anthropic.claude-3-haiku-20240307-v1:0"


updated structured data prompt to exclude images.  generating images would add too much copmplexity and i was already behind.

so I moved on to implement the vector store. I too had to instruct Q to isolate the vector file and structured data file processes, doing so significanly sped up the processsing.


#We should implement a cleanup strategy for the vector store files to prevent unnecessary buildup.

# Regular processing
python backend/vector_store.py

# Clean up old files
python backend/vector_store.py --cleanup

# Complete reset
python backend/vector_store.py --reset


another observation is that the /dev feature of que is painfully slow, compared to the speed of windsurf used by Andrew.

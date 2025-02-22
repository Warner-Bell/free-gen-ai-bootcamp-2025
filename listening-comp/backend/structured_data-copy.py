# File: backend/structured_data-copy.py

import sys
sys.stdout.reconfigure(line_buffering=True)

from typing import Optional, Dict, List
import boto3
import os
import json
import logging

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.StreamHandler(),
        logging.FileHandler('transcript_processing.log')
    ]
)
logger = logging.getLogger(__name__)

# Model ID - Claude 3 Haiku
MODEL_ID = "anthropic.claude-3-haiku-20240307-v1:0"

class QuestionFormat:
    """Constants for question formatting"""
    TEMPLATE = """
    <question>
    Introduction:
    [the situation setup in japanese]
    
    Conversation:
    [the dialogue in japanese]
    
    Question:
    [the question being asked in japanese]
    </question>
    """

    RULES = """
    Rules:
    - Only extract questions from the specified section
    - Ignore any practice examples (marked with 例)
    - Do not translate any Japanese text
    - Do not include any section descriptions or other text
    """

class SectionCriteria:
    """Criteria for different section types"""
    AUDIO_ONLY = """
    ONLY include questions that meet these criteria:
    - The answer can be determined purely from the spoken dialogue
    - No spatial/visual information is needed (like locations, layouts, or physical appearances)
    - No physical objects or visual choices need to be compared
    - DO NOT include any questions that require or reference images, pictures, diagrams, or visual choices
    
    For example, INCLUDE questions about:
    - Times and dates
    - Numbers and quantities
    - Spoken choices or decisions
    - Clear verbal directions
    
    DO NOT include questions about:
    - Physical locations that need a map or diagram
    - Visual choices between objects
    - Spatial arrangements or layouts
    - Physical appearances of people or things
    - Any question requiring an image or visual aid
    """

class TranscriptStructurer:
    def __init__(self, model_id: str = MODEL_ID):
        """Initialize the transcript structurer"""
        print("Initializing TranscriptStructurer...")
        self.bedrock_client = boto3.client('bedrock-runtime', region_name="us-east-1")
        self.model_id = model_id
        self.prompts = self._initialize_prompts()
        
        # Create necessary directories
        self.output_dir = 'backend/data/processed_questions'
        os.makedirs(self.output_dir, exist_ok=True)
        print("Initialization complete.")

    def _initialize_prompts(self) -> Dict[int, str]:
        """Initialize prompts for each section"""
        prompts = {}
        
        base_prompt = """You are a JLPT transcript analyzer. Your task is to extract and structure questions from specific sections.
        Maintain the original Japanese text exactly as it appears. Do not translate or modify any Japanese content."""
        
        # Section 1: Audio-only questions
        prompts[1] = f"""{base_prompt}
        Focus on section 問題1 of this JLPT transcript.
        {SectionCriteria.AUDIO_ONLY}
        Format each question exactly like this:
        {QuestionFormat.TEMPLATE}
        {QuestionFormat.RULES}
        - Only include questions where answers can be determined from dialogue alone
        """

        # Section 2: Audio-only questions with different context
        prompts[2] = f"""{base_prompt}
        Focus on section 問題2 of this JLPT transcript.
        {SectionCriteria.AUDIO_ONLY}
        Format each question exactly like this:
        {QuestionFormat.TEMPLATE}
        {QuestionFormat.RULES}
        - Only include questions where answers can be determined from dialogue alone
        """

        # Section 3: All questions
        prompts[3] = f"""{base_prompt}
        Focus on section 問題3 of this JLPT transcript.
        Format each question exactly like this:
        {QuestionFormat.TEMPLATE}
        {QuestionFormat.RULES}
        """

        return prompts

    def _invoke_bedrock(self, prompt: str, transcript: str) -> Optional[str]:
        """Invoke Bedrock model with the given prompt and transcript"""
        try:
            body = {
                "anthropic_version": "bedrock-2023-05-31",
                "messages": [
                    {
                        "role": "user",
                        "content": [
                            {
                                "type": "text",
                                "text": f"{prompt}\n\nTranscript:\n{transcript}"
                            }
                        ]
                    }
                ],
                "max_tokens": 2000,
                "temperature": 0,
                "top_p": 1,
                "stop_sequences": ["Human:", "Assistant:"]
            }
            
            print("Sending request to Bedrock...")
            response = self.bedrock_client.invoke_model(
                modelId=self.model_id,
                body=json.dumps(body)
            )
            
            response_body = json.loads(response['body'].read())
            print("Response received from Bedrock")
            return self._process_response(response_body)
            
        except Exception as e:
            print(f"Error invoking Bedrock: {str(e)}")
            logger.error(f"Error invoking Bedrock: {str(e)}")
            return None

    def _process_response(self, response_body: dict) -> Optional[str]:
        """Process the response from Bedrock"""
        try:
            content = ''
            if 'messages' in response_body:
                content = response_body['messages'][0]['content'][0]['text']
            elif isinstance(response_body.get('content'), list):
                content = response_body['content'][0].get('text', '')
            else:
                content = response_body.get('content', '')
            
            print(f"Extracted content length: {len(content) if content else 0}")
            return self._process_questions(content) if content else None
        except Exception as e:
            print(f"Error processing response: {str(e)}")
            logger.error(f"Error processing response: {str(e)}")
            return None

    def _process_questions(self, content: str) -> str:
        """Process and number the questions"""
        if not content:
            return ""
        
        questions = [q.strip() for q in content.split('</question>') if q.strip()]
        processed_questions = []
        
        for i, question in enumerate(questions, 1):
            question = self._clean_question_text(question)
            if question:
                question = question.replace('\n\n', '\n').strip()
                processed_question = f'<question number="{i}">\n{question}\n</question>'
                processed_questions.append(processed_question)
        
        return '\n\n'.join(processed_questions)

    def _clean_question_text(self, question: str) -> str:
        """Clean up question text by removing tags and extra whitespace"""
        for i in range(1, 6):
            question = question.replace(f'<question number="{i}">', '')
        question = question.replace('<question>', '')
        
        lines = [line.strip() for line in question.split('\n')]
        return '\n'.join(line for line in lines if line)

    def process_transcript(self, transcript_path: str) -> Dict[int, str]:
        """Process transcript and extract questions"""
        results = {}
        
        print(f"\nLoading transcript from: {transcript_path}")
        transcript = self.load_transcript(transcript_path)
        if not transcript:
            print("Failed to load transcript")
            return results
            
        for section_num in range(1, 4):
            print(f"\nProcessing section {section_num}...")
            result = self._invoke_bedrock(self.prompts[section_num], transcript)
            if result:
                print(f"Successfully processed section {section_num}")
                results[section_num] = result
            else:
                print(f"Warning: No result for section {section_num}")
        
        return results

    def load_transcript(self, filename: str) -> Optional[str]:
        """Load transcript from a file"""
        try:
            with open(filename, 'r', encoding='utf-8') as f:
                content = f.read()
                print(f"Loaded transcript: {len(content)} characters")
                return content
        except Exception as e:
            print(f"Error loading transcript: {str(e)}")
            logger.error(f"Error loading transcript: {str(e)}")
            return None

    def save_results(self, results: Dict[int, str], base_filename: str) -> bool:
        """Save processed questions to files"""
        try:
            for section_num, content in results.items():
                filename = f"{self.output_dir}/{base_filename}_section{section_num}.txt"
                with open(filename, 'w', encoding='utf-8') as f:
                    f.write(content)
                print(f"Saved section {section_num} to: {filename}")
            return True
        except Exception as e:
            print(f"Error saving results: {str(e)}")
            logger.error(f"Error saving results: {str(e)}")
            return False

def main():
    """Main execution function"""
    print("\nStarting transcript processing...")
    structurer = TranscriptStructurer()
    
    try:
        # Process transcript
        transcript_path = "backend/data/transcripts/sY7L5cfCWno.txt"
        print(f"\nProcessing transcript: {transcript_path}")
        results = structurer.process_transcript(transcript_path)
        
        # Save results
        if results:
            base_filename = os.path.splitext(os.path.basename(transcript_path))[0]
            structurer.save_results(results, base_filename)
            print("\nProcessing completed successfully!")
        else:
            print("\nWarning: No results generated")
    except Exception as e:
        print(f"\nError in main execution: {str(e)}")
        logger.error(f"Error in main execution: {str(e)}")

if __name__ == "__main__":
    main()

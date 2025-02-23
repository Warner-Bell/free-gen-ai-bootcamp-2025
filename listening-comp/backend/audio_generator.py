import boto3
import json
import tempfile
from pathlib import Path
from typing import List, Tuple, Union, Dict, Optional

class AudioGenerator:
    def __init__(self):
        self.bedrock = boto3.client(
            service_name='bedrock-runtime',
            region_name='us-east-1'
        )
        print("Initializing audio generator...")
        
    def _invoke_bedrock(self, prompt: str) -> str:
        """Call Bedrock with the given prompt."""
        try:
            body = json.dumps({
                "prompt": prompt,
                "max_tokens_to_sample": 500,
                "temperature": 0.1,
                "top_p": 0.9,
            })
            
            response = self.bedrock.invoke_model(
                modelId="anthropic.claude-v2",
                body=body
            )
            response_body = json.loads(response.get('body').read())
            return response_body.get('completion', '')
        except Exception as e:
            print(f"Error calling Bedrock: {str(e)}")
            raise

    def parse_conversation(self, question: Union[str, Dict], max_retries: int = 3) -> List[Tuple[str, str, str]]:
        """Parse the conversation using Bedrock Claude."""
        # Convert Dict input to string if necessary
        if isinstance(question, dict):
            conversation = question.get("Conversation", "")
            if not conversation:
                raise ValueError("No conversation found in input dictionary")
            question = conversation

        def normalize_dialogue(text: str) -> Tuple[str, List[str]]:
            """Extract and normalize dialogue lines from text."""
            dialogue_lines = []
            print("\nProcessing raw input:")
            print("----------------------------------------")
            print(text)
            print("----------------------------------------\n")
            
            lines = []
            for line in text.split('\n'):
                line = ' '.join(line.strip().split())
                if line:
                    lines.append(line)
            
            for line in lines:
                if ':' not in line:
                    continue
                    
                if any(line.lower().startswith(p) for p in ['question', '質問', '会話', 'next', '次の']):
                    continue
                
                try:
                    parts = line.split(':', 1)
                    if len(parts) != 2:
                        continue
                        
                    speaker = parts[0].strip()
                    text = parts[1].strip()
                    
                    if not (speaker and text):
                        continue
                        
                    if any(p in speaker.lower() for p in ['question', '質問', '会話']):
                        continue
                    
                    dialogue_lines.append(f"{speaker}: {text}")
                except Exception as e:
                    print(f"Skipping line due to error: {str(e)}")
                    continue
                    
            result = '\n'.join(dialogue_lines)
            print("\nNormalized dialogue:")
            print("----------------------------------------")
            print(result)
            print("----------------------------------------\n")
            return result, dialogue_lines
        
        # Normalize input text and get dialogue lines
        normalized_text, dialogue_lines = normalize_dialogue(question)

        # Validate we have dialogue to process
        if not dialogue_lines:
            raise ValueError("No valid dialogue lines found in input")

        prompt = """会話をJSONに変換:

例:
田中: すみません、駅はどこですか。
駅員: ここは新宿駅です。

出力:
[
  {"speaker": "田中", "text": "すみません、駅はどこですか。", "gender": "male"},
  {"speaker": "駅員", "text": "ここは新宿駅です。", "gender": "male"}
]

会話:
{input_text}

出力を以下の形式の正確なJSONで提供してください。回答は JSON 配列のみを含めてください:"""

        formatted_prompt = prompt.format(input_text=normalized_text)
            
        for attempt in range(max_retries):
            try:
                response = self._invoke_bedrock(formatted_prompt)
                if not response:
                    raise ValueError("Empty response from API")

                start_idx = response.find('[')
                end_idx = response.rfind(']') + 1
                if start_idx == -1 or end_idx == 0:
                    raise ValueError("No JSON array found in response")

                parts = json.loads(response[start_idx:end_idx])
                
                for part in parts:
                    if part["gender"] not in ["male", "female"]:
                        part["gender"] = "male"
                
                result = [(p["speaker"], p["text"], p["gender"]) for p in parts]
                if self.validate_conversation_parts(result):
                    return result
                raise ValueError("Invalid conversation structure")
                    
            except Exception as e:
                print(f"\nAttempt {attempt + 1} failed: {str(e)}")
                if attempt < max_retries - 1:
                    print("Retrying...")
                    continue
                
                print("\nInput that caused final failure:")
                print("----------------------------------------")
                print(normalized_text)
                print("----------------------------------------")
                if 'response' in locals():
                    print("\nRaw response that caused failure:")
                    print("----------------------------------------")
                    print(response)
                    print("----------------------------------------")
                raise ValueError("Failed to parse conversation after all retry attempts")

    def validate_conversation_parts(self, parts: List[Tuple[str, str, str]]) -> bool:
        """Validate that conversation parts are well-formed."""
        if not parts:
            return False
            
        for speaker, text, gender in parts:
            if not all([speaker, text, gender]):
                return False
            if gender not in ["male", "female"]:
                return False
                
        return True
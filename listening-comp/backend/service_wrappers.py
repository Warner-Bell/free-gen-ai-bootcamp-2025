import boto3
import json
import os
import tempfile
import subprocess
import time
from typing import Dict, List, Optional, Tuple, Union
import pathlib

class BedrockServiceWrapper:
    def __init__(self, model_id: str):
        """Initialize the Bedrock service wrapper."""
        self.bedrock = boto3.client('bedrock-runtime')
        self.model_id = model_id

    def invoke_model(self, prompt: str) -> str:
        """Invoke Bedrock model with the given prompt."""
        try:
            response = self.bedrock.invoke_model(
                modelId=self.model_id,
                body=json.dumps({
                    "prompt": prompt,
                    "max_tokens_to_sample": 1000,
                    "temperature": 0.1,
                    "top_p": 0.9,
                })
            )
            response_body = json.loads(response.get('body').read())
            return response_body.get('completion', '')
        except Exception as e:
            print(f"\n🚨 [ERROR] Error calling Bedrock: {str(e)}")
            return ""

class PollyServiceWrapper:
    def __init__(self):
        """Initialize Polly service wrapper."""
        self.polly = boto3.client('polly')

    def synthesize_speech(self, text: str, voice_id: str) -> Optional[bytes]:
        """Generate speech from text using Amazon Polly."""
        try:
            response = self.polly.synthesize_speech(
                Engine='neural',
                LanguageCode='ja-JP',
                OutputFormat='mp3',
                Text=text,
                TextType='text',
                VoiceId=voice_id
            )
            if "AudioStream" in response:
                return response["AudioStream"].read()
            return None
        except Exception as e:
            print(f"\n🚨 [ERROR] Polly synthesis failed: {str(e)}")
            return None

class AudioGenerator:
    def __init__(self):
        """Initialize the audio generator."""
        self.bedrock = BedrockServiceWrapper("anthropic.claude-v2")
        self.polly = PollyServiceWrapper()
        print("🟢 Initializing audio generator...")
        self.audio_dir = os.path.join(
            os.path.dirname(os.path.dirname(os.path.abspath(__file__))), 
            'frontend', 'public', 'audio'
        )
        os.makedirs(self.audio_dir, exist_ok=True)

    def parse_conversation(self, question: Union[str, Dict], max_retries: int = 3) -> List[Tuple[str, str, str]]:
        """Parse the conversation using Bedrock Claude.
        Returns a list of tuples containing (speaker, text, gender)."""
        
        # Convert Dict input to string if necessary
        if isinstance(question, dict):
            conversation = question.get("Conversation", "")
            if not conversation:
                raise ValueError("No conversation found in input dictionary")
            question = conversation

        # Extract dialogue lines
        print("\n🟡 [INFO] Processing raw input:")
        print("----------------------------------------")
        print(question)
        print("----------------------------------------\n")
            
        # Split into lines and clean up
        lines = [line.strip() for line in question.split("\n") if line.strip()]
        
        # Process each line
        dialogue_lines = []
        for line in lines:
            if ":" not in line:
                continue
                
            if any(line.lower().startswith(p) for p in ['question', '質問', '会話', 'next', '次の']):
                continue
            
            try:
                speaker, text = line.split(":", 1)  # Split at first colon
                speaker = speaker.strip()
                text = text.strip()

                if not text or not speaker:
                    continue
                
                if any(p in speaker.lower() for p in ['question', '質問', '会話']):
                    continue
                
                dialogue_lines.append((speaker, text))

            except Exception as e:
                print(f"\n⚠️ [WARNING] Skipping line due to error: {str(e)}")
                continue
                
        if not dialogue_lines:
            raise ValueError("No valid dialogue lines found in input")

        print("\n🟢 [INFO] Extracted conversations:")
        for speaker, text in dialogue_lines:
            print(f"{speaker}: {text}")
        print("----------------------------------------\n")

        # Create prompt for gender inference
        dialogue_text = "\n".join(f"{speaker}: {text}" for speaker, text in dialogue_lines)
        prompt = f"""会話の各行について、話者の性別を推測してJSON形式で出力してください:

例:
田中: すみません、駅はどこですか。
駅員: ここは新宿駅です。

出力:
[
  {{"speaker": "田中", "text": "すみません、駅はどこですか。", "gender": "male"}},
  {{"speaker": "駅員", "text": "ここは新宿駅です。", "gender": "male"}}
]

会話:
{dialogue_text}

上記の会話を同じJSON形式で出力してください。性別は話者から推測してください。"""

        # Get gender inference from Claude
        for attempt in range(max_retries):
            try:
                print("\n🔄 [INFO] Calling Bedrock for gender inference (attempt {}/{})".format(attempt + 1, max_retries))
                
                # Get response from Bedrock
                response = self.bedrock.invoke_model(prompt)
                if not response:
                    raise ValueError("Empty response from API")

                # Find and extract JSON array
                start_idx = response.find('[')
                end_idx = response.rfind(']') + 1
                if start_idx == -1 or end_idx == 0:
                    print("\n🔍 [DEBUG] Raw API response:")
                    print(response)
                    raise ValueError("No JSON array found in response")

                # Parse JSON response
                try:
                    json_str = response[start_idx:end_idx]
                    parts = json.loads(json_str)
                except json.JSONDecodeError as e:
                    print(f"\n🚨 [ERROR] Failed to decode JSON response: {str(e)}")
                    print(f"\n🔍 [DEBUG] Failed JSON string: {json_str}")
                    raise ValueError("Invalid JSON response")

                # Validate response structure
                if not isinstance(parts, list):
                    raise ValueError("Invalid JSON structure - expected a list")

                # Convert JSON objects to tuples
                result = []
                for part in parts:
                    if isinstance(part, dict):
                        speaker = part.get("speaker", "Unknown")
                        text = part.get("text", "")
                        gender = part.get("gender", "male")
                        if text and speaker:
                            result.append((speaker, text, gender))

                # Validate final result
                if not result:
                    raise ValueError("No valid dialogue parts found in response")

                if not self.validate_conversation_parts(result):
                    raise ValueError("Invalid conversation structure")

                print("\n✅ [INFO] Successfully parsed conversation with {} parts".format(len(result)))
                return result

                

            except Exception as e:
                print(f"\n⚠️ [WARNING] Attempt {attempt + 1} failed: {str(e)}")
                if attempt < max_retries - 1:
                    print("🔁 Retrying...")
                    continue
                
                print("\n❌ [ERROR] Final failure after all retry attempts.")
                if 'response' in locals():
                    print("\n📤 [DEBUG] Raw response that caused failure:")
                    print("----------------------------------------")
                    print(response)
                    print("----------------------------------------")
                raise ValueError("Failed to parse conversation after all retry attempts")

    def validate_conversation_parts(self, parts: List[Tuple[str, str, str]]) -> bool:
        """Validate that conversation parts are well-formed.
        Each part should be a tuple of (speaker, text, gender)."""
        if not isinstance(parts, list) or not parts:
            print("\n⚠️ [WARNING] Empty or invalid parts list")
            return False
        
        for part in parts:
            if not isinstance(part, tuple) or len(part) != 3:
                print(f"\n⚠️ [WARNING] Invalid part format: {part}")
                return False
            speaker, text, gender = part
            if not all(isinstance(x, str) for x in (speaker, text, gender)):
                print(f"\n⚠️ [WARNING] Invalid data types in part: {part}")
                return False
            if not all(x.strip() for x in (speaker, text, gender)):
                print(f"\n⚠️ [WARNING] Empty strings in part: {part}")
                return False
        return True

    def generate_audio(self, question: Union[str, Dict]) -> str:
        """Generate audio from text using AWS Polly."""
        try:
            parts = self.parse_conversation(question)
            if not parts:
                raise ValueError("No conversation parts to process")

            audio_files = []
            timestamp = int(time.time())
            
            for speaker, text, gender in parts:
                voice_id = self.get_voice_for_gender(gender)
                audio_file = self.generate_audio_part(text, voice_id)
                if audio_file:
                    audio_files.append(audio_file)
                else:
                    print(f"\n⚠️ [WARNING] Failed to generate audio for: {text}")

            if not audio_files:
                raise ValueError("No audio files were generated")

            # Combine all audio files
            output_file = os.path.join(self.audio_dir, f'conversation_{timestamp}.mp3')
            if self.combine_audio_files(audio_files, output_file):
                return output_file
            raise ValueError("Failed to combine audio files")

        except Exception as e:
            print(f"\n🚨 [ERROR] Audio generation failed: {str(e)}")
            return ""

    def get_voice_for_gender(self, gender: str) -> str:
        """Get appropriate voice ID based on gender."""
        return "Kazuha" if gender.lower() == "female" else "Takumi"

    def generate_audio_part(self, text: str, voice_name: str) -> Optional[str]:
        """Generate audio for a single part of the conversation."""
        try:
            audio_data = self.polly.synthesize_speech(text, voice_name)
            if not audio_data:
                return None
                
            with tempfile.NamedTemporaryFile(suffix='.mp3', delete=False) as temp_file:
                temp_file.write(audio_data)
                return temp_file.name
        except Exception as e:
            print(f"\n🚨 [ERROR] Failed to generate audio part: {str(e)}")
            return None

    def combine_audio_files(self, audio_files: List[str], output_file: str) -> bool:
        """Combine multiple audio files into one."""
        try:
            with tempfile.NamedTemporaryFile('w', suffix='.txt', delete=False) as f:
                for audio_file in audio_files:
                    f.write(f"file '{audio_file}'\n")
                concat_list = f.name

            subprocess.run([
                'ffmpeg', '-f', 'concat', '-safe', '0',
                '-i', concat_list,
                '-c', 'copy', output_file
            ], check=True, capture_output=True)

            # Clean up temporary files
            os.unlink(concat_list)
            for file in audio_files:
                os.unlink(file)

            return True
        except Exception as e:
            print(f"\n🚨 [ERROR] Failed to combine audio files: {str(e)}")
            return False
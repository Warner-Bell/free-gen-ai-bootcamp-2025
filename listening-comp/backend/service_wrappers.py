import boto3
import json
import os
import logging
from typing import Dict, List, Tuple, Optional
from datetime import datetime
import tempfile
import subprocess

# Service Wrappers for AWS Bedrock and Polly
class BedrockServiceWrapper:
    def __init__(self, model_id: str):
        """Initialize the AWS Bedrock runtime client."""
        self.bedrock_runtime = boto3.client('bedrock-runtime')
        self.model_id = model_id

    def converse(self, prompt: str) -> str:
        """
        Handles conversation with Claude model through Amazon Bedrock.
        Ensures response is in valid JSON format.
        """
        try:
            body = {
                "anthropic_version": "bedrock-2023-05-31",
                "max_tokens": 4096,
                "messages": [
                    {
                        "role": "user",
                        "content": prompt
                    }
                ]
            }

            response = self.bedrock_runtime.invoke_model(
                modelId=self.model_id,
                body=json.dumps(body)
            )
            
            response_body = json.loads(response['body'].read().decode('utf-8'))

            # 🚀 Print the full raw response from Bedrock
            print(f"\nDEBUG: Full Bedrock Response:\n{json.dumps(response_body, indent=2)}\n")

            if "content" in response_body and isinstance(response_body["content"], list):
                json_content = response_body["content"][0].get("text", "").strip()

                # 🚀 Print the extracted content (should be JSON)
                print(f"\nDEBUG: Extracted Content:\n{json_content}\n")

                if json_content.startswith("[") and json_content.endswith("]"):
                    return json_content
                else:
                    raise Exception(f"Invalid JSON format received: {json_content}")

            raise Exception(f"Unexpected response format from Bedrock: {response_body}")

        except Exception as e:
            raise Exception(f"Bedrock converse error: {str(e)}")

    def invoke(self, prompt: str) -> str:
        """Alias for converse method to maintain compatibility."""
        return self.converse(prompt)


class PollyServiceWrapper:
    def __init__(self):
        """Initialize the Amazon Polly client."""
        self.polly = boto3.client('polly')
        
    def synthesize_speech(self, text: str, voice_id: str) -> Optional[bytes]:
        """Generate speech using Amazon Polly and return audio bytes."""
        try:
            response = self.polly.synthesize_speech(
                Engine='neural',
                Text=text,
                OutputFormat='mp3',
                VoiceId=voice_id
            )
            
            if 'AudioStream' in response:
                return response['AudioStream'].read()
            else:
                raise Exception(f"Polly did not return an AudioStream: {response}")

        except Exception as e:
            print(f"Error synthesizing speech: {str(e)}")
            return None


# Main Audio Generator Class
class AudioGenerator:
    def __init__(self):
        """Initialize the audio generator with AWS clients."""
        self.bedrock = BedrockServiceWrapper("anthropic.claude-3-haiku-20240307-v1:0")
        self.polly = PollyServiceWrapper()
        
        # Define Japanese neural voices by gender
        self.voices = {
            'male': ['Takumi'],
            'female': ['Kazuha'],
            'announcer': 'Takumi'  # Default announcer voice
        }
        
        # Create audio output directory
        self.audio_dir = os.path.join(
            os.path.dirname(os.path.dirname(os.path.abspath(__file__))),
            "frontend/static/audio"
        )
        os.makedirs(self.audio_dir, exist_ok=True)

    def parse_conversation(self, question: Dict) -> List[Tuple[str, str, str]]:
        """
        Convert question into a format for audio generation.
        Returns a list of (speaker, text, gender) tuples.
        """
        max_retries = 3
        for attempt in range(max_retries):
            try:
                response = self.bedrock.converse(question["Conversation"])

                print(f"DEBUG: Received response from Bedrock: {response}")  # Debugging 🚀
                
                parts = json.loads(response)

                # Ensure all required keys are present
                for i, part in enumerate(parts):
                    if "speaker" not in part or "text" not in part or "gender" not in part:
                        raise Exception(f"Invalid format in part {i}: {part}")

                return [(p["speaker"], p["text"], p["gender"]) for p in parts]

            except json.JSONDecodeError as e:
                print(f"Attempt {attempt + 1} failed: Invalid JSON response - {e}")
            except KeyError as e:
                print(f"Attempt {attempt + 1} failed: Missing key in response - {e}")
            except Exception as e:
                print(f"Attempt {attempt + 1} failed: {str(e)}")

        raise Exception("Failed to parse conversation after multiple attempts")

    def generate_audio(self, question: Dict) -> str:
        """
        Generate audio for the entire question.
        Returns the path to the generated audio file.
        """
        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        output_file = os.path.join(self.audio_dir, f"question_{timestamp}.mp3")
        
        try:
            # Parse conversation into parts
            parts = self.parse_conversation(question)

            # Generate audio for each part
            audio_parts = []
            for speaker, text, gender in parts:
                voice = self.get_voice_for_gender(gender)
                print(f"Using voice {voice} for {speaker} ({gender})")
                audio_file = self.generate_audio_part(text, voice)
                if audio_file:
                    audio_parts.append(audio_file)

            # Combine all parts into final audio
            if self.combine_audio_files(audio_parts, output_file):
                return output_file
            else:
                raise Exception("Failed to combine audio files")

        except Exception as e:
            if os.path.exists(output_file):
                os.unlink(output_file)
            raise Exception(f"Audio generation failed: {str(e)}")

    def get_voice_for_gender(self, gender: str) -> str:
        """Get an appropriate voice for the given gender."""
        return self.voices.get(gender, ["Takumi"])[0]  # Default to 'Takumi' if unknown

    def generate_audio_part(self, text: str, voice_name: str) -> Optional[str]:
        """Generate audio for a single part using Amazon Polly."""
        audio_data = self.polly.synthesize_speech(text, voice_name)
        if not audio_data:
            return None
        
        with tempfile.NamedTemporaryFile(suffix='.mp3', delete=False) as temp_file:
            temp_file.write(audio_data)
            return temp_file.name

    def combine_audio_files(self, audio_files: List[str], output_file: str) -> bool:
        """Combine multiple audio files using ffmpeg."""
        try:
            with tempfile.NamedTemporaryFile('w', suffix='.txt', delete=False) as f:
                for audio_file in audio_files:
                    f.write(f"file '{audio_file}'\n")
                file_list = f.name
            
            subprocess.run([
                'ffmpeg', '-f', 'concat', '-safe', '0',
                '-i', file_list,
                '-c', 'copy',
                output_file
            ], check=True)

            return True
        except Exception as e:
            print(f"Error combining audio files: {str(e)}")
            return False
        finally:
            os.unlink(file_list)
            for audio_file in audio_files:
                if os.path.exists(audio_file):
                    os.unlink(audio_file)

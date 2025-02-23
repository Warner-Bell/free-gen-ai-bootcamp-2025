# Create BedrockChat
# bedrock_chat.py
import boto3
import streamlit as st
import json
from typing import Optional, Dict, Any


# Model ID
MODEL_ID = "anthropic.claude-3-haiku-20240307-v1:0"



class BedrockChat:
    def __init__(self, model_id: str = MODEL_ID):
        """Initialize Bedrock chat client"""
        self.bedrock_client = boto3.client('bedrock-runtime', region_name="us-east-1")
        self.model_id = model_id

    def generate_response(self, message: str, inference_config: Optional[Dict[str, Any]] = None) -> Optional[str]:
        """Generate a response using Amazon Bedrock"""
        if inference_config is None:
            inference_config = {"temperature": 0.7}

        try:
            body = {
                "anthropic_version": "bedrock-2023-05-31",
                "messages": [{
                    "role": "user",
                    "content": message
                }],
                "max_tokens": 4096,
                **inference_config
            }
            
            response = self.bedrock_client.invoke_model(
                modelId=self.model_id,
                body=json.dumps(body)
            )
            response_body = json.loads(response['body'].read())
            return response_body['content'][0]['text']
                
        except Exception as e:
            st.error(f"Error generating response: {str(e)}")
            return None


if __name__ == "__main__":
    chat = BedrockChat()
    while True:
        user_input = input("You: ")
        if user_input.lower() == '/exit':
            break
        response = chat.generate_response(user_input)
        print("Bot:", response)

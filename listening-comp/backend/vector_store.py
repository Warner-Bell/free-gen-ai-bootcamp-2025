# File: backend/vector_store.py

import chromadb
from chromadb.utils import embedding_functions
import json
import os
import sys  # Add this import
import boto3
from typing import Dict, List, Optional, Any
import logging
import shutil
from datetime import datetime, timedelta

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.StreamHandler(),
        logging.FileHandler('vector_store.log')
    ]
)
logger = logging.getLogger(__name__)

class BedrockEmbeddingFunction(embedding_functions.EmbeddingFunction):
    def __init__(self, model_id: str = "amazon.titan-embed-text-v1"):
        self.bedrock_client = boto3.client('bedrock-runtime', region_name="us-east-1")
        self.model_id = model_id
        self.dimensions = 1536

    def __call__(self, texts: List[str]) -> List[List[float]]:
        if not texts:
            return []
            
        embeddings = []
        for text in texts:
            try:
                if not isinstance(text, str):
                    print(f"Skipping non-string input: {type(text)}")
                    embeddings.append([0.0] * self.dimensions)
                    continue
                    
                response = self.bedrock_client.invoke_model(
                    modelId=self.model_id,
                    body=json.dumps({
                        "inputText": text
                    })
                )
                response_body = json.loads(response['body'].read())
                embedding = response_body['embedding']
                embeddings.append(embedding)
            except Exception as e:
                print(f"Error generating embedding: {str(e)}")
                logger.error(f"Error generating embedding: {str(e)}")
                embeddings.append([0.0] * self.dimensions)
        return embeddings

class QuestionVectorStore:
    def __init__(self, persist_directory: str = "backend/data/vectorstore"):
        """Initialize the vector store"""
        self.persist_directory = persist_directory
        self._ensure_directory_structure()
        
        try:
            self.client = chromadb.PersistentClient(path=persist_directory)
            self.embedding_fn = BedrockEmbeddingFunction()
            self.collections = self._initialize_collections()
            print("Successfully initialized QuestionVectorStore")
            logger.info("Successfully initialized QuestionVectorStore")
        except Exception as e:
            error_msg = f"Failed to initialize QuestionVectorStore: {str(e)}"
            print(error_msg)
            logger.error(error_msg)
            raise

    def _ensure_directory_structure(self) -> None:
        """Ensure the vector store directory exists"""
        try:
            os.makedirs(self.persist_directory, exist_ok=True)
            print(f"Ensured directory structure at {self.persist_directory}")
            logger.info(f"Ensured directory structure at {self.persist_directory}")
        except Exception as e:
            error_msg = f"Failed to create directory structure: {str(e)}"
            print(error_msg)
            logger.error(error_msg)
            raise

    def _initialize_collections(self) -> Dict[str, Any]:
        """Initialize or get collections"""
        collections = {}
        section_configs = {
            "section2": {
                "name": "section2_questions",
                "description": "JLPT listening comprehension questions - Section 2"
            },
            "section3": {
                "name": "section3_questions",
                "description": "JLPT phrase matching questions - Section 3"
            }
        }
        
        for section, config in section_configs.items():
            try:
                # Try to get existing collection
                try:
                    collection = self.client.get_collection(
                        name=config["name"],
                        embedding_function=self.embedding_fn
                    )
                    # If collection exists, delete it
                    self.client.delete_collection(name=config["name"])
                    print(f"Deleted existing collection: {config['name']}")
                    logger.info(f"Deleted existing collection: {config['name']}")
                except Exception:
                    # Collection doesn't exist, which is fine
                    pass

                # Create new collection
                collection = self.client.create_collection(
                    name=config["name"],
                    embedding_function=self.embedding_fn,
                    metadata={"description": config["description"]}
                )
                
                collections[section] = collection
                print(f"Created new collection: {config['name']}")
                logger.info(f"Created new collection: {config['name']}")
                
            except Exception as e:
                error_msg = f"Failed to initialize collection {section}: {str(e)}"
                print(error_msg)
                logger.error(error_msg)
                raise
                
        return collections

    def add_questions(self, section_num: int, questions: List[Dict], video_id: str) -> None:
        """Add questions to the vector store"""
        if not questions:
            print("No questions provided")
            logger.warning("No questions provided")
            return
            
        if section_num not in [2, 3]:
            error_msg = "Only sections 2 and 3 are currently supported"
            print(error_msg)
            logger.error(error_msg)
            raise ValueError(error_msg)
            
        collection = self.collections.get(f"section{section_num}")
        if not collection:
            error_msg = f"Collection not found for section {section_num}"
            print(error_msg)
            logger.error(error_msg)
            raise ValueError(error_msg)
            
        try:
            # Prepare questions for adding
            ids = []
            documents = []
            metadatas = []
            
            for idx, question in enumerate(questions):
                question_id = f"{video_id}_{section_num}_{idx}"
                ids.append(question_id)
                
                metadatas.append({
                    "video_id": video_id,
                    "section": section_num,
                    "question_index": idx,
                    "full_structure": json.dumps(question)
                })
                
                doc_parts = []
                for key in ['Introduction', 'Conversation', 'Question']:
                    if key in question:
                        doc_parts.append(f"{key}: {question[key]}")
                documents.append("\n".join(doc_parts))
            
            # Add new questions
            collection.add(
                ids=ids,
                documents=documents,
                metadatas=metadatas
            )
            
            msg = f"Added {len(questions)} questions to section {section_num}"
            print(msg)
            logger.info(msg)
                
        except Exception as e:
            error_msg = f"Error adding questions: {str(e)}"
            print(error_msg)
            logger.error(error_msg)
            raise

    def search_similar_questions(self, section_num: int, query: str, n_results: int = 5) -> List[Dict]:
        """Search for similar questions in the vector store"""
        if not query:
            print("Empty query provided")
            logger.warning("Empty query provided")
            return []
            
        collection = self.collections.get(f"section{section_num}")
        if not collection:
            error_msg = f"Collection not found for section {section_num}"
            print(error_msg)
            logger.error(error_msg)
            raise ValueError(error_msg)
            
        try:
            results = collection.query(
                query_texts=[query],
                n_results=n_results
            )
            return self._format_search_results(results)
        except Exception as e:
            error_msg = f"Error searching questions: {str(e)}"
            print(error_msg)
            logger.error(error_msg)
            return []

    def _format_search_results(self, results: Dict) -> List[Dict]:
        """Format search results"""
        formatted_results = []
        try:
            for idx, metadata in enumerate(results['metadatas'][0]):
                try:
                    question_data = json.loads(metadata['full_structure'])
                    question_data['similarity_score'] = results['distances'][0][idx]
                    formatted_results.append(question_data)
                except json.JSONDecodeError:
                    error_msg = f"Error decoding question structure for index {idx}"
                    print(error_msg)
                    logger.error(error_msg)
                    continue
        except Exception as e:
            error_msg = f"Error formatting results: {str(e)}"
            print(error_msg)
            logger.error(error_msg)
            
        return formatted_results

    def index_questions_file(self, filename: str, section_num: int):
        """Index all questions from a file into the vector store"""
        # Extract video ID from filename
        video_id = os.path.basename(filename).split('_section')[0]
        
        # Parse questions from file
        questions = self.parse_questions_from_file(filename)
        
        # Add to vector store
        if questions:
            self.add_questions(section_num, questions, video_id)
            print(f"Indexed {len(questions)} questions from {filename}")
        else:
            print(f"No questions found in {filename}")

    def parse_questions_from_file(self, filename: str) -> List[Dict]:
        """Parse questions from a processed file"""
        try:
            with open(filename, 'r', encoding='utf-8') as f:
                content = f.read()
                
            questions = []
            current_question = {}
            current_section = None
            
            for line in content.split('\n'):
                line = line.strip()
                
                if line.startswith('<question'):
                    if current_question:
                        questions.append(current_question)
                        current_question = {}
                    continue
                    
                if line.startswith('Introduction:'):
                    current_section = 'Introduction'
                    current_question[current_section] = ''
                elif line.startswith('Conversation:'):
                    current_section = 'Conversation'
                    current_question[current_section] = ''
                elif line.startswith('Question:'):
                    current_section = 'Question'
                    current_question[current_section] = ''
                elif line and current_section:
                    current_question[current_section] = current_question[current_section] + '\n' + line if current_question[current_section] else line
                    
            if current_question:
                questions.append(current_question)
                
            return questions
        except Exception as e:
            print(f"Error parsing questions from {filename}: {str(e)}")
            logger.error(f"Error parsing questions from {filename}: {str(e)}")
            return []
        
        
    def cleanup_old_vectorstore(self):
        """Clean up old vector store files while preserving recent ones"""
        try:
            # Keep only the most recent N days of files
            MAX_DAYS = 7  # Adjust this value as needed
            
            # Get current timestamp for comparison
            import time
            from datetime import datetime, timedelta
            current_time = datetime.now()
            
            # List all files in the vector store directory
            for root, dirs, files in os.walk(self.persist_directory):
                for file in files:
                    file_path = os.path.join(root, file)
                    file_time = datetime.fromtimestamp(os.path.getmtime(file_path))
                    
                    # If file is older than MAX_DAYS, delete it
                    if current_time - file_time > timedelta(days=MAX_DAYS):
                        try:
                            os.remove(file_path)
                            print(f"Removed old vector store file: {file}")
                            logger.info(f"Removed old vector store file: {file}")
                        except Exception as e:
                            print(f"Error removing file {file}: {str(e)}")
                            logger.error(f"Error removing file {file}: {str(e)}")
        except Exception as e:
            print(f"Error during vector store cleanup: {str(e)}")
            logger.error(f"Error during vector store cleanup: {str(e)}")

    def reset_vectorstore(self):
        """Complete reset of vector store"""
        try:
            import shutil
            if os.path.exists(self.persist_directory):
                shutil.rmtree(self.persist_directory)
                print(f"Removed existing vector store directory: {self.persist_directory}")
                logger.info(f"Removed existing vector store directory: {self.persist_directory}")
            os.makedirs(self.persist_directory)
            print(f"Created fresh vector store directory: {self.persist_directory}")
            logger.info(f"Created fresh vector store directory: {self.persist_directory}")
        except Exception as e:
            print(f"Error resetting vector store: {str(e)}")
            logger.error(f"Error resetting vector store: {str(e)}")

if __name__ == "__main__":
    import argparse
    
    parser = argparse.ArgumentParser(description='Vector Store Management')
    parser.add_argument('--reset', action='store_true', help='Reset the entire vector store')
    parser.add_argument('--cleanup', action='store_true', help='Clean up old vector store files')
    args = parser.parse_args()
    
    try:
        store = QuestionVectorStore()
        
        # Print initial state
        print(f"\nChecking vector store directory: {store.persist_directory}")
        if os.path.exists(store.persist_directory):
            files = os.listdir(store.persist_directory)
            print(f"Initial vector store contents: {files}")
        else:
            print("Vector store directory does not exist!")
        
        if args.reset:
            store.reset_vectorstore()
            print("Vector store reset completed.")
            sys.exit(0)
            
        if args.cleanup:
            store.cleanup_old_vectorstore()
            print("Vector store cleanup completed.")
            sys.exit(0)
        
        # Regular processing continues here
        questions_dir = "backend/data/processed_questions"
        if not os.path.exists(questions_dir):
            os.makedirs(questions_dir, exist_ok=True)
        
        question_files = [
            ("backend/data/processed_questions/sY7L5cfCWno_section2.txt", 2),
            ("backend/data/processed_questions/sY7L5cfCWno_section3.txt", 3)
        ]
        
        # Process files
        for filename, section_num in question_files:
            if os.path.exists(filename):
                print(f"\nProcessing {filename}...")
                store.index_questions_file(filename, section_num)
            else:
                print(f"File not found: {filename}")
        
        # Final verification
        print("\nFinal vector store state:")
        if os.path.exists(store.persist_directory):
            files = os.listdir(store.persist_directory)
            print(f"Vector store contents: {files}")
        else:
            print("Vector store directory still does not exist!")
            
    except Exception as e:
        print(f"Error in main execution: {str(e)}")
        logger.error(f"Error in main execution: {str(e)}")
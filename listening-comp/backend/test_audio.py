import sys
import os
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from backend.audio_generator import AudioGenerator

# Test question data
test_question = {
    "Introduction": "次の会話を聞いて、質問に答えてください。",
    "Conversation": """
    男性: すみません、この電車は新宿駅に止まりますか。
    女性: はい、次の駅が新宿です。
    男性: ありがとうございます。何分くらいかかりますか。
    女性: そうですね、5分くらいです。
    """,
    "Question": "新宿駅まで何分かかりますか。",
    "Options": [
        "3分です。",
        "5分です。",
        "10分です。",
        "15分です。"
    ]
}

def test_audio_generation():
    print("Initializing audio generator...")
    generator = AudioGenerator()

    # Test question
    test_question = """
    次の会話を聞いて、質問に答えてください。
    
    田中: すみません、この電車は新宿駅に止まりますか。
    駅員: はい、止まります。でも、次の電車の方が速く着きますよ。
    田中: そうですか。何分後に来ますか。
    駅員: 5分後です。
    
    質問: 田中さんはどうすれば良いですか。
    """

    print("\nParsing conversation...")
    try:
        parts = generator.parse_conversation(test_question)
        print("\nParsed parts:")
        for speaker, text, gender in parts:
            print(f"Speaker: {speaker} ({gender})")
            print(f"Text: {text}\n")

        print("\nGenerating audio...")
        audio_file = generator.generate_audio(test_question)
        print(f"Audio generated: {audio_file}")

    except Exception as e:
        print(f"\nError during test: {str(e)}")

if __name__ == "__main__":
    test_audio_generation()

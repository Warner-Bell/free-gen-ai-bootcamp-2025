import sys
import os
import json
import traceback

# Ensure the backend directory is in the import path
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from backend.audio_generator import AudioGenerator  # Import your audio generator class

def test_audio_generation():
    print("🟢 [INFO] Starting test_audio.py execution...")

    try:
        print("🟡 [INFO] Initializing audio generator...")
        generator = AudioGenerator()

        test_conversation = """
        次の会話を聞いて、質問に答えてください。
        
        田中: すみません、この電車は新宿駅に止まりますか。
        駅員: はい、止まります。でも、次の電車の方が速く着きますよ。
        田中: そうですか。何分後に来ますか。
        駅員: 5分後です。
        
        質問: 田中さんはどうすれば良いですか。
        """

        print("\n🟡 [INFO] Calling `parse_conversation()`...")

        parts = generator.parse_conversation({"Conversation": test_conversation})

        print("\n🔍 [DEBUG] `parse_conversation()` Output:")
        if parts:
            # Print each conversation part on a new line
            print("[")
            for speaker, text, gender in parts:
                print(f"  {{")
                print(f'    "speaker": "{speaker}",')
                print(f'    "text": "{text}",')
                print(f'    "gender": "{gender}"')
                print("  }},")
            print("]")
        else:
            print("[]")

        if parts is None:
            print("\n🚨 [ERROR] `parse_conversation()` returned None!")
            return

        if not isinstance(parts, list) or len(parts) == 0:
            print("\n🚨 [ERROR] `parse_conversation()` returned an empty response!")
            print(f"Returned data type: {type(parts)}")
            return

        print("\n✅ [INFO] Successfully parsed conversation!")

        print("\n🔵 [INFO] Processing parsed conversation parts:")
        for part in parts:
            speaker = "Unknown"
            text = "No text"
            gender = "Unknown"
            
            if not isinstance(part, tuple) or len(part) != 3:
                print(f"⚠️ [WARNING] Invalid part format: {part}")
                continue

            if isinstance(part, dict) and "speaker" in part:
                speaker = part.get("speaker", "Unknown")
                text = part.get("text", "No text")
                gender = part.get("gender", "Unknown")
            elif isinstance(part, tuple):
                if len(part) == 3:
                    speaker, text, gender = part
                elif len(part) == 2:
                    speaker, text = part
                    gender = "Unknown"
                else:
                    speaker, text, gender = "Unknown", str(part), "Unknown"
            else:
                print(f"⚠️ [WARNING] Unexpected data format: {part}")
                speaker, text, gender = "Unknown", str(part), "Unknown"

            print(f"🗣 Speaker: {speaker} ({gender})")
            print(f"📝 Text: {text}\n")

        print("\n🎵 [INFO] Calling `generate_audio()`...")
        audio_file = generator.generate_audio({"Conversation": test_conversation})

        if audio_file:
            print(f"✅ [SUCCESS] Audio successfully generated: {audio_file}")
        else:
            print("❌ [ERROR] Audio generation failed.")

    except Exception as e:
        print("\n🚨 [ERROR] An unexpected error occurred:")
        traceback.print_exc()  # <-- Print full error traceback

if __name__ == "__main__":
    print("🟢 [INFO] Running test_audio_generation() function...\n")
    test_audio_generation()

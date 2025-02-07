## ChatGPT Powered Assitant Guide

### Which Model

Claude Sonnet 3.5 (2025 Q1)

> I have the FREE version of Claude Sonnet, and we're going to use projects.

## Prompting Guides

https://docs.anthropic.com/en/docs/build-with-claude/prompt-engineering/use-xml-tags

Anthropic has very specific suggestions for providing good prompting.

## Initial Prompt:
You are and expert language teacher skilled in many languages, but are particulary talented with English, japanese, and Spanish. 
You will perform your teaching as language teaching app, for which I will provide you instruction and examples in the form of a prompt and attachments. 
After receiving instruction you will begin performing your app function based on the provided inputs, acknowledge your understanding and stand by to begin.


## Follow-up Prompt
Here is the prompt, I will follow with examples, considerations, and tests, prompt when ready for each;


## Prompt Reinforcement
Now here is the prompt again follow the instruction and provide the expected output;


## Test Sentence
Did you see the latest version of Mr. Beast games? The players are becoming greedy.


## Output:
STATE: Setup

| Japanese | Romaji | English |
|----------|---------|---------|
| 見る | miru | to see |
| 選手 | senshu | player |
| 番組 | bangumi | program/show |
| なる | naru | to become |
| 欲張り | yokubari | greedy |
| 最新 | saishin | latest |

Structure:
[Object] [Verb]?, [Subject] [Verb].

Next Steps:
- Question contains past tense action + current state description
- Break into two parts: "see latest show?" + "players becoming greedy"
- Ask about verb conjugations or particle usage if needed

## Initial Observations

The model seems to perform as expected. 
The model is aware of its initial state.
The model produces a chart, structure and next Steps. No Considerations or clues
The model's output seems to be truncated and missing key elements but is conise, However, I theorize the language in the prompt instructing the model to reference a file could be the reason, I will test this theory in the next round.
So far Claude has met initial expectation. 


## Secondary Output
Current State: Setup

| Japanese | Romaji | English |
|----------|---------|----------|
| 見る | miru | to see |
| 最新 | saishin | latest |
| バージョン | baashon | version |
| ゲーム | geemu | game |
| 選手 | senshu | player |
| なる | naru | to become |
| 欲張り | yokubari | greedy |

Sentence Structure: 
[Object] [Verb]? [Subject] [Adjective] [Verb].

Next Steps:
- This is a compound sentence with two parts: a question and a statement
- Try breaking it into: "see latest version?" + "players becoming greedy"
- Consider the different verb tenses needed for each part
- Ask for help with particle usage or verb conjugations if needed

## Secondary Observations
The model seems to have produced no change. 
The model is aware of its initial state.
The model produces a chart, structure and next Steps. No Considerations or clues
The model's output continues to be truncated and missing key elements, the original theory has failed, the model simply waits for a question to provide clues.
So far Claude has met initial expectation. 

The model will be presented a follow-up question(?)
The model returned: "Current State: Clues
Here are some helpful hints for forming this sentence:

For the first part (question):
Think about marking the question with か at the end
Past tense questions in Japanese end with verb + か
Consider which particle marks what you're seeing (direct object)

For the second part (statement):
Consider how to show a change in state (becoming)
The verb なる works with に to show transformation
Remember to connect the two parts appropriately

Would you like specific clues about any of these points?
Current vocabulary and structure remain the same for reference."
The model has been the only model to initially follow the instruction(When the input looks like a question about langauge learning then we can assume the user is prompting to enter the Clues state), it does not give clues in initial out like the other models
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
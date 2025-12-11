# Straico-Cli
![Hero image of straico-cli output asking about the history of the French Quarter](https://github.com/user-attachments/assets/62d68bb4-daf5-42ab-9354-fd1950ad5e0f)

## Usage

Once in straico-cli, the text placeholder will tell you:
- The current LLM model, 
- Scroll percentage 
- Current buffer.
- Session coin usage
```text
┃ Ask the LLM... (openai/gpt-4o-mini) (%100) (1) (1.23)
┃
┃
```

The following actions are available:
- Buffer Switching: Press `F1` - `F9`
- Buffer Erase: Press `F12`
- Buffer Move: Press `Shift + Right Arrow` or `Shift + Left Arrow`.  
  For example, if you have a buffer at location `1` and want to move it to `2`, press `F1`, `Shift + Right Arrow`

```bash
Usage of straico-cli:
      --file-url strings      --file-url link1 --file-url link2
  -l, --list-models           List models
  -m, --model string          Model to use (default "anthropic/claude-3-haiku:beta")
      --save-key string       Straico API key
      --save-model            Use the model listed by -m for future queries
      --youtube-url strings   --youtube-url link1 --youtube-url link2
```

### Save your [API key](https://documenter.getpostman.com/view/5900072/2s9YyzddrR)
```bash
straico-cli --save-key YourAPIKey123
```

### Save your model
```bash
straico-cli --save-model -m "anthropic/claude-3-haiku:beta" 
```

## Resources
- [Models](https://straico.com/multimodel/)
- [API Doc - Getting API Key](https://documenter.getpostman.com/view/5900072/2s9YyzddrR)

# Straico-Cli
![Hero image of straico-cli output asking about the history of the French Quarter](https://github.com/user-attachments/assets/076050be-87c8-4bea-985d-7e25ec625400)

## Usage

Once in straico-cli, the text placeholder will tell you the current LLM model, 
 the scroll percentage and the current buffer.
```text
┃ Ask the LLM... (openai/gpt-4o-mini) (%100) (1)
┃
┃
```

To switch between buffers, push `F1` - `F8`. Pushing `F12` will erase the current buffer.

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
straico-cli --save-key Xl-YourAPIKey123
```

### Save your model
```bash
straico-cli --save-model -m "anthropic/claude-3-haiku:beta" 
```

## Resources
- [Models](https://straico.com/multimodel/)
- [API Doc - Getting API Key](https://documenter.getpostman.com/view/5900072/2s9YyzddrR)

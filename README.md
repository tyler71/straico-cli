# Straico-Cli
![Hero image of straico-cli output asking about the history of the French Quarter](https://github.com/user-attachments/assets/076050be-87c8-4bea-985d-7e25ec625400)

## Usage

```bash
Usage of straico-cli:
      --file-url strings      --file-url link1 --file-url link2
  -l, --list-models           List models
  -m, --model string          Model to use (default "anthropic/claude-3-haiku:beta")
      --save-key string       Straico API key
      --save-model            Use the model listed by -m for future queries
      --youtube-url strings   --youtube-url link1 --youtube-url link2
```

```bash
straico-cli 'What is the history behind the French Quarter?'
```

### Save your [API key](https://documenter.getpostman.com/view/5900072/2s9YyzddrR)
```bash
straico-cli --save-key Xl-YourAPIKey123
```

## Resources
- [Models](https://straico.com/multimodel/)
- [API Doc - Getting API Key](https://documenter.getpostman.com/view/5900072/2s9YyzddrR)

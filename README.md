<a id="readme-top"></a>

<!-- PROJECT SHIELDS -->
[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![project_license][license-shield]][license-url]
[![Go][go-shield]][go-url]



<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://mergram.dev">
    <img src="assets/logo.jpg" alt="MerGram Logo" width="200" height="200">
  </a>

  <h3 align="center">Mermaid Diagram Bot for Telegram</h3>

  <p align="center">
    A Telegram bot that automatically detects Mermaid diagram syntax in group conversations and converts it to rendered images, making technical discussions more visual and accessible
    <br />
    <a href="https://mergram.dev"><strong>Visit Website »</strong></a>
    <br />
    <br />
    <a href="https://t.me/MerGramBot">Try the Bot</a>
    &middot;
    <a href="https://github.com/capcom6/mergram-tg-bot">View Source</a>
    &middot;
    <a href="https://github.com/capcom6/mergram-tg-bot/issues/new?labels=bug&template=bug-report---.md">Report Bug</a>
    &middot;
    <a href="https://github.com/capcom6/mergram-tg-bot/issues/new?labels=enhancement&template=feature-request---.md">Request Feature</a>
  </p>

  <img src="assets/banner.jpg" alt="MerGram Banner">
</div>



<!-- TABLE OF CONTENTS -->
- [About The Project](#about-the-project)
  - [Key Features](#key-features)
  - [Why MerGram?](#why-mergram)
  - [Built With](#built-with)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Docker Installation (Alternative)](#docker-installation-alternative)
- [Usage](#usage)
  - [Auto-Detection (Group Messages)](#auto-detection-group-messages)
  - [Manual Commands](#manual-commands)
  - [Bot Commands](#bot-commands)
  - [Supported Diagram Types](#supported-diagram-types)
- [Contributing](#contributing)
  - [Top contributors:](#top-contributors)
- [License](#license)
- [Contact](#contact)
- [Acknowledgments](#acknowledgments)


<!-- ABOUT THE PROJECT -->
## About The Project

MerGram is a Telegram bot designed to enhance technical discussions by automatically converting Mermaid diagram syntax into rendered images directly within group conversations. Whether you're sharing system architectures, flowcharts, or sequence diagrams, MerGram eliminates the need for manual copy-pasting to external renderers.

**Powered by mermaid.ink**: The bot uses mermaid.ink service under the hood for rendering, ensuring reliable and high-quality diagram generation.

### Key Features

- **🔍 Auto-Detection**: Automatically detects Mermaid code blocks in group messages
- **🎯 Smart Rendering**: Converts syntax to PNG images for seamless viewing
- **⚡ Real-time Processing**: Fast diagram generation with error handling
- **🔒 Privacy-Focused**: No message storage
- **📊 Multi-Diagram Support**: All standard Mermaid diagram types supported

### Why MerGram?

Technical teams using Telegram often share Mermaid diagram code, but recipients must manually copy-paste to renderers. This disrupts conversation flow and creates friction in collaborative diagramming. MerGram solves this by providing instant visual feedback, making technical discussions more engaging and accessible.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



### Built With

* [![Go][Go]][Go-url]
* [![Telegram Bot API][Telegram]][Telegram-url]
* [![Mermaid][Mermaid]][Mermaid-url]
* [![mermaid.ink][mermaid-ink]][mermaid-ink-url]
* [![Docker][Docker]][Docker-url]
* [![GitHub Actions][GitHub-Actions]][GitHub-Actions-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Getting Started

This guide will help you get MerGram up and running on your local machine for development and testing purposes.

### Prerequisites

- **Go 1.25.5+**: Download from [golang.org](https://golang.org/dl/)
- **Git**: For version control
- **Telegram Bot Token**: Create a bot via [@BotFather](https://t.me/BotFather)

### Installation

1. **Clone the repository**
   ```sh
   git clone https://github.com/capcom6/mergram-tg-bot.git
   cd mergram-tg-bot
   ```

2. **Install dependencies**
   ```sh
   go mod download
   ```

3. **Environment Setup**
   Create a `.env` file in the root directory:
   ```bash
   TELEGRAM__TOKEN=your_bot_token_here
   RENDERER__TIMEOUT=10s
   ```

4. **Build the application**
   ```sh
   make build
   ```

5. **Run the bot**
   ```sh
   make run
   ```

### Docker Installation (Alternative)

```bash
# Run with environment variables
docker run -e TELEGRAM__TOKEN=your_token_here ghcr.io/capcom6/mergram-tg-bot:latest
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- USAGE EXAMPLES -->
## Usage

Once the bot is running and added to your Telegram group, here's how to use MerGram:

### Auto-Detection (Group Messages)

Simply paste Mermaid code blocks in your group chat:

````text
  User: Here's our system architecture:

  ```mermaid
  graph TD
      A[Client] --> B[Load Balancer]
      B --> C[Server 1]
      B --> D[Server 2]
      C --> E[Database]
      D --> E
  ```

  Bot: (Automatically replies with rendered diagram image)
````

### Manual Commands

**Generate a diagram from inline code:**
```
/mermaid sequenceDiagram
    Alice->>John: Hello John, how are you?
    John-->>Alice: I am good thanks!
```


### Bot Commands

| Command           | Description                            | Example                     |
| ----------------- | -------------------------------------- | --------------------------- |
| `/start`          | Welcome message and basic instructions | `/start`                    |
| `/help`           | Detailed help with examples            | `/help`                     |
| `/mermaid <code>` | Create diagram from code               | `/mermaid graph TD; A-->B;` |

### Supported Diagram Types

**Powered by mermaid.ink**: All diagram types supported by mermaid.ink are available through the bot, ensuring comprehensive coverage of Mermaid syntax.

*For the complete list of supported diagram types and syntax, refer to the [mermaid.ink documentation](https://mermaid.ink/).*

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Top contributors:

<a href="https://github.com/capcom6/mergram-tg-bot/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=capcom6/mergram-tg-bot" alt="contrib.rocks image" />
</a>



<!-- LICENSE -->
## License

Distributed under the Apache License 2.0. See `LICENSE` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTACT -->
## Contact

**Project Maintainer**: Capcom6

- Website: [mergram.dev](https://mergram.dev)
- Try the bot: [@MerGramBot](https://t.me/MerGramBot)
- GitHub: [@capcom6](https://github.com/capcom6)
- Project Link: [https://github.com/capcom6/mergram-tg-bot](https://github.com/capcom6/mergram-tg-bot)
- Issues: [Report Bug](https://github.com/capcom6/mergram-tg-bot/issues) | [Request Feature](https://github.com/capcom6/mergram-tg-bot/issues)

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- ACKNOWLEDGMENTS -->
## Acknowledgments

* **[go-telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api)** - Excellent Go SDK for Telegram Bot API
* **[Mermaid](https://mermaid.js.org/)** - Amazing diagramming and charting tool that works with markdown
* **[mermaid.ink](https://mermaid.ink/)** - Reliable online Mermaid diagram rendering service
* **[Go Community](https://golang.org/)** - For the robust and efficient Go programming language
* **[Telegram](https://telegram.org/)** - For providing a powerful Bot API platform
* **[Docker](https://www.docker.com/)** - For simplified deployment and containerization

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/capcom6/mergram-tg-bot.svg?style=for-the-badge
[contributors-url]: https://github.com/capcom6/mergram-tg-bot/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/capcom6/mergram-tg-bot.svg?style=for-the-badge
[forks-url]: https://github.com/capcom6/mergram-tg-bot/network/members
[stars-shield]: https://img.shields.io/github/stars/capcom6/mergram-tg-bot.svg?style=for-the-badge
[stars-url]: https://github.com/capcom6/mergram-tg-bot/stargazers
[issues-shield]: https://img.shields.io/github/issues/capcom6/mergram-tg-bot.svg?style=for-the-badge
[issues-url]: https://github.com/capcom6/mergram-tg-bot/issues
[license-shield]: https://img.shields.io/github/license/capcom6/mergram-tg-bot.svg?style=for-the-badge
[license-url]: https://github.com/capcom6/mergram-tg-bot/blob/main/LICENSE
[go-shield]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white

<!-- Shields.io badges. You can a comprehensive list with many more badges at: https://github.com/inttter/md-badges -->
[Go]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white
[Go-url]: https://golang.org/
[Telegram]: https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white
[Telegram-url]: https://telegram.org/
[Mermaid]: https://img.shields.io/badge/Mermaid-FF3670?style=for-the-badge&logo=mermaid&logoColor=white
[Mermaid-url]: https://mermaid.js.org/
[Docker]: https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white
[Docker-url]: https://www.docker.com/
[GitHub-Actions]: https://img.shields.io/badge/GitHub_Actions-2088FF?style=for-the-badge&logo=github-actions&logoColor=white
[GitHub-Actions-url]: https://github.com/features/actions
[mermaid-ink]: https://img.shields.io/badge/mermaid.ink-FF3670?style=for-the-badge&logo=mermaid&logoColor=white
[mermaid-ink-url]: https://mermaid.ink/
[go-url]: https://golang.org/ 
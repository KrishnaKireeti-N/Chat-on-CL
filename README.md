# Chat-on-CL
- A personal project for emulating communication (similar to famous chatting-apps) over the command line
- IMPORTANT: The current version send raw-text over tcp and isn't secure at all

### How-To-Use
- Build the project with `go build ./cmd/chat` (to pre-download all dependencies without building, use `go mod download`)
- To use:
```sh
# In one terminal (this is not a daemon so you need the terminal open)
$ ./chat -s <username>
# Listens on localhost:23456
# <username> is your username which you want to be known with!

# In some other terminal (possible other machine)
$ ./chat -c localhost <username>
# <username> is your username which you want to be known with!
# Connects to localhost:23456 and starts the UI for chatting

# Use localhost for any testing till security updates
```

### TODO
### Security
- Add a ssl/tls encrypted channel over the tcp for secure communication
- Currently we only need a IP address (and technically port but it is predefined in the `config` struct) to connect to a host and it is entirely possible that the end-user might not be someone we wish to chat with (and more importantly let them pipe malicious go commands, which is because of my poor handling of messages).
	- So, I want to add means of verifying a user with `public-key cryptography` and `challenge-response authentication`, I totally did not copy these words from archwiki page of ssh keys :). Basically add a **key-based authentication** feature
### Features + Code
- Fix the issue when quitting the application
- Change the Send and Recieve methods (they are horrendous right now)
- Change how the current strucutre of messeages are being transmitted (raw-content -> header + body)
- Add customization features like message-color... (yea I am not creative enough to list down more)


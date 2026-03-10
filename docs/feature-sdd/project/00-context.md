## Build a new CLI Project as a command catalog

I would like create a CLI project to run commands in a simple way.

## Requirements

1. I would like a catalog of commands to run on linux primary.
  - Example: To kill a process that is running in a specific port. I uses the command `sudo kill -9 $(sudo lsof -t -i:3040)`.
  - What I want:
    - Create personalized commands like `cs create "kill port $port" "sudo kill -9 $(sudo lsof -t -i:$port)"` where will save the command in a local file.
    - Now when I run  `cs kill port 3040`, under the hood will run `sudo kill -9 $(sudo lsof -t -i:$port)` defining the params.
    - Above is a example, I would like customize any command with params. 
    - To store the data, save it in a local file with the structure of id:key:value.
    - Command to list the catalog `cs list`. Fields (id, key, value)
    - Command to delete one command `cs delete $id`.
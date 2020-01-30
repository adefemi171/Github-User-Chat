
# Welcome to Github User Chatroom

Creating a Github Chat for user and Private Chat room among users using Golang. This app basically allow users to interact between themselves.

# How to setup project and run locally

## Clone the repository

> git clone https://github.com/adefemi171/
> cd into ``
> then cd into `chat` and run `go build -o githubuserChat`

# Top-level directory layout

    📦githubChat
        📦chat
            ┣ 📜client.go
            ┣ 📜githubuserChat(Built file)
            ┣ 📜main.go
            ┗ 📜room.go
        📦templates
            ┗ 📜test.html
        📦trace
            ┣ 📜tracer.go
            ┗ 📜tracer_test.go
        ┣ 📜README.md

## Test Driven Development Approach was used

Code for Unit test and Red/Green test can be found in 📦trace subdirectory

# Running Trace

Uncomment `r.tracer = trace.New(os.Stdout)` in `main.go` under the 📦chat subdirectory and also comment `tracer:  trace.Off(),` in `room.go` under the 📦chat subdirectory and run `go build -o githubuserChat` and then run `./githubuserChat`, watch your terminal to see the output of the trace information

# NOTE: `githubuserChat` is just the output name it could be changed to any name of your choice
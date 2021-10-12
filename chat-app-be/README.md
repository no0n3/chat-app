## Description
A simple chat app backend REST API with realtime chatting (using websockets) and optional sending images in chat messages.

The main goal is to test how a chat system would work on a distributed system running on more than one Node (server) and to get used to Golang.

## Run locally

* `export CHAT_APP_ENDPOINT=http://localhost:81`
* `docker build . -t chat-app-be --build-arg CHAT_APP_ENDPOINT=$CHAT_APP_ENDPOINT`
* `docker-compose up`
* create database `chat_app`
* `psql postgresql://postgres:example@localhost:5433/chat_app < ./db/db.sql`
* `psql postgresql://postgres:example@localhost:5433/chat_app < ./db/seed.sql`

## App architecture

![alt text](./images/chat-diagram.drawio.png)

In this demo there are 3 nodes created (docker containers). There is a load balancer (nginx) that distributs the traffic among the 3 nodes.

Every Node holds a map of websocket connections in format `{<userId>: [<list of connections>], ...}`.

As the load is distributed among the 3 nodes there's the possibility that a user has several websocket connections among the nodes. For example:

user1 has 2 websocket connections:
1 on BE Node1 and 1 on BE Node2

user2 has 1 websocket connections:
1 on BE Node3

When user2 sends a message to user1 the message needs to be sent to every websocket connection of user1 and user2.

When user1 sends a message to user2 or vice versa:
1) created message
2) store message in database
3) get user ids of the chat members of the chat (user1 and user2 in this case)
4) publish a message to the topics `user/<user1_id>` and `user/<user2_id>`
5) all nodes that have active websocket connections are subscribed to the topics with the user ids
6) when the a message is published to `user/<user1_id>` and `user/<user2_id>` topics all 3 Nodes (as all 3 users have websocket connections on them) will be botified and receive the chat message from the topics
7) send response to all user1 and user2 connections on all Nodes

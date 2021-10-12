# Simple chat app

## Description

A demo chat app when you chan sign up and login, add contacts and chat realtime with them and send images with chat messages.

The consists of client SPA (built with React) and a backend REST API (build with golang). The backend is dockerized and 3 instances of the backend are spawned and the all requests are load balanced to the 3 nodes. The idea is to see how a distributed chat system would work.

The main goal of this app is get used to React and golang.

For more info about each project - check out `chat-app-spa` and `chat-app-be` repos.

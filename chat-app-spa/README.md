# Simple chat app SPA using React

## Description

SPA (Single Page App) app that lets you sign up, browse users, add them as contact, realtime chat with your contacts and optionally send images in chat messages.

The main goal of this app is to get used to React.

## Run localy

* `export REACT_APP_WS_ENDPOINT=ws://localhost:81`
* `export REACT_APP_ENDPOINT=http://localhost:81`
* `npm start`

## Build and serve

* `export REACT_APP_WS_ENDPOINT=ws://localhost:81`
* `export REACT_APP_ENDPOINT=http://localhost:81`
* `npm run build`
* `http-server ./build/`

## Features

* sign up (with Email, Name, Username, Password)
* login (with Email and Password)
* logout
* browse all users
    * you can filter by Name and Username
    * user item has `Add contact` or `Remove contact` button
    * clicking on user leads to his profile
* browse your contacts
    * user item has `Remove contact` button
    * clicking on user leads to his profile
* view user profile. Profile has only profile image, name, username, description and add contact button
* if you're viewing your own profile then you can change profile image which opens a popup to crop it
* messages page
    * shows your chat history with your contacts
    * you can filter by Name and Username
    * shows last message
    * clicking on user leads to the chat
* chat page
    * chat realtime with a user
    * you can upload images and send them
    * if a message has image and it's clicked then a gallery popup is opened

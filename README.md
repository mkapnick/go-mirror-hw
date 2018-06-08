## Running locally
- `cd server`
- `GOOS=linux go build -o main`
- `docker run -p 9000:9000 go-mirror`

## Deployment
- Deployed on Kubernetes/GCP
- Can be found here: [http://35.232.131.106](http://35.232.131.106)

## DockerHub
- Can be found on docker hub here: [https://hub.docker.com/r/mdotm/go-mirror](https://hub.docker.com/r/mdotm/go-mirror)
- `docker pull mdotm/go-mirror`
- `docker run -p 9000:9000 mdotm/go-mirror`

## Context
- Every user has a `userId` and a `channelId`
- The user's `channelId` is a public channel. Anyone can subscribe to
that channel and write to it. Upon login, I automatically subscribe the
logged in user to their channelId so they can see any messages coming in on
that channel from other users
- The user's `userId` is used to DM a user.

### The flow of how DM's work
- User logs in and a websocket connection is established with the
server
- Once the websocket connection is established, the server sends an
acknowledgment (`ACK`) to the client which basically says: hey, you are now
registered correctly on the server, here is your `websocketId`.
- A user initiates a DM request to another user.
- The client sends a message to the server over the websocket connection that
includes the current `channelId`, the client's `userId`, and the
client's `websocketId`, and the message `code`
- A DM channel is 2 user id's concatanted together by a `.`. So a DM channel
looks like `userId1.userId2`. The 2nd userId (`userId2`) is the receiving user.
The first userId (`userId1`) is the sending user.
- The server receives this message and recognizes that it's a DM request from
the code.
- The server loops through all the websocket clients, and checks to see if
the client userId matches `userId2`. If it matches then the server knows
to relay this message to that client over that particular websocket connection.
Te server then sends a message over the ws connection with a specific code
`subscribe.direct`.
- The client listens for incoming messages on the ws, sees the incoming message
of type `subscribe.direct`, and displays the message to that user.


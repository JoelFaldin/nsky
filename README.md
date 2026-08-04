# nsky

A minimal recreation of ![ngrok](https://ngrok.com/) internals. Implemented using ![golang](https://go.dev/). It allows to expose a locally running app (for example, an app running on `:3000`) towards internet. No public IP or port router forwarding needed.

_This project is a part of my systems and networks learning journey._

## How does it work

The project is separated in two components: a *server* (with public IP, or accesible from the local network) and a *client* (runs alongside the local service that you want to expose).

The server exposes 3 ports with different duties:

| Port | Role | Description |
| ---- | ---- | ----------- |
| `:4443` | Control | Persistent connection between client and server. Used only for coordination. |
| `:4444` | Join | Port where the client opens a new TCP connection when the server needs to forward a visitor's traffic. |
| `:8080` | Public | Port that captures real traffic (for example, a `curl` request, or a browser). |

## Full Flow

1. A visitor makes a request to `:8080` (via `curl` or the browser).
2. The server generates a unique `id` for that request, and registers it in a pending connections map. It is associated to a go `channel`.
3. The server notifies the client, using the `control` channel, that it needs a new data connection (`{ "type": "new_conn", "id": "..." }`).
4. The client receives the notification, opens a new connection towards the join port (`:4444`), and informs `id` as the first message.
5. The server receives that connection, searches for the `id` on the `pending` map, and delivers the connections to the corresponding channel.
6. The server _proxies_ the original visitor's connection, copying bytes in both directions.
7. At the same time, the client connects that same connection with a local connection towards the exposed service (`localhost:3000`), also copying bytes in both directions.

## Project Arquitecture

.
├── go.mod
└── internal/
    ├── client/
    │   └── main.go
    ├── server/
    │   └── main.go
    ├── protocol/
    │   └── protocol.go
    └── utils/
        ├── pipe.go
        └── proxy.go

## Requirements

* Golang 1.22 or higher.
* A local service running in the port you want to expose (as default, this project uses `localhost:3000`).

## Running the Project Locally

1. Clone the repo and download dependencies:

```
go mod download
```

2. Run the server:

```
go run ./internal/server/main.go
```

3. Run the client:

```
go run ./internal/client/main.go
```

4. Run the service you want to expose:
_(For example, a Nextjs app)_

```
npm run dev
```

5. Make a request:

```
curl -v localhost:8080
```

## Final considerations

This is merely a learning project, used to learn networking concepts and concurrent programming in golang.

(Dont ask me to explain the inners of the project, as I won't be able to!!)

# Chat App Workshop Playbook

This playbook outlines the feature requests and implementation steps for the Chat App project that user asks Antigravity to do

## 1. Multi-User Support
- **Backend:** "Update to handle chat sessions from multiple users simultaneously."
- **Frontend:** "Implement a username prompt before users join a session."

## 2. Communication Protocol Optimization
"Evaluate the most optimal way to handle chat messages (beyond standard GET/POST)."
### Response summary:
  - **WebSockets:** The "Gold Standard" for real-time, bidirectional communication.
  - **Server-Sent Events (SSE):** Suitable for server-to-client streaming.
  - **gRPC-Web / Protobuf:** Performance-oriented binary protocol.

## 3. WebSocket Implementation
"Implement WebSockets to replace standard HTTP polling for real-time message delivery."

## 4. Redis Persistence
- "Transition from in-memory message storage to **Redis**."
- "Include a `docker-compose.yaml` to orchestrate a local Redis cluster for development."

## 5. Backend Containerization
"Create a `Dockerfile` for the Go backend to enable deployment to Google Cloud Run."

## 6. Cloud Infrastructure
- **Memorystore:** "Explain how to set up a managed Redis instance on Google Cloud."
- **Connectivity:** "Explain how to configure the Cloud Run backend to securely connect to Memorystore I just setup."

## 7. Dynamic Frontend Configuration
"Implement a mechanism to handle dynamic backend URLs in the Flutter frontend during the build process."

---
## You can check the working final version in `post-workshop-version/` folder.

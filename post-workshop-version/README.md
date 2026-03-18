# Chat App Workshop with Antigravity IDE

Welcome to the Chat App Workshop! In this session, you'll be using Google's Antigravity IDE to build a real-time chat application with a Go backend and a Flutter frontend.

## Prerequisites
- [Antigravity IDE](https://antigravity.google) installed.
- Flutter SDK installed and configured.
- Go installed.
- Docker and Docker Compose installed.

## Getting Started

### 1. Backend & Redis Setup
The project uses Docker Compose to run the Go backend and a Redis instance simultaneously.

1. From the project's root directory, start the infrastructure:
   ```bash
   docker compose up --build -d
   ```
   The backend server will start and become available at `http://localhost:8080`.
   *(Optional: The backend is configured to respect `PORT` and `REDIS_URL` environment variables if you deploy or run it manually).*

### 2. Frontend Setup
1. Navigate to the `frontend` directory.
2. In your environment or configuration, set `BACKEND_URL` to point to your server (e.g., `http://localhost:8080` for local testing).
3. Fetch dependencies:
   ```bash
   flutter pub get
   ```
4. Run the app:
   ```bash
   flutter run -d chrome
   ```
   (Or use your preferred device/emulator)

## Vibe Coding with Skills
The `.agents/skills` directory contains specialized instructions for Antigravity. You can use these to accelerate your development:

- **"Refresh the vibe"**: Ask Antigravity to improve the Flutter UI.
- **"Harden the backend"**: Ask Antigravity to add security and robustness to your Go server.
- **Describe a new feature**: Use natural language to ask for new capabilities like "Add message emojis".

## Workshop Goals
1. **Explore**: Understand the boilerplate code.
2. **Experiment**: Use Antigravity to modify the UI and backend
3. **Build**: Make the chat multi-user. The boilerplate code doesn't support that.
4. **Enhance**: Create your own "skills" in `skills.md` to automate repetitive tasks or complex logic.
5. **Improve**: Make messages persist somewhere better than in the backend memory - `[]Message{}`.
6. **Deploy**: Deploy the app to a cloud platform.
7. **Secure**: Add authentication and authorization to the app
8. **Expand**: Make an Android app build!

### **Happy Coding with Google Antigravity IDE!**
[Antigravity IDE](https://antigravity.google)

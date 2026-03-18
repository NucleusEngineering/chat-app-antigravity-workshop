import 'package:flutter/material.dart';
import 'dart:convert';
import 'dart:async';
import 'package:http/http.dart' as http;
import 'package:web_socket_channel/web_socket_channel.dart';

void main() {
  runApp(const ChatApp());
}

class ChatApp extends StatelessWidget {
  const ChatApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Workshop Chat',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.deepPurple),
        useMaterial3: true,
      ),
      initialRoute: '/',
      routes: {
        '/': (context) => const JoinScreen(),
      },
      onGenerateRoute: (settings) {
        if (settings.name == '/chat') {
          final username = settings.arguments as String;
          return MaterialPageRoute(
            builder: (context) => ChatScreen(username: username),
          );
        }
        return null;
      },
    );
  }
}

class JoinScreen extends StatefulWidget {
  const JoinScreen({super.key});

  @override
  State<JoinScreen> createState() => _JoinScreenState();
}

class _JoinScreenState extends State<JoinScreen> {
  final TextEditingController _usernameController = TextEditingController();

  void _joinSession() {
    if (_usernameController.text.isNotEmpty) {
      Navigator.pushReplacementNamed(
        context,
        '/chat',
        arguments: _usernameController.text,
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Join Chat'),
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Text(
              'Welcome to Workshop Chat!',
              style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 32),
            TextField(
              controller: _usernameController,
              decoration: const InputDecoration(
                labelText: 'Enter your username',
                border: OutlineInputBorder(),
              ),
              onSubmitted: (_) => _joinSession(),
            ),
            const SizedBox(height: 16),
            ElevatedButton(
              onPressed: _joinSession,
              child: const Text('Join'),
            ),
          ],
        ),
      ),
    );
  }
}

class ChatScreen extends StatefulWidget {
  final String username;
  const ChatScreen({super.key, required this.username});

  @override
  State<ChatScreen> createState() => _ChatScreenState();
}

class _ChatScreenState extends State<ChatScreen> {
  final TextEditingController _controller = TextEditingController();
  final List<dynamic> _messages = [];
  // Use String.fromEnvironment to allow setting the backend URL during build/run.
  // Example: flutter run --dart-define=BACKEND_URL=http://your-backend-url.com
  static const String _backendBase = String.fromEnvironment(
    'BACKEND_URL',
    defaultValue: 'http://localhost:8080',
  );

  final String _historyUrl = '$_backendBase/messages';
  final String _wsUrl = _backendBase.replaceFirst('http', 'ws') + '/ws';
  WebSocketChannel? _channel;

  @override
  void initState() {
    super.initState();
    _fetchHistory();
    _connectWebSocket();
  }

  void _connectWebSocket() {
    _channel = WebSocketChannel.connect(Uri.parse(_wsUrl));
    _channel!.stream.listen((message) {
      setState(() {
        final decoded = json.decode(message);
        // Only add if not already present (avoid duplicates if history and WS collide)
        if (!_messages.any((m) => m['id'] == decoded['id'])) {
          _messages.add(decoded);
        }
      });
    }, onError: (error) {
      debugPrint('WebSocket error: $error');
      // Simple reconnection logic
      Future.delayed(const Duration(seconds: 2), _connectWebSocket);
    }, onDone: () {
      debugPrint('WebSocket closed');
    });
  }

  @override
  void dispose() {
    _channel?.sink.close();
    super.dispose();
  }

  Future<void> _fetchHistory() async {
    try {
      final response = await http.get(Uri.parse(_historyUrl));
      if (response.statusCode == 200) {
        if (mounted) {
          setState(() {
            _messages.clear();
            _messages.addAll(json.decode(response.body));
          });
        }
      }
    } catch (e) {
      debugPrint('Error fetching history: $e');
    }
  }

  Future<void> _sendMessage() async {
    if (_controller.text.isEmpty) return;

    final text = _controller.text;
    _controller.clear();

    try {
      // Instead of HTTP POST, we can either use WS or stay with POST.
      // Since we want full WS, let's use the REST API for sending to trigger the broadcast
      // OR send via WS. The backend Hub broadcasts to everyone.
      // If we send via WS, we need a Reader in the backend.
      // Current backend only broadcasts when POSTed to /messages or when Hub gets a message.
      // Let's keep POST for simplicity and WS for real-time receiving,
      // or update backend to read from WS too.

      // Let's stick with POST for sending to keep it consistent with the profanity filter
      // and history storage, but the response will be received via WS.
      final response = await http.post(
        Uri.parse(_historyUrl),
        headers: {'Content-Type': 'application/json'},
        body: json.encode({'sender': widget.username, 'text': text}),
      );

      if (response.statusCode != 201) {
        debugPrint('Failed to send message: ${response.statusCode}');
      }
    } catch (e) {
      debugPrint('Error sending message: $e');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Chatting as ${widget.username}'),
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
        actions: [
          IconButton(
            icon: const Icon(Icons.history),
            onPressed: _fetchHistory,
          ),
          IconButton(
            icon: const Icon(Icons.logout),
            onPressed: () => Navigator.pushReplacementNamed(context, '/'),
          ),
        ],
      ),
      body: Column(
        children: [
          Expanded(
            child: ListView.builder(
              itemCount: _messages.length,
              itemBuilder: (context, index) {
                final msg = _messages[index];
                final isMe = msg['sender'] == widget.username;
                return ListTile(
                  title: Text(
                    msg['sender'] ?? 'Unknown',
                    style: TextStyle(
                      fontWeight: isMe ? FontWeight.bold : FontWeight.normal,
                      color: isMe ? Colors.deepPurple : Colors.black,
                    ),
                  ),
                  subtitle: Text(msg['text'] ?? ''),
                  trailing: Text(msg['timestamp'] != null
                      ? DateTime.parse(msg['timestamp'])
                          .toLocal()
                          .toString()
                          .split(' ')[1]
                          .split('.')[0]
                      : ''),
                );
              },
            ),
          ),
          Padding(
            padding: const EdgeInsets.all(8.0),
            child: Row(
              children: [
                Expanded(
                  child: TextField(
                    controller: _controller,
                    decoration: const InputDecoration(
                      hintText: 'Enter message...',
                      border: OutlineInputBorder(),
                    ),
                    onSubmitted: (_) => _sendMessage(),
                  ),
                ),
                const SizedBox(width: 8),
                IconButton(
                  icon: const Icon(Icons.send),
                  onPressed: _sendMessage,
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

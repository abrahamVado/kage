import 'dart:async';
import 'dart:convert';

import 'package:web_socket_channel/web_socket_channel.dart';
import 'package:web_socket_channel/status.dart' as status;

class WebSocketService {
  WebSocketService({this.baseUrl = 'ws://localhost:8080', WebSocketChannel Function(Uri uri)? channelBuilder})
      : _channelBuilder = channelBuilder ?? WebSocketChannel.connect;

  final String baseUrl;
  final WebSocketChannel Function(Uri uri) _channelBuilder;
  WebSocketChannel? _channel;
  final _outgoingQueue = StreamController<String>.broadcast();
  StreamSubscription<String>? _outgoingSubscription;

  Stream<dynamic> get stream => _channel?.stream ?? const Stream.empty();

  Stream<String> get outgoing => _outgoingQueue.stream;

  Future<void> connect(String token) async {
    //1.- Close any existing connection to avoid duplicate sockets.
    await disconnect();

    //2.- Establish a new WebSocket connection against the Go backend.
    final uri = Uri.parse('$baseUrl/ws?token=$token');
    _channel = _channelBuilder(uri);

    //3.- Pipe queued messages into the active channel.
    await _outgoingSubscription?.cancel();
    _outgoingSubscription = _outgoingQueue.stream.listen((message) {
      _channel?.sink.add(message);
    });
  }

  void sendJson(Map<String, dynamic> payload) {
    //1.- Encode the payload and push it into the outgoing stream.
    final encoded = jsonEncode(payload);
    if (_channel != null) {
      _channel!.sink.add(encoded);
    } else {
      _outgoingQueue.add(encoded);
    }
  }

  Future<void> disconnect() async {
    //1.- Gracefully close the channel when the feature flow is disposed.
    await _channel?.sink.close(status.normalClosure);
    _channel = null;
  }

  void dispose() {
    //1.- Cleanup resources to avoid memory leaks.
    _outgoingQueue.close();
    _outgoingSubscription?.cancel();
    _channel?.sink.close(status.goingAway);
  }
}

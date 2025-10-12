import 'dart:convert';

import 'package:http/http.dart' as http;

import '../../features/auth/auth_model.dart';

class AuthSession {
  const AuthSession({required this.token, required this.userId});

  final String token;
  final String userId;
}

class AuthService {
  AuthService({http.Client? client, this.baseUrl = 'http://localhost:8080'})
      : _client = client ?? http.Client();

  final http.Client _client;
  final String baseUrl;

  Future<AuthSession> login(AuthCredentials credentials) async {
    //1.- Prepare the HTTP request payload for the Go backend.
    final payload = jsonEncode({
      'email': credentials.email,
      'password': credentials.password,
    });

    //2.- Execute the POST request against the login endpoint.
    final response = await _client.post(
      Uri.parse('$baseUrl/api/v1/auth/login'),
      headers: const {'Content-Type': 'application/json'},
      body: payload,
    );

    //3.- Validate the response and map it to an AuthSession.
    if (response.statusCode >= 200 && response.statusCode < 300) {
      final body = jsonDecode(response.body) as Map<String, dynamic>;
      final token = body['token'] as String?;
      final userId = body['user_id'] as String?;
      if (token != null && userId != null) {
        return AuthSession(token: token, userId: userId);
      }
    }

    throw const AuthException('Unable to authenticate with the provided credentials.');
  }
}

class AuthException implements Exception {
  const AuthException(this.message);

  final String message;

  @override
  String toString() => 'AuthException: $message';
}

import 'package:flutter/foundation.dart';

import '../../core/services/auth_service.dart';
import '../../core/services/location_service.dart';
import '../../core/services/secure_storage_service.dart';
import '../../core/services/websocket_service.dart';
import 'auth_model.dart';

class AuthViewModel extends ChangeNotifier {
  AuthViewModel({
    required this.authService,
    required this.secureStorageService,
    required this.webSocketService,
    required this.locationService,
  });

  final AuthService authService;
  final SecureStorageService secureStorageService;
  final WebSocketService webSocketService;
  final LocationService locationService;

  String email = '';
  String password = '';
  AuthViewState _state = const AuthViewState();

  AuthViewState get state => _state;

  void updateEmail(String value) {
    email = value;
    notifyListeners();
  }

  void updatePassword(String value) {
    password = value;
    notifyListeners();
  }

  Future<AuthSession?> login() async {
    //1.- Guard against missing credentials to avoid empty requests.
    if (email.isEmpty || password.isEmpty) {
      _state = _state.copyWith(errorMessage: 'Email and password are required.');
      notifyListeners();
      return null;
    }

    //2.- Trigger loading state and notify listeners.
    _state = _state.copyWith(isLoading: true, errorMessage: null);
    notifyListeners();

    try {
      //3.- Delegate authentication to the service layer.
      final credentials = AuthCredentials(email: email, password: password);
      final session = await authService.login(credentials);

      //4.- Persist credentials, open WebSocket, and start location publishing.
      await secureStorageService.saveToken(session.token);
      await webSocketService.connect(session.token);
      await locationService.startPublishing(webSocketService: webSocketService);

      //5.- Complete the flow by clearing the loading state.
      _state = _state.copyWith(isLoading: false);
      notifyListeners();
      return session;
    } catch (error) {
      //6.- Surface any issues back to the UI for display.
      _state = _state.copyWith(
        isLoading: false,
        errorMessage: error is AuthException ? error.message : 'Login failed. Please retry.',
      );
      notifyListeners();
      return null;
    }
  }

  Future<void> logout() async {
    //1.- Clear persisted data and disconnect socket/location streams.
    await locationService.stopPublishing();
    await webSocketService.disconnect();
    await secureStorageService.clearToken();
  }
}

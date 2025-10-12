import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class SecureStorageService {
  SecureStorageService({FlutterSecureStorage? storage})
      : _storage = storage ?? const FlutterSecureStorage();

  final FlutterSecureStorage _storage;

  Future<void> saveToken(String token) {
    //1.- Persist the token using platform secure storage mechanisms.
    return _storage.write(key: 'auth_token', value: token);
  }

  Future<String?> readToken() {
    //1.- Retrieve the persisted token for authenticated calls.
    return _storage.read(key: 'auth_token');
  }

  Future<void> clearToken() {
    //1.- Remove the token when the session is invalidated.
    return _storage.delete(key: 'auth_token');
  }
}

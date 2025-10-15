import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';

import 'package:flutter_app/features/auth/auth_view_model.dart';
import 'package:flutter_app/core/services/auth_service.dart';
import 'package:flutter_app/core/services/location_service.dart';
import 'package:flutter_app/core/services/secure_storage_service.dart';
import 'package:flutter_app/core/services/websocket_service.dart';

class MockAuthService extends Mock implements AuthService {}

class MockSecureStorageService extends Mock implements SecureStorageService {}

class MockWebSocketService extends Mock implements WebSocketService {}

class MockLocationService extends Mock implements LocationService {}

void main() {
  late MockAuthService authService;
  late MockSecureStorageService secureStorageService;
  late MockWebSocketService webSocketService;
  late MockLocationService locationService;
  late AuthViewModel viewModel;

  setUp(() {
    authService = MockAuthService();
    secureStorageService = MockSecureStorageService();
    webSocketService = MockWebSocketService();
    locationService = MockLocationService();

    viewModel = AuthViewModel(
      authService: authService,
      secureStorageService: secureStorageService,
      webSocketService: webSocketService,
      locationService: locationService,
    );
  });

  test('startDemoSession provisions a local session and side effects', () async {
    when(() => secureStorageService.saveToken(any())).thenAnswer((_) async {});
    when(() => webSocketService.connect(any())).thenAnswer((_) async {});
    when(() => locationService.startPublishing(webSocketService: webSocketService))
        .thenAnswer((_) async {});

    final session = await viewModel.startDemoSession();

    expect(session, isNotNull);
    expect(session!.token, equals('demo-token'));
    expect(session.userId, equals('demo-user'));
    expect(viewModel.state.isLoading, isFalse);
    expect(viewModel.state.errorMessage, isNull);
    verify(() => secureStorageService.saveToken('demo-token')).called(1);
    verify(() => webSocketService.connect('demo-token')).called(1);
    verify(() => locationService.startPublishing(webSocketService: webSocketService))
        .called(1);
  });

  test('startDemoSession surfaces errors when setup fails', () async {
    when(() => secureStorageService.saveToken(any())).thenAnswer((_) async {});
    when(() => webSocketService.connect(any())).thenAnswer((_) async {});
    when(() => locationService.startPublishing(webSocketService: webSocketService))
        .thenThrow(Exception('failure'));

    final session = await viewModel.startDemoSession();

    expect(session, isNull);
    expect(viewModel.state.isLoading, isFalse);
    expect(
      viewModel.state.errorMessage,
      equals('Unable to start the demo right now. Please retry.'),
    );
  });
}

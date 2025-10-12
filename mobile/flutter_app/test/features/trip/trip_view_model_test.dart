import 'dart:async';

import 'package:fake_async/fake_async.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:flutter_app/core/services/websocket_service.dart';
import 'package:flutter_app/features/trip/trip_model.dart';
import 'package:flutter_app/features/trip/trip_view_model.dart';

class _TestWebSocketService extends WebSocketService {
  _TestWebSocketService() : super(baseUrl: 'ws://test');

  final StreamController<dynamic> _controller = StreamController<dynamic>();
  final List<Map<String, dynamic>> sent = [];

  @override
  Stream<dynamic> get stream => _controller.stream;

  @override
  void sendJson(Map<String, dynamic> payload) {
    sent.add(payload);
  }
}

void main() {
  test('timer supports start, pause, and resume', () {
    fakeAsync((async) {
      final service = _TestWebSocketService();
      final viewModel = TripViewModel(webSocketService: service);

      viewModel.startTrip();
      async.elapse(const Duration(seconds: 3));
      expect(viewModel.state.elapsed.inSeconds, 3);
      expect(viewModel.state.status, TripStatus.enRoute);

      viewModel.pauseTimer();
      final pausedSeconds = viewModel.state.elapsed.inSeconds;
      async.elapse(const Duration(seconds: 2));
      expect(viewModel.state.elapsed.inSeconds, pausedSeconds);

      viewModel.resumeTimer();
      async.elapse(const Duration(seconds: 4));
      expect(viewModel.state.elapsed.inSeconds, pausedSeconds + 4);

      viewModel.completeTrip();
      expect(viewModel.state.status, TripStatus.completed);
      expect(service.sent.map((event) => event['type']).toList(),
          containsAll(['trip_start', 'trip_pause', 'trip_resume', 'trip_complete']));
    });
  });
}

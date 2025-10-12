import 'dart:async';

import 'package:flutter/foundation.dart';

import '../../core/services/websocket_service.dart';
import 'trip_model.dart';

class TripViewModel extends ChangeNotifier {
  TripViewModel({required this.webSocketService});

  final WebSocketService webSocketService;
  TripState _state = const TripState();
  Timer? _timer;

  TripState get state => _state;

  void startTrip() {
    if (_state.status == TripStatus.enRoute && _state.isTimerRunning) {
      return;
    }

    //1.- Reset timer state and inform backend that trip has started.
    _state = const TripState(status: TripStatus.enRoute, isTimerRunning: true);
    notifyListeners();
    _startTicker();
    webSocketService.sendJson({'type': 'trip_start'});
  }

  void pauseTimer() {
    if (!_state.isTimerRunning) {
      return;
    }

    //1.- Pause the periodic timer and emit pause event.
    _timer?.cancel();
    _timer = null;
    _state = _state.copyWith(isTimerRunning: false);
    notifyListeners();
    webSocketService.sendJson({'type': 'trip_pause'});
  }

  void resumeTimer() {
    if (_state.isTimerRunning || _state.status == TripStatus.completed) {
      return;
    }

    //1.- Restart the ticker and notify backend about the resume action.
    _state = _state.copyWith(isTimerRunning: true);
    notifyListeners();
    _startTicker();
    webSocketService.sendJson({'type': 'trip_resume'});
  }

  void markArrived() {
    if (_state.status == TripStatus.arrived || _state.status == TripStatus.completed) {
      return;
    }

    //1.- Set status to arrived while keeping timer state.
    _state = _state.copyWith(status: TripStatus.arrived);
    notifyListeners();
    webSocketService.sendJson({'type': 'trip_arrived'});
  }

  void completeTrip() {
    if (_state.status == TripStatus.completed) {
      return;
    }

    //1.- Stop the timer, update status, and inform backend about completion.
    _timer?.cancel();
    _timer = null;
    _state = _state.copyWith(status: TripStatus.completed, isTimerRunning: false);
    notifyListeners();
    webSocketService.sendJson({'type': 'trip_complete', 'payload': {'duration_seconds': _state.elapsed.inSeconds}});
  }

  void _startTicker() {
    //1.- Schedule a periodic task to increment elapsed time.
    _timer?.cancel();
    _timer = Timer.periodic(const Duration(seconds: 1), (_) {
      _state = _state.copyWith(elapsed: _state.elapsed + const Duration(seconds: 1));
      notifyListeners();
    });
  }

  @override
  void dispose() {
    _timer?.cancel();
    super.dispose();
  }
}

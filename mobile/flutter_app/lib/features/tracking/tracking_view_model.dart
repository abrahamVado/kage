import 'dart:async';
import 'dart:convert';

import 'package:flutter/foundation.dart';

import '../../core/services/websocket_service.dart';
import 'tracking_model.dart';

class TrackingViewModel extends ChangeNotifier {
  TrackingViewModel({required this.webSocketService});

  final WebSocketService webSocketService;
  TrackingState _state = const TrackingState();
  StreamSubscription<dynamic>? _subscription;

  TrackingState get state => _state;

  void startListening() {
    if (_state.isListening) {
      return;
    }

    _state = _state.copyWith(isListening: true);
    notifyListeners();
    _subscription = webSocketService.stream.listen(_handleEvent);
  }

  void _handleEvent(dynamic event) {
    //1.- Decode the WebSocket payload and ensure it is JSON.
    final data = event is String ? jsonDecode(event) as Map<String, dynamic> : event as Map<String, dynamic>;
    if (data['type'] != 'rider_update') {
      return;
    }

    //2.- Map the rider payloads into strongly typed models.
    final riders = (data['payload'] as List<dynamic>? ?? [])
        .map((entry) => NearbyRider(
              id: entry['id'] as String,
              name: entry['name'] as String,
              distanceMeters: (entry['distance_meters'] as num).toDouble(),
            ))
        .toList()
      ..sort((a, b) => a.distanceMeters.compareTo(b.distanceMeters));

    //3.- Update state with the sorted list for UI rendering.
    _state = _state.copyWith(riders: riders);
    notifyListeners();
  }

  Future<void> stopListening() async {
    await _subscription?.cancel();
    _subscription = null;
    _state = _state.copyWith(isListening: false);
    notifyListeners();
  }

  @override
  void dispose() {
    _subscription?.cancel();
    super.dispose();
  }
}

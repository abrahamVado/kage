import 'dart:async';

import 'package:flutter_test/flutter_test.dart';

import 'package:flutter_app/core/services/websocket_service.dart';
import 'package:flutter_app/features/bidding/bidding_view_model.dart';

class _TestWebSocketService extends WebSocketService {
  _TestWebSocketService(this.controller) : super(baseUrl: 'ws://test');

  final StreamController<dynamic> controller;
  final List<Map<String, dynamic>> sent = [];

  @override
  Stream<dynamic> get stream => controller.stream;

  @override
  void sendJson(Map<String, dynamic> payload) {
    sent.add(payload);
  }
}

void main() {
  test('sorts incoming bids by proximity', () async {
    final controller = StreamController<dynamic>.broadcast();
    final service = _TestWebSocketService(controller);
    addTearDown(() async {
      await controller.close();
      service.dispose();
    });
    final viewModel = BiddingViewModel(webSocketService: service);

    viewModel.initialize();
    controller.add({
      'type': 'bid_offer',
      'payload': [
        {'id': 'b1', 'rider_id': 'r1', 'price': 15, 'distance_meters': 550.0},
        {'id': 'b2', 'rider_id': 'r2', 'price': 16, 'distance_meters': 120.0},
        {'id': 'b3', 'rider_id': 'r3', 'price': 13, 'distance_meters': 320.0},
      ],
    });

    await Future<void>.delayed(Duration.zero);

    expect(viewModel.state.bids.map((bid) => bid.id).toList(), ['b2', 'b3', 'b1']);
  });

  test('submitting bid forwards payload to backend', () async {
    final controller = StreamController<dynamic>.broadcast();
    final service = _TestWebSocketService(controller);
    addTearDown(() async {
      await controller.close();
      service.dispose();
    });
    final viewModel = BiddingViewModel(webSocketService: service);

    await viewModel.submitBid(riderId: 'r1', price: 10.5);

    expect(service.sent.length, 1);
    expect(service.sent.first['type'], 'submit_bid');
    expect(service.sent.first['payload'], {'rider_id': 'r1', 'price': 10.5});
  });
}

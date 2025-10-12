import 'dart:async';
import 'dart:convert';
import 'package:flutter/foundation.dart';

import '../../core/services/websocket_service.dart';
import 'bidding_model.dart';

class BiddingViewModel extends ChangeNotifier {
  BiddingViewModel({required this.webSocketService});

  final WebSocketService webSocketService;
  BiddingState _state = const BiddingState();
  StreamSubscription<dynamic>? _subscription;

  BiddingState get state => _state;

  void initialize() {
    _subscription ??= webSocketService.stream.listen(_handleEvent);
  }

  void _handleEvent(dynamic event) {
    //1.- Decode WebSocket payloads to handle bid updates.
    final data = event is String ? jsonDecode(event) as Map<String, dynamic> : event as Map<String, dynamic>;
    switch (data['type']) {
      case 'bid_offer':
        final bids = (data['payload'] as List<dynamic>? ?? [])
            .map((entry) => Bid(
                  id: entry['id'] as String,
                  riderId: entry['rider_id'] as String,
                  price: (entry['price'] as num).toDouble(),
                  distanceMeters: (entry['distance_meters'] as num).toDouble(),
                ))
            .toList()
          ..sort((a, b) => a.distanceMeters.compareTo(b.distanceMeters));
        _state = _state.copyWith(bids: bids);
        notifyListeners();
        break;
      case 'bid_selected':
        _state = _state.copyWith(selectedBidId: data['payload'] as String?);
        notifyListeners();
        break;
      default:
    }
  }

  Future<void> submitBid({required String riderId, required double price}) async {
    //1.- Toggle loading state while building the outgoing payload.
    _state = _state.copyWith(isSubmitting: true);
    notifyListeners();

    try {
      //2.- Send the bid payload through the WebSocket service.
      webSocketService.sendJson({
        'type': 'submit_bid',
        'payload': {
          'rider_id': riderId,
          'price': price,
        },
      });
    } finally {
      //3.- Regardless of network result, return UI to idle state.
      _state = _state.copyWith(isSubmitting: false);
      notifyListeners();
    }
  }

  void selectBid(String bidId) {
    //1.- Optimistically select the bid and inform the backend.
    _state = _state.copyWith(selectedBidId: bidId);
    notifyListeners();
    webSocketService.sendJson({
      'type': 'select_bid',
      'payload': {'bid_id': bidId},
    });
  }

  @override
  void dispose() {
    _subscription?.cancel();
    super.dispose();
  }
}

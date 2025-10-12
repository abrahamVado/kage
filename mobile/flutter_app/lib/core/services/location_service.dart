import 'dart:async';

import 'package:geolocator/geolocator.dart';

import 'websocket_service.dart';

class LocationService {
  LocationService({
    GeolocatorPlatform? geolocator,
    this.positionStreamOverride,
    this.permissionRequestOverride,
  }) : _geolocator = geolocator ?? GeolocatorPlatform.instance;

  final GeolocatorPlatform _geolocator;
  final Stream<Position>? positionStreamOverride;
  final Future<LocationPermission> Function()? permissionRequestOverride;
  StreamSubscription<Position>? _subscription;

  Future<void> startPublishing({
    required WebSocketService webSocketService,
    Duration interval = const Duration(seconds: 5),
  }) async {
    //1.- Ensure location permissions are granted before streaming positions.
    final permission = await _requestPermission();
    if (permission == LocationPermission.deniedForever ||
        permission == LocationPermission.denied) {
      throw const LocationException('Location permissions are not granted.');
    }

    //2.- Create the stream either from override (tests) or Geolocator API.
    final stream = positionStreamOverride ??
        _geolocator.getPositionStream(
          locationSettings: LocationSettings(
            distanceFilter: 10,
            timeLimit: interval,
            accuracy: LocationAccuracy.high,
          ),
        );

    //3.- Forward each position to the Go backend via WebSocket.
    await _subscription?.cancel();
    _subscription = stream.listen((position) {
      webSocketService.sendJson({
        'type': 'location_update',
        'lat': position.latitude,
        'lng': position.longitude,
        'timestamp': DateTime.now().toIso8601String(),
      });
    });
  }

  Future<LocationPermission> _requestPermission() async {
    if (permissionRequestOverride != null) {
      return permissionRequestOverride!.call();
    }

    //1.- Check existing permission state and prompt if needed.
    var permission = await _geolocator.checkPermission();
    if (permission == LocationPermission.denied) {
      permission = await _geolocator.requestPermission();
    }
    return permission;
  }

  Future<void> stopPublishing() async {
    //1.- Cancel the active stream subscription when the trip stops.
    await _subscription?.cancel();
    _subscription = null;
  }
}

class LocationException implements Exception {
  const LocationException(this.message);

  final String message;

  @override
  String toString() => 'LocationException: $message';
}

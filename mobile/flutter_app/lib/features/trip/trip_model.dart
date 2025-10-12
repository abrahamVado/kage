enum TripStatus { idle, enRoute, arrived, completed }

class TripState {
  const TripState({
    this.status = TripStatus.idle,
    this.elapsed = Duration.zero,
    this.isTimerRunning = false,
  });

  final TripStatus status;
  final Duration elapsed;
  final bool isTimerRunning;

  TripState copyWith({TripStatus? status, Duration? elapsed, bool? isTimerRunning}) {
    return TripState(
      status: status ?? this.status,
      elapsed: elapsed ?? this.elapsed,
      isTimerRunning: isTimerRunning ?? this.isTimerRunning,
    );
  }
}

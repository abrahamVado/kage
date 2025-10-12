import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../shared_widgets/primary_button.dart';
import '../../shared_widgets/responsive_scaffold.dart';
import 'trip_view_model.dart';

class TripView extends StatelessWidget {
  const TripView({super.key});

  String _formatDuration(Duration duration) {
    final minutes = duration.inMinutes.remainder(60).toString().padLeft(2, '0');
    final seconds = duration.inSeconds.remainder(60).toString().padLeft(2, '0');
    final hours = duration.inHours.toString().padLeft(2, '0');
    return '$hours:$minutes:$seconds';
  }

  @override
  Widget build(BuildContext context) {
    return Consumer<TripViewModel>(
      builder: (context, viewModel, _) {
        final state = viewModel.state;
        return ResponsiveScaffold(
          title: 'Trip Control',
          body: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Text(
                'Elapsed time: ' + _formatDuration(state.elapsed),
                style: Theme.of(context).textTheme.headlineSmall,
              ),
              const SizedBox(height: 12),
              Text('Status: ' + state.status.name),
              const SizedBox(height: 24),
              Wrap(
                spacing: 12,
                runSpacing: 12,
                children: [
                  PrimaryButton(
                    label: 'Start Trip',
                    onPressed: viewModel.startTrip,
                    isLoading: false,
                  ),
                  PrimaryButton(
                    label: state.isTimerRunning ? 'Pause Timer' : 'Resume Timer',
                    onPressed: state.isTimerRunning ? viewModel.pauseTimer : viewModel.resumeTimer,
                    isLoading: false,
                  ),
                  PrimaryButton(
                    label: 'Arrived',
                    onPressed: viewModel.markArrived,
                    isLoading: false,
                  ),
                  PrimaryButton(
                    label: 'Complete Trip',
                    onPressed: viewModel.completeTrip,
                    isLoading: false,
                  ),
                ],
              ),
            ],
          ),
        );
      },
    );
  }
}

import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../shared_widgets/responsive_scaffold.dart';
import 'tracking_view_model.dart';

class TrackingView extends StatefulWidget {
  const TrackingView({super.key});

  @override
  State<TrackingView> createState() => _TrackingViewState();
}

class _TrackingViewState extends State<TrackingView> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (!mounted) return;
      Provider.of<TrackingViewModel>(context, listen: false).startListening();
    });
  }

  @override
  void dispose() {
    Provider.of<TrackingViewModel>(context, listen: false).stopListening();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Consumer<TrackingViewModel>(
      builder: (context, viewModel, _) {
        return ResponsiveScaffold(
          title: 'Nearby Riders',
          body: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                'Active riders near you',
                style: Theme.of(context).textTheme.titleLarge,
              ),
              const SizedBox(height: 12),
              Expanded(
                child: ListView.separated(
                  itemCount: viewModel.state.riders.length,
                  separatorBuilder: (_, __) => const Divider(),
                  itemBuilder: (context, index) {
                    final rider = viewModel.state.riders[index];
                    final initials = rider.name.isNotEmpty
                        ? rider.name[0].toUpperCase()
                        : '?';
                    return ListTile(
                      title: Text(rider.name),
                      subtitle: Text('${rider.distanceMeters.toStringAsFixed(0)} meters away'),
                      leading: CircleAvatar(child: Text(initials)),
                    );
                  },
                ),
              ),
            ],
          ),
        );
      },
    );
  }
}

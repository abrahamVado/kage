import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../shared_widgets/primary_button.dart';
import '../../shared_widgets/responsive_scaffold.dart';
import 'auth_view_model.dart';

class AuthView extends StatelessWidget {
  const AuthView({super.key});

  @override
  Widget build(BuildContext context) {
    return Consumer<AuthViewModel>(
      builder: (context, viewModel, _) {
        return ResponsiveScaffold(
          title: 'Driver Login',
          body: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              TextField(
                key: const Key('emailField'),
                decoration: const InputDecoration(labelText: 'Email'),
                onChanged: viewModel.updateEmail,
                keyboardType: TextInputType.emailAddress,
              ),
              const SizedBox(height: 16),
              TextField(
                key: const Key('passwordField'),
                decoration: const InputDecoration(labelText: 'Password'),
                obscureText: true,
                onChanged: viewModel.updatePassword,
              ),
              const SizedBox(height: 24),
              if (viewModel.state.errorMessage != null)
                Text(
                  viewModel.state.errorMessage!,
                  style: TextStyle(color: Theme.of(context).colorScheme.error),
                ),
              const SizedBox(height: 8),
              PrimaryButton(
                label: 'Sign In',
                isLoading: viewModel.state.isLoading,
                onPressed: () async {
                  final session = await viewModel.login();
                  if (session != null && context.mounted) {
                    Navigator.of(context).pushReplacementNamed('/dashboard');
                  }
                },
              ),
              const SizedBox(height: 12),
              OutlinedButton(
                onPressed: viewModel.state.isLoading
                    ? null
                    : () async {
                        //1.- Attempt to start a local demo session for quick access.
                        final session = await viewModel.startDemoSession();
                        //2.- Mirror the authenticated navigation when the demo succeeds.
                        if (session != null && context.mounted) {
                          Navigator.of(context)
                              .pushReplacementNamed('/dashboard');
                        }
                      },
                child: const Text('Explore Demo'),
              ),
            ],
          ),
        );
      },
    );
  }
}

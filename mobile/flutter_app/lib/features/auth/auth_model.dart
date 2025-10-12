class AuthCredentials {
  const AuthCredentials({required this.email, required this.password});

  final String email;
  final String password;
}

class AuthViewState {
  const AuthViewState({
    this.isLoading = false,
    this.errorMessage,
  });

  final bool isLoading;
  final String? errorMessage;

  AuthViewState copyWith({bool? isLoading, String? errorMessage}) {
    return AuthViewState(
      isLoading: isLoading ?? this.isLoading,
      errorMessage: errorMessage,
    );
  }
}

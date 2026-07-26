import 'package:flutter/material.dart';
import '../services/storage_service.dart';

class AppProvider with ChangeNotifier {
  bool _onboardingCompleted = false;
  String _selectedCity = 'La Paz';
  bool _isOnline = true;
  bool _isPremium = false;

  bool get onboardingCompleted => _onboardingCompleted;
  String get selectedCity => _selectedCity;
  bool get isOnline => _isOnline;
  bool get isPremium => _isPremium;

  AppProvider() {
    _init();
  }

  Future<void> _init() async {
    _onboardingCompleted = await StorageService.isOnboardingCompleted();
    _selectedCity = await StorageService.getSelectedCity();
    _isPremium = await StorageService.isPremiumUser();
    notifyListeners();
  }

  Future<void> completeOnboarding() async {
    _onboardingCompleted = true;
    await StorageService.setOnboardingCompleted(true);
    notifyListeners();
  }

  Future<void> setCity(String city) async {
    _selectedCity = city;
    await StorageService.setSelectedCity(city);
    notifyListeners();
  }

  Future<void> setPremiumStatus(bool status) async {
    _isPremium = status;
    await StorageService.setPremiumUser(status);
    notifyListeners();
  }

  void setOnlineStatus(bool online) {
    _isOnline = online;
    notifyListeners();
  }
}

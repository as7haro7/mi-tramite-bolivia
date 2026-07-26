import 'dart:convert';
import 'package:shared_preferences/shared_preferences.dart';
import '../models/checklist_item.dart';
import '../models/tramite.dart';

class StorageService {
  static const String _kOnboardingCompleted = 'onboarding_completed';
  static const String _kSelectedCity = 'selected_city';
  static const String _kFavorites = 'favorites_slugs';
  static const String _kChecklists = 'user_checklists';
  static const String _kRecentSearches = 'recent_searches';
  static const String _kCachedTramites = 'cached_tramites';
  static const String _kIsPremium = 'is_premium_user';

  // Onboarding
  static Future<bool> isOnboardingCompleted() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getBool(_kOnboardingCompleted) ?? false;
  }

  static Future<void> setOnboardingCompleted(bool completed) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setBool(_kOnboardingCompleted, completed);
  }

  // Selected City
  static Future<String> getSelectedCity() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString(_kSelectedCity) ?? 'La Paz';
  }

  static Future<void> setSelectedCity(String city) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_kSelectedCity, city);
  }

  // Premium Subscription (Demo)
  static Future<bool> isPremiumUser() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getBool(_kIsPremium) ?? false;
  }

  static Future<void> setPremiumUser(bool isPremium) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setBool(_kIsPremium, isPremium);
  }

  // Favorites
  static Future<List<String>> getFavorites() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getStringList(_kFavorites) ?? [];
  }

  static Future<void> toggleFavorite(String slug) async {
    final prefs = await SharedPreferences.getInstance();
    final list = prefs.getStringList(_kFavorites) ?? [];
    if (list.contains(slug)) {
      list.remove(slug);
    } else {
      list.add(slug);
    }
    await prefs.setStringList(_kFavorites, list);
  }

  // Checklists
  static Future<List<TramiteChecklist>> getChecklists() async {
    final prefs = await SharedPreferences.getInstance();
    final raw = prefs.getString(_kChecklists);
    if (raw == null) return [];
    try {
      final List decoded = jsonDecode(raw);
      return decoded.map((item) => TramiteChecklist.fromJson(item)).toList();
    } catch (_) {
      return [];
    }
  }

  static Future<void> saveChecklists(List<TramiteChecklist> checklists) async {
    final prefs = await SharedPreferences.getInstance();
    final raw = jsonEncode(checklists.map((c) => c.toJson()).toList());
    await prefs.setString(_kChecklists, raw);
  }

  // Recent Searches
  static Future<List<String>> getRecentSearches() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getStringList(_kRecentSearches) ?? [];
  }

  static Future<void> addRecentSearch(String term) async {
    if (term.trim().isEmpty) return;
    final prefs = await SharedPreferences.getInstance();
    final list = prefs.getStringList(_kRecentSearches) ?? [];
    list.remove(term);
    list.insert(0, term);
    if (list.length > 8) list.removeLast();
    await prefs.setStringList(_kRecentSearches, list);
  }

  static Future<void> clearRecentSearches() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove(_kRecentSearches);
  }

  // Offline Cache for Tramites
  static Future<List<Tramite>> getCachedTramites() async {
    final prefs = await SharedPreferences.getInstance();
    final raw = prefs.getString(_kCachedTramites);
    if (raw == null) return [];
    try {
      final List decoded = jsonDecode(raw);
      return decoded.map((item) => Tramite.fromJson(item)).toList();
    } catch (_) {
      return [];
    }
  }

  static Future<void> cacheTramites(List<Tramite> tramites) async {
    final prefs = await SharedPreferences.getInstance();
    final raw = jsonEncode(tramites.map((t) => t.toJson()).toList());
    await prefs.setString(_kCachedTramites, raw);
  }

  // Clear all local data
  static Future<void> clearAll() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.clear();
  }
}

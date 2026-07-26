import 'package:flutter/material.dart';
import '../services/storage_service.dart';

class FavoritesProvider with ChangeNotifier {
  Set<String> _favoriteSlugs = {};

  Set<String> get favoriteSlugs => _favoriteSlugs;

  FavoritesProvider() {
    _loadFavorites();
  }

  Future<void> _loadFavorites() async {
    final list = await StorageService.getFavorites();
    _favoriteSlugs = list.toSet();
    notifyListeners();
  }

  bool isFavorite(String slug) => _favoriteSlugs.contains(slug);

  Future<void> toggleFavorite(String slug) async {
    if (_favoriteSlugs.contains(slug)) {
      _favoriteSlugs.remove(slug);
    } else {
      _favoriteSlugs.add(slug);
    }
    await StorageService.toggleFavorite(slug);
    notifyListeners();
  }
}

import 'dart:io';
import 'package:flutter/foundation.dart' show kIsWeb;

/// Configuración centralizada de IDs de anuncios para Google AdMob
/// 
/// IMPORTANTE: Antes de publicar en producción, reemplaza estos IDs de prueba
/// con tus propios IDs de anuncios de AdMob
class AdConfig {
  // IDs de prueba de Google AdMob (para desarrollo)
  // Estos IDs no generan ingresos pero permiten probar la integración
  static const String _testBannerAdUnitIdAndroid = 'ca-app-pub-3940256099942544/6300978111';
  static const String _testBannerAdUnitIdIOS = 'ca-app-pub-3940256099942544/2934735716';
  
  // TODO: Reemplazar con tus IDs reales de producción
  // static const String _prodBannerAdUnitIdAndroid = 'ca-app-pub-XXXXXXXXXXXXXXXX/YYYYYYYYYY';
  // static const String _prodBannerAdUnitIdIOS = 'ca-app-pub-XXXXXXXXXXXXXXXX/YYYYYYYYYY';

  /// Obtiene el ID de anuncio de banner según la plataforma
  static String get bannerAdUnitId {
    // Web no soporta AdMob
    if (kIsWeb) {
      return '';
    }
    
    // En modo debug, usar IDs de prueba
    // En producción, descomentar y usar IDs reales
    if (Platform.isAndroid) {
      return _testBannerAdUnitIdAndroid;
      // return _prodBannerAdUnitIdAndroid; // Usar en producción
    } else if (Platform.isIOS) {
      return _testBannerAdUnitIdIOS;
      // return _prodBannerAdUnitIdIOS; // Usar en producción
    }
    throw UnsupportedError('Plataforma no soportada');
  }

  /// Indica si se deben mostrar anuncios (útil para versión premium sin ads)
  static bool get shouldShowAds {
    // No mostrar anuncios en web (no soportado)
    if (kIsWeb) {
      return false;
    }
    
    // TODO: Integrar con lógica de suscripción premium
    // Por ahora siempre muestra ads en móviles (versión demo/free)
    return true;
  }
}

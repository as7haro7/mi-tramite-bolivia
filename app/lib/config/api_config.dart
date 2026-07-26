import 'dart:io';
import 'package:flutter/foundation.dart';

class ApiConfig {
  // Base URL resolution for Web, Android Emulator, and Desktop/iOS
  static String get baseUrl {
    if (kIsWeb) {
      return 'http://localhost:8080';
    }
    if (Platform.isAndroid) {
      // 10.0.2.2 maps to host machine localhost in Android Emulator
      return 'http://10.0.2.2:8080';
    }
    return 'http://localhost:8080';
  }

  static const Duration timeout = Duration(seconds: 12);

  // Endpoint routes
  static const String tramites = '/api/v1/tramites';
  static const String categorias = '/api/v1/categorias';
  static const String instituciones = '/api/v1/instituciones';
  static const String oficinas = '/api/v1/oficinas';
  static const String conversaciones = '/api/v1/chat/conversaciones';
  static const String reportes = '/api/v1/reportes';

  static String tramiteDetail(String slug) => '/api/v1/tramites/$slug';
  static String tramiteOficinas(String slug) => '/api/v1/tramites/$slug/oficinas';
  static String chatMensajes(String conversacionId) => '/api/v1/chat/conversaciones/$conversacionId/mensajes';
  static String mensajeFeedback(String mensajeId) => '/api/v1/mensajes/$mensajeId/feedback';
}

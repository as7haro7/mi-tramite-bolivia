import 'package:flutter/material.dart';
import '../models/chat_message.dart';
import '../services/api_service.dart';

class AssistantProvider with ChangeNotifier {
  List<ChatMessage> _messages = [];
  String? _conversacionId;
  bool _isLoading = false;

  List<ChatMessage> get messages => _messages;
  bool get isLoading => _isLoading;

  AssistantProvider() {
    _initConversation();
  }

  Future<void> _initConversation() async {
    // Welcome message
    _messages = [
      ChatMessage(
        id: 'welcome',
        content: '¡Hola! Soy tu asistente de **Mi Trámite Bolivia**. '
            '¿Qué trámite necesitas preparar o consultar hoy?',
        isUser: false,
        timestamp: DateTime.now(),
        fuentes: [
          FuenteRef(
            titulo: 'Fuentes Oficiales Verificadas',
            url: 'https://gob.bo',
          ),
        ],
      ),
    ];
    _conversacionId = await ApiService.crearConversacion();
    notifyListeners();
  }

  Future<void> sendMessage(String text, {String? city}) async {
    if (text.trim().isEmpty || _isLoading) return;

    final userMsg = ChatMessage(
      id: DateTime.now().millisecondsSinceEpoch.toString(),
      content: text.trim(),
      isUser: true,
      timestamp: DateTime.now(),
    );

    _messages.add(userMsg);
    _isLoading = true;
    notifyListeners();

    _conversacionId ??= await ApiService.crearConversacion();

    final response = await ApiService.enviarMensaje(
      conversacionId: _conversacionId ?? 'default',
      prompt: text.trim(),
      municipio: city,
    );

    _messages.add(response);
    _isLoading = false;
    notifyListeners();
  }

  Future<void> sendFeedback(String messageId, bool isPositive) async {
    final index = _messages.indexWhere((m) => m.id == messageId);
    if (index >= 0) {
      _messages[index] = _messages[index].copyWith(feedback: isPositive);
      notifyListeners();
      await ApiService.enviarFeedback(messageId, isPositive);
    }
  }

  void clearChat() {
    _messages.clear();
    _initConversation();
  }
}

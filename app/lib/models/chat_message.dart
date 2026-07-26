class ChatMessage {
  final String id;
  final String content;
  final bool isUser;
  final DateTime timestamp;
  final List<FuenteRef> fuentes;
  final bool? feedback; // true = 👍, false = 👎, null = no feedback yet

  ChatMessage({
    required this.id,
    required this.content,
    required this.isUser,
    required this.timestamp,
    this.fuentes = const [],
    this.feedback,
  });

  ChatMessage copyWith({
    String? id,
    String? content,
    bool? isUser,
    DateTime? timestamp,
    List<FuenteRef>? fuentes,
    bool? feedback,
  }) {
    return ChatMessage(
      id: id ?? this.id,
      content: content ?? this.content,
      isUser: isUser ?? this.isUser,
      timestamp: timestamp ?? this.timestamp,
      fuentes: fuentes ?? this.fuentes,
      feedback: feedback ?? this.feedback,
    );
  }

  factory ChatMessage.fromJson(Map<String, dynamic> json) {
    return ChatMessage(
      id: json['id'] ?? DateTime.now().millisecondsSinceEpoch.toString(),
      content: json['contenido'] ?? json['content'] ?? json['respuesta'] ?? '',
      isUser: json['rol'] == 'usuario' || json['is_user'] == true,
      timestamp: json['timestamp'] != null
          ? DateTime.parse(json['timestamp'])
          : DateTime.now(),
      fuentes: json['fuentes'] != null
          ? (json['fuentes'] as List).map((f) => FuenteRef.fromJson(f)).toList()
          : [],
      feedback: json['feedback'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'contenido': content,
      'is_user': isUser,
      'timestamp': timestamp.toIso8601String(),
      'fuentes': fuentes.map((f) => f.toJson()).toList(),
      'feedback': feedback,
    };
  }
}

class FuenteRef {
  final String titulo;
  final String url;
  final String? tramiteSlug;
  final String? verificadoEn;

  FuenteRef({
    required this.titulo,
    required this.url,
    this.tramiteSlug,
    this.verificadoEn,
  });

  factory FuenteRef.fromJson(Map<String, dynamic> json) {
    return FuenteRef(
      titulo: json['titulo'] ?? json['nombre'] ?? 'Fuente Oficial',
      url: json['url'] ?? '',
      tramiteSlug: json['tramite_slug'] ?? json['slug'],
      verificadoEn: json['verificado_en'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'titulo': titulo,
      'url': url,
      'tramite_slug': tramiteSlug,
      'verificado_en': verificadoEn,
    };
  }
}

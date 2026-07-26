class RequisitoCheckable {
  final int id;
  final String nombre;
  final String? descripcion;
  bool isCompleted;
  String personalNote;

  RequisitoCheckable({
    required this.id,
    required this.nombre,
    this.descripcion,
    this.isCompleted = false,
    this.personalNote = '',
  });

  factory RequisitoCheckable.fromJson(Map<String, dynamic> json) {
    return RequisitoCheckable(
      id: json['id'] is int ? json['id'] : int.parse((json['id'] ?? 0).toString()),
      nombre: json['nombre'] ?? '',
      descripcion: json['descripcion'],
      isCompleted: json['is_completed'] ?? false,
      personalNote: json['personal_note'] ?? '',
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'nombre': nombre,
      'descripcion': descripcion,
      'is_completed': isCompleted,
      'personal_note': personalNote,
    };
  }
}

class TramiteChecklist {
  final String tramiteSlug;
  final String tramiteTitulo;
  final String institucionNombre;
  final DateTime createdAt;
  final List<RequisitoCheckable> items;

  TramiteChecklist({
    required this.tramiteSlug,
    required this.tramiteTitulo,
    required this.institucionNombre,
    required this.createdAt,
    required this.items,
  });

  int get completedCount => items.where((i) => i.isCompleted).length;
  int get totalCount => items.length;
  double get progress => totalCount == 0 ? 0.0 : completedCount / totalCount;
  bool get isFullyCompleted => totalCount > 0 && completedCount == totalCount;

  factory TramiteChecklist.fromJson(Map<String, dynamic> json) {
    return TramiteChecklist(
      tramiteSlug: json['tramite_slug'] ?? '',
      tramiteTitulo: json['tramite_titulo'] ?? '',
      institucionNombre: json['institucion_nombre'] ?? '',
      createdAt: json['created_at'] != null ? DateTime.parse(json['created_at']) : DateTime.now(),
      items: json['items'] != null
          ? (json['items'] as List).map((i) => RequisitoCheckable.fromJson(i)).toList()
          : [],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'tramite_slug': tramiteSlug,
      'tramite_titulo': tramiteTitulo,
      'institucion_nombre': institucionNombre,
      'created_at': createdAt.toIso8601String(),
      'items': items.map((i) => i.toJson()).toList(),
    };
  }
}

class Categoria {
  final int id;
  final String codigo;
  final String nombre;
  final String? descripcion;
  final String? icono;

  Categoria({
    required this.id,
    required this.codigo,
    required this.nombre,
    this.descripcion,
    this.icono,
  });

  factory Categoria.fromJson(Map<String, dynamic> json) {
    return Categoria(
      id: json['id'] is int ? json['id'] : int.parse(json['id'].toString()),
      codigo: json['codigo'] ?? '',
      nombre: json['nombre'] ?? '',
      descripcion: json['descripcion'],
      icono: json['icono'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'codigo': codigo,
      'nombre': nombre,
      'descripcion': descripcion,
      'icono': icono,
    };
  }
}

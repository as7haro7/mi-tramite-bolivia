class Institucion {
  final int id;
  final String codigo;
  final String nombre;
  final String? sigla;
  final String? tipo;
  final String? portalOficial;
  final String? logoUrl;

  Institucion({
    required this.id,
    required this.codigo,
    required this.nombre,
    this.sigla,
    this.tipo,
    this.portalOficial,
    this.logoUrl,
  });

  factory Institucion.fromJson(Map<String, dynamic> json) {
    return Institucion(
      id: json['id'] is int ? json['id'] : int.parse(json['id'].toString()),
      codigo: json['codigo'] ?? '',
      nombre: json['nombre'] ?? '',
      sigla: json['sigla'],
      tipo: json['tipo'],
      portalOficial: json['portal_oficial'],
      logoUrl: json['logo_url'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'codigo': codigo,
      'nombre': nombre,
      'sigla': sigla,
      'tipo': tipo,
      'portal_oficial': portalOficial,
      'logo_url': logoUrl,
    };
  }
}

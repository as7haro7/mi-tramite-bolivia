import 'categoria.dart';
import 'institucion.dart';

class Tramite {
  final int id;
  final String slug;
  final String titulo;
  final String? resumen;
  final String? objetivo;
  final String? plazoEstimado;
  final String? verificadoEn;
  final String? proximaRevisionEn;
  final bool verificado;
  final Categoria? categoria;
  final Institucion? institucion;
  final List<Requisito> requisitos;
  final List<Paso> pasos;
  final List<Costo> costos;
  final List<Oficina> oficinas;
  final List<Modalidad> modalidades;
  final List<Fuente> fuentes;

  Tramite({
    required this.id,
    required this.slug,
    required this.titulo,
    this.resumen,
    this.objetivo,
    this.plazoEstimado,
    this.verificadoEn,
    this.proximaRevisionEn,
    this.verificado = true,
    this.categoria,
    this.institucion,
    this.requisitos = const [],
    this.pasos = const [],
    this.costos = const [],
    this.oficinas = const [],
    this.modalidades = const [],
    this.fuentes = const [],
  });

  factory Tramite.fromJson(Map<String, dynamic> json) {
    return Tramite(
      id: json['id'] is int ? json['id'] : int.parse((json['id'] ?? 0).toString()),
      slug: json['slug'] ?? '',
      titulo: json['titulo'] ?? json['nombre'] ?? '',
      resumen: json['resumen'] ?? json['descripcion'],
      objetivo: json['objetivo'],
      plazoEstimado: json['plazo_estimado'] ?? json['duracion_estimada'],
      verificadoEn: json['verificado_en'] ?? json['actualizado_en'],
      proximaRevisionEn: json['proxima_revision_en'],
      verificado: json['verificado'] ?? true,
      categoria: json['categoria'] != null ? Categoria.fromJson(json['categoria']) : null,
      institucion: json['institucion'] != null ? Institucion.fromJson(json['institucion']) : null,
      requisitos: json['requisitos'] != null
          ? (json['requisitos'] as List).map((i) => Requisito.fromJson(i)).toList()
          : [],
      pasos: json['pasos'] != null
          ? (json['pasos'] as List).map((i) => Paso.fromJson(i)).toList()
          : [],
      costos: json['costos'] != null
          ? (json['costos'] as List).map((i) => Costo.fromJson(i)).toList()
          : [],
      oficinas: json['oficinas'] != null
          ? (json['oficinas'] as List).map((i) => Oficina.fromJson(i)).toList()
          : [],
      modalidades: json['modalidades'] != null
          ? (json['modalidades'] as List).map((i) => Modalidad.fromJson(i)).toList()
          : [],
      fuentes: json['fuentes'] != null
          ? (json['fuentes'] as List).map((i) => Fuente.fromJson(i)).toList()
          : [],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'slug': slug,
      'titulo': titulo,
      'resumen': resumen,
      'objetivo': objetivo,
      'plazo_estimado': plazoEstimado,
      'verificado_en': verificadoEn,
      'proxima_revision_en': proximaRevisionEn,
      'verificado': verificado,
      'categoria': categoria?.toJson(),
      'institucion': institucion?.toJson(),
      'requisitos': requisitos.map((r) => r.toJson()).toList(),
      'pasos': pasos.map((p) => p.toJson()).toList(),
      'costos': costos.map((c) => c.toJson()).toList(),
      'oficinas': oficinas.map((o) => o.toJson()).toList(),
      'modalidades': modalidades.map((m) => m.toJson()).toList(),
      'fuentes': fuentes.map((f) => f.toJson()).toList(),
    };
  }
}

class Requisito {
  final int id;
  final String nombre;
  final String? descripcion;
  final bool esObligatorio;
  final String? tipoDocumento;
  final String? aplicaSi;

  Requisito({
    required this.id,
    required this.nombre,
    this.descripcion,
    this.esObligatorio = true,
    this.tipoDocumento,
    this.aplicaSi,
  });

  factory Requisito.fromJson(dynamic json) {
    if (json is String) {
      return Requisito(id: 0, nombre: json);
    }
    final map = json is Map ? json : <String, dynamic>{};
    return Requisito(
      id: map['id'] is int ? map['id'] : int.parse((map['id'] ?? 0).toString()),
      nombre: map['nombre'] ?? map['titulo'] ?? map.toString(),
      descripcion: map['descripcion'],
      esObligatorio: map['es_obligatorio'] ?? map['obligatorio'] ?? true,
      tipoDocumento: map['tipo_documento'],
      aplicaSi: map['aplica_si']?.toString(),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'nombre': nombre,
      'descripcion': descripcion,
      'es_obligatorio': esObligatorio,
      'tipo_documento': tipoDocumento,
      'aplica_si': aplicaSi,
    };
  }
}

class Paso {
  final int orden;
  final String titulo;
  final String? descripcion;
  final String? medioAtencion;

  Paso({
    required this.orden,
    required this.titulo,
    this.descripcion,
    this.medioAtencion,
  });

  factory Paso.fromJson(dynamic json) {
    if (json is String) {
      return Paso(orden: 1, titulo: json);
    }
    final map = json is Map ? json : <String, dynamic>{};
    return Paso(
      orden: map['orden'] is int ? map['orden'] : int.parse((map['orden'] ?? 1).toString()),
      titulo: map['titulo'] ?? map['nombre'] ?? '',
      descripcion: map['descripcion'],
      medioAtencion: map['medio_atencion'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'orden': orden,
      'titulo': titulo,
      'descripcion': descripcion,
      'medio_atencion': medioAtencion,
    };
  }
}

class Costo {
  final String concepto;
  final double monto;
  final String moneda;

  Costo({
    required this.concepto,
    required this.monto,
    this.moneda = 'BOB',
  });

  factory Costo.fromJson(Map<String, dynamic> json) {
    return Costo(
      concepto: json['concepto'] ?? json['nombre'] ?? 'Tasa general',
      monto: (json['monto'] is num) ? (json['monto'] as num).toDouble() : double.parse((json['monto'] ?? 0).toString()),
      moneda: json['moneda'] ?? 'BOB',
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'concepto': concepto,
      'monto': monto,
      'moneda': moneda,
    };
  }
}

class Oficina {
  final int id;
  final String nombre;
  final String? municipio;
  final String? direccion;
  final String? telefono;
  final String? horario;
  final double? latitud;
  final double? longitud;

  Oficina({
    required this.id,
    required this.nombre,
    this.municipio,
    this.direccion,
    this.telefono,
    this.horario,
    this.latitud,
    this.longitud,
  });

  factory Oficina.fromJson(Map<String, dynamic> json) {
    return Oficina(
      id: json['id'] is int ? json['id'] : int.parse((json['id'] ?? 0).toString()),
      nombre: json['nombre'] ?? json['oficina'] ?? '',
      municipio: json['municipio'] ?? json['ciudad'],
      direccion: json['direccion'],
      telefono: json['telefono'],
      horario: json['horario'] ?? json['horario_atencion'],
      latitud: json['latitud'] != null ? double.parse(json['latitud'].toString()) : null,
      longitud: json['longitud'] != null ? double.parse(json['longitud'].toString()) : null,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'nombre': nombre,
      'municipio': municipio,
      'direccion': direccion,
      'telefono': telefono,
      'horario': horario,
      'latitud': latitud,
      'longitud': longitud,
    };
  }
}

class Modalidad {
  final String codigo; // presencial, en_linea, mixta
  final String nombre;
  final String? urlCanal;

  Modalidad({
    required this.codigo,
    required this.nombre,
    this.urlCanal,
  });

  factory Modalidad.fromJson(Map<String, dynamic> json) {
    return Modalidad(
      codigo: json['codigo'] ?? json['tipo'] ?? 'presencial',
      nombre: json['nombre'] ?? 'Presencial',
      urlCanal: json['url_canal'] ?? json['url'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'codigo': codigo,
      'nombre': nombre,
      'url_canal': urlCanal,
    };
  }
}

class Fuente {
  final String titulo;
  final String url;
  final String? verificadoEn;

  Fuente({
    required this.titulo,
    required this.url,
    this.verificadoEn,
  });

  factory Fuente.fromJson(Map<String, dynamic> json) {
    return Fuente(
      titulo: json['titulo'] ?? json['nombre'] ?? 'Portal Oficial',
      url: json['url'] ?? '',
      verificadoEn: json['verificado_en'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'titulo': titulo,
      'url': url,
      'verificado_en': verificadoEn,
    };
  }
}

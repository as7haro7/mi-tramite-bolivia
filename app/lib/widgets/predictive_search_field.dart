import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_typeahead/flutter_typeahead.dart';
import '../config/theme.dart';
import '../models/tramite.dart';
import '../services/api_service.dart';

/// Widget de búsqueda con autocompletado predictivo
/// Muestra sugerencias en tiempo real mientras el usuario escribe
class PredictiveSearchField extends StatefulWidget {
  final String? initialValue;
  final String hintText;
  final Function(String) onSubmitted;
  final Function(Tramite)? onTramiteSelected;
  final bool autofocus;

  const PredictiveSearchField({
    super.key,
    this.initialValue,
    this.hintText = 'Buscar trámites, licencias, NIT...',
    required this.onSubmitted,
    this.onTramiteSelected,
    this.autofocus = false,
  });

  @override
  State<PredictiveSearchField> createState() => _PredictiveSearchFieldState();
}

class _PredictiveSearchFieldState extends State<PredictiveSearchField> {
  final TextEditingController _controller = TextEditingController();
  final FocusNode _focusNode = FocusNode();

  @override
  void initState() {
    super.initState();
    if (widget.initialValue != null) {
      _controller.text = widget.initialValue!;
    }
  }

  @override
  void dispose() {
    _controller.dispose();
    _focusNode.dispose();
    super.dispose();
  }

  Future<List<Tramite>> _getSuggestions(String query) async {
    if (query.isEmpty || query.length < 2) {
      return [];
    }

    try {
      // Fetch de trámites con el query
      final tramites = await ApiService.getTramites(query: query);
      
      // Retornar top 5 resultados más relevantes
      return tramites.take(5).toList();
    } catch (e) {
      // En caso de error, retornar lista vacía
      return [];
    }
  }

  @override
  Widget build(BuildContext context) {
    return TypeAheadField<Tramite>(
      controller: _controller,
      focusNode: _focusNode,
      
      // Builder para el campo de texto
      builder: (context, controller, focusNode) {
        return Container(
          decoration: BoxDecoration(
            color: AppTheme.iosLightSearchBg,
            borderRadius: BorderRadius.circular(10),
          ),
          padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 8),
          child: Row(
            children: [
              const Icon(
                CupertinoIcons.search,
                color: AppTheme.iosLightSubtext,
                size: 20,
              ),
              const SizedBox(width: 8),
              Expanded(
                child: TextField(
                  controller: controller,
                  focusNode: focusNode,
                  autofocus: widget.autofocus,
                  onSubmitted: (val) {
                    if (val.isNotEmpty) {
                      HapticFeedback.lightImpact();
                      widget.onSubmitted(val);
                    }
                  },
                  style: const TextStyle(
                    fontSize: 16,
                    color: AppTheme.iosLightText,
                  ),
                  decoration: InputDecoration(
                    hintText: widget.hintText,
                    hintStyle: const TextStyle(
                      color: AppTheme.iosLightSubtext,
                      fontSize: 16,
                    ),
                    border: InputBorder.none,
                    isDense: true,
                    contentPadding: EdgeInsets.zero,
                  ),
                ),
              ),
              if (controller.text.isNotEmpty)
                GestureDetector(
                  onTap: () {
                    controller.clear();
                    setState(() {});
                  },
                  child: const Icon(
                    CupertinoIcons.clear_thick_circled,
                    color: AppTheme.iosLightSubtext,
                    size: 18,
                  ),
                ),
            ],
          ),
        );
      },

      // Sugerencias basadas en el texto
      suggestionsCallback: (pattern) async {
        return await _getSuggestions(pattern);
      },

      // Builder para cada item de sugerencia
      itemBuilder: (context, Tramite tramite) {
        return ListTile(
          dense: true,
          leading: Container(
            width: 40,
            height: 40,
            decoration: BoxDecoration(
              color: AppTheme.primaryBlue.withAlpha(20),
              borderRadius: BorderRadius.circular(8),
            ),
            child: const Icon(
              CupertinoIcons.doc_text,
              color: AppTheme.primaryBlue,
              size: 20,
            ),
          ),
          title: Text(
            tramite.titulo,
            style: const TextStyle(
              fontSize: 14,
              fontWeight: FontWeight.w600,
            ),
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
          ),
          subtitle: Text(
            tramite.institucion?.nombre ?? 'Sin institución',
            style: const TextStyle(
              fontSize: 12,
              color: AppTheme.iosLightSubtext,
            ),
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
          ),
          trailing: const Icon(
            CupertinoIcons.arrow_up_left,
            size: 16,
            color: AppTheme.iosLightSubtext,
          ),
        );
      },

      // Cuando se selecciona una sugerencia
      onSelected: (Tramite tramite) {
        HapticFeedback.lightImpact();
        
        // Si hay callback de trámite seleccionado, usarlo
        if (widget.onTramiteSelected != null) {
          widget.onTramiteSelected!(tramite);
        } else {
          // Sino, hacer búsqueda del título
          _controller.text = tramite.titulo;
          widget.onSubmitted(tramite.titulo);
        }
      },

      // Configuración del dropdown
      decorationBuilder: (context, child) {
        return Material(
          elevation: 4,
          borderRadius: BorderRadius.circular(12),
          child: child,
        );
      },

      // Offset del dropdown
      offset: const Offset(0, 8),

      // Constraints del dropdown
      constraints: const BoxConstraints(
        maxHeight: 300,
      ),

      // Builder para cuando no hay sugerencias
      emptyBuilder: (context) {
        return Padding(
          padding: const EdgeInsets.all(16),
          child: Row(
            children: const [
              Icon(
                CupertinoIcons.search,
                color: AppTheme.iosLightSubtext,
                size: 20,
              ),
              SizedBox(width: 12),
              Expanded(
                child: Text(
                  'Escribe al menos 2 caracteres...',
                  style: TextStyle(
                    fontSize: 13,
                    color: AppTheme.iosLightSubtext,
                  ),
                ),
              ),
            ],
          ),
        );
      },

      // Builder para estado de carga
      loadingBuilder: (context) {
        return Padding(
          padding: const EdgeInsets.all(16),
          child: Row(
            children: const [
              SizedBox(
                width: 16,
                height: 16,
                child: CircularProgressIndicator(
                  strokeWidth: 2,
                  valueColor: AlwaysStoppedAnimation<Color>(AppTheme.primaryBlue),
                ),
              ),
              SizedBox(width: 12),
              Text(
                'Buscando trámites...',
                style: TextStyle(
                  fontSize: 13,
                  color: AppTheme.iosLightSubtext,
                ),
              ),
            ],
          ),
        );
      },

      // Builder para estado de error
      errorBuilder: (context, error) {
        return Padding(
          padding: const EdgeInsets.all(16),
          child: Row(
            children: const [
              Icon(
                CupertinoIcons.exclamationmark_triangle,
                color: AppTheme.alertRed,
                size: 20,
              ),
              SizedBox(width: 12),
              Expanded(
                child: Text(
                  'Error al buscar. Intenta nuevamente.',
                  style: TextStyle(
                    fontSize: 13,
                    color: AppTheme.alertRed,
                  ),
                ),
              ),
            ],
          ),
        );
      },

      // Duración del debounce (espera antes de buscar)
      debounceDuration: const Duration(milliseconds: 400),

      // Animación del dropdown
      transitionBuilder: (context, animation, child) {
        return FadeTransition(
          opacity: CurvedAnimation(
            parent: animation,
            curve: Curves.easeOut,
          ),
          child: child,
        );
      },
    );
  }
}

/// Widget compacto de búsqueda predictiva (sin bordes, para usar inline)
class PredictiveSearchFieldCompact extends StatelessWidget {
  final TextEditingController controller;
  final String hintText;
  final Function(String) onSubmitted;
  final Function(Tramite)? onTramiteSelected;

  const PredictiveSearchFieldCompact({
    super.key,
    required this.controller,
    this.hintText = 'Buscar...',
    required this.onSubmitted,
    this.onTramiteSelected,
  });

  Future<List<Tramite>> _getSuggestions(String query) async {
    if (query.isEmpty || query.length < 2) return [];
    
    try {
      final tramites = await ApiService.getTramites(query: query);
      return tramites.take(5).toList();
    } catch (e) {
      return [];
    }
  }

  @override
  Widget build(BuildContext context) {
    return TypeAheadField<Tramite>(
      controller: controller,
      
      builder: (context, controller, focusNode) {
        return TextField(
          controller: controller,
          focusNode: focusNode,
          onSubmitted: onSubmitted,
          decoration: InputDecoration(
            hintText: hintText,
            prefixIcon: const Icon(CupertinoIcons.search),
            suffixIcon: controller.text.isNotEmpty
                ? IconButton(
                    icon: const Icon(CupertinoIcons.clear),
                    onPressed: () => controller.clear(),
                  )
                : null,
          ),
        );
      },

      suggestionsCallback: _getSuggestions,

      itemBuilder: (context, tramite) {
        return ListTile(
          dense: true,
          leading: const Icon(CupertinoIcons.doc_text, size: 20),
          title: Text(
            tramite.titulo,
            style: const TextStyle(fontSize: 14),
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
          ),
          subtitle: Text(
            tramite.institucion?.nombre ?? 'Sin institución',
            style: const TextStyle(fontSize: 12),
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
          ),
        );
      },

      onSelected: (tramite) {
        if (onTramiteSelected != null) {
          onTramiteSelected!(tramite);
        } else {
          controller.text = tramite.titulo;
          onSubmitted(tramite.titulo);
        }
      },

      debounceDuration: const Duration(milliseconds: 400),
    );
  }
}

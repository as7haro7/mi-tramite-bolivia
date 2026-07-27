import 'package:flutter/material.dart';
import '../config/theme.dart';
import '../models/tramite.dart';
import '../services/api_service.dart';
import '../widgets/predictive_search_field.dart';
import '../widgets/tramite_card.dart';
import '../widgets/shimmer_loading.dart';

class SearchResultsScreen extends StatefulWidget {
  final String initialQuery;
  final ValueChanged<Tramite> onSelectTramite;
  final VoidCallback onOpenAssistant;

  const SearchResultsScreen({
    super.key,
    required this.initialQuery,
    required this.onSelectTramite,
    required this.onOpenAssistant,
  });

  @override
  State<SearchResultsScreen> createState() => _SearchResultsScreenState();
}

class _SearchResultsScreenState extends State<SearchResultsScreen> {
  late TextEditingController _searchController;
  List<Tramite> _results = [];
  bool _isLoading = true;
  String _selectedModalidad = 'todas';

  @override
  void initState() {
    super.initState();
    _searchController = TextEditingController(text: widget.initialQuery);
    _performSearch(widget.initialQuery);
  }

  Future<void> _performSearch(String query) async {
    setState(() => _isLoading = true);
    final list = await ApiService.getTramites(
      query: query,
      modalidad: _selectedModalidad,
    );
    setState(() {
      _results = list;
      _isLoading = false;
    });
  }

  Future<void> _handleRefresh() async {
    return _performSearch(_searchController.text);
  }

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Resultados de Búsqueda'),
      ),
      body: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              children: [
                PredictiveSearchFieldCompact(
                  controller: _searchController,
                  hintText: 'Buscar por trámite, sigla o institución...',
                  onSubmitted: _performSearch,
                  onTramiteSelected: (tramite) {
                    widget.onSelectTramite(tramite);
                  },
                ),
                const SizedBox(height: 10),
                Row(
                  children: [
                    const Text('Modalidad: ', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 13)),
                    const SizedBox(width: 8),
                    ChoiceChip(
                      label: const Text('Todas'),
                      selected: _selectedModalidad == 'todas',
                      onSelected: (val) {
                        if (val) {
                          setState(() => _selectedModalidad = 'todas');
                          _performSearch(_searchController.text);
                        }
                      },
                    ),
                    const SizedBox(width: 6),
                    ChoiceChip(
                      label: const Text('Presencial'),
                      selected: _selectedModalidad == 'presencial',
                      onSelected: (val) {
                        if (val) {
                          setState(() => _selectedModalidad = 'presencial');
                          _performSearch(_searchController.text);
                        }
                      },
                    ),
                    const SizedBox(width: 6),
                    ChoiceChip(
                      label: const Text('En Línea'),
                      selected: _selectedModalidad == 'en_linea',
                      onSelected: (val) {
                        if (val) {
                          setState(() => _selectedModalidad = 'en_linea');
                          _performSearch(_searchController.text);
                        }
                      },
                    ),
                  ],
                ),
              ],
            ),
          ),
          Expanded(
            child: RefreshIndicator(
              onRefresh: _handleRefresh,
              color: AppTheme.primaryBlue,
              child: _isLoading
                  ? ListView.builder(
                      padding: const EdgeInsets.symmetric(horizontal: 16),
                      itemCount: 4,
                      itemBuilder: (context, index) => const TramiteCardShimmer(),
                    )
                  : _results.isEmpty
                      ? SingleChildScrollView(
                          physics: const AlwaysScrollableScrollPhysics(),
                          child: Padding(
                            padding: const EdgeInsets.all(24.0),
                            child: Column(
                              mainAxisAlignment: MainAxisAlignment.center,
                              children: [
                                const Icon(Icons.search_off_rounded, size: 60, color: Colors.grey),
                                const SizedBox(height: 16),
                                Text(
                                  'No encontramos un trámite verificado para "${_searchController.text}".',
                                  textAlign: TextAlign.center,
                                  style: Theme.of(context).textTheme.titleMedium,
                                ),
                                const SizedBox(height: 20),
                                ElevatedButton.icon(
                                  icon: const Icon(Icons.auto_awesome_rounded),
                                  label: const Text('Preguntar al Asistente IA'),
                                  onPressed: widget.onOpenAssistant,
                                ),
                              ],
                            ),
                          ),
                        )
                      : ListView.builder(
                          padding: const EdgeInsets.symmetric(horizontal: 16),
                          itemCount: _results.length,
                          itemBuilder: (context, index) {
                            final item = _results[index];
                            return TramiteCard(
                              tramite: item,
                              onTap: () => widget.onSelectTramite(item),
                            );
                          },
                        ),
            ),
          ),
        ],
      ),
    );
  }
}

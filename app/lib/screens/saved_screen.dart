import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../config/theme.dart';
import '../models/tramite.dart';
import '../providers/checklist_provider.dart';
import '../providers/favorites_provider.dart';
import '../services/api_service.dart';
import '../widgets/premium_cta_card.dart';
import '../widgets/tramite_card.dart';
import 'checklist_screen.dart';

class SavedScreen extends StatefulWidget {
  final ValueChanged<Tramite> onSelectTramite;
  final VoidCallback onOpenPremium;

  const SavedScreen({
    super.key,
    required this.onSelectTramite,
    required this.onOpenPremium,
  });

  @override
  State<SavedScreen> createState() => _SavedScreenState();
}

class _SavedScreenState extends State<SavedScreen> with SingleTickerProviderStateMixin {
  late TabController _tabController;
  List<Tramite> _favoriteTramites = [];
  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
    _loadFavorites();
  }

  Future<void> _loadFavorites() async {
    setState(() => _isLoading = true);
    final favSlugs = Provider.of<FavoritesProvider>(context, listen: false).favoriteSlugs;
    final all = await ApiService.getTramites();
    setState(() {
      _favoriteTramites = all.where((t) => favSlugs.contains(t.slug)).toList();
      _isLoading = false;
    });
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final checklistProvider = Provider.of<ChecklistProvider>(context);
    final checklists = checklistProvider.checklists;

    return Scaffold(
      appBar: AppBar(
        title: const Text('Mis Guardados'),
        bottom: TabBar(
          controller: _tabController,
          labelColor: AppTheme.primaryBlue,
          indicatorColor: AppTheme.primaryBlue,
          tabs: const [
            Tab(text: 'Checklists Activos'),
            Tab(text: 'Trámites Favoritos'),
          ],
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: [
          // Checklists Tab
          ListView(
            padding: const EdgeInsets.all(16),
            children: [
              PremiumCtaCard(
                onTap: widget.onOpenPremium,
                title: 'Sincronizar en la Nube (Premium)',
                subtitle: 'Accede a tus checklists desde cualquier dispositivo móvil o web.',
              ),
              const SizedBox(height: 16),
              if (checklists.isEmpty)
                const Padding(
                  padding: EdgeInsets.all(32.0),
                  child: Column(
                    children: [
                      Icon(Icons.checklist_rtl_rounded, size: 60, color: Colors.grey),
                      SizedBox(height: 12),
                      Text(
                        'No tienes checklists activos.',
                        style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
                      ),
                      SizedBox(height: 4),
                      Text(
                        'Abre la ficha de cualquier trámite y toca "Preparar mi Checklist Local".',
                        textAlign: TextAlign.center,
                        style: TextStyle(color: Colors.grey, fontSize: 13),
                      ),
                    ],
                  ),
                )
              else
                ...checklists.map((ch) {
                  return Card(
                    margin: const EdgeInsets.only(bottom: 12),
                    child: InkWell(
                      borderRadius: BorderRadius.circular(16),
                      onTap: () {
                        Navigator.push(
                          context,
                          MaterialPageRoute(builder: (_) => ChecklistScreen(checklist: ch)),
                        );
                      },
                      child: Padding(
                        padding: const EdgeInsets.all(16),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Row(
                              mainAxisAlignment: MainAxisAlignment.spaceBetween,
                              children: [
                                Text(
                                  ch.institucionNombre,
                                  style: const TextStyle(
                                    color: AppTheme.primaryBlue,
                                    fontWeight: FontWeight.bold,
                                    fontSize: 12,
                                  ),
                                ),
                                Text(
                                  '${(ch.progress * 100).toInt()}% completado',
                                  style: TextStyle(
                                    color: ch.isFullyCompleted
                                        ? AppTheme.successGreen
                                        : AppTheme.primaryBlue,
                                    fontWeight: FontWeight.bold,
                                    fontSize: 12.5,
                                  ),
                                ),
                              ],
                            ),
                            const SizedBox(height: 6),
                            Text(
                              ch.tramiteTitulo,
                              style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 15),
                            ),
                            const SizedBox(height: 10),
                            ClipRRect(
                              borderRadius: BorderRadius.circular(4),
                              child: LinearProgressIndicator(
                                value: ch.progress,
                                minHeight: 8,
                                backgroundColor: Colors.grey.shade200,
                                color: ch.isFullyCompleted
                                    ? AppTheme.successGreen
                                    : AppTheme.primaryBlue,
                              ),
                            ),
                          ],
                        ),
                      ),
                    ),
                  );
                }),
            ],
          ),

          // Favoritos Tab
          RefreshIndicator(
            onRefresh: _loadFavorites,
            child: ListView(
              padding: const EdgeInsets.all(16),
              children: [
                if (_isLoading)
                  const Padding(padding: EdgeInsets.all(20), child: Center(child: CircularProgressIndicator()))
                else if (_favoriteTramites.isEmpty)
                  const Padding(
                    padding: EdgeInsets.all(32.0),
                    child: Column(
                      children: [
                        Icon(Icons.star_outline_rounded, size: 60, color: Colors.grey),
                        SizedBox(height: 12),
                        Text(
                          'No tienes trámites en favoritos.',
                          style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
                        ),
                        SizedBox(height: 4),
                        Text(
                          'Toca la estrella ⭐ en cualquier trámite para guardarlo aquí.',
                          textAlign: TextAlign.center,
                          style: TextStyle(color: Colors.grey, fontSize: 13),
                        ),
                      ],
                    ),
                  )
                else
                  ..._favoriteTramites.map((t) {
                    return TramiteCard(
                      tramite: t,
                      onTap: () => widget.onSelectTramite(t),
                    );
                  }),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

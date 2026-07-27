import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:url_launcher/url_launcher.dart';
import '../config/theme.dart';
import '../models/tramite.dart';
import '../providers/checklist_provider.dart';
import '../providers/favorites_provider.dart';
import '../widgets/cost_card.dart';
import '../widgets/disclaimer_banner.dart';
import '../widgets/source_card.dart';
import '../widgets/verification_seal.dart';
import 'checklist_screen.dart';
import 'report_form_screen.dart';

class TramiteDetailScreen extends StatefulWidget {
  final Tramite tramite;
  final VoidCallback onChecklistCreated;

  const TramiteDetailScreen({
    super.key,
    required this.tramite,
    required this.onChecklistCreated,
  });

  @override
  State<TramiteDetailScreen> createState() => _TramiteDetailScreenState();
}

class _TramiteDetailScreenState extends State<TramiteDetailScreen> with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 4, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  Future<void> _openMap(double? lat, double? lng, String? dir) async {
    final query = lat != null && lng != null ? '$lat,$lng' : Uri.encodeComponent(dir ?? 'La Paz Bolivia');
    final url = Uri.parse('https://www.google.com/maps/search/?api=1&query=$query');
    if (await canLaunchUrl(url)) {
      await launchUrl(url, mode: LaunchMode.externalApplication);
    }
  }

  @override
  Widget build(BuildContext context) {
    final t = widget.tramite;
    final favProvider = Provider.of<FavoritesProvider>(context);
    final checklistProvider = Provider.of<ChecklistProvider>(context);
    final isFav = favProvider.isFavorite(t.slug);
    final hasChecklist = checklistProvider.getChecklistBySlug(t.slug) != null;

    final instSigla = t.institucion?.sigla ?? t.institucion?.nombre ?? 'Institución';

    return Scaffold(
      appBar: AppBar(
        title: Text(instSigla),
        actions: [
          IconButton(
            icon: Icon(
              isFav ? Icons.star_rounded : Icons.star_outline_rounded,
              color: isFav ? Colors.amber : null,
            ),
            onPressed: () => favProvider.toggleFavorite(t.slug),
          ),
          IconButton(
            icon: const Icon(Icons.report_problem_outlined),
            tooltip: 'Reportar dato desactualizado',
            onPressed: () {
              Navigator.push(
                context,
                MaterialPageRoute(
                  builder: (_) => ReportFormScreen(tramiteSlug: t.slug, tramiteTitulo: t.titulo),
                ),
              );
            },
          ),
        ],
      ),
      body: Column(
        children: [
          // Header Summary
          Container(
            width: double.infinity,
            padding: const EdgeInsets.all(16),
            color: Theme.of(context).cardTheme.color,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                VerificationSeal(date: t.verificadoEn, isVerified: t.verificado),
                const SizedBox(height: 8),
                Text(
                  t.titulo,
                  style: Theme.of(context).textTheme.headlineMedium?.copyWith(
                        fontSize: 20,
                      ),
                ),
                if (t.resumen != null && t.resumen!.isNotEmpty) ...[
                  const SizedBox(height: 6),
                  Text(
                    t.resumen!,
                    style: Theme.of(context).textTheme.bodyMedium,
                  ),
                ],
                const SizedBox(height: 12),
                Row(
                  children: [
                    if (t.plazoEstimado != null) ...[
                      const Icon(Icons.schedule_rounded, size: 16, color: AppTheme.secondaryTeal),
                      const SizedBox(width: 4),
                      Text(
                        'Plazo: ${t.plazoEstimado}',
                        style: const TextStyle(fontSize: 12.5, fontWeight: FontWeight.w600),
                      ),
                      const SizedBox(width: 16),
                    ],
                    const Icon(Icons.location_city_rounded, size: 16, color: AppTheme.primaryBlue),
                    const SizedBox(width: 4),
                    Expanded(
                      child: Text(
                        t.institucion?.nombre ?? 'Institución Pública',
                        style: const TextStyle(fontSize: 12.5, fontWeight: FontWeight.w600),
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 12),
                SizedBox(
                  width: double.infinity,
                  child: ElevatedButton.icon(
                    icon: Icon(hasChecklist ? Icons.check_circle_rounded : Icons.add_task_rounded),
                    label: Text(hasChecklist ? 'Ver mi Checklist Activo' : 'Preparar mi Checklist Local'),
                    style: ElevatedButton.styleFrom(
                      backgroundColor: hasChecklist ? AppTheme.successGreen : AppTheme.primaryBlue,
                    ),
                    onPressed: () async {
                      await checklistProvider.createOrUpdateFromTramite(t);
                      widget.onChecklistCreated();
                      final ch = checklistProvider.getChecklistBySlug(t.slug);
                      if (ch != null && context.mounted) {
                        Navigator.push(
                          context,
                          MaterialPageRoute(
                            builder: (_) => ChecklistScreen(checklist: ch),
                          ),
                        );
                      }
                    },
                  ),
                ),
              ],
            ),
          ),

          // Tabs Header
          TabBar(
            controller: _tabController,
            labelColor: AppTheme.primaryBlue,
            indicatorColor: AppTheme.primaryBlue,
            tabs: const [
              Tab(text: 'Requisitos'),
              Tab(text: 'Pasos'),
              Tab(text: 'Costos'),
              Tab(text: 'Oficinas'),
            ],
          ),

          // Tabs View Content
          Expanded(
            child: TabBarView(
              controller: _tabController,
              children: [
                // Requisitos Tab
                ListView(
                  padding: const EdgeInsets.all(16),
                  children: [
                    const DisclaimerBanner(isCompact: true),
                    const SizedBox(height: 16),
                    if (t.requisitos.isEmpty)
                      const Text('No se detallaron requisitos específicos.')
                    else
                      ...t.requisitos.map((req) {
                        return Container(
                          margin: const EdgeInsets.only(bottom: 10),
                          padding: const EdgeInsets.all(12),
                          decoration: BoxDecoration(
                            color: Theme.of(context).cardTheme.color,
                            borderRadius: BorderRadius.circular(12),
                            border: Border.all(color: AppTheme.lightBorder),
                          ),
                          child: Row(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              const Icon(Icons.check_box_outline_blank_rounded, color: AppTheme.primaryBlue, size: 20),
                              const SizedBox(width: 12),
                              Expanded(
                                child: Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    Text(
                                      req.nombre,
                                      style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 14),
                                    ),
                                    if (req.descripcion != null && req.descripcion!.isNotEmpty) ...[
                                      const SizedBox(height: 4),
                                      Text(
                                        req.descripcion!,
                                        style: TextStyle(fontSize: 13, color: Theme.of(context).colorScheme.onSurface.withAlpha(180)),
                                      ),
                                    ],
                                  ],
                                ),
                              ),
                            ],
                          ),
                        );
                      }),
                  ],
                ),

                // Pasos Tab
                ListView(
                  padding: const EdgeInsets.all(16),
                  children: [
                    if (t.pasos.isEmpty)
                      const Text('No se detallaron pasos para este trámite.')
                    else
                      ...t.pasos.map((paso) {
                        return Container(
                          margin: const EdgeInsets.only(bottom: 12),
                          padding: const EdgeInsets.all(12),
                          decoration: BoxDecoration(
                            color: Theme.of(context).cardTheme.color,
                            borderRadius: BorderRadius.circular(12),
                            border: Border.all(color: AppTheme.lightBorder),
                          ),
                          child: Row(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              CircleAvatar(
                                radius: 14,
                                backgroundColor: AppTheme.primaryBlue,
                                child: Text(
                                  '${paso.orden}',
                                  style: const TextStyle(color: Colors.white, fontWeight: FontWeight.bold, fontSize: 12),
                                ),
                              ),
                              const SizedBox(width: 12),
                              Expanded(
                                child: Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    Text(
                                      paso.titulo,
                                      style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 14),
                                    ),
                                    if (paso.descripcion != null && paso.descripcion!.isNotEmpty) ...[
                                      const SizedBox(height: 4),
                                      Text(
                                        paso.descripcion!,
                                        style: TextStyle(fontSize: 13, color: Theme.of(context).colorScheme.onSurface.withAlpha(180)),
                                      ),
                                    ],
                                  ],
                                ),
                              ),
                            ],
                          ),
                        );
                      }),
                  ],
                ),

                // Costos Tab
                ListView(
                  padding: const EdgeInsets.all(16),
                  children: [
                    if (t.costos.isEmpty)
                      Container(
                        padding: const EdgeInsets.all(14),
                        decoration: BoxDecoration(
                          color: AppTheme.successGreen.withAlpha(20),
                          borderRadius: BorderRadius.circular(12),
                        ),
                        child: const Text('Este trámite está catalogado como GRATUITO.'),
                      )
                    else
                      ...t.costos.map((costo) => CostCard(costo: costo)),
                    const SizedBox(height: 16),
                    if (t.fuentes.isNotEmpty) ...[
                      const Text('Fuentes Oficiales Consultadas:', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 13)),
                      const SizedBox(height: 6),
                      ...t.fuentes.map((f) => SourceCard(title: f.titulo, url: f.url, date: f.verificadoEn)),
                    ],
                  ],
                ),

                // Oficinas Tab
                ListView(
                  padding: const EdgeInsets.all(16),
                  children: [
                    if (t.oficinas.isEmpty)
                      const Text('Atención disponible en las oficinas generales de la institución.')
                    else
                      ...t.oficinas.map((oficina) {
                        return Container(
                          margin: const EdgeInsets.only(bottom: 12),
                          padding: const EdgeInsets.all(14),
                          decoration: BoxDecoration(
                            color: Theme.of(context).cardTheme.color,
                            borderRadius: BorderRadius.circular(12),
                            border: Border.all(color: AppTheme.lightBorder),
                          ),
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Row(
                                children: [
                                  const Icon(Icons.business_rounded, color: AppTheme.primaryBlue, size: 20),
                                  const SizedBox(width: 8),
                                  Expanded(
                                    child: Text(
                                      oficina.nombre,
                                      style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 14),
                                    ),
                                  ),
                                ],
                              ),
                              if (oficina.direccion != null) ...[
                                const SizedBox(height: 6),
                                Row(
                                  children: [
                                    const Icon(Icons.location_on_outlined, size: 16, color: AppTheme.iosLightSubtext),
                                    const SizedBox(width: 4),
                                    Expanded(
                                      child: Text(
                                        oficina.direccion!,
                                        style: const TextStyle(fontSize: 13),
                                      ),
                                    ),
                                  ],
                                ),
                              ],
                              if (oficina.horario != null) ...[
                                const SizedBox(height: 4),
                                Row(
                                  children: [
                                    const Icon(Icons.access_time_rounded, size: 15, color: AppTheme.iosLightSubtext),
                                    const SizedBox(width: 4),
                                    Expanded(
                                      child: Text(
                                        oficina.horario!,
                                        style: const TextStyle(fontSize: 12.5, color: AppTheme.iosLightSubtext),
                                      ),
                                    ),
                                  ],
                                ),
                              ],
                              const SizedBox(height: 10),
                              OutlinedButton.icon(
                                icon: const Icon(Icons.map_rounded, size: 16),
                                label: const Text('Abrir Ubicación en Mapas'),
                                onPressed: () => _openMap(oficina.latitud, oficina.longitud, oficina.direccion),
                              ),
                            ],
                          ),
                        );
                      }),
                  ],
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

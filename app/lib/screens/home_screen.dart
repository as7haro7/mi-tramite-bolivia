import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';
import '../config/ad_config.dart';
import '../config/theme.dart';
import '../models/categoria.dart';
import '../models/tramite.dart';
import '../providers/app_provider.dart';
import '../providers/checklist_provider.dart';
import '../services/api_service.dart';
import '../services/error_handler.dart';
import '../widgets/ad_banner_widget.dart';
import '../widgets/disclaimer_banner.dart';
import '../widgets/premium_cta_card.dart';
import '../widgets/predictive_search_field.dart';
import '../widgets/shimmer_loading.dart';
import '../widgets/tramite_card.dart';

class HomeScreen extends StatefulWidget {
  final Function(String search) onSearchSubmitted;
  final VoidCallback onOpenAssistant;
  final ValueChanged<Tramite> onSelectTramite;
  final VoidCallback onOpenPremium;

  const HomeScreen({
    super.key,
    required this.onSearchSubmitted,
    required this.onOpenAssistant,
    required this.onSelectTramite,
    required this.onOpenPremium,
  });

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  List<Categoria> _categorias = [];
  List<Tramite> _tramitesFrecuentes = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    setState(() => _isLoading = true);
    
    try {
      final cats = await ApiService.getCategorias();
      final trs = await ApiService.getTramites();
      
      if (mounted) {
        setState(() {
          _categorias = cats;
          _tramitesFrecuentes = trs;
          _isLoading = false;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() => _isLoading = false);
        // No mostramos SnackBar en carga inicial, solo en UI
      }
    }
  }

  Future<void> _handleRefresh() async {
    try {
      final cats = await ApiService.getCategorias();
      final trs = await ApiService.getTramites();
      
      setState(() {
        _categorias = cats;
        _tramitesFrecuentes = trs;
      });
      
      if (mounted) {
        ErrorHandler.showSuccess(context, 'Datos actualizados');
      }
    } catch (e) {
      if (mounted) {
        ErrorHandler.showError(
          context,
          ErrorHandler.parseError(e),
          onRetry: _handleRefresh,
        );
      }
    }
  }

  IconData _getCategoryIconData(String codigo) {
    switch (codigo.toLowerCase()) {
      case 'identidad':
        return CupertinoIcons.person_crop_square;
      case 'negocio':
        return CupertinoIcons.briefcase;
      case 'impuestos':
        return CupertinoIcons.doc_plaintext;
      case 'vehiculos':
        return CupertinoIcons.car_detailed;
      case 'salud':
        return CupertinoIcons.heart;
      case 'educacion':
        return CupertinoIcons.book;
      default:
        return CupertinoIcons.folder;
    }
  }

  @override
  void dispose() {
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final appProvider = Provider.of<AppProvider>(context);
    final checklistProvider = Provider.of<ChecklistProvider>(context);
    final activeChecklists = checklistProvider.checklists;
    final isDark = Theme.of(context).brightness == Brightness.dark;

    return Scaffold(
      body: SafeArea(
        bottom: false,
        child: RefreshIndicator(
          onRefresh: _handleRefresh,
          color: AppTheme.primaryBlue,
          backgroundColor: Colors.white,
          child: CustomScrollView(
            physics: const BouncingScrollPhysics(),
            slivers: [
            // iOS Cupertino Large Title Navigation Bar Header
            SliverToBoxAdapter(
              child: Padding(
                padding: const EdgeInsets.fromLTRB(16, 16, 16, 8),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Row(
                          children: [
                            const Icon(CupertinoIcons.location_fill, size: 14, color: AppTheme.tintColor),
                            const SizedBox(width: 4),
                            Text(
                              appProvider.selectedCity.toUpperCase(),
                              style: const TextStyle(
                                fontWeight: FontWeight.w600,
                                fontSize: 12,
                                color: AppTheme.iosLightSubtext,
                                letterSpacing: -0.2,
                              ),
                            ),
                          ],
                        ),
                        CupertinoButton(
                          padding: EdgeInsets.zero,
                          minimumSize: Size.zero,
                          onPressed: () {
                            HapticFeedback.lightImpact();
                            widget.onOpenPremium();
                          },
                          child: const Row(
                            children: [
                              Icon(CupertinoIcons.star_circle_fill, color: AppTheme.systemOrange, size: 22),
                              SizedBox(width: 4),
                              Text(
                                'Premium',
                                style: TextStyle(
                                  color: AppTheme.systemOrange,
                                  fontWeight: FontWeight.w600,
                                  fontSize: 14,
                                ),
                              ),
                            ],
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 6),
                    Text(
                      'Mi Trámite',
                      style: Theme.of(context).textTheme.headlineLarge,
                    ),
                  ],
                ),
              ),
            ),

            SliverPadding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              sliver: SliverList(
                delegate: SliverChildListDelegate([
                  const DisclaimerBanner(isCompact: true),
                  const SizedBox(height: 14),

                  // iOS HIG Search Bar con Autocompletado Predictivo
                  PredictiveSearchField(
                    hintText: 'Buscar trámites, licencias, NIT...',
                    onSubmitted: (query) {
                      HapticFeedback.lightImpact();
                      widget.onSearchSubmitted(query);
                    },
                    onTramiteSelected: (tramite) {
                      HapticFeedback.lightImpact();
                      widget.onSelectTramite(tramite);
                    },
                  ),
                  const SizedBox(height: 18),

                  // AI Assistant iOS Row Callout
                  GestureDetector(
                    onTap: () {
                      HapticFeedback.lightImpact();
                      widget.onOpenAssistant();
                    },
                    child: Container(
                      decoration: BoxDecoration(
                        color: AppTheme.tintColor.withAlpha(isDark ? 30 : 15),
                        borderRadius: BorderRadius.circular(14),
                        border: Border.all(color: AppTheme.tintColor.withAlpha(40), width: 0.8),
                      ),
                      padding: const EdgeInsets.all(14),
                      child: Row(
                        children: [
                          Container(
                            padding: const EdgeInsets.all(8),
                            decoration: const BoxDecoration(
                              color: AppTheme.tintColor,
                              shape: BoxShape.circle,
                            ),
                            child: const Icon(CupertinoIcons.sparkles, color: Colors.white, size: 18),
                          ),
                          const SizedBox(width: 12),
                          Expanded(
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: const [
                                Text(
                                  'Asistente IA RAG',
                                  style: TextStyle(fontWeight: FontWeight.bold, fontSize: 15),
                                ),
                                SizedBox(height: 2),
                                Text(
                                  'Consulta dudas legales con fuentes oficiales.',
                                  style: TextStyle(fontSize: 13, color: AppTheme.iosLightSubtext),
                                ),
                              ],
                            ),
                          ),
                          const Icon(CupertinoIcons.chevron_right, color: AppTheme.tintColor, size: 18),
                        ],
                      ),
                    ),
                  ),
                  const SizedBox(height: 20),

                  // Active Checklists Widget
                  if (activeChecklists.isNotEmpty) ...[
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Text('Checklists Pendientes', style: Theme.of(context).textTheme.titleLarge),
                        Text(
                          '${activeChecklists.length} activos',
                          style: const TextStyle(fontSize: 13, color: AppTheme.tintColor, fontWeight: FontWeight.w600),
                        ),
                      ],
                    ),
                    const SizedBox(height: 10),
                    ...activeChecklists.take(2).map((ch) {
                      return Card(
                        margin: const EdgeInsets.only(bottom: 8),
                        child: Padding(
                          padding: const EdgeInsets.all(14),
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Row(
                                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                                children: [
                                  Expanded(
                                    child: Text(
                                      ch.tramiteTitulo,
                                      style: const TextStyle(fontWeight: FontWeight.w600, fontSize: 15),
                                      maxLines: 1,
                                      overflow: TextOverflow.ellipsis,
                                    ),
                                  ),
                                  Text(
                                    '${(ch.progress * 100).toInt()}%',
                                    style: const TextStyle(fontWeight: FontWeight.bold, color: AppTheme.systemGreen, fontSize: 13),
                                  ),
                                ],
                              ),
                              const SizedBox(height: 8),
                              ClipRRect(
                                borderRadius: BorderRadius.circular(4),
                                child: LinearProgressIndicator(
                                  value: ch.progress,
                                  backgroundColor: isDark ? AppTheme.iosDarkBorder : AppTheme.iosLightSearchBg,
                                  color: AppTheme.systemGreen,
                                  minHeight: 6,
                                ),
                              ),
                            ],
                          ),
                        ),
                      );
                    }),
                    const SizedBox(height: 16),
                  ],

                  // Premium CTA Card
                  PremiumCtaCard(onTap: widget.onOpenPremium),
                  const SizedBox(height: 20),

                  // Ad Banner después del CTA Premium
                  if (AdConfig.shouldShowAds)
                    LabeledAdBanner(
                      adUnitId: AdConfig.bannerAdUnitId,
                      margin: const EdgeInsets.symmetric(vertical: 8),
                    ),

                  // Categories Header
                  Text('Categorías', style: Theme.of(context).textTheme.titleLarge),
                  const SizedBox(height: 10),
                  if (_isLoading)
                    const CategoryChipShimmer()
                  else
                    SizedBox(
                      height: 38,
                      child: ListView.builder(
                        scrollDirection: Axis.horizontal,
                        itemCount: _categorias.length,
                        itemBuilder: (context, index) {
                          final cat = _categorias[index];
                          return Padding(
                            padding: const EdgeInsets.only(right: 8),
                            child: CupertinoButton(
                              padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 6),
                              color: isDark ? AppTheme.iosDarkSearchBg : AppTheme.iosLightSearchBg,
                              borderRadius: BorderRadius.circular(18),
                              onPressed: () {
                                HapticFeedback.lightImpact();
                                widget.onSearchSubmitted(cat.nombre);
                              },
                              child: Row(
                                children: [
                                  Icon(
                                    _getCategoryIconData(cat.codigo),
                                    size: 16,
                                    color: AppTheme.tintColor,
                                  ),
                                  const SizedBox(width: 6),
                                  Text(
                                    cat.nombre,
                                    style: TextStyle(
                                      fontSize: 13,
                                      fontWeight: FontWeight.w500,
                                      color: isDark ? AppTheme.iosDarkText : AppTheme.iosLightText,
                                    ),
                                  ),
                                ],
                              ),
                            ),
                          );
                        },
                      ),
                    ),
                  const SizedBox(height: 24),

                  // Frequent Trámites Header
                  Text('Trámites Frecuentes', style: Theme.of(context).textTheme.titleLarge),
                  const SizedBox(height: 12),
                  if (_isLoading)
                    ...List.generate(3, (index) => const TramiteCardShimmer())
                  else if (_tramitesFrecuentes.isEmpty)
                    const Padding(
                      padding: EdgeInsets.all(20),
                      child: Center(
                        child: Text('No hay trámites en la base de datos.'),
                      ),
                    )
                  else
                    Column(
                      children: _tramitesFrecuentes.map((t) {
                        return TramiteCard(
                          tramite: t,
                          onTap: () => widget.onSelectTramite(t),
                        );
                      }).toList(),
                    ),
                  const SizedBox(height: 90), // Space for bottom translucent tab bar
                ]),
              ),
            ),
          ],
        ),
      ),
    ),
  );
}
}

import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../config/theme.dart';
import '../providers/app_provider.dart';
import '../widgets/qr_modal.dart';

class PremiumPlanScreen extends StatefulWidget {
  const PremiumPlanScreen({super.key});

  @override
  State<PremiumPlanScreen> createState() => _PremiumPlanScreenState();
}

class _PremiumPlanScreenState extends State<PremiumPlanScreen> {
  void _openQrModal() {
    showDialog(
      context: context,
      builder: (ctx) => QrPaymentModal(
        onPaymentSuccess: () async {
          final provider = Provider.of<AppProvider>(context, listen: false);
          await provider.setPremiumStatus(true);
          if (!mounted) return;
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(
              backgroundColor: AppTheme.successGreen,
              content: Row(
                children: [
                  Icon(Icons.stars_rounded, color: Colors.white),
                  SizedBox(width: 10),
                  Expanded(child: Text('¡Suscripción Premium activada exitosamente!')),
                ],
              ),
            ),
          );
        },
      ),
    );
  }

  void _downgradePlan() async {
    final provider = Provider.of<AppProvider>(context, listen: false);
    await provider.setPremiumStatus(false);
    if (!mounted) return;
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Has vuelto al Plan Gratuito.')),
    );
  }

  @override
  Widget build(BuildContext context) {
    final appProvider = Provider.of<AppProvider>(context);
    final isPremium = appProvider.isPremium;

    return Scaffold(
      appBar: AppBar(
        title: const Text('Planes y Suscripción'),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.center,
          children: [
            // Hero Icon & Header
            Container(
              padding: const EdgeInsets.all(18),
              decoration: BoxDecoration(
                color: AppTheme.tintColor,
                shape: BoxShape.circle,
                boxShadow: [
                  BoxShadow(color: AppTheme.tintColor.withAlpha(40), blurRadius: 10, offset: const Offset(0, 4)),
                ],
              ),
              child: const Icon(
                Icons.workspace_premium_rounded,
                color: Colors.amberAccent,
                size: 54,
              ),
            ),
            const SizedBox(height: 18),
            Text(
              'Organiza mejor tus trámites',
              textAlign: TextAlign.center,
              style: Theme.of(context).textTheme.headlineMedium,
            ),
            const SizedBox(height: 6),
            const Text(
              'Demostración interactiva del modelo de suscripción (Bs. 20/mes)',
              textAlign: TextAlign.center,
              style: TextStyle(color: Colors.grey, fontSize: 13),
            ),
            const SizedBox(height: 24),

            // Plan Comparison Cards
            Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Free Plan Card
                Expanded(
                  child: Container(
                    padding: const EdgeInsets.all(16),
                    decoration: BoxDecoration(
                      color: Theme.of(context).cardTheme.color,
                      borderRadius: BorderRadius.circular(16),
                      border: Border.all(
                        color: !isPremium ? AppTheme.primaryBlue : AppTheme.lightBorder,
                        width: !isPremium ? 2 : 1,
                      ),
                    ),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Row(
                          children: [
                            const Text(
                              'Gratis',
                              style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
                            ),
                            if (!isPremium) ...[
                              const Spacer(),
                              const Icon(Icons.check_circle_rounded, color: AppTheme.primaryBlue, size: 18),
                            ],
                          ],
                        ),
                        const SizedBox(height: 4),
                        const Text(
                          'Bs. 0/mes',
                          style: TextStyle(fontWeight: FontWeight.bold, fontSize: 18, color: AppTheme.primaryBlue),
                        ),
                        const Divider(height: 20),
                        const FeatureRow(text: 'Información y fuentes oficiales', active: true),
                        const FeatureRow(text: 'Checklist local offline', active: true),
                        const FeatureRow(text: 'Asistente IA RAG', active: true),
                        const FeatureRow(text: 'Sincronización en la nube', active: false),
                        const FeatureRow(text: 'Perfiles familiares', active: false),
                        const FeatureRow(text: 'Alertas de cambios', active: false),
                      ],
                    ),
                  ),
                ),
                const SizedBox(width: 12),

                // Premium Plan Card
                Expanded(
                  child: Container(
                    padding: const EdgeInsets.all(16),
                    decoration: BoxDecoration(
                      color: isPremium ? AppTheme.successGreen.withAlpha(15) : AppTheme.primaryBlue.withAlpha(10),
                      borderRadius: BorderRadius.circular(16),
                      border: Border.all(
                        color: isPremium ? AppTheme.successGreen : AppTheme.secondaryTeal,
                        width: 2,
                      ),
                    ),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Row(
                          children: [
                            const Text(
                              'Premium',
                              style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
                            ),
                            if (isPremium) ...[
                              const Spacer(),
                              const Icon(Icons.check_circle_rounded, color: AppTheme.successGreen, size: 18),
                            ],
                          ],
                        ),
                        const SizedBox(height: 4),
                        const Text(
                          'Bs. 20/mes',
                          style: TextStyle(fontWeight: FontWeight.bold, fontSize: 18, color: AppTheme.secondaryTeal),
                        ),
                        const Divider(height: 20),
                        const FeatureRow(text: 'Información y fuentes oficiales', active: true),
                        const FeatureRow(text: 'Checklist local offline', active: true),
                        const FeatureRow(text: 'Asistente IA RAG', active: true),
                        const FeatureRow(text: 'Sincronización en la nube', active: true),
                        const FeatureRow(text: 'Perfiles familiares', active: true),
                        const FeatureRow(text: 'Alertas de cambios', active: true),
                      ],
                    ),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 24),

            // Action Buttons
            if (!isPremium) ...[
              SizedBox(
                width: double.infinity,
                child: ElevatedButton.icon(
                  icon: const Icon(Icons.qr_code_scanner_rounded),
                  label: const Text('Simular Pago QR / Suscribir (Bs. 20/mes)'),
                  style: ElevatedButton.styleFrom(
                    backgroundColor: AppTheme.secondaryTeal,
                    padding: const EdgeInsets.symmetric(vertical: 16),
                  ),
                  onPressed: _openQrModal,
                ),
              ),
              const SizedBox(height: 10),
              TextButton(
                onPressed: () => Navigator.pop(context),
                child: const Text('Continuar con el Plan Gratuito', style: TextStyle(color: Colors.grey)),
              ),
            ] else ...[
              Container(
                padding: const EdgeInsets.all(16),
                decoration: BoxDecoration(
                  color: AppTheme.successGreen.withAlpha(20),
                  borderRadius: BorderRadius.circular(16),
                  border: Border.all(color: AppTheme.successGreen.withAlpha(60)),
                ),
                child: Row(
                  children: const [
                    Icon(Icons.stars_rounded, color: AppTheme.successGreen, size: 30),
                    SizedBox(width: 12),
                    Expanded(
                      child: Text(
                        '¡Suscripción Premium Activa! Sincronización en la nube y alertas habilitadas.',
                        style: TextStyle(fontWeight: FontWeight.bold, fontSize: 13, color: AppTheme.successGreen),
                      ),
                    ),
                  ],
                ),
              ),
              const SizedBox(height: 14),
              OutlinedButton(
                onPressed: _downgradePlan,
                child: const Text('Volver al Plan Gratuito'),
              ),
            ],
          ],
        ),
      ),
    );
  }
}

class FeatureRow extends StatelessWidget {
  final String text;
  final bool active;

  const FeatureRow({super.key, required this.text, required this.active});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: Row(
        children: [
          Icon(
            active ? Icons.check_rounded : Icons.close_rounded,
            size: 14,
            color: active ? AppTheme.successGreen : Colors.grey,
          ),
          const SizedBox(width: 6),
          Expanded(
            child: Text(
              text,
              style: TextStyle(
                fontSize: 11.5,
                color: active ? Theme.of(context).colorScheme.onSurface : Colors.grey,
                decoration: active ? null : TextDecoration.lineThrough,
              ),
            ),
          ),
        ],
      ),
    );
  }
}

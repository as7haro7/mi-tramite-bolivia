import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../config/theme.dart';
import '../providers/app_provider.dart';
import '../services/storage_service.dart';
import '../widgets/disclaimer_banner.dart';

class ProfileScreen extends StatelessWidget {
  final VoidCallback onOpenPremium;

  const ProfileScreen({super.key, required this.onOpenPremium});

  final List<String> _cities = const [
    'La Paz',
    'El Alto',
    'Cochabamba',
    'Santa Cruz de la Sierra',
    'Oruro',
    'Sucre',
    'Potosí',
    'Tarija',
    'Trinidad',
    'Cobija',
    'Prefiero no indicar',
  ];

  void _showFamilyProfilesModal(BuildContext context, bool isPremium) {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(20)),
        title: Row(
          children: const [
            Icon(Icons.family_restroom_rounded, color: AppTheme.primaryBlue),
            SizedBox(width: 8),
            Text('Perfiles Familiares'),
          ],
        ),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              'Organiza los documentos y trámites de tus familiares (hijos, padres, cónyuge) en un solo lugar.',
              style: TextStyle(fontSize: 13.5),
            ),
            const SizedBox(height: 12),
            Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: isPremium ? AppTheme.successGreen.withAlpha(20) : Colors.amber.withAlpha(20),
                borderRadius: BorderRadius.circular(10),
              ),
              child: Row(
                children: [
                  Icon(
                    isPremium ? Icons.check_circle_rounded : Icons.lock_outline_rounded,
                    color: isPremium ? AppTheme.successGreen : Colors.amber.shade900,
                  ),
                  const SizedBox(width: 8),
                  Expanded(
                    child: Text(
                      isPremium
                          ? 'Función activa en tu Plan Premium.'
                          : 'Función disponible exclusivamente en el Plan Premium (Bs. 20/mes).',
                      style: TextStyle(
                        fontSize: 12,
                        fontWeight: FontWeight.bold,
                        color: isPremium ? AppTheme.successGreen : Colors.amber.shade900,
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
        actions: [
          if (!isPremium)
            ElevatedButton(
              onPressed: () {
                Navigator.pop(ctx);
                onOpenPremium();
              },
              child: const Text('Ver Plan Premium'),
            ),
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: const Text('Cerrar'),
          ),
        ],
      ),
    );
  }

  void _showAlertsModal(BuildContext context, bool isPremium) {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(20)),
        title: Row(
          children: const [
            Icon(Icons.notifications_active_rounded, color: AppTheme.secondaryTeal),
            SizedBox(width: 8),
            Text('Alertas de Cambios'),
          ],
        ),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              'Recibe notificaciones automáticas cuando una institución modifique requisitos, costos, cuentas bancarias u horarios.',
              style: TextStyle(fontSize: 13.5),
            ),
            const SizedBox(height: 12),
            Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: isPremium ? AppTheme.successGreen.withAlpha(20) : Colors.amber.withAlpha(20),
                borderRadius: BorderRadius.circular(10),
              ),
              child: Row(
                children: [
                  Icon(
                    isPremium ? Icons.check_circle_rounded : Icons.lock_outline_rounded,
                    color: isPremium ? AppTheme.successGreen : Colors.amber.shade900,
                  ),
                  const SizedBox(width: 8),
                  Expanded(
                    child: Text(
                      isPremium
                          ? 'Alertas activas para tus trámites guardados.'
                          : 'Función disponible en el Plan Premium (Bs. 20/mes).',
                      style: TextStyle(
                        fontSize: 12,
                        fontWeight: FontWeight.bold,
                        color: isPremium ? AppTheme.successGreen : Colors.amber.shade900,
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
        actions: [
          if (!isPremium)
            ElevatedButton(
              onPressed: () {
                Navigator.pop(ctx);
                onOpenPremium();
              },
              child: const Text('Ver Plan Premium'),
            ),
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: const Text('Cerrar'),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final appProvider = Provider.of<AppProvider>(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Perfil y Ajustes'),
      ),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          // User Card
          Container(
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              color: Theme.of(context).cardTheme.color,
              borderRadius: BorderRadius.circular(16),
              border: Border.all(color: AppTheme.lightBorder),
            ),
            child: Row(
              children: [
                CircleAvatar(
                  radius: 26,
                  backgroundColor: appProvider.isPremium ? Colors.amber : AppTheme.primaryBlue,
                  child: Icon(
                    appProvider.isPremium ? Icons.stars_rounded : Icons.person_rounded,
                    color: Colors.white,
                    size: 30,
                  ),
                ),
                const SizedBox(width: 14),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        children: [
                          Text(
                            appProvider.isPremium ? 'Usuario Premium' : 'Usuario Gratuito',
                            style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
                          ),
                          const SizedBox(width: 6),
                          Container(
                            padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                            decoration: BoxDecoration(
                              color: appProvider.isPremium ? Colors.amber.withAlpha(40) : AppTheme.primaryBlue.withAlpha(20),
                              borderRadius: BorderRadius.circular(10),
                            ),
                            child: Text(
                              appProvider.isPremium ? 'Bs. 20/mes' : 'Demo',
                              style: TextStyle(
                                fontSize: 10.5,
                                fontWeight: FontWeight.bold,
                                color: appProvider.isPremium ? Colors.amber.shade900 : AppTheme.primaryBlue,
                              ),
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 2),
                      const Text(
                        'Orientación ciudadana anónima (sin registro obligatorio)',
                        style: TextStyle(fontSize: 12, color: Colors.grey),
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(height: 16),

          // Subscription Upgrade Banner
          Card(
            color: appProvider.isPremium ? AppTheme.successGreen.withAlpha(15) : AppTheme.primaryBlue.withAlpha(15),
            child: ListTile(
              leading: Icon(
                appProvider.isPremium ? Icons.workspace_premium_rounded : Icons.star_rounded,
                color: appProvider.isPremium ? AppTheme.successGreen : AppTheme.primaryBlue,
              ),
              title: Text(
                appProvider.isPremium ? 'Suscripción Premium Activa' : 'Mejorar a Plan Premium',
                style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 14),
              ),
              subtitle: Text(
                appProvider.isPremium
                    ? 'Sincronización en la nube y perfiles familiares habilitados.'
                    : 'Prueba la experiencia Premium por Bs. 20/mes (Demostración comercial).',
                style: const TextStyle(fontSize: 12),
              ),
              trailing: const Icon(Icons.arrow_forward_ios_rounded, size: 16),
              onTap: onOpenPremium,
            ),
          ),
          const SizedBox(height: 20),

          // Extra Features Shortcuts
          const Text('Funciones Avanzadas', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 13, color: Colors.grey)),
          const SizedBox(height: 8),
          ListTile(
            leading: const Icon(Icons.family_restroom_rounded, color: AppTheme.primaryBlue),
            title: const Text('Perfiles Familiares'),
            subtitle: const Text('Gestión de trámites para dependientes', style: TextStyle(fontSize: 12)),
            trailing: const Icon(Icons.chevron_right_rounded),
            onTap: () => _showFamilyProfilesModal(context, appProvider.isPremium),
          ),
          ListTile(
            leading: const Icon(Icons.notifications_active_rounded, color: AppTheme.secondaryTeal),
            title: const Text('Alertas de Cambios Institucionales'),
            subtitle: const Text('Notificaciones sobre nuevos requisitos o costos', style: TextStyle(fontSize: 12)),
            trailing: const Icon(Icons.chevron_right_rounded),
            onTap: () => _showAlertsModal(context, appProvider.isPremium),
          ),
          const SizedBox(height: 16),

          // Territory Preference
          const Text('Preferencia Territorial', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 13, color: Colors.grey)),
          const SizedBox(height: 8),
          Container(
            decoration: BoxDecoration(
              color: Theme.of(context).cardTheme.color,
              borderRadius: BorderRadius.circular(12),
              border: Border.all(color: AppTheme.lightBorder),
            ),
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
            child: DropdownButtonHideUnderline(
              child: DropdownButton<String>(
                value: appProvider.selectedCity,
                isExpanded: true,
                items: _cities.map((c) => DropdownMenuItem(value: c, child: Text('Municipio: $c'))).toList(),
                onChanged: (val) {
                  if (val != null) appProvider.setCity(val);
                },
              ),
            ),
          ),
          const SizedBox(height: 20),

          // Transparency & Privacy
          const Text('Transparencia y Privacidad', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 13, color: Colors.grey)),
          const SizedBox(height: 8),
          const DisclaimerBanner(),
          const SizedBox(height: 16),

          // Actions
          ListTile(
            leading: const Icon(Icons.cleaning_services_rounded, color: AppTheme.alertRed),
            title: const Text('Borrar historial y datos locales'),
            subtitle: const Text('Elimina la caché de búsquedas y checklists locales', style: TextStyle(fontSize: 12)),
            onTap: () async {
              final confirm = await showDialog<bool>(
                context: context,
                builder: (ctx) => AlertDialog(
                  title: const Text('¿Borrar Datos Locales?'),
                  content: const Text('Se limpiarán tus búsquedas recientes y checklists almacenados.'),
                  actions: [
                    TextButton(onPressed: () => Navigator.pop(ctx, false), child: const Text('Cancelar')),
                    TextButton(
                      onPressed: () => Navigator.pop(ctx, true),
                      child: const Text('Borrar Todo', style: TextStyle(color: AppTheme.alertRed)),
                    ),
                  ],
                ),
              );
              if (confirm == true) {
                await StorageService.clearAll();
                if (context.mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(content: Text('Datos locales borrados correctamente.')),
                  );
                }
              }
            },
          ),
          const Divider(),
          const Padding(
            padding: EdgeInsets.symmetric(vertical: 12),
            child: Center(
              child: Text(
                'Mi Trámite Bolivia v1.0 (Demo MVP)\nDesarrollado para UMSA Emprendimiento',
                textAlign: TextAlign.center,
                style: TextStyle(fontSize: 12, color: Colors.grey),
              ),
            ),
          ),
        ],
      ),
    );
  }
}

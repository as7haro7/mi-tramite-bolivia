import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../config/theme.dart';
import '../models/checklist_item.dart';
import '../providers/checklist_provider.dart';
import '../widgets/requirement_tile.dart';

class ChecklistScreen extends StatelessWidget {
  final TramiteChecklist checklist;

  const ChecklistScreen({super.key, required this.checklist});

  @override
  Widget build(BuildContext context) {
    final provider = Provider.of<ChecklistProvider>(context);
    final currentChecklist = provider.getChecklistBySlug(checklist.tramiteSlug) ?? checklist;

    return Scaffold(
      appBar: AppBar(
        title: Text(currentChecklist.tramiteTitulo),
        actions: [
          IconButton(
            icon: const Icon(Icons.delete_outline_rounded, color: AppTheme.alertRed),
            tooltip: 'Eliminar Checklist',
            onPressed: () async {
              final confirm = await showDialog<bool>(
                context: context,
                builder: (ctx) => AlertDialog(
                  title: const Text('¿Eliminar Checklist?'),
                  content: const Text('Esta acción eliminará el progreso guardado localmente.'),
                  actions: [
                    TextButton(onPressed: () => Navigator.pop(ctx, false), child: const Text('Cancelar')),
                    TextButton(
                      onPressed: () => Navigator.pop(ctx, true),
                      child: const Text('Eliminar', style: TextStyle(color: AppTheme.alertRed)),
                    ),
                  ],
                ),
              );
              if (confirm == true) {
                provider.deleteChecklist(currentChecklist.tramiteSlug);
                if (context.mounted) Navigator.pop(context);
              }
            },
          ),
        ],
      ),
      body: Column(
        children: [
          // Progress Header
          Container(
            padding: const EdgeInsets.all(16),
            color: Theme.of(context).cardTheme.color,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(
                      'Progreso de preparación',
                      style: Theme.of(context).textTheme.titleMedium,
                    ),
                    Text(
                      '${currentChecklist.completedCount} de ${currentChecklist.totalCount} (${(currentChecklist.progress * 100).toInt()}%)',
                      style: TextStyle(
                        fontWeight: FontWeight.bold,
                        color: currentChecklist.isFullyCompleted
                            ? AppTheme.successGreen
                            : AppTheme.primaryBlue,
                        fontSize: 14,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 10),
                ClipRRect(
                  borderRadius: BorderRadius.circular(6),
                  child: LinearProgressIndicator(
                    value: currentChecklist.progress,
                    minHeight: 10,
                    backgroundColor: Colors.grey.shade200,
                    color: currentChecklist.isFullyCompleted
                        ? AppTheme.successGreen
                        : AppTheme.primaryBlue,
                  ),
                ),
                if (currentChecklist.isFullyCompleted) ...[
                  const SizedBox(height: 10),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                    decoration: BoxDecoration(
                      color: AppTheme.successGreen.withAlpha(20),
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Row(
                      children: const [
                        Icon(Icons.check_circle_rounded, color: AppTheme.successGreen, size: 18),
                        SizedBox(width: 8),
                        Text(
                          '¡Felicitaciones! Tienes todos los requisitos listos.',
                          style: TextStyle(color: AppTheme.successGreen, fontWeight: FontWeight.bold, fontSize: 12.5),
                        ),
                      ],
                    ),
                  ),
                ],
              ],
            ),
          ),

          // Items List
          Expanded(
            child: ListView.builder(
              padding: const EdgeInsets.all(16),
              itemCount: currentChecklist.items.length,
              itemBuilder: (context, index) {
                final item = currentChecklist.items[index];
                return RequirementTile(
                  item: item,
                  onChanged: (val) {
                    provider.toggleItem(currentChecklist.tramiteSlug, item.id);
                  },
                  onNoteChanged: (note) {
                    provider.updateNote(currentChecklist.tramiteSlug, item.id, note);
                  },
                );
              },
            ),
          ),
        ],
      ),
    );
  }
}

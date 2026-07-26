import 'package:flutter/material.dart';
import '../config/theme.dart';
import '../models/checklist_item.dart';

class RequirementTile extends StatefulWidget {
  final RequisitoCheckable item;
  final ValueChanged<bool?> onChanged;
  final ValueChanged<String> onNoteChanged;

  const RequirementTile({
    super.key,
    required this.item,
    required this.onChanged,
    required this.onNoteChanged,
  });

  @override
  State<RequirementTile> createState() => _RequirementTileState();
}

class _RequirementTileState extends State<RequirementTile> {
  bool _showNoteInput = false;
  late TextEditingController _noteController;

  @override
  void initState() {
    super.initState();
    _noteController = TextEditingController(text: widget.item.personalNote);
  }

  @override
  void dispose() {
    _noteController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.only(bottom: 10),
      decoration: BoxDecoration(
        color: widget.item.isCompleted
            ? AppTheme.successGreen.withAlpha(12)
            : Theme.of(context).cardTheme.color,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: widget.item.isCompleted
              ? AppTheme.successGreen.withAlpha(80)
              : AppTheme.lightBorder,
        ),
      ),
      child: Column(
        children: [
          CheckboxListTile(
            value: widget.item.isCompleted,
            activeColor: AppTheme.successGreen,
            shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(4)),
            title: Text(
              widget.item.nombre,
              style: TextStyle(
                fontWeight: FontWeight.w600,
                fontSize: 14.5,
                decoration: widget.item.isCompleted ? TextDecoration.lineThrough : null,
                color: widget.item.isCompleted ? Colors.grey : null,
              ),
            ),
            subtitle: widget.item.descripcion != null && widget.item.descripcion!.isNotEmpty
                ? Text(
                    widget.item.descripcion!,
                    style: const TextStyle(fontSize: 12.5),
                  )
                : null,
            onChanged: widget.onChanged,
          ),
          if (widget.item.personalNote.isNotEmpty || _showNoteInput) ...[
            Padding(
              padding: const EdgeInsets.only(left: 16, right: 16, bottom: 12),
              child: TextField(
                controller: _noteController,
                style: const TextStyle(fontSize: 13),
                decoration: InputDecoration(
                  hintText: 'Agregar nota personal (ej. fotocopia lista)...',
                  isDense: true,
                  filled: true,
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(8),
                    borderSide: BorderSide.none,
                  ),
                ),
                onSubmitted: widget.onNoteChanged,
                onChanged: widget.onNoteChanged,
              ),
            ),
          ] else ...[
            Align(
              alignment: Alignment.centerLeft,
              child: TextButton.icon(
                icon: const Icon(Icons.note_add_outlined, size: 16),
                label: const Text('Nota personal', style: TextStyle(fontSize: 12)),
                style: TextButton.styleFrom(
                  padding: const EdgeInsets.only(left: 16, bottom: 8),
                  minimumSize: Size.zero,
                  tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                ),
                onPressed: () {
                  setState(() {
                    _showNoteInput = true;
                  });
                },
              ),
            ),
          ],
        ],
      ),
    );
  }
}

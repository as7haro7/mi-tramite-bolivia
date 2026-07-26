import 'package:flutter/material.dart';
import '../models/checklist_item.dart';
import '../models/tramite.dart';
import '../services/storage_service.dart';

class ChecklistProvider with ChangeNotifier {
  List<TramiteChecklist> _checklists = [];

  List<TramiteChecklist> get checklists => _checklists;

  ChecklistProvider() {
    _loadChecklists();
  }

  Future<void> _loadChecklists() async {
    _checklists = await StorageService.getChecklists();
    notifyListeners();
  }

  TramiteChecklist? getChecklistBySlug(String slug) {
    try {
      return _checklists.firstWhere((c) => c.tramiteSlug == slug);
    } catch (_) {
      return null;
    }
  }

  Future<void> createOrUpdateFromTramite(Tramite tramite) async {
    final existingIndex = _checklists.indexWhere((c) => c.tramiteSlug == tramite.slug);
    
    final items = tramite.requisitos.map((r) => RequisitoCheckable(
      id: r.id,
      nombre: r.nombre,
      descripcion: r.descripcion,
    )).toList();

    final newChecklist = TramiteChecklist(
      tramiteSlug: tramite.slug,
      tramiteTitulo: tramite.titulo,
      institucionNombre: tramite.institucion?.sigla ?? tramite.institucion?.nombre ?? 'Institución',
      createdAt: DateTime.now(),
      items: items,
    );

    if (existingIndex >= 0) {
      _checklists[existingIndex] = newChecklist;
    } else {
      _checklists.insert(0, newChecklist);
    }

    await StorageService.saveChecklists(_checklists);
    notifyListeners();
  }

  Future<void> toggleItem(String slug, int reqId) async {
    final checklist = getChecklistBySlug(slug);
    if (checklist == null) return;

    try {
      final item = checklist.items.firstWhere((i) => i.id == reqId);
      item.isCompleted = !item.isCompleted;
      await StorageService.saveChecklists(_checklists);
      notifyListeners();
    } catch (_) {}
  }

  Future<void> updateNote(String slug, int reqId, String note) async {
    final checklist = getChecklistBySlug(slug);
    if (checklist == null) return;

    try {
      final item = checklist.items.firstWhere((i) => i.id == reqId);
      item.personalNote = note;
      await StorageService.saveChecklists(_checklists);
      notifyListeners();
    } catch (_) {}
  }

  Future<void> deleteChecklist(String slug) async {
    _checklists.removeWhere((c) => c.tramiteSlug == slug);
    await StorageService.saveChecklists(_checklists);
    notifyListeners();
  }
}

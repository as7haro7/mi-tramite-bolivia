import 'package:flutter/material.dart';
import '../config/theme.dart';
import '../services/api_service.dart';

class ReportFormScreen extends StatefulWidget {
  final String tramiteSlug;
  final String tramiteTitulo;

  const ReportFormScreen({
    super.key,
    required this.tramiteSlug,
    required this.tramiteTitulo,
  });

  @override
  State<ReportFormScreen> createState() => _ReportFormScreenState();
}

class _ReportFormScreenState extends State<ReportFormScreen> {
  final _formKey = GlobalKey<FormState>();
  String _tipoReporte = 'costo_distinto';
  final TextEditingController _descController = TextEditingController();
  bool _esAnonimo = true;
  bool _isSubmitting = false;

  final Map<String, String> _tiposMap = {
    'costo_distinto': 'Costo cobrado es diferente al publicado',
    'requisito_omitido': 'Me pidieron un requisito no listado',
    'oficina_cerrada': 'La oficina o dirección está desactualizada',
    'horario_incorrecto': 'El horario de atención cambió',
    'otro': 'Otro error de información',
  };

  void _submitReport() async {
    if (!_formKey.currentState!.validate() || _isSubmitting) return;

    setState(() => _isSubmitting = true);
    final success = await ApiService.enviarReporte(
      tramiteSlug: widget.tramiteSlug,
      tipoReporte: _tipoReporte,
      descripcion: _descController.text,
      esAnonimo: _esAnonimo,
    );
    setState(() => _isSubmitting = false);

    if (!mounted) return;
    if (success) {
      showDialog(
        context: context,
        builder: (ctx) => AlertDialog(
          title: const Text('¡Reporte Recibido!'),
          content: const Text(
            'Gracias por ayudar a mantener la información actualizada. El equipo editorial verificará el reporte con la fuente oficial.',
          ),
          actions: [
            TextButton(
              onPressed: () {
                Navigator.pop(ctx);
                Navigator.pop(context);
              },
              child: const Text('Entendido'),
            ),
          ],
        ),
      );
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('No se pudo enviar el reporte. Inténtalo más tarde.')),
      );
    }
  }

  @override
  void dispose() {
    _descController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Reportar Cambio o Error'),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(20),
        child: Form(
          key: _formKey,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                'Reportar para: ${widget.tramiteTitulo}',
                style: Theme.of(context).textTheme.titleMedium,
              ),
              const SizedBox(height: 12),

              // Non-sensitive PII Warning
              Container(
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: AppTheme.alertRed.withAlpha(15),
                  borderRadius: BorderRadius.circular(10),
                  border: Border.all(color: AppTheme.alertRed.withAlpha(50)),
                ),
                child: Row(
                  children: const [
                    Icon(Icons.shield_outlined, color: AppTheme.alertRed, size: 20),
                    SizedBox(width: 10),
                    Expanded(
                      child: Text(
                        'Por tu seguridad, NO adjuntes números de carnet, fotos de documentos ni datos personales.',
                        style: TextStyle(fontSize: 12, fontWeight: FontWeight.bold, color: AppTheme.alertRed),
                      ),
                    ),
                  ],
                ),
              ),
              const SizedBox(height: 20),

              const Text('Tipo de Observación:', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 13)),
              const SizedBox(height: 6),
              DropdownButtonFormField<String>(
                initialValue: _tipoReporte,
                decoration: InputDecoration(
                  border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
                  contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
                ),
                items: _tiposMap.entries.map((e) {
                  return DropdownMenuItem(value: e.key, child: Text(e.value));
                }).toList(),
                onChanged: (val) {
                  if (val != null) setState(() => _tipoReporte = val);
                },
              ),
              const SizedBox(height: 16),

              const Text('Detalle de lo ocurrido:', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 13)),
              const SizedBox(height: 6),
              TextFormField(
                controller: _descController,
                maxLines: 4,
                decoration: InputDecoration(
                  hintText: 'Ej. En la oficina del centro cobraron 15 Bs por la fotocopia legalizada...',
                  border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
                ),
                validator: (val) {
                  if (val == null || val.trim().length < 10) {
                    return 'Describe brevemente la observación (mínimo 10 caracteres).';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 16),

              SwitchListTile(
                contentPadding: EdgeInsets.zero,
                title: const Text('Enviar reporte de forma anónima'),
                value: _esAnonimo,
                activeThumbColor: AppTheme.primaryBlue,
                onChanged: (val) => setState(() => _esAnonimo = val),
              ),
              const SizedBox(height: 24),

              SizedBox(
                width: double.infinity,
                child: ElevatedButton.icon(
                  icon: _isSubmitting
                      ? const SizedBox(width: 16, height: 16, child: CircularProgressIndicator(color: Colors.white, strokeWidth: 2))
                      : const Icon(Icons.send_rounded),
                  label: Text(_isSubmitting ? 'Enviando...' : 'Enviar Reporte Ciudadano'),
                  onPressed: _isSubmitting ? null : _submitReport,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

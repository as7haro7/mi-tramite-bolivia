import 'package:flutter/material.dart';
import '../config/theme.dart';

class VerificationSeal extends StatelessWidget {
  final String? date;
  final bool isVerified;

  const VerificationSeal({
    super.key,
    this.date,
    this.isVerified = true,
  });

  @override
  Widget build(BuildContext context) {
    final color = isVerified ? AppTheme.successGreen : AppTheme.accentAmber;
    final text = isVerified
        ? (date != null ? 'Verificado: $date' : 'Información Verificada')
        : 'Revisión Pendiente';

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
      decoration: BoxDecoration(
        color: color.withAlpha(25),
        borderRadius: BorderRadius.circular(20),
        border: Border.all(color: color.withAlpha(80)),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(
            isVerified ? Icons.check_circle_rounded : Icons.history_toggle_off_rounded,
            color: color,
            size: 14,
          ),
          const SizedBox(width: 5),
          Text(
            text,
            style: TextStyle(
              color: color,
              fontSize: 11,
              fontWeight: FontWeight.bold,
            ),
          ),
        ],
      ),
    );
  }
}

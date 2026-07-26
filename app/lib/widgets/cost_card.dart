import 'package:flutter/material.dart';
import '../config/theme.dart';
import '../models/tramite.dart';

class CostCard extends StatelessWidget {
  final Costo costo;

  const CostCard({super.key, required this.costo});

  @override
  Widget build(BuildContext context) {
    final isGratis = costo.monto == 0;

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
      margin: const EdgeInsets.only(bottom: 8),
      decoration: BoxDecoration(
        color: isGratis ? AppTheme.successGreen.withAlpha(15) : AppTheme.primaryBlue.withAlpha(15),
        borderRadius: BorderRadius.circular(10),
        border: Border.all(
          color: isGratis ? AppTheme.successGreen.withAlpha(60) : AppTheme.primaryBlue.withAlpha(60),
        ),
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Expanded(
            child: Text(
              costo.concepto,
              style: const TextStyle(fontWeight: FontWeight.w600, fontSize: 13.5),
            ),
          ),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
            decoration: BoxDecoration(
              color: isGratis ? AppTheme.successGreen : AppTheme.primaryBlue,
              borderRadius: BorderRadius.circular(20),
            ),
            child: Text(
              isGratis ? 'GRATUITO' : '${costo.moneda} ${costo.monto.toStringAsFixed(2)}',
              style: const TextStyle(
                color: Colors.white,
                fontWeight: FontWeight.bold,
                fontSize: 12.5,
              ),
            ),
          ),
        ],
      ),
    );
  }
}

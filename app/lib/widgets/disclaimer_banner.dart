import 'package:flutter/material.dart';
import '../config/theme.dart';

class DisclaimerBanner extends StatelessWidget {
  final bool isCompact;

  const DisclaimerBanner({super.key, this.isCompact = false});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: EdgeInsets.symmetric(
        horizontal: 14,
        vertical: isCompact ? 8 : 12,
      ),
      decoration: BoxDecoration(
        color: AppTheme.primaryBlue.withAlpha(20),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: AppTheme.primaryBlue.withAlpha(50)),
      ),
      child: Row(
        children: [
          const Icon(
            Icons.verified_user_outlined,
            color: AppTheme.primaryBlue,
            size: 22,
          ),
          const SizedBox(width: 10),
          Expanded(
            child: Text(
              'Plataforma de orientación independiente de instituciones públicas. No ejecuta trámites.',
              style: TextStyle(
                color: Theme.of(context).colorScheme.onSurface,
                fontSize: isCompact ? 11 : 12.5,
                height: 1.3,
                fontWeight: FontWeight.w500,
              ),
            ),
          ),
        ],
      ),
    );
  }
}

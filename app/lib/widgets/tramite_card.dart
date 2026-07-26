import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';
import '../config/theme.dart';
import '../models/tramite.dart';
import '../providers/favorites_provider.dart';
import 'verification_seal.dart';

class TramiteCard extends StatelessWidget {
  final Tramite tramite;
  final VoidCallback onTap;

  const TramiteCard({
    super.key,
    required this.tramite,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final favProvider = Provider.of<FavoritesProvider>(context);
    final isFav = favProvider.isFavorite(tramite.slug);
    final instSigla = tramite.institucion?.sigla ?? tramite.institucion?.nombre ?? 'Institución';
    final isDark = Theme.of(context).brightness == Brightness.dark;

    return Container(
      margin: const EdgeInsets.only(bottom: 10),
      child: Material(
        color: isDark ? AppTheme.iosDarkCard : AppTheme.iosLightCard,
        borderRadius: BorderRadius.circular(16),
        child: InkWell(
          borderRadius: BorderRadius.circular(16),
          onTap: () {
            HapticFeedback.lightImpact();
            onTap();
          },
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  crossAxisAlignment: CrossAxisAlignment.center,
                  children: [
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                      decoration: BoxDecoration(
                        color: AppTheme.tintColor.withAlpha(isDark ? 30 : 15),
                        borderRadius: BorderRadius.circular(8),
                      ),
                      child: Text(
                        instSigla,
                        style: const TextStyle(
                          color: AppTheme.tintColor,
                          fontWeight: FontWeight.w600,
                          fontSize: 12,
                        ),
                      ),
                    ),
                    const Spacer(),
                    CupertinoButton(
                      padding: EdgeInsets.zero,
                      minimumSize: Size.zero,
                      onPressed: () {
                        HapticFeedback.selectionClick();
                        favProvider.toggleFavorite(tramite.slug);
                      },
                      child: Icon(
                        isFav ? CupertinoIcons.star_fill : CupertinoIcons.star,
                        color: isFav ? AppTheme.systemOrange : AppTheme.iosLightSubtext,
                        size: 22,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 10),
                Text(
                  tramite.titulo,
                  style: TextStyle(
                    fontWeight: FontWeight.bold,
                    fontSize: 17,
                    letterSpacing: -0.4,
                    color: isDark ? AppTheme.iosDarkText : AppTheme.iosLightText,
                  ),
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                ),
                if (tramite.resumen != null && tramite.resumen!.isNotEmpty) ...[
                  const SizedBox(height: 4),
                  Text(
                    tramite.resumen!,
                    style: const TextStyle(
                      fontSize: 14,
                      color: AppTheme.iosLightSubtext,
                      height: 1.3,
                    ),
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                  ),
                ],
                const SizedBox(height: 12),
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    VerificationSeal(date: tramite.verificadoEn, isVerified: tramite.verificado),
                    const Row(
                      children: [
                        Text(
                          'Ficha',
                          style: TextStyle(
                            color: AppTheme.tintColor,
                            fontWeight: FontWeight.w600,
                            fontSize: 14,
                          ),
                        ),
                        SizedBox(width: 2),
                        Icon(
                          CupertinoIcons.chevron_right,
                          color: AppTheme.tintColor,
                          size: 16,
                        ),
                      ],
                    ),
                  ],
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}

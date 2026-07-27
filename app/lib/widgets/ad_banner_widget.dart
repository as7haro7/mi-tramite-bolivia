import 'package:flutter/material.dart';
import 'package:google_mobile_ads/google_mobile_ads.dart';
import '../config/theme.dart';

class AdBannerWidget extends StatefulWidget {
  final String adUnitId;
  final AdSize adSize;
  final EdgeInsets margin;

  const AdBannerWidget({
    super.key,
    required this.adUnitId,
    this.adSize = AdSize.banner,
    this.margin = const EdgeInsets.symmetric(vertical: 16, horizontal: 16),
  });

  @override
  State<AdBannerWidget> createState() => _AdBannerWidgetState();
}

class _AdBannerWidgetState extends State<AdBannerWidget> {
  BannerAd? _bannerAd;
  bool _isAdLoaded = false;
  bool _hasError = false;

  @override
  void initState() {
    super.initState();
    _loadAd();
  }

  void _loadAd() {
    _bannerAd = BannerAd(
      adUnitId: widget.adUnitId,
      size: widget.adSize,
      request: const AdRequest(),
      listener: BannerAdListener(
        onAdLoaded: (ad) {
          if (mounted) {
            setState(() {
              _isAdLoaded = true;
              _hasError = false;
            });
          }
        },
        onAdFailedToLoad: (ad, error) {
          debugPrint('Ad failed to load: $error');
          ad.dispose();
          if (mounted) {
            setState(() {
              _hasError = true;
              _isAdLoaded = false;
            });
          }
        },
      ),
    );

    _bannerAd?.load();
  }

  @override
  void dispose() {
    _bannerAd?.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    // Si hay error, no mostrar nada
    if (_hasError) {
      return const SizedBox.shrink();
    }

    // Si está cargando, mostrar placeholder con shimmer effect
    if (!_isAdLoaded) {
      return Container(
        margin: widget.margin,
        height: widget.adSize.height.toDouble(),
        decoration: BoxDecoration(
          color: AppTheme.iosLightSearchBg,
          borderRadius: BorderRadius.circular(12),
          border: Border.all(
            color: AppTheme.iosLightBorder.withOpacity(0.3),
            width: 1,
          ),
        ),
        child: Center(
          child: Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              SizedBox(
                width: 16,
                height: 16,
                child: CircularProgressIndicator(
                  strokeWidth: 2,
                  valueColor: AlwaysStoppedAnimation<Color>(
                    AppTheme.iosLightSubtext,
                  ),
                ),
              ),
              const SizedBox(width: 12),
              Text(
                'Cargando anuncio...',
                style: TextStyle(
                  color: AppTheme.iosLightSubtext,
                  fontSize: 12,
                  fontWeight: FontWeight.w500,
                ),
              ),
            ],
          ),
        ),
      );
    }

    // Mostrar el anuncio cargado
    return Container(
      margin: widget.margin,
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: AppTheme.iosLightBorder.withOpacity(0.3),
          width: 1,
        ),
      ),
      clipBehavior: Clip.antiAlias,
      child: SizedBox(
        width: _bannerAd!.size.width.toDouble(),
        height: _bannerAd!.size.height.toDouble(),
        child: AdWidget(ad: _bannerAd!),
      ),
    );
  }
}

/// Widget helper para mostrar un anuncio con una etiqueta de "Publicidad"
class LabeledAdBanner extends StatelessWidget {
  final String adUnitId;
  final AdSize adSize;
  final EdgeInsets margin;

  const LabeledAdBanner({
    super.key,
    required this.adUnitId,
    this.adSize = AdSize.banner,
    this.margin = const EdgeInsets.symmetric(vertical: 16, horizontal: 16),
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.center,
      children: [
        Padding(
          padding: const EdgeInsets.only(bottom: 4),
          child: Text(
            'Publicidad',
            style: TextStyle(
              fontSize: 10,
              color: AppTheme.iosLightSubtext.withOpacity(0.6),
              fontWeight: FontWeight.w500,
              letterSpacing: 0.5,
            ),
          ),
        ),
        AdBannerWidget(
          adUnitId: adUnitId,
          adSize: adSize,
          margin: margin,
        ),
      ],
    );
  }
}
